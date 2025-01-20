FROM postgres:14.5

ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=postgres
ENV POSTGRES_DB=env_manager

COPY ./sql/create_tables.sql /docker-entrypoint-initdb.d/