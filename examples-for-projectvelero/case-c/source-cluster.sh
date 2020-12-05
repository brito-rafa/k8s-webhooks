kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v1.0.3/cert-manager.yaml
while  [ "$(kubectl get pods -n cert-manager | grep -i running | grep '1/1' |  wc -l | awk '{print $1}')" != "3" ]; do echo "INFO: Waiting cert-manager..." && kubectl get pods -n cert-manager && sleep 10 ; done
kubectl apply --validate=false -f https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-a/source/case-a-source.yaml
while  [ "$(kubectl get pods -n music-system | grep -i running | grep '2/2' |  wc -l | awk '{print $1}')" != "1" ]; do echo "INFO: Waiting music-system..." && kubectl get pods -n music-system && sleep 10 ; done
# creating the testing CRs
kubectl create namespace rockbands-v1
kubectl create --validate=false -n rockbands-v1 -f https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-a/source/music/config/samples/music_v1_rockband.yaml
kubectl create namespace rockbands-v1alpha1
kubectl create --validate=false -n rockbands-v1alpha1 -f https://raw.githubusercontent.com/brito-rafa/k8s-webhooks/master/examples-for-projectvelero/case-a/source/music/config/samples/music_v1alpha1_rockband.yaml
# Displaying all CRs
kubectl get rockbands -A -o yaml
