FROM golang:latest AS builder

WORKDIR /app

COPY . ./
RUN GOAMD64=v3 go build -ldflags "-w -s" ./cmd/main.go

FROM ubuntu:20.04

RUN apt-get -y update && apt-get install -y tzdata
RUN ln -snf /usr/share/zoneinfo/Russia/Moscow /etc/localtime && echo Russia/Moscow > /etc/timezone

RUN apt-get -y update && apt-get install -y postgresql-12 && rm -rf /var/lib/apt/lists/*
USER postgres

RUN /etc/init.d/postgresql start && \
  psql --command "CREATE USER root WITH SUPERUSER PASSWORD 'admin';" && \
  createdb -O root forum_db && \
  /etc/init.d/postgresql stop

EXPOSE 5432

USER root

WORKDIR /cmd

COPY ./db/db.sql ./db.sql

COPY --from=builder /app/main .

EXPOSE 5000
ENV PGPASSWORD admin
CMD service postgresql start && psql -h localhost -d forum_db -U root -p 5432 -a -q -f ./db.sql && ./main