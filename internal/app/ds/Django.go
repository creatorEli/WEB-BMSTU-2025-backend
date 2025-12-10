package ds

// Структуры для запроса в Django
type DjangoRequest struct {
	TravelTimeID int        `json:"travel_time_id"`
	Distance     int        `json:"distance"`
	Biom         string     `json:"biom"`
	CreatorID    int        `json:"creator_id"`
	ModeratorID  *int       `json:"moderator_id,omitempty"`
	Armies       []ArmyData `json:"armies"`
}

type ArmyData struct {
	ArmyID            int    `json:"ArmyID"`
	NameArmy          string `json:"NameArmy"`
	ClassArmy         string `json:"ClassArmy"`
	MinPlainSpeed     int    `json:"MinPlainSpeed"`
	MaxPlainSpeed     int    `json:"MaxPlainSpeed"`
	MinDesertSpeed    int    `json:"MinDesertSpeed"`
	MaxDesertSpeed    int    `json:"MaxDesertSpeed"`
	MinRiverSpeed     int    `json:"MinRiverSpeed"`
	MaxRiverSpeed     int    `json:"MaxRiverSpeed"`
	MinForestSpeed    int    `json:"MinForestSpeed"`
	MaxForestSpeed    int    `json:"MaxForestSpeed"`
	MinMountSpeed     int    `json:"MinMountSpeed"`
	MaxMountSpeed     int    `json:"MaxMountSpeed"`
	KmPerDayChronical int    `json:"KmPerDayChronical"`
}

// Структура для ответа callback от Django
type DjangoCallback struct {
	TravelTimeID    int    `json:"travel_time_id"`
	Success         bool   `json:"success"`
	ResultMinTT     int    `json:"result_min_tt,omitempty"`
	ResultMaxTT     int    `json:"result_max_tt,omitempty"`
	ResultChronical int    `json:"result_chronical,omitempty"`
	ErrorMessage    string `json:"error_message,omitempty"`
}
