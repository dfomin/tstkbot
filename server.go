package main

import (
    "fmt"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "qwe")
    fmt.Println(r)
}

func main() {
    http.HandleFunc("/tstkbot", handler)
    err := http.ListenAndServeTLS(":8443", "/etc/letsencrypt/live/pigowl.com/fullchain.pem", "/etc/letsencrypt/live/pigowl.com/privkey.pem", nil)
    if err != nil {
        fmt.Println("ListenAndServe: ", err)
    }
}
