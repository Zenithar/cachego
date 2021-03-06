package cachego

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Mongo struct {
	collection *mgo.Collection
}

type MongoContent struct {
	Duration int64
	Key      string `bson:"_id"`
	Value    string
}

func (m *Mongo) Contains(key string) bool {
	if _, err := m.Fetch(key); err != nil {
		return false
	}

	return true
}

func (m *Mongo) Delete(key string) error {
	return m.collection.Remove(bson.M{"_id": key})
}

func (m *Mongo) Fetch(key string) (string, error) {
	content := &MongoContent{}

	err := m.collection.Find(bson.M{"_id": key}).One(content)

	if err != nil {
		return "", err
	}

	if content.Duration == 0 {
		return content.Value, nil
	}

	if content.Duration <= time.Now().Unix() {
		m.Delete(key)
		return "", errors.New("Cache expired")
	}

	return content.Value, nil
}

func (m *Mongo) FetchMulti(keys []string) map[string]string {
	result := make(map[string]string)

	iter := m.collection.Find(bson.M{"_id": bson.M{"$in": keys}}).Iter()

	content := &MongoContent{}

	for iter.Next(content) {
		result[content.Key] = content.Value
	}

	return result
}

func (m *Mongo) Flush() error {
	_, err := m.collection.RemoveAll(bson.M{})

	return err
}

func (m *Mongo) Save(key string, value string, lifeTime time.Duration) error {
	duration := int64(0)

	if lifeTime > 0 {
		duration = time.Now().Unix() + int64(lifeTime.Seconds())
	}

	content := &MongoContent{
		duration,
		key,
		value,
	}

	if _, err := m.collection.Upsert(bson.M{"_id": key}, content); err != nil {
		return err
	}

	return nil
}
