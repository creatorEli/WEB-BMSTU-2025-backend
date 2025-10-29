package ds

import (
	"time_of_armies/internal/app/role"

	"github.com/google/uuid"
)

type Historian struct {
	HistorianUUID uuid.UUID `gorm:"type:uuid"`
	HistorianID   uint      `gorm:"primary_key;not null" json:"id"`
	HisLogin      string    `gorm:"type:varchar(25);unique;not null" json:"login"`
	HisPassword   string    `gorm:"type:varchar(100);not null" json:"-"`
	//HisIsModerator bool      `gorm:"type:boolean;default:false" json:"is_moderator"`
	Role role.Role `sql:"type:string;"`
}

// var Creator = Historian_tat{
// 	HistorianID_tat:      1,
// 	Login_tat:       "Eli",
// 	Password_tat:    "12345",
// 	IsModerator_tat: false,
// }

// var Moderator = Historian_tat{
// 	HistorianID_tat:      2,
// 	Login_tat:       "Anton",
// 	Password_tat:    "qwerty",
// 	IsModerator_tat: true,
// }

var CurrentHistorian Historian // немного костыльный способ для сохранения инфы об автроизованном прям сейчас пользователе
