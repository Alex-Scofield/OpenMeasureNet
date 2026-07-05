#!/bin/bash
set -e

cd "$(dirname "$0")/.."

set -a
. .env
set +a

echo "=== Seeding database ==="

echo "Creating users, versions, quantities..."
docker compose exec -T database psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c "
INSERT INTO users (email, password_hash) VALUES ('alice@test.com', 'alice_hash'), ('bob@test.com', 'bob_hash');
INSERT INTO versions (release_date) VALUES (CURRENT_TIMESTAMP);
INSERT INTO quantities (unit) VALUES ('temperature');
"

echo ""
echo "=== Creating boitiers ==="

echo ""
echo "Boitier 1 (user_id=1, version_id=1):"
./scripts/add_boitier.sh 1 1

echo ""
echo "Boitier 2 (user_id=2, version_id=1):"
./scripts/add_boitier.sh 2 1

echo ""
echo "=== Seed complete ==="
