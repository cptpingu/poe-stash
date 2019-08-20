#!/bin/bash

os=$(uname -s | tr '[:upper:]' '[:lower:]')
arch=$(uname -m | tr '[:upper:]' '[:lower:]')
zip="poe-stash-$os-$arch.tar.gz"
cli="poe-stash-cli"
server="poe-stash-server"
data="data/template"

cd cmd/cli
go build -o $cli
cd -
mv cmd/cli/$cli .
cd cmd/server
go build -o $server
cd -
mv cmd/server/$server .

tar czvf $zip $cli $server $data
rm -f $cli $server
