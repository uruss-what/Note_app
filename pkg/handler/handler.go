package handler

import (
	"ToDoApp/pkg/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.LoadHTMLGlob("templates/*")
	router.Static("/css", "./css")

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.GET("/sign-in", func(c *gin.Context) {
			c.HTML(http.StatusOK, "sign_in.html", gin.H{})
		})
		auth.GET("/sign-up", func(c *gin.Context) {
			c.HTML(http.StatusOK, "sign_up.html", gin.H{})
		})
	}

	api := router.Group("/api", h.userIdentity)
	{
		lists := api.Group("/lists")
		{
			lists.POST("/", h.createList)
			lists.GET("/", h.getAllLists)

			lists.GET("/:id", h.getListById)
			lists.PUT("/:id", h.updateList)
			lists.DELETE("/:id", h.deleteList)

			items := lists.Group(":id/items")
			{
				items.POST("/", h.createItem)
				items.GET("/", h.getAllItems)
			}
		}

		items := api.Group("items")
		{
			items.GET("/:id", h.getItemById)
			items.PUT("/:id", h.updateItem)
			items.DELETE("/:id", h.deleteItem)
		}
		main := api.Group("/main")
		{
			main.GET("/", h.mainPage)
		}

		users := api.Group("/users")
		{
			users.GET("/", h.usersPage)
		}
	}

	return router
}
