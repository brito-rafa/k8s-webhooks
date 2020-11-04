kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v1.0.3/cert-manager.yaml
sleep 10
echo "INFO: Waiting for cert-manager to start..."
sleep 10
kubectl apply --validate=false -f https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-a/source/case-a-source.yaml
sleep 10
echo "INFO: Waiting for controller to start..."
sleep 10
# creating the testing CRs
kubectl create --validate=false -f https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-a/source/music/config/samples/music_v1_rockband.yaml
kubectl create --validate=false -f https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-a/source/music/config/samples/music_v1alpha1_rockband.yaml
# if you want to test
kubectl get rockbands -A -o yaml
