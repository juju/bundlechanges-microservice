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
	expected := []params.Change{
		params.Change{
			Args:     []interface{}{"cs:precise/mongodb-21"},
			Id:       "addCharm-0",
			Requires: []string{},
			Method:   "addCharm",
		}, params.Change{
			Args: []interface{}{
				"$addCharm-0",
				"mongodb",
				map[string]interface{}{},
				"mem=2G cpu-cores=1",
				map[string]string{},
				map[string]string{},
			},
			Id:       "deploy-1",
			Requires: []string{"addCharm-0"},
			Method:   "deploy",
		}, params.Change{
			Args: []interface{}{
				"$deploy-1", "application",
				map[string]string{
					"gui-x": "940.5",
					"gui-y": "388.7698359714502",
				},
			},
			Id:       "setAnnotations-2",
			Requires: []string{"deploy-1"},
			Method:   "setAnnotations",
		}, params.Change{
			Args: []interface{}{
				"$deploy-1",
				interface{}(nil),
			},
			Id:       "addUnit-3",
			Requires: []string{"deploy-1"},
			Method:   "addUnit"},
	}
	c.Assert(response.Changes, gc.DeepEquals, expected)
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
	c.Assert(err, gc.ErrorMatches, `error reading bundle data: cannot unmarshal bundle data: yaml: unmarshal errors:\n.*`)
}
