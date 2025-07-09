package models

import (
	"time"

	"gorm.io/gorm"
)

type GenderType string

const (
	GenderMale   GenderType = "male"
	GenderFemale GenderType = "female"
	GenderOther  GenderType = "other"
)

type Person struct {
	gorm.Model

	FirstName  string `json:"first_name" gorm:"not null"`
	LastName   string `json:"last_name" gorm:"not null"`
	MiddleName string `json:"middle_name"`
	MaidenName string `json:"maiden_name"` // Девичья фамилия для женщин

	Gender     GenderType `json:"gender" gorm:"type:varchar(10);check:gender IN ('male','female','other')"`
	BirthDate  *time.Time `json:"birth_date"`
	BirthPlace string     `json:"birth_place"`
	DeathDate  *time.Time `json:"death_date"`
	DeathPlace string     `json:"death_place"`
	IsAlive    bool       `json:"is_alive" gorm:"default:true"`

	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`

	Biography string `json:"biography" gorm:"type:text"`
	Notes     string `json:"notes" gorm:"type:text"`

	PhotoURL string `json:"photo_url"`

	CreatedByUserID uint `json:"created_by_user_id"` // ID пользователя, создавшего запись

	// Связи
	RelationshipsAsPerson1 []Relationship `json:"relationships_as_person1" gorm:"foreignKey:Person1ID"`
	RelationshipsAsPerson2 []Relationship `json:"relationships_as_person2" gorm:"foreignKey:Person2ID"`
}

func (p *Person) GetFullName() string {
	if p.MiddleName != "" {
		return p.FirstName + " " + p.MiddleName + " " + p.LastName
	}
	return p.FirstName + " " + p.LastName
}

func (p *Person) GetAge() *int {
	if p.BirthDate == nil {
		return nil
	}

	var endDate time.Time
	if p.DeathDate != nil {
		endDate = *p.DeathDate
	} else {
		endDate = time.Now()
	}

	age := int(endDate.Sub(*p.BirthDate).Hours() / 24 / 365.25)
	return &age
}

// IsDeceased проверяет, умер ли человек
func (p *Person) IsDeceased() bool {
	return p.DeathDate != nil || !p.IsAlive
}
