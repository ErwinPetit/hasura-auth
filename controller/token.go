package controller

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

func (ctrl *Controller) RefreshToken(ctx *gin.Context) {
	var req RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ctrl.hasuraConn.GetUserByRefreshToken(ctx.Request.Context(), req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 1 in 10 request will delete expired refresh tokens
	// TODO: CRONJOB in the future.
	if rand.Intn(10) == 0 {
		_ = ctrl.hasuraConn.DeleteExpiredRefreshToken(ctx.Request.Context())
	}

	if err := ctrl.hasuraConn.UpdateUserLastSeen(ctx.Request.Context(), user.ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": user.ID,
			"iat": now.UnixMilli(),
			"exp": now.Add(15 * time.Minute).UnixMilli(),
			"iss": "hasura-auth",
			"https://hasura.io/jwt/claims": map[string]any{
				"x-hasura-allowed-roles":     user.Roles,
				"x-hasura-default-role":      user.DefaultRole,
				"x-hasura-user-id":           user.ID,
				"x-hasura-user-is-anonymous": fmt.Sprintf("%v", user.IsAnonymous),
			},
		},
	)

	tokenString, err := token.SignedString([]byte(ctrl.jwtSecret.Key))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"accessToken":          tokenString,
			"accessTokenExpiresIn": 900,
			"refreshToken":         req.RefreshToken,
			"user":                 user,
		},
	)
}
