package services

import (
	"errors"

	"gitlab.com/geno-tree/go-back/internal/models"
	"gitlab.com/geno-tree/go-back/internal/repositories"
)

type TreeService struct {
	personRepository       *repositories.PersonRepository
	relationshipRepository *repositories.RelationshipRepository
}

func NewTreeService(personRepository *repositories.PersonRepository, relationshipRepository *repositories.RelationshipRepository) *TreeService {
	return &TreeService{
		personRepository:       personRepository,
		relationshipRepository: relationshipRepository,
	}
}

// CRUD операции для персон
func (s *TreeService) CreatePerson(person *models.Person) (models.Person, error) {
	return s.personRepository.CreatePerson(person)
}

func (s *TreeService) UpdatePersonPartial(id uint, updateData map[string]interface{}) (models.Person, error) {
	person, err := s.personRepository.UpdatePerson(id, updateData)
	if err != nil {
		return models.Person{}, err
	}
	return *person, nil
}

func (s *TreeService) GetPerson(id uint) (*models.Person, error) {
	return s.personRepository.GetPerson(id)
}

func (s *TreeService) GetPersons(userID uint) ([]models.Person, error) {
	return s.personRepository.GetPersons(userID)
}

func (s *TreeService) DeletePerson(id uint) error {
	return s.personRepository.DeletePerson(id)
}

// CRUD операции для связей
func (s *TreeService) CreateRelationship(relationship *models.Relationship) (models.Relationship, error) {
	return s.relationshipRepository.CreateRelationship(relationship)
}

func (s *TreeService) UpdateRelationship(relationship *models.Relationship) (models.Relationship, error) {
	return s.relationshipRepository.UpdateRelationship(relationship)
}

func (s *TreeService) GetRelationships(userID uint) ([]models.Relationship, error) {
	return s.relationshipRepository.GetRelationships(userID)
}

func (s *TreeService) GetRelationshipsByPersonID(personID uint) ([]models.Relationship, error) {
	return s.relationshipRepository.GetRelationshipsByPersonID(personID)
}

func (s *TreeService) DeleteRelationship(id uint) error {
	return s.relationshipRepository.DeleteRelationship(id)
}

// SearchPersonsByName ищет персон по имени
func (s *TreeService) SearchPersonsByName(userID uint, query string) ([]models.Person, error) {
	return s.personRepository.SearchPersonsByName(userID, query)
}

// GetFamilyTree строит полное семейное дерево от главной персоны
func (s *TreeService) GetFamilyTree(mainPersonID, userID uint) (*models.FamilyTree, error) {
	// Получаем главную персону
	mainPerson, err := s.GetPerson(mainPersonID)
	if err != nil {
		return nil, err
	}

	// Проверяем права доступа
	if mainPerson.CreatedByUserID != userID {
		return nil, errors.New("unauthorized")
	}

	// Получаем всех персон пользователя
	allPersons, err := s.GetPersons(userID)
	if err != nil {
		return nil, err
	}

	// Получаем все связи пользователя
	allRelationships, err := s.GetRelationships(userID)
	if err != nil {
		return nil, err
	}

	// Строим карту персон
	personsMap := make(map[uint]*models.Person)
	for i := range allPersons {
		personsMap[allPersons[i].ID] = &allPersons[i]
	}

	// Строим карту связей с учетом всех типов
	connectionsMap := make(map[uint][]models.Connection)

	for _, rel := range allRelationships {
		// Определяем тип связи для Person1 -> Person2
		connType1 := s.getConnectionType(rel.Type, true)
		conn1 := models.Connection{
			TargetPersonID: rel.Person2ID,
			TargetPerson:   personsMap[rel.Person2ID],
			Relationship:   rel,
			Type:           connType1,
			Direction:      models.DirectionTo,
			IsDirect:       true,
		}
		connectionsMap[rel.Person1ID] = append(connectionsMap[rel.Person1ID], conn1)

		// Определяем тип связи для Person2 -> Person1
		connType2 := s.getConnectionType(rel.Type, false)
		conn2 := models.Connection{
			TargetPersonID: rel.Person1ID,
			TargetPerson:   personsMap[rel.Person1ID],
			Relationship:   rel,
			Type:           connType2,
			Direction:      models.DirectionFrom,
			IsDirect:       true,
		}
		connectionsMap[rel.Person2ID] = append(connectionsMap[rel.Person2ID], conn2)
	}

	// Добавляем косвенные связи (через других людей)
	s.addIndirectConnections(connectionsMap, personsMap)

	// Определяем поколения
	generations := s.buildGenerations(mainPersonID, personsMap, connectionsMap)

	// Строим статистику
	statistics := s.buildStatistics(allPersons, allRelationships, generations)

	tree := &models.FamilyTree{
		RootPerson:    mainPerson,
		Persons:       personsMap,
		Relationships: allRelationships,
		Connections:   connectionsMap,
		Generations:   generations,
		Statistics:    statistics,
	}

	return tree, nil
}

// GetTreeStatistics получает статистику семейного дерева
func (s *TreeService) GetTreeStatistics(userID uint) (*models.TreeStatistics, error) {
	persons, err := s.GetPersons(userID)
	if err != nil {
		return nil, err
	}

	relationships, err := s.GetRelationships(userID)
	if err != nil {
		return nil, err
	}

	// Строим временную карту поколений для статистики
	generations := make(map[int][]*models.Person)

	stats := &models.TreeStatistics{
		TotalPersons:       len(persons),
		TotalRelationships: len(relationships),
		Generations:        len(generations),
		RelationshipStats:  make(map[models.RelationshipType]int),
	}

	// Подсчитываем живых и умерших
	for _, person := range persons {
		if person.IsDeceased() {
			stats.DeceasedPersons++
		} else {
			stats.LivingPersons++
		}
	}

	// Подсчитываем статистику по типам связей
	for _, rel := range relationships {
		stats.RelationshipStats[rel.Type]++
	}

	return stats, nil
}

// getConnectionType определяет тип связи с учетом всех ваших типов
func (s *TreeService) getConnectionType(relType models.RelationshipType, isFromPerson1 bool) models.ConnectionType {
	switch relType {
	case models.RelationshipParent:
		if isFromPerson1 {
			return models.ConnectionParent
		}
		return models.ConnectionChild

	case models.RelationshipSpouse:
		return models.ConnectionSpouse

	case models.RelationshipSibling:
		return models.ConnectionSibling

	case models.RelationshipGrandparent:
		if isFromPerson1 {
			return models.ConnectionGrandparent
		}
		return models.ConnectionGrandchild

	case models.RelationshipGrandchild:
		if isFromPerson1 {
			return models.ConnectionGrandchild
		}
		return models.ConnectionGrandparent

	case models.RelationshipUncle:
		if isFromPerson1 {
			return models.ConnectionUncle
		}
		return models.ConnectionCousin // Для племянника

	case models.RelationshipAunt:
		if isFromPerson1 {
			return models.ConnectionAunt
		}
		return models.ConnectionCousin // Для племянника

	case models.RelationshipCousin:
		return models.ConnectionCousin

	case models.RelationshipMarriage:
		return models.ConnectionMarriage

	case models.RelationshipDivorce:
		return models.ConnectionDivorce

	case models.RelationshipEngagement:
		return models.ConnectionEngagement

	case models.RelationshipPartnership:
		return models.ConnectionPartnership

	case models.RelationshipFriend:
		return models.ConnectionFriend

	case models.RelationshipColleague:
		return models.ConnectionColleague

	case models.RelationshipNeighbor:
		return models.ConnectionNeighbor

	default:
		return models.ConnectionOther
	}
}

func (s *TreeService) addIndirectConnections(connections map[uint][]models.Connection, persons map[uint]*models.Person) {
	// Находим дядей и тетей через родителей и их братьев/сестер
	for personID, personConns := range connections {
		for _, conn := range personConns {
			if conn.Type == models.ConnectionParent {
				// Ищем братьев/сестер родителя
				parentConns := connections[conn.TargetPersonID]
				for _, parentConn := range parentConns {
					if parentConn.Type == models.ConnectionSibling {
						// Определяем пол родителя для правильного типа связи
						parent := persons[conn.TargetPersonID]
						var uncleAuntType models.ConnectionType
						if parent.Gender == models.GenderMale {
							uncleAuntType = models.ConnectionUncle
						} else {
							uncleAuntType = models.ConnectionAunt
						}

						// Добавляем связь дядя/тетя
						uncleConn := models.Connection{
							TargetPersonID: parentConn.TargetPersonID,
							TargetPerson:   parentConn.TargetPerson,
							Relationship:   parentConn.Relationship,
							Type:           uncleAuntType,
							Direction:      models.DirectionTo,
							IsDirect:       false,
						}
						connections[personID] = append(connections[personID], uncleConn)
					}
				}
			}
		}
	}
}

// buildGenerations строит карту поколений для семейного дерева
func (s *TreeService) buildGenerations(
	rootID uint,
	persons map[uint]*models.Person,
	connections map[uint][]models.Connection,
) map[int][]*models.Person {
	generations := make(map[int][]*models.Person)
	visited := make(map[uint]bool)

	type queueItem struct {
		personID   uint
		generation int
	}

	queue := []queueItem{{personID: rootID, generation: 0}}
	visited[rootID] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		person := persons[current.personID]
		generations[current.generation] = append(generations[current.generation], person)

		for _, conn := range connections[current.personID] {
			if !visited[conn.TargetPersonID] {
				visited[conn.TargetPersonID] = true
				nextGen := current.generation

				// Определяем поколение на основе типа связи
				switch conn.Type {
				case models.ConnectionParent:
					nextGen = current.generation - 1
				case models.ConnectionChild:
					nextGen = current.generation + 1
				case models.ConnectionSpouse, models.ConnectionSibling:
					// Остаются в том же поколении
				}
				queue = append(queue, queueItem{personID: conn.TargetPersonID, generation: nextGen})
			}
		}
	}

	return generations
}

// buildStatistics строит статистику по дереву
func (s *TreeService) buildStatistics(
	persons []models.Person,
	relationships []models.Relationship,
	generations map[int][]*models.Person,
) *models.TreeStatistics {
	stats := &models.TreeStatistics{
		TotalPersons:       len(persons),
		TotalRelationships: len(relationships),
		Generations:        len(generations),
		RelationshipStats:  make(map[models.RelationshipType]int),
	}

	maxDepth := 0
	for gen := range generations {
		if gen > maxDepth {
			maxDepth = gen
		}
	}
	stats.MaxDepth = maxDepth

	for _, person := range persons {
		if person.IsDeceased() {
			stats.DeceasedPersons++
		} else {
			stats.LivingPersons++
		}
	}

	for _, rel := range relationships {
		stats.RelationshipStats[rel.Type]++
	}

	return stats
}
