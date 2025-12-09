package repository

import (
	"fmt"
	"strings"
	"time_of_armies/internal/app/ds"

	"github.com/sirupsen/logrus"
)

func (r *Repository) GetArmies() ([]ds.Army, error) {

	var armies []ds.Army
	err := r.db.Find(&armies).Error

	if err != nil {
		return nil, err
	}

	if len(armies) == 0 {
		return nil, fmt.Errorf("Армий нет!")
	}

	return armies, nil
}

func (r *Repository) GetArmy(id int) (ds.Army, error) {
	army := ds.Army{}
	err := r.db.Where("army_id = ?", id).First(&army).Error
	if err != nil {
		return ds.Army{}, err
	}
	return army, nil
}

func (r *Repository) GetArmyByTitle(title string) ([]ds.Army, error) {
	var armies []ds.Army
	err := r.db.Where("Name_Army ILIKE ?", "%"+title+"%").Find(&armies).Error
	if err != nil {
		return nil, err
	}
	return armies, nil
}

func (r *Repository) GetArmiesByClass(title string) ([]ds.Army, error) {
	var armies []ds.Army
	err := r.db.Where("Class_Army ILIKE ?", "%"+strings.ToLower(title)+"%").Find(&armies).Error
	if err != nil {
		return nil, err
	}
	return armies, nil
}

func (r *Repository) GetArmiesByTTid(TTid int) ([]ds.Army, error) { // Для отображения в заявке принадлежащих ей армий
	var ArmiesIdsFromConArmyTT []ds.ConnectArmyTT
	// сперва собираем все записи с нашей заявкой
	// проверить на ошибки !!!
	err := r.db.Where("con_travel_time_id = ?", TTid).Find(&ArmiesIdsFromConArmyTT).Error
	if err != nil {
		return nil, err
	}

	// имеем массив, теперь надо из него извлечь все армии, принадлежащие данной заявке
	armiesRes := make([]ds.Army, 0)
	for _, curCATT := range ArmiesIdsFromConArmyTT {
		curArmy, err := r.GetArmy(curCATT.ConArmyID)
		if err != nil {
			return nil, err
		}
		armiesRes = append(armiesRes, curArmy)
	}

	return armiesRes, err
}

func (r *Repository) GetConnsByTTid(TTid int) ([]ds.ConnectArmyTT, error) {
	var ConnsATT []ds.ConnectArmyTT
	err := r.db.Where("con_travel_time_id = ?", TTid).Find(&ConnsATT).Error
	if err != nil {
		return nil, err
	}
	return ConnsATT, err
}

func (r *Repository) AddArmyToTT(idArmy int, idTT int) error {
	var CATT = ds.ConnectArmyTT{
		ConArmyID:       idArmy,
		ConTravelTimeID: idTT,
		KmPerDayIRL:     0,
	}
	logrus.Info("adding post army 3,1123")
	err := r.db.Create(&CATT).Error
	logrus.Info("adding post army 3,1234")
	if err != nil {
		fmt.Errorf("ошибка при добавлении армии в расчёт!: %w", err)
		return err
	}
	return nil
}

func (r *Repository) CheckIfArmyAlreadyInTT(idArmy int, idTT int) bool {
	var ConnsATT []ds.ConnectArmyTT
	err := r.db.Where("con_travel_time_id = ? AND con_army_id = ?", idTT, idArmy).Find(&ConnsATT).Error
	if err != nil {
		return false
	}
	return len(ConnsATT) != 0
}
