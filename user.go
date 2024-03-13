package todo

type User struct {
	Id       int    `form:"-" db:"id"`
	Name     string `form:"name" binding:"required"`
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}
