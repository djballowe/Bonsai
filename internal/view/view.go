package view

import (
	"bonsai/internal/printer"
	"fmt"
)

type StatusView struct {
	NozzleTemp    string
	BedTemp       string
	ChamberTemp   string
	GcodeState    string
	PrintPercent  string
	RemainingTime string
	TotalLayerNum string
	LayerNum      string
	PrintError    string
	WifiSignal    string
}

func NewStatusView(s *printer.PrinterState) StatusView {
	v := StatusView{
		NozzleTemp:    "--",
		BedTemp:       "--",
		ChamberTemp:   "--",
		GcodeState:    "--",
		PrintPercent:  "--",
		RemainingTime: "--",
		TotalLayerNum: "--",
		LayerNum:      "--",
		PrintError:    "none",
		WifiSignal:    "--",
	}

	if s == nil {
		return v
	}

	if s.NozzleTemp != nil {
		v.NozzleTemp = fmt.Sprintf("%.1f°C", *s.NozzleTemp)
	}
	if s.BedTemp != nil {
		v.BedTemp = fmt.Sprintf("%.1f°C", *s.BedTemp)
	}
	if s.ChamberTemp != nil {
		v.ChamberTemp = fmt.Sprintf("%.1f°C", *s.ChamberTemp)
	}
	if s.GcodeState != nil {
		v.GcodeState = *s.GcodeState
	}
	if s.PrintPercent != nil {
		v.PrintPercent = fmt.Sprintf("%d%%", *s.PrintPercent)
	}
	if s.RemainingTime != nil {
		v.RemainingTime = fmt.Sprintf("%d min", *s.RemainingTime)
	}
	if s.TotalLayerNum != nil {
		v.TotalLayerNum = fmt.Sprintf("%d", *s.TotalLayerNum)
	}
	if s.LayerNum != nil {
		v.LayerNum = fmt.Sprintf("%d", *s.LayerNum)
	}
	if s.PrintError != nil && *s.PrintError != 0 {
		v.PrintError = fmt.Sprintf("%d", *s.PrintError)
	}
	if s.WifiSignal != nil {
		v.WifiSignal = *s.WifiSignal
	}

	return v
}
