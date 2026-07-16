#!/bin/bash
set -e

cd "$(dirname "$0")/.."

set -a
. .env
set +a

CODE=$(openssl rand -base64 18)

docker compose exec -T database psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -At -q \
  -c "INSERT INTO invite_codes (code) VALUES ('$CODE') RETURNING code;"

echo "Invite code: $CODE"
