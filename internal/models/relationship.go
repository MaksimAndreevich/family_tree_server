package models

import (
	"gorm.io/gorm"
)

type RelationshipType string

const (
	// Семейные связи
	RelationshipParent      RelationshipType = "parent"      // Родитель-ребенок
	RelationshipSpouse      RelationshipType = "spouse"      // Супруг/супруга
	RelationshipSibling     RelationshipType = "sibling"     // Брат/сестра
	RelationshipGrandparent RelationshipType = "grandparent" // Дедушка/бабушка
	RelationshipGrandchild  RelationshipType = "grandchild"  // Внук/внучка
	RelationshipUncle       RelationshipType = "uncle"       // Дядя
	RelationshipAunt        RelationshipType = "aunt"        // Тетя
	RelationshipCousin      RelationshipType = "cousin"      // Двоюродный брат/сестра

	// Романтические связи
	RelationshipMarriage    RelationshipType = "marriage"    // Брак
	RelationshipDivorce     RelationshipType = "divorce"     // Развод
	RelationshipEngagement  RelationshipType = "engagement"  // Помолвка
	RelationshipPartnership RelationshipType = "partnership" // Партнерство

	// Другие связи
	RelationshipFriend    RelationshipType = "friend"    // Друг
	RelationshipColleague RelationshipType = "colleague" // Коллега
	RelationshipNeighbor  RelationshipType = "neighbor"  // Сосед
	RelationshipOther     RelationshipType = "other"     // Другое
)

type Relationship struct {
	gorm.Model

	Person1ID uint `json:"person1_id" gorm:"not null;index"`
	Person2ID uint `json:"person2_id" gorm:"not null;index"`

	Type RelationshipType `json:"type" gorm:"not null;index"`

	Place string `json:"place"`

	Description string `json:"description" gorm:"type:text"`
	Notes       string `json:"notes" gorm:"type:text"`

	CreatedByUserID uint `json:"created_by_user_id"`

	// Связи с другими таблицами
	Person1 Person `json:"person1" gorm:"foreignKey:Person1ID"`
	Person2 Person `json:"person2" gorm:"foreignKey:Person2ID"`
	User    User   `json:"user" gorm:"foreignKey:CreatedByUserID"`
}
