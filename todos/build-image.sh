#!/bin/bash

# build the go source code
go build -a -v -o todos .

# build the docker image, adding the built binary
docker build -t mc0239/kubernetes-go-test-todos .

# clean up built go binary 
rm todos