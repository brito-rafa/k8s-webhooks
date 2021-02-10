# Multiple API Groups for Project Velero Testing

This work is to support the following change:
https://github.com/vmware-tanzu/velero/issues/2551

Pull requests:

https://github.com/vmware-tanzu/velero/pull/2373 # backup

https://github.com/vmware-tanzu/velero/pull/3133 # restore


Design doc:
https://github.com/vmware-tanzu/velero/pull/3050

For such, we need an API group with multiple versions. For the target clusters, please use k8s 1.18 and later.

Evolution of RockBand.music.example.io schema across API Group versions:

- `RockBandv1alpha1` : Fields `Spec.Genre`, `Spec.NumberComponents`
- `RockBandv1` : all previous plus `Spec.LeadSinger`
- `RockBandv2beta1` : all previous plus `Spec.LeadGuitar`
- `RockBandv2beta2` : all previous plus `Spec.Drummer`
- `RockBandv2` : all previous plus `Spec.Bass`

Example of the yaml:

```yaml music_v2_rockband.yaml
apiVersion: music.example.io/v2
kind: RockBand
metadata:
  name: beatles
spec:
  # Add fields here
  genre: '60s rock'
  numberComponents: 4
  leadSinger: John
  leadGuitar: George
  drummer: Ringo
  bass: Paul
```

Running the controller and installing the music.example.io API Groups:

```bash
# installing the controller & API Group - using case C target: RockBandv2 (preferred) & RockBandv1 (supported)
$ curl -k -s https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-c/target-cluster.sh | bash
(...)
NAME                                        READY   STATUS    RESTARTS   AGE
music-controller-manager-84d898799b-7xd54   2/2     Running   0          11s
```

Now, creating a RockBand example:

```bash
$ kubectl create --validate=false -f https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-c/target/music/config/samples/music_v2_rockband.yaml 
rockband.music.example.io/beatles created
```

Showing the Custom Resource (CR) created across the two versions. The field mutations and conversions were for testing / debugging purposes:

```yaml
# $ kubectl get rockband beatles -o yaml
apiVersion: music.example.io/v2
kind: RockBand
metadata:
  creationTimestamp: "2020-11-09T16:30:42Z"
  generation: 1
  name: beatles
  namespace: default
  resourceVersion: "2178"
  selfLink: /apis/music.example.io/v2/namespaces/default/rockbands/beatles
  uid: 51ea2460-7015-4605-96f3-4d23cd72eae7
spec:
  bass: Paul McCartney
  drummer: Ringo Starr
  genre: 60s rock
  leadGuitar: George Harrison
  leadSinger: John Lennon
  numberComponents: 4
status:
  lastPlayed: "2020"

# $ kubectl get rockband.v1.music.example.io beatles -o yaml
apiVersion: music.example.io/v1
kind: RockBand
metadata:
  annotations:
    rockbands.v2.music.example.io/bass: Paul McCartney
    rockbands.v2beta1.music.example.io/leadGuitar: George Harrison
    rockbands.v2beta2.music.example.io/drummer: Ringo Starr
  creationTimestamp: "2020-11-09T16:30:42Z"
  generation: 1
  name: beatles
  namespace: default
  resourceVersion: "2178"
  selfLink: /apis/music.example.io/v1/namespaces/default/rockbands/beatles
  uid: 51ea2460-7015-4605-96f3-4d23cd72eae7
spec:
  genre: 60s rock
  leadSinger: John Lennon
  numberComponents: 4
status:
  lastPlayed: "2020"
```

We will test 4 cases with different combination of API group versions (each case will be a subdirectory).

If you want to learn how to create this custom controller and API Group, recommend you to follow the [README.md here.](/README.md)

## Cases

As per Velero current design, it will try the following priority list to use the API group version for restore:

- **Priority 0** (User override). Users determine restore version priority using a config map. To test this use case, one can override case d.
- **Priority 1**. Target preferred version can be used.
- **Priority 2**. Source preferred version can be used.
- **Priority 3**. A common supported version can be used. This means
  - target supported version == source supported version
  - if multiple support versions intersect, choose the version using the [Kubernetes’ version prioritization system](https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definition-versioning/#version-priority)

Note that Kubernetes has a version priority list, so stable releases (v1, v2, etc) are forced to be preferred versions in detriment of beta and alpha releases:
https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definition-versioning/#version-priority


## Case A - Priority 1

[case-a directory](/examples-for-projectvelero/case-a/): code and instructions how to setup.

Rationale: `target preferred version == source preferred version`

RockBand on Source cluster: v1 (preferred), v1alpha1

RockBand on Target cluster: v1 (preferred), v2beta1

Expected result: `v1` being used during restore.

## Case B - Priority 1

[case-b directory](/examples-for-projectvelero/case-b/): code and instructions how to setup.

**Attention**: this use case is very unlikely to happen because the target preferred version should be a stable (v1, v2, etc) version.

Rationale: `target preferred version != source preferred version; target preferred version belongs in source supported version array`

RockBand on Source cluster: v1 (preferred), v2beta2, v2beta1

RockBand on Target cluster: v2beta2 (preferred), v2beta1

Expected result: `v2beta2` being used during restore.


## Case C - Priority 2

[case-c directory](/examples-for-projectvelero/case-c/): code and instructions how to setup.

Rationale: `target preferred version != source preferred version; target preferred version does not belong in source supported version; source preferred version belongs in target supported version array`

RockBand on Source cluster: v1 (preferred), v1alpha1

RockBand on Target cluster: v2 (preferred), v1

Expected result: `v1` being used during restore.

## Case D - Priority 3

[case-d directory](/examples-for-projectvelero/case-d/): code and instructions how to setup.

Rationale: `target preferred version != source preferred version; target preferred version does not belong in source supported version; source preferred version does not belong in target supported version array; use intersection of supported arrays`

RockBand on Source cluster: v1 (preferred), v2beta2, v2beta1

RockBand on Target cluster: v2 (preferred), v2beta2, v2beta1

Expected result: v2beta1 and v2beta2 are common. `v2beta2` to be used during restore. 

## Case D - Priority 0 with User Override

To test priority 0 (user-defined), one can use the same as Case D, but configuring velero to use v2beta1 instead of v2beta2.

For such, create a configmap with the following content (note the sequence of the API Group versions).

```cm
rockbands.music.example.io=v2beta1,v2beta2
```

 `kubectl create configmap enableapigroupversions --from-file=<absolute path>/restoreResourcesVersionPriority -n velero`

## Known Issues

If you encounter the following errors, make sure you run the controller on a k8s cluster 1.18 and later.

```
mutatingwebhookconfiguration.admissionregistration.k8s.io/music-mutating-webhook-configuration created
validatingwebhookconfiguration.admissionregistration.k8s.io/music-validating-webhook-configuration created
Error from server (InternalError): error when creating "https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-a/target/case-a-target.yaml": Internal error occurred: failed calling webhook "webhook.cert-manager.io": Post https://cert-manager-webhook.cert-manager.svc:443/mutate?timeout=10s: dial tcp 10.96.80.208:443: connect: connection refused
Error from server (InternalError): error when creating "https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-a/target/case-a-target.yaml": Internal error occurred: failed calling webhook "webhook.cert-manager.io": Post https://cert-manager-webhook.cert-manager.svc:443/mutate?timeout=10s: dial tcp 10.96.80.208:443: connect: connection refused
```
