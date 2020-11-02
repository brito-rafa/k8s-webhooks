# Second Example: Two Versions of an API Group

If you have not seen yet, please go to the [first example with a single version](/single-gvk/README.md) because we will start from the code in there.

You can code multiple versions in one shot, but we found easier starting with one API Group version and  making sure the webhook is operational.

## CRD Design Decisions

For this example, I am coding a backward compability case, which is, 

I chose to make an earlier version - v1alpha1 - of the API group than the existent one (v1).

If you follow the K8s API standards, you will see v1alpha1 comes before than v1.

Since this is an example, you are free to create a newer version of the API group.

At this time, you need to define two other things on your API Groups:

- which will be the prefered version (aka storage version)
- what schema changes will be between the versions

 So, v1 is my prefered version and my `RockBandv1alpha1` will **not** have all fields in  this case, it does not have the field `Spec.LeadSinger`. 


## Creating a second API Group version with kubebuilder

Run the following command:

```
kubebuilder create api --group music --version v1alpha1 --kind RockBand
# answer "y" for "Create Resource" and "n" for the "Create Resource"
```

You now should see a second directory under `music/api` with the version `v1alpha1`.


Let's edit `music/api/v1alpha1/rockband_types.go` and add the same fields as v1 but without LeadSinger.

```go music/api/v1alpha1/rockband_types.go

// RockBandSpec defines the desired state of RockBand
type RockBandSpec struct {
	// +kubebuilder:validation:Optional
	Genre string `json:"genre"`
	// +kubebuilder:validation:Optional
	NumberComponents int32 `json:"numberComponents"`
}

// RockBandStatus defines the observed state of RockBand
type RockBandStatus struct {
	LastPlayed string `json:"lastPlayed"`
}
```

Edit the type files and pick `v1` to be the [preferred version](https://book.kubebuilder.io/multiversion-tutorial/api-changes.html#storage-versions):


``` go
// +kubebuilder:storageversion
```

Test the introduction of `v1alpha` changes running the following:

```
make manifests && make generate
```

Check the file `music/config/crd/bases/music.example.io_rockbands.yaml` for both versions.

Then run:

```
$ make install
(...)
kustomize build config/crd | kubectl apply -f -
customresourcedefinition.apiextensions.k8s.io/rockbands.music.example.io configured
```

Check both versions are installed and preferred version is v1:

```
$ kubectl get --raw /apis/music.example.io | jq -r
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
      "groupVersion": "music.example.io/v1alpha1",
      "version": "v1alpha1"
    }
  ],
  "preferredVersion": {
    "groupVersion": "music.example.io/v1",
    "version": "v1"
  }
```

Then, check the fields of the `v1alpha1`:

```
$ kubectl get crd rockbands.music.example.io -o yaml

(...)

  - name: v1alpha1
    schema:
      openAPIV3Schema:
        description: RockBand is the Schema for the rockbands API
        properties:
          apiVersion:
            description: 'APIVersion defines the versioned schema of this representation
              of an object. Servers should convert recognized schemas to the latest
              internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
            type: string
          kind:
            description: 'Kind is a string value representing the REST resource this
              object represents. Servers may infer this from the endpoint the client
              submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
            type: string
          metadata:
            type: object
          spec:
            description: RockBandSpec defines the desired state of RockBand
            properties:
              genre:
                type: string
              numberComponents:
                format: int32
                type: integer
            type: object
          status:
            description: RockBandStatus defines the observed state of RockBand
            properties:
              lastPlayed:
                type: string
            required:
            - lastPlayed
            type: object
        type: object
(...)
```

### Coding API Group conversion

Kubebuilder book reference: https://book.kubebuilder.io/multiversion-tutorial/conversion.html

Create the following files:

Under preferred version api directory, you will need the Hub

- `music/api/v1/rockband_conversion.go`:

```
package v1

// conversion - v1 is the Hub
// https://book.kubebuilder.io/multiversion-tutorial/conversion.html
func (*RockBand) Hub() {}

```

On other versions, you will need to add the conversion logic. In this case, if one creates the RockBandv1alpha1, we will add "TBD Converter" as LeadSinger.
And if you retrieve the RockBand object over v1alpha1, we will add leadSinger on the annotation.

- `music/api/v1alpha1/rockband_conversion.go`, entire file here:

```
var (
	leadSingerAnnotation = "rockband.music.example.io/lead-singer"
)

func (src *RockBand) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1.RockBand)
	dst.Spec.LeadSinger = "TBD Converter"
	return nil
}

func (dst *RockBand) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1.RockBand)
	dst.Spec.NumberComponents = src.Spec.NumberComponents
	dst.Spec.Genre = src.Spec.Genre
	// Retaining the LeadSinger as an annotation
	dst.Annotations[leadSingerAnnotation] = src.Spec.LeadSinger
	return nil
}
```

## Deploying and testing the conversion 

As usual, rebuild and deploy the controller:

```
export IMG=quay.io/brito_rafa/music-controller:multiple-gvk-v0.1
make docker-build
docker push quay.io/brito_rafa/music-controller:multiple-gvk-v0.1
make deploy IMG=quay.io/brito_rafa/music-controller:multiple-gvk-v0.1
```


## Common Errors during this code


### webhook-service not found

```
$ kubectl get rockbands.v1alpha1.music.example.io -o yaml
apiVersion: v1
items: []
kind: List
metadata:
  resourceVersion: ""
  selfLink: ""
Error from server: conversion webhook for music.example.io/v1alpha1, Kind=RockBandList failed: Post https://webhook-service.system.svc:443/convert?timeout=30s: service "webhook-service" not found
```















