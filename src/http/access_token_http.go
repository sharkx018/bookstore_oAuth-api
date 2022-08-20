package http

import (
	"github.com/gin-gonic/gin"
	"github.com/sharkx018/bookstore_oauth-api/src/domain/access_token"
	access_token2 "github.com/sharkx018/bookstore_oauth-api/src/service/access_token"
	errors "github.com/sharkx018/bookstore_oauth-api/src/utils/errors"
	"net/http"
	"strings"
)

type AccessTokenHandler interface {
	GetById(ctx *gin.Context)
	Create(ctx *gin.Context)
}

type accessTokenHandler struct {
	service access_token2.Service
}

func NewHandler(service access_token2.Service) AccessTokenHandler {
	return &accessTokenHandler{
		service: service,
	}
}

func (h *accessTokenHandler) GetById(c *gin.Context) {
	accessTokenId := strings.TrimSpace(c.Param("access_token_id"))
	accessToken, err := h.service.GetById(accessTokenId)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, accessToken)
}

func (h *accessTokenHandler) Create(c *gin.Context) {
	var atr access_token.AccessTokenRequest
	//var at access_token.AccessToken
	//var err errors.RestErr
	if err := c.ShouldBindJSON(&atr); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	at, err := h.service.CreateToken(atr)

	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusCreated, at)
}
