kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v1.0.3/cert-manager.yaml
sleep 10
echo "INFO: Waiting for cert-manager to start..."
sleep 10
kubectl apply --validate=false -f https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-a/target/case-a-target.yaml
sleep 10
echo "INFO: Waiting for controller to start..."
sleep 10
kubectl get pods -n music-system
echo "INFO: Run a Velero Restore or create the testing CRs running:"
echo "    kubectl create --validate=false -f https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-a/target/music/config/samples/music_v1_rockband.yaml"
echo "    kubectl create --validate=false -f https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-a/target/music/config/samples/music_v2beta1_rockband.yaml"
