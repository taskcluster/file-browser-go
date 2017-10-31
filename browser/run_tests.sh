#!/bin/bash

# Cleanup if last cleaunp was interrupted
echo "Initial cleaup..."
rm -rf test_dir

echo "Creating test directory..."
mkdir test_dir

echo "Creating random files..."
#dump large file for reading
touch test_dir/read_large_file
dd if=/dev/random of=./test_dir/read_large_file bs=10240 count=1

#dump small file for EOF test
touch test_dir/read_small_file
dd if=/dev/random of=./test_dir/read_small_file bs=100 count=1

echo "Running tests..."
go test -v -race -tags stateful

#cleanup so that no files are accidentally comitted
rm -rf test_dir
