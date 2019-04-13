package recipes_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/ashkarin/ashkarin-api-test/pkg/recipe"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	baseUrl               = "http://localhost:9090"
	timeout               = 4 * time.Second
	recipeID              = ""
	averageRating float64 = 4
	ratingsCount  int64   = 3
	score         int8    = 3
)
var loc *time.Location

func CreateHTTPRequest(method, url string, body interface{}) *http.Request {
	var req *http.Request
	if body != nil {
		req, _ = http.NewRequest(method, url, strings.NewReader(body.(string)))
	} else {
		req, _ = http.NewRequest(method, url, nil)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return req
}

var _ = Describe("RecipesService", func() {
	It("should get time location properly", func() {
		var err error
		loc, err = time.LoadLocation("UTC")
		if err != nil {
			GinkgoWriter.Write([]byte(err.Error()))
			Expect(err).NotTo(HaveOccurred())
		}
	})

	It("should return alive when requesting the root", func() {
		req := CreateHTTPRequest("GET", baseUrl+"/", nil)
		client := &http.Client{Timeout: time.Duration(timeout)}
		res, err := client.Do(req)

		if err != nil {
			GinkgoWriter.Write([]byte(err.Error()))
			Expect(err).NotTo(HaveOccurred())
		} else {
			Expect(res.StatusCode).To(Equal(200))
		}
	})

	It("should create a recipe", func() {
		timestamp := time.Date(2019, 04, 11, 9, 00, 20, 0, loc)
		expected := recipe.Recipe{
			Name:          "My Recipe",
			PrepTime:      timestamp,
			Difficulty:    recipe.Easy,
			Vegetarian:    true,
			AverageRating: averageRating,
			RatingsCount:  ratingsCount,
		}

		obtained := recipe.Recipe{}
		params, _ := json.Marshal(expected)

		req := CreateHTTPRequest("POST", baseUrl+"/recipes", string(params))
		client := &http.Client{Timeout: time.Duration(timeout)}
		res, err := client.Do(req)

		if err != nil {
			GinkgoWriter.Write([]byte(err.Error()))
			Expect(err).NotTo(HaveOccurred())
		} else {
			Expect(res.StatusCode).To(Equal(http.StatusCreated))
			body, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(body, &obtained)
			Expect(expected.Name).To(Equal(obtained.Name))
			Expect(expected.PrepTime).To(Equal(obtained.PrepTime))
			Expect(int(expected.Difficulty)).To(Equal(int(obtained.Difficulty)))
			Expect(expected.Vegetarian).To(Equal(obtained.Vegetarian))
			Expect(expected.AverageRating).To(Equal(obtained.AverageRating))
			Expect(expected.RatingsCount).To(Equal(obtained.RatingsCount))
		}
	})

	It("should obtain a recipe", func() {
		req := CreateHTTPRequest("GET", baseUrl+"/recipes/0/5", nil)
		client := &http.Client{Timeout: time.Duration(timeout)}
		res, err := client.Do(req)
		obtained := []recipe.Recipe{}

		if err != nil {
			GinkgoWriter.Write([]byte(err.Error()))
			Expect(err).NotTo(HaveOccurred())
		} else {
			Expect(res.StatusCode).To(Equal(http.StatusOK))
			body, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(body, &obtained)
			recipeID = obtained[0].ID.(string)
			Expect(obtained).To(HaveLen(1))
		}
	})

	It("should rate a recipe", func() {
		req := CreateHTTPRequest("POST", baseUrl+"/recipes/"+fmt.Sprintf("%s/rate/%d", recipeID, score), nil)
		client := &http.Client{Timeout: time.Duration(timeout)}
		res, err := client.Do(req)

		if err != nil {
			GinkgoWriter.Write([]byte(err.Error()))
		} else {
			Expect(res.StatusCode).To(Equal(http.StatusOK))
		}
	})

	It("should return a single recipe by id", func() {
		req := CreateHTTPRequest("GET", baseUrl+"/recipes/"+recipeID, nil)
		client := &http.Client{Timeout: time.Duration(timeout)}
		res, err := client.Do(req)
		obtained := recipe.Recipe{}

		if err != nil {
			GinkgoWriter.Write([]byte(err.Error()))
			Expect(err).NotTo(HaveOccurred())
		} else {
			Expect(res.StatusCode).To(Equal(http.StatusOK))
			body, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(body, &obtained)
			Expect(obtained.ID.(string)).To(Equal(recipeID))

			expectedRating := averageRating + (float64(score)-averageRating)/float64(ratingsCount+1)
			Expect(obtained.RatingsCount).To(Equal(ratingsCount + 1))
			Expect(obtained.AverageRating).To(Equal(expectedRating))
		}
	})

	It("should update a single recipe by id", func() {
		timestamp := time.Date(2019, 04, 11, 9, 00, 20, 0, loc)
		expected := recipe.Recipe{
			Name:          "Updated",
			PrepTime:      timestamp,
			Difficulty:    1,
			Vegetarian:    true,
			AverageRating: averageRating,
			RatingsCount:  ratingsCount,
		}

		params, _ := json.Marshal(expected)
		obtained := recipe.Recipe{}

		req := CreateHTTPRequest("PUT", baseUrl+"/recipes/"+recipeID, string(params))
		client := &http.Client{Timeout: time.Duration(timeout)}
		res, err := client.Do(req)

		if err != nil {
			GinkgoWriter.Write([]byte(err.Error()))
		} else {
			Expect(res.StatusCode).To(Equal(http.StatusOK))
			body, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(body, &obtained)
			Expect(expected.Name).To(Equal(obtained.Name))
			Expect(expected.PrepTime).To(Equal(obtained.PrepTime))
			Expect(int(expected.Difficulty)).To(Equal(int(obtained.Difficulty)))
			Expect(expected.Vegetarian).To(Equal(obtained.Vegetarian))
			Expect(expected.AverageRating).To(Equal(obtained.AverageRating))
			Expect(expected.RatingsCount).To(Equal(obtained.RatingsCount))
		}
	})

	It("should search through recipes", func() {
		req := CreateHTTPRequest("GET", baseUrl+"/recipes/search/ated", nil)
		client := &http.Client{Timeout: time.Duration(timeout)}
		res, err := client.Do(req)

		obtained := []recipe.Recipe{}

		if err != nil {
			GinkgoWriter.Write([]byte(err.Error()))
		} else {
			Expect(res.StatusCode).To(Equal(http.StatusOK))
			body, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(body, &obtained)
			Expect(len(obtained)).Should(BeNumerically(">=", 1))
		}
	})

	It("should delete a recipe", func() {
		req := CreateHTTPRequest("DELETE", baseUrl+"/recipes/"+recipeID, nil)
		client := &http.Client{Timeout: time.Duration(timeout)}
		res, err := client.Do(req)

		if err != nil {
			GinkgoWriter.Write([]byte(err.Error()))
		} else {
			Expect(res.StatusCode).To(Equal(http.StatusOK))
		}
	})
})
