package params

import (
	"github.com/juju/httprequest"

	"github.com/juju/bundlechanges"
)

// ChangesResponse contains the results of parsing a bundle into a list of
// changes.
type ChangesResponse struct {
	// TODO This should be an API specific type
	Changes []bundlechanges.Change `json:"changes"`
}

// ChangesRequest contains the bundle as a YAML-encoded string which is to be
// parsed into a list of changes.
type ChangesRequest struct {
	Bundle string `json:"bundle"`
}

// ChangesFromYAMLParams contains the parameters required for passing a bundle
// to the API and recieving a list of changes in return
type ChangesFromYAMLParams struct {
	httprequest.Route `httprequest:"POST /bundlechanges/fromYAML"`
	NicelyFormatted   bool           `httprequest:"nice,form"`
	Body              ChangesRequest `httprequest:",body"`
}
