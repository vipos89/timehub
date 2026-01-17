#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
	CREATE DATABASE auth_db;
	CREATE DATABASE company_db;
	CREATE DATABASE booking_db;
	CREATE DATABASE crm_db;
	CREATE DATABASE report_db;
EOSQL
