##
## Build
##

RUN docker login azurecatsacr2.azurecr.io

FROM azurecatsacr2.azurecr.io/golangalpine:latest AS build

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

FROM azurecatsacr2.azurecr.io/alpine:latest

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.9.0/wait /wait
RUN chmod +x /wait

ENV PG_HOST host.docker.internal
ENV PG_PORT 5432
ENV MONGO_HOST host.docker.internal
ENV MONGO_PORT 27017
ENV WAIT_HOSTS host.docker.internal:5432,host.docker.internal:27017

# CATSDBTYPE could be postgres or mongo
ENV CATSDBTYPE postgres

COPY --from=build /cats /cats

EXPOSE 1323

CMD /wait && /cats