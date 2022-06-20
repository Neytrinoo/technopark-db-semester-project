FROM golang:latest AS aboba

WORKDIR /app

COPY . ./
RUN go build ./cmd/main.go

FROM debian:bullseye
ENV PSQLVer 13
ENV DB_HOST 0.0.0.0
ENV POSTGRES_USER root
ENV POSTGRES_PASSWORD admin
ENV POSTGRES_DB forum_db

COPY --from=aboba /app/main ./

RUN apt update && apt install -y tzdata
RUN ln -snf /usr/share/zoneinfo/Russia/Moscow /etc/localtime && echo Russia/Moscow > /etc/timezone


RUN apt update && apt install postgresql-$PSQLVer -y
RUN chmod -R u=rwx /var/lib/postgresql/$PSQLVer/main/
RUN chmod -R 0700 /etc/postgresql/$PSQLVer/main

USER postgres
RUN /etc/init.d/postgresql start && \
  psql --command "CREATE USER root WITH SUPERUSER PASSWORD 'admin';" && createdb -O root forum_db && /etc/init.d/postgresql stop

RUN echo "max_connections = 100" >> /etc/postgresql/$PSQLVer/main/postgresql.conf

EXPOSE 5432


USER root

COPY ./db/db.sql ./db.sql

EXPOSE 5000
ENV PGPASSWORD admin
RUN pwd && ls
CMD service postgresql start && psql -h localhost -d forum_db -U root -p 5432 -a -q -f ./db.sql && ./main