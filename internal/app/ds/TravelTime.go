package ds

import "time"

type TravelTime struct {
	TtID int `gorm:"primaryKey;not null"`
	//Armies     []int // надо будет удалить, переброска в M-M
	// 5 статусов: черновик, удален, сформирован, завершен, отклонен
	StatusTT       string    `gorm:"type:varchar(15);not null"`
	DateCreateTT   time.Time `gorm:"not null"`
	DateUpdateTT   time.Time
	DateFinishTT   time.Time
	CreatorID_TT   uint `gorm:"not null"`
	ModeratorID_TT uint

	CreatorTT   User_tat `gorm:"foreignKey:CreatorID_TT"`
	ModeratorTT User_tat `gorm:"foreignKey:ModeratorID_TT"`

	ChosenBiomTT string `gorm:"type:varchar(20);"`
	DistanceTT   int
	ResultMinTT  int
	ResultMaxTT  int
}
