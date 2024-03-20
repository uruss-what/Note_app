package handler

import (
	todo "ToDoApp"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"net/http"
)

var store = sessions.NewCookieStore([]byte("your-secret-key"))

func (h *Handler) signUp(c *gin.Context) {
	var input todo.User

	//метод контекста чтоб парсить значения из json в поля с аналогичными тегами
	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.Redirect(http.StatusSeeOther, "/api/sign-in")
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})

}

type signInInput struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input signInInput

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	// Set the JWT token in a cookie
	c.SetCookie("jwt", token, 3600, "/", "localhost", false, true)

	// Извлекаем пользователя из базы данных
	name, err := h.services.Authorization.GetName(input.Username, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Сохраняем имя пользователя в сессии
	session, _ := store.Get(c.Request, "session-name")
	session.Values["username"] = name
	session.Save(c.Request, c.Writer)

	// Redirect the user to the main page
	c.Redirect(http.StatusSeeOther, "/api/main")

}
