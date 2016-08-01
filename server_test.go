package main

import (
	"testing"

	"github.com/juju/bundleservice/params"

	gc "gopkg.in/check.v1"
)

func Test(t *testing.T) { gc.TestingT(t) }

type bundleServiceSuite struct{}

var _ = gc.Suite(&bundleServiceSuite{})

func (s *bundleServiceSuite) TestGetChangesForBundle(c *gc.C) {
	var bundle = `
applications:
  mongodb:
    charm: "cs:precise/mongodb-21"
    num_units: 1
    annotations:
      "gui-x": "940.5"
      "gui-y": "388.7698359714502"
    constraints: "mem=2G cpu-cores=1"
series: precise
`
	request := params.ChangesFromYAMLParams{
		Body: params.ChangesRequest{
			Bundle: bundle,
		},
	}
	h := handler{}
	response, err := h.GetChangesFromYAML(&request)
	c.Assert(err, gc.IsNil)
	c.Assert(len(response.Changes), gc.Equals, 4)
}

func (s *bundleServiceSuite) TestGetChangesForBundleError(c *gc.C) {
	request := params.ChangesFromYAMLParams{
		Body: params.ChangesRequest{
			Bundle: "bad-wolf",
		},
	}
	h := handler{}
	response, err := h.GetChangesFromYAML(&request)
	c.Assert(response, gc.DeepEquals, params.ChangesResponse{})
	c.Assert(err.Error(), gc.Equals, "error reading bundle data: cannot unmarshal bundle data: yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `bad-wolf` into charm.legacyBundleData")
}
