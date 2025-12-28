package handler

import (
	//"encoding/json"
	"net/http"
	"strconv"
	"time"
	"time_of_armies/internal/app/ds"

	_ "time_of_armies/docs"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Paginationn struct {
	Page       int
	Limit      int
	Total      int64
	TotalPages int64
}

type ResArmies struct {
	Armies         []ds.Army
	Pagination     Paginationn
	QueryTimeMs    int64
	QueryWithIndex bool
}

// GetArmies godoc
// @Summary      Получить список армий
// @Description  Получить список армий, включая фильтрацию (param "class") и поисковый запрос (param "searchNameArmy") а также пагинацию (param "page") и число записей на страцие (param "limit")
// @Tags         Requests
// @Accept 		 json
// @Produce      json
// Param [name] [type] [dataType] [required] [description]
// @Param		 class query string false "фильтрация"
// @Param 		 searchNameArmy query string false "поиск армии"
// @Param		 page query string false "страница"
// @Param		 limit query string false "число записей на страницу"
// @Param		 withIndexation query bool false "число записей на страницу"
// @Success      200 {object} ResArmies
// @Router       /armies [get]
func (h *Handler) GetArmies(c *gin.Context) {
	logrus.Info("GetArmies!")
	startTime := time.Now()

	withIndexation := c.Query("withIndexation")
	logrus.Info("withIndexation = ", withIndexation)

	// Параметры пагинации
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "12"))

	offset := (page - 1) * limit

	if limit > 120 {
		limit = 120 // ограничение на количество записей
	}

	searchArmyQuery := c.Query("searchNameArmy")
	filter := c.Query("class")
	var totalCount int64

	if withIndexation == "true" {
		var armies []ds.Army
		var err error

		if searchArmyQuery != "" {
			armies, err = h.Repository.GetArmyByTitle(searchArmyQuery, offset, limit, &totalCount)
			if err != nil {
				logrus.Error(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			duration := time.Since(startTime)
			pagRes := Paginationn{
				page, limit, totalCount, (totalCount + int64(limit) - 1) / int64(limit),
			}
			armRes := ResArmies{
				armies, pagRes, duration.Milliseconds(), true,
			}
			c.JSON(http.StatusOK, armRes)
			return
		}

		if filter != "" {
			armies, err = h.Repository.GetArmiesByClass(filter, offset, limit, &totalCount)
			if err != nil {
				logrus.Error(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			duration := time.Since(startTime)
			pagRes := Paginationn{
				page, limit, totalCount, (totalCount + int64(limit) - 1) / int64(limit),
			}
			armRes := ResArmies{
				armies, pagRes, duration.Milliseconds(), true,
			}
			c.JSON(http.StatusOK, armRes)
			// c.JSON(http.StatusOK, gin.H{
			// 	"armiess": armies,
			// 	"pagination": gin.H{
			// 		"page":       page,
			// 		"limit":      limit,
			// 		"total":      totalCount,
			// 		"totalPages": (totalCount + int64(limit) - 1) / int64(limit),
			// 	},
			// 	"query_time_ms":    duration.Milliseconds(),
			// 	"query_with_index": true, // для тестирования производительности
			// })
			return
		}
		armies, err = h.Repository.GetArmies(offset, limit, &totalCount)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		duration := time.Since(startTime)
		pagRes := Paginationn{
			page, limit, totalCount, (totalCount + int64(limit) - 1) / int64(limit),
		}
		armRes := ResArmies{
			armies, pagRes, duration.Milliseconds(), true,
		}
		c.JSON(http.StatusOK, armRes)
		// c.JSON(http.StatusOK, gin.H{
		// 	"armies": armies,
		// 	"pagination": gin.H{
		// 		"page":       page,
		// 		"limit":      limit,
		// 		"total":      totalCount,
		// 		"totalPages": (totalCount + int64(limit) - 1) / int64(limit),
		// 	},
		// 	"query_time_ms":    duration.Milliseconds(),
		// 	"query_with_index": true, // для тестирования производительности
		// })
	} else { // ЕСЛИ ОТКЛЮЧИЛИ ИНДЕКСАЦИЮ
		armies, err := h.Repository.GetArmiesWithoutIndexation(searchArmyQuery, filter, offset, limit, &totalCount)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		//logrus.Info("totalCount = ", totalCount)

		duration := time.Since(startTime)
		pagRes := Paginationn{
			page, limit, totalCount, (totalCount + int64(limit) - 1) / int64(limit),
		}
		armRes := ResArmies{
			armies, pagRes, duration.Milliseconds(), false,
		}
		c.JSON(http.StatusOK, armRes)
		// c.JSON(http.StatusOK, gin.H{
		// 	"data": armies,
		// 	"pagination": gin.H{
		// 		"page":       page,
		// 		"limit":      limit,
		// 		"total":      totalCount,
		// 		"totalPages": (totalCount + int64(limit) - 1) / int64(limit),
		// 	},
		// 	"query_time_ms":    duration.Milliseconds(),
		// 	"query_with_index": false,
		// })

	}

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
// @Param		 ArmyID formData int true "ID добавляемой армии"
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
	if int(ds.CurrentHistorian.HistorianID) == 0 {
		c.JSON(401, gin.H{
			"error":   err,
			"message": "Вы не авторизованы!",
		})
		return
	}
	logrus.Info("HisID current: ", int(ds.CurrentHistorian.HistorianID))
	ttDraft, errr := h.getOrCreateDraftTime(int(ds.CurrentHistorian.HistorianID))

	//logrus.Info("adding post army 2")

	if errr != nil {
		// logrus.Error(err)
		// logrus.Info(ttDraft.TtID)
		c.JSON(400, gin.H{
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
