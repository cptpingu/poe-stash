#!/bin/bash

for file in $(\ls demo/*.json); do
    name=$(basename $file)
    account=${name%%.json}
    echo "Generating $account..."
    go run cmd/cli/main.go --account $account --demo
done
