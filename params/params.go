package params

import (
	"github.com/juju/bundlechanges"
	"github.com/juju/httprequest"
)

// ChangesResponse contains the results of parsing a bundle into a list of
// changes.
type ChangesResponse struct {
	Changes []bundlechanges.Change
}

// ChangesRequest contains the bundle which is to be parsed into a list of
// changes.
type ChangesRequest struct {
	Bundle string
}

// ChangesFromYAMLParams contains the parameters required for passing a bundle
// to the API and recieving a list of changes in return
type ChangesFromYAMLParams struct {
	httprequest.Route `httprequest:"POST /bundlechanges/fromYAML"`
	NicelyFormatted   bool           `httprequest:"nice,form"`
	Body              ChangesRequest `httprequest:",body"`
}
