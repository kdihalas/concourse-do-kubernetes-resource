package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/digitalocean/godo"
	"github.com/go-playground/validator"
	"github.com/kdihalas/concourse-do-kubernetes-resource/internal/errors"
	"github.com/kdihalas/concourse-do-kubernetes-resource/internal/types"
	"github.com/rs/xid"
)

var (
	ctx      = context.Background()
	validate = validator.New()
	guid     = xid.New()
	delete   = false

	cluster_id string
	output     types.InOutput
)

func main() {
	// set logs to stderr
	log.SetOutput(os.Stderr)

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

	// create a new DO client using the token from the source
	client := godo.NewFromToken(payload.Source.Token)

	clusters, _, err := client.Kubernetes.List(ctx, &godo.ListOptions{})
	errors.E(err)

	for _, cluster := range clusters {
		if payload.Params.Name == cluster.Name {
			cluster_id = cluster.ID
		}
	}

	// set the delete value form the delete param
	if payload.Params.Delete {
		delete = payload.Params.Delete
	}

	// if delete then destroy the cluster
	if delete {
		log.Println("Detected delete command")
		_, err := client.Kubernetes.Delete(ctx, cluster_id)
		errors.E(err)
		log.Println("Delete was successfull")
	} else {
		log.Println("Detected cluster update")
		pools, _, err := client.Kubernetes.ListNodePools(ctx, cluster_id, &godo.ListOptions{})
		errors.E(err)
		for _, pool := range payload.Params.NodePools {
			for _, npool := range pools {
				if pool.Name == npool.Name {
					log.Printf("Updating pool %s with ID: %s", npool.Name, npool.ID)
					_, _, err := client.Kubernetes.UpdateNodePool(ctx, cluster_id, npool.ID, &godo.KubernetesNodePoolUpdateRequest{
						Count: &pool.Count,
					})
					errors.E(err)
					log.Println("Update was successfull")
				}
			}
		}
	}
	// prepare output
	output.Version.Ref = guid.String()
	output.Metadata = append(output.Metadata, types.NameValue{
		Name:  "cluster_id",
		Value: cluster_id,
	})

	// parse output to json string
	result, err := json.Marshal(output)
	errors.E(err)

	// print the results
	fmt.Println(string(result))

}
