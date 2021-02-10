# Workshop on API Group Versions and Webhooks

This is a placeholder if there is the need of a workshop on Kubernetes API Group Versions and kubebuilder with webhooks.


## Part 1: Kubernetes API Group Versions

- Audience Requirements: basic to medium understanding on Kubernetes. What is a kind cluster, control planel, worker nodes, service, pods, deployment, etc.
- What is an API Group Version?
- Step back: Native K8s objects versus CRDs
- Slide deck (just for reference): https://docs.google.com/presentation/d/1PszHcfRs_o02Azsb98pdg_o-w19i88bgc9pcOA92_TM/edit?ts=5e822888#slide=id.g829b865cd1_2_12
- `kubectl api-resources`
- `kubectl get crds`
- `kubectl get --raw /apis | jq`
- Native K8s example of Multi API group versions: HorizontalPodAutoscaler
- Resource belongs to an API group. example: horizontalautoscaling (resource) belongs to autoscaling (group)
- `kubectl get --raw /apis/autoscaling | jq`
- `kubectl apply --validate=false -f https://raw.githubusercontent.com/brito-rafa/velero-test/master/myexample-test.yaml`
- kubectl get hpa -A # it creates a v2beta2 of hpa but the command gets me v1
- Introduce the concept of preferred + stable version for Kubernetes API server
- API deprecation: https://kubernetes.io/docs/reference/using-api/deprecation-policy/
- API Version Priority: https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definition-versioning/#version-priority



## Part 2: Real world problem with API Group Versions
- Business Problem: you have deployed an API Group Version of your app in production but now you need to change the schema
- Example of hpa v2beta2 and v1 on autoscaling
- kubectl get hpa php-apache-autoscaler -n myexample -o yaml # v1
- kubectl get hpa.v2beta2.autoscaling php-apache-autoscaler -n myexample -o yaml
- Rockbands schema
- Concept of webhook
- Supporting multiple API Group Versions
- Multiple examples for Velero: https://github.com/brito-rafa/k8s-webhooks/blob/master/examples-for-projectvelero/README.md
- Going over the step-by-step on the code if there is time
- If there is no time:
- deploy case C target and show the creation of objects
- Show the type and webhook code.