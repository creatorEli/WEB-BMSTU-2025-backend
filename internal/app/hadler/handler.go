package handler

import (
	"time_of_armies/internal/app/repository"
	"time_of_armies/internal/app/role"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	_ "time_of_armies/docs" // !! оч нужно чтобы swagger ui работал

	cfg "time_of_armies/internal/app/config"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"time_of_armies/internal/app/redis"
)

type Handler struct {
	Repository  *repository.Repository
	config      *cfg.Config
	redis       *redis.Client
	DjangoURL   string // URL Django сервиса
	DjangoToken string // Токен для авторизации
}

func NewHandler(r *repository.Repository, cfg *cfg.Config, redisClient *redis.Client) *Handler {
	return &Handler{
		Repository:  r,
		config:      cfg,
		redis:       redisClient,
		DjangoURL:   "http://localhost:8085",
		DjangoToken: "12345678",
	}
}

func (h *Handler) RegisterHandler(router *gin.Engine) { // маршрутизация
	// router.GET("/armies", h.GetArmies)                  // выводим все армии
	// router.GET("/one_army/:id", h.GetArmy)              // выводим страницу с отдельной армией
	// router.GET("/travel_time/:ttid", h.GetTravelTime)   // выводим страницу с текущим расчётом
	// router.POST("/add_army_to_tt/:aaid", h.addArmyToTT) // добавляем армию в текущий расчет

	// router.POST("/delete_tt/:ttid", h.deleteTT) // удаляем армию из текущего расчёта

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// armies

	//timeToTravel
	//ПОЛУЧАТЬ ОТВЕТ ОБ ОТСТУТСТВИИ ЧЕРНОВИКА И БЕЗ ID

	// Historian

	// крч сперва публичные маршруты, а потом уже и защищённые

	// действия, доступные всем (не используем проверку на роли):
	//router.Use(h.WithAuthCheck(role.Historian, role.Moderator)).GET("/api/armies", h.GetArmies) // список с фильтрацией
	router.POST("/api/travel_time/:ttid/update_calc", h.CalculationCallback)

	router.GET("/api/armies", h.GetArmies)                 // список с фильтрацией
	router.GET("/api/army/:id", h.GetArmy)                 // одна запись
	router.POST("/api/historian/reg", h.RegisterHistorian) //регистрация
	router.POST("/api/historian/auth", h.AuthHistorian)    // аутентификация
	//router.Use(h.WithAuthCheck(role.Historian, role.Moderator)).GET("/api/travel_time", h.GetTravelTimeDraft) // получить id черновика и кол-во услуг в нем
	router.GET("/api/travel_time", h.GetTravelTimeDraft) // получить id черновика и кол-во услуг в нем

	// действия, доступные авторизованным пользователям:
	router.Use(h.WithAuthCheck(role.Historian, role.Moderator)).POST("/api/army/add_to_travel", h.addArmyToTT) // добавление армии в расчёт черновик
	router.Use(h.WithAuthCheck(role.Historian, role.Moderator)).GET("/api/historian", h.GetHistorianInfo)      // GET полей пользователя после аутентификации (для личного кабинета)
	router.Use(h.WithAuthCheck(role.Historian, role.Moderator)).PUT("/api/historian", h.UpdateHistorianInfo)   // PUT пользователя (личный кабинет)
	router.Use(h.WithAuthCheck(role.Historian, role.Moderator)).POST("/api/historian/exit", h.DeauthHistorian) // деавторизация

	// M:M
	router.Use(h.WithAuthCheck(role.Historian, role.Moderator)).DELETE("/api/travel_time/delete_army", h.DeleteConn)   // удаление из заявки (без PK м-м)
	router.Use(h.WithAuthCheck(role.Historian, role.Moderator)).PUT("/api/travel_time/update_army", h.UpdateConn)      // PUT изменение количества/порядка/значения в м-м (без PK м-м)
	router.Use(h.WithAuthCheck(role.Historian, role.Moderator)).GET("/api/travel_times", h.GetTravelTimes)             // список (кроме удаленных и черновика, поля модератора и создателя через логины) с фильтрацией по диапазону даты формирования и статусу
	router.Use(h.WithAuthCheck(role.Historian, role.Moderator)).GET("/api/travel_time/:ttid", h.GetTravelTime)         //одна запись (поля заявки + ее услуги). При получении заявки возвращется список ее услуг с картинками
	router.Use(h.WithAuthCheck(role.Historian, role.Moderator)).PUT("/api/travel_time/:ttid", h.UpdateTravelTime)      //изменения полей заявки по теме
	router.Use(h.WithAuthCheck(role.Historian, role.Moderator)).PUT("/api/travel_time/:ttid/form", h.ToFormTravelTime) // сформировать создателем (дата формирования). Происходит проверка на обязательные поля
	router.Use(h.WithAuthCheck(role.Historian, role.Moderator)).DELETE("/api/travel_time/:ttid", h.deleteTT)           // удаление расчёта (дата формирования)

	// действия, доступные модератору:
	router.Use(h.WithAuthCheck(role.Moderator)).POST("/api/army", h.AddArmy)                                   // Добавление без изображения
	router.Use(h.WithAuthCheck(role.Moderator)).PUT("/api/army/:id", h.UpdateArmy)                             // изменение армии
	router.Use(h.WithAuthCheck(role.Moderator)).DELETE("/api/army/:id", h.DeleteArmy)                          // "удаление" Армии - смена статуса услуги на "удален"
	router.Use(h.WithAuthCheck(role.Moderator)).POST("/api/army/:id/upload_image", h.uploadArmyImage)          // добавление изображения к армии
	router.Use(h.WithAuthCheck(role.Moderator)).PUT("/api/travel_time/:ttid/moderate", h.ToModerateTravelTime) // завершить/отклонить модератором
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
