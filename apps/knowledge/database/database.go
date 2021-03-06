package database

import (
	"errors"

	"github.com/cstdev/knowledge-hub/apps/knowledge/types"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	log "github.com/sirupsen/logrus"
)

// Database provides an interface to define the methods required of a database
// connector.
type Database interface {
	Create(r types.Record) (string, error)
	Search(query types.SearchQuery) ([]types.Record, error)
	Update(id string, r types.Record) error
	Delete(id string) error
	Fields() ([]types.Field, error)
	UpdateFields(fields []types.Field) error
	DeleteField(id string) error
}

// MongoDB provides access and methods to talk to Mongo
type MongoDB struct {
	URL             string
	Database        string
	Collection      string
	FieldCollection string
}

// Create takes a record and writes it to the Mongo database
func (db *MongoDB) Create(r types.Record) (string, error) {
	session, err := GetSession(db.URL)
	defer session.Close()
	if err != nil {
		return "", err
	}

	if len(r.Location.Coordinates) == 0 {
		r.Location.Coordinates = []float64{r.Location.Lng, r.Location.Lat}
		r.Location.Type = "Point"
		log.WithFields(log.Fields{
			"coordinates": r.Location.Coordinates,
		}).Debug("Added coordinates for easy searching")
	}

	//Use DB from URL
	c := session.DB("").C(db.Collection)

	id := bson.NewObjectId()
	r.ID = id.Hex()
	_, err = c.UpsertId(id, r)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Failed to instert into database.")
		return "", err
	}

	log.WithFields(log.Fields{
		"_id": r.ID,
	}).Debug("Inserted Id")

	return r.ID, nil
}

func boundsPresent(query types.SearchQuery) bool {
	log.Debug("Checking bounds exist")
	return query.MinLat != 0 && query.MaxLat != 0 && query.MinLng != 0 && query.MaxLng != 0
}

// Search takes a query and returns all records that match
func (db *MongoDB) Search(query types.SearchQuery) ([]types.Record, error) {
	log.WithFields(log.Fields{
		"query":  query.Query,
		"minLat": query.MinLat,
		"maxLat": query.MaxLat,
		"minLng": query.MinLng,
		"maxLng": query.MaxLng,
	}).Debug("Searching DB")
	session, err := GetSession(db.URL)
	defer session.Close()
	if err != nil {
		return nil, err
	}
	c := session.DB("").C(db.Collection)

	var records []types.Record

	log.WithFields(log.Fields{
		"boundsPresent": boundsPresent(query),
	}).Debug("Checking for bounds")

	if boundsPresent(query) {
		err = c.Find(bson.M{"$and": []interface{}{bson.M{
			"location.coordinates": bson.M{
				"$geoWithin": bson.M{
					"$box": []interface{}{
						[]interface{}{query.MinLng, query.MinLat},
						[]interface{}{query.MaxLng, query.MaxLat},
					},
				},
			},
		}, bson.M{"deleted": bson.M{"$ne": true}}}}).All(&records)

		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error("Failed to get any records from the database.")
			return nil, err
		}

	}
	return records, nil
}

// Update takes and id of the record to update and a Record object
// containing any changes and writes it to Mongo
func (db *MongoDB) Update(id string, r types.Record) error {
	if r.ID != id {
		log.WithFields(log.Fields{
			"pathID":   id,
			"recordID": r.ID,
		}).Debug("Record ID does not match URL Path ID")
		return errors.New("Record ID does not match URL Path ID")
	}

	r.Location.Coordinates = []float64{r.Location.Lng, r.Location.Lat}

	session, err := GetSession(db.URL)
	defer session.Close()
	if err != nil {
		return err
	}

	c := session.DB("").C(db.Collection)

	err = c.Update(bson.M{"id": id}, r)
	if err != nil {
		log.WithFields(log.Fields{
			"id":    id,
			"error": err.Error(),
		}).Error("Failed to update database.")
		return err
	}

	return nil
}

// Delete marks the matching record in the database as deleted
func (db *MongoDB) Delete(id string) error {

	log.WithFields(log.Fields{
		"recordID": id,
	}).Debug("Marking record as deleted.")

	session, err := GetSession(db.URL)
	defer session.Close()
	if err != nil {
		return err
	}

	c := session.DB("").C(db.Collection)

	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"deleted": true}})
	if err != nil {
		log.WithFields(log.Fields{
			"id":    id,
			"error": err.Error(),
		}).Error("Failed to mark as deleted in database.")
		return err
	}
	return nil
}

// Fields retrieves all the set fields that can be use for entering information
func (db *MongoDB) Fields() ([]types.Field, error) {
	session, err := GetSession(db.URL)
	defer session.Close()
	if err != nil {
		return nil, err
	}

	c := session.DB("").C(db.FieldCollection)

	var fields []types.Field

	err = c.Find(bson.M{"deleted": bson.M{"$ne": true}}).All(&fields)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Failed to get the fields from the database.")
		return nil, err
	}

	return fields, nil
}

// UpdateFields writes the fields objects to the database, updating any that already exist.
func (db *MongoDB) UpdateFields(fields []types.Field) error {
	session, err := GetSession(db.URL)
	defer session.Close()
	if err != nil {
		return err
	}

	c := session.DB("").C(db.FieldCollection)

	bulk := c.Bulk()

	for _, field := range fields {
		bulk.Upsert(bson.M{"id": field.ID}, field)
	}

	_, err = bulk.Run()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Failed to update field in the database.")
		return err
	}

	return nil
}

//DeleteField takes an id of a field and marks it as deleted
func (db *MongoDB) DeleteField(id string) error {

	session, err := GetSession(db.URL)
	defer session.Close()
	if err != nil {
		return err
	}

	c := session.DB("").C(db.FieldCollection)

	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"deleted": true}})
	if err != nil {
		if err == mgo.ErrNotFound {
			log.WithFields(log.Fields{
				"id":    id,
				"error": err.Error(),
			}).Warn("Field with id doesn't exist.")
			return &types.FieldNotFoundError{id, "Field does not exist in the database."}
		}
		log.WithFields(log.Fields{
			"id":    id,
			"error": err.Error(),
		}).Error("Failed to delete field in the database.")
		return err
	}

	return nil
}

type FakeDB struct {
}

func (f *FakeDB) Create(r types.Record) (string, error) {
	return "", nil
}

func (f *FakeDB) Search(query types.SearchQuery) ([]types.Record, error) {
	return nil, nil
}

func (f *FakeDB) Update(id string, r types.Record) error {
	return nil
}

func (f *FakeDB) Delete(id string) error {
	return nil
}

func (f *FakeDB) Fields() ([]types.Field, error) {
	return nil, nil
}

func (f *FakeDB) UpdateFields(fields []types.Field) error {
	return nil
}
