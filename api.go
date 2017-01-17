package main

import (
	"fmt"
	"io"
	"net/http"
	"text/template"
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

// rootHandler just renders the root.tmpl that explains the api calls usage
func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t := template.Must(template.ParseFiles("templates/root.tmpl"))
	t.Execute(w, nil)
}

// reportsHandler takes care of the path /api/reports that lists all the checks
// on all the hosts, each elements defined based on the reports_element.tmpl
// template
func reportsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	t := template.Must(template.ParseFiles("templates/reports_element.tmpl"))
	io.WriteString(w, "[")
	i := 1
	for host, svcs := range cache {
		j := 1
		for svc, chk := range svcs {
			c := map[string]map[string]string{
				"check": map[string]string{"host": host, "name": svc, "status": statusString(chk.state), "message": chk.output, "timestamp": fmt.Sprint(chk.timestamp)},
				// additions will be used to inject custom-defined fields
				"additions": map[string]string{"Addition field example": "and my value"},
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

// initAPIServer starts the API HTTP server. This is where the routes are
// defined
func initAPIServer(listenerIP string, port uint) {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/api/reports", reportsHandler)
	http.ListenAndServe(fmt.Sprint(listenerIP, ":", port), nil)
}
