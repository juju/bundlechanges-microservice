package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/juju/bundlechanges"
	"gopkg.in/juju/charm.v6-unstable"

	"github.com/juju/httprequest"
	"github.com/julienschmidt/httprouter"
)

// The handler type contains all of the handlers for the server.
type handler struct{}

type errorResponse struct {
	Message string
}

var errorMapper httprequest.ErrorMapper = func(err error) (int, interface{}) {
	return http.StatusInternalServerError, &errorResponse{
		Message: err.Error(),
	}
}

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

type changesResponse struct {
	Changes []bundlechanges.Change
}

type changesRequest struct {
	Bundle string
}

// Retrieving changes from YAML
type changesFromYAMLParams struct {
	httprequest.Route `httprequest:"POST /bundlechanges/fromYAML"`
	NicelyFormatted   bool           `httprequest:"nice,form"`
	Body              changesRequest `httprequest:",body"`
}

func (h *handler) GetChangesFromYAML(p *changesFromYAMLParams) (changesResponse, error) {
	changes, err := getChanges(p.Body.Bundle)
	if err != nil {
		return changesResponse{}, err
	}
	return changesResponse{
		Changes: changes,
	}, nil
}

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
