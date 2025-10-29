package handler

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time_of_armies/internal/app/ds"
	"time_of_armies/internal/app/role"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
)

const jwtPrefix = "Bearer "

func (h *Handler) WithAuthCheck(assignRoles ...role.Role) func(c *gin.Context) {
	return func(c *gin.Context) {
		logrus.Info("got in WithAuthCheck")
		jwtStr := c.GetHeader("Authorization")
		if !strings.HasPrefix(jwtStr, jwtPrefix) { // если нет префикса то нас дурят!
			c.AbortWithStatus(http.StatusForbidden) // отдаем что нет доступа
			return                                  // завершаем обработку
		}

		// отрезаем префикс
		jwtStr = jwtStr[len(jwtPrefix):]

		// ПРОВЕРЯЕМ ЧЕРНЫЙ СПИСОК REDIS ПЕРВЫМ ДЕЛОМ
		err := h.redis.CheckJWTInBlacklist(c.Request.Context(), jwtStr)
		if err == nil {
			// Токен найден в черном списке
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Токен более не действителен",
			})
			return
		}

		// Если ошибка не nil и не "ключ не найден" - это реальная ошибка Redis
		if !errors.Is(err, redis.Nil) {
			logrus.Errorf("Redis error in CheckJWTInBlacklist: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Ошибка сервера - нет доступа к redis!",
			})
			return
		}

		token, err := jwt.ParseWithClaims(jwtStr, &ds.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(h.config.JWT.Token), nil
		})

		if err != nil {
			c.JSON(403, gin.H{
				"error":   err,
				"token":   token,
				"message": "Доступ запрещён!",
			})
			logrus.Error(err)
			return
		}

		// Проверяем валидность токена
		if !token.Valid {
			c.AbortWithStatusJSON(403, gin.H{
				"error":   "Invalid token",
				"message": "Доступ запрещён!",
			})
			return
		}

		if claims, ok := token.Claims.(*ds.JWTClaims); ok {
			c.Set("userID", claims.UserUUID)
			c.Set("userLogin", claims.Login)
			c.Set("userScopes", claims.Scopes)
		}

		myClaims := token.Claims.(*ds.JWTClaims)

		roleFound := false
		for _, oneOfAssignedRole := range assignRoles {
			if myClaims.Role == oneOfAssignedRole {
				roleFound = true
				break
			}
		}

		if !roleFound {
			c.AbortWithStatus(http.StatusForbidden)
			log.Printf("role %v is not assigned in %v", myClaims.Role, assignRoles)
			return
		}
		logrus.Info(h.config.JWT.ExpiresIn)
		// Если роль найдена в разрешенных - продолжаем выполнение
		c.Next()
	}
}
