# Kubernetes Custom Resource Definitions and Webhooks

This is a repository to give a simple example in how to create k8s Custom Resource Definitions (CRD) and webhooks using kubebuilder.
One can just create CRDs and Custom Resources (CRs) on a k8s cluster, but they will do not much without a controller with webhooks (required for mutation and validation).
The code and steps here will focus in how to setup the the whole shebang CRD/CR/webhook/controller from scratch.

## Pre-reqs

### Development Software

Please see https://github.com/embano1/codeconnect-vm-operator#developer-software 

### Kubernetes Cluster

Please see https://github.com/embano1/codeconnect-vm-operator#kubernetes-cluster

### Other K8s Components

On your k8s cluster where your controller will run, you will need cert-manager. 
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

## High-Level Design of the CRD

We will use an example of a CRD that will go over life cycle, which means having multiple versions over time.
The webhook on the controller will handle the changes across versions.

Our example will use these following parameters:
- domain: example.io
- group: music
- kind: rockband
- versions: v1beta1, v1, v2alpha1, etc

## First Example: A single Group Version Kind (GVK)

This is the first example and start with this directory.
There is one group, one version and one kind.
