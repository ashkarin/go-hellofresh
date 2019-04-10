package usecases

import (
	"github.com/ashkarin/ashkarin-api-test/pkg/recipe"
)

// UpdateRecipe update recipe entry in the storage
func UpdateRecipe(s recipe.StorageGateway, r *recipe.Recipe) error {
	return s.Update(r)
}
