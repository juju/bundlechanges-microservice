package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"gopkg.in/juju/charm.v6-unstable"

	"github.com/juju/bundlechanges"
	"github.com/julienschmidt/httprouter"
)

func getChangesFromStore(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	bundleURL := params.ByName("bundleURL")[1:]
	fmt.Fprintf(w, "Not implemented; bundle requested: %s", bundleURL)
}

func getChangesFromYAML(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	bundleYAML := r.FormValue("bundleYAML")
	if bundleYAML == "" {
		http.Error(w, "Bundle is empty", 400)
		return
	}
	bundle, err := charm.ReadBundleData(strings.NewReader(bundleYAML))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading bundle data: %v", err), 422)
		return
	}
	err = bundle.Verify(nil, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error verifying bundle data: %v", err), 422)
		return
	}
	changes := bundlechanges.FromData(bundle)
	changesJSON, err := json.Marshal(changes)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error marshalling JSON: %v", err), 500)
		return
	}
	fmt.Fprint(w, string(changesJSON))
}

func main() {
	router := httprouter.New()
	// Add handlers
	// TODO:
	//router.GET("/bundlesvg/fromStore/*bundleURL", getSVGFromStore)
	//router.POST("/bundlesvg/fromYAML/", getSVGFromYAML)
	//router.GET("/bundlechanges/fromStore/*bundleURL", getChangesFromStore)
	router.POST("/bundlechanges/fromYAML/", getChangesFromYAML)
	log.Fatal(http.ListenAndServe(":8000", router))
}
