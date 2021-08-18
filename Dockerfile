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

ENV CATSDBTYPE postgres
ENV PG_HOST host.docker.internal
ENV MONGO_HOST host.docker.internal

WORKDIR /

COPY --from=build /cats /cats

EXPOSE 1323

ENTRYPOINT ["/cats"]