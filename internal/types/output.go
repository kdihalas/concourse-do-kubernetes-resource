package types

type CheckOutput struct {
	Ref string `json:"ref"`
}

type InOutput struct {
	Version  CheckOutput `json:"version"`
	Metadata []NameValue `json:"metadata,omitempty"`
}

type NameValue struct {
	Name  string `json:"name"`
	Value string `json:"value,omitempty"`
}
