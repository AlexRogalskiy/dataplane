package tests

import (
	tests "dataplane/Tests"
	"dataplane/auth"
	"dataplane/routes"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/bxcodec/faker/v3"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
)

/*
go test -p 1 -v -count=1 -run TestUserAuth dataplane/Tests/users
* Create user
* Login
* Validate token
* Exchange refresh token for access token
* Logout
*/
func TestUserAuth(t *testing.T) {
	app := routes.Setup()

	graphQLUrl := "http://localhost:9000/public/graphql"

	testUser := faker.Email()
	testPassword := faker.Password()

	//--------- Create user ------------
	createUser := `mutation {
			createUser(input: {
			first_name: "` + faker.FirstName() + `",
			last_name: "` + faker.LastName() + `",
			email: "` + testUser + `",
			password: "` + testPassword + `",
			timezone: " ` + faker.Timezone() + ` ",
			}) {
			user_id
			first_name
			last_name
			email
			timezone
			}
			}`

	createUserResponse, httpResponse := tests.GraphQLRequestPublic(createUser, "{}", graphQLUrl, t, app)

	log.Println(string(createUserResponse))

	if strings.Contains(string(createUserResponse), `"errors":`) {
		t.Errorf("Error in graphql response")
	}

	assert.Equalf(t, http.StatusOK, httpResponse.StatusCode, "Create user 200 status code")

	//--------- Login ------------

	loginUser := `{
		loginUser(
		  username: "` + testUser + `",  
		  password: "` + testPassword + `",
		) {
		  access_token
		  refresh_token
		}
	  }`

	loginUserResponse, httpLoginResponse := tests.GraphQLRequestPublic(loginUser, "{}", graphQLUrl, t, app)

	log.Println(string(loginUserResponse))

	if strings.Contains(string(loginUserResponse), `"errors":`) {
		t.Errorf("Error in graphql response")
	}

	assert.Equalf(t, http.StatusOK, httpLoginResponse.StatusCode, "Login user 200 status code")

	// --------- Validate token --------------
	accessToken := jsoniter.Get(loginUserResponse, "data", "loginUser", "access_token").ToString()
	refreshToken := jsoniter.Get(loginUserResponse, "data", "loginUser", "refresh_token").ToString()
	validatetoken, _ := auth.ValidateAccessToken(accessToken)
	assert.Equalf(t, true, validatetoken, "Access token validation")

	//--------- Exchange refresh token for access token ------------
	// log.Println(refreshToken)

	reqQuery := `
	{
		"refresh_token": "` + refreshToken + `"
	}
	`

	url := "http://localhost:9000/refreshtoken"
	exchangeUserResponse, httpExchangeResponse := tests.RestRequestPublic(reqQuery, "POST", url, t, app)

	// log.Println(string(exchangeUserResponse))
	// log.Println(httpExchangeResponse)
	assert.Equalf(t, http.StatusOK, httpExchangeResponse.StatusCode, "Exchange token 200 status code")
	accessTokenExchange := jsoniter.Get(exchangeUserResponse, "access_token").ToString()

	// log.Println("Exchanged token: ", accessTokenExchange)
	validatetokenExchange, _ := auth.ValidateAccessToken(accessTokenExchange)
	assert.Equalf(t, true, validatetokenExchange, "Exchange access token validation")

	//--------- Logout ------------
	logoutUser := `{
		logoutUser
	  }`
	graphQLPrivateUrl := "http://localhost:9000/private/graphql"
	logoutUserResponse, httpLogoutResponse := tests.GraphQLRequestPrivate(logoutUser, accessTokenExchange, "{}", graphQLPrivateUrl, t, app)
	assert.Equalf(t, http.StatusOK, httpLogoutResponse.StatusCode, "Logout 200 status code")
	assert.Equalf(t, "Logged out", jsoniter.Get(logoutUserResponse, "data", "logoutUser").ToString(), "Logout correct response")

}
