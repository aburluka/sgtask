package event

type Event struct {
	ClientTime string `json:"client_time"`
	DeviceID   string `json:"device_id"`
	DeviceOs   string `json:"device_os"`
	Session    string `json:"session"`
	Sequence   int    `json:"sequence"`
	Event      string `json:"event"`
	ParamInt   int    `json:"param_int"`
	ParamStr   string `json:"param_str"`
	IP         string `json:"ip"`
	ServerTime string `json:"server_time"`
}

func (e *Event) Enrich() {
	e.IP = "8.8.8.8"
	e.ServerTime = "2020-12-01 23:53:00"
}
