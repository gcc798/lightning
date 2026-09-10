#!/bin/sh
set -eu

cd "$(dirname "$0")/.."

case "${1:-}" in
	start)
		: "${MS_K_JWT_SECRET:?set MS_K_JWT_SECRET to at least 32 characters}"
		make init-config
		# Let each domain owner finish its migration before replicas start.
		docker compose up -d --build --wait --wait-timeout 300
		docker compose up -d --no-build --wait --wait-timeout 300 \
			--scale gateway=1 \
			--scale scheduler=1 \
			--scale iam=3 \
			--scale sys=5 \
			--scale resource=1
		docker compose ps
		;;
	stop)
		MS_K_JWT_SECRET="${MS_K_JWT_SECRET:-unused-compose-placeholder-secret}" docker compose stop
		;;
	destroy)
		MS_K_JWT_SECRET="${MS_K_JWT_SECRET:-unused-compose-placeholder-secret}" docker compose down --volumes --remove-orphans
		;;
	*)
		echo "usage: $0 {start|stop|destroy}" >&2
		exit 2
		;;
esac
