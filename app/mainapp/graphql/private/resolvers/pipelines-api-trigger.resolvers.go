package privateresolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/dataplane-app/dataplane/app/mainapp/auth"
	permissions "github.com/dataplane-app/dataplane/app/mainapp/auth_permissions"
	dpconfig "github.com/dataplane-app/dataplane/app/mainapp/config"
	"github.com/dataplane-app/dataplane/app/mainapp/database"
	"github.com/dataplane-app/dataplane/app/mainapp/database/models"
	"github.com/dataplane-app/dataplane/app/mainapp/logging"
	"gorm.io/gorm/clause"
)

// GeneratePipelineTrigger is the resolver for the generatePipelineTrigger field.
func (r *mutationResolver) GeneratePipelineTrigger(ctx context.Context, pipelineID string, environmentID string, triggerID string, apiKeyActive bool, publicLive bool, privateLive bool, dataSizeLimit float64, dataTTL float64) (string, error) {
	currentUser := ctx.Value("currentUser").(string)
	platformID := ctx.Value("platformID").(string)

	// log.Println("data ttl", dataSizeLimit, dataTTL)

	// ----- Permissions
	perms := []models.Permissions{
		{Subject: "user", SubjectID: currentUser, Resource: "admin_platform", ResourceID: platformID, Access: "write", EnvironmentID: "d_platform"},
		{Subject: "user", SubjectID: currentUser, Resource: "admin_environment", ResourceID: environmentID, Access: "write", EnvironmentID: environmentID},
		{Subject: "user", SubjectID: currentUser, Resource: "environment_edit_all_pipelines", ResourceID: environmentID, Access: "write", EnvironmentID: environmentID},
		{Subject: "user", SubjectID: currentUser, Resource: "specific_pipeline", ResourceID: pipelineID, Access: "write", EnvironmentID: environmentID},
	}

	permOutcome, _, _, _ := permissions.MultiplePermissionChecks(perms)

	if permOutcome == "denied" {
		return "", errors.New("Requires permission")
	}

	trigger := models.PipelineApiTriggers{
		TriggerID:     triggerID,
		PipelineID:    pipelineID,
		EnvironmentID: environmentID,
		APIKeyActive:  apiKeyActive,
		PublicLive:    publicLive,
		PrivateLive:   privateLive,
		DataSizeLimit: dataSizeLimit,
		DataTTL:       dataTTL,
	}

	err := database.DBConn.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "trigger_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"public_live", "private_live", "api_key_active", "data_size_limit", "data_ttl"}),
	}).Create(&trigger).Error

	if err != nil {
		if dpconfig.Debug == "true" {
			logging.PrintSecretsRedact(err)
		}

		return "", err
	}

	return "Success", nil
}

// GenerateDeploymentTrigger is the resolver for the generateDeploymentTrigger field.
func (r *mutationResolver) GenerateDeploymentTrigger(ctx context.Context, deploymentID string, environmentID string, triggerID string, apiKeyActive bool, publicLive bool, privateLive bool, dataSizeLimit float64, dataTTL float64) (string, error) {
	currentUser := ctx.Value("currentUser").(string)
	platformID := ctx.Value("platformID").(string)

	// ----- Permissions
	perms := []models.Permissions{
		{Subject: "user", SubjectID: currentUser, Resource: "admin_platform", ResourceID: platformID, Access: "write", EnvironmentID: "d_platform"},
		{Subject: "user", SubjectID: currentUser, Resource: "admin_environment", ResourceID: environmentID, Access: "write", EnvironmentID: environmentID},
	}

	permOutcome, _, _, _ := permissions.MultiplePermissionChecks(perms)

	if permOutcome == "denied" {
		return "", errors.New("Requires permission")
	}

	trigger := models.DeploymentApiTriggers{
		TriggerID:     triggerID,
		DeploymentID:  deploymentID,
		EnvironmentID: environmentID,
		APIKeyActive:  apiKeyActive,
		PublicLive:    publicLive,
		PrivateLive:   privateLive,
		DataSizeLimit: dataSizeLimit,
		DataTTL:       dataTTL,
	}

	err := database.DBConn.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "trigger_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"public_live", "private_live", "api_key_active", "data_size_limit", "data_ttl"}),
	}).Create(&trigger).Error

	if err != nil {
		if dpconfig.Debug == "true" {
			logging.PrintSecretsRedact(err)
		}

		return "", err
	}

	return "Success", nil
}

// AddPipelineAPIKey is the resolver for the addPipelineApiKey field.
func (r *mutationResolver) AddPipelineAPIKey(ctx context.Context, triggerID string, apiKey string, pipelineID string, environmentID string, expiresAt *time.Time) (string, error) {
	currentUser := ctx.Value("currentUser").(string)
	platformID := ctx.Value("platformID").(string)

	// ----- Permissions
	perms := []models.Permissions{
		{Subject: "user", SubjectID: currentUser, Resource: "admin_platform", ResourceID: platformID, Access: "write", EnvironmentID: "d_platform"},
		{Subject: "user", SubjectID: currentUser, Resource: "admin_environment", ResourceID: environmentID, Access: "write", EnvironmentID: environmentID},
		{Subject: "user", SubjectID: currentUser, Resource: "environment_edit_all_pipelines", ResourceID: environmentID, Access: "write", EnvironmentID: environmentID},
		{Subject: "user", SubjectID: currentUser, Resource: "specific_pipeline", ResourceID: pipelineID, Access: "write", EnvironmentID: environmentID},
	}

	permOutcome, _, _, _ := permissions.MultiplePermissionChecks(perms)

	if permOutcome == "denied" {
		return "", errors.New("Requires permission")
	}

	//  Hash API key
	hashedApiKey, err := auth.Encrypt(apiKey)
	if err != nil {
		if dpconfig.Debug == "true" {
			logging.PrintSecretsRedact(err)
		}
		return "", errors.New("unable to hash api key")
	}

	keys := models.PipelineApiKeys{
		TriggerID:     triggerID,
		APIKey:        hashedApiKey,
		APIKeyTail:    strings.Split(apiKey, "-")[3],
		PipelineID:    pipelineID,
		EnvironmentID: environmentID,
		ExpiresAt:     expiresAt,
	}

	err = database.DBConn.Create(&keys).Error
	if err != nil {
		if dpconfig.Debug == "true" {
			logging.PrintSecretsRedact(err)
		}
		return "", errors.New("Register database error.")
	}

	return "Success", nil
}

// AddDeploymentAPIKey is the resolver for the addDeploymentApiKey field.
func (r *mutationResolver) AddDeploymentAPIKey(ctx context.Context, triggerID string, apiKey string, deploymentID string, environmentID string, expiresAt *time.Time) (string, error) {
	currentUser := ctx.Value("currentUser").(string)
	platformID := ctx.Value("platformID").(string)

	// ----- Permissions
	perms := []models.Permissions{
		{Subject: "user", SubjectID: currentUser, Resource: "admin_platform", ResourceID: platformID, Access: "write", EnvironmentID: "d_platform"},
		{Subject: "user", SubjectID: currentUser, Resource: "admin_environment", ResourceID: environmentID, Access: "write", EnvironmentID: environmentID},
	}

	permOutcome, _, _, _ := permissions.MultiplePermissionChecks(perms)

	if permOutcome == "denied" {
		return "", errors.New("Requires permission")
	}

	//  Hash API key
	hashedApiKey, err := auth.Encrypt(apiKey)
	if err != nil {
		if dpconfig.Debug == "true" {
			logging.PrintSecretsRedact(err)
		}
		return "", errors.New("unable to hash api key")
	}

	keys := models.DeploymentApiKeys{
		TriggerID:     triggerID,
		APIKey:        hashedApiKey,
		APIKeyTail:    strings.Split(apiKey, "-")[3],
		DeploymentID:  deploymentID,
		EnvironmentID: environmentID,
		ExpiresAt:     expiresAt,
	}

	err = database.DBConn.Create(&keys).Error
	if err != nil {
		if dpconfig.Debug == "true" {
			logging.PrintSecretsRedact(err)
		}
		return "", errors.New("Register database error.")
	}

	return "Success", nil
}

// DeletePipelineAPIKey is the resolver for the deletePipelineApiKey field.
func (r *mutationResolver) DeletePipelineAPIKey(ctx context.Context, apiKey string, pipelineID string, environmentID string) (string, error) {
	currentUser := ctx.Value("currentUser").(string)
	platformID := ctx.Value("platformID").(string)

	// ----- Permissions
	perms := []models.Permissions{
		{Subject: "user", SubjectID: currentUser, Resource: "admin_platform", ResourceID: platformID, Access: "write", EnvironmentID: "d_platform"},
		{Subject: "user", SubjectID: currentUser, Resource: "admin_environment", ResourceID: environmentID, Access: "write", EnvironmentID: environmentID},
		{Subject: "user", SubjectID: currentUser, Resource: "environment_edit_all_pipelines", ResourceID: environmentID, Access: "write", EnvironmentID: environmentID},
		{Subject: "user", SubjectID: currentUser, Resource: "specific_pipeline", ResourceID: pipelineID, Access: "write", EnvironmentID: environmentID},
	}

	permOutcome, _, _, _ := permissions.MultiplePermissionChecks(perms)

	if permOutcome == "denied" {
		return "", errors.New("Requires permission")
	}

	k := models.PipelineApiKeys{}

	query := database.DBConn.Where("pipeline_id = ? and environment_id = ? and api_key = ?",
		pipelineID, environmentID, apiKey).Delete(&k)
	if query.Error != nil {
		if dpconfig.Debug == "true" {
			logging.PrintSecretsRedact(query.Error)
		}
		return "", errors.New("Delete pipeline key database error.")
	}
	if query.RowsAffected == 0 {
		return "", errors.New("Delete pipeline key database error.")
	}

	return "Success", nil
}

// DeleteDeploymentAPIKey is the resolver for the deleteDeploymentApiKey field.
func (r *mutationResolver) DeleteDeploymentAPIKey(ctx context.Context, apiKey string, deploymentID string, environmentID string) (string, error) {
	currentUser := ctx.Value("currentUser").(string)
	platformID := ctx.Value("platformID").(string)

	// ----- Permissions
	perms := []models.Permissions{
		{Subject: "user", SubjectID: currentUser, Resource: "admin_platform", ResourceID: platformID, Access: "write", EnvironmentID: "d_platform"},
		{Subject: "user", SubjectID: currentUser, Resource: "admin_environment", ResourceID: environmentID, Access: "write", EnvironmentID: environmentID},
	}

	permOutcome, _, _, _ := permissions.MultiplePermissionChecks(perms)

	if permOutcome == "denied" {
		return "", errors.New("Requires permission")
	}

	k := models.DeploymentApiKeys{}

	query := database.DBConn.Where("deployment_id = ? and environment_id = ? and api_key = ?",
		deploymentID, environmentID, apiKey).Delete(&k)
	if query.Error != nil {
		if dpconfig.Debug == "true" {
			logging.PrintSecretsRedact(query.Error)
		}
		return "", errors.New("Delete deployment key database error.")
	}
	if query.RowsAffected == 0 {
		return "", errors.New("Delete deployment key database error.")
	}

	return "Success", nil
}

// GetPipelineTrigger is the resolver for the getPipelineTrigger field.
func (r *queryResolver) GetPipelineTrigger(ctx context.Context, pipelineID string, environmentID string) (*models.PipelineApiTriggers, error) {
	currentUser := ctx.Value("currentUser").(string)
	platformID := ctx.Value("platformID").(string)

	// ----- Permissions
	perms := []models.Permissions{
		{Subject: "user", SubjectID: currentUser, Resource: "admin_platform", ResourceID: platformID, Access: "write", EnvironmentID: "d_platform"},
		{Subject: "user", SubjectID: currentUser, Resource: "admin_environment", ResourceID: environmentID, Access: "write", EnvironmentID: environmentID},
		{Subject: "user", SubjectID: currentUser, Resource: "environment_edit_all_pipelines", ResourceID: environmentID, Access: "write", EnvironmentID: environmentID},
		{Subject: "user", SubjectID: currentUser, Resource: "specific_pipeline", ResourceID: pipelineID, Access: "write", EnvironmentID: environmentID},
	}

	permOutcome, _, _, _ := permissions.MultiplePermissionChecks(perms)

	if permOutcome == "denied" {
		return nil, errors.New("Requires permission")
	}

	e := models.PipelineApiTriggers{}

	err := database.DBConn.Where("pipeline_id = ? and environment_id = ?", pipelineID, environmentID).First(&e).Error
	if err != nil {
		if err.Error() == "record not found" {
			return nil, errors.New("record not found")
		}
		if dpconfig.Debug == "true" {
			logging.PrintSecretsRedact(err)
		}
		return nil, errors.New("Retrive pipeline trigger database error.")
	}

	return &e, nil
}

// GetDeploymentTrigger is the resolver for the getDeploymentTrigger field.
func (r *queryResolver) GetDeploymentTrigger(ctx context.Context, deploymentID string, environmentID string) (*models.DeploymentApiTriggers, error) {
	currentUser := ctx.Value("currentUser").(string)
	platformID := ctx.Value("platformID").(string)

	// ----- Permissions
	perms := []models.Permissions{
		{Subject: "user", SubjectID: currentUser, Resource: "admin_platform", ResourceID: platformID, Access: "write", EnvironmentID: "d_platform"},
		{Subject: "user", SubjectID: currentUser, Resource: "admin_environment", ResourceID: environmentID, Access: "write", EnvironmentID: environmentID},
		{Subject: "user", SubjectID: currentUser, Resource: "environment_run_all_pipelines", ResourceID: environmentID, Access: "write", EnvironmentID: environmentID},
		{Subject: "user", SubjectID: currentUser, Resource: "specific_pipeline", ResourceID: deploymentID, Access: "run", EnvironmentID: environmentID},
		{Subject: "user", SubjectID: currentUser, Resource: "specific_pipeline", ResourceID: deploymentID, Access: "deploy", EnvironmentID: environmentID},
	}

	permOutcome, _, _, _ := permissions.MultiplePermissionChecks(perms)

	if permOutcome == "denied" {
		return nil, errors.New("Requires permission")
	}

	e := models.DeploymentApiTriggers{}

	err := database.DBConn.Where("deployment_id = ? and environment_id = ?", deploymentID, environmentID).First(&e).Error
	if err != nil {
		if err.Error() == "record not found" {
			return nil, errors.New("record not found")
		}
		if dpconfig.Debug == "true" {
			logging.PrintSecretsRedact(err)
		}
		return nil, errors.New("Retrive deployment trigger database error.")
	}

	return &e, nil
}

// GetPipelineAPIKeys is the resolver for the getPipelineApiKeys field.
func (r *queryResolver) GetPipelineAPIKeys(ctx context.Context, pipelineID string, environmentID string) ([]*models.PipelineApiKeys, error) {
	currentUser := ctx.Value("currentUser").(string)
	platformID := ctx.Value("platformID").(string)

	// ----- Permissions
	perms := []models.Permissions{
		{Subject: "user", SubjectID: currentUser, Resource: "admin_platform", ResourceID: platformID, Access: "write", EnvironmentID: "d_platform"},
		{Subject: "user", SubjectID: currentUser, Resource: "admin_environment", ResourceID: environmentID, Access: "write", EnvironmentID: environmentID},
	}

	permOutcome, _, _, _ := permissions.MultiplePermissionChecks(perms)

	if permOutcome == "denied" {
		return nil, errors.New("Requires permission")
	}

	e := []*models.PipelineApiKeys{}

	err := database.DBConn.Where("pipeline_id = ? and environment_id = ?", pipelineID, environmentID).Find(&e).Error
	if err != nil {
		if dpconfig.Debug == "true" {
			logging.PrintSecretsRedact(err)
		}
		return nil, errors.New("Retrive pipeline trigger database error.")
	}
	return e, nil
}

// GetDeploymentAPIKeys is the resolver for the getDeploymentApiKeys field.
func (r *queryResolver) GetDeploymentAPIKeys(ctx context.Context, deploymentID string, environmentID string) ([]*models.DeploymentApiKeys, error) {
	currentUser := ctx.Value("currentUser").(string)
	platformID := ctx.Value("platformID").(string)

	// ----- Permissions
	perms := []models.Permissions{
		{Subject: "user", SubjectID: currentUser, Resource: "admin_platform", ResourceID: platformID, Access: "write", EnvironmentID: "d_platform"},
		{Subject: "user", SubjectID: currentUser, Resource: "admin_environment", ResourceID: environmentID, Access: "write", EnvironmentID: environmentID},
	}

	permOutcome, _, _, _ := permissions.MultiplePermissionChecks(perms)

	if permOutcome == "denied" {
		return nil, errors.New("Requires permission")
	}

	e := []*models.DeploymentApiKeys{}

	err := database.DBConn.Where("deployment_id = ? and environment_id = ?", deploymentID, environmentID).Find(&e).Error
	if err != nil {
		if dpconfig.Debug == "true" {
			logging.PrintSecretsRedact(err)
		}
		return nil, errors.New("Retrive deployment trigger database error.")
	}
	return e, nil
}
