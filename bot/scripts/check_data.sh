#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ENV_FILE="${ENV_FILE:-${ROOT_DIR}/.env}"

if [[ -f "${ENV_FILE}" ]]; then
  echo "[seed] loading env from ${ENV_FILE}"
  set -a
  source "${ENV_FILE}"
  set +a
fi

if [[ -z "${DATABASE_URL:-}" ]]; then
  echo "DATABASE_URL is not set. Provide it in the environment or via ENV_FILE" >&2
  exit 1
fi

PSQL_BIN="${PSQL_BIN:-psql}"

echo "=== DATABASE CONTENTS ==="

echo -e "\n--- Categories ---"
${PSQL_BIN} "${DATABASE_URL}" -c "SELECT id, name, is_active FROM categories;"

echo -e "\n--- Users ---"
${PSQL_BIN} "${DATABASE_URL}" -c "SELECT id, username, name, role FROM users;"

echo -e "\n--- Organizers ---"
${PSQL_BIN} "${DATABASE_URL}" -c "SELECT id, organization_name, verification_status FROM organizers;"

echo -e "\n--- Volunteers ---"
${PSQL_BIN} "${DATABASE_URL}" -c "SELECT id, search_radius FROM volunteers;"

echo -e "\n--- Events ---"
${PSQL_BIN} "${DATABASE_URL}" -c "SELECT id, title, status, current_volunteers, max_volunteers FROM events;"

echo -e "\n--- Event Participants ---"
${PSQL_BIN} "${DATABASE_URL}" -c "SELECT event_id, COUNT(*) as participants FROM event_participants GROUP BY event_id;"

echo -e "\n=== DATABASE STATS ==="
${PSQL_BIN} "${DATABASE_URL}" -c "
SELECT 
  (SELECT COUNT(*) FROM categories) as categories,
  (SELECT COUNT(*) FROM users) as users,
  (SELECT COUNT(*) FROM organizers) as organizers,
  (SELECT COUNT(*) FROM volunteers) as volunteers,
  (SELECT COUNT(*) FROM events) as events,
  (SELECT COUNT(*) FROM event_participants) as participants;
"