package db

import (
	"errors"
	"github.com/gocql/gocql"
	"github.com/sharkx018/bookstore_oauth-api/clients/cassandra"
	"github.com/sharkx018/bookstore_oauth-api/src/domain/access_token"
	"github.com/sharkx018/bookstore_utils-go/rest_errors"
)

const (
	queryGetAccessToken    = "SELECT access_token, user_id, client_id, expires FROM access_tokens WHERE access_token=?;"
	queryCreateAccessToken = "INSERT INTO access_tokens(access_token, user_id, client_id, expires) VALUES(?,?,?,?);"
	queryUpdateExpires     = "UPDATE access_tokens SET expires=? WHERE access_token=?;"
)

func NewRepository() DbRepository {
	return &dbRepository{}
}

type dbRepository struct{}

type DbRepository interface {
	GetById(string) (*access_token.AccessToken, rest_errors.RestErr)
	Create(at access_token.AccessToken) rest_errors.RestErr
	UpdateExpirationTime(at access_token.AccessToken) rest_errors.RestErr
}

func (r *dbRepository) GetById(id string) (*access_token.AccessToken, rest_errors.RestErr) {

	var result access_token.AccessToken
	if err := cassandra.GetSession().Query(queryGetAccessToken, id).Scan(
		&result.AccessToken,
		&result.UserId,
		&result.ClientID,
		&result.Expires); err != nil {
		if err == gocql.ErrNotFound {
			return nil, rest_errors.NewNotFoundError("no access token found with given id")
		}

		return nil, rest_errors.NewInternalServerError("error when trying to get current id:", errors.New("database error"))
	}

	return &result, nil
}

func (r *dbRepository) Create(at access_token.AccessToken) rest_errors.RestErr {

	if err := cassandra.GetSession().Query(queryCreateAccessToken,
		at.AccessToken,
		at.UserId,
		at.ClientID,
		at.Expires).Exec(); err != nil {
		return rest_errors.NewInternalServerError("error when trying to save access token in database", err)
	}

	return nil

}

func (r *dbRepository) UpdateExpirationTime(at access_token.AccessToken) rest_errors.RestErr {

	if err := cassandra.GetSession().Query(queryUpdateExpires,
		at.Expires,
		at.AccessToken,
	).Exec(); err != nil {
		return rest_errors.NewInternalServerError("error when trying ()to update currnet usererr.Error()", errors.New("database error"))
	}

	return nil

}
