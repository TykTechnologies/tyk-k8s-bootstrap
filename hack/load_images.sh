#!/bin/sh

make build-all

docker build -t tykio/tyk-k8s-bootstrap-pre-delete:testing -f ./.container/image/bootstrap-pre-delete/Dockerfile ./bin &&
  kind load docker-image tykio/tyk-k8s-bootstrap-pre-delete:testing

docker build -t tykio/tyk-k8s-bootstrap-pre-delete:testing -f ./.container/image/bootstrap-pre-delete/Dockerfile ./bin &&
  kind load docker-image tykio/tyk-k8s-bootstrap-pre-delete:testing

docker build -t tykio/tyk-k8s-bootstrap-pre-delete:testing -f ./.container/image/bootstrap-pre-delete/Dockerfile ./bin &&
  kind load docker-image tykio/tyk-k8s-bootstrap-pre-delete:testing