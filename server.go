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
		// error
		fmt.Print("Bundle is empty")
	}
	bundle, err := charm.ReadBundleData(strings.NewReader(bundleYAML))
	if err != nil {
		// error
		fmt.Printf("Error reading bundle data: %v", err)
	}
	err = bundle.Verify(nil, nil)
	if err != nil {
		// error
		fmt.Printf("Error verifying bundle data: %v", err)
	}
	changes := bundlechanges.FromData(bundle)
	changesJSON, err := json.Marshal(changes)
	if err != nil {
		// error
		fmt.Printf("Error marshalling JSON: %v", err)
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
