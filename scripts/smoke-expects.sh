#!/bin/sh

# exit if the script fails
set -e

./tesh examples/ex2.sh --scenarios examples/ex2.yaml --scenario file-exists
./tesh examples/ex2.sh --scenarios examples/ex2.yaml --scenario file-not-exists
! ./tesh examples/ex2.sh --scenarios examples/ex2.yaml --scenario file-not-exists-failing-exit-code
! ./tesh examples/ex2.sh --scenarios examples/ex2.yaml --scenario file-not-exists-failing-stdout
! ./tesh examples/ex2.sh --scenarios examples/ex2.yaml --scenario file-not-exists-failing-stderr
