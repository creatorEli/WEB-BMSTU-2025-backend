package ds

type ConnectArmyTT struct { // проверить на ошибки создания БД!!!
	catt_ID int `gorm:"primaryKey"`

	// здесь создаем Unique key, указывая общий uniqueIndex
	ConArmyID       int `gorm:"not null;uniqueIndex:idx_contt_army"`
	ConTravelTimeID int `gorm:"not null;uniqueIndex:idx_contt_army"`

	Army       Army       `gorm:"foreignKey:ConArmyID"`
	TravelTime TravelTime `gorm:"foreignKey:ConTravelTimeID"`

	KmPerDayIRL int
}
