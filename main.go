package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	policyv1beta1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	nodename := flag.String("nodename", "", "Name of the node you want to drain.")
	nodesfile := flag.String("nodefile", "/tmp/nodesfile", "absolute path to a text file containing the names of all nodes to be drained. Default: /etc/nodesfile")
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	if *nodename != "" {
		node := *nodename
		log.Printf("Draining %s", node)
		k8NodeCordon(node, clientset)
		// deleteNodePods(node, clientset)
		evictNodePods(node, clientset)
		deleteNode(node, clientset)
	} else {
		file, err := os.Open(*nodesfile)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			node := scanner.Text()
			log.Printf("Draining %s", node)
			k8NodeCordon(node, clientset)
			// deleteNodePods(node, clientset)
			evictNodePods(node, clientset)
			deleteNode(node, clientset)
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}
}

func k8NodeCordon(nodeInstance string, clientSet *kubernetes.Clientset) {
	type patchStringValue struct {
		Op    string `json:"op"`
		Path  string `json:"path"`
		Value bool   `json:"value"`
	}

	payload := []patchStringValue{{
		Op:    "replace",
		Path:  "/spec/unschedulable",
		Value: true,
	}}
	payloadBytes, _ := json.Marshal(payload)

	_, err := clientSet.CoreV1().Nodes().Patch(context.TODO(), nodeInstance, types.JSONPatchType, payloadBytes, metav1.PatchOptions{})
	if err != nil {
		fmt.Printf("Node %s couldn't be cordoned\n", nodeInstance)
		panic(err.Error())
	} else {
		fmt.Printf("Node %s cordoned\n", nodeInstance)
	}
}

// func deleteNodePods(nodeInstance string, client *kubernetes.Clientset) {
// 	pods, err := client.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{
// 		FieldSelector: "spec.nodeName=" + nodeInstance,
// 	})

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	for _, i := range pods.Items {
// 		if i.Namespace == "kube-system" || i.Namespace == "local-path-storage" {
// 			continue
// 		} else {
// 			fmt.Printf("Deleting pod %s\n", i.Name)
// 			err := client.CoreV1().Pods(i.Namespace).Delete(context.TODO(), i.Name, metav1.DeleteOptions{})
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 		}
// 	}
// }

func evictNodePods(nodeInstance string, client *kubernetes.Clientset) {
	pods, err := client.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{
		FieldSelector: "spec.nodeName=" + nodeInstance,
	})

	if err != nil {
		log.Fatal(err)
	}

	for _, i := range pods.Items {
		if i.Namespace == "kube-system" || i.Namespace == "local-path-storage" {
			continue
		} else {
			eviction := &policyv1beta1.Eviction{
				// TypeMeta: metav1.TypeMeta{
				// 	APIVersion: policyGroupVersion,
				// 	Kind:       EvictionKind,
				// },
				ObjectMeta: metav1.ObjectMeta{
					Name:      i.Name,
					Namespace: i.Namespace,
				},
				// DeleteOptions: &metav1.DeleteOptions{
				// 	GracePeriodSeconds: &gracePeriodSeconds,
				// },
			}
			fmt.Printf("Evicting pod %s\n", i.Name)
			// err := client.PolicyV1beta1().Evictions(i.Namespace).Evict(context.TODO(), eviction)
			for {
				err := client.PolicyV1beta1().Evictions(i.Namespace).Evict(context.TODO(), eviction)
				if err != nil {
					log.Printf("pod %s evicted\n", i.Name)
					break
				}
				// log.Println(err)
				// log.Println("Retrying")
				time.Sleep(5 * time.Second)
			}
		}
		// log.Printf("pod %s evicted\n", i.Name)
	}
}

func deleteNode(nodeInstance string, client *kubernetes.Clientset) {
	err := client.CoreV1().Nodes().Delete(context.TODO(), nodeInstance, metav1.DeleteOptions{})
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Node %s deleted\n", nodeInstance)
	}
}
