package server

import (
	"bonsai/internal/printer"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var updates = make(chan *printer.PrinterState, 1)

func Broadcast(state *printer.PrinterState) {
	select {
	case updates <- state:
	default:
	}
}

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
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	log.Printf("sse client connected: %s", r.RemoteAddr)

	for {
		select {
		case state := <-updates:
			data, err := json.Marshal(state)
			if err != nil {
				log.Printf("could not marshal json: %s", err)
				continue
			}

			fmt.Fprintf(w, "event: printerUpdate\ndata: %s\n\n", data)
			flusher.Flush()
		case <-r.Context().Done():
			log.Printf("sse client disconnected: %s", r.RemoteAddr)
			return
		}
	}
}
