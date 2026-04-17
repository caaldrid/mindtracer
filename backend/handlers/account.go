package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/caaldrid/mindtracer/backend/auth"
	"github.com/caaldrid/mindtracer/backend/models"
	"github.com/caaldrid/mindtracer/backend/storage"
)

type AccountHandler struct {
	users         storage.UserStorage
	secret        string
	tokenLifespan int
}

type authRegister struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email"    binding:"required"`
}

type authLogin struct {
	Email    string `json:"email"    binding:"required"`
	Password string `json:"password" binding:"required"`
}

func NewAccountHandler(users storage.UserStorage, secret string, tokenLifespan int) *AccountHandler {
	return &AccountHandler{
		users:         users,
		secret:        secret,
		tokenLifespan: tokenLifespan,
	}
}

func (a *AccountHandler) RegisterRoutes(group *gin.RouterGroup) {
	group.POST("/register", a.register)
	group.POST("/login", a.login)
}

func (a *AccountHandler) register(ctx *gin.Context) {
	var authInput authRegister

	if err := ctx.ShouldBindJSON(&authInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(authInput.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newUser := models.User{
		UserName: authInput.Username,
		Password: string(passwordHash),
		Email:    authInput.Email,
	}

	if err := a.users.Create(ctx.Request.Context(), &newUser); err != nil {
		if errors.Is(err, storage.ErrUserAlreadyExists) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, fmt.Sprintf("%s has been registered", newUser.UserName))
}

func (a *AccountHandler) login(ctx *gin.Context) {
	var authInput authLogin

	if err := ctx.ShouldBindJSON(&authInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	foundUser, err := a.users.FindByEmail(ctx.Request.Context(), authInput.Email)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email"})
		return
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(foundUser.Password),
		[]byte(authInput.Password),
	); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Password"})
		return
	}

	token, err := auth.CreateToken(foundUser.ID.String(), a.secret, a.tokenLifespan)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": fmt.Sprintf("Failed to generate token: %s", err.Error())},
		)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"accessToken": token})
}
