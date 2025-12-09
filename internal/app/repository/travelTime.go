package repository

import (
	"fmt"
	"time_of_armies/internal/app/ds"

	"github.com/sirupsen/logrus"
)

func (r *Repository) GetTravelTimes() ([]ds.TravelTime, error) {
	var times []ds.TravelTime
	err := r.db.Find(&times).Error

	if err != nil {
		return nil, err
	}

	if len(times) == 0 {
		return nil, fmt.Errorf("Расчётов путей не найдено!")
	}

	return times, nil
}

func (r *Repository) GetTravelTime(id int) (ds.TravelTime, error) {
	tt := ds.TravelTime{}
	err := r.db.Where("tt_id = ?", id).First(&tt).Error
	if err != nil {
		return ds.TravelTime{}, err
	}
	return tt, nil
}

// func (r *Repository) GetLastTime() int { // старый код остался для запихивания расчёта в конец таблицу
// 	travelTimes, err := r.GetTravelTimes()
// 	if err != nil {
// 		return 0 // тут у нас уже есть кастомная ошибка из нашего метода, поэтому мы можем просто вернуть ее
// 	}
// 	return len(travelTimes) // пока забиваем в самый конец списка заявок.
// }

func (r *Repository) GetTTDraft(creatorID int) (ds.TravelTime, error) {

	resultTT := ds.TravelTime{}
	err := r.db.Model(&ds.TravelTime{}).Where("Creator_ID_tt = ? AND Status_tt = ?", creatorID, "черновик").First(&resultTT).Error
	if err != nil {
		return ds.TravelTime{}, err
	}
	return resultTT, nil
}

func (r *Repository) AddTTInDB(tt *ds.TravelTime) (ds.TravelTime, error) {
	//	creatorID := 1
	// createrID уже есть в tt
	err := r.db.Model(&ds.TravelTime{}).Create(tt).Error

	if err != nil {
		return ds.TravelTime{}, err
	}

	// получм теперь этот же черновик (единственный такой на одного юзера)
	resultTT, err := r.GetTTDraft(int(tt.CreatorID_TT))
	if err != nil {
		return ds.TravelTime{}, err
	}
	return resultTT, nil
}

// для получения количества услуг в заявке
func (r *Repository) CountArmiesInTime(idTT int) int64 {
	//var timeID uint
	var countArmies int64
	creatorID := 1
	// пока что мы захардкодили id создателя заявки, в последующем сделаем авторизацию и будем получать его из JWT

	err := r.db.Model(&ds.TravelTime{}).Where("Creator_ID_tt = ? AND Status_tt = ?", creatorID, "черновик").Select("tt_id").First(&idTT).Error
	if err != nil {
		logrus.Println("Ошибка, не удалось найти текущий расчёт: ", err)
		return 0
	}

	err = r.db.Model(&ds.ConnectArmyTT{}).Where("con_travel_time_id = ?", idTT).Count(&countArmies).Error
	if err != nil {
		logrus.Println("Ошибка, не удалось подсчитать число армий в расчёте: ", err)
		return 0
	}

	return countArmies
}

func (r *Repository) DeleteTT(TTid int) error {
	// надо без ORM - юзаем курсоры
	query := `UPDATE Travel_Times SET Status_tt = 'удален' WHERE Tt_ID = $1`

	err := r.db.Exec(query, TTid).Error
	if err != nil {
		return fmt.Errorf("ошибка при удалении расчета c id %d: %w", TTid, err)
	}
	return nil
}
