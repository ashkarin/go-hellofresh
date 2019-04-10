package gateways

import (
	"github.com/ashkarin/ashkarin-api-test/pkg/recipe"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type mgoGateway struct {
	session    *mgo.Session
	db         *mgo.Database
	collection *mgo.Collection
}

// NewMongoDbGateway create a storage gateway to the MongoDB
func NewMongoDbGateway(server, port, username, password, database, collection string) (recipe.StorageGateway, error) {
	session, err := mgo.Dial(server)
	if err != nil {
		return nil, err
	}
	session.SetMode(mgo.Monotonic, true)
	db := session.DB(database)
	gw := &mgoGateway{
		session:    session,
		db:         db,
		collection: db.C(collection),
	}
	return gw, nil
}

func (s *mgoGateway) GetRange(start, limit uint64) ([]*recipe.Recipe, error) {
	var recipes []*recipe.Recipe
	err := s.collection.Find(nil).Skip(int(start)).Limit(int(limit)).All(&recipes)
	return recipes, err
}

func (s *mgoGateway) GetByID(id string) (*recipe.Recipe, error) {
	recipe := &recipe.Recipe{}
	query := bson.M{"_id": bson.ObjectIdHex(id)}
	if err := s.collection.Find(query).One(&recipe); err != nil {
		return nil, err
	}
	return recipe, nil
}

func (s *mgoGateway) DeleteByID(id string) error {
	query := bson.M{"_id": bson.ObjectIdHex(id)}
	return s.collection.Remove(query)
}

func (s *mgoGateway) Store(r *recipe.Recipe) error {
	return s.collection.Insert(r)
}

func (s *mgoGateway) Update(r *recipe.Recipe) error {
	var ID string
	switch v := r.ID.(type) {
	case string:
		ID = v
	case bson.ObjectId:
		ID = v.Hex()
	}

	query := bson.M{"_id": bson.ObjectIdHex(ID)}
	change := bson.M{"$set": bson.M{
		"name":          r.Name,
		"prepTime":      r.PrepTime,
		"difficulty":    r.Difficulty,
		"vegetarian":    r.Vegetarian,
		"averageRating": r.AverageRating,
		"ratingsCount":  r.RatingsCount,
	}}
	return s.collection.Update(query, change)
}

func (s *mgoGateway) Delete(r *recipe.Recipe) error {
	query := bson.M{"_id": bson.ObjectIdHex(r.ID.(string))}
	return s.collection.Remove(query)
}

func (s *mgoGateway) Search(pattern string) ([]*recipe.Recipe, error) {
	var recipes []*recipe.Recipe
	regex := bson.M{"$regex": bson.RegEx{Pattern: pattern}}
	if err := s.collection.Find(bson.M{"name": regex}).All(&recipes); err != nil {
		return nil, err
	}
	return recipes, nil
}
