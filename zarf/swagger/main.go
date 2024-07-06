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

	swaggerHandler := swaggerui.Handler(spec)

	http.Handle("/", swaggerHandler)
	log.Println("serving on  http://localhost:8086")
	log.Fatal(http.ListenAndServe("localhost:8086", nil))
}

//
//func byteHandler(b []byte) http.HandlerFunc {
//	return func(w http.ResponseWriter, _ *http.Request) {
//		w.Write(b)
//	}
//}
//
//// Handler returns a handler that will serve a self-hosted Swagger UI with your spec embedded
//func Handler(spec []byte) http.Handler {
//	// render the index template with the proper spec name inserted
//	static, _ := fs.Sub(swagfs, "embed")
//	mux := http.NewServeMux()
//	mux.HandleFunc("/swagger_spec", byteHandler(spec))
//	mux.Handle("/", http.FileServer(http.FS(static)))
//	return mux
//}
