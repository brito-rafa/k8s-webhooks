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

Run the following commands:

```
$ kubebuilder create api --group music --version v1alpha1 --kind RockBand --resource=true --controller=false

$ kubebuilder create webhook --group music --version v1alpha1 --kind RockBand --conversion
Writing scaffold for you to edit...
Webhook server has been set up for you.
You need to implement the conversion.Hub and conversion.Convertible interfaces for your CRD types.
api/v1alpha1/rockband_webhook.go
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


``` go v1
// +kubebuilder:storageversion

// RockBand is the Schema for the rockbands API
type RockBand struct {
	metav1.TypeMeta   `json:",inline"`
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

- `main.go`

Because of the command `kubebuilder create webhook ... --conversion`, your `main.go` should have a second `SetupWebhookWithManager` call, this time to the RockBandv1alpha`:

```go main.go

	if err = (&musicv1.RockBand{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "RockBand")
		os.Exit(1)
	}
	if err = (&musicv1alpha1.RockBand{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "RockBand")
		os.Exit(1)
	}

```

Now, we need to write code for the conversion using the [Kubebuilder book references](https://book.kubebuilder.io/multiversion-tutorial/conversion.html).


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


## Common Errors Found during this code


### webhook-service not found

Error: 
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

Solution: bad `webhook_in_rockbands.yaml` - I had to forcily put the correct namespace `music-system` on this file. The origin of this bug is due to the kubebuilder was creating files with crd v1beta1 extensions.

``` yaml
        service:
          namespace: music-system     # <<<<< FIX
          name: music-webhook-service #<<<<< FIX
          path: /convert
```


### No v1alpha1 webhook

Solution: the main.go was missing the `(&musicv1alpha1.RockBand{}).SetupWebhookWithManager`  - lack of running the `kubectl create webhook .... --conversion`

The correct logs should look like this:

```
2020-11-02T22:28:04.016Z        INFO    controller-runtime.webhook      registering webhook     {"path": "/convert"}
2020-11-02T22:28:04.016Z        INFO    controller-runtime.builder      conversion webhook enabled      {"object": {"metadata":{"creationTimestamp":null},"spec":{"genre":"","numberComponents":0,"leadSinger":""},"status":{"lastPlayed":""}}}
2020-11-02T22:28:04.017Z        INFO    controller-runtime.builder      skip registering a mutating webhook, admission.Defaulter interface is not implemented {"GVK": "music.example.io/v1alpha1, Kind=RockBand"}
2020-11-02T22:28:04.017Z        INFO    controller-runtime.builder      skip registering a validating webhook, admission.Validator interface is not implemented       {"GVK": "music.example.io/v1alpha1, Kind=RockBand"}
```

### Convert Object Panic

Solution: check your ConvertTo/From code. The ObjectMeta must exist, the annotation field is not well coded, etc.

Error:
```
2020-11-02T22:55:50.172Z        DEBUG   controller-runtime.webhook.webhooks     wrote response       {"webhook": "/validate-music-example-io-v1-rockband", "UID": "06c0e871-1200-40e4-a2ba-ced5a12c827d", "allowed": true, "result": {}, "resultError": "got runtime.Object without object metadata: &Status{ListMeta:ListMeta{SelfLink:,ResourceVersion:,Continue:,RemainingItemCount:nil,},Status:,Message:,Reason:,Details:nil,Code:200,}"}
2020/11/02 22:55:50 http: panic serving 10.244.0.1:54282: assignment to entry in nil map
goroutine 412 [running]:
net/http.(*conn).serve.func1(0xc000122780)
        /usr/local/go/src/net/http/server.go:1795 +0x139
panic(0x13d7bc0, 0x172a4e0)
        /usr/local/go/src/runtime/panic.go:679 +0x1b2
music/api/v1alpha1.(*RockBand).ConvertFrom(0xc0000f1400, 0x1770820, 0xc0002d8580, 0x1770820, 0xc0002d8580)
        /workspace/api/v1alpha1/rockband_conversion.go:24 +0xc1
sigs.k8s.io/controller-runtime/pkg/webhook/conversion.(*Webhook).convertObject(0xc0003a1430, 0x174ba20, 0xc0002d8580, 0x174baa0, 0xc0000f1400, 0x174baa0, 0xc0000f1400)
        /go/pkg/mod/sigs.k8s.io/controller-runtime@v0.5.0/pkg/webhook/conversion/conversion.go:142 +0x7bc
sigs.k8s.io/controller-runtime/pkg/webhook/conversion.(*Webhook).handleConvertRequest(0xc0003a1430, 0xc0004f3580, 0xc000547fb0, 0x0, 0x0)
        /go/pkg/mod/sigs.k8s.io/controller-runtime@v0.5.0/pkg/webhook/conversion/conversion.go:107 +0x1f8
sigs.k8s.io/controller-runtime/pkg/webhook/conversion.(*Webhook).ServeHTTP(0xc0003a1430, 0x1770b20, 0xc00064c1c0, 0xc00050aa00)
        /go/pkg/mod/sigs.k8s.io/controller-runtime@v0.5.0/pkg/webhook/conversion/conversion.go:74 +0x10b
sigs.k8s.io/controller-runtime/pkg/webhook.instrumentedHook.func1(0x1770b20, 0xc00064c1c0, 0xc00050aa00)
        /go/pkg/mod/sigs.k8s.io/controller-runtime@v0.5.0/pkg/webhook/server.go:123 +0xfc
net/http.HandlerFunc.ServeHTTP(0xc000616b70, 0x1770b20, 0xc00064c1c0, 0xc00050aa00)
        /usr/local/go/src/net/http/server.go:2036 +0x44
net/http.(*ServeMux).ServeHTTP(0xc000343f00, 0x1770b20, 0xc00064c1c0, 0xc00050aa00)
        /usr/local/go/src/net/http/server.go:2416 +0x1bd
net/http.serverHandler.ServeHTTP(0xc0001aa460, 0x1770b20, 0xc00064c1c0, 0xc00050aa00)
        /usr/local/go/src/net/http/server.go:2831 +0xa4
net/http.(*conn).serve(0xc000122780, 0x17750a0, 0xc00063af40)
        /usr/local/go/src/net/http/server.go:1919 +0x875
created by net/http.(*Server).Serve
        /usr/local/go/src/net/http/server.go:2957 +0x384
2020-11-02T22:55:50.241Z        INFO    controllers.RockBand    received reconcile request for "beatles" (namespace: "default")      {"rockband": "default/beatles"}
2020-11-02T22:55:50.248Z        DEBUG   controller-runtime.controller   Successfully Reconciled      {"controller": "rockband", "request": "default/beatles"}
2020-11-02T22:55:50.248Z        INFO    controllers.RockBand    received reconcile request for "beatles" (namespace: "default")      {"rockband": "default/beatles"}
2020-11-02T22:55:50.248Z        DEBUG   controller-runtime.controller   Successfully Reconciled      {"controller": "rockband", "request": "default/beatles"}
```













