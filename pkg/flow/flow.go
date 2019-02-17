package flow

import (
	"io"
	"lumo"
)

// Service is the usecase when we want to work within the domains flows.
type Service struct {
	zip     lumo.Zipper
	decoder lumo.Decoder
}

// NewService creates a new service.
func NewService(z lumo.Zipper, d lumo.Decoder) *Service {
	return &Service{z, d}
}

// DecodeZipFile takes a reader to a zipfile and then attempts to decrypt it with the given password.
func (s *Service) DecodeZipFile(r io.Reader, password string) (lumo.Flows, error) {
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

// EncodeFlowToZip takes flows and encrypts it into a zip file that contains an encrypted file from the decoder.
func (s *Service) EncodeFlowToZip(f lumo.Flows, name, password string) (io.Reader, error) {
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
