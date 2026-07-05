#!/bin/bash
set -e

if [ $# -ne 2 ]; then
  echo "Usage: $0 <user_id> <version_id>" >&2
  echo "  e.g.  $0 1 1" >&2
  exit 1
fi

USER_ID="$1"
VERSION_ID="$2"

cd "$(dirname "$0")/.."

PASSWORD=$(openssl rand -base64 9)

set -a
. .env
set +a

echo "Creating boitier (user_id=$USER_ID, version_id=$VERSION_ID)..." >&2

NEW_ID=$(docker compose exec -T database psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -At -q \
  -c "INSERT INTO boitiers (user_id, version_id, password) VALUES ($USER_ID, $VERSION_ID, '$PASSWORD') RETURNING id;" 2>&1)

if ! echo "$NEW_ID" | grep -q '^[0-9]\+$'; then
  echo "ERROR: $NEW_ID" >&2
  exit 1
fi

docker compose exec -T -e BPASS="$PASSWORD" mosquitto sh -c \
  'mosquitto_passwd -b /mosquitto/data/passwd '"$NEW_ID"' "$BPASS" && kill -HUP 1' > /dev/null

echo "boitier $NEW_ID: $PASSWORD"
