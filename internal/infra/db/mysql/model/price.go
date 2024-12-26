package model

type Price struct {
	ID                 int32  `gorm:"primaryKey;autoIncrement"`
	Original           uint32 `gorm:"not null"`
	Final              uint32 `gorm:"not null"`
	DiscountPercentage uint8  `gorm:"not null"`
	Currency           string `gorm:"index;type:varchar(10);not null"`
}

func (Price) TableName() string {
	return "prices"
}
