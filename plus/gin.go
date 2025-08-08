package main

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// Controller
type UserController struct {
	session sessions.Session
}

func (c *UserController) Login(ctx *gin.Context) {
	// Handle login logic here
	c.session.Set("user", "John Doe")
	c.session.Save()
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Logged in successfully",
	})
}

func (c *UserController) Logout(ctx *gin.Context) {
	// Handle logout logic here
	c.session.Delete("user")
	c.session.Save()
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

// Middleware (Filter)
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		user := session.Get("user")
		if user == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Next()
	}
}

func main() {
	r := gin.Default()

	// Session management
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	// Controllers
	userController := &UserController{
		session: sessions.Default(r),
	}

	// Routing
	r.POST("/login", userController.Login)
	r.POST("/logout", userController.Logout)
	r.GET("/protected", AuthMiddleware(), func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Protected route",
		})
	})

	// Error handling
	r.Use(func(ctx *gin.Context) {
		ctx.Next()
		if len(ctx.Errors) > 0 {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": ctx.Errors.ByType(gin.ErrorTypePrivate).String(),
			})
		}
	})

	r.Run()
}
