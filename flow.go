package flow

import (
	"io"
)

// Flows is a slice of flows.
type Flows []Flow

// Flow is the domain structure of a flow.
type Flow struct {
	Type   string   `json:"type"`
	Fields []string `json:"fields"`
}

// Decoder is an interface that describe how to take a flow and make it into a zip file.
type Decoder interface {
	Decode(r io.Reader) (Flows, error)
	Encode(flows Flows, name string) (io.Reader, error)
}
