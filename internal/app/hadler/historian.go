package handler

import (
	"net/http"
	"strconv"
	"time_of_armies/internal/app/ds"

	"github.com/gin-gonic/gin"
)

func (h *Handler) RegisterHistorian(c *gin.Context) {
	loginHistorian := c.PostForm("loginHistorian")
	passwordHistorian := c.PostForm("passwordHistorian")

	if loginHistorian == "" || passwordHistorian == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "логин и пароль не могут быть пустыми!",
		})
		return
	}

	historian, err := h.Repository.GetHistorianByLogin(loginHistorian)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			//"error":   err,
			"message": "Пользователь с таким именем уже существет в системе!",
		})
		return
	}

	historian.HisLogin = loginHistorian
	historian.HisPassword = passwordHistorian

	res, err := h.Repository.RegisterHistorian(historian)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err,
			"message": "Не удалось зарегистрировать пользователя!",
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message":   "Регистрация прошла успешно!",
		"historian": res,
	})

	// надо бы сделать проверку на наличие в БД
	// логин может быть единственным!
	// а вот пароль может быть любым
}

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

func (h *Handler) UpdateHistorianInfo(c *gin.Context) {
	loginHistorian := c.PostForm("loginHistorian")
	passwordHistorian := c.PostForm("passwordHistorian")
	HisIsModerator, err := strconv.ParseBool(c.PostForm("HisIsModerator"))

	// userToUpdate := ds.User_tat{}
	// err = c.ShouldBindJSON(&userToUpdate) // считаем, что принимаем данные в формате JSON
	// if err != nil {
	// 	c.JSON(400, gin.H{"Error": err.Error()})
	// }

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Нужно указать, будет ли пользователь модератором или нет!",
		})
		return
	}

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

	historianToUpdate.HisIsModerator = HisIsModerator
	historianToUpdate.HisPassword = passwordHistorian
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

func (h *Handler) AuthHistorian(c *gin.Context) {
	loginHistorian := c.PostForm("loginHistorian")
	passwordHistorian := c.PostForm("passwordHistorian")

	if loginHistorian == "" || passwordHistorian == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "логин и пароль не могут быть пустыми!",
		})
		return
	}
	historian, err := h.Repository.AuthHistorian(loginHistorian, passwordHistorian)
	if err != nil {
		c.JSON(403, gin.H{
			"error":   err,
			"message": "Логин или пароль неверны!",
		})
		return
	}
	ds.CurrentHistorian = historian
	c.JSON(http.StatusOK, gin.H{
		"message":   "Авторизация прошла успешно!",
		"historian": historian,
	})
}

func (h *Handler) DeauthHistorian(c *gin.Context) {
	if ds.CurrentHistorian.HisLogin != "" {
		ds.CurrentHistorian = ds.Historian{} // деавторизация
		c.JSON(http.StatusAccepted, gin.H{
			"message": "Вы успешно деавторизовались!",
		})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"message": "Вы уже деавторизованы!",
	})
}

// func (h *Handler) testGlob(c *gin.Context) {
// 	logrus.Info(ds.Creator.IsModerator_tat)
// }
// func (h *Handler) UpdateGlob(c *gin.Context) {
// 	ds.Creator.IsModerator_tat = !ds.Creator.IsModerator_tat
// }
