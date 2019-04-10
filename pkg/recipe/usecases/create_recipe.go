package usecases

import (
	"github.com/ashkarin/ashkarin-api-test/pkg/recipe"
)

// CreateRecipe create the recipe entry in the storage
func CreateRecipe(s recipe.StorageGateway, r *recipe.Recipe) error {
	return s.Store(r)
}
