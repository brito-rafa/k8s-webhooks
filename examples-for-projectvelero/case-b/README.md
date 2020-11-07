# Case B

Rationale: `target preferred version != source preferred version; target preferred version belongs in source supported version array`

RockBand on Source cluster: v1 (preferred), v2beta2, v2beta1

RockBand on Target cluster: v2beta2 (preferred), v2beta1

Expected result: v2beta2 being used during restore.

Attention: this use case is very unlikely to happen because the target preferred version should be have a stable version.

## RockBand

Evolution of RockBand.music.example.io schema across versions:

- `RockBandv1alpha1` : Fields `Spec.Genre`, `Spec.NumberComponents`
- `RockBandv1` : all previous plus `Spec.LeadSinger`
- `RockBandv2beta1` : all previous plus `Spec.LeadGuitar`
- `RockBandv2beta2` : all previous plus `Spec.Drummer`
- `RockBandv2` : all previous plus `Spec.Bass`


## Deploying the Case B

### Quick Deploy on Source Cluster

Run:

```bash
curl -k -s https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-b/source-cluster.sh | bash
```

### Quick Deploy on Target Cluster

Run:

```bash
curl -k -s https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-b/target-cluster.sh | bash
```

### Step-by-step Deployment on Source Cluster

Run:

```bash
kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v1.0.3/cert-manager.yaml
sleep 20
kubectl apply --validate=false -f https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-b/source/case-b-source-manually-added-mutations.yaml
sleep 10
# creating the testing CRs
kubectl create --validate=false -f https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-b/source/music/config/samples/music_v2beta2_rockband.yaml
kubectl create --validate=false -f https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-b/source/music/config/samples/music_v2beta1_rockband.yaml
# if you want to test
kubectl get rockbands -A -o yaml
```


### Step-by-step Deployment on Target Cluster

Run:

```bash
kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v1.0.3/cert-manager.yaml
sleep 10
kubectl apply --validate=false -f https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-b/target/case-b-target-manually-added-mutations.yaml
sleep 10
# no need to create testing CR - should be from Velero backup
```

## Coding and Scaffolding

### Source Cluster

Copied from the case-a/target and added v2beta2:

```bash
cp -pr ../../../case-a/target/music/* .
kubebuilder create api --group music --version v2beta2 --kind RockBand --resource=true --controller=false
kubebuilder create webhook --group music --version v2beta2 --kind RockBand --conversion --programmatic-validation
```

And change type, converter and webhook accordingly for each api version.

Change `conversionReviewVersions:` field on `webhook_in_rockbands.yaml`


Installing the CRDs:

```bash
make manifests && make generate && make install
# eyeball the supported and preferred versions
kubectl get --raw /apis/music.example.io | jq 

# it should be this:
{
  "kind": "APIGroup",
  "apiVersion": "v1",
  "name": "music.example.io",
  "versions": [
    {
      "groupVersion": "music.example.io/v1",
      "version": "v1"
    },
    {
      "groupVersion": "music.example.io/v2beta2",
      "version": "v2beta2"
    },
    {
      "groupVersion": "music.example.io/v2beta1",
      "version": "v2beta1"
    }
  ],
  "preferredVersion": {
    "groupVersion": "music.example.io/v1",
    "version": "v1"
  }
}
```

Compiling and Deploying the Code

```bash
export IMG=quay.io/brito_rafa/music-controller:case-b-source-v0.1
make docker-build
make docker-push
kustomize build config/default > ../case-b-source.yaml
```

***ATTENTION:***
There is a bug on kubebuilder that `case-b-source.yaml` is not generating mutating webhook for `v1` and `v2beta1`.

Please see manually edited file `case-b-source-manually-added-mutations.yaml`

### Target cluster

Created from scratch since the preferred version is now v2beta2

```bash
mkdir music
cd music/
go mod init music
kubebuilder init --domain example.io
kubebuilder create api --group music --version v2beta2 --kind RockBand --resource=true --controller=true
kubebuilder create webhook --group music --version v2beta2 --kind RockBand --defaulting --programmatic-validation
kubebuilder create api --group music --version v2beta1 --kind RockBand --resource=true --controller=false
kubebuilder create webhook --group music --version v2beta1 --kind RockBand --conversion

# setting to use CRD v1 - using the source
cp ../../source/music/Makefile Makefile

# Make v2beta2/rockband_types.go - preferred version
# adding the kubebuilder tag // +kubebuilder:storageversion

# re-using v2beta1 types and webhooks
cp ../../source/music/api/v2beta1/rockband_types.go api/v2beta1/rockband_types.go

# re-using v2beta2 webhook
cp ../../source/music/api/v2beta2/rockband_webhook.go api/v2beta2/rockband_webhook.go 

# setting up the hub
# do not forget to change the package name to v2beta2
cp ../../source/music/api/v1/rockband_conversion.go api/v2beta2/

# coding the v2beta1 conversion to v2beta2
cp ../../source/music/api/v2beta1/rockband_conversion.go api/v2beta1/rockband_conversion.go
# make new code here, do not forget to change import and package name


# Enabling webhooks and fixing kubebuilder
cp ../../source/music/config/default/kustomization.yaml config/default/kustomization.yaml
cp ../../source/music/config/certmanager/certificate.yaml config/certmanager/certificate.yaml 
cp ../../source/music/config/crd/kustomization.yaml config/crd/kustomization.yaml
cp ../../source/music/config/crd/patches/* config/crd/patches/

make manifests && make generate && make install

# copying/generating the samples
cp ../../source/music/config/samples/* config/samples/

# edit the samples to match your case

# $ kubectl get --raw /apis/music.example.io | jq
{
  "kind": "APIGroup",
  "apiVersion": "v1",
  "name": "music.example.io",
  "versions": [
    {
      "groupVersion": "music.example.io/v2beta2",
      "version": "v2beta2"
    },
    {
      "groupVersion": "music.example.io/v2beta1",
      "version": "v2beta1"
    }
  ],
  "preferredVersion": {
    "groupVersion": "music.example.io/v2beta2",
    "version": "v2beta2"
  }
}

# creating the image
export IMG=quay.io/brito_rafa/music-controller:case-b-target-v0.1
make docker-build && make docker-push & kustomize build config/default > ../case-b-target.yaml
make deploy IMG=quay.io/brito_rafa/music-controller:case-b-target-v0.1
```



## Generating the YAML and Image

***ATTENTION:***
There is a bug on kubebuilder that is not generating more than one mutating webhook. You will need to generate the yaml as specified in this steps and then add the other mutations manually.

Please see manually edited files `case-b-source-manually-added-mutations.yaml` and 
`case-b-target-manually-added-mutations.yaml`.

### Source
With the code:
```bash
export IMG=quay.io/brito_rafa/music-controller:case-b-source-v0.1
make docker-build
make docker-push
kustomize build config/default > ../case-b-source.yaml

# if you want to test the image 
make deploy IMG=quay.io/brito_rafa/music-controller:case-b-source-v0.1
```

### Target
With the code:
```bash
export IMG=quay.io/brito_rafa/music-controller:case-b-target-v0.1
make docker-build
make docker-push
kustomize build config/default > ../case-b-target.yaml

# if you want to test the image 
make deploy IMG=quay.io/brito_rafa/music-controller:case-b-target-v0.1
```