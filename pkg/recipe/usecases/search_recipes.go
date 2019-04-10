package usecases

import (
	"github.com/ashkarin/ashkarin-api-test/pkg/recipe"
)

// SearchRecipes get recipes which names match the given pattern
func SearchRecipes(s recipe.StorageGateway, namePattern string) ([]*recipe.Recipe, error) {
	return s.Search(namePattern)
}
