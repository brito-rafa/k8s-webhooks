# How to create CRDs

## Conceptual

Define the following parameters:
- domain: example.io
- group: music
- kind: rockband
- versions: v1beta1, v1, v2alpha1, etc

References:
https://book.kubebuilder.io/multiversion-tutorial/tutorial.html


## Generating

How to create the skeleton:

```
mkdir music
go mod init music
kubebuilder init --domain example.io
kubebuilder create api --group music --version v1alpha1 --kind RockBand
kubebuilder create api --group music --version v1 --kind RockBand
kubebuilder create webhook --group music --version v1 --kind RockBand --defaulting --programmatic-validation
```

Edit the type files and pick one to be the (preferred version)[https://book.kubebuilder.io/multiversion-tutorial/api-changes.html#storage-versions]:

```
// +kubebuilder:storageversion
```

Edit Makefile and set CRD line:
```
CRD_OPTIONS ?= "crd:preserveUnknownFields=false,crdVersions=v1,trivialVersions=true"
```

Finally, then run

```
make manifests && make generate
```

## Installing on the cluster

```
make install
```

## CRD Conversion and Webhooks

Theory here:
https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/

### Code

The main.go starts the webhook server.
On your preferred version api directory, you should have rockband_webhook.go and rockband_conversion.go.
For each other supported version, you will need functions ConvertTo and ConvertFrom functions. They should be rockband_conversion.go.

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

For this demo, we came up with a couple silly rules just to make a point how to mutate and validate objects.

Mutation: 
1. if LeadSinger is not specified, set it as "TBD".
2. if RockBand name is "beatles" and Lead Singer is "John", set it to "John Lennon" (disclaimer: some Beatles fans (including me) will argue The Beatles did not have "a" Lead Singer).

Validation:
1. We can't create RockBands CRs on "kube-system" namespace
2. We can't set as "John" as Lead Singer of Object Name "beatles" (similar to mutation, but during at update time. Spoiler: this condition will never met because of the mutation)
3. We can't update the Lead Singer as "Ringo" of Object Name "beatles" 

Note that during validation, one CAN'T change fields, only generate errors - see the kubebuilder option mutating=false.

Code snippet for Mutator. Remember, this will execute for each CR request.

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
Listing specific versions of the object:

```
kubectl get rockband.v1alpha1.music.example.io beatles -o yaml
```
