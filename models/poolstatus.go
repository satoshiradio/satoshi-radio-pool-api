package models

type PoolStatus struct {
	Runtime      int     `json:"runtime"`
	LastUpdate   int64   `json:"lastupdate"`
	Users        int     `json:"Users"`
	Workers      int     `json:"Workers"`
	Idle         int     `json:"Idle"`
	Disconnected int     `json:"Disconnected"`
	Hashrate1m   string  `json:"hashrate1m"`
	Hashrate5m   string  `json:"hashrate5m"`
	Hashrate15m  string  `json:"hashrate15m"`
	Hashrate1hr  string  `json:"hashrate1hr"`
	Hashrate6hr  string  `json:"hashrate6hr"`
	Hashrate1d   string  `json:"hashrate1d"`
	Hashrate7d   string  `json:"hashrate7d"`
	Diff         float64 `json:"diff"`
	Accepted     int     `json:"accepted"`
	Rejected     int     `json:"rejected"`
	BestShare    int     `json:"bestshare"`
	SPS1m        float64 `json:"SPS1m"`
	SPS5m        float64 `json:"SPS5m"`
	SPS15m       float64 `json:"SPS15m"`
	SPS1h        float64 `json:"SPS1h"`
}
