package ds

type Army struct { // поменять по услугам!
	ArmyID          int    `gorm:"primaryKey" json:"ArmyID"`
	NameArmy        string `gorm:"type:varchar(150);not null" json:"NameArmy"`
	DescriptionArmy string `gorm:"type:text;default:null" json:"DescriptionArmy"`
	// status: удален/действует
	StatusArmy   string `gorm:"type:varchar(15);not null" json:"StatusArmy"`
	ImageArmyUrl string `gorm:"type:varchar(100);default:null" json:"ImageArmyUrl"`
	ClassArmy    string `gorm:"type:varchar(30);not null" json:"ClassArmy"` // для фильтрации
	//скорости для каждого типа местности:
	MinPlainSpeed  int `json:"MinPlainSpeed"`  // Равнина минимальная
	MaxPlainSpeed  int `json:"MaxPlainSpeed"`  // Равнина Максимальная
	MinMountSpeed  int `json:"MinMountSpeed"`  // Горы минимальная
	MaxMountSpeed  int `json:"MaxMountSpeed"`  // Горы Максимальная
	MinForestSpeed int `json:"MinForestSpeed"` // Лес минимальная
	MaxForestSpeed int `json:"MaxForestSpeed"` // Лес Максимальная
	MinRiverSpeed  int `json:"MinRiverSpeed"`  // Река минимальная
	MaxRiverSpeed  int `json:"MaxRiverSpeed"`  // Река Максимальная
	MinDesertSpeed int `json:"MinDesertSpeed"` // Пустыня минимальная
	MaxDesertSpeed int `json:"MaxDesertSpeed"` // Пустыня Максимальная
}
