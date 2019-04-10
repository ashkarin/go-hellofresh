package usecases

import (
	"github.com/ashkarin/ashkarin-api-test/pkg/recipe"
)

// GetRecipe get recipe by its ID
func GetRecipe(s recipe.StorageGateway, ID string) (*recipe.Recipe, error) {
	return s.GetByID(ID)
}
