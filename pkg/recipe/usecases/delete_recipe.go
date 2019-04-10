package usecases

import (
	"github.com/ashkarin/ashkarin-api-test/pkg/recipe"
)

// DeleteRecipe delete the recipe entry from the storage
func DeleteRecipe(s recipe.StorageGateway, r *recipe.Recipe) error {
	return s.Delete(r)
}

// DeleteRecipeByID delete the recipe entry from the storage
func DeleteRecipeByID(s recipe.StorageGateway, id string) error {
	return s.DeleteByID(id)
}
