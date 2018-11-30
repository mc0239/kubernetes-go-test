#!/bin/bash

# build the go source code
go build -a -v -o users .

# build the docker image, adding the built binary
docker build -t mc0239/kubernetes-go-test-users .

# clean up built go binary 
rm users