package model

import "time"

// OperationLog 操作日志模型
type OperationLog struct {
	ID          uint      `json:"id" gorm:"primaryKey;comment:主键ID"`
	AdminID     uint      `json:"admin_id" gorm:"comment:管理员ID"`
	Username    string    `json:"username" gorm:"size:50;comment:管理员用户名"`
	Method      string    `json:"method" gorm:"size:10;comment:请求方法"`
	Path        string    `json:"path" gorm:"size:255;comment:请求路径"`
	Action      string    `json:"action" gorm:"size:50;comment:操作描述"`
	TargetID    string    `json:"target_id" gorm:"size:50;comment:操作目标ID"`
	RequestBody string    `json:"request_body" gorm:"type:longtext;comment:请求数据"`
	IP          string    `json:"ip" gorm:"size:50;comment:请求IP"`
	UserAgent   string    `json:"user_agent" gorm:"size:255;comment:用户代理"`
	Status      int       `json:"status" gorm:"comment:响应状态码"`
	CostTime    int64     `json:"cost_time" gorm:"comment:耗时(毫秒)"`
	CreatedAt   time.Time `json:"created_at" gorm:"comment:操作时间"`
}

// TableName 指定表名
func (OperationLog) TableName() string {
	return "operation_logs"
}
