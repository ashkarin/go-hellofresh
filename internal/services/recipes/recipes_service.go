package recipes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/ashkarin/ashkarin-api-test/internal/utils"
	"github.com/ashkarin/ashkarin-api-test/pkg/recipe"
	"github.com/ashkarin/ashkarin-api-test/pkg/recipe/usecases"
)

// Service provides a set of HTTP handlers for work with recipes
type Service struct {
	storage recipe.StorageGateway
	router  *mux.Router
}

// NewService creates a service to work with recipes
func NewService(s recipe.StorageGateway, router *mux.Router) *Service {
	service := &Service{
		storage: s,
		router:  router,
	}
	service.initializeRoutes()

	return service
}

func (s *Service) initializeRoutes() {
	s.router.HandleFunc("/", s.IsAlive).Methods("GET")

	// POST [create recipe] ?/recipes
	s.router.HandleFunc("/recipes", s.СreateRecipe).Methods("POST")

	// GET [get recipe] ?/recipes/{id}
	s.router.HandleFunc("/recipes/{id}", s.GetRecipe).Methods("GET")

	// PUT [update recipe] ?/recipes/{id}
	s.router.HandleFunc("/recipes/{id}", s.UpdateRecipe).Methods("PUT")

	// DELETE [delete recipe] ?/recipes/{id}
	s.router.HandleFunc("/recipes/{id}", s.DeleteRecipe).Methods("DELETE")

	// GET [get recipes list] ?/recipes/{start:[0-9]+}/{limit:[0-9]+}
	s.router.HandleFunc("/recipes/{start:[0-9]+}/{limit:[0-9]+}", s.ListRecipes).Methods("GET")

	// POST [rate recipe] ?/recipes/{id}/rate/{score:[1-5]}
	s.router.HandleFunc("/recipes/{id}/rate/{score:[1-5]}", s.RateRecipe).Methods("POST")

	// GET [search recipes by name] ?/recipes/search/{name}
	s.router.HandleFunc("/recipes/search/{search:.+}", s.SearchRecipes).Methods("GET")
}

// IsAlive is the HTTP handler to check whether the service is alive or not
func (s *Service) IsAlive(w http.ResponseWriter, r *http.Request) {
	utils.ResponseWithJSON(w, http.StatusOK, "alive")
}

// СreateRecipe is the HTTP handler to create the recipe entry in the storage
func (s *Service) СreateRecipe(w http.ResponseWriter, r *http.Request) {
	var recipe recipe.Recipe
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&recipe); err != nil {
		log.Errorf("CreateRecipe: %v", err)
		utils.ResponseWithError(w, http.StatusBadRequest, "Invalid request payload JSON format")
		return
	}
	defer r.Body.Close()

	// Create the recipe in the storage
	if err := usecases.CreateRecipe(s.storage, &recipe); err != nil {
		log.Errorf("CreateRecipe: %v", err)
		utils.ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.ResponseWithJSON(w, http.StatusCreated, recipe)
}

// GetRecipe is the HTTP handler to get the recipe from storage by ID
func (s *Service) GetRecipe(w http.ResponseWriter, r *http.Request) {
	// Get the recipe ID
	vars := mux.Vars(r)
	id := vars["id"]

	// Get the recipe
	recipe, err := usecases.GetRecipe(s.storage, id)
	if err != nil {
		log.Errorf("GetRecipe: %v", err)
		utils.ResponseWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.ResponseWithJSON(w, http.StatusOK, recipe)
}

// UpdateRecipe is the HTTP handler to update the recipe entry in the storage
func (s *Service) UpdateRecipe(w http.ResponseWriter, r *http.Request) {
	// Get the recipe ID
	vars := mux.Vars(r)
	recipe := recipe.Recipe{ID: vars["id"]}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&recipe); err != nil {
		log.Errorf("UpdateRecipe: %v", err)
		utils.ResponseWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()

	// Update the recipe in the storage
	if err := usecases.UpdateRecipe(s.storage, &recipe); err != nil {
		log.Errorf("UpdateRecipe: %v", err)
		utils.ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.ResponseWithJSON(w, http.StatusOK, recipe)
}

// DeleteRecipe is the HTTP handler to delete the recipe from the storage
func (s *Service) DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	// Get the recipe ID
	vars := mux.Vars(r)
	id := vars["id"]

	// Delete the recipe
	if err := usecases.DeleteRecipeByID(s.storage, id); err != nil {
		log.Errorf("DeleteRecipe: %v", err)
		utils.ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.ResponseWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

// ListRecipes is the HTTP handler to list all the recipes in the storage
func (s *Service) ListRecipes(w http.ResponseWriter, r *http.Request) {
	// Get get the range of requested entries
	vars := mux.Vars(r)
	start, err := strconv.ParseUint(vars["start"], 10, 64)
	if err != nil {
		start = 0
	}

	limit, err := strconv.ParseUint(vars["limit"], 10, 64)
	if err != nil {
		limit = 10
	}

	// Call the related usecase
	recipes, err := usecases.ListRecipes(s.storage, start, limit)
	if err != nil {
		log.Errorf("ListRecipes: %v", err)
		utils.ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.ResponseWithJSON(w, http.StatusOK, recipes)
}

// RateRecipe is the HTTP handler to rate the recipe
func (s *Service) RateRecipe(w http.ResponseWriter, r *http.Request) {
	// Get the recipe ID and given score
	vars := mux.Vars(r)
	id := vars["id"]
	score, err := strconv.ParseUint(vars["score"], 10, 8)

	if err != nil {
		log.Errorf("RateRecipe: %v", err)
		utils.ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Rate recipe
	if err := usecases.RateRecipeByID(s.storage, id, uint8(score)); err != nil {
		log.Errorf("RateRecipe: %v", err)
		utils.ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.ResponseWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

// SearchRecipes is the HTTP handler to search the recipes
func (s *Service) SearchRecipes(w http.ResponseWriter, r *http.Request) {
	// Get the name pattern
	vars := mux.Vars(r)
	search := vars["search"]
	if search == "" {
		log.Errorf("SearchRecipes: No search pattern given")
		utils.ResponseWithError(w, http.StatusInternalServerError, "No search pattern given")
		return
	}

	// Search the recipes
	recipes, err := usecases.SearchRecipes(s.storage, search)
	if err != nil {
		log.Errorf("SearchRecipes: %v", err)
		utils.ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.ResponseWithJSON(w, http.StatusOK, recipes)
}
