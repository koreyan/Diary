package main

import (
	"diary/network"
	"net/http"

	"github.com/urfave/negroni"
)

func main() {
	n := negroni.Classic()
	n.UseHandler(network.MakeHandler())
	http.ListenAndServe(":8080", n)
}
