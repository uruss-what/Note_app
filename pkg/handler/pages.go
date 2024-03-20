package handler

import (
	todo "ToDoApp"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"net/http"
)

func (h *Handler) mainPage(c *gin.Context) {
	// Получаем имя пользователя из сессии
	session, _ := store.Get(c.Request, "session-name")
	name := session.Values["username"].(string)

	tmpl, err := template.ParseFiles("templates/main.html")
	if err != nil {
		log.Println(err)
		return
	}

	data := map[string]interface{}{
		"name": name,
	}

	err = tmpl.Execute(c.Writer, data)
	if err != nil {
		log.Println(err)
		return
	}
}

func (h *Handler) listsPage(c *gin.Context) {
	// Получаем имя пользователя из сессии
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	tmpl, err := template.ParseFiles("templates/notes.html")
	if err != nil {
		log.Println(err)
		return
	}

	type ListsPageData struct {
		Lists []todo.TodoList
		Items []todo.TodoItem
	}

	lists, err := h.services.TodoList.GetAll(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var allItems []todo.TodoItem
	for _, list := range lists {
		items, err := h.services.TodoItem.GetAll(userId, list.Id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		allItems = append(allItems, items...)
	}

	data := ListsPageData{
		Lists: lists,
		Items: allItems,
	}

	err = tmpl.Execute(c.Writer, data)
	if err != nil {
		log.Println(err)
		return
	}
}

func (h *Handler) usersPage(c *gin.Context) {
	names, err := h.services.Authorization.GetAllNames()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	statuses, err := h.services.Status.GetAllStatuses()
	c.HTML(http.StatusOK, "users.html", gin.H{"names": names})
	for _, status := range statuses {
		fmt.Println("Name: ", status[0], " Status: ", status[1])
	}
}
