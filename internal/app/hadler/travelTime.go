package handler

import (
	"math"
	"net/http"
	"strconv"
	"time"
	"time_of_armies/internal/app/ds"
	"time_of_armies/internal/app/role"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	_ "time_of_armies/docs"
)

type resultsStr struct {
	TtID int `gorm:"primaryKey;not null"`
	// 5 статусов: черновик, удален, сформирован, завершен, отклонен
	StatusTT        string
	DateCreateTT    time.Time
	DateUpdateTT    time.Time
	DateFinishTT    time.Time
	CreatorLogin    string
	ModeratorLogin  string
	ChosenBiomTT    string
	DistanceTT      int
	ResultMinTT     int
	ResultMaxTT     int
	ResultChronical int
}

func (h *Handler) ttToresultsStr(oneTime ds.TravelTime) (resultsStr, error) {
	resultTTformated := resultsStr{
		TtID:            oneTime.TtID,
		StatusTT:        oneTime.StatusTT,
		DateCreateTT:    oneTime.DateCreateTT,
		DateUpdateTT:    oneTime.DateUpdateTT,
		DateFinishTT:    oneTime.DateFinishTT,
		ChosenBiomTT:    oneTime.ChosenBiomTT,
		DistanceTT:      oneTime.DistanceTT,
		ResultMinTT:     oneTime.ResultMinTT,
		ResultMaxTT:     oneTime.ResultMaxTT,
		ResultChronical: oneTime.ResultChronical,
	}

	tmpHis, err := h.Repository.GetHistorianByID(int(oneTime.CreatorID_TT))
	if err != nil {
		return resultsStr{}, err
	}
	resultTTformated.CreatorLogin = tmpHis.HisLogin

	if int(*oneTime.ModeratorID_TT) != 0 {
		tmpHis, err = h.Repository.GetHistorianByID(int(*oneTime.ModeratorID_TT))
		if err != nil {
			return resultsStr{}, err
		}
		resultTTformated.ModeratorLogin = tmpHis.HisLogin
	} else {
		resultTTformated.ModeratorLogin = ""
	}

	return resultTTformated, nil
}

// deleteTT godoc
// @Summary      Удалить расчёт
// @Description  Удалить расчёт по его id
// @Tags         Requests
// @Produce      json
// Param [name] [type] [dataType] [required] [description]
// @Param		 ttid path int true "id расчёта"
// @Success      200  {object} message
// @Router       /travel_time/{ttid} [delete]
// @Security BearerAuth
func (h *Handler) deleteTT(c *gin.Context) {
	// "удалить" заявку=ТТ - (статус меняется на удален) с помощью выполнения SQL запроса UPDATE, без ORM
	// в третьей лабе поменять на ORM!
	idStr := c.Param("ttid")
	idTTd, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error("error deleting TT: " + err.Error())
	}
	logrus.Info("deleting tt")
	// тут надо теперь вызвать метод Модели - на замену статуса
	err = h.Repository.DeleteTT(idTTd)
	if err != nil {
		logrus.Error(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "расчёт удален!",
	})
}

func (h *Handler) getOrCreateDraftTime(creatorID int) (ds.TravelTime, error) {
	// функция для удобства - возвращает текущую заявку черновик, если же таковой нету, то создает черновик и возвращает его
	traveltime, err := h.Repository.GetTTDraft(creatorID)
	if err == nil {
		return traveltime, nil
	}
	var modID uint
	modID = 2 // ну пусть пока так...
	logrus.Info("нету черновика тут!")
	if (traveltime == ds.TravelTime{}) {
		logrus.Info("создаем черновик тут!")

		// создаем новую заявку со статусом черновик, если нету таковой в БД
		traveltime = ds.TravelTime{
			StatusTT:       "черновик",
			DateCreateTT:   time.Now(),
			CreatorID_TT:   uint(creatorID),
			ModeratorID_TT: &modID,
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

type ttToresultsStrPlusArmies struct {
	Time_travel resultsStr
	List_armies []ds.Army
}

// GetTravelTime godoc
// @Summary      Получить расчёт
// @Description  Получить расчёт по его id
// @Tags         Requests
// @Produce      json
// Param [name] [type] [dataType] [required] [description]
// @Param		 ttid path int true "id расчёта"
// @Success      200  {object} ttToresultsStrPlusArmies
// @Router       /travel_time/{ttid} [get]
// @Security BearerAuth
func (h *Handler) GetTravelTime(c *gin.Context) {

	// а эта функция дает нам любой расчёт - универсальность
	idStr := c.Param("ttid")
	idTT, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error(err)
	}

	oneTime, err := h.Repository.GetTravelTime(idTT)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "расчёт не найден!"})
		return
	}

	// if oneTime.StatusTT == "удален" {
	// 	c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "расчёт удален!"})
	// 	return
	// }

	ConnsRes, err := h.Repository.GetConnsByTTid(idTT)
	if err != nil {
		logrus.Error(err)
	}

	// чтобы вставить в шаблон в один цикл данные из двух моделей, использую доп.структуру
	type ArmyPlusIRLdistance struct {
		Army              ds.Army
		KmPerDayChronical int
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
		tmpAPIRLd.KmPerDayChronical = curCATT.KmPerDayChronical
		resArrForTT = append(resArrForTT, tmpAPIRLd)
	}

	resultTTformated, err := h.ttToresultsStr(oneTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "ошибка при формировании результата",
			"error":   err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"time_travel": resultTTformated,
		"list_armies": resArrForTT,
	})

	// var countArmiesCurTravelTime = h.Repository.CountArmiesInTime(idTT)

	// c.HTML(http.StatusOK, "traveltime.html", gin.H{
	// 	"reqArmy":            resArrForTT,
	// 	"timeTravel":         oneTime,
	// 	"countArmiesTimeBTN": countArmiesCurTravelTime,
	// })
}

type idDraftCountArmies struct {
	CountArmies int
	TTid        int
}

// GetTravelTimeDraft godoc
// @Summary		 Получить id черновика и количество армий в нём
// @Description  Получить id черновика и количество армий в нём (id вычисляется для текущего пользователя)
// @Tags         Requests
// @Produce      json
// Param [name] [type] [dataType] [required] [description]
// @Success      200  {object} idDraftCountArmies
// @Router       /travel_time [get]
// @Security BearerAuth
func (h *Handler) GetTravelTimeDraft(c *gin.Context) {
	//logrus.Info("started GetTravelTimeDraft")
	// creatorLogin := c.Query("userLogin")

	// if creatorLogin == "" {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": "Некорректно введён логин пользователя!",
	// 	})
	// 	return
	// }

	// считаю, что вход я уже сделал, и по логину вычислили ID пользователя
	//creatorID := 1

	historian, err := h.Repository.GetHistorianByLogin("Andy")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Не удалось найти пользователя, которому принадлежит заявка",
		})
	}

	traveltime, err := h.Repository.GetTTDraft(int(historian.HistorianID))
	if err != nil {
		if err.Error() == "404" {
			c.JSON(404, gin.H{
				"message": "Черновик отсутствует!",
			})
			return
		}
		//logrus.Error("ошибка получения черновика - " + err.Error())
		c.JSON(403, gin.H{
			"Error":   err,
			"message": "ошибка получения черновика!",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"TTid":        traveltime.TtID,
		"CountArmies": int(h.Repository.CountArmiesInTime(traveltime.TtID)),
	})
}

// GetTravelTimes godoc
// @Summary      Получить список Расчётов
// @Description  Получить список Расчётов, включая фильтрацию по диапазону дат и статусу
// @Tags         Requests
// @Accept 		 json
// @Produce      json
// Param [name] [type] [dataType] [required] [description]
// @Param		 dateFromTT query string false "фильтр даты начала" format(date)
// @Param		 dateToTT query string false "фильтр даты окончания" format(date)
// @Param		 statusTT query string false "статус расчёта"
// @Success      200  {array} ds.TravelTime
// @Router       /travel_times [get]
// @Security BearerAuth
func (h *Handler) GetTravelTimes(c *gin.Context) {
	//список (кроме удаленных и черновика, поля модератора и создателя через логины)
	//с фильтрацией по диапазону даты формирования и статусу

	var travel_times []ds.TravelTime

	var err error

	dateFromTT := c.Query("dateFromTT")
	dateToTT := c.Query("dateToTT")
	statusTT := c.Query("statusTT")
	//logrus.Info("status from param = ", statusTT)

	var arrStatuses []string

	if statusTT == "" {
		arrStatuses = append(arrStatuses, []string{"сформирован", "завершен", "отклонен"}...)
	} else {
		arrStatuses = append(arrStatuses, statusTT)
	}

	if dateFromTT == "" {
		dateFromTT = "1970-01-01"
	}

	if dateToTT == "" {
		dateToTT = "2038-01-18"
	}

	logrus.Error(dateFromTT)
	logrus.Error(dateToTT)

	dateF, err := time.Parse("2006-01-02", dateFromTT)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"Ошибка обработки даты": err,
		})
		return
	}

	dateT, err := time.Parse("2006-01-02", dateToTT)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"Ошибка обработки даты": err,
		})
		return
	}

	travel_times, err = h.Repository.GetTravelTimesByDatesAndStatuses(dateF, dateT, arrStatuses)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"Ошибка! не удалось получить список расчётов! ": err,
		})
		return
	}
	var results1 []resultsStr
	if ds.CurrentHistorian.Role == role.Historian {
		for _, time := range travel_times {
			logrus.Info(time)
			//var curRes resultsStr
			//logrus.Info(string(time.CreatorTT.HistorianID) + " = " + string(ds.CurrentHistorian.HistorianID))

			if time.CreatorID_TT == ds.CurrentHistorian.HistorianID {
				curRes, err := h.ttToresultsStr(time)

				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"message": "ошибка при формировании результата",
						"error":   err,
					})
					return
				}

				results1 = append(results1, curRes)
			}
		}
	} else {
		for _, time := range travel_times {
			//var curRes resultsStr

			curRes, err := h.ttToresultsStr(time)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "ошибка при формировании результата",
					"error":   err,
				})
				return
			}

			results1 = append(results1, curRes)
		}
	}

	// нужно передавтаь массив статусов

	// теперь надо из каждой отдельной заявки убрать Moderator и Creator, а так же их ID
	// нужны только логины.
	//var results []resultsStr

	c.JSON(http.StatusOK, gin.H{
		"travel_times": results1,
	})
}

// UpdateTravelTime godoc
// @Summary      Обновить расчёт
// @Description  Обновить расчёт по его id
// @Tags         Requests
// @Produce      json
// Param [name] [type] [dataType] [required] [description]
// @Param		 ttid path int true "id расчёта"
// @Param		 TimeToAdd body ds.TravelTime true "Обновленные данные расчёта"
// @Success      200  {object} ds.TravelTime
// @Router       /travel_time/{ttid} [put]
// @Security BearerAuth
func (h *Handler) UpdateTravelTime(c *gin.Context) {
	idStr := c.Param("ttid")
	idTT, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error(err)
	}

	var TimeToAdd ds.TravelTime
	err = c.ShouldBindJSON(&TimeToAdd) // считаем, что принимаем данные в формате JSON

	if err != nil {
		if err.Error() == "EOF" { // костыль на случай, если мы ничего не отправили
			return
		}
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}

	TimeTravelToChange, err := h.Repository.GetTravelTime(idTT)
	if err != nil {
		logrus.Error(err)
	}

	if TimeToAdd.StatusTT != "" {
		TimeTravelToChange.StatusTT = TimeToAdd.StatusTT
	}
	TimeTravelToChange.DateUpdateTT = time.Now()
	TimeTravelToChange.ChosenBiomTT = TimeToAdd.ChosenBiomTT
	TimeTravelToChange.DistanceTT = TimeToAdd.DistanceTT
	//TimeTravelToChange.ResultMinTT = TimeToAdd.ResultMinTT
	//TimeTravelToChange.ResultMaxTT = TimeToAdd.ResultMaxTT

	result, err := h.Repository.EditTravelTime(TimeTravelToChange)
	if err != nil {
		logrus.Error("не удалось обновить армию!")
		c.JSON(400, gin.H{
			"Error":    err.Error(),
			"messaage": "не удалось обновить армию!",
		})
		return
	}
	logrus.Info("status TT " + result.StatusTT)
	//c.JSON(http.StatusOK, result)
	resultTTformated, err := h.ttToresultsStr(result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "ошибка при формировании результата",
			"error":   err,
		})
		return
	}

	c.JSON(http.StatusOK, resultTTformated)
}

// ToFormTravelTime godoc
// @Summary      Сформировать расчёт
// @Description  Сформировать расчёт по его id, происходит проверка на обязательные поля
// @Tags         Requests
// @Produce      json
// Param [name] [type] [dataType] [required] [description]
// @Param		 ttid path int true "id расчёта"
// @Success      200  {object} ds.TravelTime
// @Router       /travel_time/{ttid}/form [put]
// @Security BearerAuth
func (h *Handler) ToFormTravelTime(c *gin.Context) {
	idStr := c.Param("ttid")
	idTT, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error(err)
	}
	h.UpdateTravelTime(c) // обновляем перед формированием поля

	logrus.Info("we are here! 1")

	// надо проверить все поля!
	// т.к. уже вызвали метод Update + уже на бэкенде и все данные уже есть в БД, то просто вызову метод модели
	//logrus.Info(idTT)
	err = h.Repository.CheckFieldsTravelTime(idTT)
	if err != nil {
		c.JSON(400, gin.H{
			"Error":    err.Error(),
			"messaage": "Некорректно заполнены поля!",
		})
		return
	}
	logrus.Info("we are here! 2")
	// по сути пользователю нужно только заполнить поля ChosenBiomTT и DistanceTT и хотя бы одну запись в М:М
	res, err := h.Repository.SetTravelTimeStatus(idTT, "сформирован")
	logrus.Info("we are here! 3")
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"Error":    err.Error(),
			"messaage": "Не удалось уставновить статус расчёта!",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"messaage":    "Расчёт сформирован!",
		"travel_time": res,
	})

}

type TTplusMessage struct {
	Travel_time ds.TravelTime
	Message     string
}

// ToModerateTravelTime godoc
// @Summary      Завершить/отклонить расчёт
// @Description  Завершить/отклонить расчёт по его id, происходит расчёт и оценка времени, требуемого на переход заданной дистанции
// @Tags         Requests
// @Produce      json
// Param [name] [type] [dataType] [required] [description]
// @Param		 ttid path int true "id расчёта"
// @Param		 statusTT formData string true "Выбор действия (завершить или отклонить)" Enums(завершить,отклонить)
// @Success      200  {object} TTplusMessage
// @Router       /travel_time/{ttid}/moderate [put]
func (h *Handler) ToModerateTravelTime(c *gin.Context) {
	//PUT завершить/отклонить модератором. При завершить/отклонении заявки
	// проставляется модератор и дата завершения. Одно из доп. полей заявки
	// или м-м рассчитывается (реализовать формулу представленную в лаб-2)
	// при завершении заявки (вычисление стоимости заказа, даты доставки
	// в течении месяца, вычисления в м-м).
	idStr := c.Param("ttid")
	idTT, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Некорректный TTid!",
		})
		return
	}
	newStatus := c.PostForm("statusTT")
	logrus.Info("status = " + newStatus)

	if newStatus != "завершен" && newStatus != "отклонен" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Некорректный статус!",
		})
		return
	}

	if newStatus == "отклонен" {
		tmp, err := h.Repository.SetTravelTimeStatus(idTT, "отклонен")
		if err != nil {
			c.JSON(500, gin.H{
				"error": err,
			})
		}
		c.JSON(200, gin.H{
			"Message":     "Расчёт отклонён!",
			"travel_time": tmp,
		})
		return
	}

	// расчёт полей
	logrus.Info("asdfaliohgpao")
	time_travel, err := h.Repository.GetTravelTime(idTT)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err,
			"message": "не удалось найти расчёт для изменения!",
		})
		return
	}

	// теперь надо рассчитать время в пути.
	// как это сделать? у меня есть несколько армий, в них уже есть времена в пути
	// также у меня уже есть бион.
	// можно при помощи switch выбрать нужные min и max, для соотв. Биома у всех армий в расчёте
	// далее найдя нужные min и max и взяв KmPerDayChronical рассчитать время в пути.

	ConnsInTT, err := h.Repository.GetConnsByTTid(idTT)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err,
			"message": "не удалось найти список армий для расчёта!",
		})
		return
	}

	// надо взять самое большое из всех армий сперва - универсальность на будущее
	// но можно просто взять очень большое число

	minSpeed := 100000000
	maxSpeed := -1 // будет сюда подставлять в зависимости от minSpeed, а не сомостоятельно
	ChronicleSpeed := 0

	switch time_travel.ChosenBiomTT {
	case "Plain":
		{
			//logrus.Info("switch plain")
			for _, conn := range ConnsInTT {
				army, err := h.Repository.GetArmy(conn.ConArmyID)
				if err != nil {
					c.JSON(500, gin.H{
						"error":   err,
						"message": "не удалось получить армии при расчёте!",
					})
					return
				}
				if minSpeed > army.MinPlainSpeed {
					minSpeed = army.MinPlainSpeed
					maxSpeed = army.MaxPlainSpeed
					ChronicleSpeed = conn.KmPerDayChronical
				}
			}
			//logrus.Info(minSpeed)
		}
	case "Desert":
		{
			logrus.Info("switch desert")

			for _, conn := range ConnsInTT {
				army, err := h.Repository.GetArmy(conn.ConArmyID)
				if err != nil {
					c.JSON(500, gin.H{
						"error":   err,
						"message": "не удалось получить армии при расчёте!",
					})
					return
				}
				if minSpeed > army.MinDesertSpeed {
					minSpeed = army.MinDesertSpeed
					maxSpeed = army.MaxDesertSpeed
					ChronicleSpeed = conn.KmPerDayChronical
				}
			}
		}
	case "River":
		{
			for _, conn := range ConnsInTT {
				army, err := h.Repository.GetArmy(conn.ConArmyID)
				if err != nil {
					c.JSON(500, gin.H{
						"error":   err,
						"message": "не удалось получить армии при расчёте!",
					})
					return
				}
				if minSpeed > army.MinRiverSpeed {
					minSpeed = army.MinRiverSpeed
					maxSpeed = army.MaxRiverSpeed
					ChronicleSpeed = conn.KmPerDayChronical
				}
			}
		}
	case "Forest":
		{
			for _, conn := range ConnsInTT {
				army, err := h.Repository.GetArmy(conn.ConArmyID)
				if err != nil {
					c.JSON(500, gin.H{
						"error":   err,
						"message": "не удалось получить армии при расчёте!",
					})
					return
				}
				if minSpeed > army.MinForestSpeed {
					minSpeed = army.MinForestSpeed
					maxSpeed = army.MaxForestSpeed
					ChronicleSpeed = conn.KmPerDayChronical
				}
			}
		}
	case "Mount":
		{
			for _, conn := range ConnsInTT {
				army, err := h.Repository.GetArmy(conn.ConArmyID)
				if err != nil {
					c.JSON(500, gin.H{
						"error":   err,
						"message": "не удалось получить армии при расчёте!",
					})
					return
				}
				if minSpeed > army.MinMountSpeed {
					minSpeed = army.MinMountSpeed
					maxSpeed = army.MaxMountSpeed
					ChronicleSpeed = conn.KmPerDayChronical
				}
			}
		}
	default:
		{

			c.JSON(400, gin.H{
				"error": "Не выбран биом!",
			})
			return
		}
	}
	// достали самую медленную армию для выбранного биома
	// теперь берем дистанию и находим кол-во дней (округляем дни в большую сторону в обоих случаях)
	logrus.Info(ChronicleSpeed)

	chronicalDays := 0

	if ChronicleSpeed != 0 {
		chronicalDays = int(math.Ceil(float64(time_travel.DistanceTT) / float64(ChronicleSpeed)))
	}

	minDays := int(math.Ceil(float64(time_travel.DistanceTT) / float64(maxSpeed)))
	maxDays := int(math.Ceil(float64(time_travel.DistanceTT) / float64(minSpeed)))

	// logrus.Info(minDays)
	// logrus.Info(maxDays)
	// logrus.Info(chronicalDays)

	time_travel.ResultMinTT = minDays
	time_travel.ResultMaxTT = maxDays
	time_travel.ResultChronical = chronicalDays

	result, err := h.Repository.EditTravelTime(time_travel)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err,
			"message": "не удалось выполнить расчёт и обновить данные!",
		})
		return
	}

	result, err = h.Repository.SetTravelTimeStatus(idTT, "завершен")

	c.JSON(http.StatusOK, gin.H{
		"message":     "Расчёт завершен!",
		"travel_time": result,
	})
}
