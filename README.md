# Kubernetes Custom Resource Definitions and Webhooks

This is a repository to give a simple example in how to create k8s Custom Resource Definitions (CRD) and webhooks using kubebuilder.

One can just create CRDs and Custom Resources (CRs) on a k8s cluster, but they will do not much without a controller with webhooks (required for mutation and validation).

The code and steps here will focus in how to setup the the whole shebang CRD/CR/webhook/controller from scratch.

## Pre-reqs

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

If you want, you can use the example-io-self-signed.yaml file here to test if cert-manager is operational.

```
$ kubectl create -f example-io-self-signed.yaml 
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


## High-Level Design of the CRD

Our CRD example will go over a life cycle, which means having multiple versions over time.
And the webhook on the controller will handle the changes across versions.

Our example:
- domain: example.io
- group: music
- kind: rockband
- versions: v1beta1, v1, v2alpha1, etc

So our CRD will be rockbands.music.example.io. The CRD will have information of a given rockband under the API group music.example.io.

```
$ kubectl get crd rockbands.music.example.io
NAME                         CREATED AT
rockbands.music.example.io   2020-10-28T19:51:25Z
```
We will talk about the specs on the examples.

## First Example: A single Group Version Kind (GVK)

This is the first example.
Start following the README.md there.
We will start with one version of the group and kind.
