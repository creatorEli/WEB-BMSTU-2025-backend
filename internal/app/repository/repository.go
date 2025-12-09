package repository

import (
	"fmt"
	"strings"
)

type Repository struct {
}

func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

type Army struct {
	Id        int
	NameArmy  string
	ClassArmy string // для фильтрации
	/*скорости для каждого типа местности:
	Равнина минимальная speeds[0]
	Равнина Максимальная speeds[1]
	Горы минимальная speeds[2]
	Горы Максимальная speeds[3]
	Лес минимальная speeds[4]
	Лес Максимальная speeds[5]
	Река минимальная speeds[6]
	Река Максимальная speeds[7]
	Пустыня минимальная speeds[8]
	Пустыня Максимальная speeds[9]
	*/
	MinPlain  int
	MaxPlain  int
	MinMount  int
	MaxMount  int
	MinForest int
	MaxForest int
	MinRiver  int
	MaxRiver  int
	MinDesert int
	MaxDesert int
	ImageLink string
	//shortDecription  string // не понадобится
}

type Calculation struct {
	Id         int
	Armies     []int
	ChosenBiom string
	Distance   int
	ResultMin  int
	ResultMax  int
}

type TimeOfWay struct {
	Id         int
	Armies     []int // надо будет удалить, переброска в M-M
	ChosenBiom string
	Distance   int
	ResultMin  int
	ResultMax  int
}

func (r *Repository) GetArmies() ([]Army, error) {
	// имитируем работу с БД. Типа мы выполнили sql запрос и получили эти строки из БД
	armies := []Army{
		{
			Id:        1,
			NameArmy:  "Пехота",
			ClassArmy: "step",
			//speedsBioms: [10]int{20, 30, 10, 15, 8, 12, 5, 15, 15, 25},
			MinPlain:  20,
			MaxPlain:  30,
			MinMount:  10,
			MaxMount:  15,
			MinForest: 8,
			MaxForest: 12,
			MinRiver:  5,
			MaxRiver:  15,
			MinDesert: 15,
			MaxDesert: 25,
			ImageLink: "http://127.0.0.1:9000/armies/infantry.jpg",
		},
		{
			Id:        2,
			NameArmy:  "Кавалерия",
			ClassArmy: "horse",
			//speedsBioms: [10]int{40, 60, 20, 30, 15, 25, 5, 15, 30, 50},
			MinPlain:  40,
			MaxPlain:  60,
			MinMount:  20,
			MaxMount:  30,
			MinForest: 15,
			MaxForest: 25,
			MinRiver:  5,
			MaxRiver:  15,
			MinDesert: 30,
			MaxDesert: 50,
			ImageLink: "http://127.0.0.1:9000/armies/Horses.jpg",
		},
		{
			Id:        3,
			NameArmy:  "Обоз",
			ClassArmy: "wheel",
			MinPlain:  15,
			MaxPlain:  20,
			MinMount:  5,
			MaxMount:  10,
			MinForest: 4,
			MaxForest: 8,
			MinRiver:  2,
			MaxRiver:  10,
			MinDesert: 10,
			MaxDesert: 15,
			ImageLink: "http://127.0.0.1:9000/armies/wagon.jpg",
		},
		{
			Id:        4,
			NameArmy:  "Артиллерия",
			ClassArmy: "wheel",
			MinPlain:  10,
			MaxPlain:  15,
			MinMount:  2,
			MaxMount:  5,
			MinForest: 3,
			MaxForest: 6,
			MinRiver:  1,
			MaxRiver:  3,
			MinDesert: 8,
			MaxDesert: 12,
			ImageLink: "http://127.0.0.1:9000/armies/artillery.jpg",
		},
		{
			Id:        5,
			NameArmy:  "Верблюжья кавалерия",
			ClassArmy: "horse",
			MinPlain:  30,
			MaxPlain:  50,
			MinMount:  15,
			MaxMount:  25,
			MinForest: 10,
			MaxForest: 15,
			MinRiver:  5,
			MaxRiver:  15,
			MinDesert: 40,
			MaxDesert: 60,
			ImageLink: "http://127.0.0.1:9000/armies/camels.jpg",
		},
		{
			Id:        6,
			NameArmy:  "Егеря",
			ClassArmy: "step",
			MinPlain:  25,
			MaxPlain:  35,
			MinMount:  15,
			MaxMount:  20,
			MinForest: 12,
			MaxForest: 18,
			MinRiver:  5,
			MaxRiver:  15,
			MinDesert: 15,
			MaxDesert: 25,
			ImageLink: "go",
		},
	}

	if len(armies) == 0 {
		return nil, fmt.Errorf("Армий нет!")
	}

	return armies, nil
}

func (r *Repository) GetArmy(id int) (Army, error) {
	// тут у вас будет логика получения нужной армии, тоже через цикл в первой лабе, и через запрос к БД начиная со второй
	armies, err := r.GetArmies()
	if err != nil {
		return Army{}, err // тут у нас уже есть кастомная ошибка из нашего метода, поэтому мы можем просто вернуть ее
	}

	for _, army := range armies {
		if army.Id == id {
			return army, nil
		}
	}
	return Army{}, fmt.Errorf("Армия не найдена") // тут нужна кастомная ошибка, чтобы понимать на каком этапе возникла ошибка и что произошло
}

// func (r *Repository) GetArmiesByIds(id []int) ([]Army, error) {
// 	// тут у вас будет логика получения нужной армии, тоже через цикл в первой лабе, и через запрос к БД начиная со второй
// 	armiesRes := make([]repository.Army, 0)

// 	armies, err := r.GetArmies()
// 	if err != nil {
// 		return []Army{}, err // тут у нас уже есть кастомная ошибка из нашего метода, поэтому мы можем просто вернуть ее
// 	}

// 	for _, idArmy := range id {
// 		armiesRes = append(armiesRes, armies[idArmy])
// 	}

// 	return armiesRes, nil
// }

func (r *Repository) GetArmyByTitle(title string) ([]Army, error) {
	armies, err := r.GetArmies()
	if err != nil {
		return []Army{}, err
	}

	var result []Army

	for _, army := range armies {
		if strings.Contains(strings.ToLower(army.NameArmy), strings.ToLower(title)) {
			result = append(result, army)
		}
	}

	return result, nil
}

func (r *Repository) GetArmiesByClass(title string) ([]Army, error) {
	armies, err := r.GetArmies()
	if err != nil {
		return []Army{}, err
	}

	var result []Army

	for _, army := range armies {
		if strings.Contains(strings.ToLower(army.ClassArmy), strings.ToLower(title)) {
			result = append(result, army)
		}
	}

	return result, nil
}

func (r *Repository) GetTravelTimes() ([]TimeOfWay, error) {
	// имитируем обращение к БД и выдачу всех наших расчетов (заявок)
	timesOfWays := []TimeOfWay{ //т.к. в первой лабе работаем только с текущей заявкой, то и в коллекции оставляю только одну заявку
		{
			Id:         1,
			Armies:     []int{1, 2, 3},
			ChosenBiom: "Desert",
			Distance:   120,
			ResultMin:  8,
			ResultMax:  12,
		},
	}

	if len(timesOfWays) == 0 {
		return nil, fmt.Errorf("Заявок нет!")
	}

	return timesOfWays, nil
}

func (r *Repository) GetTravelTime(id int) (TimeOfWay, error) {
	timesOfWays, err := r.GetTravelTimes()
	if err != nil {
		return TimeOfWay{}, err // тут у нас уже есть кастомная ошибка из нашего метода, поэтому мы можем просто вернуть ее
	}

	for _, timeW := range timesOfWays {
		if timeW.Id == id {
			return timeW, nil // если нашли, то просто возвращаем найденный заказ (услугу) без ошибок
		}
	}
	return TimeOfWay{}, fmt.Errorf("Расчет не найден") // тут нужна кастомная ошибка, чтобы понимать на каком этапе возникла ошибка и что произошло
}

func (r *Repository) CountArmiesInTime(id int) int {
	timesOfWays, err := r.GetTravelTimes()
	if err != nil {
		return 0 // тут у нас уже есть кастомная ошибка из нашего метода, поэтому мы можем просто вернуть ее
	}

	for _, timeW := range timesOfWays {
		if timeW.Id == id {
			return len(timeW.Armies) // если нашли, то просто возвращаем найденный заказ (услугу) без ошибок
		}
	}
	return 0 // тут нужна кастомная ошибка, чтобы понимать на каком этапе возникла ошибка и что произошло
}

func (r *Repository) GetLastTime() int {
	travelTimes, err := r.GetTravelTimes()
	if err != nil {
		return 0 // тут у нас уже есть кастомная ошибка из нашего метода, поэтому мы можем просто вернуть ее
	}
	return len(travelTimes) // пока забиваем в самый конец списка заявок.
}
