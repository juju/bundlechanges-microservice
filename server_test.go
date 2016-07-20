package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"

	gc "gopkg.in/check.v1"
	//"gopkg.in/juju/charm.v6-unstable"
)

func Test(t *testing.T) { gc.TestingT(t) }

type newSuite struct{}

var _ = gc.Suite(&newSuite{})

func (s *newSuite) TestGetChangesForBundle(c *gc.C) {
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

	rec := httptest.NewRecorder()
	handler := func(_ http.ResponseWriter, r *http.Request) {
		getChangesFromYAML(rec, r, httprouter.Params{})
	}

	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	_, err := http.NewRequest("POST", ts.URL, strings.NewReader(fmt.Sprintf("bundleYAML=%s", bundle)))
	if err != nil {
		log.Fatal(err)
	}

	c.Logf("%v", rec)
	c.Fail()
}
