package handler

import (
	"time_of_armies/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
}

//var currentTimeId = 1 // ПЕРЕДЕЛАТЬ НА ЗАПРОС ПО ЗАЯВКЕ ТЕКУЩЕЙ!

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) RegisterHandler(router *gin.Engine) { // маршрутизация
	router.GET("/armies", h.GetArmies)                  // выводим все армии
	router.GET("/one_army/:id", h.GetArmy)              // выводим страницу с отдельной армией
	router.GET("/travel_time/:ttid", h.GetTravelTime)   // выводим страницу с текущим расчётом
	router.POST("/add_army_to_tt/:aaid", h.addArmyToTT) // добавляем армию в текущий расчет
	// router.POST("/add_army_to_tt/:aaid", func(ctx *gin.Context) {
	// 	ctx.JSON(http.StatusOK, gin.H{ // отправляем ответ в виде json
	// 		"message": "added so post army",
	// 	})
	// })
	router.POST("/delete_tt/:ttid", h.deleteTT) // удаляем армию из текущего расчёта
}

func (h *Handler) RegisterStatic(router *gin.Engine) { // changed here
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./resources") //W:\BM57U\WEB\for lab 1\lab1\resources\img\logo_ozon_1.webp
	// слева название папки, в которую выгрузится наша статика
	// справа путь к папке, в которой лежит статика
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) { // changed here
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
