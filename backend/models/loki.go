package models

import (
	"time"
)

type LokiMessage struct {
	Timestamp string
	Message   string
}

type LokiMessages []LokiMessage

func (msg LokiMessages) Len() int { return len(msg) }

func (msg LokiMessages) Less(i, j int) bool {
	return msg[i].Timestamp > msg[j].Timestamp
}

func (msg LokiMessages) Swap(i, j int) { msg[i], msg[j] = msg[j], msg[i] }

// LogHistory [...]
type LogHistory struct {
	ID         int       `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	LabelJSON  string    `gorm:"column:label_json;type:text" json:"label_json"`
	CreateAt   time.Time `gorm:"column:create_at;type:datetime" json:"create_at"`
	FilterJSON string    `gorm:"column:filter_json;type:text" json:"filter_json"`
	UserID     int       `gorm:"index;column:user_id;" json:"user_id"`
	User       User      `gorm:"foreignkey:UserID;" json:"user"`
	LogQL      string    `gorm:"column:log_ql;type:varchar(2048);" json:"log_ql"`
}

func (LogHistory) TableName() string {
	return "log_history"
}

// LogSnapshot [...]
type LogSnapshot struct {
	ID          int       `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Name        string    `gorm:"column:name;type:varchar(128);not null" json:"name"`
	Count       int       `gorm:"column:count;type:int(11)" json:"count"`
	CreateAt    time.Time `gorm:"column:create_at;type:datetime" json:"create_at"`
	DownloadURL string    `gorm:"column:download_url;type:varchar(512)" json:"download_url"`
	UserID      int       `gorm:"index;column:user_id;" json:"user_id"`
	User        User      `gorm:"foreignkey:UserID;" json:"user"`
	StartTime   time.Time `gorm:"column:start_time;type:datetime" json:"start_time"`
	EndTime     time.Time `gorm:"column:end_time;type:datetime" json:"end_time"`
	Dir         string    `gorm:"column:dir;type:varchar(128)" json:"dir"`
}

func (LogSnapshot) TableName() string {
	return "log_snapshot"
}
