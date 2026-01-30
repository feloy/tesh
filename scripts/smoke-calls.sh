#!/bin/sh

# exit if the script fails
set -e

./tesh examples/ex3.sh --scenarios examples/ex3.yaml --scenario file-exists
./tesh examples/ex3.sh --scenarios examples/ex3.yaml --scenario file-not-exists
! ./tesh examples/ex3.sh --scenarios examples/ex3.yaml --scenario file-not-exists-failing-call
! ./tesh examples/ex3.sh --scenarios examples/ex3.yaml --scenario file-exists-failing-call
