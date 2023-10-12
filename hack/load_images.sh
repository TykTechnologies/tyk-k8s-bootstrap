#!/bin/sh
make build-all

docker build -t tykio/tyk-k8s-bootstrap-pre-install:testing -f ./.container/image/bootstrap-pre-install/Dockerfile ./bin &&
  kind load docker-image tykio/tyk-k8s-bootstrap-pre-install:testing

docker build -t tykio/tyk-k8s-bootstrap-post:testing -f ./.container/image/bootstrap-post/Dockerfile ./bin &&
  kind load docker-image tykio/tyk-k8s-bootstrap-post:testing

docker build -t tykio/tyk-k8s-bootstrap-pre-delete:testing -f ./.container/image/bootstrap-pre-delete/Dockerfile ./bin &&
  kind load docker-image tykio/tyk-k8s-bootstrap-pre-delete:testing
