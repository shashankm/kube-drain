# kube-drain
This tool allows you to pass a text file with the line seprated names of kubernetes nodes you want to drain.

It first cordones the nodes and then evicts pods on those nodes.

It then deletes the nodes thus scaling down the cluster.
