package weberr

const (
	ALARM_LEVEL_NOTICE   = "notice"
	ALARM_LEVEL_NORMAL   = "normal"
	ALARM_LEVEL_CRITICAL = "critical"
	ALARM_LEVEL_DEADLY   = "deadly"
)

type PublicAlarm struct {
	ID      int         `json:"id"`
	Type    string      `json:"type"`
	Level   string      `json:"level"`
	Detail  string      `json:"detail"`
}

var (
	ALARM_SYS_FILE = PublicAlarm{ID: 1001, Type: "OS", Level: ALARM_LEVEL_CRITICAL, Detail: "filesystem not available"}
	ALARM_SYS_CPU  = PublicAlarm{ID: 1002, Type: "OS", Level: ALARM_LEVEL_NORMAL, Detail: "CPU overload"}
	ALARM_SYS_MEM  = PublicAlarm{ID: 1003, Type: "OS", Level: ALARM_LEVEL_CRITICAL, Detail: "memory overload"}
	ALARM_SYS_DISK  = PublicAlarm{ID: 1004, Type: "OS", Level: ALARM_LEVEL_CRITICAL, Detail: "disk overload"}
	ALARM_PROC_ETCD  = PublicAlarm{ID: 2001, Type: "SYS", Level: ALARM_LEVEL_CRITICAL, Detail: "etcd not available"}
)
