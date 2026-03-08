package server

import (
	"fmt"
	"log"
	"net/http"
)

func Start() {
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/events", handleEvents)

	go func() {
		log.Println("listening on http://localhost:3100")
		if err := http.ListenAndServe(":3100", nil); err != nil {
			log.Fatalf("http: %v", err)
		}
	}()
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>bonsai dashboard</h1><p>coming soon</p>")
}

func handleEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	log.Printf("sse client connected: %s", r.RemoteAddr)

	<-r.Context().Done()
	log.Printf("sse client disconnected: %s", r.RemoteAddr)
}
