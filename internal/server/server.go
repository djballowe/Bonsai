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
var last *printer.PrinterState

func Broadcast(state *printer.PrinterState) {
	last = merge(last, state)
	select {
	case updates <- last:
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
			var buf bytes.Buffer
			err := statusTmpl.Execute(&buf, state)
			if err != nil {
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

func merge(last *printer.PrinterState, current *printer.PrinterState) *printer.PrinterState {
	if last == nil {
		return current
	}
	merged := *last

	if current.NozzleTemp != nil {
		merged.NozzleTemp = current.NozzleTemp
	}
	if current.BedTemp != nil {
		merged.BedTemp = current.BedTemp
	}
	if current.ChamberTemp != nil {
		merged.ChamberTemp = current.ChamberTemp
	}
	if current.GcodeState != nil {
		merged.GcodeState = current.GcodeState
	}
	if current.PrintPercent != nil {
		merged.PrintPercent = current.PrintPercent
	}
	if current.RemainingTime != nil {
		merged.RemainingTime = current.RemainingTime
	}
	if current.TotalLayerNum != nil {
		merged.TotalLayerNum = current.TotalLayerNum
	}
	if current.PrintError != nil {
		merged.PrintError = current.PrintError
	}
	if current.WifiSignal != nil {
		merged.WifiSignal = current.WifiSignal
	}
	return &merged
}
