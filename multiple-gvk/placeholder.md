
# Ignore this file for now

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
