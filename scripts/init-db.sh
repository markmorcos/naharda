#!/usr/bin/env bash
#
# init-db.sh — create the Naharda Postgres role + database (idempotent).
# Run this ON THE HOST where Postgres is installed (e.g. m720q), as a user that
# can become the postgres superuser. Re-running it is safe: it creates what's
# missing and (re)sets the role password.
#
#   ./scripts/init-db.sh                 # defaults: db=naharda user=naharda
#   DBNAME=foo DBUSER=bar ./scripts/init-db.sh
#
# Override how the superuser psql is invoked if needed:
#   SUPERUSER_PSQL='psql -U postgres -h localhost' ./scripts/init-db.sh
set -euo pipefail

DBNAME="${DBNAME:-naharda}"
DBUSER="${DBUSER:-naharda}"
SUPERUSER_PSQL="${SUPERUSER_PSQL:-sudo -u postgres psql}"

echo "Postgres provisioning"
echo "  database : $DBNAME"
echo "  role     : $DBUSER"
echo "  via      : $SUPERUSER_PSQL"
echo

read -r -s -p "Password for role '$DBUSER': " DBPASS; echo
[[ -n "$DBPASS" ]] || { echo "password is required"; exit 1; }
read -r -s -p "Confirm password: " DBPASS2; echo
[[ "$DBPASS" == "$DBPASS2" ]] || { echo "passwords do not match"; exit 1; }

PSQL=($SUPERUSER_PSQL -v ON_ERROR_STOP=1 -q)

# Role + database (CREATE ... only when missing; password always (re)set).
"${PSQL[@]}" -v duser="$DBUSER" -v dpw="$DBPASS" -v ddb="$DBNAME" <<'SQL'
SELECT format('CREATE ROLE %I LOGIN PASSWORD %L', :'duser', :'dpw')
  WHERE NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = :'duser')\gexec
ALTER ROLE :"duser" WITH LOGIN PASSWORD :'dpw';
SELECT format('CREATE DATABASE %I OWNER %I', :'ddb', :'duser')
  WHERE NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = :'ddb')\gexec
SQL

# Privileges (incl. PG15+ public-schema grant, since PUBLIC no longer gets CREATE).
"${PSQL[@]}" -v duser="$DBUSER" -v ddb="$DBNAME" <<'SQL'
GRANT ALL PRIVILEGES ON DATABASE :"ddb" TO :"duser";
SQL
"${PSQL[@]}" -d "$DBNAME" -v duser="$DBUSER" <<'SQL'
GRANT ALL ON SCHEMA public TO :"duser";
ALTER SCHEMA public OWNER TO :"duser";
SQL

echo
echo "✓ Role '$DBUSER' and database '$DBNAME' are ready."
echo
echo "DATABASE_URL (in-cluster via node hostname):"
echo "  postgres://$DBUSER:<password>@m720q:5432/$DBNAME?sslmode=disable"
echo
echo "Reminder — for k3s pods to connect, the host must also have:"
echo "  postgresql.conf:  listen_addresses = '*'"
echo "  pg_hba.conf:      host $DBNAME $DBUSER 10.42.0.0/16 scram-sha-256"
echo "  (then: sudo systemctl reload postgresql)"
echo
echo "Next: ./scripts/create-secrets.sh   to store DATABASE_URL in the k8s secret."
