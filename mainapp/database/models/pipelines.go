package models

import (
	"time"

	"gorm.io/datatypes"
)

func (Pipelines) IsEntity() {}

func (Pipelines) TableName() string {
	return "pipelines"
}

type Pipelines struct {
	PipelineID  string         `gorm:"PRIMARY_KEY;type:varchar(64);" json:"pipeline_id"`
	Name        string         `gorm:"type:varchar(255);index:idx_pipelines,unique;" json:"name"`
	Version     string         `gorm:"type:varchar(125);index:idx_pipelines,unique;" json:"version"`
	YAMLHash    string         `json:"yaml_hash"`
	Description string         `json:"description"`
	Active      bool           `json:"active"`
	Online      bool           `json:"online"`
	FileJSON    datatypes.JSON `json:"file_json"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   *time.Time     `json:"updated_at"`
	DeletedAt   *time.Time     `json:"deleted_at,omitempty"`
}

func (PipelinesHistory) IsEntity() {}

func (PipelinesHistory) TableName() string {
	return "pipelines_history"
}

type PipelinesHistory struct {
	PipelineID  string         `gorm:"PRIMARY_KEY;type:varchar(64);" json:"pipeline_id"`
	Name        string         `gorm:"type:varchar(255);index:idx_pipelines_history,unique;" json:"name"`
	Version     string         `gorm:"type:varchar(125);index:idx_pipelines_history,unique;" json:"version"`
	YAMLHash    string         `json:"yaml_hash"`
	Description string         `json:"description"`
	Active      bool           `json:"active"`
	Online      bool           `json:"online"`
	FileJSON    datatypes.JSON `json:"file_json"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   *time.Time     `json:"updated_at"`
	DeletedAt   *time.Time     `json:"deleted_at,omitempty"`
}
