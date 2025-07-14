package models

// FamilyTree представляет полное семейное дерево
type FamilyTree struct {
	RootPerson    *Person               `json:"root_person"`
	Persons       map[uint]*Person      `json:"persons"`       // Все персоны
	Relationships []Relationship        `json:"relationships"` // Все связи
	Connections   map[uint][]Connection `json:"connections"`   // Связи по ID персоны
	Generations   map[int][]*Person     `json:"generations"`   // Персоны по поколениям
	Statistics    *TreeStatistics       `json:"statistics"`
}

// Connection представляет связь между персонами
type Connection struct {
	TargetPersonID uint                `json:"target_person_id"`
	TargetPerson   *Person             `json:"target_person"`
	Relationship   Relationship        `json:"relationship"`
	Type           ConnectionType      `json:"type"`
	Direction      ConnectionDirection `json:"direction"`
	IsDirect       bool                `json:"is_direct"` // Прямая связь или через другого человека
}

type ConnectionType string

const (
	// Прямые семейные связи
	ConnectionParent      ConnectionType = "parent"      // Родитель
	ConnectionChild       ConnectionType = "child"       // Ребенок
	ConnectionSpouse      ConnectionType = "spouse"      // Супруг/супруга
	ConnectionSibling     ConnectionType = "sibling"     // Брат/сестра
	ConnectionGrandparent ConnectionType = "grandparent" // Дедушка/бабушка
	ConnectionGrandchild  ConnectionType = "grandchild"  // Внук/внучка
	ConnectionUncle       ConnectionType = "uncle"       // Дядя
	ConnectionAunt        ConnectionType = "aunt"        // Тетя
	ConnectionCousin      ConnectionType = "cousin"      // Двоюродный брат/сестра

	// Романтические связи
	ConnectionMarriage    ConnectionType = "marriage"    // Брак
	ConnectionDivorce     ConnectionType = "divorce"     // Развод
	ConnectionEngagement  ConnectionType = "engagement"  // Помолвка
	ConnectionPartnership ConnectionType = "partnership" // Партнерство

	// Другие связи
	ConnectionFriend    ConnectionType = "friend"    // Друг
	ConnectionColleague ConnectionType = "colleague" // Коллега
	ConnectionNeighbor  ConnectionType = "neighbor"  // Сосед
	ConnectionOther     ConnectionType = "other"     // Другое
)

type ConnectionDirection string

const (
	DirectionTo   ConnectionDirection = "to"   // От текущей персоны к целевой
	DirectionFrom ConnectionDirection = "from" // От целевой персоны к текущей
	DirectionBoth ConnectionDirection = "both" // Двусторонняя связь
)

// TreeStatistics статистика дерева
type TreeStatistics struct {
	TotalPersons       int `json:"total_persons"`
	TotalRelationships int `json:"total_relationships"`
	LivingPersons      int `json:"living_persons"`
	DeceasedPersons    int `json:"deceased_persons"`
	Generations        int `json:"generations"`
	MaxDepth           int `json:"max_depth"`

	// Детальная статистика по типам связей
	RelationshipStats map[RelationshipType]int `json:"relationship_stats"`
}
