# Kubernetes Custom Resource Definitions and Webhooks

This is a repository to give a simple example in how to create a k8s Custom Resource Definition (CRD) and a controller with [webhooks](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/) using [kubebuilder](https://go.kubebuilder.io/). 

The example here is a great starting point for you to learn how to validate and mutate your CRDs across multiple [API Group Versions](https://kubernetes.io/docs/concepts/overview/kubernetes-api/#api-groups-and-versioning) and a springboard for a simple controller.


## What are We Trying to Achieve?

The code and step-by-step here are focused in how to setup from scratch the whole combination CRD/CR/webhook/controller.

One can just create CRDs and Custom Resources (CRs) on a k8s cluster, but they will do not much without a controller with reconciling code and webhooks (required for mutation and validation).

The controller itself in this example does not have any busines logic but you will learn the validation and the mutator to support multiple [API Group Versions](https://kubernetes.io/docs/concepts/overview/kubernetes-api/#api-groups-and-versioning) using kubebuilder. One use case for the mutator webhook is a CRD in production that require a new schema without causing any disruption.


I created this example to aid the coding of [Project Velero to support multiple API Groups during backup and restore](https://github.com/vmware-tanzu/velero/issues/2551).

The sample CRD is named `rockbands.music.example.io` and the controller will be called "music".

Our example sample CRD:
- domain: example.io
- group: music
- kind: rockband
- versions: v1 and v1alpha1
- `RockBandv1` : Fields `Spec.Genre`, `Spec.NumberComponents`, `Spec.LeadSinger` and `Status.LastPlayed`
- `RockBandv1alpha1` : Fields `Spec.Genre`, `Spec.NumberComponents`, and `Status.LastPlayed`

The controller will present the CRs back and forth between versions `v1` and `v1alpha1`.

```
$ kubectl get crd rockbands.music.example.io
NAME                         CREATED AT
rockbands.music.example.io   2020-10-28T19:51:25Z
```

I found easier to code this example in major two steps: 
1. "First Example": Creating the first CRD version and validator webhook. 
2. "Second Example": Creating the second API Group version and conversion webhook. 


## First Example: A single Group Version Kind (GVK)

This is the first example and we will start with one version of the group and kind: `RockBandv1`.
Please start at [README.md](/single-gvk/README.md).

## Second Example: Multiple Group Version Kind (GVK)

This is the second example and it is built upon the first example. It creates the `RockBandv1alpha1`.
Please refer at [README.md](/multiple-gvk/README.md).


## Pre-reqs for this Development and Testing

### Development Software

Please see https://github.com/embano1/codeconnect-vm-operator#developer-software 

### Kubernetes Cluster

For all examples, we will use a Kind cluster.

Please see https://github.com/embano1/codeconnect-vm-operator#kubernetes-cluster

### Other K8s Components

TLS is a critical component of webhooks. You will need cert-manager running on your K8s cluster: 
https://cert-manager.io/docs/installation/kubernetes/

Install it running:
```
kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v1.0.3/cert-manager.yaml
```

The result should look like this:

```
kubectl get pods --namespace cert-manager
NAME                                       READY   STATUS    RESTARTS   AGE
cert-manager-59dbb7958b-t4w68              1/1     Running   0          24s
cert-manager-cainjector-5df5cf79bf-j9h8m   1/1     Running   1          24s
cert-manager-webhook-8557565b68-hpp5f      1/1     Running   0          24s
```

If you want, you can use the [example-io-self-signed.yaml](/other/example-io-self-signed.yaml) file here to test if cert-manager is operational. You can delete the cert after testing it (kubebuilder creates its own.)

```
$ kubectl create -f other/example-io-self-signed.yaml 
issuer.cert-manager.io/example-io-selfsigned created
certificate.cert-manager.io/example-io-selfsigned created

$ kubectl get issuers,certificates -n kube-system
NAME                                           READY   AGE
issuer.cert-manager.io/example-io-selfsigned   True    32s

NAME                                                READY   SECRET                           AGE
certificate.cert-manager.io/example-io-selfsigned   True    example-io-selfsigned-cert-tls   32s

$ kubectl delete -f example-io-self-signed.yaml
issuer.cert-manager.io "example-io-selfsigned" deleted
certificate.cert-manager.io "example-io-selfsigned" deleted
```