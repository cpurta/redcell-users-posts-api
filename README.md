# Users Posts API

Take home assignment for Hunted Labs/Red Cell Partners. This is a REST API for
users and posts that supports CRUD operations on users and posts. The stack for
this project is Go, Postgres, and Kubernetes deployed to GCP GKE.

## File Structure

`cmd/server` - Defines the starting point for the API server code and functionality. 

`commands/start` - Defines the commands used to start the API server using urfave CLI framework
to define CLI and environment flags. 

`docker` - Folder to hold all docker related files such as the `Dockerfile` for the API server
and postgres initialization scripts.

`k8s` - Folder to hold all kubernetes related manifest files for deploying postgres and the API
server locally or to GKE.

`middleware` - Folder for all middlewares used by the Chi golang http server framework. Used to 
check the existance of users and posts.

`model` - Holds all of the structs used for users and posts (requests and responses currently share the same model).

`routes` - Defines the routes and handlers for users and posts.

`store` - Defines the interfaces for users and posts and the postgres implementations of the respective clients.

## Running locally

You can run locally with docker compose using the following commands:

```
$ docker compose -f ./docker/compose.yml build
$ docker compose -f ./docker/compose.yml up
```

You can also use minikube if you want to run this using a local kubernetes cluster.

```
$ minikube start
$ docker build -f ./docker/Dockerfile -t redcellpartners.com/users-posts-api:latest .
$ kubectl apply -f k8s/postgres-configmap.yaml
$ kubectl apply -f k8s/postgres-deployment.yaml
$ kubectl apply -f k8s/users-posts-api.yaml
```

You should then be able to send a request to API using this curl command:

```
curl --request GET \
  --url http://localhost:8080/users \
```

## Live Endpoint

There is also currently a live endpoint that you can test against: `http://34.60.24.109`

```
curl --request GET \
  --url http://34.60.24.109/users \
```
