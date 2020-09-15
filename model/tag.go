package model

type Tag struct {
	DBFullModel

	Name string `json:"name" gorm:"unique" sql:"index"`
	Description string `json:"description" gorm:"type:text"`
}
