package types

import "github.com/digitalocean/godo"

type Payload struct {
	Source Source `json:"source" validate:"required"`
	Params Params `josn:"params,omitempty"`
}

type Source struct {
	Token     string                                  `json:"token" validate:"required"`
	Region    string                                  `json:"region" validate:"required"`
	Tags      []string                                `json:"tags" validate:"required"`
	NodePools []*godo.KubernetesNodePoolCreateRequest `json:"node_pools" validate:"required"`
}
