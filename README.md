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
  $ k get po
+ kubectl get po
NAME   READY   STATUS    RESTARTS   AGE
sise   1/1     Running   0          10s
  ```
  
  ```
  $ ./kube-drain --nodefile hostsfile
2021/01/06 21:08:35 Draining minikube-control-plane
Node minikube-control-plane cordoned
Deleting pod sise
  ```
  
  ```
  $ k get po
+ kubectl get po
No resources found in default namespace.
  ```
  
  ```
  + kubectl get nodes
NAME                     STATUS                     ROLES    AGE   VERSION
minikube-control-plane   Ready,SchedulingDisabled   master   16d   v1.19.1
  ```
  
  you can also pass a hostname instead of a file:
  
  ```
  $ ./kube-drain --nodename minikube-control-plane
2021/01/06 21:12:15 Draining minikube-control-plane
Node minikube-control-plane cordoned
  ```
  
