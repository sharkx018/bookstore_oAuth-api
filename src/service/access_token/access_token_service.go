package access_token

import (
	"github.com/sharkx018/bookstore_oauth-api/src/domain/access_token"
	"github.com/sharkx018/bookstore_oauth-api/src/repository/db"
	"github.com/sharkx018/bookstore_oauth-api/src/repository/rest"
	"github.com/sharkx018/bookstore_oauth-api/src/utils/errors"
)

type Service interface {
	GetById(string) (*access_token.AccessToken, *errors.RestErr)
	CreateToken(access_token.AccessTokenRequest) (*access_token.AccessToken, *errors.RestErr)
	UpdateExpirationTime(access_token.AccessToken) *errors.RestErr
}

type service struct {
	dbRepo   db.DbRepository
	restRepo rest.RestUsersRepository
}

func NewService(dbRepo db.DbRepository, restRepo rest.RestUsersRepository) Service {
	return &service{
		dbRepo:   dbRepo,
		restRepo: restRepo,
	}
}

func (s service) GetById(accessTokenId string) (*access_token.AccessToken, *errors.RestErr) {

	if len(accessTokenId) == 0 {
		return nil, errors.NewBadRequestError("invalid access token id")
	}

	accessToken, err := s.dbRepo.GetById(accessTokenId)
	if err != nil {
		return nil, err
	}

	return accessToken, nil
}

func (s service) CreateToken(request access_token.AccessTokenRequest) (*access_token.AccessToken, *errors.RestErr) {

	//Authenticate the user against the USers Api
	user, err := s.restRepo.LoginUser(request.Username, request.Password)
	if err != nil {
		return nil, err
	}

	// Generate a new access token
	at := access_token.GetNewAccessToken(user.Id)
	at.Generate()

	// Save the new access token in Cassandra
	if err := s.dbRepo.Create(at); err != nil {
		return nil, err
	}

	return &at, nil
}

func (s service) UpdateExpirationTime(at access_token.AccessToken) *errors.RestErr {
	if err := at.Validate(); err != nil {
		return err
	}
	return s.dbRepo.UpdateExpirationTime(at)
}
