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
)

func main() {
	// set logs to stderr
	log.SetOutput(os.Stderr)

	//  create a unique refID
	refID := guid.String()

	// create a decoder from the stdin
	decoder := json.NewDecoder(os.Stdin)
	decoder.DisallowUnknownFields()

	// parse payload as JSON
	var payload = types.Payload{}
	err := decoder.Decode(&payload)
	errors.E(err)

	err = validate.Struct(payload)
	errors.E(err)

	client := godo.NewFromToken(payload.Source.Token)
	_, _, err = client.Kubernetes.List(ctx, &godo.ListOptions{})
	errors.E(err)

	var output []types.CheckOutput
	output = append(output, types.CheckOutput{
		Ref: refID,
	})

	result, err := json.Marshal(output)
	errors.E(err)

	fmt.Println(string(result))

}
