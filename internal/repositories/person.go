package repositories

import (
	"time"

	"gitlab.com/geno-tree/go-back/internal/database"
	"gitlab.com/geno-tree/go-back/internal/models"
)

type PersonRepository struct {
	db *database.Db
}

func NewPersonRepository(db *database.Db) *PersonRepository {
	return &PersonRepository{db: db}
}

func (repo *PersonRepository) CreatePerson(person *models.Person) (models.Person, error) {
	result := repo.db.Create(person)

	if result.Error != nil {
		return models.Person{}, result.Error
	}

	return *person, nil
}

func (repo *PersonRepository) UpdatePerson(id uint, updateData map[string]interface{}) (*models.Person, error) {
	// Фильтруем пустые значения
	filteredUpdates := make(map[string]interface{})

	for key, value := range updateData {
		// Пропускаем nil значения
		if value == nil {
			continue
		}

		// Проверяем разные типы данных
		switch v := value.(type) {
		case string:
			if v != "" {
				filteredUpdates[key] = v
			}
		case *time.Time:
			if v != nil {
				filteredUpdates[key] = v
			}
		case bool:
			// Для bool всегда добавляем, так как false - валидное значение
			filteredUpdates[key] = v
		case int, int32, int64, uint, uint32, uint64:
			// Для чисел всегда добавляем, так как 0 - валидное значение
			filteredUpdates[key] = v
		default:
			// Для остальных типов добавляем как есть
			filteredUpdates[key] = v
		}
	}

	// Обновляем только непустые поля
	result := repo.db.Model(&models.Person{}).Where("id = ?", id).Updates(filteredUpdates)

	if result.Error != nil {
		return nil, result.Error
	}

	// Получаем обновленную персону
	var person models.Person
	if err := repo.db.Where("id = ?", id).First(&person).Error; err != nil {
		return nil, err
	}

	return &person, nil
}

func (repo *PersonRepository) GetPerson(id uint) (*models.Person, error) {
	var person models.Person

	result := repo.db.Where("id = ?", id).First(&person)

	if result.Error != nil {
		return nil, result.Error
	}

	return &person, nil

}

func (repo *PersonRepository) GetPersons(userID uint) ([]models.Person, error) {
	var persons []models.Person

	result := repo.db.Where("created_by_user_id = ?", userID).Find(&persons)

	if result.Error != nil {
		return nil, result.Error
	}

	return persons, nil

}

func (repo *PersonRepository) DeletePerson(id uint) error {
	result := repo.db.Delete(&models.Person{}, id)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// SearchPersonsByName ищет персон по имени
func (repo *PersonRepository) SearchPersonsByName(userID uint, query string) ([]models.Person, error) {
	var persons []models.Person

	result := repo.db.Where("created_by_user_id = ? AND (first_name ILIKE ? OR last_name ILIKE ? OR middle_name ILIKE ?)",
		userID, "%"+query+"%", "%"+query+"%", "%"+query+"%").Find(&persons)

	if result.Error != nil {
		return nil, result.Error
	}

	return persons, nil
}
