package view

import (
	"bonsai/internal/printer"
	"fmt"
	"strconv"
)

type StatusView struct {
	// Temperatures
	NozzleTemp       string
	NozzleTargetTemp string
	BedTemp          string
	BedTargetTemp    string
	ChamberTemp      string
	// Progress
	GcodeState        string
	PrintPercent      string
	PrintPercentInt   int
	RemainingTime     string
	RemainingTimeHHMM string
	LayerNum          string
	TotalLayerNum     string
	PrintType         string
	PrintStage        string
	PrintError        string
	// File / Task
	GcodeFile   string
	SubtaskName string
	TaskID      string
	// Fans
	CoolingFanSpeed   string
	CoolingFanPct     int
	HeatbreakFanSpeed string
	HeatbreakFanPct   int
	BigFan1Speed      string
	BigFan1Pct        int
	BigFan2Speed      string
	BigFan2Pct        int
	// Speed
	SpeedMagnitude string
	SpeedLevel     string
	// Hardware
	NozzleDiameter string
	NozzleType     string
	SDCard         string
	WifiSignal     string
	// Queue
	QueueNumber string
	QueueTotal  string
	// Filament (VTTray)
	FilamentType    string
	FilamentColor   string
	FilamentTempMin string
	FilamentTempMax string
	// Lights
	ChamberLight string
}

func fanPct(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil || val < 0 {
		return 0
	}
	if val > 15 {
		val = 15
	}
	return val * 100 / 15
}

func formatMinutes(min int) string {
	return fmt.Sprintf("%02d:%02d", min/60, min%60)
}

func NewStatusView(s *printer.PrinterState) StatusView {
	v := StatusView{
		NozzleTemp:        "--",
		NozzleTargetTemp:  "--",
		BedTemp:           "--",
		BedTargetTemp:     "--",
		ChamberTemp:       "--",
		GcodeState:        "--",
		PrintPercent:      "--",
		PrintPercentInt:   0,
		RemainingTime:     "--",
		RemainingTimeHHMM: "--",
		LayerNum:          "--",
		TotalLayerNum:     "--",
		PrintType:         "--",
		PrintStage:        "--",
		PrintError:        "none",
		GcodeFile:         "--",
		SubtaskName:       "--",
		TaskID:            "--",
		CoolingFanSpeed:   "--",
		CoolingFanPct:     0,
		HeatbreakFanSpeed: "--",
		HeatbreakFanPct:   0,
		BigFan1Speed:      "--",
		BigFan1Pct:        0,
		BigFan2Speed:      "--",
		BigFan2Pct:        0,
		SpeedMagnitude:    "--",
		SpeedLevel:        "--",
		NozzleDiameter:    "--",
		NozzleType:        "--",
		SDCard:            "--",
		WifiSignal:        "--",
		QueueNumber:       "--",
		QueueTotal:        "--",
		FilamentType:      "--",
		FilamentColor:     "--",
		FilamentTempMin:   "--",
		FilamentTempMax:   "--",
		ChamberLight:      "--",
	}

	if s == nil {
		return v
	}

	// Temperatures
	if s.NozzleTemp != nil {
		v.NozzleTemp = fmt.Sprintf("%.1f°C", *s.NozzleTemp)
	}
	if s.NozzleTargetTemp != nil {
		v.NozzleTargetTemp = fmt.Sprintf("%.1f°C", *s.NozzleTargetTemp)
	}
	if s.BedTemp != nil {
		v.BedTemp = fmt.Sprintf("%.1f°C", *s.BedTemp)
	}
	if s.BedTargetTemp != nil {
		v.BedTargetTemp = fmt.Sprintf("%.1f°C", *s.BedTargetTemp)
	}
	if s.ChamberTemp != nil {
		v.ChamberTemp = fmt.Sprintf("%.1f°C", *s.ChamberTemp)
	}
	// Progress
	if s.GcodeState != nil {
		v.GcodeState = *s.GcodeState
	}
	if s.PrintPercent != nil {
		v.PrintPercent = fmt.Sprintf("%d%%", *s.PrintPercent)
		v.PrintPercentInt = *s.PrintPercent
	}
	if s.RemainingTime != nil {
		v.RemainingTime = fmt.Sprintf("%d min", *s.RemainingTime)
		v.RemainingTimeHHMM = formatMinutes(*s.RemainingTime)
	}
	if s.LayerNum != nil {
		v.LayerNum = fmt.Sprintf("%d", *s.LayerNum)
	}
	if s.TotalLayerNum != nil {
		v.TotalLayerNum = fmt.Sprintf("%d", *s.TotalLayerNum)
	}
	if s.PrintType != nil {
		v.PrintType = *s.PrintType
	}
	if s.PrintStage != nil {
		v.PrintStage = *s.PrintStage
	}
	if s.PrintError != nil && *s.PrintError != 0 {
		v.PrintError = fmt.Sprintf("%d", *s.PrintError)
	}
	// File / Task
	if s.GcodeFile != nil {
		v.GcodeFile = *s.GcodeFile
	}
	if s.SubtaskName != nil {
		v.SubtaskName = *s.SubtaskName
	}
	if s.TaskID != nil {
		v.TaskID = *s.TaskID
	}
	// Fans
	if s.CoolingFanSpeed != nil {
		v.CoolingFanSpeed = *s.CoolingFanSpeed
		v.CoolingFanPct = fanPct(*s.CoolingFanSpeed)
	}
	if s.HeatbreakFanSpeed != nil {
		v.HeatbreakFanSpeed = *s.HeatbreakFanSpeed
		v.HeatbreakFanPct = fanPct(*s.HeatbreakFanSpeed)
	}
	if s.BigFan1Speed != nil {
		v.BigFan1Speed = *s.BigFan1Speed
		v.BigFan1Pct = fanPct(*s.BigFan1Speed)
	}
	if s.BigFan2Speed != nil {
		v.BigFan2Speed = *s.BigFan2Speed
		v.BigFan2Pct = fanPct(*s.BigFan2Speed)
	}
	// Speed
	if s.SpeedMagnitude != nil {
		v.SpeedMagnitude = fmt.Sprintf("%d%%", *s.SpeedMagnitude)
	}
	if s.SpeedLevel != nil {
		v.SpeedLevel = fmt.Sprintf("%d", *s.SpeedLevel)
	}
	// Hardware
	if s.NozzleDiameter != nil {
		v.NozzleDiameter = *s.NozzleDiameter
	}
	if s.NozzleType != nil {
		v.NozzleType = *s.NozzleType
	}
	if s.SDCard != nil {
		if *s.SDCard {
			v.SDCard = "yes"
		} else {
			v.SDCard = "no"
		}
	}
	if s.WifiSignal != nil {
		v.WifiSignal = *s.WifiSignal
	}
	// Queue
	if s.QueueNumber != nil {
		v.QueueNumber = fmt.Sprintf("%d", *s.QueueNumber)
	}
	if s.QueueTotal != nil {
		v.QueueTotal = fmt.Sprintf("%d", *s.QueueTotal)
	}
	// Filament (VTTray)
	if s.VTTray.TrayType != "" {
		v.FilamentType = s.VTTray.TrayType
		v.FilamentTempMin = s.VTTray.NozzleTempMin + "°C"
		v.FilamentTempMax = s.VTTray.NozzleTempMax + "°C"
		if len(s.VTTray.TrayColor) >= 6 {
			v.FilamentColor = "#" + s.VTTray.TrayColor[:6]
		}
	}
	// Lights
	for _, light := range s.LightsReport {
		if light.Node == "chamber_light" {
			v.ChamberLight = light.Mode
		}
	}

	return v
}
