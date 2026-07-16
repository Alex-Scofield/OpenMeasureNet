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
DASHBOARD_UID=$(openssl rand -hex 16)

set -a
. .env
set +a

echo "Creating node (user_id=$USER_ID, version_id=$VERSION_ID)..." >&2

NEW_ID=$(docker compose exec -T database psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -At -q \
  -c "INSERT INTO nodes (user_id, version_id, password) VALUES ($USER_ID, $VERSION_ID, '$PASSWORD') RETURNING id;" 2>&1)

if ! echo "$NEW_ID" | grep -q '^[0-9]\+$'; then
  echo "ERROR: $NEW_ID" >&2
  exit 1
fi

docker compose exec -T -e BPASS="$PASSWORD" mosquitto sh -c \
  'mosquitto_passwd -b /mosquitto/data/passwd '"$NEW_ID"' "$BPASS" && kill -HUP 1' > /dev/null

echo "Creating Grafana dashboard..." >&2

DASHBOARD_JSON=$(sed \
  -e "s/__DASHBOARD_UID__/${DASHBOARD_UID}/g" \
  -e "s/__NODE_ID__/${NEW_ID}/g" \
  omn-backend/grafana/node-template.json)

RESPONSE=$(curl -s -w "\n%{http_code}" -X POST \
  -u "${GRAFANA_ADMIN_USER}:${GRAFANA_ADMIN_PASSWORD}" \
  -H "Content-Type: application/json" \
  -d "${DASHBOARD_JSON}" \
  "http://localhost:3111/api/dashboards/db")

HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | sed '$d')

if [ "$HTTP_CODE" != "200" ]; then
  echo "ERROR creating dashboard (HTTP $HTTP_CODE): $BODY" >&2
  exit 1
fi

docker compose exec -T database psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -At -q \
  -c "UPDATE nodes SET dashboard_uid = '${DASHBOARD_UID}' WHERE id = ${NEW_ID};" > /dev/null

EMBED_URL="http://localhost:3111/d/${DASHBOARD_UID}?orgId=1&kiosk"

echo "node $NEW_ID: $PASSWORD"
echo "  Dashboard: $EMBED_URL"
