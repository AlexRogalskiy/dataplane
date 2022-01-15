package privateresolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"dataplane/auth_permissions"
	"dataplane/database"
	"dataplane/database/models"
	"dataplane/logging"
	"errors"
	"os"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *mutationResolver) CreateAccessGroup(ctx context.Context, environmentID string, name string) (string, error) {
	currentUser := ctx.Value("currentUser").(string)
	platformID := ctx.Value("platformID").(string)

	// ----- Permissions
	perms := []models.Permissions{
		{Resource: "admin_platform", ResourceID: platformID, Access: "write", Subject: "user", SubjectID: currentUser, EnvironmentID: "d_platform"},
		{Resource: "admin_environment", ResourceID: environmentID, Access: "write", Subject: "user", SubjectID: currentUser, EnvironmentID: environmentID},
		{Resource: "environment_permissions", ResourceID: environmentID, Access: "write", Subject: "user", SubjectID: currentUser, EnvironmentID: environmentID},
	}

	permOutcome, _, _, _ := permissions.MultiplePermissionChecks(perms)

	if permOutcome == "denied" {
		return "", errors.New("Requires permissions.")
	}

	e := models.PermissionsAccessGroups{
		AccessGroupID: uuid.New().String(),
		Name:          name,
		EnvironmentID: environmentID,
		Active:        true,
	}

	err := database.DBConn.Create(&e).Error
	if err != nil {
		if os.Getenv("debug") == "true" {
			logging.PrintSecretsRedact(err)
		}
		return "", errors.New("Add access group database error.")
	}

	return e.AccessGroupID, nil
}

func (r *mutationResolver) DeleteAccessGroup(ctx context.Context, accessGroupID string, environmentID string) (string, error) {
	currentUser := ctx.Value("currentUser").(string)
	platformID := ctx.Value("platformID").(string)

	// ----- Permissions
	perms := []models.Permissions{
		{Resource: "admin_platform", ResourceID: platformID, Access: "write", Subject: "user", SubjectID: currentUser, EnvironmentID: "d_platform"},
		{Resource: "admin_environment", ResourceID: environmentID, Access: "write", Subject: "user", SubjectID: currentUser, EnvironmentID: environmentID},
		{Resource: "environment_permissions", ResourceID: environmentID, Access: "write", Subject: "user", SubjectID: currentUser, EnvironmentID: environmentID},
	}

	permOutcome, _, _, _ := permissions.MultiplePermissionChecks(perms)

	if permOutcome == "denied" {
		return "", errors.New("Requires permissions.")
	}

	// e := models.PermissionsAccessGroups{
	// 	AccessGroupID: accessGroupID,
	// 	EnvironmentID: environmentID,
	// }

	err := database.DBConn.Where("access_group_id =? and environment_id=?", accessGroupID, environmentID).Delete(&models.PermissionsAccessGroups{})

	if err.RowsAffected == 0 {
		return "", errors.New("User to access group relationship not found.")
	}
	if err.Error != nil {
		if os.Getenv("debug") == "true" {
			logging.PrintSecretsRedact(err)
		}
		return "", errors.New("Add access group database error.")
	}

	return accessGroupID, nil
}

func (r *mutationResolver) UpdatePermissionToAccessGroup(ctx context.Context, environmentID string, resource string, resourceID string, access string, accessGroupID string) (string, error) {
	currentUser := ctx.Value("currentUser").(string)
	platformID := ctx.Value("platformID").(string)

	/* Requires admin rights to assign admin rights either at platform or environment level */
	perms := []models.Permissions{}
	// ----- Permissions
	if resource == "admin_platform" || resource == "admin_environment" {
		perms = []models.Permissions{
			{Resource: "admin_platform", ResourceID: platformID, Access: "write", Subject: "user", SubjectID: currentUser, EnvironmentID: "d_platform"},
		}

	} else {
		perms = []models.Permissions{
			{Resource: "admin_platform", ResourceID: platformID, Access: "write", Subject: "user", SubjectID: currentUser, EnvironmentID: "d_platform"},
			{Resource: "admin_environment", ResourceID: environmentID, Access: "write", Subject: "user", SubjectID: currentUser, EnvironmentID: environmentID},
			{Resource: "environment_permissions", ResourceID: environmentID, Access: "write", Subject: "user", SubjectID: currentUser, EnvironmentID: environmentID},
		}

	}

	permOutcome, _, _, _ := permissions.MultiplePermissionChecks(perms)

	if permOutcome == "denied" {
		return "", errors.New("Requires permissions.")
	}

	perm, err := permissions.CreatePermission(
		"access_group",
		accessGroupID,
		resource,
		resourceID,
		access,
		environmentID,
		false,
	)

	if err != nil {
		if os.Getenv("debug") == "true" {
			logging.PrintSecretsRedact(err)
		}
		return "", errors.New("Add access group permission database error.")
	}

	return perm.ID, nil
}

func (r *mutationResolver) UpdatePermissionToUser(ctx context.Context, environmentID string, resource string, resourceID string, access string, userID string) (string, error) {
	currentUser := ctx.Value("currentUser").(string)
	platformID := ctx.Value("platformID").(string)

	/* Requires admin rights to assign admin rights either at platform or environment level */
	perms := []models.Permissions{}
	// ----- Permissions
	if resource == "admin_platform" || resource == "admin_environment" {
		perms = []models.Permissions{
			{Resource: "admin_platform", ResourceID: platformID, Access: "write", Subject: "user", SubjectID: currentUser, EnvironmentID: "d_platform"},
		}

	} else {
		perms = []models.Permissions{
			{Resource: "admin_platform", ResourceID: platformID, Access: "write", Subject: "user", SubjectID: currentUser, EnvironmentID: "d_platform"},
			{Resource: "admin_environment", ResourceID: environmentID, Access: "write", Subject: "user", SubjectID: currentUser, EnvironmentID: environmentID},
			{Resource: "environment_permissions", ResourceID: environmentID, Access: "write", Subject: "user", SubjectID: currentUser, EnvironmentID: environmentID},
		}

	}

	permOutcome, _, _, _ := permissions.MultiplePermissionChecks(perms)

	if permOutcome == "denied" {
		return "", errors.New("Requires permissions.")
	}

	perm, err := permissions.CreatePermission(
		"user",
		userID,
		resource,
		resourceID,
		access,
		environmentID,
		false,
	)

	if err != nil {
		if os.Getenv("debug") == "true" {
			logging.PrintSecretsRedact(err)
		}
		return "", errors.New("Add permission to user database error.")
	}

	return perm.ID, nil
}

func (r *mutationResolver) DeletePermissionToUser(ctx context.Context, userID string, permissionID string, environmentID string) (string, error) {
	currentUser := ctx.Value("currentUser").(string)
	platformID := ctx.Value("platformID").(string)

	// ----- Permissions
	perms := []models.Permissions{
		{Resource: "admin_platform", ResourceID: platformID, Access: "write", Subject: "user", SubjectID: currentUser, EnvironmentID: "d_platform"},
		{Resource: "admin_environment", ResourceID: environmentID, Access: "write", Subject: "user", SubjectID: currentUser, EnvironmentID: environmentID},
		{Resource: "environment_permissions", ResourceID: environmentID, Access: "write", Subject: "user", SubjectID: currentUser, EnvironmentID: environmentID},
	}

	permOutcome, _, _, _ := permissions.MultiplePermissionChecks(perms)

	if permOutcome == "denied" {
		return "", errors.New("Requires permissions.")
	}

	e := models.Permissions{
		ID:            permissionID,
		SubjectID:     userID,
		EnvironmentID: environmentID,
	}

	err := database.DBConn.Where("id =? and subject_id = ? and environment_id=?", permissionID, userID, environmentID).Delete(&models.Permissions{})

	if err.RowsAffected == 0 {
		return "", errors.New("User to permission relationship not found.")
	}
	if err.Error != nil {
		if os.Getenv("debug") == "true" {
			logging.PrintSecretsRedact(err)
		}
		return "", errors.New("Add access group database error.")
	}

	return e.ID, nil
}

func (r *mutationResolver) UpdateUserToAccessGroup(ctx context.Context, environmentID string, userID string, accessGroupID string) (string, error) {
	currentUser := ctx.Value("currentUser").(string)
	platformID := ctx.Value("platformID").(string)

	// ----- Permissions
	perms := []models.Permissions{
		{Resource: "admin_platform", ResourceID: platformID, Access: "write", Subject: "user", SubjectID: currentUser, EnvironmentID: "d_platform"},
		{Resource: "admin_environment", ResourceID: environmentID, Access: "write", Subject: "user", SubjectID: currentUser, EnvironmentID: environmentID},
		{Resource: "environment_permissions", ResourceID: environmentID, Access: "write", Subject: "user", SubjectID: currentUser, EnvironmentID: environmentID},
	}

	permOutcome, _, _, _ := permissions.MultiplePermissionChecks(perms)

	if permOutcome == "denied" {
		return "", errors.New("Requires permissions.")
	}

	e := models.PermissionsAccessGUsers{
		AccessGroupID: accessGroupID,
		UserID:        userID,
		EnvironmentID: environmentID,
		Active:        true,
	}

	err := database.DBConn.Create(&e).Error
	if err != nil {
		if os.Getenv("debug") == "true" {
			logging.PrintSecretsRedact(err)
		}
		return "", errors.New("Add access group database error.")
	}

	return "success", nil
}

func (r *mutationResolver) RemoveUserFromAccessGroup(ctx context.Context, userID string, accessGroupID string, environmentID string) (string, error) {
	currentUser := ctx.Value("currentUser").(string)
	platformID := ctx.Value("platformID").(string)

	// ----- Permissions
	perms := []models.Permissions{
		{Resource: "admin_platform", ResourceID: platformID, Access: "write", Subject: "user", SubjectID: currentUser, EnvironmentID: "d_platform"},
		{Resource: "admin_environment", ResourceID: environmentID, Access: "write", Subject: "user", SubjectID: currentUser, EnvironmentID: environmentID},
		{Resource: "environment_permissions", ResourceID: environmentID, Access: "write", Subject: "user", SubjectID: currentUser, EnvironmentID: environmentID},
	}

	permOutcome, _, _, _ := permissions.MultiplePermissionChecks(perms)

	if permOutcome == "denied" {
		return "", errors.New("Requires permissions.")
	}

	e := models.PermissionsAccessGUsers{
		AccessGroupID: accessGroupID,
		UserID:        userID,
		EnvironmentID: environmentID,
	}

	err := database.DBConn.Where("access_group_id =? and user_id = ? and environment_id=?", accessGroupID, userID, environmentID).Delete(&models.PermissionsAccessGUsers{})

	if err.RowsAffected == 0 {
		return "", errors.New("User to access group relationship not found.")
	}
	if err.Error != nil {
		if os.Getenv("debug") == "true" {
			logging.PrintSecretsRedact(err)
		}
		return "", errors.New("Add access group database error.")
	}

	return e.AccessGroupID, nil
}

func (r *queryResolver) AvailablePermissions(ctx context.Context, environmentID string) ([]*models.ResourceTypeStruct, error) {
	platformID := ctx.Value("platformID").(string)

	var Permissions []*models.ResourceTypeStruct

	err := database.DBConn.Raw(
		`
		(select 
		p.code,
		p.label,
		p.level,
		p.access
		from 
		permissions_resource_types p
		)
`,
		//direct
	).Scan(
		&Permissions,
	).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.New("Error retrieving permissions")
	}

	// Set resource ids
	for _, p := range Permissions {
		// log.Print(p.Level)
		if p.Level == "platform" {
			p.ResourceID = platformID
		} else if p.Level == "environment" {
			p.ResourceID = environmentID
		}
	}

	return Permissions, nil
}

func (r *queryResolver) MyPermissions(ctx context.Context) ([]*models.PermissionsOutput, error) {
	currentUser := ctx.Value("currentUser").(string)

	var PermissionsOutput []*models.PermissionsOutput

	err := database.DBConn.Raw(
		`
		(select 
		p.id,
		p.access,
		p.subject,
		p.subject_id,
		p.resource,
		p.resource_id,
		p.environment_id,
		p.active,
		pt.level,
		pt.label
		from 
		permissions p, permissions_resource_types pt
		where 
		p.resource = pt.code and
		p.subject = 'user' and 
		p.subject_id = ? and
		p.active = true
		)
		union
		(
		select
		p.id,
		p.access,
		p.subject,
		p.subject_id,
		p.resource,
		p.resource_id,
		p.environment_id,
		p.active,
		pt.level,
		pt.label
		from 
		permissions p, permissions_accessg_users agu, permissions_resource_types pt
		where 
		p.resource = pt.code and
		p.subject = 'access_group' and 
		p.subject_id = agu.user_id and
		p.subject_id = ? and
		p.active = true
		)
`,
		//direct
		currentUser,
		currentUser,
	).Scan(
		&PermissionsOutput,
	).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.New("Error retrieving permissions")
	}

	return PermissionsOutput, nil
}

func (r *queryResolver) UserPermissions(ctx context.Context, userID string, environmentID string) ([]*models.PermissionsOutput, error) {
	currentUser := ctx.Value("currentUser").(string)
	platformID := ctx.Value("platformID").(string)

	// ----- Permissions
	perms := []models.Permissions{
		{Resource: "admin_platform", ResourceID: platformID, Access: "write", Subject: "user", SubjectID: currentUser, EnvironmentID: "d_platform"},
		{Resource: "admin_environment", ResourceID: environmentID, Access: "write", Subject: "user", SubjectID: currentUser, EnvironmentID: environmentID},
		{Resource: "environment_permissions", ResourceID: environmentID, Access: "write", Subject: "user", SubjectID: currentUser, EnvironmentID: environmentID},
		{Resource: "environment_users", ResourceID: environmentID, Access: "write", Subject: "user", SubjectID: currentUser, EnvironmentID: environmentID},
		{Resource: "environment_users", ResourceID: environmentID, Access: "read", Subject: "user", SubjectID: currentUser, EnvironmentID: environmentID},
	}

	permOutcome, _, _, _ := permissions.MultiplePermissionChecks(perms)

	if permOutcome == "denied" {
		return nil, errors.New("Requires permissions.")
	}

	var PermissionsOutput []*models.PermissionsOutput

	err := database.DBConn.Raw(
		`
		(select 
		p.id,
		p.access,
		p.subject,
		p.subject_id,
		p.resource,
		p.resource_id,
		p.environment_id,
		p.active,
		pt.level,
		pt.label
		from 
		permissions p, permissions_resource_types pt
		where 
		p.resource = pt.code and
		p.subject = 'user' and 
		p.subject_id = ? and
		p.active = true		
		)
		union
		(
		select
		p.id,
		p.access,
		p.subject,
		p.subject_id,
		p.resource,
		p.resource_id,
		p.environment_id,
		p.active,
		pt.level,
		pt.label
		from 
		permissions p, permissions_accessg_users agu, permissions_resource_types pt
		where 
		p.resource = pt.code and
		p.subject = 'access_group' and 
		p.subject_id = agu.user_id and
		p.subject_id = ? and
		p.active = true		
		)
`,
		//direct
		userID,
		userID,
	).Scan(
		&PermissionsOutput,
	).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.New("Error retrieving permissions")
	}

	return PermissionsOutput, nil
}

func (r *queryResolver) GetAccessGroups(ctx context.Context, userID string, environmentID string) ([]*models.PermissionsAccessGroups, error) {
	// NEEDS PERMISSIONS <==================

	e := []*models.PermissionsAccessGroups{}

	err := database.DBConn.Find(&e).Error
	if err != nil {
		if os.Getenv("debug") == "true" {
			logging.PrintSecretsRedact(err)
		}
		return nil, errors.New("Retrive users database error.")
	}
	return e, nil
}
