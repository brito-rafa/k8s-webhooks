# Case C

Rationale: `target preferred version != source preferred version; target preferred version does not belong in source supported version; source preferred version belongs in target supported version array`

RockBand on Source cluster: v1 (preferred), v1alpha1

RockBand on Target cluster: v2 (preferred), v1

Expected result: v1 being used during restore.


## Deploying the Case C

### Quick Deploy

#### Quick Deploy on Source Cluster

Reminding case-c source cluster is the same as case-a source cluster:
Run:

```bash
curl -k -s https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-c/source-cluster.sh | bash
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
sleep 10
kubectl apply --validate=false -f https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-a/source/case-a-source.yaml
sleep 10
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
sleep 10
kubectl apply --validate=false -f https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-c/target/case-c-target-manually-.yaml
sleep 10
# no need to create testing CR - should be from Velero backup
```

## Coding and Scaffolding

### Source Cluster: See Case A Source cluster

They are the same steps case a source cluster. No need to reproduce here. 

### Target Cluster

Created from scratch since the preferred version is now v2.

```bash
mkdir music
cd music/
go mod init music
kubebuilder init --domain example.io
kubebuilder create api --group music --version v2 --kind RockBand --resource=true --controller=true
kubebuilder create webhook --group music --version v2 --kind RockBand --defaulting --programmatic-validation
kubebuilder create api --group music --version v1 --kind RockBand --resource=true --controller=false
kubebuilder create webhook --group music --version v1 --kind RockBand --conversion

# setting to use CRD v1 - using the source
cp ../../../case-b/target/music/Makefile .

# Make v2/rockband_types.go - preferred version
# adding the kubebuilder tag // +kubebuilder:storageversion

# re-using v1 types and webhooks
# do not forget to remove the +kubebuilder:storageversion
# from rockband_types.go from the v1


# setting up the hub
# do not forget to change the package name to v2beta2

# coding the v1 conversion to v2

# make new code here, do not forget to change import and package name
# re-using code v1beta1/conversion to v2beta2 - do not forget to change the package name
cp ../../../case-b/target/music/api/v2beta1/rockband_conversion.go api/v1/
# change the code
# same for the controller

# Enabling webhooks and fixing kubebuilder
cp ../../../case-b/target/music/config/default/kustomization.yaml config/default/kustomization.yaml
cp ../../../case-b/target/music/certmanager/certificate.yaml config/certmanager/certificate.yaml 
cp ../../../case-b/target/music/config/crd/kustomization.yaml config/crd/kustomization.yaml
cp ../../../case-b/target/music/config/crd/patches/* config/crd/patches/

# change patches/*webhook for v1 and v2

make manifests && make generate && make install

# copying/generating the samples
cp ../../source/music/config/samples/* config/samples/

# edit the samples to match your case

$ kubectl get --raw /apis/music.example.io | jq
{
  "kind": "APIGroup",
  "apiVersion": "v1",
  "name": "music.example.io",
  "versions": [
    {
      "groupVersion": "music.example.io/v2",
      "version": "v2"
    },
    {
      "groupVersion": "music.example.io/v1",
      "version": "v1"
    }
  ],
  "preferredVersion": {
    "groupVersion": "music.example.io/v2",
    "version": "v2"
  }
}

# creating the image
export IMG=quay.io/brito_rafa/music-controller:case-c-target-v0.1
make docker-build && make docker-push & kustomize build config/default > ../case-c-target.yaml
make deploy IMG=quay.io/brito_rafa/music-controller:case-c-target-v0.1
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