package rest

import (
	"github.com/mercadolibre/golang-restclient/rest"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	rest.StartMockupServer()
	os.Exit(m.Run())
}

func TestLoginUserTimeoutFromApi(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "http://localhost:8080/users/login",
		HTTPMethod:   http.MethodPost,
		ReqBody:      `{"email":"email@gmail.com","password":"abc123"}`,
		RespHTTPCode: -1,
		RespBody:     `{}`,
	})

	repository := userRepository{}
	user, err := repository.LoginUser("email@gmail.com", "abc123")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status)
	assert.Equal(t, "invalid rest client response when trying to login user", err.Message)

}

func TestLoginUserInvalidErrorInterface(t *testing.T) {

	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "http://localhost:8080/users/login",
		HTTPMethod:   http.MethodPost,
		ReqBody:      `{"email":"email@gmail.com","password":"abc123"}`,
		RespHTTPCode: 500,
		RespBody:     `{"message":"not found","status":"500", "error":"internal_server_error"}`,
	})

	repository := userRepository{}
	user, err := repository.LoginUser("email@gmail.com", "abc123")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status)
	assert.Equal(t, "invalid err interface when trying to login user", err.Message)
}

func TestLoginUserInvalidLoginCreds(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "http://localhost:8080/users/login",
		HTTPMethod:   http.MethodPost,
		ReqBody:      `{"email":"email@gmail.com","password":"abc123"}`,
		RespHTTPCode: 401,
		RespBody:     `{"message":"invalid user credentials","status":401, "error":"unauthorized"}`,
	})

	repository := userRepository{}
	user, err := repository.LoginUser("email@gmail.com", "abc123")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusUnauthorized, err.Status)
	assert.Equal(t, "invalid user credentials", err.Message)
}

func TestLoginUserInvalidUserJsonResponse(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "http://localhost:8080/users/login",
		HTTPMethod:   http.MethodPost,
		ReqBody:      `{"email":"email@gmail.com","password":"abc123"}`,
		RespHTTPCode: 200,
		RespBody:     `{"id":"1","first_name":"mukul", "last_name":"verma","email":"mukul@email.com","status":"active"}`,
	})

	repository := userRepository{}
	user, err := repository.LoginUser("email@gmail.com", "abc123")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status)
	assert.Equal(t, "error when trying to unmarshall users response", err.Message)
}

func TestLoginUserNoError(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "http://localhost:8080/users/login",
		HTTPMethod:   http.MethodPost,
		ReqBody:      `{"email":"email@gmail.com","password":"abc123"}`,
		RespHTTPCode: 200,
		RespBody:     `{"id":1,"first_name":"mukul", "last_name":"verma","email":"mukul@email.com","status":"active"}`,
	})

	repository := userRepository{}
	user, err := repository.LoginUser("email@gmail.com", "abc123")

	assert.NotNil(t, user)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), user.Id)
	assert.Equal(t, "mukul", user.FirstName)
	assert.Equal(t, "verma", user.LastName)
	assert.Equal(t, "mukul@email.com", user.Email)
	assert.Equal(t, "active", user.Status)
}
