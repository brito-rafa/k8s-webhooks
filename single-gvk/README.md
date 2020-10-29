# First Example: A Single Version of an API Group with one Kind.

As I mentioned before, this code is focused only to create the combination of CRD/webhook/controller from scratch.
We will do minimum code to update or instrument the CRD, but the focus is managing the schema changes of CRD and how webhook will help us. 

But, of course, you can use this tutorial to spring-board to a more serious controller for your business needs. 
In addition, I recommend you to look at our other controller example https://github.com/embano1/codeconnect-vm-operator.

The final result of this code is under music/ subdirectory. You can use it as reference, but we will start from an empty directory.

## Scaffolding

First step is using the kubebuilder scaffolding:

```bash
mkdir music; cd music
go mod init music
kubebuilder init --domain example.io
kubebuilder create api --group music --version v1 --kind RockBand
# press y for both "Create Resource" and "Create Controller"
```

It creates a directory structure.

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

As you can see, this is really a simple and silly example. RockBand has genre, number of components and lead singer as part of the Spec struct.
And they are all optional. On Status struct, I decided to add on arbitrarily field called lastPlayed.
In this example, the Status field does not have a Spec counterpart. In real world, the Spec and Status fields are similar and controller reconciles them.

We want to use the latest CRD features, for such, please edit `Makefile` and set to the following line:

```
CRD_OPTIONS ?= "crd:preserveUnknownFields=false,crdVersions=v1,trivialVersions=true"
```

*Attention*: I had issues with the `Makefile` default CRD line, so make sure you change to the above.

## Generating and Installing the CRD

Generate the CRD running:

```bash
$ make manifests && make generate
```

If you have issues, please refer to the original code of this directory.

Now inspect the file `config/crd/bases/music.example.io_rockbands.yaml`. 
If you apply this files as is on your cluster, you will deploy the CRD.

Alternatively, you can run:

```bash
$ make install
(...)
customresourcedefinition.apiextensions.k8s.io/rockbands.music.example.io created
```

Once installed, you can check the CRDs on your cluster running:

```bash
$ kubectl get crds
NAME                         CREATED AT
rockbands.music.example.io   2020-10-29T02:37:21Z
```

Create one custom resource (CR) named "beatles" from the example [here](/single-gvk/music/config/samples/music_v1_rockband.yaml).

```bash
$ cat music_v1_rockband.yaml 
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

Let's create this CR:

```bash
$ kubectl create -f music_v1_rockband.yaml -n default
rockband.music.example.io/beatles created
```

Let's list the CR:

```bash
$ kubectl get rockband beatles -o yaml
apiVersion: music.example.io/v1
kind: RockBand
metadata:
  creationTimestamp: "2020-10-29T02:50:16Z"
  generation: 1
  managedFields:
  - apiVersion: music.example.io/v1
    fieldsType: FieldsV1
    fieldsV1:
      f:spec:
        .: {}
        f:genre: {}
        f:leadSinger: {}
        f:numberComponents: {}
    manager: kubectl
    operation: Update
    time: "2020-10-29T02:50:16Z"
  name: beatles
  namespace: default
  resourceVersion: "3167"
  selfLink: /apis/music.example.io/v1/namespaces/default/rockbands/beatles
  uid: 0a7e5845-1f79-42f9-82d1-de5d3b45999b
spec:
  genre: 60s rock
  leadSinger: John
  numberComponents: 4
```

Please note that there is no Status set. This is because we do not have yet a controller. See next section.

## Compiling and Starting the Controller

Since we do not care much about the busines logic on the controller, the only thing my controller will do is setting the lastPlayed Status field with the current year.
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
	return ctrl.Result{}, nil
```

At this time, the controller should be able to be compiled without errors running `make docker-build`.
I will build an image from the controller and name as `quay.io/brito_rafa/music-controller:single-gvk-v0.1`.

```
export IMG=quay.io/brito_rafa/music-controller:single-gvk-v0.1
make docker-build
(...)
Successfully built 1262ea2263df
Successfully tagged quay.io/brito_rafa/music-controller:single-gvk-v0.1
```

Let's test the controller before coding any webhooks.

Run:

```
export IMG=quay.io/brito_rafa/music-controller:single-gvk-v0.1
make docker-build
docker push quay.io/brito_rafa/music-controller:single-gvk-v0.1
make deploy IMG=quay.io/brito_rafa/music-controller:single-gvk-v0.1
```

You might see an error of missing namespace `music-system`:

```bash
Error from server (NotFound): error when creating "STDIN": namespaces "music-system" not found
```
For some reason kubebuilder+kustomize fails to set the namespace correctly (it creates namespace `system` instead).

Create the `music-system` namespace and try again:

```bash
$ kubectl create namespace music-system
$ make deploy IMG=quay.io/brito_rafa/music-controller:single-gvk-v0.1
(...)
deployment.apps/music-controller-manager created
```

Look at the log of the controller:

```bash
$ kubectl logs music-controller-manager-867b6f899c-7vw2q -n music-system manager
2020-10-29T03:30:44.858Z	INFO	controller-runtime.metrics	metrics server is starting to listen	{"addr": "127.0.0.1:8080"}
2020-10-29T03:30:44.858Z	INFO	setup	starting manager
I1029 03:30:44.859074       1 leaderelection.go:242] attempting to acquire leader lease  music-system/9f9c4fd4.example.io...
2020-10-29T03:30:44.859Z	INFO	controller-runtime.manager	starting metrics server	{"path": "/metrics"}
I1029 03:31:02.369750       1 leaderelection.go:252] successfully acquired lease music-system/9f9c4fd4.example.io
2020-10-29T03:31:02.370Z	DEBUG	controller-runtime.manager.events	Normal	{"object": {"kind":"ConfigMap","namespace":"music-system","name":"9f9c4fd4.example.io","uid":"727f114e-e7f3-4a79-b01f-82a2db2c1e8c","apiVersion":"v1","resourceVersion":"12506"}, "reason": "LeaderElection", "message": "music-controller-manager-867b6f899c-7vw2q_932dbf62-d492-4ca2-86f7-f9b98811f6d2 became leader"}
2020-10-29T03:31:02.370Z	INFO	controller-runtime.controller	Starting EventSource	{"controller": "rockband", "source": "kind source: /, Kind="}
2020-10-29T03:31:02.472Z	INFO	controller-runtime.controller	Starting Controller	{"controller": "rockband"}
2020-10-29T03:31:02.472Z	INFO	controller-runtime.controller	Starting workers	{"controller": "rockband", "worker count": 1}
2020-10-29T03:31:02.472Z	INFO	controllers.RockBand	received reconcile request for "beatles" (namespace: "default")	{"rockband": "default/beatles"}
2020-10-29T03:31:02.485Z	DEBUG	controller-runtime.controller	Successfully Reconciled	{"controller": "rockband", "request": "default/beatles"}
2020-10-29T03:31:02.487Z	INFO	controllers.RockBand	received reconcile request for "beatles" (namespace: "default")	{"rockband": "default/beatles"}
2020-10-29T03:31:02.487Z	DEBUG	controller-runtime.controller	Successfully Reconciled	{"controller": "rockband", "request": "default/beatles"}
```

If you see errors, please refer to the section "common errors" later this page.

Note from the logs that controller already reconciled the CR `default/beatles`.

Let's look again the CR `beatles` under `default` namespace:

```bash
$ kubectl get rockband beatles -n default -o yaml
apiVersion: music.example.io/v1
kind: RockBand
metadata:
(...)
spec:
  genre: 60s rock
  leadSinger: John
  numberComponents: 4
status:
  lastPlayed: "2020"
```

Note that we have a status field now. Let's now code mutator and validator webhooks.


## CRD Mutation and Validation

For this demo, we came up with a couple silly rules just to make a point how to mutate and validate RockBands objects.

Mutation: 
1. if LeadSinger is not specified, set it as "TBD".
2. if RockBand CR name is "beatles" and Lead Singer is "John", set it to "John Lennon" (disclaimer: some Beatles fans (including me) will argue The Beatles did not have "the" Lead Singer).

Validation:
1. We can't create RockBands CRs on "kube-system" namespace
2. We can't update Lead Singer as "John" if RockBand CR name is "beatles" (similar to mutation, but during at update time. Spoiler: this condition will never met because of the mutation logic)
3. We can't update Lead Singer as "Ringo" if RockBand CR name is  "beatles" .


### Code

Let's do the scaffolding:

```bash
kubebuilder create webhook --group music --version v1 --kind RockBand --defaulting --programmatic-validation
```

Scaffolding should have created the file  `music/api/v1/rockband_webhook.go` AND edited `main.go`.


#### main.go

Let's mak sure the following webhook call is on main.go, otherwise the webhook will not run:

```go main.go
	if err = (&musicv1.RockBand{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "RockBand")
		os.Exit(1)
	}
```

#### rockband_webhook.go

This where your mutation and validation logic will reside.

The default `music/api/v1/rockband_webhook.go` has multiple methods with kubebuilder tags.

Let's start with the mutator, which is the Default. This has priority than any other validator method and this will execute for each CR request.
Note the `mutating=true` kubebuilder tag.
This means the RockBand CR is mutable during this method. This is where you will want to setup all Spec fields with default values.

In our example below, we set leadSinger as "TBD" if empty. And it sets "John Lennon" if leadSinger is "John" and CR name is "beatles".

Note that in Default, there is no API error being returned.

```go
func (r *RockBand) Default() {
	rockbandlog.Info("mutator default", "name", r.Name, "namespace", r.Namespace)

	// TODO(user): fill in your defaulting logic.

	// LeadSinger is an optional field on RockBandv1
	// Adding "TBD" if it is empty
	if r.Spec.LeadSinger == "" {
		r.Spec.LeadSinger = "TBD"
	}

	// Silly mutation:
	// if the rockband name is beatles and leadSinger is John, set it as John Lennon
	if r.Name == "beatles" && r.Spec.LeadSinger == "John" {
		r.Spec.LeadSinger = "John Lennon"
	}
}
```

Let's look the validator calls. There are three: creation, update and deletion.

*Attention:* Note that during validation calls, one CAN'T change Spec fields (see the kubebuilder tag `mutating=false`). 

The validator methods only generate API errors.

Let's look at the creation validation, which forbids the creation of RockBand on `kube-system` namespace:

```go
func (r *RockBand) ValidateCreate() error {
	rockbandlog.Info("validate create", "name", r.Name, "namespace", r.Namespace, "lead singer", r.Spec.LeadSinger)

	// TODO(user): fill in your validation logic upon object creation.

	var allErrs field.ErrorList

	// Just an example of validation: one cannot create rockbands under kube-system namespace
	if r.Namespace == "kube-system" {
		err := field.Invalid(field.NewPath("metadata").Child("namespace"), r.Namespace, "is forbidden to have rockbands.")
		allErrs = append(allErrs, err)
	}

	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		schema.GroupKind{Group: "music.example.io", Kind: "RockBand"},
		r.Name, allErrs)
}
```

Let's look at the update validation. During update, "Ringo" is not allowed to be set as leadSinger if CR name is "beatles".
The other validation is not let the user to change the leadSinger to "John" if CR name is "beatles"

```go
func (r *RockBand) ValidateUpdate(old runtime.Object) error {
	rockbandlog.Info("validate update", "name", r.Name, "namespace", r.Namespace, "lead singer", r.Spec.LeadSinger)

	// TODO(user): fill in your validation logic upon object update.

	var allErrs field.ErrorList

	// Disclaimer: The following condition will never be met because of the Default mutation
	if r.Name == "beatles" && r.Spec.LeadSinger == "John" {
		err := field.Invalid(field.NewPath("spec").Child("leadSinger"), r.Spec.LeadSinger, "has the shortname of the singer.")
		allErrs = append(allErrs, err)
	}

	// Silly validation
	if r.Name == "beatles" && r.Spec.LeadSinger == "Ringo" {
		err := field.Invalid(field.NewPath("spec").Child("leadSinger"), r.Spec.LeadSinger, "was the drummer. Suggest you to pick John or Paul.")
		allErrs = append(allErrs, err)
	}

	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		schema.GroupKind{Group: "music.example.io", Kind: "RockBand"},
		r.Name, allErrs)
}
```

At this time, you should be able to compile the code with the webhooks.

Run:

```bash
make docker-build IMG=quay.io/brito_rafa/music-controller:single-gvk-v0.1
docker push quay.io/brito_rafa/music-controller:single-gvk-v0.1
```

#### Configuring Kubebuilder to Enable Webhook

There are multiple places that you will need to set to tell kubebuilder+kustomize to deploy the webhooks.

The main file is `config/default/kustomization.yaml` and you must uncomment all sections in regards `WEBHOOK` and `CERTMANAGER`.
The final result should look like this [kustomization.yaml](/single-gvk/music/config/default/kustomization.yaml).

```
# Adds namespace to all resources.
namespace: music-system

# Value of this field is prepended to the
# names of all resources, e.g. a deployment named
# "wordpress" becomes "alices-wordpress".
# Note that it should also match with the prefix (text before '-') of the namespace
# field above.
namePrefix: music-

# Labels to add to all resources and selectors.
#commonLabels:
#  someName: someValue

bases:
- ../crd
- ../rbac
- ../manager
# [WEBHOOK] To enable webhook, uncomment all the sections with [WEBHOOK] prefix including the one in 
# crd/kustomization.yaml
- ../webhook
# [CERTMANAGER] To enable cert-manager, uncomment all sections with 'CERTMANAGER'. 'WEBHOOK' components are required.
- ../certmanager
# [PROMETHEUS] To enable prometheus monitor, uncomment all sections with 'PROMETHEUS'. 
#- ../prometheus

patchesStrategicMerge:
  # Protect the /metrics endpoint by putting it behind auth.
  # If you want your controller-manager to expose the /metrics
  # endpoint w/o any authn/z, please comment the following line.
- manager_auth_proxy_patch.yaml

# [WEBHOOK] To enable webhook, uncomment all the sections with [WEBHOOK] prefix including the one in 
# crd/kustomization.yaml
- manager_webhook_patch.yaml

# [CERTMANAGER] To enable cert-manager, uncomment all sections with 'CERTMANAGER'.
# Uncomment 'CERTMANAGER' sections in crd/kustomization.yaml to enable the CA injection in the admission webhooks.
# 'CERTMANAGER' needs to be enabled to use ca injection
- webhookcainjection_patch.yaml

# the following config is for teaching kustomize how to do var substitution
vars:
# [CERTMANAGER] To enable cert-manager, uncomment all sections with 'CERTMANAGER' prefix.
- name: CERTIFICATE_NAMESPACE # namespace of the certificate CR
  objref:
    kind: Certificate
    group: cert-manager.io
    version: v1alpha2
    name: serving-cert # this name should match the one in certificate.yaml
  fieldref:
    fieldpath: metadata.namespace
- name: CERTIFICATE_NAME
  objref:
    kind: Certificate
    group: cert-manager.io
    version: v1alpha2
    name: serving-cert # this name should match the one in certificate.yaml
- name: SERVICE_NAMESPACE # namespace of the service
  objref:
    kind: Service
    version: v1
    name: webhook-service
  fieldref:
    fieldpath: metadata.namespace
- name: SERVICE_NAME
  objref:
    kind: Service
    version: v1
    name: webhook-service
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

I want to test the next section testing the creation of the CR, so I will delete the current "beatles" CR:

```bash
$ kubectl delete rockband beatles
rockband.music.example.io "beatles" deleted
```


## Testing the webhook

*Attention*: if you have not installed the cert-manager, the time is now.

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





## Common Errors

### Lack of certs or cert-manager

If the controller has the webhooks enabled but there is no certs from cert-manager:

```bash
$ kubectl logs music-controller-manager-78949d85d7-gtmhz -n music-system 
error: a container name must be specified for pod music-controller-manager-78949d85d7-gtmhz, choose one of: [manager kube-rbac-proxy]
MacBook-Pro:music rbrito$ kubectl logs music-controller-manager-78949d85d7-gtmhz -n music-system manager
2020-10-29T03:07:42.209Z	INFO	controller-runtime.metrics	metrics server is starting to listen	{"addr": "127.0.0.1:8080"}
2020-10-29T03:07:42.210Z	INFO	controller-runtime.builder	Registering a mutating webhook	{"GVK": "music.example.io/v1, Kind=RockBand", "path": "/mutate-music-example-io-v1-rockband"}
2020-10-29T03:07:42.211Z	INFO	controller-runtime.webhook	registering webhook	{"path": "/mutate-music-example-io-v1-rockband"}
2020-10-29T03:07:42.211Z	INFO	controller-runtime.builder	Registering a validating webhook	{"GVK": "music.example.io/v1, Kind=RockBand", "path": "/validate-music-example-io-v1-rockband"}
2020-10-29T03:07:42.211Z	INFO	controller-runtime.webhook	registering webhook	{"path": "/validate-music-example-io-v1-rockband"}
2020-10-29T03:07:42.211Z	INFO	setup	starting manager
I1029 03:07:42.212603       1 leaderelection.go:242] attempting to acquire leader lease  music-system/9f9c4fd4.example.io...
2020-10-29T03:07:42.214Z	INFO	controller-runtime.manager	starting metrics server	{"path": "/metrics"}
2020-10-29T03:07:42.308Z	INFO	controller-runtime.webhook.webhooks	starting webhook server
2020-10-29T03:07:42.309Z	DEBUG	controller-runtime.manager	non-leader-election runnable finished	{"runnable type": "*webhook.Server"}
2020-10-29T03:07:42.309Z	ERROR	setup	problem running manager	{"error": "open /tmp/k8s-webhook-server/serving-certs/tls.crt: no such file or directory"}
github.com/go-logr/zapr.(*zapLogger).Error
	/go/pkg/mod/github.com/go-logr/zapr@v0.1.0/zapr.go:128
main.main
	/workspace/main.go:85
runtime.main
	/usr/local/go/src/runtime/proc.go:203

```

Solution: troubleshoot the `make deploy` and look for the yaml file under `config/certmanager/certificate.yaml`.

### Lack of music-system namespace

Kubebuilder and kustomize do not setup the controller namespace according.
```bash
make deploy IMG=quay.io/brito_rafa/music-controller:single-gvk-v0.1
(...)
Error from server (NotFound): error when creating "STDIN": namespaces "music-system" not found
```

Solution: create the `music-system` namespace manually before deploying the controller.

### mutate=true at validation

Make deploy generated the following error:

```
Error: no matches for OriginalId admissionregistration.k8s.io_v1beta1_ValidatingWebhookConfiguration|~X|validating-webhook-configuration; no matches for CurrentId admissionregistration.k8s.io_v1beta1_ValidatingWebhookConfiguration|~X|validating-webhook-configuration; failed to find unique target for patch admissionregistration.k8s.io_v1beta1_ValidatingWebhookConfiguration|validating-webhook-configuration
```

### Controller running locally

If you are testing the webhook locally using `make run`, you will need certs under the directory `/tmp/k8s-webhook-server/serving-certs`.

Create and install the key/cert running:

```
openssl req -new -newkey rsa:4096 -x509 -sha256 -days 365 -nodes -out tls.crt -keyout tls.key
# answer all the questions
mkdir -p /tmp/k8s-webhook-server/serving-certs
mv tls.* /tmp/k8s-webhook-server/serving-certs
```

During the execution of the controller, possible the webhook will complain to a different directory than /tmp/k8s-webhook-server/serving-certs .
In this case, just copy the tsl.* files to the directory.

It is possible that running locally the webhook is never called - that was my case. My webhook only worked properly once running as a pod on my cluster.
In this case, you will need to deploy them on the cluster (see previous section).


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

Listing resources:

```
kubectl get api-resources
```

Listing preferred version:

```
kubectl get --raw /apis/music.example.io | jq -r
```
