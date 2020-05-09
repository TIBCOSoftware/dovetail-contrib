#!/bin/bash
# delete chaincode docker container and images
# Usage: ./cleanup [view | delete]
# use 'view' to preview which docker container and image will be deleted

CMD=${1:-"view"}

# printPrivilegedPodYaml <nodeName>
# e.g. printPrivilegedPodYaml aks-fab-27733545-vmss000000
function printPrivilegedPodYaml {
  echo "
apiVersion: v1
kind: Pod
metadata:
  name: privileged-pod
  namespace: default
spec:
  containers:
  - name: busybox
    image: busybox
    resources:
      limits:
        cpu: 200m
        memory: 100Mi
      requests:
        cpu: 100m
        memory: 50Mi
    stdin: true
    securityContext:
      privileged: true
    volumeMounts:
    - name: host-root-volume
      mountPath: /host
      readOnly: true
  volumes:
  - name: host-root-volume
    hostPath:
      path: /
  nodeSelector:
    kubernetes.io/hostname: ${1}
  hostNetwork: true
  hostPID: true
  restartPolicy: Never"
}

# cleanup all Fabric chaincode docker containers and images
function cleanup {
  local peerNodes=$(kubectl get pods -o=jsonpath='{.items[?(@.metadata.labels.app=="peer")].spec.nodeName}')
  local nodes=(${peerNodes})
  for n in "${nodes[@]}"; do
    echo "on node ${n}"
    printPrivilegedPodYaml ${n} > privileged-pod.yaml
    kubectl create -f privileged-pod.yaml

    # wait until pod is running
    local stat=$(kubectl -n default get pod privileged-pod -o jsonpath='{.status.phase}')
    until [ "${stat}" == "Running" ]; do
      echo "wait 5s for privileged-pod ..."
      sleep 5
      stat=$(kubectl -n default get pod privileged-pod -o jsonpath='{.status.phase}')
    done

    # force remove chaincode containers, i.e. docker rm -f
    local cid=$(kubectl -n default exec privileged-pod -- sh -c "chroot /host/ docker ps -q -f 'label=org.hyperledger.fabric.chaincode.id.name'")
    if [ "${CMD}" == "delete" ]; then
      echo "removing chaincode container ${cid}"
      kubectl -n default exec privileged-pod -- sh -c "chroot /host/ docker rm -f ${cid}"
    else
      echo "use 'delete' to remove chaincode container ${cid}"
      kubectl -n default exec privileged-pod -- sh -c "chroot /host/ docker inspect -f '{{.Name}} {{.State.Status}}' ${cid}"
    fi

    local iid=$(kubectl -n default exec privileged-pod -- sh -c "chroot /host/ docker images -q -f 'label=org.hyperledger.fabric.chaincode.id.name'")
    if [ "${CMD}" == "delete" ]; then
      echo "removing chaincode image ${iid}"
      kubectl -n default exec privileged-pod -- sh -c "chroot /host/ docker rmi -f ${iid}"
    else
      echo "use 'delete' to remove chaincode image ${iid}"
      kubectl -n default exec privileged-pod -- sh -c "chroot /host/ docker image inspect -f '{{.RepoTags}}' ${iid}"
    fi

    kubectl -n default delete pod privileged-pod
    rm privileged-pod.yaml
  done
}

cleanup