package handler

// import (
// 	"net/http"
// 	"strconv"
// 	"time"
// 	"time_of_armies/internal/app/ds"

// 	"github.com/gin-gonic/gin"
// )

// // Получение армий с пагинацией и фильтрацией
// func (h *Handler) GetArmiesIndPag(c *gin.Context) {
// 	startTime := time.Now()

// 	// Параметры пагинации
// 	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
// 	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "18"))
// 	offset := (page - 1) * limit

// 	if limit > 120 {
// 		limit = 120 // ограничение на количество записей
// 	}

// 	// Параметры фильтрации
// 	class := c.Query("class")
// 	status := c.Query("status")
// 	minSpeed, _ := strconv.Atoi(c.Query("min_speed"))
// 	maxSpeed, _ := strconv.Atoi(c.Query("max_speed"))

// 	// Построение запроса
// 	query := h.db.Model(&Army{})

// 	if class != "" {
// 		query = query.Where("class_army = ?", class)
// 	}
// 	// Получение общего количества для пагинации
// 	var totalCount int64
// 	query.Count(&totalCount)

// 	// Получение данных с пагинацией
// 	var armies []ds.Army
// 	err := query.Offset(offset).Limit(limit).Order("army_id").Find(&armies).Error

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Время выполнения запроса
// 	duration := time.Since(startTime)

// 	c.JSON(http.StatusOK, gin.H{
// 		"data": armies,
// 		"pagination": gin.H{
// 			"page":       page,
// 			"limit":      limit,
// 			"total":      totalCount,
// 			"totalPages": (totalCount + int64(limit) - 1) / int64(limit),
// 		},
// 		"query_time_ms":    duration.Milliseconds(),
// 		"query_with_index": true, // для тестирования производительности
// 	})
// }

// // Метод без использования индекса (для демонстрации)
// func (h *ArmyHandler) GetArmiesWithoutIndex(c *gin.Context) {
// 	startTime := time.Now()

// 	// Принудительное отключение использования индексов
// 	query := h.db.Raw(`
// 		SELECT set enable_indexscan = off;
// 		SELECT set enable_bitmapscan = off;

// 		SELECT * FROM armies
// 		WHERE class_army = ? AND status_army = ?
// 		ORDER BY army_id
// 		LIMIT ? OFFSET ?
// 	`,
// 		c.Query("class"),
// 		c.Query("status"),
// 		c.DefaultQuery("limit", "20"),
// 		(c.DefaultQuery("page", "1")-1)*c.DefaultQuery("limit", "20"),
// 	)

// 	var armies []ds.Army
// 	err := query.Scan(&armies).Error

// 	duration := time.Since(startTime)

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"data":             armies,
// 		"query_time_ms":    duration.Milliseconds(),
// 		"query_with_index": false,
// 	})
// }
