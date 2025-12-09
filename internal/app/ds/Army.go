package ds

type Army struct { // поменять по услугам!
	ArmyID          int    `gorm:"primaryKey"`
	NameArmy        string `gorm:"type:varchar(150);not null"`
	DescriptionArmy string `gorm:"type:text;default:null"`
	// status: удален/действует
	StatusArmy   string `gorm:"type:varchar(15);not null"`
	ImageArmyUrl string `gorm:"type:varchar(100);default:null"`
	ClassArmy    string `gorm:"type:varchar(30);not null"` // для фильтрации
	//скорости для каждого типа местности:
	MinPlainSpeed  int // Равнина минимальная
	MaxPlainSpeed  int // Равнина Максимальная
	MinMountSpeed  int // Горы минимальная
	MaxMountSpeed  int // Горы Максимальная
	MinForestSpeed int // Лес минимальная
	MaxForestSpeed int // Лес Максимальная
	MinRiverSpeed  int // Река минимальная
	MaxRiverSpeed  int // Река Максимальная
	MinDesertSpeed int // Пустыня минимальная
	MaxDesertSpeed int // Пустыня Максимальная
}
