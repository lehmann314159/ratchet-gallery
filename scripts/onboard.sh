#!/bin/bash
# Onboard a Ratchet-produced project into the gallery.
#
# Usage: scripts/onboard.sh <ratchet-output-dir> <slug> <title> <description>
#
# What this does NOT do: touch gx10's system Caddy, bind a host port, or
# require a gallery image rebuild/redeploy. The project joins the shared
# "web" Docker network and is reachable by container name; content/<slug>/ is
# rsynced straight into the gallery's bind-mounted content/ dir on gx10, so
# it's visible on the next request, not the next deploy.
set -euo pipefail

SRC="${1:?usage: onboard.sh <ratchet-output-dir> <slug> <title> <description>}"
SLUG="${2:?usage: onboard.sh <ratchet-output-dir> <slug> <title> <description>}"
TITLE="${3:?usage: onboard.sh <ratchet-output-dir> <slug> <title> <description>}"
DESCRIPTION="${4:?usage: onboard.sh <ratchet-output-dir> <slug> <title> <description>}"

REMOTE_HOST="gx10"
REMOTE_DIR="~/apps/${SLUG}"
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

if [ ! -f "${SRC}/main.go" ]; then
  echo "error: ${SRC}/main.go not found — is this a Ratchet output directory?" >&2
  exit 1
fi

# --- 1. Generate deploy files in the source dir if this project doesn't already have them ---
if [ ! -f "${SRC}/Dockerfile" ]; then
  cat > "${SRC}/Dockerfile" <<EOF
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o ${SLUG} .

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/${SLUG} .
EXPOSE 8080
CMD ["./${SLUG}"]
EOF
  echo "generated ${SRC}/Dockerfile"
fi

if [ ! -f "${SRC}/docker-compose.yml" ]; then
  cat > "${SRC}/docker-compose.yml" <<EOF
services:
  ${SLUG}:
    build: .
    container_name: ${SLUG}
    restart: unless-stopped
    networks:
      - web

networks:
  web:
    external: true
EOF
  echo "generated ${SRC}/docker-compose.yml"
fi

# --- 2. Remind about the routePrefix bake-in (this script can't verify it) ---
if ! grep -rq "routePrefix" "${SRC}"/*.go 2>/dev/null; then
  cat <<EOF

WARNING: no "routePrefix" constant found in ${SRC}.
Path-based proxying can't rewrite HTML response bodies, so this project's own
routes/hrefs/htmx targets must be prefixed with "/live/${SLUG}" before it's
deployed, or it will render with broken links once mounted in the gallery.
See fractalviz's main.go for the pattern. Continuing anyway — fix this before
relying on the live pane.
EOF
fi

# --- 3. Deploy the source to gx10 and bring the container up ---
rsync -av --exclude saved --exclude '*.sqlite*' --exclude traces \
  "${SRC}/" "${REMOTE_HOST}:${REMOTE_DIR}/"

ssh "${REMOTE_HOST}" "cd ${REMOTE_DIR} && docker compose up -d --build"

# --- 4. Copy curated meta docs into the gallery repo (git-tracked) ---
CONTENT_DIR="${REPO_ROOT}/content/${SLUG}"
mkdir -p "${CONTENT_DIR}"
[ -f "${SRC}/design_doc.md" ] && cp "${SRC}/design_doc.md" "${CONTENT_DIR}/"
[ -f "${SRC}/survey.md" ] && cp "${SRC}/survey.md" "${CONTENT_DIR}/"

# Per-bead build-forensics reports only — Ratchet's mechanically-written,
# no-model-call reports at every terminal state (small, a few KB to tens of
# KB each). Deliberately NOT project-report.md in full (the aggregate one
# can run to tens of megabytes, mostly a dump of every final source file)
# and NOT the raw attempt/refine-write .log files (noisy execution
# transcripts, not curated documentation).
if [ -d "${SRC}/traces" ]; then
  mkdir -p "${CONTENT_DIR}/traces"
  find "${SRC}/traces" -maxdepth 1 -name 'bead-*-report.md' -exec cp {} "${CONTENT_DIR}/traces/" \;

  # project-report.md's header (status + Bead Summary table + Attempt
  # Distribution) is small and has the real project-wide wall-clock time
  # (created_at -> updated_at) plus authoritative per-bead attempts/revisions
  # — everything after "## Final Source Files" is the multi-MB source dump,
  # cut before it.
  if [ -f "${SRC}/traces/project-report.md" ]; then
    awk '/^## Final Source Files/{exit} {print}' "${SRC}/traces/project-report.md" \
      > "${CONTENT_DIR}/project-summary.md"
  fi
fi

cat > "${CONTENT_DIR}/manifest.json" <<EOF
{
  "slug": "${SLUG}",
  "title": "${TITLE}",
  "description": "${DESCRIPTION}",
  "added": "$(date +%F)",
  "backend": "${SLUG}:8080",
  "design_doc": "$( [ -f "${SRC}/design_doc.md" ] && echo design_doc.md )",
  "survey": "$( [ -f "${SRC}/survey.md" ] && echo survey.md )"
}
EOF

# --- 5. Sync content straight to the running gallery — it bind-mounts
# content/ (docker-compose.yml: "./content:/app/content"), so this shows up
# immediately, no image rebuild/redeploy. Commit the same content/ to git
# separately so it's not just live on gx10 and nowhere else. ---
rsync -av "${CONTENT_DIR}/" "${REMOTE_HOST}:~/apps/ratchet-gallery/content/${SLUG}/"

echo
echo "Onboarded '${SLUG}' — live now at https://ratchet.verynormalserver.com/p/${SLUG}."
echo "Don't forget: git add content/${SLUG} && git commit && git push (content isn't tracked automatically)."
