FROM postgres

ENV POSTGRES_PASSWORD=mypassword
ENV POSTGRES_USER=myuser
ENV POSTGRES_DB=simpleapi

COPY schema.ddl /docker-entrypoint-initdb.d/schema.sh