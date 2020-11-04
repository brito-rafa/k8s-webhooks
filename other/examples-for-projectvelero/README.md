# Multiple API Groups for Project Velero Testing

This work is to support the following change:
https://github.com/vmware-tanzu/velero/issues/2551

For such, we need an API group with multiple versions.

Evolution of RockBand.music.example.io schema across versions:

- `RockBandv1alpha1` : Fields `Spec.Genre`, `Spec.NumberComponents`
- `RockBandv1` : all previous plus `Spec.LeadSinger`
- `RockBandv2alpha1` : all previous plus `Spec.LeadGuitar`
- `RockBandv2beta1` : all previous plus `Spec.Bass`
- `RockBandv2` : all previous plus `Spec.Drummer`
- `RockBandv3alpha1` : all previous but `Spec.NumberComponents` = int64


We will test 4 cases, each case will be a subdirectory here.

## Cases

See design doc for more details.

## Case A

RockBand on Source cluster: v1 (preferred), v1alpha1

RockBand on Target cluster: v1 (preferred), v2beta1

Expected result: v1 being used during restore.

## Case B

RockBand on Source cluster: v1 (preferred), v2beta1, v2

RockBand on Target cluster: v2 (preferred), v1, v2beta1

Expected result: v2 being used during restore.

## Case C

RockBand on Source cluster: v1 (preferred), v1alpha1, v2beta1

RockBand on Target cluster: v2 (preferred), v1, v2beta1

Expected result: v1 being used during restore.

## Case D

RockBand on Source cluster: v1 (preferred), v2beta1, v2alpha1, v1alpha1

RockBand on Target cluster: v2 (preferred), v3alpha1, v2beta1, v2alpha1

Expected result: v2beta1 and v2alpha1 are common and v2beta1 to be used during restore.

