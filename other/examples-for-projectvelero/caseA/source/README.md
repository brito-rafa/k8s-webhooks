# Second Example: Two Versions of an API Group

If you have not seen it yet, please make sure you follow [our first example](/single-gvk/README.md) with the a single API Group version.  This first example is a building block for the code in here.

Anyone can code multiple API Group versions from scratch using kubebuilder, but for the sake of academic purposes, I found easier to start with a single API Group version and making sure the webhook is fully operational before adding a second API Group.

## CRD Design Decisions

For this example, I am coding the use case of a backward compability of the existent API Group `v1`. So, I am adding an earlier version `v1alpha1` of the same API group while keeping `v1` as the default version. You might want to read and follow the convention of the [K8s API versioning.](https://kubernetes.io/docs/reference/using-api/#api-versioning)

Again, this is an academic example, you are free to create a newer version of the API group instead of an older one.

At this time, you need to define two things on your API Group:

- Which API Group Version will be set as the [storage version](https://book.kubebuilder.io/multiversion-tutorial/api-changes.html#storage-versions) A.K.A preferred version, which in this example will be  `v1` (same API Group Version as the [first example.](/single-gvk/README.md)
- What schema changes will be between the `v1alpha1` and `v1` versions, in order words, what changed between versions that my controller will need to adapt.

 In this academic example, my `RockBandv1alpha1` will **not** have one particular field  `Spec.LeadSinger` if compared to `RockBandv1`:

 ```go v1alpha1
 // RockBandSpec defines the desired state of RockBand
type RockBandSpec struct {
	// +kubebuilder:validation:Optional
	Genre string `json:"genre"`
	// +kubebuilder:validation:Optional
	NumberComponents int32 `json:"numberComponents"`
}
 ```

 My webhook's job is converting the custom resources (CRs) back and forth from `v1` and `v1alpha1` while doing its best to support the lack of `Spec.LeadSinger` field.

 The final code of the CRD+controller+webhook is already at [multiple-gvk/music/](multiple-gvk/music/) directory.

 You can skip this step-by-step and deploy the final result of this controller and CRDs running:

 ``` bash
 # ONLY IF YOU HAVE NOT DONE SO
 # cert-manager is required for any kubebuilder-created webhook
 # $ kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v1.0.3/cert-manager.yaml

 # deploying multiple-gvk music controller and CRD
 $ kubectl apply --validate=false -f https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/multiple-gvk/multiple-gvk-v0.1.yaml

 # checking existent CRs with the v1alpha1 API
 $ kubectl get rockbands.v1alpha1.music.example.io -A -o yaml
 (...)
 - apiVersion: music.example.io/v1alpha1
  kind: RockBand
  metadata:
    annotations:
      rockbands.v1.music.example.io/leadSinger: John Lennon
    name: beatles
 ```

Follow the next sections for the step-by-step.


## Creating a second API Group version with kubebuilder

Again, starting from the where [our first example](/single-gvk/README.md) left off. Run the following commands:

```
$ kubebuilder create api --group music --version v1alpha1 --kind RockBand --resource=true --controller=false

$ kubebuilder create webhook --group music --version v1alpha1 --kind RockBand --conversion
Writing scaffold for you to edit...
Webhook server has been set up for you.
You need to implement the conversion.Hub and conversion.Convertible interfaces for your CRD types.
api/v1alpha1/rockband_webhook.go
```

You now should see a second directory under `music/api` with the version `v1alpha1`.


Let's edit `music/api/v1alpha1/rockband_types.go` and add the same fields as `v1` but without `Spec.LeadSinger`.

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

Make sure the  `v1/rockband_types.go` has the kubebuilder tag to be the [storage version](https://book.kubebuilder.io/multiversion-tutorial/api-changes.html#storage-versions), A.K.A. the API Group [preferred version](https://github.com/kubernetes/apimachinery/blob/master/pkg/apis/meta/v1/types.go):


``` go v1
// +kubebuilder:storageversion

// RockBand is the Schema for the rockbands API
type RockBand struct {
	metav1.TypeMeta   `json:",inline"`
```

Test the introduction of `v1alpha1` changes running the following:

```
make manifests && make generate
```

Check the file `music/config/crd/bases/music.example.io_rockbands.yaml` for both versions.

Then run:

```bash
$ make install
(...)
kustomize build config/crd | kubectl apply -f -
customresourcedefinition.apiextensions.k8s.io/rockbands.music.example.io configured
```

Check both versions are installed and preferred version is `v1`. This is a very useful command:

```bash
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

Then, check the CRD fields of the `v1alpha1`:

```bash
$ kubectl get crd rockbands.music.example.io -o yaml

(...)
  - name: v1alpha1
    schema:
 (...)
          spec:
            description: RockBandSpec defines the desired state of RockBand
            properties:
              genre:
                type: string
              numberComponents:
                format: int32
                type: integer
            type: object
(...)
```

### Coding API Group conversion

Now that our CRD supports both versions `v1` and `v1alpha` (with `v1` being the preferred), it is time to code the webhook to convert the object back and forth the versions.

- `main.go`

The command `kubebuilder create webhook ... --conversion`, should have added a second `SetupWebhookWithManager` call on your existing `main.go`. 
Please note the difference the second call points to the `RockBandv1alpha1`:

```go main.go

	if err = (&musicv1.RockBand{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "RockBand")
		os.Exit(1)
    }
    // New call, automatically added by kubebuilder

	if err = (&musicv1alpha1.RockBand{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "RockBand")
		os.Exit(1)
	}

```

Now, we need to write code for the conversion using the [Kubebuilder book references](https://book.kubebuilder.io/multiversion-tutorial/conversion.html).


- `music/api/v1/rockband_conversion.go`:

Under preferred version api directory, you will only need the Hub function:

```go /multiple-gvk/music/api/v1/rockband_conversion.go
package v1

// conversion - v1 is the Hub
// https://book.kubebuilder.io/multiversion-tutorial/conversion.html
func (*RockBand) Hub() {}

```

- `music/api/v1alpha1/rockband_conversion.go`, entire file [here](/multiple-gvk/music/api/v1alpha1/rockband_conversion.go):

This is the "work-horse" of the conversion logic. Remember that `v1` is the storage version, so for every non-storage version, you will need a conversion.go file with `ConvertTo` and `ConvertFrom` functions.

In my academic example, if one creates a CR using the API RockBandv1alpha1, we will add a default leadSinger (in our case, as variable `defaultValueLeadSingerConverter` set with string `Converted from v1alpha1`). 

Additionally, I did something clever: I am leveraging the [Annotation struct](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/) to present v1 fields back and forth on v1alpha1 objects. Annotation is a way to attach arbitrary data on any K8s object.

See below the snippet of the code, with the entire file [here](/multiple-gvk/music/api/v1alpha1/rockband_conversion.go):

```go
var (
	// this is the annotation key to keep leadSinger value if converted from v1 to v1alpha1
	leadSingerAnnotation = "rockbands.v1.music.example.io/leadSinger"
	// default leadSinger string to be used when converting from v1alpha1 to v1
	defaultValueLeadSingerConverter = "Converted from v1alpha1"
)

// ConvertTo converts this RockBand v1alpha1 to the Hub version (v1)
func (src *RockBand) ConvertTo(dstRaw conversion.Hub) error {
    dst := dstRaw.(*v1.RockBand)
    (...)
        dst.Spec.LeadSinger = defaultValueLeadSingerConverter
    (...)
}

// ConvertFrom converts from the Hub version (v1) to this version valpha1
func (dst *RockBand) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1.RockBand)
(...)

	// Retaining the v1 LeadSinger values as an annotation
	// Saving fields as annotation is a great way
	// to keep information back and forth between legacy and modern API Groups

	// if the leadSinger is already is set as the default value from v1alpha1 (see ConvertTo)
	// do not bother to create an annotation

	if src.Spec.LeadSinger != defaultValueLeadSingerConverter {
		annotations := dst.GetAnnotations()

		if annotations == nil {
			annotations = make(map[string]string)
		}

		annotations[leadSingerAnnotation] = src.Spec.LeadSinger

		dst.SetAnnotations(annotations)

    }
    (...)
}
```

## Deploying and testing the API Group Version conversion 

Let's build and deploy the controller:

```bash
$ export IMG=quay.io/brito_rafa/music-controller:multiple-gvk-v0.1
$ make docker-build
$ make docker-push
# do not forget to make the image public (or create a pull secret)
$ make deploy IMG=quay.io/brito_rafa/music-controller:multiple-gvk-v0.1
$ $ kubectl get pods -n music-system
NAME                                        READY   STATUS    RESTARTS   AGE
music-controller-manager-54468879c9-rdznh   2/2     Running   0          2m27s
```
Let's see the controller logs and pay attention on `registering webhook     {"path": "/convert"}`:

```bash
$ kubectl logs music-controller-manager-54468879c9-rdznh -n music-system manager
2020-11-03T17:26:55.596Z        INFO    controller-runtime.metrics      metrics server is starting to listen {"addr": "127.0.0.1:8080"}
2020-11-03T17:26:55.596Z        INFO    controller-runtime.builder      Registering a mutating webhook       {"GVK": "music.example.io/v1, Kind=RockBand", "path": "/mutate-music-example-io-v1-rockband"}
2020-11-03T17:26:55.596Z        INFO    controller-runtime.webhook      registering webhook     {"path": "/mutate-music-example-io-v1-rockband"}
2020-11-03T17:26:55.596Z        INFO    controller-runtime.builder      Registering a validating webhook     {"GVK": "music.example.io/v1, Kind=RockBand", "path": "/validate-music-example-io-v1-rockband"}
2020-11-03T17:26:55.596Z        INFO    controller-runtime.webhook      registering webhook     {"path": "/validate-music-example-io-v1-rockband"}
2020-11-03T17:26:55.596Z        INFO    controller-runtime.webhook      registering webhook     {"path": "/convert"}
2020-11-03T17:26:55.596Z        INFO    controller-runtime.builder      conversion webhook enabled  {"object": {"metadata":{"creationTimestamp":null},"spec":{"genre":"","numberComponents":0,"leadSinger":""},"status":{"lastPlayed":""}}}
2020-11-03T17:26:55.597Z        INFO    controller-runtime.builder      skip registering a mutating webhook, admission.Defaulter interface is not implemented        {"GVK": "music.example.io/v1alpha1, Kind=RockBand"}
2020-11-03T17:26:55.597Z        INFO    controller-runtime.builder      skip registering a validating webhook, admission.Validator interface is not implemented      {"GVK": "music.example.io/v1alpha1, Kind=RockBand"}
2020-11-03T17:26:55.597Z        INFO    controller-runtime.builder      conversion webhook enabled  {"object": {"metadata":{"creationTimestamp":null},"spec":{"genre":"","numberComponents":0},"status":{"lastPlayed":""}}}
2020-11-03T17:26:55.597Z        INFO    setup   starting manager
I1103 17:26:55.597731       1 leaderelection.go:242] attempting to acquire leader lease  music-system/9f9c4fd4.example.io...
2020-11-03T17:26:55.598Z        INFO    controller-runtime.manager      starting metrics server {"path": "/metrics"}
2020-11-03T17:26:55.598Z        INFO    controller-runtime.webhook.webhooks     starting webhook server
2020-11-03T17:26:55.599Z        INFO    controller-runtime.certwatcher  Updated current TLS certificate
I1103 17:26:55.702653       1 leaderelection.go:252] successfully acquired lease music-system/9f9c4fd4.example.io
2020-11-03T17:26:55.793Z        INFO    controller-runtime.webhook      serving webhook server  {"host": "", "port": 9443}
2020-11-03T17:26:55.702Z        DEBUG   controller-runtime.manager.events       Normal  {"object": {"kind":"ConfigMap","namespace":"music-system","name":"9f9c4fd4.example.io","uid":"80b81443-de6d-4383-ba43-70c97c43b553","apiVersion":"v1","resourceVersion":"1081"}, "reason": "LeaderElection", "message": "music-controller-manager-54468879c9-rdznh_4ca3cf7a-1a61-4b71-abac-ed1d40f090ad became leader"}
2020-11-03T17:26:55.794Z        INFO    controller-runtime.certwatcher  Starting certificate watcher
2020-11-03T17:26:55.795Z        INFO    controller-runtime.controller   Starting EventSource    {"controller": "rockband", "source": "kind source: /, Kind="}
2020-11-03T17:26:55.896Z        INFO    controller-runtime.controller   Starting Controller     {"controller": "rockband"}
2020-11-03T17:26:55.896Z        INFO    controller-runtime.controller   Starting workers        {"controller": "rockband", "worker count": 1}
```


I added three examples on the sample directory in how to create and display the CRs.

### Creating a v1 CR and displaying as v1alpha1

Let's create the same sample `beatles` CR using `RockBandv1` that we created on the first example.
Then we will list it using `RockBandv1alpha1` over the command `kubectl get rockbands.v1alpha1.music.example.io -o yaml`:

```bash
$ kubectl create -f multiple-gvk/music/config/samples/music_v1_rockband.yaml 
rockband.music.example.io/beatles created

$ kubectl get rockbands.v1alpha1.music.example.io -o yaml
(...)
- apiVersion: music.example.io/v1alpha1
  kind: RockBand
  metadata:
    annotations:
      rockbands.v1.music.example.io/leadSinger: John Lennon
      (...)
  spec:
    genre: 60s rock
    numberComponents: 4
  status:
    lastPlayed: "2020"

```

First of all, see the the API version is `music.example.io/v1alpha1`. Then see that `Spec.LeadSinger` disappeared but there is an annotation field with `rockbands.v1.music.example.io/leadSinger` set as `John Lennon` (our validation code from first exampled kicked in as well).

### Creating v1alpha1 CR as displaying as v1

The second example is creating a CR using `RockBandv1alpha1` API which does not have a `Spec.LeadSinger`:

```bash
$ cat multiple-gvk/music/config/samples/music_v1alpha1_rockband.yaml 
apiVersion: music.example.io/v1alpha1
kind: RockBand
metadata:
  name: pearljam
spec:
  # Add fields here
  genre: Grunge
  numberComponents: 5

$ kubectl create -f  multiple-gvk/music/config/samples/music_v1alpha1_rockband.yaml 
rockband.music.example.io/pearljam created

$ kubectl get rockbands.v1.music.example.io -o yaml
(...)
- apiVersion: music.example.io/v1
  kind: RockBand
  metadata:
(...)
    name: pearljam
    namespace: default
(...)
  spec:
    genre: Grunge
    leadSinger: Converted from v1alpha1
    numberComponents: 5
  status:
    lastPlayed: "2020"
```

### Creating v1alpha1 CR with v1 field on annotation

This last example is how to create a CR using `v1alpha1` schema but still adding `v1` fields as part of annotation.

```bash
$ cat multiple-gvk/music/config/samples/music_v1alpha1_rockband-with-v1-field-as-annotation.yaml 
apiVersion: music.example.io/v1alpha1
kind: RockBand
metadata:
  name: ramones
  annotations:
    rockbands.v1.music.example.io/leadSinger: Joey
spec:
  # Add fields here
  genre: Punk
  numberComponents: 4

$ kubectl create -f multiple-gvk/music/config/samples/music_v1alpha1_rockband-with-v1-field-as-annotation.yaml 
rockband.music.example.io/ramones created

$ kubectl get rockbands.v1.music.example.io -o yaml
(...)
- apiVersion: music.example.io/v1
  kind: RockBand
  metadata:
    annotations:
      rockbands.v1.music.example.io/leadSinger: Joey
(...)
    name: ramones
    namespace: default
 (...)
  spec:
    genre: Punk
    leadSinger: Joey
    numberComponents: 4
  status:
    lastPlayed: "2020"
```

To conclude, we showed here how to create a CRD+controller+webhook using kubebuilder while supporting multiple API Group versions.



## Errors Found during this code


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

Solution: bad `webhook_in_rockbands.yaml` - I had to manually fix the the correct namespace `music-system` and service name on this file. Kubebuilder was creating files with crd v1beta1 extensions and I had to troubleshoot and adapt this template to crd v1.

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

Solution: check your ConvertTo/From code. The ObjectMeta must exist, the annotation field was not defined initially.

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
```

### cert-manager CRDs

Error:
```
unable to recognize "multiple-gvk/multiple-gvk-v0.1.yaml": no matches for kind "Certificate" in version "cert-manager.io/v1alpha2"
unable to recognize "multiple-gvk/multiple-gvk-v0.1.yaml": no matches for kind "Issuer" in version "cert-manager.io/v1alpha2"
```

Solution: install cert-manager











