package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/geno-tree/go-back/internal/models"
	"gitlab.com/geno-tree/go-back/internal/services"
)

type PersonController struct {
	treeService *services.TreeService
}

func NewPersonController(treeService *services.TreeService) *PersonController {
	return &PersonController{
		treeService: treeService,
	}
}

// CreatePerson создает новую персону
func (pc *PersonController) CreatePerson(c *gin.Context) {
	var person models.Person

	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные: " + err.Error()})
		return
	}

	// Получаем ID пользователя из контекста
	userID := c.GetUint("user_id")
	person.CreatedByUserID = userID

	createdPerson, err := pc.treeService.CreatePerson(&person)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания персоны: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"person": createdPerson})
}

// GetPerson получает персону по ID
func (pc *PersonController) GetPerson(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID персоны"})
		return
	}

	person, err := pc.treeService.GetPerson(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Персона не найдена"})
		return
	}

	// Проверяем права доступа
	userID := c.GetUint("user_id")
	if person.CreatedByUserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Нет доступа к этой персоне"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"person": person})
}

// GetPersons получает всех персон пользователя
func (pc *PersonController) GetPersons(c *gin.Context) {
	userID := c.GetUint("user_id")

	persons, err := pc.treeService.GetPersons(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения персон: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"persons": persons})
}

// UpdatePerson обновляет персону
func (pc *PersonController) UpdatePerson(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID персоны"})
		return
	}

	// Получаем существующую персону для проверки прав доступа
	existingPerson, err := pc.treeService.GetPerson(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Персона не найдена"})
		return
	}

	// Проверяем права доступа
	userID := c.GetUint("user_id")
	if existingPerson.CreatedByUserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Нет доступа к этой персоне"})
		return
	}

	// Привязываем JSON в map для частичного обновления
	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные: " + err.Error()})
		return
	}

	// Удаляем системные поля, которые не должны обновляться
	delete(updateData, "id")
	delete(updateData, "created_at")
	delete(updateData, "updated_at")
	delete(updateData, "deleted_at")
	delete(updateData, "created_by_user_id")

	// Обновляем только переданные поля
	updatedPerson, err := pc.treeService.UpdatePersonPartial(uint(id), updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления персоны: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"person": updatedPerson})
}

// DeletePerson удаляет персону
func (pc *PersonController) DeletePerson(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID персоны"})
		return
	}

	// Проверяем права доступа
	person, err := pc.treeService.GetPerson(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Персона не найдена"})
		return
	}

	userID := c.GetUint("user_id")
	if person.CreatedByUserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Нет доступа к этой персоне"})
		return
	}

	err = pc.treeService.DeletePerson(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления персоны: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Персона успешно удалена"})
}

// SearchPersons ищет персон по имени
func (pc *PersonController) SearchPersons(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Параметр поиска 'q' обязателен"})
		return
	}

	userID := c.GetUint("user_id")
	persons, err := pc.treeService.SearchPersonsByName(userID, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка поиска: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"persons": persons})
}
