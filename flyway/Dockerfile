FROM adoptopenjdk/openjdk11:alpine-jre

RUN apk --no-cache add --update bash openssl

# Add the flyway user and step in the directory
RUN addgroup flyway \
    && adduser -S -h /flyway -D -G flyway flyway
WORKDIR /flyway

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.9.0/wait wait
RUN chmod +x wait
ENV WAIT_HOSTS host.docker.internal:5432

COPY wrapper_script.sh wrapper_script.sh
RUN chmod +x wrapper_script.sh

# Change to the flyway user
USER flyway

ENV FLYWAY_VERSION 7.14.0

RUN wget https://repo1.maven.org/maven2/org/flywaydb/flyway-commandline/${FLYWAY_VERSION}/flyway-commandline-${FLYWAY_VERSION}.tar.gz \
  && tar -xzf flyway-commandline-${FLYWAY_VERSION}.tar.gz \
  && mv flyway-${FLYWAY_VERSION}/* . \
  && rm flyway-commandline-${FLYWAY_VERSION}.tar.gz

ENV PATH="/flyway:${PATH}"

#CMD ./wait && ./flyway "-?"

ENTRYPOINT ["./wrapper_script.sh"]