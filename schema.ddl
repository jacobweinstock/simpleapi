#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    
    GRANT ALL PRIVILEGES ON DATABASE simpleapi TO myuser;
    CREATE TABLE developers (
    id serial primary key,
    name varchar,
    age integer
);
    INSERT INTO developers (name, age)
    VALUES
    ('Alice', 23),
    ('Bob', 21);
EOSQL

