package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time_of_armies/internal/app/ds"

	_ "time_of_armies/docs"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

//var creatorID = 1 // хардкод пока нету функцонала юзера
// переделать на singleton юзера!

// GetArmies godoc
// @Summary      Получить список армий
// @Description  Получить список армий, включая фильтрацию (param "class") и поисковый запрос (param "searchNameArmy")
// @Tags         Requests
// @Accept 		 json
// @Produce      json
// Param [name] [type] [dataType] [required] [description]
// @Param		 class query string false "фильтрация"
// @Param 		 searchNameArmy query string false "поиск армии"
// @Success      200  {array} ds.Army
// @Router       /armies [get]
func (h *Handler) GetArmies(c *gin.Context) {

	logrus.Info("GetArmies!")

	// если есть расчёт черновик у юзера, то добавляем ссылку, иначе впихиваем якорь и нуль число

	// var href string = "/armies"
	// var armiesInDraft int64 = 0
	// ttDraft, errorr := h.Repository.GetTTDraft(creatorID)
	// if errorr == nil {
	// 	armiesInDraft = h.Repository.CountArmiesInTime(ttDraft.TtID)
	// 	href = "/travel_time/" + strconv.Itoa(ttDraft.TtID)
	// }

	var armies []ds.Army

	var err error
	searchArmyQuery := c.Query("searchNameArmy")
	filter := c.Query("class")

	if searchArmyQuery != "" {
		armies, err = h.Repository.GetArmyByTitle(searchArmyQuery)
		if err != nil {
			logrus.Error(err)
		}
		// c.HTML(http.StatusOK, "armies.html", gin.H{
		// 	"armies":             armies,
		// 	"armySearchQuery":    searchArmyQuery,
		// 	"countArmiesTimeBTN": armiesInDraft,
		// 	"hrefToTT":           href,
		// })

		armiesJSON, err := json.Marshal(armies)
		if err != nil {
			logrus.Error("не удалось перевести массив структур в джейсон форман")
		}
		c.JSON(http.StatusOK, gin.H{
			"armies": armiesJSON,
			//"armySearchQuery": searchArmyQuery,
		})
		return
	}

	if filter != "" {
		armies, err = h.Repository.GetArmiesByClass(filter)
		if err != nil {
			logrus.Error(err)
		}
		armiesJSON, err := json.Marshal(armies)
		if err != nil {
			logrus.Error("не удалось перевести массив структур в джейсон форман")
		}
		c.JSON(http.StatusOK, gin.H{
			"armies": armiesJSON,
			//"armySearchQuery": searchArmyQuery,
		})
		// c.HTML(http.StatusOK, "armies.html", gin.H{
		// 	"armies":             armies,
		// 	"countArmiesTimeBTN": armiesInDraft,
		// 	"hrefToTT":           href,
		// })
		return
	}
	armies, err = h.Repository.GetArmies()
	if err != nil {
		logrus.Error(err)
	}
	// armiesJSON, err := json.Marshal(armies)
	// if err != nil {
	// 	logrus.Error("не удалось перевести массив структур в джейсон форман")
	// }
	// var narm []ds.Army
	// re := json.Unmarshal(armiesJSON, &narm)
	// if re != nil {
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{
		"armies": armies,
		//"armySearchQuery": searchArmyQuery,
	})
	// c.HTML(http.StatusOK, "armies.html", gin.H{
	// 	"armies":             armies,
	// 	"countArmiesTimeBTN": armiesInDraft,
	// 	"hrefToTT":           href,
	// })
}

// GetArmy godoc
// @Summary      Получить одну армию
// @Description  Получить одну армию по её id
// @Tags         Requests
// @Produce      json
// Param [name] [type] [dataType] [required] [description]
// @Param		 id path int true "id армии"
// @Success      200  {object} ds.Army
// @Router       /army/{id} [get]
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

	c.JSON(http.StatusOK, oneArmy)
}

type ArmyResSw struct {
	Army ds.Army
}

// AddArmy godoc
// @Summary      Создать армию
// @Description  Создать одну армию
// @Tags         Requests
// @Accept 		 json
// @Produce      json
// Param [name] [type] [dataType] [required] [description]
// @Param		 army body ds.Army true "данные армии"
// @Success      200  {object} ArmyResSw
// @Router       /army [post]
func (h *Handler) AddArmy(c *gin.Context) {
	var armyToAdd ds.Army

	err := c.ShouldBindJSON(&armyToAdd) // считаем, что принимаем данные в формате JSON
	if err != nil {
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}

	resArmy, err := h.Repository.AddArmy(armyToAdd)
	if err != nil {
		logrus.Error("Не удалось добавить армию в БД!")
	}
	c.JSON(http.StatusOK, gin.H{
		"army": resArmy,
	})
}

// UpdateArmy godoc
// @Summary      Изменить армию
// @Description  Изменить одну армию, зная её идентификатор и изменяя существующие поля
// @Tags         Requests
// @Accept 		 json
// @Produce      json
// Param [name] [type] [dataType] [required] [description]
// @Param		 id path int true "ID изменяемой армии"
// @Param		 army body ds.Army true "данные армии"
// @Success      200  {object} ArmyResSw
// @Router       /army/{id} [put]
func (h *Handler) UpdateArmy(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error(err)
	}

	armyToChange, err := h.Repository.GetArmy(id)
	if err != nil {
		logrus.Error(err)
	}
	logrus.Info(armyToChange)

	var armyToAdd ds.Army
	err = c.ShouldBindJSON(&armyToAdd) // считаем, что принимаем данные в формате JSON
	if err != nil {
		c.JSON(400, gin.H{"Error": err.Error()})
	}

	armyToChange.NameArmy = armyToAdd.NameArmy
	armyToChange.DescriptionArmy = armyToAdd.DescriptionArmy
	armyToChange.StatusArmy = armyToAdd.StatusArmy
	//armyToChange.ImageArmyUrl = armyToAdd.ImageArmyUrl //тут изображение не меняем
	armyToChange.ClassArmy = armyToAdd.ClassArmy
	armyToChange.MinPlainSpeed = armyToAdd.MinPlainSpeed
	armyToChange.MaxPlainSpeed = armyToAdd.MaxPlainSpeed
	armyToChange.MinMountSpeed = armyToAdd.MinMountSpeed
	armyToChange.MaxMountSpeed = armyToAdd.MaxMountSpeed
	armyToChange.MinForestSpeed = armyToAdd.MinForestSpeed
	armyToChange.MaxForestSpeed = armyToAdd.MaxForestSpeed
	armyToChange.MinRiverSpeed = armyToAdd.MinRiverSpeed
	armyToChange.MaxRiverSpeed = armyToAdd.MaxRiverSpeed
	armyToChange.MinDesertSpeed = armyToAdd.MinDesertSpeed
	armyToChange.MaxDesertSpeed = armyToAdd.MaxDesertSpeed

	resArmy, err := h.Repository.EditArmy(armyToChange)
	if err != nil {
		logrus.Error("не удалось обновить армию!")
		c.JSON(500, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"army": resArmy,
	})
}

type message struct {
	Message string
}

// DeleteArmy godoc
// @Summary      Удалить армию
// @Description  Удалить одну армию, зная её идентификатор и изменяя статус (запрещено удалять запись из БД)
// @Tags         Requests
// @Produce      json
// Param [name] [type] [dataType] [required] [description]
// @Param		 id path int true "ID удаляемой армии"
// @Success      200  {object} message
// @Router       /army/{id} [delete]
func (h *Handler) DeleteArmy(c *gin.Context) {
	// Антон Игоревич говорит, что удалять можно только записи из М:М, посему меняю только статус
	// удален
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		c.JSON(500, gin.H{
			"error":   err,
			"message": "Некорректно введён ID !",
		})
		return
	}

	err = h.Repository.DeleteArmy(id)

	if err != nil {
		//logrus.Error("Не удалось удалить армию!")
		c.JSON(500, gin.H{
			"error":   err,
			"message": "Не удалось удалить армию!",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Армия 'удалена'!",
	})
}

type ResDraftSw struct {
	Message     string
	Travel_time ds.TravelTime
}

// addArmyToTT godoc
// @Summary		 Добавить армию в расчёт
// @Description  Добавить армию в расчёт, зная её идентификатор и при необходимости создавая новый черновик
// @Tags         Requests
// @Produce      json
// Param [name] [type] [dataType] [required] [description]
// @Param		 id formData int true "ID добавляемой армии"
// @Success      200  {object} ResDraftSw
// @Router       /army/add_to_travel [post]
// @Security BearerAuth
func (h *Handler) addArmyToTT(c *gin.Context) {
	// получаем данные в виде form-data
	//logrus.Info("adding post army 0")
	idStr := c.PostForm("ArmyID")
	//logrus.Info(idStr)
	idArmy, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		c.JSON(400, gin.H{
			"error":   err,
			"message": "Некорректно введён ID !",
		})
		return
	}

	//logrus.Info("adding post army 1")

	// ПОКА ХАРДКОД 1, Т.К. НЕТУ ФУНКЦИОНАЛА АУТЕНТИФИКАЦИИ !!!

	ttDraft, errr := h.getOrCreateDraftTime(int(ds.CurrentHistorian.HistorianID))

	//logrus.Info("adding post army 2")

	if errr != nil {
		// logrus.Error(err)
		// logrus.Info(ttDraft.TtID)
		c.JSON(500, gin.H{
			"error":   err,
			"message": "Не удалось найти или создать черновик!",
		})
		return
	}

	var idDraft = ttDraft.TtID
	//logrus.Info("adding post army 3.1")

	ifArmyInTT := h.Repository.CheckIfArmyAlreadyInTT(idArmy, idDraft)
	if ifArmyInTT {
		//c.Redirect(http.StatusFound, "/armies")
		//logrus.Info("gone to redirect in adding!")
		c.JSON(208, gin.H{
			"message": "Армия уже присутсвует в расчёте!",
		})
		return
	}

	//logrus.Info("adding post army 3.3")

	err = h.Repository.AddArmyToTT(idArmy, idDraft, 0)
	//logrus.Info("adding post army 4")
	if err != nil {
		c.JSON(500, gin.H{
			"error":   err,
			"message": "Не удалось добавить армию в расчёт",
		})
		return
	}

	resDraft, err := h.Repository.GetTravelTime(idDraft)
	if err != nil {
		c.JSON(500, gin.H{
			"error":   err,
			"message": "Не удалось получить черновик расчёта",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Армия добавлена в расчёт",
		"travel_time": resDraft,
	})
}

// func (h *Handler) AddArmyImage(c *gin.Context) {

// }
