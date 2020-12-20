package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var (
	sha1ver      string
	buildTime    string
)

func handleDebug(w http.ResponseWriter, r *http.Request) {
	s := fmt.Sprintf("url: %s %s", r.Method, r.RequestURI)
	a := []string{s}

	a = append(a, "Headers:")
	for k, v := range r.Header {
		if len(v) == 0 {
			a = append(a, k)
		} else if len(v) == 1 {
			s = fmt.Sprintf("  %s: %v", k, v[0])
			a = append(a, s)
		} else {
			a = append(a, "  "+k+":")
			for _, v2 := range v {
				a = append(a, "    "+v2)
			}
		}
	}

	a = append(a, "")
	a = append(a, fmt.Sprintf("ver: https://github.com/naiduarvind/serialized-jira/commit/%s", sha1ver))
	a = append(a, fmt.Sprintf("built on: %s", buildTime))

	s = strings.Join(a, "\n")
	servePlainText(w, s)
}

func servePlainText(w http.ResponseWriter, s string) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.Itoa(len(s)))
	w.WriteHeader(http.StatusOK)
	write, err := w.Write([]byte(s))
	if err != nil {
		log.Printf("Failed to write %v with error %s!", write, err)
	}
}
