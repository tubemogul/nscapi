package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	cFields  customFields
	tmplRoot string
)

// statusString returns the corresponding string to the nagios status
func statusString(state int16) string {
	switch state {
	case 0:
		return "OK"
	case 1:
		return "Warning"
	case 2:
		return "Critical"
	}
	return "Unknown"
}

// ToJSONString is a function that takes an interface as output and formats it for
// json output in return. Its purpose is to be used in the templates which is
// the reason wly it is exported
func ToJSONString(v interface{}) string {
	bytesOutput, _ := json.Marshal(v)
	return string(bytesOutput)
}

// We need to remove newline characters
// to have valid JSON output
// because it breaks it in other case
func sanitizeJSONString(str string) string {
	return strings.Replace(str, "\n", " ", -1)
}

// rootHandler just renders the root.tmpl that explains the api calls usage
func rootHandler(w http.ResponseWriter, r *http.Request) {
	tmplPath := filepath.Join(tmplRoot, "root.tmpl")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t := template.Must(template.ParseFiles(tmplPath))
	t.Execute(w, nil)
}

// reportsHandler takes care of the path /api/reports that lists all the checks
// on all the hosts, each elements defined based on the reports_element.tmpl
// template
func reportsHandler(w http.ResponseWriter, r *http.Request) {
	tmplName := "reports_element.tmpl"
	tmplPath := filepath.Join(tmplRoot, tmplName)
	w.Header().Set("Content-Type", "application/json")
	fMaps := template.FuncMap{"tojson": ToJSONString}
	t := template.Must(template.New(tmplName).Funcs(fMaps).ParseFiles(tmplPath))
	io.WriteString(w, "[")
	i := 1
	for host, svcs := range cache {
		j := 1
		for svc, chk := range svcs {
			c := map[string]map[string]interface{}{
				"check": map[string]interface{}{"host": host, "name": svc, "status": statusString(chk.state), "message": sanitizeJSONString(chk.output), "timestamp": fmt.Sprint(chk.timestamp), "statusFirstSeen": fmt.Sprint(chk.statusFirstSeen)},
				// custom will be used to inject custom-defined fields
				"custom": cFields.get(host, svc),
			}
			t.Execute(w, c)
			// This part just takes care of adding a coma or not between the elements
			// to have a correcly-formated json
			if !(i == len(cache) && j == len(svcs)) {
				io.WriteString(w, ",")
			}
			j++
		}
		i++
	}
	io.WriteString(w, "]\n")
}

// setIfPathExists validates that the given path exists before assigning it to
// the given variable
func setIfPathExists(dir string, varToSet *string) error {
	if _, err := os.Stat(dir); err != nil {
		return err
	}
	*varToSet = dir
	return nil
}

// initAPIServer starts the API HTTP server. This is where the routes are
// defined. The customFieldRoot is the root of the hierarchy of yaml files used
// for the custom fields. The templatesRoot is the root directory where to find
// the templates used by the API
func initAPIServer(listenerIP string, port uint, customFieldRoot string, templatesRoot string) {
	// Init custom fields
	var customFRoot string
	setIfPathExists(customFieldRoot, &customFRoot)
	cFields.load(customFRoot)

	setIfPathExists(templatesRoot, &tmplRoot)
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/api/reports", reportsHandler)
	http.ListenAndServe(fmt.Sprint(listenerIP, ":", port), nil)
}
