package repository

import (
	"fmt"
	"time"
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

func (r *Repository) GetTravelTimesByStatuses(statusTT []string) ([]ds.TravelTime, error) {
	var times []ds.TravelTime
	for _, status := range statusTT {
		//logrus.Info(status)
		var curTT []ds.TravelTime
		err := r.db.Where("status_tt = ?", status).Find(&curTT).Error
		if err != nil {
			return nil, err
		}

		times = append(times, curTT...)
	}

	if len(times) == 0 {
		return nil, fmt.Errorf("Расчётов путей не найдено!")
	}
	return times, nil
}

func (r *Repository) GetTravelTimesByDates(dateFrom time.Time, dateTo time.Time) ([]ds.TravelTime, error) {
	var times []ds.TravelTime
	//if dateFrom == "" || dateFrom == nil || dateFrom '' // по ТЗ не требуется обрабатывать логические ошибки даты!
	err := r.db.Find(&times).Where("date_create_tt >= ? AND date_create_tt <= ?", dateFrom, dateTo).Error
	if err != nil {
		return nil, err
	}

	if len(times) == 0 {
		return nil, fmt.Errorf("Расчётов путей не найдено!")
	}
	return times, nil
}

func (r *Repository) GetTravelTimesByDatesAndStatuses(dateFrom time.Time, dateTo time.Time, statusTT []string) ([]ds.TravelTime, error) {
	var times []ds.TravelTime
	//logrus.Info(dateFrom, dateTo, statusTT)
	for _, status := range statusTT {
		//logrus.Info(status)
		var curTT []ds.TravelTime
		err := r.db.Where("status_tt = ? AND date_create_tt >= ? AND date_create_tt <= ?", status, dateFrom, dateTo).Find(&curTT).Error
		if err != nil {
			return nil, err
		}

		times = append(times, curTT...)
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

func (r *Repository) GetTTDraft(creatorID int) (ds.TravelTime, error) {

	resultTT := ds.TravelTime{}

	err := r.db.Model(&ds.TravelTime{}).Where("Creator_ID_tt = ? AND Status_tt = ?", creatorID, "черновик").First(&resultTT).Error
	if err != nil {
		if err.Error() == "record not found" {
			return ds.TravelTime{}, fmt.Errorf("404")
		} else {
			return ds.TravelTime{}, err
		}
	}
	logrus.Info(err)
	return resultTT, nil
}

func (r *Repository) AddTTInDB(tt *ds.TravelTime) (ds.TravelTime, error) {
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
	creatorID := ds.CurrentHistorian.HistorianID
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
	// в лабе 3 юзай ORM!!!
	err := r.db.Model(&ds.TravelTime{}).Where("tt_id = ?", TTid).Update("status_tt", "удален").Error
	if err != nil {
		return fmt.Errorf("ошибка при удалении расчета c id %d: %w", TTid, err)
	}
	err = r.db.Model(&ds.TravelTime{}).Where("tt_id = ?", TTid).Update("date_create_tt", time.Now()).Error
	if err != nil {
		return fmt.Errorf("ошибка при удалении расчета (изменение времени) c id %d: %w", TTid, err)
	}
	return nil
}

func (r *Repository) EditTravelTime(travel_time ds.TravelTime) (ds.TravelTime, error) {
	err := r.db.Model(&travel_time).Updates(ds.TravelTime{
		StatusTT:        travel_time.StatusTT,
		DateUpdateTT:    travel_time.DateUpdateTT,
		ChosenBiomTT:    travel_time.ChosenBiomTT,
		DistanceTT:      travel_time.DistanceTT,
		ResultMinTT:     travel_time.ResultMinTT,
		ResultMaxTT:     travel_time.ResultMaxTT,
		ResultChronical: travel_time.ResultChronical,
	}).Error
	if err != nil {
		fmt.Errorf("ошибка при обновлении расчёта пути!: %w", err)
		return ds.TravelTime{}, err
	}
	return travel_time, nil
}

func (r *Repository) CheckFieldsTravelTime(TTid int) error {
	travel_time, err := r.GetTravelTime(TTid)
	if err != nil {
		return err
	}

	if r.CountArmiesInTime(TTid) == 0 {
		return fmt.Errorf("Нужно выбрать хотя бы одну Армию для расчёта!")
	}

	ifBiom := false

	for _, curBiom := range []string{"Desert", "Plain", "Mount", "Forest", "River"} {
		if travel_time.ChosenBiomTT == curBiom {
			ifBiom = true
			break
		}
	}

	if !ifBiom {
		return fmt.Errorf("Не выбран биом или он не соответсвует списку: [Desert, Plain, Mount, Forest, River]")
	}

	if travel_time.DistanceTT <= 0 {
		return fmt.Errorf("Пройденное расстояние должно быть натуральным числом!")
	}

	return nil
}

// func (r *Repository) SetTravelTimeFormed(TTid int) error {
// 	err := r.db.Model(&ds.TravelTime{}).Where("tt_id = ?", TTid).Update("status_tt", "сформирован").Error
// 	if err != nil {
// 		return fmt.Errorf("ошибка при смене статуса расчета c id %d: %w", TTid, err)
// 	}
// 	return nil
// }

func (r *Repository) SetTravelTimeStatus(TTid int, status string) (ds.TravelTime, error) {
	travelTime := ds.TravelTime{}
	err := r.db.Model(&travelTime).Where("tt_id = ?", TTid).Update("status_tt", status).Error
	if err != nil {
		return ds.TravelTime{}, fmt.Errorf("ошибка при смене статуса расчета c id %d: %w", TTid, err)
	}
	travelTime, err = r.GetTravelTime(TTid)
	if err != nil {
		return ds.TravelTime{}, fmt.Errorf("ошибка при получении изменённого расчёта")
	}
	//logrus.Info(travelTime)
	return travelTime, nil
}
