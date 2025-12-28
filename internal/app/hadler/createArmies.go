package handler

import (
	"math/rand"
	//"net/http"
	"time"
	"time_of_armies/internal/app/ds"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Базовая конфигурация для разных типов армий
var armyConfigs = map[string]struct {
	speedRanges map[string][2]int
}{
	"step": {
		speedRanges: map[string][2]int{
			"plain":  {20, 35},
			"mount":  {10, 20},
			"forest": {8, 18},
			"river":  {5, 15},
			"desert": {15, 25},
		},
	},
	"horse": {
		speedRanges: map[string][2]int{
			"plain":  {30, 60},
			"mount":  {15, 30},
			"forest": {10, 25},
			"river":  {5, 15},
			"desert": {30, 60},
		},
	},
	"wheel": {
		speedRanges: map[string][2]int{
			"plain":  {10, 20},
			"mount":  {2, 10},
			"forest": {3, 8},
			"river":  {1, 10},
			"desert": {8, 15},
		},
	},
}

var ArmiesClasses = []string{"step", "horse", "wheel"}

// Генератор случайных названий армий

// // Генерация случайного названия армии
// func generateArmyName(r *rand.Rand) string {
// 	parts := []string{}

// 	// Всегда добавляем тип
// 	parts = append(parts, namesArmies[rand.Intn(len(namesArmies))])

// 	// 30% шанс добавить номер
// 	if r.Float32() < 0.3 {
// 		parts = append(parts, fmt.Sprintf("%d", r.Intn(100)+1))
// 	}

// 	return strings.Join(parts, " ")
// }

// Генерация скорости с небольшим разбросом
func generateSpeed(baseMin, baseMax int, r *rand.Rand) (int, int) {
	variation := r.Intn(5) - 2 // случайное отклонение от -2 до +2
	min := baseMin + variation
	max := baseMax + variation

	// Гарантируем, что min <= max и значения не отрицательные
	if min < 1 {
		min = 1
	}
	if max < min {
		max = min + r.Intn(5) + 1
	}
	if max > 100 {
		max = 100
	}

	return min, max
}

var namesArmiesStep = []string{
	"пехота", "гвардия", "легион", "янычары",
	"дивизия пешая", "Копьеносцы", "полк", "батальон", "рота",
	"драгуны пешие", "кирасиры", "егеря", "снайперы", "Лучники",
}

var namesArmiesHorse = []string{
	"Рыцари тяжеловозы", "драгуны конные", "Боевые слоны",
	"Крылатые гусары", "Казаки",
	"кавалерия", "гусары", "уланы",
	"Кочевники",
}

var namesArmiesWheel = []string{
	"артиллерия", "колесницы", "Башни", "Катапульты", "Мортиры", "Пушки", "Гаубицы", "Обозы",
}

// Основная функция генерации армии
func GenerateArmy(r *rand.Rand) (ds.Army, error) {

	var curArmy = ds.Army{}

	curArmy.ClassArmy = ArmiesClasses[r.Intn(3)]

	config := armyConfigs[curArmy.ClassArmy]

	// Генерация скоростей
	curArmy.MinPlainSpeed, curArmy.MaxPlainSpeed = generateSpeed(
		config.speedRanges["plain"][0],
		config.speedRanges["plain"][1],
		r,
	)
	curArmy.MinMountSpeed, curArmy.MaxMountSpeed = generateSpeed(
		config.speedRanges["mount"][0],
		config.speedRanges["mount"][1],
		r,
	)
	curArmy.MinForestSpeed, curArmy.MaxForestSpeed = generateSpeed(
		config.speedRanges["forest"][0],
		config.speedRanges["forest"][1],
		r,
	)
	curArmy.MinRiverSpeed, curArmy.MaxRiverSpeed = generateSpeed(
		config.speedRanges["river"][0],
		config.speedRanges["river"][1],
		r,
	)
	curArmy.MinDesertSpeed, curArmy.MaxDesertSpeed = generateSpeed(
		config.speedRanges["desert"][0],
		config.speedRanges["desert"][1],
		r,
	)

	if curArmy.ClassArmy == "step" {
		curArmy.NameArmy = namesArmiesStep[r.Intn(len(namesArmiesStep))]
	}
	if curArmy.ClassArmy == "horse" {
		curArmy.NameArmy = namesArmiesHorse[r.Intn(len(namesArmiesHorse))]
	}
	if curArmy.ClassArmy == "wheel" {
		curArmy.NameArmy = namesArmiesWheel[r.Intn(len(namesArmiesWheel))]
	}

	curArmy.StatusArmy = "действует"
	return curArmy, nil
}

func (h *Handler) CreateArmiesss(c *gin.Context) {
	//logrus.Info("Типо создали 100500 Армий!\n")
	//var result = []ds.Army{}
	for i := range 863 {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		curArmy, err := GenerateArmy(r)
		if err != nil {
			logrus.Error("На этапе генерации армии произошла ошибка!")
			logrus.Error("Итерация: ", i)
			logrus.Error("Ошибка: ", err.Error())
			return
		}
		//result = append(result, curArmy)
		_, err = h.Repository.AddArmy(curArmy)
		if err != nil {
			logrus.Error("Не удалось добавить армию в БД!")
			logrus.Error("Итерация: ", i)
			logrus.Error("Ошибка: ", err.Error())
			return
		}
	}
	logrus.Info("Армии сгенерировались!")
	// c.JSON(http.StatusOK, gin.H{
	// 	"armies": result,
	// })
}
