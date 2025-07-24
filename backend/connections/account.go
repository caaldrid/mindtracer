package connections

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/caaldrid/mindtracer/backend/models"
	"github.com/caaldrid/mindtracer/backend/setup"
)

type accountHandler struct {
	db             *gorm.DB
	secret         string
	token_lifespan int
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

type customJWTClaim struct {
	UID string `json:"uid"`
	jwt.RegisteredClaims
}

func setupAccountHandler(db *gorm.DB, router *gin.Engine, c setup.Config) {
	account := accountHandler{
		db:             db,
		secret:         c.SecretKey,
		token_lifespan: c.TokenLifespan,
	}

	g := router.Group("/api/auth")
	g.POST("/register", account.register)
	g.POST("/login", account.login)
}

func jwtAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.Request.Header.Get("Authorization")

		split := strings.Split(bearerToken, " ")
		if len(split) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(
			split[1],
			&customJWTClaim{},
			func(t *jwt.Token) (any, error) { return []byte(secret), nil },
		)

		switch {
		case token != nil && token.Valid:
			if claims, ok := token.Claims.(*customJWTClaim); ok {
				c.Request.Header.Add("uid", claims.UID)
				c.Next()
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
		case errors.Is(err, jwt.ErrTokenMalformed):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			c.Abort()
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Signature"})
			c.Abort()
		case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
			c.JSON(
				http.StatusUnauthorized,
				gin.H{"error": "Token is either expired or not yet valid"},
			)
			c.Abort()
		default:
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
		}
	}
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

	// Create claims with multiple fields populated
	claims := customJWTClaim{
		foundUser.ID.String(),
		jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(time.Hour * time.Duration(a.token_lifespan)),
			),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        uuid.NewString(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	ss, err := token.SignedString([]byte(a.secret))
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": fmt.Sprintf("Failed to generate token: %s", err.Error())},
		)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"accessToken": ss})
}
