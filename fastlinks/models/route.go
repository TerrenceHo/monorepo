package models

import (
	"bytes"
	"encoding/gob"
	"strings"

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

type routeModelValFunc func(*Route) error

// Validate validates that a Route object is proper. The current validations are
// as follows:
//  - Key is not empty
//  - RedirectURL is not empty
//  - If ExtendedURL is not empty, then RedirectURL must contain the format
//    specifier "{}"
func (r *Route) Validate() error {
	return validateFuncs(r, hasKey, hasRedirect, validateExtendedURL)
}

func validateFuncs(r *Route, fns ...routeModelValFunc) error {
	for _, fn := range fns {
		if err := fn(r); err != nil {
			return stackerrors.Wrap(err, "validation failed")
		}
	}
	return nil
}

func hasKey(r *Route) error {
	if r.Key == "" {
		return stackerrors.New("Route must have a Key, was empty string")
	}
	return nil
}

func hasRedirect(r *Route) error {
	if r.RedirectURL == "" {
		return stackerrors.New("Route must have a RedirectURL, was empty string")
	}
	return nil
}

func validateExtendedURL(r *Route) error {
	if r.ExtendedURL != "" {
		if !strings.Contains(r.ExtendedURL, "{}") {
			return stackerrors.New("If route has ExtendedURL, it must contain '{}'")
		}
	}
	return nil
}
