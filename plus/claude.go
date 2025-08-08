package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// ========== Models ==========

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var users = []User{
	{ID: 1, Name: "John Doe", Email: "john@example.com"},
	{ID: 2, Name: "Jane Smith", Email: "jane@example.com"},
}

// ========== Base Controller ==========

type BaseController struct {
	Ctx *gin.Context
}

func (c *BaseController) Init(ctx *gin.Context) {
	c.Ctx = ctx
}

func (c *BaseController) Success(data interface{}) {
	c.Ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"data":    data,
		"message": "Operation successful",
	})
}

func (c *BaseController) Error(code int, message string) {
	c.Ctx.JSON(code, gin.H{
		"status":  "error",
		"message": message,
	})
	c.Ctx.Abort()
}

func (c *BaseController) GetSession() sessions.Session {
	return sessions.Default(c.Ctx)
}

func (c *BaseController) SetSession(key string, value interface{}) error {
	session := c.GetSession()
	session.Set(key, value)
	return session.Save()
}

func (c *BaseController) GetSessionValue(key string) interface{} {
	session := c.GetSession()
	return session.Get(key)
}

func (c *BaseController) DeleteSession(key string) error {
	session := c.GetSession()
	session.Delete(key)
	return session.Save()
}

// ========== Controllers ==========

type UserController struct {
	BaseController
}

func (u *UserController) GetAll(c *gin.Context) {
	u.Init(c)

	// Check session for user authentication
	userID := u.GetSessionValue("user_id")
	if userID == nil {
		u.Error(http.StatusUnauthorized, "Please login first")
		return
	}

	u.Success(users)
}

func (u *UserController) GetByID(c *gin.Context) {
	u.Init(c)

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		u.Error(http.StatusBadRequest, "Invalid user ID")
		return
	}

	for _, user := range users {
		if user.ID == id {
			u.Success(user)
			return
		}
	}

	u.Error(http.StatusNotFound, "User not found")
}

func (u *UserController) Create(c *gin.Context) {
	u.Init(c)

	var newUser User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		u.Error(http.StatusBadRequest, "Invalid JSON data")
		return
	}

	// Generate new ID
	newUser.ID = len(users) + 1
	users = append(users, newUser)

	u.Success(newUser)
}

func (u *UserController) Update(c *gin.Context) {
	u.Init(c)

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		u.Error(http.StatusBadRequest, "Invalid user ID")
		return
	}

	var updateUser User
	if err := c.ShouldBindJSON(&updateUser); err != nil {
		u.Error(http.StatusBadRequest, "Invalid JSON data")
		return
	}

	for i, user := range users {
		if user.ID == id {
			updateUser.ID = id
			users[i] = updateUser
			u.Success(updateUser)
			return
		}
	}

	u.Error(http.StatusNotFound, "User not found")
}

func (u *UserController) Delete(c *gin.Context) {
	u.Init(c)

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		u.Error(http.StatusBadRequest, "Invalid user ID")
		return
	}

	for i, user := range users {
		if user.ID == id {
			users = append(users[:i], users[i+1:]...)
			u.Success(gin.H{"message": "User deleted successfully"})
			return
		}
	}

	u.Error(http.StatusNotFound, "User not found")
}

type AuthController struct {
	BaseController
}

func (a *AuthController) Login(c *gin.Context) {
	a.Init(c)

	var loginData struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		a.Error(http.StatusBadRequest, "Email and password are required")
		return
	}

	// Simple authentication (in real app, check against database with hashed password)
	for _, user := range users {
		if user.Email == loginData.Email && loginData.Password == "password" {
			// Set session
			if err := a.SetSession("user_id", user.ID); err != nil {
				a.Error(http.StatusInternalServerError, "Failed to create session")
				return
			}
			if err := a.SetSession("user_email", user.Email); err != nil {
				a.Error(http.StatusInternalServerError, "Failed to create session")
				return
			}

			a.Success(gin.H{
				"message": "Login successful",
				"user":    user,
			})
			return
		}
	}

	a.Error(http.StatusUnauthorized, "Invalid credentials")
}

func (a *AuthController) Logout(c *gin.Context) {
	a.Init(c)

	a.DeleteSession("user_id")
	a.DeleteSession("user_email")

	a.Success(gin.H{"message": "Logout successful"})
}

func (a *AuthController) Profile(c *gin.Context) {
	a.Init(c)

	userID := a.GetSessionValue("user_id")
	if userID == nil {
		a.Error(http.StatusUnauthorized, "Please login first")
		return
	}

	id := userID.(int)
	for _, user := range users {
		if user.ID == id {
			a.Success(user)
			return
		}
	}

	a.Error(http.StatusNotFound, "User not found")
}

// ========== Middleware/Filters ==========

// Logging middleware
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		log.Printf("[%s] %s %s %d %v",
			c.Request.Method,
			c.Request.RequestURI,
			c.ClientIP(),
			status,
			latency,
		)
	}
}

// Authentication middleware
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")

		if userID == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Authentication required",
			})
			c.Abort()
			return
		}

		// Add user info to context
		c.Set("user_id", userID)
		c.Next()
	}
}

// CORS middleware
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Rate limiting middleware
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Simple rate limiting implementation
		// In production, use Redis or similar
		c.Next()
	}
}

// ========== Error Handler ==========

func ErrorHandler() gin.HandlerFunc {
	return gin.RecoveryWithWriter(gin.DefaultWriter, func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": fmt.Sprintf("Internal server error: %s", err),
			})
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	})
}

// ========== Router Setup ==========

func SetupRouter() *gin.Engine {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// Session setup
	store := cookie.NewStore([]byte("secret-key-change-in-production"))
	store.Options(sessions.Options{
		Path:     "/",
		Domain:   "",
		MaxAge:   86400 * 7, // 7 days
		Secure:   false,     // Set to true in production with HTTPS
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	router.Use(sessions.Sessions("session", store))

	// Global middleware
	router.Use(ErrorHandler())
	router.Use(LoggerMiddleware())
	router.Use(CORSMiddleware())
	router.Use(RateLimitMiddleware())

	// Initialize controllers
	userController := &UserController{}
	authController := &AuthController{}

	// Public routes
	public := router.Group("/api/v1")
	{
		public.POST("/auth/login", authController.Login)
		public.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status": "healthy",
				"time":   time.Now(),
			})
		})
	}

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(AuthMiddleware()) // Apply auth middleware to this group
	{
		// Auth routes
		protected.POST("/auth/logout", authController.Logout)
		protected.GET("/auth/profile", authController.Profile)

		// User routes
		protected.GET("/users", userController.GetAll)
		protected.GET("/users/:id", userController.GetByID)
		protected.POST("/users", userController.Create)
		protected.PUT("/users/:id", userController.Update)
		protected.DELETE("/users/:id", userController.Delete)
	}

	// 404 handler
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Route not found",
		})
	})

	return router
}

// ========== Main Application ==========

func main() {
	router := SetupRouter()

	log.Println("Starting server on :8080")
	log.Println("API Endpoints:")
	log.Println("POST /api/v1/auth/login    - Login")
	log.Println("POST /api/v1/auth/logout   - Logout (requires auth)")
	log.Println("GET  /api/v1/auth/profile  - Get profile (requires auth)")
	log.Println("GET  /api/v1/users         - Get all users (requires auth)")
	log.Println("GET  /api/v1/users/:id     - Get user by ID (requires auth)")
	log.Println("POST /api/v1/users         - Create user (requires auth)")
	log.Println("PUT  /api/v1/users/:id     - Update user (requires auth)")
	log.Println("DELETE /api/v1/users/:id   - Delete user (requires auth)")
	log.Println("GET  /api/v1/health        - Health check")

	// Start server
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// ========== Usage Examples ==========

/*
Test the API with curl:

1. Health check:
curl http://localhost:8080/api/v1/health

2. Login:
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "john@example.com", "password": "password"}' \
  -c cookies.txt

3. Get users (using session from login):
curl -X GET http://localhost:8080/api/v1/users \
  -b cookies.txt

4. Get profile:
curl -X GET http://localhost:8080/api/v1/auth/profile \
  -b cookies.txt

5. Create user:
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "New User", "email": "new@example.com"}' \
  -b cookies.txt

6. Update user:
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "Updated Name", "email": "updated@example.com"}' \
  -b cookies.txt

7. Delete user:
curl -X DELETE http://localhost:8080/api/v1/users/1 \
  -b cookies.txt

8. Logout:
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -b cookies.txt
*/
