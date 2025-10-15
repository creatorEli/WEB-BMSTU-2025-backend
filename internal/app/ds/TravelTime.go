package ds

import "time"

type TravelTime struct {
	TtID int `gorm:"primaryKey;not null"`
	// 5 статусов: черновик, удален, сформирован, завершен, отклонен
	StatusTT       string    `gorm:"type:varchar(15);not null" json:"StatusTT"`
	DateCreateTT   time.Time `gorm:"not null"`
	DateUpdateTT   time.Time
	DateFinishTT   time.Time
	CreatorID_TT   uint  `gorm:"not null"`
	ModeratorID_TT *uint `gorm:"default:null"`

	CreatorTT   Historian `gorm:"foreignKey:CreatorID_TT"`
	ModeratorTT Historian `gorm:"foreignKey:ModeratorID_TT;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	//"Desert", "Plain", "Mount", "Forest", "River"
	ChosenBiomTT    string `gorm:"type:varchar(20);" json:"ChosenBiomTT"`
	DistanceTT      int    `json:"DistanceTT"`
	ResultMinTT     int    `json:"ResultMinTT"`
	ResultMaxTT     int    `json:"ResultMaxTT"`
	ResultChronical int    `json:"ResultChronical"`
}
