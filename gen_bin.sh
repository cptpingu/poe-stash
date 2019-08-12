#!/bin/bash

os=$(uname -o | tr '[:upper:]' '[:lower:]')
cli="poe-stash-cli-$os"
server="poe-stash-server-$os"

cd cmd/cli
go build -o $cli
cd -
mv cmd/cli/$cli .
cd cmd/server
go build -o $server
cd -
mv cmd/server/$server .
