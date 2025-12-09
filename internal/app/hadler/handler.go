package handler

import (
	"net/http"
	"strconv"
	"time_of_armies/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
}

var currentTimeId = 1 // ПЕРЕДЕЛАТЬ НА ЗАПРОС ПО ЗАЯВКЕ ТЕКУЩЕЙ!

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) GetArmies(c *gin.Context) {
	var armies []repository.Army
	var err error
	searchArmyQuery := c.Query("searchNameArmy")
	filter := c.Query("class")

	idTime := h.Repository.GetLastTime()

	if searchArmyQuery != "" {
		armies, err = h.Repository.GetArmyByTitle(searchArmyQuery)
		if err != nil {
			logrus.Error(err)
		}
		c.HTML(http.StatusOK, "armies.html", gin.H{
			"armies":          armies,
			"armySearchQuery": searchArmyQuery,
			"lastTTid":        idTime,
		})
		return
	}

	if filter != "" {
		armies, err = h.Repository.GetArmiesByClass(filter)
		if err != nil {
			logrus.Error(err)
		}
		c.HTML(http.StatusOK, "armies.html", gin.H{
			"armies":   armies,
			"lastTTid": idTime,
		})
		return
	}

	//logrus.Info(armies)

	armies, err = h.Repository.GetArmies()
	if err != nil {
		logrus.Error(err)
	}

	var countArmiesCurTime = h.Repository.CountArmiesInTime(currentTimeId)
	c.HTML(http.StatusOK, "armies.html", gin.H{
		"armies":             armies,
		"lastTTid":           idTime,
		"countArmiesTimeBTN": countArmiesCurTime,
	})
}

func (h *Handler) GetArmy(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error(err)
	}

	oneArmy, err := h.Repository.GetArmy(id)
	if err != nil {
		logrus.Error(err)
	}

	idTime := h.Repository.GetLastTime()

	var countArmiesCurTime = h.Repository.CountArmiesInTime(currentTimeId)
	c.HTML(http.StatusOK, "onearmy.html", gin.H{
		"army":               oneArmy,
		"lastTTid":           idTime,
		"countArmiesTimeBTN": countArmiesCurTime,
	})
}

func (h *Handler) GetCurTravelTime(c *gin.Context) {
	idStr := c.Param("ttid")       // заглушка для единственной (пока что) заявки у нас
	id, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error(err)
	}

	oneTime, err := h.Repository.GetTravelTime(id)
	if err != nil {
		logrus.Error(err)
	}

	armiesRes := make([]repository.Army, 0)

	var armies []repository.Army
	armies, err = h.Repository.GetArmies()
	if err != nil {
		logrus.Error(err)
	}
	for _, idArmy := range oneTime.Armies {
		armiesRes = append(armiesRes, armies[idArmy])
	}

	var countArmiesCurTravelTime = h.Repository.CountArmiesInTime(currentTimeId)

	c.HTML(http.StatusOK, "traveltime.html", gin.H{
		"reqArmy":            armiesRes,
		"time":               oneTime,
		"countArmiesTimeBTN": countArmiesCurTravelTime,
	})
}
