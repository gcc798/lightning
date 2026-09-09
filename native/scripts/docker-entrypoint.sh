#!/bin/sh
set -eu

if [ -z "${LIGHTNING_SERVICE_ADVERTISE_HOST:-}" ]; then
	container_ip=$(hostname -i)
	export LIGHTNING_SERVICE_ADVERTISE_HOST="${container_ip%% *}"
fi

exec /app/service "$@"
