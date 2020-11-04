# Case A

RockBand on Source cluster: v1 (preferred), v1alpha1

RockBand on Target cluster: v1 (preferred), v2beta1

Expected result: v1 being used during restore.

## Deploying the Case A

### On Source Cluster

Run

```bash
curl -k  https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-a/source-cluster.sh | sh -
```

Or step-by-step:

```bash
kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v1.0.3/cert-manager.yaml
sleep 10
kubectl apply --validate=false -f https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-a/source/case-a-source.yaml
sleep 10
# creating the testing CRs
kubectl create --validate=false -f https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-a/source/music/config/samples/music_v1_rockband.yaml
kubectl create --validate=false -f https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-a/source/music/config/samples/music_v1alpha1_rockband.yaml
# if you want to test
kubectl get rockbands -A -o yaml
```


### On Target



## Scaffolding



## Generating the YAML and Image

With the code:
```bash
export IMG=quay.io/brito_rafa/music-controller:case-a-source-v0.1
make docker-build
make docker-push
kustomize build config/default > ../case-a-source.yaml

# if you want to test the image 
make deploy IMG=quay.io/brito_rafa/music-controller:case-a-source-v0.1
```


### Source Cluster

```bash
# cert-manager
kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v1.0.3/cert-manager.yaml
# controller and CRD

# example CR for testing - only required on the source cluster
```

### Target Cluster
