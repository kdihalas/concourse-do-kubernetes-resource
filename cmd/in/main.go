package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/digitalocean/godo"
	"github.com/go-playground/validator"
	"github.com/kdihalas/concourse-do-kubernetes-resource/internal/errors"
	"github.com/kdihalas/concourse-do-kubernetes-resource/internal/types"
	"github.com/rs/xid"
)

var (
	ctx      = context.Background()
	validate = validator.New()
	skip     = false
	guid     = xid.New()

	output types.InOutput
)

func main() {
	// set logs to stderr
	log.SetOutput(os.Stderr)
	// read the target
	target := os.Args[1]
	// create a decoder from the stdin
	decoder := json.NewDecoder(os.Stdin)
	decoder.DisallowUnknownFields()

	// parse payload as JSON
	var payload = types.Payload{}
	err := decoder.Decode(&payload)
	errors.E(err)

	// validate the json payload
	err = validate.Struct(payload)
	errors.E(err)

	// set the params.skip
	if payload.Params.Skip {
		skip = payload.Params.Skip
	}

	// create a new DO client using the token from the source
	client := godo.NewFromToken(payload.Source.Token)

	if !skip {
		// convert the source and params to KubernetesClusterCreateRequest
		clusterReq := &godo.KubernetesClusterCreateRequest{
			Name:        payload.Params.Name,
			RegionSlug:  payload.Source.Region,
			VersionSlug: payload.Params.Version,
			Tags:        payload.Source.Tags,
			NodePools:   payload.Source.NodePools,
		}

		log.Println("Creating the kubernetes cluster")
		// create the kubernetes cluster using the config
		cluster, _, err := client.Kubernetes.Create(ctx, clusterReq)
		errors.E(err)
		// create a file for the cluster id
		clusterIdTarget, err := os.Create(fmt.Sprintf("%s/cluster_id", target))
		defer clusterIdTarget.Close()
		errors.E(err)
		// write the cluster id to the target path
		clusterIdTarget.Write([]byte(cluster.ID))
		log.Println("Waiting 30 seconds for kubernetes cluster to start")
		// check every 30 seconds if the cluster is ready
		for {
			time.Sleep(30 * time.Second)
			cluster, _, _ = client.Kubernetes.Get(ctx, cluster.ID)
			if cluster.Status.State == "running" {
				goto CLUSTER_STATE
			}
			log.Println("Not ready yet, waiting another 30 seconds")
		}

	CLUSTER_STATE:
		// jump out of the loop when the cluster is ready
		log.Println("Kubernetse cluster creation done!")
		log.Println("Downloading kubeconfig")
		// get the kubeconfig file for the cluster
		kubeConfig, _, err := client.Kubernetes.GetKubeConfig(ctx, cluster.ID)
		errors.E(err)
		// create a file in the target path for the kubeconfig
		kubeconfigTarget, err := os.Create(fmt.Sprintf("%s/kubeconfig", target))
		defer kubeconfigTarget.Close()
		errors.E(err)
		kubeconfigTarget.Write(kubeConfig.KubeconfigYAML)

		output.Version.Ref = guid.String()
		output.Metadata = append(output.Metadata, types.NameValue{
			Name:  "cluster_id",
			Value: cluster.ID,
		})
	} else {
		output.Version.Ref = guid.String()
	}
	// parse output to json string
	result, err := json.Marshal(output)
	errors.E(err)

	// print the results
	fmt.Println(string(result))

}
