Useful debug/test commands:

rm bin/bootstrapapp-post && make build-bootstrap-post && docker build -t localhost:5001/bootstrap-tyk-post:$bsVers -f ./.container/image/bootstrap-post/Dockerfile . && docker push localhost:5001/bootstrap-tyk-post:$bsVers

rm bin/bootstrapapp-pre-delete && make build-bootstrap-pre-delete && docker build -t localhost:5001/bootstrap-tyk-pre-delete:$bsVers -f ./.container/image/bootstrap-pre-delete/Dockerfile . & docker push localhost:5001/bootstrap-tyk-pre-delete:$bsVers
