# First Example

A single version of a kind on an api group.

As I mentioned before, this code is focused only to create a CRD/webhook/controller from scratch, not doing anything business-fancy.

We will do minimum code to update or instrument the CRD, but the focus is managing the schema changes of CRD and how webhook will help us. 

Of course, you can use this tutorial to spring-board to a more serious controller. 
For such, I recommend you to look at our other example https://github.com/embano1/codeconnect-vm-operator.

## Scaffolding

First step is using the kubebuilder scaffolding:

```
mkdir music; cd music
go mod init music
kubebuilder init --domain example.io
kubebuilder create api --group music --version v1 --kind RockBand
# press y for both "Create Resource" and "Create Controller"
kubebuilder create webhook --group music --version v1 --kind RockBand --defaulting --programmatic-validation
```

It creates a directory structure. The final result of this code is under music/ subdirectory here.

## Creating our own CRD business logic

Our CRD schema is created by scaffolding but does not have any business logic.
For this version v1, I will set the following fields on RockBand.Spec and RockBand.Status on the file `music/api/v1/rockband_types.go`:

```go
// RockBandSpec defines the desired state of RockBand
type RockBandSpec struct {
	// +kubebuilder:validation:Optional
	Genre string `json:"genre"`
	// +kubebuilder:validation:Optional
	NumberComponents int32 `json:"numberComponents"`
	// +kubebuilder:validation:Optional
	LeadSinger string `json:"leadSinger"`
}

// RockBandStatus defines the observed state of RockBand
type RockBandStatus struct {
	LastPlayed string `json:"lastPlayed"`
}
```

As you can see, this is really a simple and silly example. RockBand has genre, number of components and lead singer as part of the Spec struc. And they are all optional. On Status struct, I decided to add on arbitrarily field called lastPlayed.
In this example, the Status field does not have a Spec counterpart. In real world, the Spec and Status fields are similar and controller reconciles them.

We want to use the latest CRD features, for such, please edit `Makefile` and set to the following line:

```
CRD_OPTIONS ?= "crd:preserveUnknownFields=false,crdVersions=v1,trivialVersions=true"
```

*Attention*: I had issues with the `Makefile` default CRD line, so make sure you change to the above.

## Generating and Installing the CRD

Generate the CRD running:
```
make manifests && make generate
```

If you have issues, please refer to the original code of this directory.

Now inspect the file `config/crd/bases/music.example.io_rockbands.yaml`. 
If you apply this files as is on your cluster, you will deploy the CRD.

Alternatively, you can run:

```
make install
```

Once installed, you can check the CRDs on your cluster running:

```
kubectl get crds
```

Listing specific versions of the object:

```
kubectl get rockband.v1alpha1.music.example.io beatles -o yaml
```

## Compiling and Starting the Controller

I do not care about the busines logic on the controller, so the only thing my controller will do is setting the lastPlayed Status field with the current year.
Here is the snippet of the `music/controllers/rockband_controller.go`

```go
  // your logic here

	rb := &musicv1.RockBand{}
	if err := r.Client.Get(ctx, req.NamespacedName, rb); err != nil {
		// add some debug information if it's not a NotFound error
		if !k8serr.IsNotFound(err) {
			log.Error(err, "unable to fetch RockBand")
		}
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	msg := fmt.Sprintf("received reconcile request for %q (namespace: %q)", rb.GetName(), rb.GetNamespace())
	log.Info(msg)

	if rb.Status.LastPlayed == "" {
		year := time.Now().Year()
		// Adding the year in Status filed for now
		rb.Status.LastPlayed = strconv.Itoa(year)
		if err := r.Status().Update(ctx, rb); err != nil {
			log.Error(err, "unable to update RockBand status")
			return ctrl.Result{}, err
		}
	}
```

At this time, the controller must be able to be compiled without errors. You can try to run:

```
export IMG=quay.io/brito_rafa/music-controller:single-gvk-v0.1
make docker-build
docker push quay.io/brito_rafa/music-controller:single-gvk-v0.1
make deploy
```

## CRD Mutation and Validation

For this demo, we came up with a couple silly rules just to make a point how to mutate and validate RockBands objects.

Mutation: 
1. if LeadSinger is not specified, set it as "TBD".
2. if RockBand CR name is "beatles" and Lead Singer is "John", set it to "John Lennon" (disclaimer: some Beatles fans (including me) will argue The Beatles did not have "a" Lead Singer).

Validation:
1. We can't create RockBands CRs on "kube-system" namespace
2. We can't update Lead Singer as "John" if RockBand CR name is "beatles" (similar to mutation, but during at update time. Spoiler: this condition will never met because of the mutation logic)
3. We can't update Lead Singer as "Ringo" if RockBand CR name is  "beatles" .


### Code

Scaffolding should have created the following files: `music/main.go` and `music/api/v1/rockband_webhook.go`.

The main.go starts the webhook server.

#### main.go

It must have the invokation of the webhook, otherwise the webhook will not run:

```
	if err = (&musicv1.RockBand{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "RockBand")
		os.Exit(1)
	}
```

#### rockband_webhook.go

This where your mutation and validation logic will reside.

Code snippet for Mutator. Remember, this will execute for each CR request.

Note that during validation calls, one CAN'T change fields, only generate errors - see the kubebuilder option mutating=false.

```
```
## Testing the webhook

### Creation

I used the follow example to test my first mutation and validation:

```
$ kubectl create -f beatles.yaml -n default
$ cat beatles.yaml
apiVersion: music.example.io/v1
kind: RockBand
metadata:
  name: beatles
spec:
  # Add fields here
  genre: '60s rock'
  numberComponents: 4
  leadSinger: John
```

Controller logs during creation:

```
2020-10-28T21:36:09.478Z        DEBUG   controller-runtime.webhook.webhooks     received request    {"webhook": "/mutate-music-example-io-v1-rockband", "UID": "ae7f50b2-df85-41fc-9460-9fdde49882f0", "kind": "music.example.io/v1, Kind=RockBand", "resource": {"group":"music.example.io","version":"v1","resource":"rockbands"}}
2020-10-28T21:36:09.478Z        INFO    rockband-resource       mutator default {"name": "beatles", "namespace": "default"}
2020-10-28T21:36:09.480Z        DEBUG   controller-runtime.webhook.webhooks     wrote response  {"webhook": "/mutate-music-example-io-v1-rockband", "UID": "ae7f50b2-df85-41fc-9460-9fdde49882f0", "allowed": true, "result": {}, "resultError": "got runtime.Object without object metadata: &Status{ListMeta:ListMeta{SelfLink:,ResourceVersion:,Continue:,RemainingItemCount:nil,},Status:,Message:,Reason:,Details:nil,Code:200,}"}
2020-10-28T21:36:09.487Z        DEBUG   controller-runtime.webhook.webhooks     received request    {"webhook": "/validate-music-example-io-v1-rockband", "UID": "ef366e5c-9cbb-4df5-8d8e-cea6ba28beab", "kind": "music.example.io/v1, Kind=RockBand", "resource": {"group":"music.example.io","version":"v1","resource":"rockbands"}}
2020-10-28T21:36:09.487Z        INFO    rockband-resource       validate create {"name": "beatles", "namespace": "default", "lead singer": "John Lennon"}
2020-10-28T21:36:09.487Z        DEBUG   controller-runtime.webhook.webhooks     wrote response  {"webhook": "/validate-music-example-io-v1-rockband", "UID": "ef366e5c-9cbb-4df5-8d8e-cea6ba28beab", "allowed": true, "result": {}, "resultError": "got runtime.Object without object metadata: &Status{ListMeta:ListMeta{SelfLink:,ResourceVersion:,Continue:,RemainingItemCount:nil,},Status:,Message:,Reason:,Details:nil,Code:200,}"}
2020-10-28T21:36:09.495Z        INFO    controllers.RockBand    received reconcile request for "beatles" (namespace: "default")      {"rockband": "default/beatles"}
2020-10-28T21:36:09.506Z        DEBUG   controller-runtime.controller   Successfully Reconciled {"controller": "rockband", "request": "default/beatles"}
2020-10-28T21:36:09.506Z        INFO    controllers.RockBand    received reconcile request for "beatles" (namespace: "default")      {"rockband": "default/beatles"}
2020-10-28T21:36:09.506Z        DEBUG   controller-runtime.controller   Successfully Reconciled {"controller": "rockband", "request": "default/beatles"}
```

Result of the CR (plase note the leadSinger as John Lennon):

```
 kubectl get rockband beatles -o yaml
apiVersion: music.example.io/v1
kind: RockBand
metadata:
  creationTimestamp: "2020-10-28T21:36:09Z"
  generation: 1
  name: beatles
  namespace: default
  resourceVersion: "144885"
  selfLink: /apis/music.example.io/v1/namespaces/default/rockbands/beatles
  uid: 2320f90d-cb13-4398-ba35-a78dc9709912
spec:
  genre: 60s rock
  leadSinger: John Lennon
  numberComponents: 4
status:
  lastPlayed: "2020"
```

Let's test the validation creation now:

```
$ kubectl create -f beatles -n kube-system
Error from server (RockBand.music.example.io "beatles" is invalid: metadata.namespace: Invalid value: "kube-system": is forbidden to have rockbands.): error when creating "music_v1_rockband.yaml": admission webhook "vrockband.kb.io" denied the request: RockBand.music.example.io "beatles" is invalid: metadata.namespace: Invalid value: "kube-system": is forbidden to have rockbands.
```

Here are the controller logs for the creation validation:

```
2020-10-28T21:40:32.087Z        DEBUG   controller-runtime.webhook.webhooks     received request    {"webhook": "/mutate-music-example-io-v1-rockband", "UID": "f9c03827-f1ad-4338-983f-a90d1d897dab", "kind": "music.example.io/v1, Kind=RockBand", "resource": {"group":"music.example.io","version":"v1","resource":"rockbands"}}
2020-10-28T21:40:32.087Z        INFO    rockband-resource       mutator default {"name": "beatles", "namespace": "kube-system"}
2020-10-28T21:40:32.087Z        DEBUG   controller-runtime.webhook.webhooks     wrote response  {"webhook": "/mutate-music-example-io-v1-rockband", "UID": "f9c03827-f1ad-4338-983f-a90d1d897dab", "allowed": true, "result": {}, "resultError": "got runtime.Object without object metadata: &Status{ListMeta:ListMeta{SelfLink:,ResourceVersion:,Continue:,RemainingItemCount:nil,},Status:,Message:,Reason:,Details:nil,Code:200,}"}
2020-10-28T21:40:32.090Z        DEBUG   controller-runtime.webhook.webhooks     received request    {"webhook": "/validate-music-example-io-v1-rockband", "UID": "2e6212dc-8bd6-468a-92e2-2a078eb41122", "kind": "music.example.io/v1, Kind=RockBand", "resource": {"group":"music.example.io","version":"v1","resource":"rockbands"}}
2020-10-28T21:40:32.090Z        INFO    rockband-resource       validate create {"name": "beatles", "namespace": "kube-system", "lead singer": "John Lennon"}
2020-10-28T21:40:32.090Z        DEBUG   controller-runtime.webhook.webhooks     wrote response  {"webhook": "/validate-music-example-io-v1-rockband", "UID": "2e6212dc-8bd6-468a-92e2-2a078eb41122", "allowed": false, "result": {}, "resultError": "got runtime.Object without object metadata: &Status{ListMeta:ListMeta{SelfLink:,ResourceVersion:,Continue:,RemainingItemCount:nil,},Status:,Message:,Reason:RockBand.music.example.io \"beatles\" is invalid: metadata.namespace: Invalid value: \"kube-system\": is forbidden to have rockbands.,Details:nil,Code:403,}"}
```


Let's test the update validation now. If you remember, the code does not let you to setup leadSinger as "John" if the rockband is "beatles".

```
$ kubectl edit rockband beatles -n default
(...)
spec:
  genre: 60s rock
  leadSinger: John
(...)
rockband.music.example.io/beatles edited
```

Wait. It let me. Why? Let's see the controller logs:

```
2020-10-28T21:44:32.065Z        DEBUG   controller-runtime.webhook.webhooks     received request    {"webhook": "/mutate-music-example-io-v1-rockband", "UID": "dd68946a-9b9f-4b4b-bec2-41db5116da17", "kind": "music.example.io/v1, Kind=RockBand", "resource": {"group":"music.example.io","version":"v1","resource":"rockbands"}}
2020-10-28T21:44:32.065Z        INFO    rockband-resource       mutator default {"name": "beatles", "namespace": "default"}
2020-10-28T21:44:32.065Z        DEBUG   controller-runtime.webhook.webhooks     wrote response  {"webhook": "/mutate-music-example-io-v1-rockband", "UID": "dd68946a-9b9f-4b4b-bec2-41db5116da17", "allowed": true, "result": {}, "resultError": "got runtime.Object without object metadata: &Status{ListMeta:ListMeta{SelfLink:,ResourceVersion:,Continue:,RemainingItemCount:nil,},Status:,Message:,Reason:,Details:nil,Code:200,}"}
2020-10-28T21:44:32.068Z        DEBUG   controller-runtime.webhook.webhooks     received request    {"webhook": "/validate-music-example-io-v1-rockband", "UID": "404e2e73-9110-4217-8ea1-502e24fb102f", "kind": "music.example.io/v1, Kind=RockBand", "resource": {"group":"music.example.io","version":"v1","resource":"rockbands"}}
2020-10-28T21:44:32.068Z        INFO    rockband-resource       validate update {"name": "beatles", "namespace": "default", "lead singer": "John Lennon"}
2020-10-28T21:44:32.068Z        DEBUG   controller-runtime.webhook.webhooks     wrote response  {"webhook": "/validate-music-example-io-v1-rockband", "UID": "404e2e73-9110-4217-8ea1-502e24fb102f", "allowed": true, "result": {}, "resultError": "got runtime.Object without object metadata: &Status{ListMeta:ListMeta{SelfLink:,ResourceVersion:,Continue:,RemainingItemCount:nil,},Status:,Message:,Reason:,Details:nil,Code:200,}"}
```

The answer is because the mutation kicks off before validation. So "John" will be converted to "John Lennon" before the validation.

Now, let's try the other update validation, which is changing the leadSinger to Ringo.

```
$ kubectl edit rockband beatles -n default
(...)
spec:
  genre: 60s rock
  leadSinger: Ringo
(...)
error: rockbands.music.example.io "beatles" could not be patched: admission webhook "vrockband.kb.io" denied the request: RockBand.music.example.io "beatles" is invalid: spec.leadSinger: Invalid value: "Ringo": was the drummer. Suggest you to pick John or Paul.
You can run `kubectl replace -f /var/folders/p2/vpgr25xn16777ll0y7fsmc2c0000gp/T/kubectl-edit-e3iws.yaml` to try this update again.
```

Let's see the controller logs:

```
2020-10-28T22:07:26.509Z        DEBUG   controller-runtime.webhook.webhooks     received request    {"webhook": "/mutate-music-example-io-v1-rockband", "UID": "128f1b35-6b58-4230-8704-546b5f2e1d12", "kind": "music.example.io/v1, Kind=RockBand", "resource": {"group":"music.example.io","version":"v1","resource":"rockbands"}}
2020-10-28T22:07:26.511Z        INFO    rockband-resource       mutator default {"name": "beatles", "namespace": "default"}
2020-10-28T22:07:26.512Z        DEBUG   controller-runtime.webhook.webhooks     wrote response  {"webhook": "/mutate-music-example-io-v1-rockband", "UID": "128f1b35-6b58-4230-8704-546b5f2e1d12", "allowed": true, "result": {}, "resultError": "got runtime.Object without object metadata: &Status{ListMeta:ListMeta{SelfLink:,ResourceVersion:,Continue:,RemainingItemCount:nil,},Status:,Message:,Reason:,Details:nil,Code:200,}"}
2020-10-28T22:07:26.515Z        DEBUG   controller-runtime.webhook.webhooks     received request    {"webhook": "/validate-music-example-io-v1-rockband", "UID": "36290e4a-0cbc-4db6-9e59-42915b041364", "kind": "music.example.io/v1, Kind=RockBand", "resource": {"group":"music.example.io","version":"v1","resource":"rockbands"}}
2020-10-28T22:07:26.515Z        INFO    rockband-resource       validate update {"name": "beatles", "namespace": "default", "lead singer": "Ringo"}
2020-10-28T22:07:26.516Z        DEBUG   controller-runtime.webhook.webhooks     wrote response  {"webhook": "/validate-music-example-io-v1-rockband", "UID": "36290e4a-0cbc-4db6-9e59-42915b041364", "allowed": false, "result": {}, "resultError": "got runtime.Object without object metadata: &Status{ListMeta:ListMeta{SelfLink:,ResourceVersion:,Continue:,RemainingItemCount:nil,},Status:,Message:,Reason:RockBand.music.example.io \"beatles\" is invalid: spec.leadSinger: Invalid value: \"Ringo\": was the drummer. Suggest you to pick John or Paul.,Details:nil,Code:403,}"}
2020-10-28T22:07:26.521Z        DEBUG   controller-runtime.webhook.webhooks     received request    {"webhook": "/mutate-music-example-io-v1-rockband", "UID": "10476cf0-60b5-442c-9c45-45a00cc6357c", "kind": "music.example.io/v1, Kind=RockBand", "resource": {"group":"music.example.io","version":"v1","resource":"rockbands"}}
2020-10-28T22:07:26.521Z        INFO    rockband-resource       mutator default {"name": "beatles", "namespace": "default"}
2020-10-28T22:07:26.522Z        DEBUG   controller-runtime.webhook.webhooks     wrote response  {"webhook": "/mutate-music-example-io-v1-rockband", "UID": "10476cf0-60b5-442c-9c45-45a00cc6357c", "allowed": true, "result": {}, "resultError": "got runtime.Object without object metadata: &Status{ListMeta:ListMeta{SelfLink:,ResourceVersion:,Continue:,RemainingItemCount:nil,},Status:,Message:,Reason:,Details:nil,Code:200,}"}
2020-10-28T22:07:26.524Z        DEBUG   controller-runtime.webhook.webhooks     received request    {"webhook": "/validate-music-example-io-v1-rockband", "UID": "bb035ace-c3f1-4531-9bf3-8bd6cdcdf6f7", "kind": "music.example.io/v1, Kind=RockBand", "resource": {"group":"music.example.io","version":"v1","resource":"rockbands"}}
2020-10-28T22:07:26.526Z        INFO    rockband-resource       validate update {"name": "beatles", "namespace": "default", "lead singer": "Ringo"}
2020-10-28T22:07:26.566Z        DEBUG   controller-runtime.webhook.webhooks     wrote response  {"webhook": "/validate-music-example-io-v1-rockband", "UID": "bb035ace-c3f1-4531-9bf3-8bd6cdcdf6f7", "allowed": false, "result": {}, "resultError": "got runtime.Object without object metadata: &Status{ListMeta:ListMeta{SelfLink:,ResourceVersion:,Continue:,RemainingItemCount:nil,},Status:,Message:,Reason:RockBand.music.example.io \"beatles\" is invalid: spec.leadSinger: Invalid value: \"Ringo\": was the drummer. Suggest you to pick John or Paul.,Details:nil,Code:403,}"}
```


### Certs

Webhooks require a tls cert and key.

#### Controller running locally

If you are testing the webhook locally, the certs are expected under the directory `/tmp/k8s-webhook-server/serving-certs`.

Create and install the key/cert running:

```
openssl req -new -newkey rsa:4096 -x509 -sha256 -days 365 -nodes -out tls.crt -keyout tls.key
# answer all the questions
mkdir -p /tmp/k8s-webhook-server/serving-certs
mv tls.* /tmp/k8s-webhook-server/serving-certs
```

During the execution of the controller, possible the webhook will complain to a different directory than /tmp/k8s-webhook-server/serving-certs .
In this case, just copy the tsl.* files to the directory.

It is possible that running locally the webhook is never called. In this case, you will need to deploy them on the cluster (see next section).

#### Controller running on the Cluster

https://cert-manager.io/docs/installation/kubernetes/

```
kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v1.0.3/cert-manager.yaml
```
It should like this:

```
kubectl get pods --namespace cert-manager
NAME                                       READY   STATUS    RESTARTS   AGE
cert-manager-59dbb7958b-t4w68              1/1     Running   0          24s
cert-manager-cainjector-5df5cf79bf-j9h8m   1/1     Running   1          24s
cert-manager-webhook-8557565b68-hpp5f      1/1     Running   0          24s
```

Testing the issuer and cert creation:

```
kubectl create namespace music-system
cat <<EOF > example-io-self-signed.yaml
---
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: example-io-selfsigned
  namespace: music-system
spec:
  selfSigned: {}
---
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: example-io-selfsigned
  namespace: music-system
spec:
  dnsNames:
    - example.io
  secretName: example-io-selfsigned-cert-tls
  issuerRef:
    name: example-io-selfsigned
EOF
kubectl create -f example-io-self-signed.yaml
```

Let's configure the `config/certmanager/certificate.yaml` for a self-signed issuer and a cert.
Accept the defaults but they need to be added under a namespace, the example below, it is "music-system" (instead of the default "system"):

```
apiVersion: cert-manager.io/v1alpha2
kind: Issuer
metadata:
  name: selfsigned-issuer
  namespace: music-system
spec:
  selfSigned: {}
---
apiVersion: cert-manager.io/v1alpha2
kind: Certificate
metadata:
  name: serving-cert  # this name should match the one appeared in kustomizeconfig.yaml
  namespace: music-system
```

Now we need to configure our controller to use the cert under the music-system namespace.
Go to `config/default/kustomization.yaml`

```
```

Uncomment the `crd/kustomization.yaml` for the patches and certmanager

```
patchesStrategicMerge:
# [WEBHOOK] To enable webhook, uncomment all the sections with [WEBHOOK] prefix.
# patches here are for enabling the conversion webhook for each CRD
- patches/webhook_in_rockbands.yaml
# +kubebuilder:scaffold:crdkustomizewebhookpatch

# [CERTMANAGER] To enable webhook, uncomment all the sections with [CERTMANAGER] prefix.
# patches here are for enabling the CA injection for each CRD
- patches/cainjection_in_rockbands.yaml
# +kubebuilder:scaffold:crdkustomizecainjectionpatch

# the following config is for teaching kustomize how to do kustomization for CRDs.
configurations:
- kustomizeconfig.yaml
```

I had an issue where the `config/crd/patches/*yaml` were using apiextensions.k8s.io/v1beta1 instead of piextensions.k8s.io/v1. I had to manually edit the file.

Here is the original error:

```
Error: accumulating resources: accumulateFile "accumulating resources from '../crd': '/Users/rbrito/go/src/music/config/crd' must resolve to a file", accumulateDirector: "recursed accumulation of path '/Users/rbrito/go/src/music/config/crd': no matches for OriginalId apiextensions.k8s.io_v1betav1
```

Another error, because the webhook was generated for v1beta1, the spec is .spec.conversion.webhook.clientconfig (not .spec.conversion.webclientconfig):

```
error: error validating "STDIN": error validating data: ValidationError(CustomResourceDefinition.spec.conversion): unknown field "webhookClientConfig" in io.k8s.apiextensions-apiserver.pkg.apis.apiextensions.v1.CustomResourceConversion; if you choose to ignore these errors, turn validation off with --validate=false
```

## Testing Commands

If you are using kind you do need to push the controller into a public docker registry.
Use load

```
export IMG=quay.io/brito_rafa/music-controller:latest
make docker-build # controller:latest is the default name - see Makefile
kind load docker-image controller:latest --name=crd-dev
make deploy IMG=controller:latest
```


Listing resources:

```
kubectl get api-resources
```

Listing preferred version:

```
kubectl get --raw /apis/music.example.io | jq -r
```



## Errors found during this development

### mutate=true at validation

Make deploy generated the following error:

```
Error: no matches for OriginalId admissionregistration.k8s.io_v1beta1_ValidatingWebhookConfiguration|~X|validating-webhook-configuration; no matches for CurrentId admissionregistration.k8s.io_v1beta1_ValidatingWebhookConfiguration|~X|validating-webhook-configuration; failed to find unique target for patch admissionregistration.k8s.io_v1beta1_ValidatingWebhookConfiguration|validating-webhook-configuration
```



## Multi-API

kubebuilder init --domain example.io
kubebuilder create api --group music --version v1alpha1 --kind RockBand
kubebuilder create api --group music --version v1 --kind RockBand
kubebuilder create webhook --group music --version v1 --kind RockBand --defaulting --programmatic-validation
```

Edit the type files and pick one to be the (preferred version)[https://book.kubebuilder.io/multiversion-tutorial/api-changes.html#storage-versions]:

```
// +kubebuilder:storageversion
```

On your preferred version api directory, you should have rockband_convert.go.
For each other supported version, you will need functions ConvertTo and ConvertFrom functions. They should be rockband_conversion.go.


References:
https://book.kubebuilder.io/multiversion-tutorial/tutorial.html

Edit the type files and pick one to be the (preferred version)[https://book.kubebuilder.io/multiversion-tutorial/api-changes.html#storage-versions]:

```
// +kubebuilder:storageversion
```
