package privateresolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	permissions "dataplane/mainapp/auth_permissions"
	"dataplane/mainapp/config"
	"dataplane/mainapp/database"
	"dataplane/mainapp/database/models"
	"dataplane/mainapp/filesystem"
	privategraphql "dataplane/mainapp/graphql/private"
	"dataplane/mainapp/logging"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
)

func (r *mutationResolver) CreateFolderNode(ctx context.Context, input *privategraphql.FolderNodeInput) (*models.CodeFolders, error) {
	currentUser := ctx.Value("currentUser").(string)
	platformID := ctx.Value("platformID").(string)

	// ----- Permissions
	perms := []models.Permissions{
		{Subject: "user", SubjectID: currentUser, Resource: "admin_platform", ResourceID: platformID, Access: "write", EnvironmentID: "d_platform"},
		{Subject: "user", SubjectID: currentUser, Resource: "platform_environment", ResourceID: platformID, Access: "write", EnvironmentID: input.EnvironmentID},
		{Subject: "user", SubjectID: currentUser, Resource: "environment_edit_all_pipelines", ResourceID: platformID, Access: "write", EnvironmentID: input.EnvironmentID},
		{Subject: "user", SubjectID: currentUser, Resource: "specific_pipeline", ResourceID: input.PipelineID, Access: "write", EnvironmentID: input.EnvironmentID},
	}

	permOutcome, _, _, _ := permissions.MultiplePermissionChecks(perms)

	if permOutcome == "denied" {
		return &models.CodeFolders{}, errors.New("Requires permissions.")
	}

	// ----- Add node files to database

	f := models.CodeFolders{
		EnvironmentID: input.EnvironmentID,
		PipelineID:    input.PipelineID,
		NodeID:        input.NodeID,
		FolderID:      input.FolderID,
		ParentID:      input.ParentID,
		FolderName:    input.FolderName,
		FType:         input.FType,
		Level:         uuid.NewString(),
		Active:        input.Active,
	}

	parentFolder, err := filesystem.FolderConstructByID(database.DBConn, input.ParentID)
	if err != nil {
		return &models.CodeFolders{}, errors.New("Create folder - build parent folder failed")
	}

	rFolderout, _, err := filesystem.CreateFolder(f, parentFolder)
	if err != nil {
		return &models.CodeFolders{}, errors.New("Create folder error")
	}

	return &rFolderout, nil
}

func (r *mutationResolver) MoveFolderNode(ctx context.Context, folderID string, toFolderID string) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteFolderNode(ctx context.Context, environmentID string, folderID string, nodeID string, pipelineID string) (string, error) {
	currentUser := ctx.Value("currentUser").(string)
	platformID := ctx.Value("platformID").(string)

	// ----- Permissions
	perms := []models.Permissions{
		{Subject: "user", SubjectID: currentUser, Resource: "admin_platform", ResourceID: platformID, Access: "write", EnvironmentID: "d_platform"},
		{Subject: "user", SubjectID: currentUser, Resource: "platform_environment", ResourceID: platformID, Access: "write", EnvironmentID: environmentID},
		{Subject: "user", SubjectID: currentUser, Resource: "environment_edit_all_pipelines", ResourceID: platformID, Access: "write", EnvironmentID: environmentID},
		{Subject: "user", SubjectID: currentUser, Resource: "specific_pipeline", ResourceID: pipelineID, Access: "write", EnvironmentID: environmentID},
	}

	permOutcome, _, _, _ := permissions.MultiplePermissionChecks(perms)

	if permOutcome == "denied" {
		return "", errors.New("Requires permissions.")
	}

	folderpath, _ := filesystem.FolderConstructByID(database.DBConn, folderID)

	// Make sure there is a path
	if strings.TrimSpace(folderpath) == "" {
		return "", errors.New("Missing folder path.")
	}

	// 1. ----- Delete folder and all its contents from directory
	err := os.RemoveAll(config.CodeDirectory + folderpath)
	if err != nil {
		return "", errors.New(err.Error())
	}

	// 2. ----- Delete folder and all its contents from the database

	// Folder ids of to be deleted folders are stored here. Initialized with parent folder
	foldersIDsToDelete := []string{folderID}

	// Initialize recursive function to walk the folder to be deleted
	// and get ids for all child folders
	var getFolderIdsFromDB func(folderID string) (string, error)

	fo := []models.CodeFolders{}

	getFolderIdsFromDB = func(folderID string) (string, error) {

		// Get folder's child folders
		err := database.DBConn.Where("parent_id = ?", folderID).Find(&fo).Error
		if err != nil {
			return "", errors.New(err.Error())
		}

		// Append ids of all child folders to the to be deleted slice
		for _, item := range fo {
			foldersIDsToDelete = append(foldersIDsToDelete, item.FolderID)
		}

		// Recursively re-run the function for all child folders
		for _, item := range fo {
			getFolderIdsFromDB(item.FolderID)
		}

		return "", nil

	}

	getFolderIdsFromDB(folderID)

	// File ids of to be deleted files are stored here.
	fileIDsToDelete := []string{}

	for _, item := range foldersIDsToDelete {
		f := []models.CodeFiles{}

		// Find all files in a folder to be deleted
		err := database.DBConn.Where("folder_id = ?", item).Find(&f).Error
		if err != nil {
			return "", errors.New(err.Error())
		}

		// Append ids of all child files to the to be deleted slice
		for _, item := range f {
			fileIDsToDelete = append(fileIDsToDelete, item.FileID)
		}

	}

	// Delete folders from the database
	for _, id := range foldersIDsToDelete {
		f := []models.CodeFolders{}

		err := database.DBConn.Where("folder_id = ?", id).Delete(&f).Error
		if err != nil {
			if os.Getenv("debug") == "true" {
				logging.PrintSecretsRedact(err)
			}
			return "", errors.New("Delete folder database error.")
		}
	}

	// Delete files from the database
	for _, id := range fileIDsToDelete {
		f := []models.CodeFiles{}

		err := database.DBConn.Where("file_id = ?", id).Delete(&f).Error
		if err != nil {
			if os.Getenv("debug") == "true" {
				logging.PrintSecretsRedact(err)
			}
			return "", errors.New("Delete file database error.")
		}
	}

	return "Success", nil
}

func (r *mutationResolver) RenameFolder(ctx context.Context, environmentID string, folderID string, nodeID string, pipelineID string) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UploadFileNode(ctx context.Context, environmentID string, nodeID string, pipelineID string, folderID string, file graphql.Upload) (string, error) {
	currentUser := ctx.Value("currentUser").(string)
	platformID := ctx.Value("platformID").(string)

	// ----- Permissions
	perms := []models.Permissions{
		{Subject: "user", SubjectID: currentUser, Resource: "admin_platform", ResourceID: platformID, Access: "write", EnvironmentID: "d_platform"},
		{Subject: "user", SubjectID: currentUser, Resource: "platform_environment", ResourceID: platformID, Access: "write", EnvironmentID: environmentID},
		{Subject: "user", SubjectID: currentUser, Resource: "environment_edit_all_pipelines", ResourceID: platformID, Access: "write", EnvironmentID: environmentID},
		{Subject: "user", SubjectID: currentUser, Resource: "specific_pipeline", ResourceID: pipelineID, Access: "write", EnvironmentID: environmentID},
	}

	permOutcome, _, _, _ := permissions.MultiplePermissionChecks(perms)

	if permOutcome == "denied" {
		return "", errors.New("Requires permissions.")
	}

	// Save to code-files

	p := make([]byte, file.Size)
	file.File.Read(p)

	input := models.CodeFiles{
		EnvironmentID: environmentID,
		NodeID:        nodeID,
		FileName:      file.Filename,
		Active:        true,
		Level:         "node_file",
		FType:         "file",
		FolderID:      folderID,
	}

	// Folder excludes code directory
	parentFolder, err := filesystem.FolderConstructByID(database.DBConn, folderID)
	if err != nil {
		return "", errors.New("Create folder - build parent folder failed")
	}

	_, _, err = filesystem.CreateFile(input, parentFolder, p)
	if err != nil {
		if config.Debug == "true" {
			log.Println(err)
		}
		return "", errors.New("Failed to save file.")
	}

	return "Success", nil
}

func (r *mutationResolver) DeleteFileNode(ctx context.Context, environmentID string, fileID string, nodeID string, pipelineID string) (string, error) {
	currentUser := ctx.Value("currentUser").(string)
	platformID := ctx.Value("platformID").(string)

	// ----- Permissions
	perms := []models.Permissions{
		{Subject: "user", SubjectID: currentUser, Resource: "admin_platform", ResourceID: platformID, Access: "write", EnvironmentID: "d_platform"},
		{Subject: "user", SubjectID: currentUser, Resource: "platform_environment", ResourceID: platformID, Access: "write", EnvironmentID: environmentID},
		{Subject: "user", SubjectID: currentUser, Resource: "environment_edit_all_pipelines", ResourceID: platformID, Access: "write", EnvironmentID: environmentID},
		{Subject: "user", SubjectID: currentUser, Resource: "specific_pipeline", ResourceID: pipelineID, Access: "write", EnvironmentID: environmentID},
	}

	permOutcome, _, _, _ := permissions.MultiplePermissionChecks(perms)

	if permOutcome == "denied" {
		return "", errors.New("Requires permissions.")
	}

	// Delete file from folder
	filepath, _ := filesystem.FileConstructByID(database.DBConn, fileID)

	err := os.Remove(config.CodeDirectory + filepath)
	if err != nil {
		return "", errors.New(err.Error())
	}

	// Delete file from database
	f := models.CodeFiles{}

	err = database.DBConn.Where("file_id = ?", fileID).Delete(&f).Error

	if err != nil {
		if os.Getenv("debug") == "true" {
			logging.PrintSecretsRedact(err)
		}
		return "", errors.New("Delete file database error.")
	}

	return "Success", nil
}

func (r *mutationResolver) MoveFileNode(ctx context.Context, fileID string, toFolderID string) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CodeEditorRun(ctx context.Context, environmentID string, nodeID string, pipelineID string, path string) (string, error) {
	currentUser := ctx.Value("currentUser").(string)
	platformID := ctx.Value("platformID").(string)

	// ----- Permissions
	perms := []models.Permissions{
		{Subject: "user", SubjectID: currentUser, Resource: "admin_platform", ResourceID: platformID, Access: "write", EnvironmentID: "d_platform"},
		{Subject: "user", SubjectID: currentUser, Resource: "platform_environment", ResourceID: platformID, Access: "write", EnvironmentID: environmentID},
		{Subject: "user", SubjectID: currentUser, Resource: "environment_edit_all_pipelines", ResourceID: platformID, Access: "write", EnvironmentID: environmentID},
		{Subject: "user", SubjectID: currentUser, Resource: "specific_pipeline", ResourceID: pipelineID, Access: "write", EnvironmentID: environmentID},
	}

	permOutcome, _, _, _ := permissions.MultiplePermissionChecks(perms)

	if permOutcome == "denied" {
		return "", errors.New("Requires permissions.")
	}

	log.Println("Path: ", path)

	return "Success", nil
}

func (r *queryResolver) FilesNode(ctx context.Context, environmentID string, nodeID string, pipelineID string) (*privategraphql.CodeTree, error) {
	currentUser := ctx.Value("currentUser").(string)
	platformID := ctx.Value("platformID").(string)

	// ----- Permissions
	perms := []models.Permissions{
		{Subject: "user", SubjectID: currentUser, Resource: "admin_platform", ResourceID: platformID, Access: "write", EnvironmentID: "d_platform"},
		{Subject: "user", SubjectID: currentUser, Resource: "platform_environment", ResourceID: platformID, Access: "write", EnvironmentID: environmentID},
		{Subject: "user", SubjectID: currentUser, Resource: "environment_edit_all_pipelines", ResourceID: platformID, Access: "write", EnvironmentID: environmentID},
		{Subject: "user", SubjectID: currentUser, Resource: "specific_pipeline", ResourceID: pipelineID, Access: "write", EnvironmentID: environmentID},
	}

	permOutcome, _, _, _ := permissions.MultiplePermissionChecks(perms)

	if permOutcome == "denied" {
		return nil, errors.New("Requires permissions.")
	}

	fo := []*models.CodeFolders{}

	err := database.DBConn.Where("node_id = ?", nodeID).Find(&fo).Error
	if err != nil {
		if os.Getenv("debug") == "true" {
			logging.PrintSecretsRedact(err)
		}
		return nil, errors.New("Retrive user database error.")
	}

	fi := []*models.CodeFiles{}

	err = database.DBConn.Where("node_id = ?", nodeID).Find(&fi).Error
	if err != nil {
		if os.Getenv("debug") == "true" {
			logging.PrintSecretsRedact(err)
		}
		return nil, errors.New("Retrive user database error.")
	}

	t := privategraphql.CodeTree{
		Folders: fo,
		Files:   fi,
	}

	return &t, nil
}
