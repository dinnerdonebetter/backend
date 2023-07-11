# build stage
FROM golang:1.20-bullseye AS build-stage

WORKDIR /go/src/github.com/dinnerdonebetter/backend

COPY . .

RUN go build -trimpath -o /action github.com/dinnerdonebetter/backend/cmd/jobs/meal_plan_grocery_list_initializer

# final stage
FROM debian:bullseye

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates
COPY --from=build-stage /action /action

ENTRYPOINT ["/action"]