package version

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
	kotsv1beta1 "github.com/replicatedhq/kots/kotskinds/apis/kots/v1beta1"
	"github.com/replicatedhq/kots/pkg/api/version/types"
	"github.com/replicatedhq/kots/pkg/gitops"
	"github.com/replicatedhq/kots/pkg/k8sutil"
	"github.com/replicatedhq/kots/pkg/kotsutil"
	"github.com/replicatedhq/kots/pkg/logger"
	"github.com/replicatedhq/kots/pkg/persistence"
	"github.com/replicatedhq/kots/pkg/store"
	storetypes "github.com/replicatedhq/kots/pkg/store/types"
	"github.com/replicatedhq/kots/pkg/util"
	troubleshootpreflight "github.com/replicatedhq/troubleshoot/pkg/preflight"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes/scheme"
	applicationv1beta1 "sigs.k8s.io/application/api/v1beta1"
)

type DownstreamGitOps struct {
}

func (d *DownstreamGitOps) CreateGitOpsDownstreamCommit(appID string, clusterID string, newSequence int, filesInDir string, downstreamName string) (string, error) {
	downstreamGitOps, err := gitops.GetDownstreamGitOps(appID, clusterID)
	if err != nil {
		return "", errors.Wrap(err, "failed to get downstream gitops")
	}
	if downstreamGitOps == nil {
		return "", nil
	}

	a, err := store.GetStore().GetApp(appID)
	if err != nil {
		return "", errors.Wrap(err, "failed to get app")
	}
	createdCommitURL, err := gitops.CreateGitOpsCommit(downstreamGitOps, a.Slug, a.Name, int(newSequence), filesInDir, downstreamName)
	if err != nil {
		return "", errors.Wrap(err, "failed to create gitops commit")
	}

	return createdCommitURL, nil
}

// DeployVersion deploys the version for the given sequence
func DeployVersion(appID string, sequence int64) error {

	blocked, err := isBlockedDueToStricPreFlights(appID, sequence)
	if err != nil {
		return errors.Wrap(err, "failed to validate strict preflights")
	}
	if blocked {
		return util.ActionableError{
			NoRetry: true,
			Message: "Unable to deploy as preflight check's strict analyzer has failed",
		}
	}

	db := persistence.MustGetDBSession()

	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "failed to begin")
	}
	defer tx.Rollback()

	query := `update app_downstream set current_sequence = $1 where app_id = $2`
	_, err = tx.Exec(query, sequence, appID)
	if err != nil {
		return errors.Wrap(err, "failed to update app downstream current sequence")
	}

	query = `update app_downstream_version set status = 'deployed', applied_at = $3 where sequence = $1 and app_id = $2`
	_, err = tx.Exec(query, sequence, appID, time.Now())
	if err != nil {
		return errors.Wrap(err, "failed to update app downstream version status")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit")
	}

	return nil
}

func GetRealizedLinksFromAppSpec(appID string, sequence int64) ([]types.RealizedLink, error) {
	db := persistence.MustGetDBSession()
	query := `select app_spec, kots_app_spec from app_version where app_id = $1 and sequence = $2`
	row := db.QueryRow(query, appID, sequence)

	var appSpecStr sql.NullString
	var kotsAppSpecStr sql.NullString
	if err := row.Scan(&appSpecStr, &kotsAppSpecStr); err != nil {
		if err == sql.ErrNoRows {
			return []types.RealizedLink{}, nil
		}
		return nil, errors.Wrap(err, "failed to scan")
	}

	if appSpecStr.String == "" {
		return []types.RealizedLink{}, nil
	}

	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, _, err := decode([]byte(appSpecStr.String), nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode app spec yaml")
	}
	appSpec := obj.(*applicationv1beta1.Application)

	obj, _, err = decode([]byte(kotsAppSpecStr.String), nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode kots app spec yaml")
	}
	kotsAppSpec := obj.(*kotsv1beta1.Application)

	realizedLinks := []types.RealizedLink{}
	for _, link := range appSpec.Spec.Descriptor.Links {
		rewrittenURL := link.URL
		for _, port := range kotsAppSpec.Spec.ApplicationPorts {
			if port.ApplicationURL == link.URL {
				rewrittenURL = fmt.Sprintf("http://localhost:%d", port.LocalPort)
			}
		}
		realizedLink := types.RealizedLink{
			Title: link.Description,
			Uri:   rewrittenURL,
		}
		realizedLinks = append(realizedLinks, realizedLink)
	}

	return realizedLinks, nil
}

func GetForwardedPortsFromAppSpec(appID string, sequence int64) ([]types.ForwardedPort, error) {
	appVersion, err := store.GetStore().GetAppVersion(appID, sequence)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get app version")
	}

	if appVersion.KOTSKinds == nil {
		return nil, errors.Wrap(err, "failed to get kots kinds for app")
	}

	kotsAppSpec := appVersion.KOTSKinds.KotsApplication

	if len(kotsAppSpec.Spec.ApplicationPorts) == 0 {
		return []types.ForwardedPort{}, nil
	}

	ports := []types.ForwardedPort{}

	clientset, err := k8sutil.GetClientset()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get clientset")
	}

	appNamespace := util.AppNamespace()

	// To forward the ports, we need to have the port listed
	// in the kots spec only
	for _, port := range kotsAppSpec.Spec.ApplicationPorts {
		// make a best effort to not return ports to services that are not yet ready
		// this is best effort because the service could restart at any time
		// and the RBAC persona that this api is running as does not match
		// the users RBAC persona. Finally, this will not work in a gitops-application
		// unless it's deployed to the same namespace as the admin console
		// This has always been a limitation though, we need to design for this

		svc, err := clientset.CoreV1().Services(appNamespace).Get(context.TODO(), port.ServiceName, metav1.GetOptions{})
		if err != nil {
			logger.Error(errors.Wrapf(err, "failed to get service to check status, namespace = %s", util.PodNamespace))
			continue
		}

		options := metav1.ListOptions{LabelSelector: labels.SelectorFromSet(svc.Spec.Selector).String()}
		podList, err := clientset.CoreV1().Pods(appNamespace).List(context.TODO(), options)
		if err != nil {
			logger.Error(errors.Wrap(err, "failed to list pods in service"))
			continue
		}

		hasReadyPod := false
		for _, pod := range podList.Items {
			if pod.Status.Phase == corev1.PodRunning {
				for _, status := range pod.Status.ContainerStatuses {
					if status.Ready {
						hasReadyPod = true
					}
				}
			}
		}

		if !hasReadyPod {
			logger.Info("not forwarding to service because no pods are ready", zap.String("serviceName", port.ServiceName), zap.String("namespace", util.PodNamespace))
			continue
		}

		ports = append(ports, types.ForwardedPort{
			ServiceName:    port.ServiceName,
			ServicePort:    port.ServicePort,
			LocalPort:      port.LocalPort,
			ApplicationURL: port.ApplicationURL,
		})
	}

	return ports, nil
}

func getStatus(appID string, sequence int64) (storetypes.DownstreamVersionStatus, error) {
	db := persistence.MustGetDBSession()
	var status sql.NullString
	query := `
SELECT
	adv.status
FROM
	app_downstream_version AS adv
WHERE
	adv.app_id = $1 AND
	adv.parent_sequence = $2`
	row := db.QueryRow(query, appID, sequence)
	if err := row.Scan(
		&status,
	); err != nil {
		return "", errors.Wrap(err, "failed to scan")
	}
	return storetypes.DownstreamVersionStatus(status.String), nil
}

func isBlockedDueToStricPreFlights(appID string, sequence int64) (bool, error) {
	var (
		status             sql.NullString
		preflightSkipped   sql.NullBool
		preflightSpecStr   sql.NullString
		preflightResultStr sql.NullString
	)
	db := persistence.MustGetDBSession()
	query := `
SELECT
	adv.status,
	adv.preflight_skipped,
	av.preflight_spec,
	adv.preflight_result
FROM
	app_version AS av
LEFT JOIN
	app_downstream_version AS adv
ON
	adv.app_id = av.app_id AND adv.parent_sequence = av.sequence
WHERE
	av.app_id = $1 AND
 	av.sequence = $2`

	row := db.QueryRow(query, appID, sequence)
	if err := row.Scan(
		&status,
		&preflightSkipped,
		&preflightSpecStr,
		&preflightResultStr,
	); err != nil {
		return false, errors.Wrap(err, "failed to scan")
	}

	hasStrictPreflights := false
	if preflightSpecStr.Valid && preflightSpecStr.String != "" {
		preflight, err := kotsutil.LoadPreflightFromContents([]byte(preflightSpecStr.String))
		if err != nil {
			return false, errors.Wrap(err, "failed to load preflights from spec")
		}
		hasStrictPreflights = kotsutil.HasStrictPreflights(preflight)
	}

	// if preflights were not skipped and status is pending_preflight, poll till the status gets updated
	// if preflights were skipped don't poll and check results
	if hasStrictPreflights && !preflightSkipped.Bool && storetypes.DownstreamVersionStatus(status.String) == storetypes.VersionPendingPreflight {

		err := wait.PollImmediateInfinite(2*time.Second, func() (bool, error) {
			versionStatus, err := getStatus(appID, sequence)
			if err != nil {
				return false, errors.Wrap(err, "failed get status")
			}
			if versionStatus != storetypes.VersionPendingPreflight {
				return true, nil
			}
			return false, nil
		})
		if err != nil {
			return false, errors.Wrap(err, "failed to poll for preflights results")
		}

		// fetch latest results
		row = db.QueryRow(query, appID, sequence)
		if err := row.Scan(
			&status,
			&preflightSkipped,
			&preflightSpecStr,
			&preflightResultStr,
		); err != nil {
			return false, errors.Wrap(err, "failed to scan")
		}
	}

	preflightResult := troubleshootpreflight.UploadPreflightResults{}
	if preflightResultStr.Valid && preflightResultStr.String != "" {
		if err := json.Unmarshal([]byte(preflightResultStr.String), &preflightResult); err != nil {
			return false, errors.Wrap(err, "failed to unmarshal preflightResults")
		}
	}
	return hasStrictPreflights && kotsutil.IsStrictPreflightFailing(&preflightResult), nil
}
