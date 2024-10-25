package main

const (
	default_log_size       int    = 32
	default_log_backups    int    = 10
	default_log_age        int    = 60
	default_log_to_console bool   = false
	default_service_port   string = "33470"

	RFC3339Lite     string = "2006-01-02T15:04:05"
	RFC3339LiteDate string = "2006-01-02"
	// RFC3339LiteDateNoHyphen string = "20060102"
)

// const (
// 	// 工作狀態
// 	// (0: 初始), (3: 施工中), (5: 待上傳), (11: 完工)
// 	Stat_0        int = 0
// 	Stat_Doing    int = 3
// 	Stat_Waiting  int = 5
// 	Stat_Complete int = 11
// 	Stat_Abort    int = 31
// )

const (
	DBName string = "weeping-willow"

	Order_Asc  int = 1
	Order_Desc int = -1

	// Coll_Account string = "account"
	// Coll_Node    string = "node"
)
