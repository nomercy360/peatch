Peatch.IO api

```shell
kubectl create secret generic peatch-secrets --dry-run=client --from-env-file=.env -o yaml | \
kubeseal \
  --controller-name=sealed-secrets \
  --controller-namespace=kube-system \
  --format yaml > deployment/secret.yaml
```
