package privateresolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"dataplane/auth"
	"dataplane/database"
	"dataplane/database/models"
	privategraphql "dataplane/graphql/private"
	"dataplane/logging"
	"errors"
	"os"
	"strings"

	"github.com/google/uuid"
)

func (r *mutationResolver) UpdateDeactivateUser(ctx context.Context, userid string) (*string, error) {
	currentUserID := ctx.Value("currentUser").(string)

	if currentUserID == userid {
		return nil, errors.New("User to be deactivated cannot be the same as logged in user.")
	}
	t := models.AuthRefreshTokens{}
	// We need to remove any access to deactivated user.
	database.DBConn.Where(
		&models.AuthRefreshTokens{UserID: userid},
	).Delete(&t)

	// Check if user alredy inactive
	u := models.Users{}

	err := database.DBConn.Where("user_id = ?", userid).First(&u).Error

	if err != nil {
		if os.Getenv("debug") == "true" {
			logging.PrintSecretsRedact(err)
		}
		return nil, errors.New("Retrive me database error.")
	}

	if u.Active == false {
		return nil, errors.New("User is already inactive.")
	}

	// Deactivate user
	err = database.DBConn.Where(&models.Users{UserID: userid}).Select("status", "active").
		Updates(models.Users{Status: "inactive", Active: false}).Error

	if err != nil {
		if os.Getenv("debug") == "true" {
			logging.PrintSecretsRedact(err)
		}
		return nil, errors.New("DeactivateUser database error.")
	}

	response := "User deactivated"
	return &response, nil
}

func (r *mutationResolver) UpdateDeleteUser(ctx context.Context, userid string) (*string, error) {
	currentUserID := ctx.Value("currentUser").(string)

	if currentUserID == userid {
		return nil, errors.New("User to be deleted cannot be the same as logged in user.")
	}
	t := models.AuthRefreshTokens{}

	database.DBConn.Where(
		&models.AuthRefreshTokens{UserID: userid},
	).Delete(&t)

	u := models.Users{}

	err := database.DBConn.Where(&models.Users{UserID: userid}).Delete(&u).Error

	if err != nil {
		if os.Getenv("debug") == "true" {
			logging.PrintSecretsRedact(err)
		}
		return nil, errors.New("DeleteUser database error.")
	}

	response := "User deleted"
	return &response, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input *privategraphql.AddUsersInput) (*models.Users, error) {
	// validate if the email, username and password are in correct format
	// e := auth.ValidateRegister(input)
	// if e.Err {
	// 	finalJson, _ := json.Marshal(e)
	// 	return nil, errors.New("validation failed" + string(finalJson))
	// }

	password, err := auth.Encrypt(input.Password)

	if err != nil {
		return nil, errors.New("Password hash failed.")
	}

	userData := models.Users{
		UserID:    uuid.New().String(),
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Password:  password,
		Email:     input.Email,
		Status:    "active",
		Active:    true,
		Timezone:  input.Timezone,
		Username:  input.Email,
	}

	err = database.DBConn.Create(&userData).Error

	if err != nil {
		if os.Getenv("debug") == "true" {
			logging.PrintSecretsRedact(err)
		}
		if strings.Contains(err.Error(), "duplicate key") {
			return nil, errors.New("User already exists.")
		}
		return nil, errors.New("Register database error.")
	}

	return &models.Users{
		UserID:    userData.UserID,
		FirstName: userData.FirstName,
		LastName:  userData.LastName,
		Email:     userData.Email,
		Timezone:  userData.Timezone,
	}, nil
}

func (r *mutationResolver) UpdateChangePassword(ctx context.Context, input *privategraphql.ChangePasswordInput) (*string, error) {
	userID := ctx.Value("currentUser").(string)

	password, err := auth.Encrypt(input.Password)

	if err != nil {
		return nil, errors.New("Password hash failed.")
	}

	err = database.DBConn.Where("user_id = ?", userID).Updates(models.Users{
		Password: password,
	}).Error

	if err != nil {
		if os.Getenv("debug") == "true" {
			logging.PrintSecretsRedact(err)
		}
		return nil, errors.New("database error.")
	}

	response := "success"
	return &response, nil
}

func (r *queryResolver) LogoutUser(ctx context.Context) (*string, error) {
	userID := ctx.Value("currentUser").(string)

	u := models.AuthRefreshTokens{}

	if res := database.DBConn.Where(
		&models.AuthRefreshTokens{UserID: userID},
	).Delete(&u); res.RowsAffected <= 0 {
		return nil, errors.New("Invalid credentials")
	}

	response := "Logged out"
	return &response, nil
}
