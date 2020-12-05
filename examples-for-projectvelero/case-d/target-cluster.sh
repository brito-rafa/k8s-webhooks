kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v1.0.3/cert-manager.yaml
while  [ "$(kubectl get pods -n cert-manager | grep -i running | grep '1/1' |  wc -l | awk '{print $1}')" != "3" ]; do echo "INFO: Waiting cert-manager..." && kubectl get pods -n cert-manager && sleep 10 ; done
kubectl apply --validate=false -f https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-d/target/case-d-target-manually-added-mutations.yaml
while  [ "$(kubectl get pods -n music-system | grep -i running | grep '2/2' |  wc -l | awk '{print $1}')" != "1" ]; do echo "INFO: Waiting music-system...  Break if it is taking too long..." && kubectl get pods -n music-system && sleep 10 ; done
kubectl get pods -n music-system
echo "INFO: Run a Velero Restore or create the testing CRs running:"
echo "    # Testing v2
echo "    kubectl create namespace rockbands-v2
echo "    kubectl create -n rockbands-v2 --validate=false -f https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-d/target/music/config/samples/music_v2_rockband.yaml"
echo "    # Testing v2beta2
echo "    kubectl create namespace rockbands-v2beta2
echo "    kubectl create -n rockbands-v2beta2 --validate=false -f https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-d/target/music/config/samples/music_v2beta2_rockband.yaml"
echo "    # Testing v2beta1
echo "    kubectl create namespace rockbands-v2beta1
echo "    kubectl create -n rockbands-v2beta1 --validate=false -f https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-d/target/music/config/samples/music_v2beta1_rockband.yaml"
