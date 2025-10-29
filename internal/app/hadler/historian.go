package handler

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"
	"time_of_armies/internal/app/ds"
	"time_of_armies/internal/app/role"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	_ "time_of_armies/docs"
)

type HisLPSw struct {
	LoginHistorian    string
	PasswordHistorian string
}

type MesHisLPMSw struct {
	Message   string
	Historian ds.Historian
}

type registerReq struct {
	Name string `json:"name"` // лучше назвать то же самое что login
	Pass string `json:"pass"`
}

type registerResp struct {
	Ok bool `json:"ok"`
}

// RegisterHistorian godoc
// @Summary      Регистрация историка
// @Description  Регистрация историка по его логину и паролю
// @Tags         Requests
// @Produce      json
// Param [name] [type] [dataType] [required] [description]
// @Param		 LoginHistorian formData string true "логин историка"
// @Param		 PasswordHistorian formData string true "пароль историка"
// @Success      200  {object} MesHisLPMSw
// @Router       /historian/reg [post]
func (h *Handler) RegisterHistorian(c *gin.Context) {
	loginHistorian := c.PostForm("loginHistorian")
	passwordHistorian := c.PostForm("passwordHistorian")

	if loginHistorian == "" || passwordHistorian == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "логин и пароль не могут быть пустыми!",
		})
		return
	}

	_, err := h.Repository.GetHistorianByLogin(loginHistorian)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			//"error":   err,
			"message": "Пользователь с таким именем уже существет в системе!",
		})
		return
	}

	_, err = h.Repository.RegisterHistorian(ds.Historian{
		HistorianUUID: uuid.New(),
		HisLogin:      loginHistorian,
		HisPassword:   generateHashString(passwordHistorian),
		Role:          role.Historian,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err,
			"message": "Не удалось зарегистрировать пользователя!",
		})
		return
	}

	c.JSON(http.StatusAccepted, &registerResp{
		Ok: true,
	})
}

func generateHashString(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// type HistILM struct {
// 	Id           int
// 	Login        string
// 	Is_moderator bool
// }

type HisSw struct {
	Historian ds.Historian
}

// GetHistorianInfo godoc
// @Summary      Получить информацию об историке
// @Description  Получить информацию об историке (для его личного кабинета), уже авторизованном
// @Tags         Requests
// @Produce      json
// Param [name] [type] [dataType] [required] [description]
// @Success      200  {object} HisSw
// @Router       /historian [get]
// @Security BearerAuth
func (h *Handler) GetHistorianInfo(c *gin.Context) {
	if ds.CurrentHistorian.HisLogin == "" {
		c.JSON(403, gin.H{
			"message": "В данный момент в системе нету автороизованного пользователя!",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"historian": ds.CurrentHistorian,
	})
}

// UpdateHistorianInfo godoc
// @Summary      Обновить информацию об историке
// @Description  Обновить информацию об историке по его логину
// @Tags         Requests
// @Produce      json
// Param [name] [type] [dataType] [required] [description]
// @Param		 LoginHistorian formData string true "логин историка"
// @Param		 PasswordHistorian formData string true "пароль историка"
// @Param		 HisIsModerator formData bool true "Историк - модератор?"
// @Success      200  {object} HisSw
// @Router       /historian [put]
// @Security BearerAuth
func (h *Handler) UpdateHistorianInfo(c *gin.Context) {
	loginHistorian := c.PostForm("loginHistorian")
	passwordHistorian := c.PostForm("passwordHistorian")
	//HisIsModerator, err := strconv.ParseBool(c.PostForm("HisIsModerator"))

	// userToUpdate := ds.User_tat{}
	// err = c.ShouldBindJSON(&userToUpdate) // считаем, что принимаем данные в формате JSON
	// if err != nil {
	// 	c.JSON(400, gin.H{"Error": err.Error()})
	// }

	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"message": "Нужно указать, будет ли пользователь модератором или нет!",
	// 	})
	// 	return
	// }

	if loginHistorian == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Нужно указать логин пользователя, чьи данные надо обновить!",
		})
		return
	}

	historianToUpdate, err := h.Repository.GetHistorianByLogin(loginHistorian)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err,
			"message": "Не найден пользователь, чьи данные нужно обновить!",
		})
		return
	}

	//historianToUpdate.HisIsModerator = HisIsModerator
	historianToUpdate.HisPassword = generateHashString(passwordHistorian)
	historianToUpdate.HisLogin = loginHistorian

	res, err := h.Repository.UpdateHistorian(historianToUpdate)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err,
			"message": "не удалось обновить данные пользователя!",
		})
	}
	c.JSON(http.StatusAccepted, gin.H{
		"message":   "Данные успешно обновлены!",
		"historian": res,
	})
	ds.CurrentHistorian = res
}

type loginReq struct {
	Login    string `json:"loginHistorian"`
	Password string `json:"passwordHistorian"`
}

type loginResp struct {
	ExpiresIn   int    `json:"expires_in"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

// AuthHistorian godoc
// @Summary      Аутентификация историка
// @Description  Аутентификация историка по его логину и паролю
// @Tags         Requests
// @Produce      json
// Param [name] [type] [dataType] [required] [description]
// @Param		 loginHistorian formData string true "логин историка"
// @Param		 passwordHistorian formData string true "пароль историка"
// @Success      200  {object} MesHisLPMSw
// @Router       /historian/auth [post]
func (h *Handler) AuthHistorian(c *gin.Context) {
	loginHistorian := c.PostForm("loginHistorian")
	passwordHistorian := c.PostForm("passwordHistorian")

	logrus.Info("authHistorian !")

	if loginHistorian == "" || passwordHistorian == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "логин и пароль не могут быть пустыми!",
		})
		return
	}

	historian, err := h.Repository.AuthHistorian(loginHistorian, generateHashString(passwordHistorian))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			//"error":   err,
			"message": "Логин или пароль неверны!",
		})
		return
	}
	ds.CurrentHistorian = historian

	// вот тут проверили, что успешно авторизовались, можно и токен делать
	//req := &loginReq{}
	cfg := h.config
	logrus.Info("cfg.JWT.ExpiresIn = " + string(cfg.JWT.ExpiresIn))
	claims := &ds.JWTClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(cfg.JWT.ExpiresIn).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "time_of_armies",
			Subject:   historian.HisLogin,
			//Id:        uuid.New().String(), // Уникальный ID токена
		},
		UserUUID: historian.HistorianUUID,
		Scopes:   []string{"historian"},
		Login:    historian.HisLogin,
		Role:     historian.Role,
	}

	token := jwt.NewWithClaims(cfg.JWT.GetSigningMethod(), claims)

	if token == nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Токен не был создан!"))
		return
	}

	strToken, err := token.SignedString([]byte(cfg.JWT.Token))
	if err != nil {
		logrus.Error("Token signing error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Ошибка создания токена",
		})
		return
	}

	// Преобразуем ExpiresIn в секунды (стандарт для OAuth)
	expiresInSec := int(cfg.JWT.ExpiresIn / time.Second)

	logrus.Info("Successful auth for: ", loginHistorian)

	c.JSON(http.StatusOK, loginResp{
		ExpiresIn:   expiresInSec,
		AccessToken: strToken,
		TokenType:   "Bearer",
	})

	logrus.Info()
	// c.JSON(http.StatusOK, gin.H{
	// 	"message":   "Авторизация прошла успешно!",
	// 	"historian": historian,
	// })
}

// UpdateHistorianInfo godoc
// @Summary      Деавторизация историка
// @Description  Деавторизация историка (завершение текущей сессии)
// @Tags         Requests
// @Produce      json
// Param [name] [type] [dataType] [required] [description]
// @Success      200  {object} message
// @Router       /historian/exit [post]
// @Security BearerAuth
func (h *Handler) DeauthHistorian(c *gin.Context) {
	// if ds.CurrentHistorian.HisLogin != "" {
	// 	c.JSON(http.StatusAccepted, gin.H{
	// 		"message": "Вы успешно деавторизовались!",
	// 	})
	// 	return
	// }
	// c.JSON(http.StatusBadRequest, gin.H{
	// 	"message": "Вы уже деавторизованы!",
	// })

	jwtStr := c.GetHeader("Authorization")
	if !strings.HasPrefix(jwtStr, jwtPrefix) { // если нет префикса то нас дурят!
		c.AbortWithStatus(http.StatusBadRequest) // отдаем что нет доступа
		return                                   // завершаем обработку
	}

	jwtStr = jwtStr[len("Bearer "):]

	_, err := jwt.ParseWithClaims(jwtStr, &ds.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.config.JWT.Token), nil
	})

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		logrus.Info(err)
		return
	}

	err = h.redis.WriteJWTToBlacklist(c.Request.Context(), jwtStr, h.config.JWT.ExpiresIn)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ds.CurrentHistorian = ds.Historian{} // деавторизация
	c.JSON(http.StatusOK, gin.H{
		"message": "Вы успешно деавторизовались",
	})
}

// func (h *Handler) testGlob(c *gin.Context) {
// 	logrus.Info(ds.Creator.IsModerator_tat)
// }
// func (h *Handler) UpdateGlob(c *gin.Context) {
// 	ds.Creator.IsModerator_tat = !ds.Creator.IsModerator_tat
// }
