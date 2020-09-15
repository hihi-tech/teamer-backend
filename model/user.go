package model

// User describes a user object
type User struct {
	DBSafeModel

	Email    string          `json:"email" gorm:"unique" sql:"index"`
	Password string          `json:"-"`
	Phone    *string         `json:"phone"`

	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Birthday  string `json:"birthday"`

	Schools []*School `json:"schools" gorm:"many2many:users_schools"`
	Tags    []*Tag    `json:"tags" gorm:"many2many:users_tags"`
}
