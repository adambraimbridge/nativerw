package nativerw

import (
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"labix.org/v2/mgo/bson"
	"strings"
	"time"
)

type PropertyModifier interface {
	Remove()
	Update(interface{})
	Put(string, interface{})
	Name() string
	Value() interface{}
}

type simplePropertyModifier struct {
	values map[string]interface{}
	name   string
}

func (spm simplePropertyModifier) Remove() {
	delete(spm.values, spm.name)
}

func (spm simplePropertyModifier) Update(newValue interface{}) {
	spm.values[spm.name] = newValue
	//fmt.Printf("after update for %v the values are %v\n", spm.name, spm.values)
}

func (spm simplePropertyModifier) Put(name string, value interface{}) {
	spm.values[name] = value
	//fmt.Printf("after put for %v of %v=%v the values are %v\n", spm.name, name, value, spm.values)
}

func (spm simplePropertyModifier) Name() string {
	return spm.name
}

func (spm simplePropertyModifier) Value() interface{} {
	return spm.values[spm.name]
}

type propertyConverter func(pm PropertyModifier, context map[string]interface{}, name string, value interface{}) (stop bool)

type compositePropertyConverter struct {
	converters []propertyConverter
}

func (cpc compositePropertyConverter) convert(pm PropertyModifier, context map[string]interface{}, name string, value interface{}) (stop bool) {
	for _, c := range cpc.converters {
		stop := c(pm, context, name, pm.Value())
		if stop {
			return true
		}
	}
	return false
}

func UUIDToBson(pm PropertyModifier, context map[string]interface{}, name string, value interface{}) (stop bool) {
	if name == "uuid" {
		bsonUuid := bson.Binary{Kind: 0x04, Data: []byte(uuid.Parse(value.(string)))}
		//TODO: check it parsed ok. This might not really be a uuid
		pm.Update(bsonUuid)
	}
	return false
}

func UUIDFromBson(pm PropertyModifier, context map[string]interface{}, name string, value interface{}) (stop bool) {
	if name == "uuid" {
		if bin, ok := value.(bson.Binary); ok && bin.Kind == 0x04 {
			pm.Update(uuid.UUID(bin.Data).String())
		}
	}
	return false
}

func DateToBson(pm PropertyModifier, context map[string]interface{}, name string, value interface{}) (changed bool) {
	if strings.HasSuffix(name, "Date") {
		if s, ok := value.(string); ok {
			t, err := time.Parse(time.RFC3339, s)
			if err == nil {
				pm.Update(t)
				return true
			}
		}
	}
	return false
}

func DateFromBson(pm PropertyModifier, context map[string]interface{}, name string, value interface{}) (changed bool) {
	if strings.HasSuffix(name, "Date") {
		d, ok := value.(time.Time)
		if ok {
			pm.Update(d.Format(time.RFC3339))
			return true
		}
	}
	return false
}

func MongoIdRemover(pm PropertyModifier, context map[string]interface{}, name string, value interface{}) (changed bool) {
	if name == "_id" {
		pm.Remove()
		return true
	}
	return false
}

func ApiUrlInserter(pm PropertyModifier, context map[string]interface{}, name string, value interface{}) (changed bool) {
	if name == "uuid" {
		pm.Put("apiUrl", fmt.Sprintf("http://localhost:8082/%s/%s", "content" /*FIXME*/, value))
	}
	return false
}
