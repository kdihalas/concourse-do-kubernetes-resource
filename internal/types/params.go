package types

type Params struct {
	Name      string      `json:"name"`
	Skip      bool        `json:"skip,omitempty"`
	Delete    bool        `json:"delete,omitempty"`
	Version   string      `json:"version,omitempty"`
	NodePools []NodePools `json:"node_pools,omitempty"`
}

type NodePools struct {
	Name  string `json:"name" validate:"required"`
	Count int    `json:"count" validate:"required"`
}
