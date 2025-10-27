package repository

import (
	"fmt"
	"time_of_armies/internal/app/ds"
)

func (r *Repository) GetHistorianByLogin(login string) (ds.Historian, error) {
	historian := ds.Historian{}
	err := r.db.Where("his_login = ?", login).First(&historian).Error
	if err != nil {
		return ds.Historian{}, err
	}
	return historian, nil
}

func (r *Repository) GetHistorianByID(HisID int) (ds.Historian, error) {
	historian := ds.Historian{}
	err := r.db.Where("historian_id = ?", HisID).First(&historian).Error
	if err != nil {
		return ds.Historian{}, err
	}
	return historian, nil
}

func (r *Repository) AuthHistorian(login string, password string) (ds.Historian, error) {
	historian := ds.Historian{}
	err := r.db.Where("his_login = ? AND his_password = ?", login, password).First(&historian).Error
	if err != nil {
		return ds.Historian{}, err
	}
	return historian, nil
}

func (r *Repository) UpdateHistorian(historian ds.Historian) (ds.Historian, error) {
	err := r.db.Model(&historian).Updates(ds.Historian{
		HisLogin:       historian.HisLogin,
		HisPassword:    historian.HisPassword,
		HisIsModerator: historian.HisIsModerator,
	}).Error
	if err != nil {
		fmt.Errorf("ошибка при обновлении пользователя!: %w", err)
		return ds.Historian{}, err
	}
	return historian, nil
}

func (r *Repository) RegisterHistorian(historian ds.Historian) (ds.Historian, error) {
	err := r.db.Create(&historian).Error
	if err != nil {
		fmt.Errorf("ошибка при добавлении нового пользователя в БД!: %w", err)
		return ds.Historian{}, err
	}
	return historian, nil
}
