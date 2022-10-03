package rest

import (
	"encoding/json"
	"errors"
	"github.com/mercadolibre/golang-restclient/rest"
	"github.com/sharkx018/bookstore_oauth-api/src/domain/users"
	"github.com/sharkx018/bookstore_utils-go/rest_errors"
	"time"
)

var (
	userRestClient = rest.RequestBuilder{
		BaseURL: "http://localhost:8081",
		Timeout: 100 * time.Millisecond,
	}
)

type RestUsersRepository interface {
	LoginUser(email string, password string) (*users.User, rest_errors.RestErr)
}

type userRepository struct{}

func NewRepository() RestUsersRepository {
	return &userRepository{}
}

func (u userRepository) LoginUser(email string, password string) (*users.User, rest_errors.RestErr) {

	request := users.UserLoginRequest{
		Email:    email,
		Password: password,
	}

	response := userRestClient.Post("/users/login", request)
	if response == nil || response.Response == nil {
		return nil, rest_errors.NewInternalServerError("invalid rest client response when trying to login user", errors.New("restclient error"))
	}

	if response.StatusCode > 299 {
		apiError, err := rest_errors.NewRestErrorFromBytes(response.Bytes())
		if err != nil {
			return nil, rest_errors.NewInternalServerError("invalid error interface when trying to login user", err)
		}
		return nil, apiError
	}

	var user users.User
	if err := json.Unmarshal(response.Bytes(), &user); err != nil {
		return nil, rest_errors.NewInternalServerError("error when trying to unmarshall users response", errors.New("restclient error"))
	}

	return &user, nil
}
