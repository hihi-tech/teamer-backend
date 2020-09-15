package model

type School struct {
	DBFullModel

	Name string `json:"name" gorm:"unique" sql:"index"`
	//AlternateName string `json:"alternateName"`

	// Location is the administration assignment of the current object. Example: `北京` for `北京大学`
	Location string `json:"location" sql:"index"`
	//Description   string `json:"description" gorm:"type:text"`
}
