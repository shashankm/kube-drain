# kube-drain
This tool allows you to pass a text file with the line seprated names of kubernetes nodes you want to drain.

It first cordones the nodes and then evicts pods on those nodes.

It then deletes the nodes thus scaling down the cluster.

# Usage
  -kubeconfig string
    	(optional) absolute path to the kubeconfig file (default "/Users/smathur/.kube/config")
  -nodefile string
    	absolute path to a text file containing the names of all nodes to be drained. Default: /etc/nodesfile (default "/tmp/nodesfile")
  -nodename string
    	Name of the node you want to drain.
      
# Example
  
  ```
kubectl get nodes
NAME                 STATUS   ROLES    AGE     VERSION
kind-control-plane   Ready    master   2m15s   v1.19.1
kind-worker          Ready    <none>   103s    v1.19.1
kind-worker2         Ready    <none>   103s    v1.19.1
kind-worker3         Ready    <none>   103s    v1.19.1
  ```
  
  ```
kubectl get pods -o wide
NAME                          READY   STATUS    RESTARTS   AGE   IP           NODE           NOMINATED NODE   READINESS GATES
hello-node-7567d9fdc9-9cfh9   1/1     Running   0          61s   10.244.2.2   kind-worker    <none>           <none>
hello-node-7567d9fdc9-qxptt   1/1     Running   0          53s   10.244.3.2   kind-worker2   <none>           <none>
  ```
  
  ```
cat hostsfile
kind-worker
kind-worker2
  ```
  
  ```
./kube-drain --nodefile hostsfile

Current context is kind-kind

Are you sure you want to drain these nodes?
kind-worker
kind-worker2

✔ Yes
Do you also want to delete these nodes from the cluster upon draining them?
✔ Yes
2021/01/13 19:09:00 Draining kind-worker
2021/01/13 19:09:00 Node kind-worker cordoned
2021/01/13 19:09:00 Evicting pod hello-node-7567d9fdc9-9cfh9
2021/01/13 19:09:05 pod hello-node-7567d9fdc9-9cfh9 evicted
2021/01/13 19:09:05 Deleting kind-worker
2021/01/13 19:09:05 Node kind-worker deleted
2021/01/13 19:09:05 Draining kind-worker2
2021/01/13 19:09:05 Node kind-worker2 cordoned
2021/01/13 19:09:05 Evicting pod hello-node-7567d9fdc9-qxptt
2021/01/13 19:09:05 pod hello-node-7567d9fdc9-qxptt evicted
2021/01/13 19:09:05 Deleting kind-worker2
2021/01/13 19:09:05 Node kind-worker2 deleted

Nodes currently available on the cluster are:
kind-control-plane
kind-worker3
  ```
  
```
kubectl get pods -o wide
NAME                          READY   STATUS    RESTARTS   AGE    IP           NODE           NOMINATED NODE   READINESS GATES
hello-node-7567d9fdc9-cdl2l   1/1     Running   0          55s    10.244.1.3   kind-worker3   <none>           <none>
hello-node-7567d9fdc9-ntjj5   1/1     Running   0          118s   10.244.1.2   kind-worker3   <none>           <none>
```
  
  you can also pass a hostname instead of a file:
  
  ```
  $ ./kube-drain --nodename minikube-control-plane
2021/01/06 21:12:15 Draining minikube-control-plane
Node minikube-control-plane cordoned
  ```
  
