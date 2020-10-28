package main

import (
	"regexp"
)

var html_pattern *regexp.Regexp
var url_pattern *regexp.Regexp

const resources_format = `package main

import (
	"net/http"
	"strings"
)

func Quark_server() {
	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		url := req.URL.String()
		if strings.HasSuffix(url, ".css") {
			rw.Header().Set("Content-Type", "text/css")
		}
		if strings.HasSuffix(url, ".html") {
			rw.Header().Set("Content-Type", "text/html")
		}
		if strings.HasSuffix(url, ".js") {
			rw.Header().Set("Content-Type", "text/js")
		}
		rw.Write([]byte(resources[url[1:]]))
	})

	http.ListenAndServe(":8080", nil)
}

var resources map[string][]byte = map[string][]byte {
	%s
}
`

const file_format = "\"%s\": %s,"

const html_regex = "(href|src) *= *\"(.*)\""
const url_regex = ".+:\\/\\/.+\\..+"
