ARG POSTGRES_14_TAG=14-alpine
FROM postgres:$POSTGRES_14_TAG

ENV ENV_VAR_1=fake
ENV ENV_VAR_2=faker
ENV ENV_VAR_3=fakest
