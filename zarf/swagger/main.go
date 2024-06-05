package main

import (
	_ "embed"
	"log"
	"net/http"

	"github.com/flowchartsman/swaggerui"
)

//go:embed swagger.json
var spec []byte

func main() {
	log.SetFlags(0)
	http.Handle("/", swaggerui.Handler(spec))
	log.Println("serving on  http://localhost:8000")
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
