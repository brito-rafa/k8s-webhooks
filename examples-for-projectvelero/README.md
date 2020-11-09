# Multiple API Groups for Project Velero Testing

This work is to support the following change:
https://github.com/vmware-tanzu/velero/issues/2551

Current design doc (on the PR):
https://github.com/vmware-tanzu/velero/pull/3050

For such, we need an API group with multiple versions.

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

See design doc for more details on each case.

Note that Kubernetes has a version priority list, so stable releases (v1, v2, etc) are forced to be preferred versions in detriment of beta and alpha releases:
https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definition-versioning/#version-priority


## Case A

[case-a directory](/examples-for-projectvelero/case-a/): code and instructions how to setup.

Rationale: `target preferred version == source preferred version`

RockBand on Source cluster: v1 (preferred), v1alpha1

RockBand on Target cluster: v1 (preferred), v2beta1

Expected result: `v1` being used during restore.

## Case B

[case-b directory](/examples-for-projectvelero/case-b/): code and instructions how to setup.

Rationale: `target preferred version != source preferred version; target preferred version belongs in source supported version array`

RockBand on Source cluster: v1 (preferred), v2beta2, v2beta1

RockBand on Target cluster: v2beta2 (preferred), v2beta1

Expected result: `v2beta2` being used during restore.

**Attention**: this use case is very unlikely to happen because the target preferred version should be have a stable version.

## Case C

[case-c directory](/examples-for-projectvelero/case-c/): code and instructions how to setup.

Rationale: `target preferred version != source preferred version; target preferred version does not belong in source supported version; source preferred version belongs in target supported version array`

RockBand on Source cluster: v1 (preferred), v1alpha1

RockBand on Target cluster: v2 (preferred), v1

Expected result: `v1` being used during restore.

## Case D

[case-d directory](/examples-for-projectvelero/case-d/): code and instructions how to setup.

Rationale: `target preferred version != source preferred version; target preferred version does not belong in source supported version; source preferred version does not belong in target supported version array; use intersection of supported arrays`

RockBand on Source cluster: v1 (preferred), v2beta2, v2beta1

RockBand on Target cluster: v2 (preferred), v2beta2, v2beta1

Expected result: v2beta1 and v2beta2 are common. `v2beta2` to be used during restore.
