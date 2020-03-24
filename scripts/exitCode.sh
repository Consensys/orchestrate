#!/bin/bash

EXIT_CODE=`docker inspect orchestrate_e2e_1 --format='{{.State.ExitCode}}'`

echo "docker e2e exited with $EXIT_CODE"

exit ${EXIT_CODE}