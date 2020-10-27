package main

import "regexp"

var html_pattern *regexp.Regexp
var url_pattern *regexp.Regexp

const resources_format = `package main

import "net/http"

func Quark_server() {
	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		url := req.URL.String()
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
