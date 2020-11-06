# Multiple API Groups for Project Velero Testing

This work is to support the following change:
https://github.com/vmware-tanzu/velero/issues/2551

Current design doc (on the PR):
https://github.com/vmware-tanzu/velero/pull/3050

For such, we need an API group with multiple versions.

Evolution of RockBand.music.example.io schema across versions:

- `RockBandv1alpha1` : Fields `Spec.Genre`, `Spec.NumberComponents`
- `RockBandv1` : all previous plus `Spec.LeadSinger`
- `RockBandv2beta1` : all previous plus `Spec.LeadGuitar`
- `RockBandv2beta2` : all previous plus `Spec.Drummer`
- `RockBandv2` : all previous plus `Spec.Bass`

We will test 4 cases, each case will be a subdirectory here.

Kubernetes has a version priority list:
https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definition-versioning/#version-priority

Kubernetes will always use this priority to setup preferred versions of the CRDs.

## Cases

See design doc for more details.

## Case A

Rationale: `target preferred version == source preferred version`

RockBand on Source cluster: v1 (preferred), v1alpha1

RockBand on Target cluster: v1 (preferred), v2beta1

Expected result: v1 being used during restore.

## Case B

Rationale: `target preferred version != source preferred version; target preferred version belongs in source supported version array`

RockBand on Source cluster: v1 (preferred), v2beta2, v2beta1

RockBand on Target cluster: v2beta2 (preferred), v2beta1

Expected result: v2beta2 being used during restore.

Attention: this use case is very unlikely to happen because the target preferred version should be have a stable version.

## Case C

Rationale: `target preferred version != source preferred version; target preferred version does not belong in source supported version; source preferred version belongs in target supported version array`

RockBand on Source cluster: v1 (preferred), v1alpha1

RockBand on Target cluster: v2 (preferred), v1

Expected result: v1 being used during restore.

## Case D

Rationale: `target preferred version != source preferred version; target preferred version does not belong in source supported version; source preferred version does not belong in target supported version array; use intersection of supported arrays`

RockBand on Source cluster: v1 (preferred), v2beta2, v2beta1

RockBand on Target cluster: v2 (preferred), v2beta2, v2beta1

Expected result: v2beta1 and v2beta2 are common. v2beta2 to be used during restore.
