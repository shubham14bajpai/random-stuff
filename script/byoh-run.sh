## usage: ./run.sh
## make sure the script is moved to byoh repo for this to work

#!/bin/bash

function makeSetup {
    kind create cluster
    clusterctl init --infrastructure byoh
    kubectl apply -f config/crd/bases/
}

function makeBuild {
    make docker-build
    make host-agent-binaries

    kind load docker-image gcr.io/k8s-staging-cluster-api/cluster-api-byoh-controller:dev

    kubectl -n byoh-system set image deployment/byoh-controller-manager manager=gcr.io/k8s-staging-cluster-api/cluster-api-byoh-controller:dev
    kubectl apply -f config/rbac/role.yaml
    make prepare-byoh-docker-host-image
}

function makeHosts {
    docker stop host1 host2
    docker rm host1 host2
    for i in {1..2}
    do
    echo "Creating docker container named host$i"
    docker run --detach --tty --hostname host$i --name host$i --privileged --security-opt seccomp=unconfined --tmpfs /tmp --tmpfs /run --volume /var --volume /lib/modules:/lib/modules:ro --network kind byoh/node:e2e
    done

    cp ~/.kube/config ~/.kube/management-cluster.conf
    export KIND_IP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' kind-control-plane)
    sed -i 's/    server\:.*/    server\: https\:\/\/'"$KIND_IP"'\:6443/g' ~/.kube/management-cluster.conf

    for i in {1..2}
    do
    echo "Copy agent binary to host $i"
    docker cp bin/byoh-hostagent-linux-amd64 host$i:/byoh-hostagent
    echo "Copy kubeconfig to host $i"
    docker cp ~/.kube/management-cluster.conf host$i:/management-cluster.conf
    done
}

makeSetup
makeBuild
makeHosts

echo "Run the following command in seperate shells to continue\n"

echo 'export HOST_NAME=host1
docker exec -it $HOST_NAME sh -c "chmod +x byoh-hostagent && ./byoh-hostagent --kubeconfig management-cluster.conf"'
echo "\n"

echo 'export HOST_NAME=host2
docker exec -it $HOST_NAME sh -c "chmod +x byoh-hostagent && ./byoh-hostagent --kubeconfig management-cluster.conf"'
echo "\n"

echo "Post that run the generate cluster, kubeconfig and apply cni using:\n"

echo 'kubectl get secret/byoh-cluster-kubeconfig -o json \
  | jq -r .data.value \
  | base64 --decode \
  > ./byoh-cluster.kubeconfig'
echo "\n"

echo 'KUBECONFIG=byoh-cluster.kubeconfig kubectl apply -f https://docs.projectcalico.org/v3.20/manifests/calico.yaml'
echo "\n"

