package main

// Schedule struct
type Schedule struct {
	TaskCodename string
	Expression   string
	Timeout      int64
	Arguments    string
}

// ScheduleInfo struct
type ScheduleInfo struct {
	IpAddress string `json:"ip_address"`
	StartTime string `json:"start_time"`
}
