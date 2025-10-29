package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	_ "time_of_armies/docs"
)

// DeleteConn godoc
// @Summary      Удалить армию из расчёта
// @Description  Удалить армию из расчёта по id расчёта и id армии
// @Tags         Requests
// @Produce      json
// Param [name] [type] [dataType] [required] [description]
// @Param		 ttid query int true "id расчёта"
// @Param		 ArmyID query int true "id армии"
// @Success      200  {object} message
// @Router       /travel_time/delete_army [delete]
// @Security BearerAuth
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

// UpdateConn godoc
// @Summary      Обновить летописные сведения
// @Description  Обновить летописные сведения (поле M:M) - измерение пройденного за день расстояния по летописи
// @Tags         Requests
// @Produce      json
// Param [name] [type] [dataType] [required] [description]
// @Param		 ttid query int true "id расчёта"
// @Param		 ArmyID query int true "id армии"
// @Param		 newKmPerDayChronical query int true "пройденное расстояние за день по летописи"
// @Success      200  {object} message
// @Router       /travel_time/update_army [put]
// @Security BearerAuth
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
