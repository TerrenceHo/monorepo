package models

import (
	"bytes"
	"encoding/gob"

	"github.com/TerrenceHo/monorepo/utils-go/stackerrors"
)

type Route struct {
	// Key short URL name
	Key string

	// RedirectURL is the URL to redirect to
	RedirectURL string

	// ExtendedURL holds optional appended string to the RedirectURL
	ExtendedURL string
}

func (r *Route) Encode() (*bytes.Buffer, error) {
	var b bytes.Buffer
	e := gob.NewEncoder(&b)
	if err := e.Encode(*r); err != nil {
		return nil, stackerrors.Wrap(err, "failed to serialize route to bytes")
	}
	return &b, nil
}

func (r *Route) Decode(b *bytes.Buffer) error {
	d := gob.NewDecoder(b)
	if err := d.Decode(r); err != nil {
		return stackerrors.Wrap(err, "failed to deserialize route to struct")
	}
	return nil
}
