package recipes_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/ashkarin/ashkarin-api-test/internal/services/recipes"
	"github.com/ashkarin/ashkarin-api-test/pkg/recipe/gateways"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var server *http.Server
var ctx context.Context

func TestRecipes(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Recipes Suite")
}

var _ = BeforeSuite(func() {
	// Open a gateway to the MongoDB storage
	dbhost := "mongodb"
	dbport := "27017"
	dbuser := ""
	dbpassword := ""
	dbname := "test_db"
	dbcollection := "testrecipes"
	ctx = context.Background()

	recipesStorage, err := gateways.NewMongoDbGateway(dbhost, dbport, dbuser, dbpassword, dbname, dbcollection)
	Expect(err).NotTo(HaveOccurred())

	go func() {
		// Create route and service
		router := mux.NewRouter()
		_ = recipes.NewService(recipesStorage, router)

		// Create the server
		server = &http.Server{
			Handler:      router,
			Addr:         fmt.Sprintf("%s:%s", "", "9090"),
			WriteTimeout: time.Duration(10) * time.Second,
			ReadTimeout:  time.Duration(10) * time.Second,
		}

		server.ListenAndServe()
	}()
})

var _ = AfterSuite(func() {
	server.Shutdown(ctx)
})
