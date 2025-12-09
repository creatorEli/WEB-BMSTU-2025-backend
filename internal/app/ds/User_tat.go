package ds

type User_tat struct {
	UserID_tat      uint   `gorm:"primary_key;not null" json:"id"`
	Login_tat       string `gorm:"type:varchar(25);unique;not null" json:"login"`
	Password_tat    string `gorm:"type:varchar(100);not null" json:"-"`
	IsModerator_tat bool   `gorm:"type:boolean;default:false" json:"is_moderator"`
}
