package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) DeleteConn(c *gin.Context) {
	//TTid int, ArmyID int

	idStrT := c.Query("TTid")
	idStrA := c.Query("ArmyID")
	idTT, err := strconv.Atoi(idStrT) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error("error deleting TT: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err,
			"message": "не удалось извлечь ID расчёта",
		})
		return
	}
	ArmyID, err := strconv.Atoi(idStrA) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error("error deleting TT: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err,
			"message": "не удалось извлечь ID армии",
		})
		return
	}
	err = h.Repository.DeleteConn(ArmyID, idTT)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err,
			"message": "не удалось удалить армию из расчёта",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Армя удалена из расчёта",
	})
}

func (h *Handler) UpdateConn(c *gin.Context) {
	idStrT := c.Query("TTid")
	idStrA := c.Query("ArmyID")
	speedStr := c.Query("newKmPerDayChronical")
	idTT, err := strconv.Atoi(idStrT) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		//logrus.Error("error deleting TT: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err,
			"message": "не удалось извлечь ID расчёта",
		})
		return
	}
	ArmyID, err := strconv.Atoi(idStrA) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		//logrus.Error("error deleting TT: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err,
			"message": "не удалось извлечь ID армии",
		})
		return
	}
	newKmPerDayChronical, err := strconv.Atoi(speedStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		//logrus.Error("error deleting TT: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err,
			"message": "не удалось извлечь новую летописную скорость",
		})
		return
	}
	err = h.Repository.UpdateConn(ArmyID, idTT, newKmPerDayChronical)
	if err != nil {
		//logrus.Error("error deleting TT: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err,
			"message": "не удалось изменить запись в БД",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "летописные сведения успешно обновлены!",
	})
}
