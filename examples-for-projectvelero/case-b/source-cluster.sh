kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v1.0.3/cert-manager.yaml
sleep 20
echo "INFO: Waiting for cert-manager to start..."
sleep 20
kubectl apply --validate=false -f https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-b/source/case-b-source-manually-added-mutations.yaml
sleep 10
echo "INFO: Waiting for controller to start..."
sleep 10
# creating the testing CRs
kubectl create --validate=false -f https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-b/source/music/config/samples/music_v2beta2_rockband.yaml
kubectl create --validate=false -f https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-b/source/music/config/samples/music_v2beta1_rockband.yaml
# if you want to test
kubectl get rockbands -A -o yaml
