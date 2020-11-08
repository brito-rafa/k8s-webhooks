kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v1.0.3/cert-manager.yaml
while  [ "$(kubectl get pods -n cert-manager | grep -i running | grep '1/1' |  wc -l | awk '{print $1}')" != "3" ]; do echo "INFO: Waiting cert-manager..." && kubectl get pods -n cert-manager && sleep 10 ; done
kubectl apply --validate=false -f https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-d/target/case-d-target-manually-added-mutations.yaml
while  [ "$(kubectl get pods -n music-system | grep -i running | grep '2/2' |  wc -l | awk '{print $1}')" != "1" ]; do echo "INFO: Waiting music-system...  Break if it is taking too long..." && kubectl get pods -n music-system && sleep 10 ; done
kubectl get pods -n music-system
echo "INFO: Run a Velero Restore or create the testing CRs running:"
echo "    kubectl create --validate=false -f https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-c/target/music/config/samples/music_v1_rockband.yaml"
echo "    kubectl create --validate=false -f https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-c/target/music/config/samples/music_v2_rockband.yaml"
