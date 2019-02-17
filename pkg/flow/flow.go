package flow

import (
	"flow"
	"io"
)

// Service is the usecase when we want to work within the domains flows.
type Service struct {
	zip     flow.Zipper
	decoder flow.Decoder
}

// NewService creates a new service.
func NewService(z flow.Zipper, d flow.Decoder) *Service {
	return &Service{z, d}
}

func (s *Service) DecodeZipFile(r io.Reader, password string) (flow.Flows, error) {
	r, err := s.zip.Unlock(r, password)
	if err != nil {
		return nil, err
	}
	f, err := s.decoder.Decode(r)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (s *Service) EncodeFlowToZip(f flow.Flows, name, password string) (io.Reader, error) {
	r, err := s.decoder.Encode(f, name)
	if err != nil {
		return nil, err
	}
	brd, err := s.zip.WithPassword(r, name, password)
	if err != nil {
		return nil, err
	}
	return brd, nil
}
