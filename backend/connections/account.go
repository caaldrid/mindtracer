package connections

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/caaldrid/mindtracer/backend/models"
)

type accountHandler struct {
	db *gorm.DB
}

type authRegister struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email"    binding:"required"`
}

type authLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func setupAccountHandler(db *gorm.DB, router *gin.Engine) {
	account := accountHandler{
		db: db,
	}

	g := router.Group("/api/auth")
	g.POST("/register", account.register)
	g.POST("/login", account.login)
}

func (a *accountHandler) register(ctx *gin.Context) {
	var authInput authRegister

	if err := ctx.ShouldBindJSON(&authInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var foundUser models.User

	if err := a.db.Where("user_name=?", authInput.Username).Find(&foundUser).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if foundUser.ID != uuid.Nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": "User already exists for the given Username"},
		)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(authInput.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newUser := models.User{
		UserName: authInput.Username,
		Password: string(passwordHash),
		Email:    authInput.Email,
	}

	if err := a.db.Create(&newUser).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, fmt.Sprintf("%s has been registered", newUser.UserName))
}

func (a *accountHandler) login(ctx *gin.Context) {
	var authInput authLogin

	if err := ctx.ShouldBindJSON(&authInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var foundUser models.User

	if err := a.db.Where("user_name=?", authInput.Username).Find(&foundUser).Error; err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid  UserName"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(authInput.Password)); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Password"})
		return
	}

	ctx.JSON(http.StatusOK, fmt.Sprintf("%s has been logged in", foundUser.UserName))
}
