package usecases

import (
	"github.com/ashkarin/ashkarin-api-test/pkg/recipe"
)

// ListRecipes list all recipes in the storage
func ListRecipes(s recipe.StorageGateway, start, limit uint64) ([]*recipe.Recipe, error) {
	return s.GetRange(start, limit)
}
