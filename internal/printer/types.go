package printer

type PrinterState struct {
	// Temperatures
	NozzleTemp       float64 `json:"nozzle_temper"`
	NozzleTargetTemp float64 `json:"nozzle_target_temper"`
	BedTemp          float64 `json:"bed_temper"`
	BedTargetTemp    float64 `json:"bed_target_temper"`
	ChamberTemp      float64 `json:"chamber_temper"`

	// Print progress
	GcodeState    string `json:"gcode_state"`
	PrintPercent  int    `json:"mc_percent"`
	RemainingTime int    `json:"mc_remaining_time"`
	PrintStage    string `json:"mc_print_stage"`
	PrintSubStage int    `json:"mc_print_sub_stage"`
	PrintLineNum  string `json:"mc_print_line_number"`
	LayerNum      int    `json:"layer_num"`
	TotalLayerNum int    `json:"total_layer_num"`
	PrintType     string `json:"print_type"`
	PrintError    int    `json:"print_error"`

	// File / task
	GcodeFile           string `json:"gcode_file"`
	GcodeFilePrepercent string `json:"gcode_file_prepare_percent"`
	SubtaskName         string `json:"subtask_name"`
	SubtaskID           string `json:"subtask_id"`
	TaskID              string `json:"task_id"`
	ProjectID           string `json:"project_id"`
	ProfileID           string `json:"profile_id"`

	// Fans (values are strings like "0"–"15")
	CoolingFanSpeed   string `json:"cooling_fan_speed"`
	HeatbreakFanSpeed string `json:"heatbreak_fan_speed"`
	BigFan1Speed      string `json:"big_fan1_speed"`
	BigFan2Speed      string `json:"big_fan2_speed"`
	FanGear           int    `json:"fan_gear"`

	// Speed
	SpeedMagnitude int `json:"spd_mag"`
	SpeedLevel     int `json:"spd_lvl"`

	// Hardware
	NozzleDiameter string `json:"nozzle_diameter"`
	NozzleType     string `json:"nozzle_type"`
	SDCard         bool   `json:"sdcard"`
	WifiSignal     string `json:"wifi_signal"`

	// AMS
	AMSStatus     int    `json:"ams_status"`
	AMSRFIDStatus int    `json:"ams_rfid_status"`
	AMS           AMS    `json:"ams"`
	VTTray        VTTray `json:"vt_tray"`

	// Lights
	LightsReport []LightReport `json:"lights_report"`

	// Queue
	QueueNumber int `json:"queue_number"`
	QueueTotal  int `json:"queue_total"`
}

type AMS struct {
	AMSList    []AMSUnit `json:"ams"`
	TrayNow    string    `json:"tray_now"`
	TrayTarget string    `json:"tray_tar"`
	InsertFlag bool      `json:"insert_flag"`
}

type AMSUnit struct {
	ID       string    `json:"id"`
	Humidity string    `json:"humidity"`
	Temp     string    `json:"temp"`
	Trays    []AMSTray `json:"tray"`
}

type AMSTray struct {
	ID        string `json:"id"`
	TrayType  string `json:"tray_type"`
	TrayColor string `json:"tray_color"`
	Remain    int    `json:"remain"`
}

type VTTray struct {
	TrayType      string `json:"tray_type"`
	TrayColor     string `json:"tray_color"`
	NozzleTempMin string `json:"nozzle_temp_min"`
	NozzleTempMax string `json:"nozzle_temp_max"`
	Remain        int    `json:"remain"`
}

type LightReport struct {
	Node string `json:"node"`
	Mode string `json:"mode"`
}
