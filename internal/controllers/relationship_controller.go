package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/geno-tree/go-back/internal/models"
	"gitlab.com/geno-tree/go-back/internal/services"
)

type RelationshipController struct {
	treeService *services.TreeService
}

func NewRelationshipController(treeService *services.TreeService) *RelationshipController {
	return &RelationshipController{
		treeService: treeService,
	}
}

// CreateRelationship создает новую связь
func (rc *RelationshipController) CreateRelationship(c *gin.Context) {
	var relationship models.Relationship

	if err := c.ShouldBindJSON(&relationship); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные: " + err.Error()})
		return
	}

	// Получаем ID пользователя из контекста
	userID := c.GetUint("user_id")
	relationship.CreatedByUserID = userID

	// Проверяем, что обе персоны принадлежат пользователю
	person1, err := rc.treeService.GetPerson(relationship.Person1ID)
	if err != nil || person1.CreatedByUserID != userID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Персона 1 не найдена или не принадлежит вам"})
		return
	}

	person2, err := rc.treeService.GetPerson(relationship.Person2ID)
	if err != nil || person2.CreatedByUserID != userID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Персона 2 не найдена или не принадлежит вам"})
		return
	}

	createdRelationship, err := rc.treeService.CreateRelationship(&relationship)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания связи: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"relationship": createdRelationship})
}

// GetRelationships получает все связи пользователя
func (rc *RelationshipController) GetRelationships(c *gin.Context) {
	userID := c.GetUint("user_id")

	relationships, err := rc.treeService.GetRelationships(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения связей: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"relationships": relationships})
}

// GetRelationshipsByPerson получает связи конкретной персоны
func (rc *RelationshipController) GetRelationshipsByPerson(c *gin.Context) {
	personIDStr := c.Param("id")
	personID, err := strconv.ParseUint(personIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID персоны"})
		return
	}

	userID := c.GetUint("user_id")

	// Проверяем, что персона принадлежит пользователю
	person, err := rc.treeService.GetPerson(uint(personID))
	if err != nil || person.CreatedByUserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Нет доступа к этой персоне"})
		return
	}

	relationships, err := rc.treeService.GetRelationshipsByPersonID(uint(personID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения связей: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"relationships": relationships})
}

// UpdateRelationship обновляет связь
func (rc *RelationshipController) UpdateRelationship(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID связи"})
		return
	}

	var relationship models.Relationship
	if err := c.ShouldBindJSON(&relationship); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные: " + err.Error()})
		return
	}

	relationship.ID = uint(id)
	userID := c.GetUint("user_id")
	relationship.CreatedByUserID = userID

	updatedRelationship, err := rc.treeService.UpdateRelationship(&relationship)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления связи: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"relationship": updatedRelationship})
}

// DeleteRelationship удаляет связь
func (rc *RelationshipController) DeleteRelationship(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID связи"})
		return
	}

	err = rc.treeService.DeleteRelationship(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления связи: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Связь успешно удалена"})
}
