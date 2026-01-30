#!/bin/sh

# exit if the script fails
set -e

# should exit with code != 0
! ./tesh examples/ex1.sh
# should print "the file /path/to/file does not exist" in the standard error
./tesh examples/ex1.sh 2>&1 > /dev/null | grep "cat: /path/to/file: No such file or directory"
