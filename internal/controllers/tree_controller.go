package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/geno-tree/go-back/internal/services"
)

type TreeController struct {
	treeService *services.TreeService
}

func NewTreeController(treeService *services.TreeService) *TreeController {
	return &TreeController{
		treeService: treeService,
	}
}

// GetFamilyTree получает полное семейное дерево от главной персоны
func (tc *TreeController) GetFamilyTree(c *gin.Context) {
	personIDStr := c.Param("personId")
	personID, err := strconv.ParseUint(personIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID персоны"})
		return
	}

	userID := c.GetUint("user_id")

	tree, err := tc.treeService.GetFamilyTree(uint(personID), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка построения дерева: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tree": tree})
}

// GetTreeStatistics получает статистику семейного дерева
func (tc *TreeController) GetTreeStatistics(c *gin.Context) {
	userID := c.GetUint("user_id")

	statistics, err := tc.treeService.GetTreeStatistics(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения статистики: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statistics": statistics})
}

// GetMyFamilyTree получает семейное дерево текущего пользователя
func (tc *TreeController) GetMyFamilyTree(c *gin.Context) {
	userID := c.GetUint("user_id")

	// Получаем всех персон пользователя
	persons, err := tc.treeService.GetPersons(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения персон: " + err.Error()})
		return
	}

	if len(persons) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "У вас пока нет персон в дереве"})
		return
	}

	// Используем первую персону как корневую (можно добавить логику выбора главной персоны)
	rootPerson := persons[0]

	tree, err := tc.treeService.GetFamilyTree(rootPerson.ID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка построения дерева: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tree": tree})
}
