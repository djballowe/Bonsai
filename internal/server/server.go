package server

import (
	"bonsai/internal/printer"
	"bonsai/internal/view"
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
			templateStateUI := view.NewStatusView(state)
			err := statusTmpl.Execute(&buf, templateStateUI)
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

// This is cringe but whatever for now
func merge(last *printer.PrinterState, current *printer.PrinterState) *printer.PrinterState {
	if last == nil {
		return current
	}
	merged := *last

	// Temperatures
	if current.NozzleTemp != nil {
		merged.NozzleTemp = current.NozzleTemp
	}
	if current.NozzleTargetTemp != nil {
		merged.NozzleTargetTemp = current.NozzleTargetTemp
	}
	if current.BedTemp != nil {
		merged.BedTemp = current.BedTemp
	}
	if current.BedTargetTemp != nil {
		merged.BedTargetTemp = current.BedTargetTemp
	}
	if current.ChamberTemp != nil {
		merged.ChamberTemp = current.ChamberTemp
	}
	// Print progress
	if current.GcodeState != nil {
		merged.GcodeState = current.GcodeState
	}
	if current.PrintPercent != nil {
		merged.PrintPercent = current.PrintPercent
	}
	if current.RemainingTime != nil {
		merged.RemainingTime = current.RemainingTime
	}
	if current.PrintStage != nil {
		merged.PrintStage = current.PrintStage
	}
	if current.PrintSubStage != nil {
		merged.PrintSubStage = current.PrintSubStage
	}
	if current.PrintLineNum != nil {
		merged.PrintLineNum = current.PrintLineNum
	}
	if current.LayerNum != nil {
		merged.LayerNum = current.LayerNum
	}
	if current.TotalLayerNum != nil {
		merged.TotalLayerNum = current.TotalLayerNum
	}
	if current.PrintType != nil {
		merged.PrintType = current.PrintType
	}
	if current.PrintError != nil {
		merged.PrintError = current.PrintError
	}
	// File / task
	if current.GcodeFile != nil {
		merged.GcodeFile = current.GcodeFile
	}
	if current.GcodeFilePrepercent != nil {
		merged.GcodeFilePrepercent = current.GcodeFilePrepercent
	}
	if current.SubtaskName != nil {
		merged.SubtaskName = current.SubtaskName
	}
	if current.SubtaskID != nil {
		merged.SubtaskID = current.SubtaskID
	}
	if current.TaskID != nil {
		merged.TaskID = current.TaskID
	}
	if current.ProjectID != nil {
		merged.ProjectID = current.ProjectID
	}
	if current.ProfileID != nil {
		merged.ProfileID = current.ProfileID
	}
	// Fans
	if current.CoolingFanSpeed != nil {
		merged.CoolingFanSpeed = current.CoolingFanSpeed
	}
	if current.HeatbreakFanSpeed != nil {
		merged.HeatbreakFanSpeed = current.HeatbreakFanSpeed
	}
	if current.BigFan1Speed != nil {
		merged.BigFan1Speed = current.BigFan1Speed
	}
	if current.BigFan2Speed != nil {
		merged.BigFan2Speed = current.BigFan2Speed
	}
	if current.FanGear != nil {
		merged.FanGear = current.FanGear
	}
	// Speed
	if current.SpeedMagnitude != nil {
		merged.SpeedMagnitude = current.SpeedMagnitude
	}
	if current.SpeedLevel != nil {
		merged.SpeedLevel = current.SpeedLevel
	}
	// Hardware
	if current.NozzleDiameter != nil {
		merged.NozzleDiameter = current.NozzleDiameter
	}
	if current.NozzleType != nil {
		merged.NozzleType = current.NozzleType
	}
	if current.SDCard != nil {
		merged.SDCard = current.SDCard
	}
	if current.WifiSignal != nil {
		merged.WifiSignal = current.WifiSignal
	}
	// AMS
	if current.AMSStatus != nil {
		merged.AMSStatus = current.AMSStatus
	}
	if current.AMSRFIDStatus != nil {
		merged.AMSRFIDStatus = current.AMSRFIDStatus
	}
	if len(current.AMS.AMSList) > 0 {
		merged.AMS = current.AMS
	}
	// Lights
	if len(current.LightsReport) > 0 {
		merged.LightsReport = current.LightsReport
	}
	// Queue
	if current.QueueNumber != nil {
		merged.QueueNumber = current.QueueNumber
	}
	if current.QueueTotal != nil {
		merged.QueueTotal = current.QueueTotal
	}
	return &merged
}
