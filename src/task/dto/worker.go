package dto

type WorkerInfo struct {
	Modular   string `json:"modular"`
	Address   string `json:"address"`
	WorkerId  string `json:"worker_id"`
	StartTime string `json:"start_time"`
}
