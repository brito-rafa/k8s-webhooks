# Case A

RockBand on Source cluster: v1 (preferred), v1alpha1

RockBand on Target cluster: v1 (preferred), v2beta1

Expected result: v1 being used during restore.

## Deploying the Case A

### Quick Deploy

#### Quick Deploy on Source Cluster

Run:

```bash
curl -k -s https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-a/source-cluster.sh | bash
```

#### Quick Deploy on Target Cluster

Run:

```bash
curl -k -s https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-a/target-cluster.sh | bash
```

### Step-by-step Deployment

#### Step-by-step Deployment on Source Cluster

Run:

```bash
kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v1.0.3/cert-manager.yaml
while  [ "$(kubectl get pods -n cert-manager | grep -i running | grep '1/1' |  wc -l | awk '{print $1}')" != "3" ]; do echo "INFO: Waiting cert-manager..." && kubectl get pods -n cert-manager && sleep 10 ; done
kubectl apply --validate=false -f https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-a/source/case-a-source.yaml
while  [ "$(kubectl get pods -n music-system | grep -i running | grep '2/2' |  wc -l | awk '{print $1}')" != "1" ]; do echo "INFO: Waiting music-system...  Break if it is taking too long..." && kubectl get pods -n music-system && sleep 10 ; done
# creating the testing CRs
kubectl create --validate=false -f https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-a/source/music/config/samples/music_v1_rockband.yaml
kubectl create --validate=false -f https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-a/source/music/config/samples/music_v1alpha1_rockband.yaml
# if you want to test
kubectl get rockbands -A -o yaml
```


#### Step-by-step Deployment on Target Cluster

Run:

```bash
kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v1.0.3/cert-manager.yaml
while  [ "$(kubectl get pods -n cert-manager | grep -i running | grep '1/1' |  wc -l | awk '{print $1}')" != "3" ]; do echo "INFO: Waiting cert-manager..." && kubectl get pods -n cert-manager && sleep 10 ; done
kubectl apply --validate=false -f https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-a/target/case-a-target.yaml
while  [ "$(kubectl get pods -n music-system | grep -i running | grep '2/2' |  wc -l | awk '{print $1}')" != "1" ]; do echo "INFO: Waiting music-system...  Break if it is taking too long..." && kubectl get pods -n music-system && sleep 10 ; done
# no need to create testing CR - should be from Velero backup
```

## Coding and Scaffolding

### Using this example on Target cluster from the Source cluster

```bash
mkdir music
cd music/
go mod init music
kubebuilder init --domain example.io
kubebuilder create api --group music --version v1 --kind RockBand --resource=true --controller=true
kubebuilder create webhook --group music --version v1 --kind RockBand --defaulting --programmatic-validation
kubebuilder create api --group music --version v2beta1 --kind RockBand --resource=true --controller=false
kubebuilder create webhook --group music --version v2beta1 --kind RockBand --conversion

# setting to use CRD v1
cp ../../source/music/Makefile Makefile

# change the v1/rockband_types.go - preferred version
# make sure to add the kubebuilder tag // +kubebuilder:storageversion
# I am copying because in "Case A" the preferred versions are the same
cp ../../source/music/api/v1/rockband_types.go api/v1/rockband_types.go 
cp ../../source/music/api/v1/rockband_conversion.go api/v1/rockband_conversion.go
cp ../../source/music/api/v1/rockband_webhook.go api/v1/rockband_webhook.go 

# catering v2beta1
# just need to change Spec and Status (no kubebuilder tags)
# on rockbands_types.go
# create a rockbands_conversion.go - you can copy from another version
# just remember to change the package name

# Enabling webhooks and fixing kubebuilder
cp ../../source/music/config/default/kustomization.yaml config/default/kustomization.yaml
cp ../../source/music/config/certmanager/certificate.yaml config/certmanager/certificate.yaml 
cp ../../source/music/config/crd/kustomization.yaml config/crd/kustomization.yaml
cp ../../source/music/config/crd/patches/* config/crd/patches/

make manifests && make generate && make install

# copying/generating the samples
cp ../../source/music/config/samples/* config/samples/

# edit the samples to match your 

# creating the image
export IMG=quay.io/brito_rafa/music-controller:case-a-target-v0.1

make docker-build
make docker-push
kustomize build config/default > ../case-a-target.yaml
```



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
