package model

type Rating struct {
	ID         int    `gorm:"column:id;primary_key" json:"id"`
	Player     int    `gorm:"column:player" json:"player"`
	Position   int    `gorm:"column:rating" json:"position"`
	DateUpdate string `gorm:"column:dateUpdate" json:"dateUpdate"`
	Points     int    `gorm:"column:points" json:"points"`
}

// TableName sets the insert table name for this struct type
func (p *Rating) TableName() string {
	return "rating"
}
