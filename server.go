// Package main runs the bundleservice server, a microservice which provides
// API endpoints for various bundle-related functionalities, such as retrieving
// the list of changes which the bundle describes.
package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/juju/bundlechanges"
	"github.com/juju/httprequest"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/errgo.v1"
	"gopkg.in/juju/charm.v6-unstable"

	"github.com/juju/bundleservice/params"
)

// main builds and runs the server when the run command is called.
func main() {
	router := httprouter.New()
	// Add handlers
	f := func(p httprequest.Params) (*handler, error) {
		return &handler{}, nil
	}
	for _, h := range errorMapper.Handlers(f) {
		router.Handle(h.Method, h.Path, h.Handle)
	}
	log.Fatal(http.ListenAndServe(":8000", router))
}

// The handler type contains all of the handlers for the server.
type handler struct{}

// errorResponse represents an error encountered by the server.
type errorResponse struct {
	Message string
}

// errorMapper maps an error from a handler into an HTTP server error.
// TODO map other errors
var errorMapper httprequest.ErrorMapper = func(err error) (int, interface{}) {
	status := http.StatusInternalServerError
	cause := errgo.Cause(err)
	switch cause {
	case params.ErrUnparsable:
		status = 422 // "Unprocessable Entity" - http doesn't have a const for this error.
	}
	return status, &errorResponse{
		Message: err.Error(),
	}
}

// GetChangesFromYAML receives a bundle in the request body and returns a list
// of changes.
func (h *handler) GetChangesFromYAML(p *params.ChangesFromYAMLParams) (params.ChangesResponse, error) {
	changes, err := getChanges(p.Body.Bundle)
	if err != nil {
		return params.ChangesResponse{}, err
	}
	return params.ChangesResponse{
		Changes: changes,
	}, nil
}

// getChanges recieves a bundle in YAML format as a string and returns the list
// of changes from the changeset.
func getChanges(bundleYAML string) ([]params.Change, error) {
	bundle, err := charm.ReadBundleData(strings.NewReader(bundleYAML))
	if err != nil {
		return nil, errgo.WithCausef(err, params.ErrUnparsable, "error reading bundle data")
	}
	err = bundle.Verify(nil, nil)
	if err != nil {
		return nil, errgo.WithCausef(err, params.ErrUnparsable, "error verifying bundle data")
	}
	changes := bundlechanges.FromData(bundle)
	changeSet := make([]params.Change, len(changes))
	for i, change := range changes {
		changeSet[i] = params.Change{
			Id:       change.Id(),
			Args:     change.GUIArgs(),
			Requires: change.Requires(),
			Method:   change.Method(),
		}
	}
	return changeSet, nil
}
