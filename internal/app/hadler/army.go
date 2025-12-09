package handler

import (
	"net/http"
	"strconv"
	"time_of_armies/internal/app/ds"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var creatorID = 1 // хардкод пока нету функцонала юзера

func (h *Handler) addArmyToTT(c *gin.Context) {

	logrus.Info("adding post army 0")
	idStr := c.Param("aaid")
	idArmy, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error(err)
		return
	}

	logrus.Info("adding post army 1")
	ttDraft, errr := h.getOrCreateDraftTime()
	logrus.Info("adding post army 2")

	if errr != nil {
		logrus.Error(err)
		logrus.Info(ttDraft.TtID)
		return
	}

	var idDraft = ttDraft.TtID
	logrus.Info("adding post army 3.1")

	ifArmyInTT := h.Repository.CheckIfArmyAlreadyInTT(idArmy, idDraft)
	if ifArmyInTT {
		c.Redirect(http.StatusFound, "/armies")
		logrus.Info("gone to redirect in adding!")
		return
	}

	logrus.Info("adding post army 3.3")

	err = h.Repository.AddArmyToTT(idArmy, idDraft)
	logrus.Info("adding post army 4")
	if err != nil {
		logrus.Error(err)
		return
	}

	// Добавили в заявку и переходим обратно на страницу всех услуг
	c.Redirect(http.StatusFound, "/armies")
}

func (h *Handler) GetArmies(c *gin.Context) {

	logrus.Info("GetArmies!")

	// если есть расчёт черновик у юзера, то добавляем ссылку, иначе впихиваем якорь и нуль число

	//ttDraft, errr := h.getOrCreateDraftTime()
	//logrus.Info("GetArmies 1!")

	// if errr != nil { // число логическая ошибка была
	// 	// c.HTML(http.StatusOK, "armies.html", gin.H{
	// 	// 	"armies":             []ds.Army{},
	// 	// 	"armySearchQuery":    "",
	// 	// 	"draftTTid":          -1,
	// 	// 	"countArmiesTimeBTN": 0,
	// 	// })
	// 	// return
	// 	logrus.Warn("this appears after deleting draft!")
	// }
	//logrus.Info("GetArmies 2!")

	var href string = "/armies"
	var armiesInDraft int64 = 0
	ttDraft, errorr := h.Repository.GetTTDraft(creatorID)
	if errorr == nil {
		armiesInDraft = h.Repository.CountArmiesInTime(ttDraft.TtID)
		href = "/travel_time/" + strconv.Itoa(ttDraft.TtID)
	}

	var armies []ds.Army

	var err error
	searchArmyQuery := c.Query("searchNameArmy")
	filter := c.Query("class")

	// отлично, теперь есть черновик, беерм его id и пихаем в шаблон!

	//var countArmiesCurTime = h.Repository.CountArmiesInTime(ttDraft.TtID)

	if searchArmyQuery != "" {
		armies, err = h.Repository.GetArmyByTitle(searchArmyQuery)
		if err != nil {
			logrus.Error(err)
		}
		c.HTML(http.StatusOK, "armies.html", gin.H{
			"armies":          armies,
			"armySearchQuery": searchArmyQuery,
			//"draftTTid":          ttDraft.TtID,
			"countArmiesTimeBTN": armiesInDraft,
			"hrefToTT":           href,
		})
		return
	}

	if filter != "" {
		armies, err = h.Repository.GetArmiesByClass(filter)
		if err != nil {
			logrus.Error(err)
		}
		c.HTML(http.StatusOK, "armies.html", gin.H{
			"armies": armies,
			//"draftTTid":          ttDraft.TtID,
			"countArmiesTimeBTN": armiesInDraft,
			"hrefToTT":           href,
		})
		return
	}

	//logrus.Info(armies)

	armies, err = h.Repository.GetArmies()
	if err != nil {
		logrus.Error(err)
	}

	c.HTML(http.StatusOK, "armies.html", gin.H{
		"armies": armies,
		//"draftTTid":          ttDraft.TtID,
		"countArmiesTimeBTN": armiesInDraft,
		"hrefToTT":           href,
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

	// то же самое
	// ttDraft, errr := h.getOrCreateDraftTime()
	// if errr != nil {
	// 	return
	// }

	var href string = "/armies"
	var armiesInDraft int64 = 0
	ttDraft, errorr := h.Repository.GetTTDraft(creatorID)
	if errorr == nil {
		armiesInDraft = h.Repository.CountArmiesInTime(ttDraft.TtID)
		href = "/travel_time/" + strconv.Itoa(ttDraft.TtID)
	}

	//var countArmiesCurTime = h.Repository.CountArmiesInTime(ttDraft.TtID)
	c.HTML(http.StatusOK, "onearmy.html", gin.H{
		"army": oneArmy,
		//"draftTTid":          ttDraft.TtID,
		"countArmiesTimeBTN": armiesInDraft,
		"hrefToTT":           href,
	})
}
