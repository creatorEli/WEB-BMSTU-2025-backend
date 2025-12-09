package handler

import (
	"net/http"
	"strconv"
	"time"
	"time_of_armies/internal/app/ds"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) deleteTT(c *gin.Context) {
	// "удалить" заявку=ТТ - (статус меняется на удален) с помощью выполнения SQL запроса UPDATE, без ORM

	idStr := c.Param("ttid")
	idTTd, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error(err)
	}
	logrus.Info("deleting tt")
	// тут надо теперь вызвать метод Модели - на замену статуса
	err = h.Repository.DeleteTT(idTTd)
	if err != nil {
		logrus.Error(err)
	}

	// удалили заявку и переходим обратно на страницу всех услуг
	// не подгружаются услуги - ошибка
	c.Redirect(http.StatusFound, "/armies")
}

func (h *Handler) getOrCreateDraftTime() (ds.TravelTime, error) {
	// функция для удобства - возвращает текущую заявку черновик, если же таковой нету, то создает черновик и возвращает его
	creatorID := 1
	// пока что мы захардкодили id создателя заявки, в последующем сделаем авторизацию и будем получать его из JWT
	traveltime, err := h.Repository.GetTTDraft(creatorID)
	if err == nil {
		return traveltime, nil
	}
	logrus.Info("нету черновика тут!")
	if (traveltime == ds.TravelTime{}) {
		logrus.Info("создаем черновик тут!")

		// создаем новую заявку со статусом черновик, если нету таковой в БД
		traveltime = ds.TravelTime{
			StatusTT:       "черновик",
			DateCreateTT:   time.Now(),
			CreatorID_TT:   uint(creatorID),
			ModeratorID_TT: 1,        // хардкод, пока нету функционала модератора.
			ChosenBiomTT:   "Desert", // хардкод, пока нету функционала расчета.
			ResultMinTT:    8,        // хардкод, пока нету функционала расчета.
			ResultMaxTT:    12,       // хардкод, пока нету функционала расчета.
			DistanceTT:     120,      // хардкод, пока нету функционала расчета.
		}

		traveltimerr, err := h.Repository.AddTTInDB(&traveltime)
		traveltime = traveltimerr

		if err != nil {
			logrus.Panic("не удалось создать черновик расчёта перехода армии!: " + err.Error())
			return ds.TravelTime{}, err
		}
	}

	// надо возвращать не тутошний черновик, а заново получить его из бд либо достать id черновика

	newDraft, err := h.Repository.GetTTDraft(creatorID)
	if err != nil {
		logrus.Error("ошибка при возвращении новосозданного черновика расчёта")
		//return newDraft, err
	}
	return newDraft, err

}

//func (h *Handler) getDraft

func (h *Handler) GetTravelTime(c *gin.Context) {

	// как тут переделать?
	// строго говоря, id на черновик мы уже напиихали везде в ссылки
	// а эта функция дает нам любой черновик - универсальность
	idStr := c.Param("ttid")         // надо будет переделать под создание нового id
	idTT, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error(err)
	}

	oneTime, err := h.Repository.GetTravelTime(idTT)
	if err != nil {
		logrus.Error(err)
	}

	if oneTime.StatusTT == "удален" {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "расчёт удален!"})
		return
	}

	//armiesRes, err := h.Repository.GetArmiesByTTid(idTT)
	ConnsRes, err := h.Repository.GetConnsByTTid(idTT)
	if err != nil {
		logrus.Error(err)
	}

	// чтобы вставить в шаблон в один цикл данные из двух моделей, использую доп.структуру
	type ArmyPlusIRLdistance struct {
		Army        ds.Army
		KmPerDayIRL int
	}

	resArrForTT := make([]ArmyPlusIRLdistance, 0)

	for _, curCATT := range ConnsRes {
		tmpAPIRLd := ArmyPlusIRLdistance{}
		curArmy, err := h.Repository.GetArmy(curCATT.ConArmyID)
		if err != nil {
			logrus.Error("не удалось извлечь армии для расчёта")
			continue
		}
		tmpAPIRLd.Army = curArmy
		tmpAPIRLd.KmPerDayIRL = curCATT.KmPerDayIRL
		resArrForTT = append(resArrForTT, tmpAPIRLd)
	}

	var countArmiesCurTravelTime = h.Repository.CountArmiesInTime(idTT)

	c.HTML(http.StatusOK, "traveltime.html", gin.H{
		"reqArmy":            resArrForTT,
		"timeTravel":         oneTime,
		"countArmiesTimeBTN": countArmiesCurTravelTime,
	})
}
