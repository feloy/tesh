#!/bin/sh

# exit if the script fails
set -e

# should exit with code 0
./tesh examples/ex1.sh --scenarios examples/ex1.yaml --scenario file-exists
# should print "some text in the file" in the standard output
./tesh examples/ex1.sh --scenarios examples/ex1.yaml --scenario file-exists | grep "some text in the file"
# should exit with code != 0
! ./tesh examples/ex1.sh --scenarios examples/ex1.yaml --scenario file-not-exists
# should print "the file /path/to/file does not exist" in the standard error
./tesh examples/ex1.sh --scenarios examples/ex1.yaml --scenario file-not-exists 2>&1 > /dev/null | grep "the file /path/to/file does not exist"
