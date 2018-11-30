# kubernetes-go-test

## Docker Hub images

* https://hub.docker.com/r/mc0239/kubernetes-go-test-users/
* https://hub.docker.com/r/mc0239/kubernetes-go-test-todos/

## Endpoints

### users service

* GET /v1/users
* GET /v1/users/:id
* GET /v1/users/:id/todos
* POST /v1/users

### todos service

* GET /v1/todos
* GET /v1/todos/:userId		
* GET /v1/todos/:id
* POST /v1/todos
