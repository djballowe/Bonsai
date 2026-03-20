package server

import (
	"bonsai/internal/printer"
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

var tmpl = template.Must(template.ParseFiles("templates/index.html"))
var statusTmpl = template.Must(template.ParseFiles("templates/status.html"))

var updates = make(chan *printer.PrinterState, 1)

func Broadcast(state *printer.PrinterState) {
	select {
	case updates <- state:
	default:
	}
}

func Start() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
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
	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf("template error: %v", err)
	}
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
			log.Printf("nozzle temp: %.1f", state.NozzleTemp)
			var buf bytes.Buffer
			if err := statusTmpl.Execute(&buf, state); err != nil {
				log.Printf("template error: %v", err)
				continue
			}
			html := strings.ReplaceAll(buf.String(), "\n", "")
			fmt.Fprintf(w, "event: printerUpdate\ndata: %s\n\n", html)
			flusher.Flush()
		case <-r.Context().Done():
			log.Printf("sse client disconnected: %s", r.RemoteAddr)
			return
		}
	}
}
