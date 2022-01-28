package privateresolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	permissions "dataplane/auth_permissions"
	"dataplane/database"
	"dataplane/database/models"
	privategraphql "dataplane/graphql/private"
	"dataplane/logging"
	"dataplane/worker"
	"encoding/json"
	"errors"

	"github.com/tidwall/buntdb"
)

func (r *queryResolver) GetWorkers(ctx context.Context, environmentName string) ([]*privategraphql.Workers, error) {
	var resp []*privategraphql.Workers
	var worker worker.WorkerStats

	currentUser := ctx.Value("currentUser").(string)
	platformID := ctx.Value("platformID").(string)

	e := models.Environment{}
	database.DBConn.First(&e, "name = ?", environmentName)

	// ----- Permissions
	perms := []models.Permissions{
		{Resource: "admin_platform", ResourceID: platformID, Access: "write", Subject: "user", SubjectID: currentUser, EnvironmentID: "d_platform"},
		{Resource: "admin_environment", ResourceID: e.ID, Access: "write", Subject: "user", SubjectID: currentUser, EnvironmentID: e.ID},
		{Resource: "environment_view_workers", ResourceID: e.ID, Access: "read", Subject: "user", SubjectID: currentUser, EnvironmentID: e.ID},
	}

	permOutcome, _, _, _ := permissions.MultiplePermissionChecks(perms)

	if permOutcome == "denied" {
		return nil, errors.New("Requires permissions.")
	}

	database.GoDBWorker.View(func(tx *buntdb.Tx) error {
		tx.AscendEqual("environment", `{"Env":"`+environmentName+`"}`, func(key, val string) bool {
			// fmt.Printf("Worker Groups: %s %s\n", key, val)

			err := json.Unmarshal([]byte(val), &worker)
			if err != nil {
				logging.PrintSecretsRedact(err)
			}

			resp = append(resp, &privategraphql.Workers{
				WorkerID:    key,
				WorkerGroup: worker.WorkerGroup,
				Status:      worker.Status,
				T:           worker.T,
				Interval:    worker.Interval,
				CPUPerc:     worker.CPUPerc,
				Load:        worker.Load,
				MemoryPerc:  worker.MemoryPerc,
				MemoryUsed:  worker.MemoryUsed,
				Env:         worker.Env,
				Lb:          worker.LB,
				WorkerType:  worker.WorkerType,
			})
			return true
		})
		return nil
	})

	return resp, nil
}

func (r *queryResolver) GetWorkerGroups(ctx context.Context, environmentName string) ([]*privategraphql.WorkerGroup, error) {
	var resp []*privategraphql.WorkerGroup
	var workergroup worker.WorkerGroup

	currentUser := ctx.Value("currentUser").(string)
	platformID := ctx.Value("platformID").(string)

	e := models.Environment{}
	database.DBConn.First(&e, "name = ?", environmentName)

	// ----- Permissions
	perms := []models.Permissions{
		{Resource: "admin_platform", ResourceID: platformID, Access: "write", Subject: "user", SubjectID: currentUser, EnvironmentID: "d_platform"},
		{Resource: "admin_environment", ResourceID: e.ID, Access: "write", Subject: "user", SubjectID: currentUser, EnvironmentID: e.ID},
		{Resource: "environment_view_workers", ResourceID: e.ID, Access: "read", Subject: "user", SubjectID: currentUser, EnvironmentID: e.ID},
	}

	permOutcome, _, _, _ := permissions.MultiplePermissionChecks(perms)

	if permOutcome == "denied" {
		return nil, errors.New("Requires permissions.")
	}

	database.GoDBWorkerGroup.View(func(tx *buntdb.Tx) error {
		tx.AscendEqual("environment", `{"Env":"`+environmentName+`"}`, func(key, val string) bool {
			// fmt.Printf("Worker Groups: %s %s\n", key, val)

			err := json.Unmarshal([]byte(val), &workergroup)
			if err != nil {
				logging.PrintSecretsRedact(err)
			}

			resp = append(resp, &privategraphql.WorkerGroup{
				WorkerGroup: key,
				Status:      workergroup.Status,
				T:           workergroup.T,
				Interval:    workergroup.Interval,
				Env:         workergroup.Env,
				Lb:          workergroup.LB,
				WorkerType:  workergroup.WorkerType,
			})
			return true
		})
		return nil
	})

	return resp, nil
}
