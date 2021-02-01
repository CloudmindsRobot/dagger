package models

import (
	"time"

	"gorm.io/gorm"
)

type LogLabel struct {
	ID      int     `gorm:"primary_key;AUTO_INCREMENT;" json:"-"`
	Key     string  `gorm:"column:key;type:varchar(128)" json:"key"`
	Value   string  `gorm:"column:value;type:varchar(128)" json:"value"`
	LogRule LogRule `gorm:"foreignkey:RuleID;" json:"rule"`
	RuleID  int     `gorm:"index;column:rule_id;" json:"rule_id"`
}

func (LogLabel) TableName() string {
	return "log_label"
}

type LogGroup struct {
	ID             int          `gorm:"primary_key;AUTO_INCREMENT;" json:"-"`
	LogUserGroup   LogUserGroup `gorm:"foreignkey:LogUserGroupID;" json:"log_user_group"`
	LogUserGroupID int          `gorm:"index;column:group_id;" json:"id"`
	LogRule        LogRule      `gorm:"foreignkey:RuleID;" json:"rule"`
	RuleID         int          `gorm:"index;column:rule_id;" json:"rule_id"`
}

func (LogGroup) TableName() string {
	return "log_group"
}

type LogRule struct {
	ID          int        `gorm:"primary_key;AUTO_INCREMENT;" json:"id"`
	CreateAt    time.Time  `gorm:"column:create_at;type:datetime" json:"create_at"`
	UpdateAt    time.Time  `gorm:"column:update_at;type:datetime" json:"update_at"`
	Key         string     `gorm:"column:key;type:varchar(128);not null;unique;" json:"key"`
	Name        string     `gorm:"column:name;type:varchar(64)" json:"name"`
	Description string     `gorm:"column:description;type:varchar(2056)" json:"description"`
	Summary     string     `gorm:"column:summary;type:varchar(2056)" json:"summary"`
	LogQL       string     `gorm:"column:log_ql;type:varchar(512)" json:"log_ql"`
	UserID      int        `gorm:"index;column:user_id;" json:"user_id"`
	User        User       `gorm:"foreignkey:UserID;" json:"user"`
	Labels      []LogLabel `gorm:"foreignkey:rule_id;" json:"labels"`
	Groups      []LogGroup `gorm:"foreignkey:rule_id;" json:"groups"`
	Level       string     `gorm:"column:level;type:varchar(32)" json:"level"`
}

func (LogRule) TableName() string {
	return "log_rule"
}

func (rule *LogRule) BeforeDelete(tx *gorm.DB) (err error) {
	tx.Delete(&LogLabel{}, "rule_id = ?", rule.ID)
	tx.Delete(&LogGroup{}, "rule_id = ?", rule.ID)
	tx.Delete(&LogEventDetail{}, "rule_id = ?", rule.ID)
	tx.Delete(&LogEvent{}, "rule_id = ?", rule.ID)
	return
}

func (rule *LogRule) BeforeUpdate(tx *gorm.DB) (err error) {
	tx.Delete(&LogLabel{}, "rule_id = ?", rule.ID)
	tx.Delete(&LogGroup{}, "rule_id = ?", rule.ID)
	return
}

type LogUserGroup struct {
	ID        int        `gorm:"primary_key;AUTO_INCREMENT;" json:"id"`
	CreateAt  time.Time  `gorm:"column:create_at;type:datetime" json:"create_at"`
	GroupName string     `gorm:"column:group_name;type:varchar(64)" json:"group_name"`
	Users     []LogUser  `gorm:"foreignkey:group_id;" json:"users"`
	Groups    []LogGroup `gorm:"foreignkey:group_id;" json:"rule_groups"`
	UserID    int        `gorm:"index;column:user_id;" json:"user_id"`
	User      User       `gorm:"foreignkey:UserID;" json:"user"`
}

func (LogUserGroup) TableName() string {
	return "log_user_group"
}

func (group *LogUserGroup) BeforeDelete(tx *gorm.DB) (err error) {
	tx.Delete(&LogUser{}, "group_id = ?", group.ID)
	return
}

func (group *LogUserGroup) BeforeUpdate(tx *gorm.DB) (err error) {
	tx.Delete(&LogUser{}, "group_id = ?", group.ID)
	return
}

type LogUser struct {
	ID             int          `gorm:"primary_key;AUTO_INCREMENT;" json:"-"`
	LogUserGroup   LogUserGroup `gorm:"foreignkey:LogUserGroupID;" json:"log_user_group"`
	LogUserGroupID int          `gorm:"index;column:group_id;" json:"group_id"`
	User           User         `gorm:"foreignkey:UserID;" json:"user"`
	UserID         int          `gorm:"index;column:user_id;" json:"id"`
}

func (LogUser) TableName() string {
	return "log_user"
}

type LogEventDetail struct {
	ID          int       `gorm:"primary_key;AUTO_INCREMENT;" json:"id"`
	StartsAt    time.Time `gorm:"column:starts_at;type:datetime(6);index" json:"starts_at"`
	Summary     string    `gorm:"column:summary;type:longtext" json:"summary"`
	Labels      string    `gorm:"column:labels;type:longtext" json:"labels"`
	Description string    `gorm:"column:description;type:longtext;not null" json:"description"`
	LogEvent    LogEvent  `gorm:"foreignkey:EventID;" json:"event"`
	LogRule     LogRule   `gorm:"foreignkey:RuleID;" json:"rule"`
	RuleID      *int      `gorm:"index;column:rule_id;" json:"rule_id"`
	EventID     *int      `gorm:"index;column:event_id;" json:"event_id"`
	Level       string    `gorm:"column:level;type:varchar(32)" json:"level"`
}

func (LogEventDetail) TableName() string {
	return "log_event_detail"
}

type LogEvent struct {
	ID        int              `gorm:"primary_key;AUTO_INCREMENT;" json:"id"`
	ResolveAt *time.Time       `gorm:"column:resolve_at;type:datetime(6)" json:"resolve_at"`
	CreateAt  time.Time        `gorm:"column:create_at;type:datetime" json:"create_at"`
	LogRule   LogRule          `gorm:"foreignkey:RuleID;" json:"rule"`
	RuleID    *int             `gorm:"index;column:rule_id;" json:"rule_id"`
	Details   []LogEventDetail `gorm:"foreignkey:event_id;" json:"details"`
	Status    string           `gorm:"column:status;type:varchar(24);not null" json:"status"`
	Count     int              `gorm:"column:count;type:int(11)" json:"count"`
}

func (LogEvent) TableName() string {
	return "log_event"
}
