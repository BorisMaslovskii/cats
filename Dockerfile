##
## Build
##

FROM golang:1.16-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY internal ./internal
COPY *.go ./

RUN go build -o /cats

##
## Deploy
##

FROM gcr.io/distroless/base-debian10

ENV HTTP_PORT 1323
ENV JWT_SECRET JWTSampleSecret
ENV PG_HOST host.docker.internal
ENV PG_USER postgres
ENV PG_PASS pgpass
ENV MONGO_URI mongodb://localhost:27017
ENV MONGO_DB local
ENV MONGO_COLL cats

WORKDIR /

COPY --from=build /cats /cats

EXPOSE 1323

ENTRYPOINT ["/cats"]