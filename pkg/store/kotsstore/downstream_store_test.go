package kotsstore

import (
	"testing"

	"github.com/blang/semver"
	downstreamtypes "github.com/replicatedhq/kots/pkg/api/downstream/types"
	"github.com/stretchr/testify/assert"
)

func Test_isSameUpstreamRelease(t *testing.T) {
	tests := []struct {
		name             string
		v1               *downstreamtypes.DownstreamVersion
		v2               *downstreamtypes.DownstreamVersion
		isSemverRequired bool
		expected         bool
	}{
		{
			name: "non-semver, same channel, different cursor",
			v1: &downstreamtypes.DownstreamVersion{
				ChannelID: "channel-id-1",
				Cursor:    2,
			},
			v2: &downstreamtypes.DownstreamVersion{
				ChannelID: "channel-id-1",
				Cursor:    3,
			},
			expected: false,
		},
		{
			name: "non-semver, same cursor, different channel",
			v1: &downstreamtypes.DownstreamVersion{
				ChannelID: "channel-id-1",
				Cursor:    2,
			},
			v2: &downstreamtypes.DownstreamVersion{
				ChannelID: "channel-id-2",
				Cursor:    2,
			},
			expected: false,
		},
		{
			name: "non-semver, same cursor, same channel",
			v1: &downstreamtypes.DownstreamVersion{
				ChannelID: "channel-id-1",
				Cursor:    2,
			},
			v2: &downstreamtypes.DownstreamVersion{
				ChannelID: "channel-id-1",
				Cursor:    2,
			},
			expected: true,
		},
		{
			name: "semver, same channel, same cursor, same semver",
			v1: &downstreamtypes.DownstreamVersion{
				ChannelID: "channel-id-1",
				Cursor:    2,
				Semver: &semver.Version{
					Major: 3,
				},
			},
			v2: &downstreamtypes.DownstreamVersion{
				ChannelID: "channel-id-1",
				Cursor:    2,
				Semver: &semver.Version{
					Major: 3,
				},
			},
			isSemverRequired: true,
			expected:         true,
		},
		{
			name: "semver, same channel, same cursor, different semver",
			v1: &downstreamtypes.DownstreamVersion{
				ChannelID: "channel-id-1",
				Cursor:    2,
				Semver: &semver.Version{
					Major: 3,
				},
			},
			v2: &downstreamtypes.DownstreamVersion{
				ChannelID: "channel-id-1",
				Cursor:    2,
				Semver: &semver.Version{
					Major: 4,
				},
			},
			isSemverRequired: true,
			expected:         true,
		},
		{
			name: "semver, different channel, same cursor, different semver",
			v1: &downstreamtypes.DownstreamVersion{
				ChannelID: "channel-id-1",
				Cursor:    2,
				Semver: &semver.Version{
					Major: 3,
				},
			},
			v2: &downstreamtypes.DownstreamVersion{
				ChannelID: "channel-id-2",
				Cursor:    2,
				Semver: &semver.Version{
					Major: 4,
				},
			},
			isSemverRequired: true,
			expected:         false,
		},
		{
			name: "semver, different channel, same cursor, same semver",
			v1: &downstreamtypes.DownstreamVersion{
				ChannelID: "channel-id-1",
				Cursor:    2,
				Semver: &semver.Version{
					Major: 4,
				},
			},
			v2: &downstreamtypes.DownstreamVersion{
				ChannelID: "channel-id-2",
				Cursor:    2,
				Semver: &semver.Version{
					Major: 4,
				},
			},
			isSemverRequired: true,
			expected:         true,
		},
		{
			name: "semver, same channel, different cursor, different semver",
			v1: &downstreamtypes.DownstreamVersion{
				ChannelID: "channel-id-1",
				Cursor:    2,
				Semver: &semver.Version{
					Major: 3,
				},
			},
			v2: &downstreamtypes.DownstreamVersion{
				ChannelID: "channel-id-1",
				Cursor:    3,
				Semver: &semver.Version{
					Major: 4,
				},
			},
			isSemverRequired: true,
			expected:         false,
		},
		{
			name: "semver, same channel, different cursor, same semver",
			v1: &downstreamtypes.DownstreamVersion{
				ChannelID: "channel-id-1",
				Cursor:    2,
				Semver: &semver.Version{
					Major: 3,
				},
			},
			v2: &downstreamtypes.DownstreamVersion{
				ChannelID: "channel-id-1",
				Cursor:    3,
				Semver: &semver.Version{
					Major: 3,
				},
			},
			isSemverRequired: true,
			expected:         true,
		},
		{
			name: "semver, different channel, different cursor, same semver",
			v1: &downstreamtypes.DownstreamVersion{
				ChannelID: "channel-id-1",
				Cursor:    2,
				Semver: &semver.Version{
					Major: 4,
				},
			},
			v2: &downstreamtypes.DownstreamVersion{
				ChannelID: "channel-id-2",
				Cursor:    3,
				Semver: &semver.Version{
					Major: 4,
				},
			},
			isSemverRequired: true,
			expected:         true,
		},
		{
			name: "semver, different channel, different cursor, different semver",
			v1: &downstreamtypes.DownstreamVersion{
				ChannelID: "channel-id-1",
				Cursor:    2,
				Semver: &semver.Version{
					Major: 3,
				},
			},
			v2: &downstreamtypes.DownstreamVersion{
				ChannelID: "channel-id-2",
				Cursor:    3,
				Semver: &semver.Version{
					Major: 4,
				},
			},
			isSemverRequired: true,
			expected:         false,
		},
		{
			name: "semver and non-semver, same channel, same cursor",
			v1: &downstreamtypes.DownstreamVersion{
				ChannelID: "channel-id-1",
				Cursor:    2,
				Semver: &semver.Version{
					Major: 3,
				},
			},
			v2: &downstreamtypes.DownstreamVersion{
				ChannelID: "channel-id-1",
				Cursor:    2,
			},
			isSemverRequired: true,
			expected:         true,
		},
		{
			name: "semver and non-semver, different channel, same cursor",
			v1: &downstreamtypes.DownstreamVersion{
				ChannelID: "channel-id-1",
				Cursor:    2,
				Semver: &semver.Version{
					Major: 3,
				},
			},
			v2: &downstreamtypes.DownstreamVersion{
				ChannelID: "channel-id-2",
				Cursor:    2,
			},
			isSemverRequired: true,
			expected:         false,
		},
		{
			name: "semver and non-semver, same channel, different cursor",
			v1: &downstreamtypes.DownstreamVersion{
				ChannelID: "channel-id-1",
				Cursor:    2,
				Semver: &semver.Version{
					Major: 3,
				},
			},
			v2: &downstreamtypes.DownstreamVersion{
				ChannelID: "channel-id-1",
				Cursor:    3,
			},
			isSemverRequired: true,
			expected:         false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			same := isSameUpstreamRelease(test.v1, test.v2, test.isSemverRequired)
			assert.Equal(t, test.expected, same)
		})
	}
}

func Test_isAppVersionDeployable(t *testing.T) {
	tests := []struct {
		name                 string
		version              *downstreamtypes.DownstreamVersion
		appVersions          *downstreamtypes.DownstreamVersions
		isSemverRequired     bool
		expectedIsDeployable bool
		expectedCause        string
	}{
		{
			name: "failing strict preflights",
			version: &downstreamtypes.DownstreamVersion{
				HasFailingStrictPreflights: true,
			},
			expectedIsDeployable: false,
			expectedCause:        "Deployment is disabled as a strict analyzer in this version's preflight checks has failed or has not been run.",
		},
		{
			name:                 "no version is deployed yet",
			version:              &downstreamtypes.DownstreamVersion{},
			appVersions:          &downstreamtypes.DownstreamVersions{},
			expectedIsDeployable: true,
			expectedCause:        "",
		},
		{
			name: "version is same as deployed version",
			version: &downstreamtypes.DownstreamVersion{
				Sequence: 7,
			},
			appVersions: &downstreamtypes.DownstreamVersions{
				CurrentVersion: &downstreamtypes.DownstreamVersion{
					Sequence: 7,
				},
			},
			expectedIsDeployable: true,
			expectedCause:        "",
		},
		/* ---- Non semver tests begin here ---- */
		{
			name: "non-semver -- deployed version is from a different channel, not required, no required releases in between",
			version: &downstreamtypes.DownstreamVersion{
				Sequence:  1,
				ChannelID: "channel-id-1",
			},
			appVersions: &downstreamtypes.DownstreamVersions{
				AllVersions: []*downstreamtypes.DownstreamVersion{
					{
						Sequence:  1,
						ChannelID: "channel-id-1",
					},
					{
						Sequence:  0,
						ChannelID: "channel-id-2",
					},
				},
				CurrentVersion: &downstreamtypes.DownstreamVersion{
					Sequence:  0,
					ChannelID: "channel-id-2",
				},
			},
			expectedIsDeployable: true,
			expectedCause:        "",
		},
		{
			name: "non-semver -- deployed version is from a different channel, is required, no required releases in between",
			version: &downstreamtypes.DownstreamVersion{
				Sequence:  1,
				ChannelID: "channel-id-1",
			},
			appVersions: &downstreamtypes.DownstreamVersions{
				AllVersions: []*downstreamtypes.DownstreamVersion{
					{
						Sequence:  1,
						ChannelID: "channel-id-1",
					},
					{
						Sequence:   0,
						ChannelID:  "channel-id-2",
						IsRequired: true,
					},
				},
				CurrentVersion: &downstreamtypes.DownstreamVersion{
					Sequence:   0,
					ChannelID:  "channel-id-2",
					IsRequired: true,
				},
			},
			expectedIsDeployable: true,
			expectedCause:        "",
		},
		{
			name: "non-semver -- deployed version is from a different channel, is required, required releases in between from different channel",
			version: &downstreamtypes.DownstreamVersion{
				Sequence:  3,
				ChannelID: "channel-id-1",
			},
			appVersions: &downstreamtypes.DownstreamVersions{
				AllVersions: []*downstreamtypes.DownstreamVersion{
					{
						Sequence:  3,
						ChannelID: "channel-id-1",
					},
					{
						Sequence:   2,
						ChannelID:  "channel-id-2",
						IsRequired: true,
					},
					{
						Sequence:   1,
						ChannelID:  "channel-id-2",
						IsRequired: true,
					},
					{
						Sequence:   0,
						ChannelID:  "channel-id-2",
						IsRequired: true,
					},
				},
				CurrentVersion: &downstreamtypes.DownstreamVersion{
					Sequence:   0,
					ChannelID:  "channel-id-2",
					IsRequired: true,
				},
			},
			expectedIsDeployable: true,
			expectedCause:        "",
		},
		{
			name: "non-semver -- deployed version is from a different channel, not required, required releases in between from same channel, same variants",
			version: &downstreamtypes.DownstreamVersion{
				Sequence:   3,
				ChannelID:  "channel-id-1",
				IsRequired: true,
				Cursor:     3,
			},
			appVersions: &downstreamtypes.DownstreamVersions{
				AllVersions: []*downstreamtypes.DownstreamVersion{
					{
						Sequence:   3,
						ChannelID:  "channel-id-1",
						IsRequired: true,
						Cursor:     3,
					},
					{
						Sequence:   2,
						ChannelID:  "channel-id-1",
						IsRequired: true,
						Cursor:     3,
					},
					{
						Sequence:   1,
						ChannelID:  "channel-id-1",
						IsRequired: true,
						Cursor:     3,
					},
					{
						Sequence:  0,
						ChannelID: "channel-id-2",
						Cursor:    1,
					},
				},
				CurrentVersion: &downstreamtypes.DownstreamVersion{
					Sequence:  0,
					ChannelID: "channel-id-2",
					Cursor:    1,
				},
			},
			expectedIsDeployable: true,
			expectedCause:        "",
		},
		{
			name: "non-semver -- deployed version is from a different channel, not required, required releases in between from same channel, different variants",
			version: &downstreamtypes.DownstreamVersion{
				Sequence:   3,
				ChannelID:  "channel-id-1",
				IsRequired: true,
				Cursor:     4,
			},
			appVersions: &downstreamtypes.DownstreamVersions{
				AllVersions: []*downstreamtypes.DownstreamVersion{
					{
						Sequence:   3,
						ChannelID:  "channel-id-1",
						IsRequired: true,
						Cursor:     4,
					},
					{
						Sequence:   2,
						ChannelID:  "channel-id-1",
						IsRequired: true,
						Cursor:     3,
					},
					{
						Sequence:   1,
						ChannelID:  "channel-id-1",
						IsRequired: true,
						Cursor:     2,
					},
					{
						Sequence:  0,
						ChannelID: "channel-id-2",
						Cursor:    1,
					},
				},
				CurrentVersion: &downstreamtypes.DownstreamVersion{
					Sequence:  0,
					ChannelID: "channel-id-2",
					Cursor:    1,
				},
			},
			expectedIsDeployable: true,
			expectedCause:        "",
		},
		{
			name: "non-semver -- deployed version is from same channel, not required, required releases in between from same channel, same variants as version",
			version: &downstreamtypes.DownstreamVersion{
				Sequence:   3,
				ChannelID:  "channel-id-1",
				IsRequired: true,
				Cursor:     2,
			},
			appVersions: &downstreamtypes.DownstreamVersions{
				AllVersions: []*downstreamtypes.DownstreamVersion{
					{
						Sequence:   3,
						ChannelID:  "channel-id-1",
						IsRequired: true,
						Cursor:     2,
					},
					{
						Sequence:   2,
						ChannelID:  "channel-id-1",
						IsRequired: true,
						Cursor:     2,
					},
					{
						Sequence:   1,
						ChannelID:  "channel-id-1",
						IsRequired: true,
						Cursor:     2,
					},
					{
						Sequence:  0,
						ChannelID: "channel-id-1",
						Cursor:    1,
					},
				},
				CurrentVersion: &downstreamtypes.DownstreamVersion{
					Sequence:  0,
					ChannelID: "channel-id-1",
					Cursor:    1,
				},
			},
			expectedIsDeployable: true,
			expectedCause:        "",
		},
		{
			name: "non-semver -- deployed version is from same channel, not required, required releases in between from same channel, same variants as deployed version",
			version: &downstreamtypes.DownstreamVersion{
				Sequence:   3,
				ChannelID:  "channel-id-1",
				IsRequired: true,
				Cursor:     2,
			},
			appVersions: &downstreamtypes.DownstreamVersions{
				AllVersions: []*downstreamtypes.DownstreamVersion{
					{
						Sequence:   3,
						ChannelID:  "channel-id-1",
						IsRequired: true,
						Cursor:     2,
					},
					{
						Sequence:   2,
						ChannelID:  "channel-id-1",
						IsRequired: true,
						Cursor:     1,
					},
					{
						Sequence:   1,
						ChannelID:  "channel-id-1",
						IsRequired: true,
						Cursor:     1,
					},
					{
						Sequence:   0,
						ChannelID:  "channel-id-1",
						IsRequired: true,
						Cursor:     1,
					},
				},
				CurrentVersion: &downstreamtypes.DownstreamVersion{
					Sequence:   0,
					ChannelID:  "channel-id-1",
					IsRequired: true,
					Cursor:     1,
				},
			},
			expectedIsDeployable: true,
			expectedCause:        "",
		},
		{
			name: "non-semver -- deployed version is from same channel, is required, required releases in between from same channel, different variants, lower cursor",
			version: &downstreamtypes.DownstreamVersion{
				Sequence:     4,
				ChannelID:    "channel-id-1",
				Cursor:       5,
				VersionLabel: "5.0",
			},
			appVersions: &downstreamtypes.DownstreamVersions{
				AllVersions: []*downstreamtypes.DownstreamVersion{
					{
						Sequence:     4,
						ChannelID:    "channel-id-1",
						Cursor:       5,
						VersionLabel: "5.0",
					},
					{
						Sequence:     3,
						ChannelID:    "channel-id-1",
						IsRequired:   true,
						Cursor:       4,
						VersionLabel: "4.0",
					},
					{
						Sequence:     2,
						ChannelID:    "channel-id-1",
						IsRequired:   true,
						Cursor:       3,
						VersionLabel: "3.0",
					},
					{
						Sequence:     1,
						ChannelID:    "channel-id-1",
						Cursor:       2,
						IsRequired:   true,
						VersionLabel: "2.0",
					},
					{
						Sequence:     0,
						ChannelID:    "channel-id-1",
						Cursor:       1,
						IsRequired:   true,
						VersionLabel: "1.0",
					},
				},
				CurrentVersion: &downstreamtypes.DownstreamVersion{
					Sequence:     1,
					ChannelID:    "channel-id-1",
					Cursor:       2,
					IsRequired:   true,
					VersionLabel: "2.0",
				},
			},
			expectedIsDeployable: false,
			expectedCause:        "This version cannot be deployed because versions 3.0, 4.0 are required and must be deployed first.",
		},
		{
			name: "non-semver -- deployed version is from same channel, is required, required releases in between from same channel, different variants, 2 higher cursor and 1 lower cursor from different channel",
			version: &downstreamtypes.DownstreamVersion{
				Sequence:     5,
				ChannelID:    "channel-id-1",
				Cursor:       2,
				IsRequired:   true,
				VersionLabel: "2.0",
			},
			appVersions: &downstreamtypes.DownstreamVersions{
				AllVersions: []*downstreamtypes.DownstreamVersion{
					{
						Sequence:     5,
						ChannelID:    "channel-id-1",
						Cursor:       2,
						IsRequired:   true,
						VersionLabel: "2.0",
					},
					{
						Sequence:     4,
						ChannelID:    "channel-id-1",
						IsRequired:   true,
						Cursor:       4,
						VersionLabel: "4.0",
					},
					{
						Sequence:     3,
						ChannelID:    "channel-id-2",
						IsRequired:   true,
						Cursor:       1,
						VersionLabel: "3.1",
					},
					{
						Sequence:     2,
						ChannelID:    "channel-id-1",
						IsRequired:   true,
						Cursor:       3,
						VersionLabel: "3.0",
					},
					{
						Sequence:     1,
						ChannelID:    "channel-id-1",
						Cursor:       2,
						IsRequired:   true,
						VersionLabel: "2.0",
					},
					{
						Sequence:     0,
						ChannelID:    "channel-id-1",
						Cursor:       1,
						IsRequired:   true,
						VersionLabel: "1.0",
					},
				},
				CurrentVersion: &downstreamtypes.DownstreamVersion{
					Sequence:     1,
					ChannelID:    "channel-id-1",
					Cursor:       2,
					IsRequired:   true,
					VersionLabel: "2.0",
				},
			},
			expectedIsDeployable: true,
			expectedCause:        "",
		},
		{
			name: "non-semver -- deployed version is from same channel, not required, required and non-required releases in between from same channel, different variants",
			version: &downstreamtypes.DownstreamVersion{
				Sequence:     5,
				ChannelID:    "channel-id-1",
				Cursor:       6,
				VersionLabel: "6.0",
			},
			appVersions: &downstreamtypes.DownstreamVersions{
				AllVersions: []*downstreamtypes.DownstreamVersion{
					{
						Sequence:     5,
						ChannelID:    "channel-id-1",
						Cursor:       6,
						VersionLabel: "6.0",
					},
					{
						Sequence:     4,
						ChannelID:    "channel-id-1",
						IsRequired:   true,
						Cursor:       5,
						VersionLabel: "5.0",
					},
					{
						Sequence:     3,
						ChannelID:    "channel-id-1",
						Cursor:       4,
						VersionLabel: "4.0",
					},
					{
						Sequence:     2,
						ChannelID:    "channel-id-1",
						IsRequired:   true,
						Cursor:       3,
						VersionLabel: "3.0",
					},
					{
						Sequence:     1,
						ChannelID:    "channel-id-1",
						Cursor:       2,
						VersionLabel: "2.0",
					},
					{
						Sequence:     0,
						ChannelID:    "channel-id-1",
						Cursor:       1,
						IsRequired:   true,
						VersionLabel: "1.0",
					},
				},
				CurrentVersion: &downstreamtypes.DownstreamVersion{
					Sequence:     1,
					ChannelID:    "channel-id-1",
					Cursor:       2,
					VersionLabel: "2.0",
				},
			},
			expectedIsDeployable: false,
			expectedCause:        "This version cannot be deployed because versions 3.0, 5.0 are required and must be deployed first.",
		},
		{
			name: "non-semver -- deployed version is from same channel, not required, required releases in between from same channel and same variants",
			version: &downstreamtypes.DownstreamVersion{
				Sequence:     5,
				ChannelID:    "channel-id-1",
				Cursor:       6,
				VersionLabel: "4.0",
			},
			appVersions: &downstreamtypes.DownstreamVersions{
				AllVersions: []*downstreamtypes.DownstreamVersion{
					{
						Sequence:     5,
						ChannelID:    "channel-id-1",
						Cursor:       6,
						VersionLabel: "4.0",
					},
					{
						Sequence:     4,
						ChannelID:    "channel-id-1",
						IsRequired:   true,
						Cursor:       3,
						VersionLabel: "3.0",
					},
					{
						Sequence:     3,
						ChannelID:    "channel-id-1",
						IsRequired:   true,
						Cursor:       3,
						VersionLabel: "3.0",
					},
					{
						Sequence:     2,
						ChannelID:    "channel-id-1",
						IsRequired:   true,
						Cursor:       3,
						VersionLabel: "3.0",
					},
					{
						Sequence:     1,
						ChannelID:    "channel-id-1",
						Cursor:       2,
						VersionLabel: "2.0",
					},
					{
						Sequence:     0,
						ChannelID:    "channel-id-1",
						Cursor:       1,
						IsRequired:   true,
						VersionLabel: "1.0",
					},
				},
				CurrentVersion: &downstreamtypes.DownstreamVersion{
					Sequence:     1,
					ChannelID:    "channel-id-1",
					Cursor:       2,
					VersionLabel: "2.0",
				},
			},
			expectedIsDeployable: false,
			expectedCause:        "This version cannot be deployed because version 3.0 is required and must be deployed first.",
		},
		/* ---- Non semver tests end here ---- */

		/* ---- Semver tests begin here ---- */
		{
			name: "semver -- deployed version is from a different channel, not required, no required releases in between",
			version: &downstreamtypes.DownstreamVersion{
				Sequence:     1,
				ChannelID:    "channel-id-1",
				Cursor:       1,
				VersionLabel: "2.0",
				Semver: &semver.Version{
					Major: 2,
				},
			},
			appVersions: &downstreamtypes.DownstreamVersions{
				AllVersions: []*downstreamtypes.DownstreamVersion{
					{
						Sequence:     1,
						ChannelID:    "channel-id-1",
						Cursor:       1,
						VersionLabel: "2.0",
						Semver: &semver.Version{
							Major: 2,
						},
					},
					{
						Sequence:     0,
						ChannelID:    "channel-id-2",
						Cursor:       1,
						VersionLabel: "1.0",
						Semver: &semver.Version{
							Major: 1,
						},
					},
				},
				CurrentVersion: &downstreamtypes.DownstreamVersion{
					Sequence:     0,
					ChannelID:    "channel-id-2",
					Cursor:       1,
					VersionLabel: "1.0",
					Semver: &semver.Version{
						Major: 1,
					},
				},
			},
			isSemverRequired:     true,
			expectedIsDeployable: true,
			expectedCause:        "",
		},
		{
			name: "semver -- deployed version is from a different channel, is required, no required releases in between",
			version: &downstreamtypes.DownstreamVersion{
				Sequence:     1,
				ChannelID:    "channel-id-1",
				Cursor:       1,
				VersionLabel: "2.0",
				Semver: &semver.Version{
					Major: 2,
				},
			},
			appVersions: &downstreamtypes.DownstreamVersions{
				AllVersions: []*downstreamtypes.DownstreamVersion{
					{
						Sequence:     1,
						ChannelID:    "channel-id-1",
						Cursor:       1,
						VersionLabel: "2.0",
						Semver: &semver.Version{
							Major: 2,
						},
					},
					{
						Sequence:     0,
						ChannelID:    "channel-id-2",
						IsRequired:   true,
						Cursor:       1,
						VersionLabel: "1.0",
						Semver: &semver.Version{
							Major: 1,
						},
					},
				},
				CurrentVersion: &downstreamtypes.DownstreamVersion{
					Sequence:     0,
					ChannelID:    "channel-id-2",
					IsRequired:   true,
					Cursor:       1,
					VersionLabel: "1.0",
					Semver: &semver.Version{
						Major: 1,
					},
				},
			},
			isSemverRequired:     true,
			expectedIsDeployable: true,
			expectedCause:        "",
		},
		{
			name: "semver -- deployed version is from a different channel, is required, required releases in between from different channel",
			version: &downstreamtypes.DownstreamVersion{
				Sequence:     3,
				ChannelID:    "channel-id-1",
				Cursor:       4,
				VersionLabel: "4.0",
				Semver: &semver.Version{
					Major: 4,
				},
			},
			appVersions: &downstreamtypes.DownstreamVersions{
				AllVersions: []*downstreamtypes.DownstreamVersion{
					{
						Sequence:     3,
						ChannelID:    "channel-id-1",
						Cursor:       4,
						VersionLabel: "4.0",
						Semver: &semver.Version{
							Major: 4,
						},
					},
					{
						Sequence:     2,
						ChannelID:    "channel-id-2",
						IsRequired:   true,
						Cursor:       3,
						VersionLabel: "3.0",
						Semver: &semver.Version{
							Major: 3,
						},
					},
					{
						Sequence:     1,
						ChannelID:    "channel-id-2",
						IsRequired:   true,
						Cursor:       2,
						VersionLabel: "2.0",
						Semver: &semver.Version{
							Major: 2,
						},
					},
					{
						Sequence:     0,
						ChannelID:    "channel-id-2",
						IsRequired:   true,
						Cursor:       1,
						VersionLabel: "1.0",
						Semver: &semver.Version{
							Major: 1,
						},
					},
				},
				CurrentVersion: &downstreamtypes.DownstreamVersion{
					Sequence:     0,
					ChannelID:    "channel-id-2",
					IsRequired:   true,
					Cursor:       1,
					VersionLabel: "1.0",
					Semver: &semver.Version{
						Major: 1,
					},
				},
			},
			isSemverRequired:     true,
			expectedIsDeployable: false,
			expectedCause:        "This version cannot be deployed because versions 2.0, 3.0 are required and must be deployed first.",
		},
		{
			name: "semver -- deployed version is from a different channel, not required, required releases in between from same channel, same variants",
			version: &downstreamtypes.DownstreamVersion{
				Sequence:     3,
				ChannelID:    "channel-id-1",
				IsRequired:   true,
				Cursor:       2,
				VersionLabel: "2.0",
				Semver: &semver.Version{
					Major: 2,
				},
			},
			appVersions: &downstreamtypes.DownstreamVersions{
				AllVersions: []*downstreamtypes.DownstreamVersion{
					{
						Sequence:     3,
						ChannelID:    "channel-id-1",
						IsRequired:   true,
						Cursor:       2,
						VersionLabel: "2.0",
						Semver: &semver.Version{
							Major: 2,
						},
					},
					{
						Sequence:     2,
						ChannelID:    "channel-id-1",
						IsRequired:   true,
						Cursor:       2,
						VersionLabel: "2.0",
						Semver: &semver.Version{
							Major: 2,
						},
					},
					{
						Sequence:     1,
						ChannelID:    "channel-id-1",
						IsRequired:   true,
						Cursor:       2,
						VersionLabel: "2.0",
						Semver: &semver.Version{
							Major: 2,
						},
					},
					{
						Sequence:     0,
						ChannelID:    "channel-id-2",
						Cursor:       1,
						VersionLabel: "1.0",
						Semver: &semver.Version{
							Major: 1,
						},
					},
				},
				CurrentVersion: &downstreamtypes.DownstreamVersion{
					Sequence:     0,
					ChannelID:    "channel-id-2",
					Cursor:       1,
					VersionLabel: "1.0",
					Semver: &semver.Version{
						Major: 1,
					},
				},
			},
			isSemverRequired:     true,
			expectedIsDeployable: true,
			expectedCause:        "",
		},
		{
			name: "semver -- deployed version is from a different channel, not required, required releases in between from same channel, different variants",
			version: &downstreamtypes.DownstreamVersion{
				Sequence:     3,
				ChannelID:    "channel-id-1",
				IsRequired:   true,
				Cursor:       4,
				VersionLabel: "4.0",
				Semver: &semver.Version{
					Major: 4,
				},
			},
			appVersions: &downstreamtypes.DownstreamVersions{
				AllVersions: []*downstreamtypes.DownstreamVersion{
					{
						Sequence:     3,
						ChannelID:    "channel-id-1",
						IsRequired:   true,
						Cursor:       4,
						VersionLabel: "4.0",
						Semver: &semver.Version{
							Major: 4,
						},
					},
					{
						Sequence:     2,
						ChannelID:    "channel-id-1",
						IsRequired:   true,
						Cursor:       3,
						VersionLabel: "3.0",
						Semver: &semver.Version{
							Major: 3,
						},
					},
					{
						Sequence:     1,
						ChannelID:    "channel-id-1",
						IsRequired:   true,
						Cursor:       2,
						VersionLabel: "2.0",
						Semver: &semver.Version{
							Major: 2,
						},
					},
					{
						Sequence:     0,
						ChannelID:    "channel-id-2",
						Cursor:       1,
						VersionLabel: "1.0",
						Semver: &semver.Version{
							Major: 1,
						},
					},
				},
				CurrentVersion: &downstreamtypes.DownstreamVersion{
					Sequence:     0,
					ChannelID:    "channel-id-2",
					Cursor:       1,
					VersionLabel: "1.0",
					Semver: &semver.Version{
						Major: 1,
					},
				},
			},
			isSemverRequired:     true,
			expectedIsDeployable: false,
			expectedCause:        "This version cannot be deployed because versions 2.0, 3.0 are required and must be deployed first.",
		},
		{
			name: "semver -- deployed version is from same channel, not required, required releases in between from same channel, same variants as version",
			version: &downstreamtypes.DownstreamVersion{
				Sequence:     3,
				ChannelID:    "channel-id-1",
				IsRequired:   true,
				Cursor:       2,
				VersionLabel: "2.0",
				Semver: &semver.Version{
					Major: 2,
				},
			},
			appVersions: &downstreamtypes.DownstreamVersions{
				AllVersions: []*downstreamtypes.DownstreamVersion{
					{
						Sequence:     3,
						ChannelID:    "channel-id-1",
						IsRequired:   true,
						Cursor:       2,
						VersionLabel: "2.0",
						Semver: &semver.Version{
							Major: 2,
						},
					},
					{
						Sequence:     2,
						ChannelID:    "channel-id-1",
						IsRequired:   true,
						Cursor:       2,
						VersionLabel: "2.0",
						Semver: &semver.Version{
							Major: 2,
						},
					},
					{
						Sequence:     1,
						ChannelID:    "channel-id-1",
						IsRequired:   true,
						Cursor:       2,
						VersionLabel: "2.0",
						Semver: &semver.Version{
							Major: 2,
						},
					},
					{
						Sequence:     0,
						ChannelID:    "channel-id-1",
						Cursor:       1,
						VersionLabel: "1.0",
						Semver: &semver.Version{
							Major: 1,
						},
					},
				},
				CurrentVersion: &downstreamtypes.DownstreamVersion{
					Sequence:     0,
					ChannelID:    "channel-id-1",
					Cursor:       1,
					VersionLabel: "1.0",
					Semver: &semver.Version{
						Major: 1,
					},
				},
			},
			isSemverRequired:     true,
			expectedIsDeployable: true,
			expectedCause:        "",
		},
		{
			name: "semver -- deployed version is from same channel, not required, required releases in between from same channel, same variants as deployed version",
			version: &downstreamtypes.DownstreamVersion{
				Sequence:     3,
				ChannelID:    "channel-id-1",
				IsRequired:   true,
				Cursor:       2,
				VersionLabel: "2.0",
				Semver: &semver.Version{
					Major: 2,
				},
			},
			appVersions: &downstreamtypes.DownstreamVersions{
				AllVersions: []*downstreamtypes.DownstreamVersion{
					{
						Sequence:     3,
						ChannelID:    "channel-id-1",
						IsRequired:   true,
						Cursor:       2,
						VersionLabel: "2.0",
						Semver: &semver.Version{
							Major: 2,
						},
					},
					{
						Sequence:     2,
						ChannelID:    "channel-id-1",
						IsRequired:   true,
						Cursor:       1,
						VersionLabel: "1.0",
						Semver: &semver.Version{
							Major: 1,
						},
					},
					{
						Sequence:     1,
						ChannelID:    "channel-id-1",
						IsRequired:   true,
						Cursor:       1,
						VersionLabel: "1.0",
						Semver: &semver.Version{
							Major: 1,
						},
					},
					{
						Sequence:     0,
						ChannelID:    "channel-id-1",
						IsRequired:   true,
						Cursor:       1,
						VersionLabel: "1.0",
						Semver: &semver.Version{
							Major: 1,
						},
					},
				},
				CurrentVersion: &downstreamtypes.DownstreamVersion{
					Sequence:     0,
					ChannelID:    "channel-id-1",
					IsRequired:   true,
					Cursor:       1,
					VersionLabel: "1.0",
					Semver: &semver.Version{
						Major: 1,
					},
				},
			},
			isSemverRequired:     true,
			expectedIsDeployable: true,
			expectedCause:        "",
		},
		{
			name: "semver -- deployed version is from same channel, is required, required releases in between from same channel, different variants, lower cursor",
			version: &downstreamtypes.DownstreamVersion{
				Sequence:     4,
				ChannelID:    "channel-id-1",
				Cursor:       5,
				VersionLabel: "5.0",
				Semver: &semver.Version{
					Major: 5,
				},
			},
			appVersions: &downstreamtypes.DownstreamVersions{
				AllVersions: []*downstreamtypes.DownstreamVersion{
					{
						Sequence:     4,
						ChannelID:    "channel-id-1",
						Cursor:       5,
						VersionLabel: "5.0",
						Semver: &semver.Version{
							Major: 5,
						},
					},
					{
						Sequence:     3,
						ChannelID:    "channel-id-1",
						IsRequired:   true,
						Cursor:       4,
						VersionLabel: "4.0",
						Semver: &semver.Version{
							Major: 4,
						},
					},
					{
						Sequence:     2,
						ChannelID:    "channel-id-1",
						IsRequired:   true,
						Cursor:       3,
						VersionLabel: "3.0",
						Semver: &semver.Version{
							Major: 3,
						},
					},
					{
						Sequence:     1,
						ChannelID:    "channel-id-1",
						Cursor:       2,
						IsRequired:   true,
						VersionLabel: "2.0",
						Semver: &semver.Version{
							Major: 2,
						},
					},
					{
						Sequence:     0,
						ChannelID:    "channel-id-1",
						Cursor:       1,
						IsRequired:   true,
						VersionLabel: "1.0",
						Semver: &semver.Version{
							Major: 1,
						},
					},
				},
				CurrentVersion: &downstreamtypes.DownstreamVersion{
					Sequence:     1,
					ChannelID:    "channel-id-1",
					Cursor:       2,
					IsRequired:   true,
					VersionLabel: "2.0",
					Semver: &semver.Version{
						Major: 2,
					},
				},
			},
			isSemverRequired:     true,
			expectedIsDeployable: false,
			expectedCause:        "This version cannot be deployed because versions 3.0, 4.0 are required and must be deployed first.",
		},
		{
			name: "semver -- deployed version is from same channel, is required, required releases in between from same channel, different variants, 2 higher cursor and 1 lower cursor from different channel",
			version: &downstreamtypes.DownstreamVersion{
				Sequence:     5,
				ChannelID:    "channel-id-1",
				Cursor:       2,
				IsRequired:   true,
				VersionLabel: "5.0",
				Semver: &semver.Version{
					Major: 5,
				},
			},
			appVersions: &downstreamtypes.DownstreamVersions{
				AllVersions: []*downstreamtypes.DownstreamVersion{
					{
						Sequence:     5,
						ChannelID:    "channel-id-1",
						Cursor:       2,
						IsRequired:   true,
						VersionLabel: "5.0",
						Semver: &semver.Version{
							Major: 5,
						},
					},
					{
						Sequence:     4,
						ChannelID:    "channel-id-1",
						IsRequired:   true,
						Cursor:       4,
						VersionLabel: "4.0",
						Semver: &semver.Version{
							Major: 4,
						},
					},
					{
						Sequence:     3,
						ChannelID:    "channel-id-2",
						IsRequired:   true,
						Cursor:       1,
						VersionLabel: "3.1",
						Semver: &semver.Version{
							Major: 3,
							Minor: 1,
						},
					},
					{
						Sequence:     2,
						ChannelID:    "channel-id-1",
						IsRequired:   true,
						Cursor:       3,
						VersionLabel: "3.0",
						Semver: &semver.Version{
							Major: 3,
						},
					},
					{
						Sequence:     1,
						ChannelID:    "channel-id-1",
						Cursor:       2,
						IsRequired:   true,
						VersionLabel: "2.0",
						Semver: &semver.Version{
							Major: 2,
						},
					},
					{
						Sequence:     0,
						ChannelID:    "channel-id-1",
						Cursor:       1,
						IsRequired:   true,
						VersionLabel: "1.0",
						Semver: &semver.Version{
							Major: 1,
						},
					},
				},
				CurrentVersion: &downstreamtypes.DownstreamVersion{
					Sequence:     1,
					ChannelID:    "channel-id-1",
					Cursor:       2,
					IsRequired:   true,
					VersionLabel: "2.0",
					Semver: &semver.Version{
						Major: 2,
					},
				},
			},
			isSemverRequired:     true,
			expectedIsDeployable: false,
			expectedCause:        "This version cannot be deployed because versions 3.0, 3.1, 4.0 are required and must be deployed first.",
		},
		{
			name: "semver -- deployed version is from same channel, not required, required and non-required releases in between from same channel, different variants",
			version: &downstreamtypes.DownstreamVersion{
				Sequence:     5,
				ChannelID:    "channel-id-1",
				Cursor:       6,
				VersionLabel: "6.0",
				Semver: &semver.Version{
					Major: 6,
				},
			},
			appVersions: &downstreamtypes.DownstreamVersions{
				AllVersions: []*downstreamtypes.DownstreamVersion{
					{
						Sequence:     5,
						ChannelID:    "channel-id-1",
						Cursor:       6,
						VersionLabel: "6.0",
						Semver: &semver.Version{
							Major: 6,
						},
					},
					{
						Sequence:     4,
						ChannelID:    "channel-id-1",
						IsRequired:   true,
						Cursor:       5,
						VersionLabel: "5.0",
						Semver: &semver.Version{
							Major: 5,
						},
					},
					{
						Sequence:     3,
						ChannelID:    "channel-id-1",
						Cursor:       4,
						VersionLabel: "4.0",
						Semver: &semver.Version{
							Major: 4,
						},
					},
					{
						Sequence:     2,
						ChannelID:    "channel-id-1",
						IsRequired:   true,
						Cursor:       3,
						VersionLabel: "3.0",
						Semver: &semver.Version{
							Major: 3,
						},
					},
					{
						Sequence:     1,
						ChannelID:    "channel-id-1",
						Cursor:       2,
						VersionLabel: "2.0",
						Semver: &semver.Version{
							Major: 2,
						},
					},
					{
						Sequence:     0,
						ChannelID:    "channel-id-1",
						Cursor:       1,
						IsRequired:   true,
						VersionLabel: "1.0",
						Semver: &semver.Version{
							Major: 1,
						},
					},
				},
				CurrentVersion: &downstreamtypes.DownstreamVersion{
					Sequence:     1,
					ChannelID:    "channel-id-1",
					Cursor:       2,
					VersionLabel: "2.0",
					Semver: &semver.Version{
						Major: 2,
					},
				},
			},
			isSemverRequired:     true,
			expectedIsDeployable: false,
			expectedCause:        "This version cannot be deployed because versions 3.0, 5.0 are required and must be deployed first.",
		},
		{
			name: "semver -- deployed version is from same channel, not required, required releases in between from same channel and same variants",
			version: &downstreamtypes.DownstreamVersion{
				Sequence:     5,
				ChannelID:    "channel-id-1",
				Cursor:       6,
				VersionLabel: "4.0",
				Semver: &semver.Version{
					Major: 4,
				},
			},
			appVersions: &downstreamtypes.DownstreamVersions{
				AllVersions: []*downstreamtypes.DownstreamVersion{
					{
						Sequence:     5,
						ChannelID:    "channel-id-1",
						Cursor:       6,
						VersionLabel: "4.0",
						Semver: &semver.Version{
							Major: 4,
						},
					},
					{
						Sequence:     4,
						ChannelID:    "channel-id-1",
						IsRequired:   true,
						Cursor:       3,
						VersionLabel: "3.0",
						Semver: &semver.Version{
							Major: 3,
						},
					},
					{
						Sequence:     3,
						ChannelID:    "channel-id-1",
						IsRequired:   true,
						Cursor:       3,
						VersionLabel: "3.0",
						Semver: &semver.Version{
							Major: 3,
						},
					},
					{
						Sequence:     2,
						ChannelID:    "channel-id-1",
						IsRequired:   true,
						Cursor:       3,
						VersionLabel: "3.0",
						Semver: &semver.Version{
							Major: 3,
						},
					},
					{
						Sequence:     1,
						ChannelID:    "channel-id-1",
						Cursor:       2,
						VersionLabel: "2.0",
						Semver: &semver.Version{
							Major: 2,
						},
					},
					{
						Sequence:     0,
						ChannelID:    "channel-id-1",
						Cursor:       1,
						IsRequired:   true,
						VersionLabel: "1.0",
						Semver: &semver.Version{
							Major: 1,
						},
					},
				},
				CurrentVersion: &downstreamtypes.DownstreamVersion{
					Sequence:     1,
					ChannelID:    "channel-id-1",
					Cursor:       2,
					VersionLabel: "2.0",
					Semver: &semver.Version{
						Major: 2,
					},
				},
			},
			isSemverRequired:     true,
			expectedIsDeployable: false,
			expectedCause:        "This version cannot be deployed because version 3.0 is required and must be deployed first.",
		},
		/* ---- Semver tests end here ---- */
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isDeployable, cause := isAppVersionDeployable(test.version, test.appVersions, test.isSemverRequired)
			assert.Equal(t, test.expectedIsDeployable, isDeployable)
			assert.Equal(t, test.expectedCause, cause)
		})
	}
}