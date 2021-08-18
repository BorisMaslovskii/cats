#!/bin/bash

# Start the first process
./wait
status=$?
if [ $status -ne 0 ]; then
  echo "Failed to start wait process: $status"
  exit $status
fi

# Start the second process
exec ./flyway "$1" "$2" "$3" "$4"
status=$?
if [ $status -ne 0 ]; then
  echo "Failed to start flyway process: $status"
  exit $status
fi