package handler

import (
	"time_of_armies/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) RegisterHandler(router *gin.Engine) { // маршрутизация
	// router.GET("/armies", h.GetArmies)                  // выводим все армии
	// router.GET("/one_army/:id", h.GetArmy)              // выводим страницу с отдельной армией
	// router.GET("/travel_time/:ttid", h.GetTravelTime)   // выводим страницу с текущим расчётом
	// router.POST("/add_army_to_tt/:aaid", h.addArmyToTT) // добавляем армию в текущий расчет

	// router.POST("/delete_tt/:ttid", h.deleteTT) // удаляем армию из текущего расчёта

	// armies
	router.GET("/api/armies", h.GetArmies)                       // список с фильтрацией
	router.GET("/api/army/:id", h.GetArmy)                       // одна запись
	router.POST("/api/army", h.AddArmy)                          // Добавление без изображения
	router.PUT("/api/army/:id", h.UpdateArmy)                    // изменение армии
	router.DELETE("/api/army/:id", h.DeleteArmy)                 // "удаление" Армии - смена статуса услуги на "удален"
	router.POST("/api/army/add_to_travel", h.addArmyToTT)        // добавление армии в расчёт черновик
	router.POST("/api/army/:id/upload_image", h.uploadArmyImage) // добавление изображения к армии

	//timeToTravel
	router.GET("/api/travel_time", h.GetTravelTimeDraft)                  // получить id черновика и кол-во услуг в нем
	router.GET("/api/travel_times", h.GetTravelTimes)                     // список (кроме удаленных и черновика, поля модератора и создателя через логины) с фильтрацией по диапазону даты формирования и статусу
	router.GET("/api/travel_time/:ttid", h.GetTravelTime)                 //одна запись (поля заявки + ее услуги). При получении заявки возвращется список ее услуг с картинками
	router.PUT("/api/travel_time/:ttid", h.UpdateTravelTime)              //изменения полей заявки по теме
	router.PUT("/api/travel_time/:ttid/form", h.ToFormTravelTime)         // сформировать создателем (дата формирования). Происходит проверка на обязательные поля
	router.PUT("/api/travel_time/:ttid/moderate", h.ToModerateTravelTime) // завершить/отклонить модератором
	router.DELETE("/api/travel_time/:ttid", h.deleteTT)                   // удаление (дата формирования)

	// M:M
	router.DELETE("/api/travel_time/delete_army", h.DeleteConn) // удаление из заявки (без PK м-м)
	router.PUT("/api/travel_time/update_army", h.UpdateConn)    // PUT изменение количества/порядка/значения в м-м (без PK м-м)

	// Historian
	router.POST("/api/historian/reg", h.RegisterHistorian) //регистрация
	router.GET("/api/historian", h.GetHistorianInfo)       // GET полей пользователя после аутентификации (для личного кабинета)
	router.PUT("/api/historian", h.UpdateHistorianInfo)    // PUT пользователя (личный кабинет)
	router.POST("/api/historian/auth", h.AuthHistorian)    // аутентификация
	router.POST("/api/historian/exit", h.DeauthHistorian)  // деавторизация

	// //testGlobalVar
	// router.GET("/api/testg", h.testGlob)
	// router.GET("/api/upg", h.UpdateGlob)

}

// func (h *Handler) RegisterStatic(router *gin.Engine) {
// 	router.LoadHTMLGlob("templates/*")
// 	router.Static("/static", "./resources")
// }

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
