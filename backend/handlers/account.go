package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/caaldrid/mindtracer/backend/models"
	"github.com/caaldrid/mindtracer/backend/setup"
	"github.com/caaldrid/mindtracer/backend/storage"
)

type accountHandler struct {
	users         storage.UserStorage
	secret        string
	TokenLifespan int
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

func setupAccountHandler(store storage.Storage, router *gin.Engine, c setup.Config) {
	account := accountHandler{
		users:         store.Users,
		secret:        c.SecretKey,
		TokenLifespan: c.TokenLifespan,
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

	if err := a.users.CreateIfNotExists(ctx.Request.Context(), &newUser); err != nil {
		if errors.Is(err, storage.ErrUserAlreadyExists) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, fmt.Sprintf("%s has been registered", newUser.UserName))
}

func (a *accountHandler) login(ctx *gin.Context) {
	var authInput authLogin

	if err := ctx.ShouldBindJSON(&authInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	foundUser, err := a.users.FindByUsername(ctx.Request.Context(), authInput.Username)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid UserName"})
		return
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(foundUser.Password),
		[]byte(authInput.Password),
	); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Password"})
		return
	}

	token, err := createToken(foundUser.ID.String(), a.secret, a.TokenLifespan)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to generate token: %s", err.Error())})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"accessToken": token})
}
