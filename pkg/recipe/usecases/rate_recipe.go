package usecases

import (
	"fmt"

	"github.com/ashkarin/ashkarin-api-test/pkg/recipe"
	log "github.com/sirupsen/logrus"
)

// RateRecipe rate the recipe by giving it a score from 1 to 5
func RateRecipe(s recipe.StorageGateway, r *recipe.Recipe, score uint8) error {
	return RateRecipeByID(s, r.ID.(string), score)
}

// RateRecipeByID rate the recipe by giving it a score from 1 to 5
func RateRecipeByID(s recipe.StorageGateway, id string, score uint8) error {
	if score < 1 || score > 5 {
		return fmt.Errorf("Recipe can be rated from 1 to 5")
	}
	// NOTE: We could rely on other usecases, but it will make dependencies

	// 1. Begin a transaction

	// 2. Get the recipe entry
	r, err := s.GetByID(id)
	if err != nil {
		return fmt.Errorf("Error in getting the recipe: %v", err)
	}

	// 3. Recompute average
	log.Infof("Before rate: %v", r)
	r.RatingsCount = r.RatingsCount + 1
	r.AverageRating = r.AverageRating + (float64(score)-r.AverageRating)/float64(r.RatingsCount)

	log.Infof("After rate: %v", r)

	// 4. Update entry in the storage
	if err := s.Update(r); err != nil {
		return fmt.Errorf("Error in updating the recipe: %v", err)
	}

	// 5. Close the transaction

	return nil
}
