package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/juju/bundlechanges"
	"github.com/juju/bundleservice/params"
	"gopkg.in/juju/charm.v6-unstable"

	"github.com/juju/httprequest"
	"github.com/julienschmidt/httprouter"
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
var errorMapper httprequest.ErrorMapper = func(err error) (int, interface{}) {
	return http.StatusInternalServerError, &errorResponse{
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
func getChanges(bundleYAML string) ([]bundlechanges.Change, error) {
	bundle, err := charm.ReadBundleData(strings.NewReader(bundleYAML))
	if err != nil {
		return nil, fmt.Errorf("error reading bundle data: %v", err)
	}
	err = bundle.Verify(nil, nil)
	if err != nil {
		return nil, fmt.Errorf("error verifying bundle data: %v", err)
	}
	changes := bundlechanges.FromData(bundle)
	return changes, nil
}
