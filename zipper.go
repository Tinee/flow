package lumo

import (
	"io"
)

// Zipper is an interface that describe how to take a flow and make it into a zip file.
type Zipper interface {
	WithPassword(r io.Reader, name, password string) (io.Reader, error)
	Unlock(r io.Reader, password string) (io.Reader, error)
}
