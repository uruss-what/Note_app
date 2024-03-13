package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	// Read the JWT token from the cookie
	//!!!!!!!
	cookies := c.Request.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "jwt" {
			// Use the token from the cookie
			token := cookie.Value
			userId, err := h.services.Authorization.ParseToken(token)
			if err != nil {
				newErrorResponse(c, http.StatusUnauthorized, err.Error())
				return
			}
			c.Set(userCtx, userId)
			return
		}
	}
	//!!!!!!!!!
	// If we reach here, the JWT token was not found in the cookie, so we fall back to the Authorization header

	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}
	userId, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(userCtx, userId)
}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return 0, errors.New("user is not found")
	}

	idInt, ok := id.(int)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id is of invalid type")
		return 0, errors.New("user id is not found")
	}
	return idInt, nil
}
