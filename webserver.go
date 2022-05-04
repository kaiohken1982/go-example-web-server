package main

import (
	"fmt"
  "log"
	"net/http"
  "sync"
)

var mu sync.Mutex
var count int

func main() {
  registerRoutes()
  log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func registerRoutes() {
  http.HandleFunc("/", routeMain) // each request calls handler
  http.HandleFunc("/count", routeCount)
  http.HandleFunc("/details", routeRequestDetails)
}

func routeMain(w http.ResponseWriter, r *http.Request) {
  /**
   * Gestisce la race conditions, 
   * ovvero N chiamate contemporanee.
   * Solo una goroutine che gestisce la chiamata accede a questa variabile 
   */
  mu.Lock()
  count++
  mu.Unlock()
  fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}

func routeCount(w http.ResponseWriter, r *http.Request) {
  mu.Lock()
  fmt.Fprintf(w, "Count %d\n", count)
  mu.Unlock()
}

func routeRequestDetails(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL, r.Proto)
  for k, v := range r.Header {
    fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
  }

  fmt.Fprintf(w, "Host = %q\n", r.Host)
  fmt.Fprintf(w, "RemoteAddr = %q\n", r.RemoteAddr)
  if err := r.ParseForm(); err != nil {
    log.Print(err)
  }
  for k, v := range r.Form {
    fmt.Fprintf(w, "Form[%q] = %q\n", k, v)
  }
}
