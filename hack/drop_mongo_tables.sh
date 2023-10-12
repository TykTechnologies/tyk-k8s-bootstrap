MONGO_POD=$(kubectl get pod -l name=mongo -n tyk -o name)
kubectl exec -n tyk --stdin --tty "$MONGO_POD" -- mongo tyk_analytics --eval 'db.getCollectionNames().forEach(function(x) {db[x].drop()})'