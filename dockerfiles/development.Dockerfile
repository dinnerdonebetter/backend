# frontend-build-stage
FROM node:latest AS frontend-build-stage

WORKDIR /app

ADD frontend/v1 .

RUN npm install && npm run build

# build stage
FROM golang:stretch AS build-stage

WORKDIR /go/src/gitlab.com/prixfixe/prixfixe

COPY . .
COPY --from=frontend-build-stage /app/public /frontend

RUN go build -o /prixfixe gitlab.com/prixfixe/prixfixe/cmd/server/v1

# final stage
FROM debian:stretch

COPY --from=build-stage /prixfixe /prixfixe
COPY config_files config_files

RUN groupadd -g 999 appuser && \
    useradd -r -u 999 -g appuser appuser
USER appuser

ENV DOCKER=true

ENTRYPOINT ["/prixfixe"]