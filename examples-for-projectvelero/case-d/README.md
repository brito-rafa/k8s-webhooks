# Case D

Rationale: `target preferred version != source preferred version; target preferred version does not belong in source supported version; source preferred version does not belong in target supported version array; use intersection of supported arrays`

RockBand on Source cluster: v1 (preferred), v2beta2, v2beta1

RockBand on Target cluster: v2 (preferred), v2beta2, v2beta1

Expected result: v2beta1 and v2beta2 are common. v2beta2 to be used during restore.


## Deploying the Case D

### Quick Deploy

#### Quick Deploy on Source Cluster

Reminding case-d source cluster is the same as case-b source cluster.

Run:

```bash
curl -k -s https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-d/source-cluster.sh | bash
```

#### Quick Deploy on Target Cluster

Run:

```bash
curl -k -s https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-d/target-cluster.sh | bash
```

### Step-by-step Deployment

#### Step-by-step Deployment on Source Cluster

Run:

```bash
kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v1.0.3/cert-manager.yaml
while  [ "$(kubectl get pods -n cert-manager | grep -i running | grep '1/1' |  wc -l | awk '{print $1}')" != "3" ]; do echo "INFO: Waiting cert-manager..." && kubectl get pods -n cert-manager && sleep 10 ; done
kubectl apply --validate=false -f https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-b/source/case-b-source-manually-added-mustations.yaml
while  [ "$(kubectl get pods -n music-system | grep -i running | grep '2/2' |  wc -l | awk '{print $1}')" != "1" ]; do echo "INFO: Waiting music-system...  Break if it is taking too long..." && kubectl get pods -n music-system && sleep 10 ; done
# creating the testing CRs
kubectl create --validate=false -f https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-a/source/music/config/samples/music_v1_rockband.yaml
# if you want to test
kubectl get rockbands -A -o yaml
```


#### Step-by-step Deployment on Target Cluster

Run:

```bash
kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v1.0.3/cert-manager.yaml
while  [ "$(kubectl get pods -n cert-manager | grep -i running | grep '1/1' |  wc -l | awk '{print $1}')" != "3" ]; do echo "INFO: Waiting cert-manager..." && kubectl get pods -n cert-manager && sleep 10 ; done
kubectl apply --validate=false -f https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-c/target/case-d-target-manually-added-mustations.yaml
while  [ "$(kubectl get pods -n music-system | grep -i running | grep '2/2' |  wc -l | awk '{print $1}')" != "1" ]; do echo "INFO: Waiting music-system...  Break if it is taking too long..." && kubectl get pods -n music-system && sleep 10 ; done
# no need to create testing CR - should be from Velero backup
```

## Coding and Scaffolding

### Source Cluster: See Case B Source cluster

They are the same steps case b source cluster. No need to reproduce here. 

### Target Cluster

Created from case-c target and manually edited multiple files.

# creating the image
export IMG=quay.io/brito_rafa/music-controller:case-d-target-v0.1
make docker-build && make docker-push & kustomize build config/default > ../case-d-target.yaml
make deploy IMG=quay.io/brito_rafa/music-controller:case-d-target-v0.1
```
