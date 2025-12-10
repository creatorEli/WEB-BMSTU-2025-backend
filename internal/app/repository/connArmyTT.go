package repository

import (
	"time_of_armies/internal/app/ds"

	"github.com/sirupsen/logrus"
)

func (r *Repository) DeleteConn(idArmy int, TTid int) error {
	err := r.db.Where("con_army_id = ? AND con_travel_time_id = ?", idArmy, TTid).Delete(&ds.ConnectArmyTT{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateConn(idArmy int, TTid int, kmPerDayChrono int) error {
	logrus.Info("UpdateConn")
	err := r.db.Model(&ds.ConnectArmyTT{}).Where("con_army_id = ? AND con_travel_time_id = ?", idArmy, TTid).Update("km_per_day_chronical", kmPerDayChrono).Error
	if err != nil {
		return err
	}
	return nil
}
