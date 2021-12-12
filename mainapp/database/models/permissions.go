package models

import "time"

// Constraints

/*
Subject types are users or access groups to which they belong.
Who is requesting access?
*/
var SubjectType = []string{"user", "access_group", "server"}

/*
Resource types are resources that users want access to.
Requesting access to what?
*/

type ResourceTypeStruct struct {
	Code  string
	Label string
	Level string
}

var ResourceType = []*ResourceTypeStruct{

	// Platform level
	{Code: "admin_platform", Level: "platform", Label: "Admin"},
	{Code: "platform_environment", Level: "platform", Label: "Manage environments"},

	// Environment level
	{Code: "admin_environment", Level: "environment", Label: "Admin"},
	// To add an admin user - you will need admin rights
	{Code: "environment_users", Level: "environment", Label: "Manage users"},
	{Code: "environment_permissions", Level: "environment", Label: "Manage permissions"},
	{Code: "environment_all_pipelines", Level: "environment", Label: "View all pipelines"},
	{Code: "environment_secrets", Level: "environment", Label: "Manage secrets"},
	{Code: "environment_edit_workers", Level: "environment", Label: "Manage workers"},

	// Specific level
	{Code: "specific_worker", Level: "specific", Label: "Worker - ${{worker_name}}"},
	{Code: "specific_pipeline", Level: "specific", Label: "Pipeline - ${{pipeline_name}}"},
}

/* Access: what type of access does the user have to the resource - read, write */
var AccessTypes = []string{"read", "write"}

// -------------- Permissions

func (Permissions) IsEntity() {}

func (Permissions) TableName() string {
	return "permissions"
}

type Permissions struct {
	ID string `gorm:"PRIMARY_KEY;type:varchar(64);" json:"id" validate:"required"`

	// Who requires access - user, server, access_group
	Subject   string `gorm:"index:idx_permissions,unique;type:varchar(64);" json:"subject" validate:"required"`
	SubjectID string `gorm:"index:idx_permissions,unique;type:varchar(64);" json:"subject_id" validate:"required"`

	// To which resource
	Resource   string `gorm:"index:idx_permissions,unique;type:varchar(64);" json:"resource" validate:"required"`
	ResourceID string `gorm:"index:idx_permissions,unique;type:varchar(64);" json:"resource_id" validate:"required"`

	// Type of access - read or write
	Access        string    `gorm:"index:idx_permissions,unique;type:varchar(64);" json:"access" validate:"required"`
	Active        bool      `json:"active" validate:"required"`
	EnvironmentID string    `json:"environment_id"`
	Test          string    `json:"test"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// -------------- Access groups

func (PermissionsAccessGroups) IsEntity() {}

func (PermissionsAccessGroups) TableName() string {
	return "permissions_access_groups"
}

type PermissionsAccessGroups struct {
	AccessGroupID string `gorm:"PRIMARY_KEY;type:varchar(64);" json:"access_group_id" validate:"required"`

	// Who requires access - user, server, access_group
	Name          string    `gorm:"index:idx_ag,unique;type:varchar(255);" json:"name" validate:"required"`
	Active        bool      `json:"active" validate:"required"`
	EnvironmentID string    `json:"environment_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// -------------- The mapping of access groups to permissions
func (PermissionsAccessGPerms) IsEntity() {}

func (PermissionsAccessGPerms) TableName() string {
	return "permissions_accessg_perms"
}

type PermissionsAccessGPerms struct {
	AccessGroupID string    `gorm:"PRIMARY_KEY;type:varchar(64);" json:"access_group_id" validate:"required"`
	PermissionID  string    `gorm:"PRIMARY_KEY;type:varchar(64);" json:"permission_id" validate:"required"`
	Active        bool      `json:"active" validate:"required"`
	EnvironmentID string    `json:"environment_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// -------------- The mapping of access groups to users
func (PermissionsAccessGUsers) IsEntity() {}

func (PermissionsAccessGUsers) TableName() string {
	return "permissions_accessg_users"
}

type PermissionsAccessGUsers struct {
	AccessGroupID string    `gorm:"PRIMARY_KEY;type:varchar(64);" json:"access_group_id" validate:"required"`
	UserID        string    `gorm:"PRIMARY_KEY;type:varchar(64);" json:"user_id" validate:"required"`
	Active        bool      `json:"active" validate:"required"`
	EnvironmentID string    `json:"environment_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
