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

// server represents a bundleservice HTTP server
type server struct {
	router *httprouter.Router
}

// ServeHTTP implements http.Handler.Handle.
func (srv *server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	header := w.Header()
	ao := "*"
	if o := req.Header.Get("Origin"); o != "" {
		ao = o
	}
	header.Set("Access-Control-Allow-Origin", ao)
	header.Set("Access-Control-Allow-Headers", "Bakery-Protocol-Version, Macaroons, X-Requested-With, Content-Type")
	header.Set("Access-Control-Allow-Credentials", "true")
	header.Set("Access-Control-Cache-Max-Age", "600")
	// TODO: in handlers, look up methods for this request path and return only those methods here.
	header.Set("Access-Control-Allow-Methods", "POST,OPTIONS")
	header.Set("Access-Control-Expose-Headers", "WWW-Authenticate")
	srv.router.ServeHTTP(w, req)
}

// options is a no-op handler that responds to OPTIONS requests for each path.
func (srv *server) options(http.ResponseWriter, *http.Request, httprouter.Params) {
	// no-op
}

// main builds and runs the server when the run command is called.
func main() {
	srv := &server{
		router: httprouter.New(),
	}
	// Add handlers
	f := func(p httprequest.Params) (*handler, error) {
		return &handler{}, nil
	}
	for _, h := range errorMapper.Handlers(f) {
		srv.router.Handle(h.Method, h.Path, h.Handle)
		srv.router.OPTIONS(h.Path, srv.options)
	}
	httpServer := &http.Server{
		Addr:    "0.0.0.0:8000",
		Handler: srv,
	}
	log.Fatal(httpServer.ListenAndServe())
}

// The handler type contains all of the handlers for the server.
type handler struct{}

// errorMapper maps an error from a handler into an HTTP server error.
// TODO map other errors
var errorMapper httprequest.ErrorMapper = func(err error) (int, interface{}) {
	status := http.StatusInternalServerError
	cause := errgo.Cause(err)
	code := cause.(params.ErrorCode)
	switch cause {
	case params.ErrUnparsable, params.ErrVerificationFailure:
		status = params.StatusUnprocessableEntity
	}
	return status, &params.ErrorResponse{
		Message: err.Error(),
		Code:    code,
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
		return nil, errgo.WithCausef(err, params.ErrVerificationFailure, "error verifying bundle data")
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
