# Case A

RockBand on Source cluster: v1 (preferred), v1alpha1

RockBand on Target cluster: v1 (preferred), v2beta1

Expected result: v1 being used during restore.

## Deploying the Case A

### On Source Cluster

Run

```bash
curl -k -s https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-a/source-cluster.sh | bash
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



## Coding and Scaffolding

### On Target


```bash
mkdir music
cd music/
go mod init music
kubebuilder init --domain example.io
kubebuilder create api --group music --version v1 --kind RockBand --resource=true --controller=true
kubebuilder create webhook --group music --version v1 --kind RockBand --defaulting --programmatic-validation
kubebuilder create api --group music --version v2beta1 --kind RockBand --resource=true --controller=false
kubebuilder create webhook --group music --version v2beta1 --kind RockBand --conversion

# setting to use CRD
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

make manifests && make generate

# copying/generating the samples


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


### Source Cluster

```bash
# cert-manager
kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v1.0.3/cert-manager.yaml
# controller and CRD

# example CR for testing - only required on the source cluster
```

### Target Cluster
