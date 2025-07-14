package repositories

import (
	"gitlab.com/geno-tree/go-back/internal/database"
	"gitlab.com/geno-tree/go-back/internal/models"
)

type RelationshipRepository struct {
	db *database.Db
}

func NewRelationshipRepository(db *database.Db) *RelationshipRepository {
	return &RelationshipRepository{db: db}
}

func (repo *RelationshipRepository) CreateRelationship(relationship *models.Relationship) (models.Relationship, error) {
	result := repo.db.Create(relationship)

	if result.Error != nil {
		return models.Relationship{}, result.Error
	}

	return *relationship, nil
}

func (repo *RelationshipRepository) UpdateRelationship(relationship *models.Relationship) (models.Relationship, error) {
	result := repo.db.Save(relationship)

	if result.Error != nil {
		return models.Relationship{}, result.Error
	}

	return *relationship, nil
}

func (repo *RelationshipRepository) GetRelationships(userID uint) ([]models.Relationship, error) {
	var relationships []models.Relationship

	result := repo.db.Where("created_by_user_id = ?", userID).Find(&relationships)

	if result.Error != nil {
		return nil, result.Error
	}

	return relationships, nil
}

func (repo *RelationshipRepository) GetRelationshipsByPersonID(personID uint) ([]models.Relationship, error) {
	var relationships []models.Relationship

	result := repo.db.Where("person1_id = ? OR person2_id = ?", personID, personID).Find(&relationships)

	if result.Error != nil {
		return nil, result.Error
	}

	return relationships, nil
}

func (repo *RelationshipRepository) DeleteRelationship(id uint) error {
	result := repo.db.Delete(&models.Relationship{}, id)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
