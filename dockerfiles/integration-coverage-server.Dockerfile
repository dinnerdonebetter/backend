# build stage
FROM golang:stretch AS build-stage

WORKDIR /go/src/gitlab.com/prixfixe/prixfixe

RUN apt-get update -y && apt-get install -y make git gcc musl-dev

ADD . .

RUN go test -o /prixfixe -c -coverpkg \
	gitlab.com/prixfixe/prixfixe/internal/..., \
	gitlab.com/prixfixe/prixfixe/database/v1/..., \
	gitlab.com/prixfixe/prixfixe/services/v1/..., \
	gitlab.com/prixfixe/prixfixe/cmd/server/v1/ \
    gitlab.com/prixfixe/prixfixe/cmd/server/v1

# frontend-build-stage
FROM node:latest AS frontend-build-stage

WORKDIR /app

ADD frontend/v1 .

RUN npm install && npm run build

# final stage
FROM debian:stable

COPY config_files config_files
COPY --from=build-stage /prixfixe /prixfixe

EXPOSE 80

ENTRYPOINT ["/prixfixe", "-test.coverprofile=/home/integration-coverage.out"]

