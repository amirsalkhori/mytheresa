package model

import "mytheresa/internal/domain"

type Category struct {
	ID   uint32 `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"index;type:varchar(255);not null"`
}

func (Category) TableName() string {
	return "categories"
}

func ToModelCategory(domainCategory domain.Category) Category {
	return Category{
		ID:   domainCategory.ID,
		Name: domainCategory.Name,
	}
}

func ToDomainCategory(modelCategory Category) domain.Category {
	return domain.Category{
		ID:   modelCategory.ID,
		Name: modelCategory.Name,
	}
}
