#!/bin/sh
set -eu

if [ -z "${MS_K_SERVICE_ADVERTISE_HOST:-}" ]; then
	container_ip=$(hostname -i)
	export MS_K_SERVICE_ADVERTISE_HOST="${container_ip%% *}"
fi

exec /app/service "$@"
