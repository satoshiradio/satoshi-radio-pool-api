package models

type User struct {
	Hashrate1m  string  `json:"hashrate1m"`
	Hashrate5m  string  `json:"hashrate5m"`
	Hashrate1hr string  `json:"hashrate1hr"`
	Hashrate1d  string  `json:"hashrate1d"`
	Hashrate7d  string  `json:"hashrate7d"`
	LastShare   int64   `json:"lastshare"`
	Workers     int     `json:"workers"`
	Shares      int     `json:"shares"`
	BestShare   float64 `json:"bestshare"`
	BestEver    int     `json:"bestever"`
	Authorised  int64   `json:"authorised"`
	Worker      []struct {
		WorkerName  string  `json:"workername"`
		Hashrate1m  string  `json:"hashrate1m"`
		Hashrate5m  string  `json:"hashrate5m"`
		Hashrate1hr string  `json:"hashrate1hr"`
		Hashrate1d  string  `json:"hashrate1d"`
		Hashrate7d  string  `json:"hashrate7d"`
		LastShare   int64   `json:"lastshare"`
		Shares      int     `json:"shares"`
		BestShare   float64 `json:"bestshare"`
		BestEver    int     `json:"bestever"`
	} `json:"worker"`
}
