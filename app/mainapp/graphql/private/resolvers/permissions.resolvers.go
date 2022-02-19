package privateresolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	permissions "dataplane/mainapp/auth_permissions"
	"dataplane/mainapp/database"
	"dataplane/mainapp/database/models"
	"dataplane/mainapp/logging"
	"errors"
	"log"
	"os"

	"gorm.io/gorm"
)

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

func (r *mutationResolver) PipelinePermissionsToUser(ctx context.Context, environmentID string, resource string, resourceID string, access string, userID string, checked string) (string, error) {
	currentUser := ctx.Value("currentUser").(string)
	platformID := ctx.Value("platformID").(string)

	// ----- Permissions
	perms := []models.Permissions{
		{Subject: "user", SubjectID: currentUser, Resource: "admin_platform", ResourceID: platformID, Access: "write", EnvironmentID: "d_platform"},
		{Subject: "user", SubjectID: currentUser, Resource: "platform_environment", ResourceID: platformID, Access: "write", EnvironmentID: environmentID},
		{Subject: "user", SubjectID: currentUser, Resource: "environment_edit_all_pipelines", ResourceID: platformID, Access: "write", EnvironmentID: environmentID},
		{Subject: "user", SubjectID: currentUser, Resource: "specific_pipeline", ResourceID: resourceID, Access: "write", EnvironmentID: environmentID},
	}

	permOutcome, _, _, _ := permissions.MultiplePermissionChecks(perms)

	if permOutcome == "denied" {
		return "", errors.New("Requires permissions.")
	}

	if checked == "no" {
		err := database.DBConn.Where("subject_id = ? and resource = ? and resource_id = ? and access = ? and environment_id = ?",
			userID, resource, resourceID, access, environmentID).Delete(&models.Permissions{})

		// Fix below!!!
		log.Println(err)
		// if err != nil {
		// 	if os.Getenv("debug") == "true" {
		// 		logging.PrintSecretsRedact(err)
		// 	}
		// 	return "", errors.New("Delete pipelines permission to user database error.")
		// }

		return "Permission deleted.", nil
	}

	if checked == "yes" {

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

	return "", errors.New("Check must be yes or no")
}

func (r *mutationResolver) PipelinePermissionsToAccessGroup(ctx context.Context, environmentID string, resource string, resourceID string, access string, accessGroupID string, checked string) (string, error) {
	currentUser := ctx.Value("currentUser").(string)
	platformID := ctx.Value("platformID").(string)

	// ----- Permissions
	perms := []models.Permissions{
		{Subject: "user", SubjectID: currentUser, Resource: "admin_platform", ResourceID: platformID, Access: "write", EnvironmentID: "d_platform"},
		{Subject: "user", SubjectID: currentUser, Resource: "platform_environment", ResourceID: platformID, Access: "write", EnvironmentID: environmentID},
		{Subject: "user", SubjectID: currentUser, Resource: "environment_edit_all_pipelines", ResourceID: platformID, Access: "write", EnvironmentID: environmentID},
		{Subject: "user", SubjectID: currentUser, Resource: "specific_pipeline", ResourceID: resourceID, Access: "write", EnvironmentID: environmentID},
	}

	permOutcome, _, _, _ := permissions.MultiplePermissionChecks(perms)

	if permOutcome == "denied" {
		return "", errors.New("Requires permissions.")
	}

	if checked == "no" {
		err := database.DBConn.Where("subject_id = ? and resource = ? and resource_id = ? and access = ?",
			accessGroupID, resource, resourceID, access).Delete(&models.Permissions{})

		// Fix below!!!
		log.Println(err)
		// if err != nil {
		// 	if os.Getenv("debug") == "true" {
		// 		logging.PrintSecretsRedact(err)
		// 	}
		// 	return "", errors.New("Delete pipelines permission to user database error.")
		// }

		return "Permission deleted.", nil
	}

	if checked == "yes" {

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
			return "", errors.New("Add permission to user database error.")
		}

		return perm.ID, nil
	}

	return "", errors.New("Check must be yes or no")
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
