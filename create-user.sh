#!/bin/bash
# Add user to k8s 1.5 using service account, no RBAC (unsafe)

if [[ -z "$1" ]] ;then
  echo "usage: $0 <username>"
  exit 1
fi

user=$1
kubectl create sa ${user} -n ${user}
secret=$(kubectl get sa ${user} -n ${user} -o json | jq -r .secrets[].name)
echo "secret = ${secret}"

kubectl create rolebinding ${user}-admin \
  --clusterrole=cluster-admin \
  --serviceaccount=${user}:${user} \
  --namespace=${user}

kubectl create rolebinding ${user}-view \
  --clusterrole=view\
  --serviceaccount=${user}:${user}

kubectl get secret ${secret} -n ${user} -o json | jq -r '.data["ca.crt"]' | base64 -D > users/${user}.ca.crt
user_token=$(kubectl get secret ${secret} -n ${user} -o json | jq -r '.data["token"]' | base64 -D)
echo "token = ${user_token}"

c=`kubectl config current-context`
echo "context = $c"

cluster_name=`kubectl config get-contexts $c | awk '{print $3}' | tail -n 1`
echo "cluster_name= ${cluster_name}"

endpoint=`kubectl config view -o jsonpath="{.clusters[?(@.name == \"${cluster_name}\")].cluster.server}"`
echo "endpoint = ${endpoint}"

# Set up the config
KUBECONFIG=k8s-${user}-conf kubectl config set-cluster ${cluster_name} \
    --embed-certs=true \
    --server=${endpoint} \
    --certificate-authority=./users/${user}.ca.crt
echo ">>>>>>>>>>>>ca.crt"
cat ./users/${user}.ca.crt
echo "<<<<<<<<<<<<ca.crt"
echo ">>>>>>>>>>>>${user}-setup.sh"
echo kubectl config set-cluster ${cluster_name} \
    --embed-certs=true \
    --server=${endpoint} \
    --certificate-authority=./users/${user}.ca.crt
echo kubectl config set-credentials ${user}-${cluster_name#cluster-} --token=${user_token}
echo kubectl config set-context ${user}-${cluster_name#cluster-} \
    --cluster=${cluster_name} \
    --user=${user}-${cluster_name#cluster-}
echo kubectl config use-context ${user}-${cluster_name#cluster-}
echo "<<<<<<<<<<<<${user}-setup.sh"

echo "...preparing k8s-${user}-conf"
KUBECONFIG=k8s-${user}-conf kubectl config set-credentials ${user}-${cluster_name#cluster-} --token=${user_token}
KUBECONFIG=k8s-${user}-conf kubectl config set-context ${user}-${cluster_name#cluster-} \
    --cluster=${cluster_name} \
    --user=${user}-${cluster_name#cluster-}
KUBECONFIG=k8s-${user}-conf kubectl config use-context ${user}-${cluster_name#cluster-}

echo "done! Test with: "
echo "KUBECONFIG=k8s-${user}-conf kubectl get no"