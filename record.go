package sfdc

import (
	"encoding/json"
	"errors"
	"strings"
)

const (
	// RecordAttributes is the attribute map from the record JSON
	RecordAttributes = "attributes"
	recordAttrType   = "type"
	recordAttrURL    = "url"
)

// Record is a representation of a Salesforce
// record.
type Record struct {
	sobject string
	url     string
	fields  map[string]interface{}
}

// RecordFromJSONMap creates a recrod from a JSON map.
func RecordFromJSONMap(jsonMap map[string]interface{}) (*Record, error) {
	if jsonMap == nil {
		return nil, errors.New("record: map can not be nil")
	}
	r := &Record{}
	r.fromJSONMap(jsonMap)
	return r, nil
}

// UnmarshalJSON provides a custom unmarshaling of a
// JSON byte array.
func (r *Record) UnmarshalJSON(data []byte) error {
	if r == nil {
		return errors.New("record: can't unmarshal to a nil struct")
	}
	if data == nil {
		return errors.New("record: can't unmarshal to a nil byte array")
	}
	var jsonMap map[string]interface{}
	err := json.Unmarshal(data, &jsonMap)
	if err != nil {
		return err
	}

	r.fromJSONMap(jsonMap)
	return nil
}

func (r *Record) fromJSONMap(jsonMap map[string]interface{}) {
	r.fields = make(map[string]interface{})

	for k, v := range jsonMap {
		if k == RecordAttributes {
			if attr, ok := v.(map[string]interface{}); ok {
				if obj, ok := attr[recordAttrType]; ok {
					if sobj, ok := obj.(string); ok {
						r.sobject = sobj
					}
				}
				if obj, ok := attr[recordAttrURL]; ok {
					if url, ok := obj.(string); ok {
						r.url = url
					}
				}
			}
		} else {
			if v != nil {
				if m, is := v.(map[string]interface{}); is {
					if _, has := m["records"]; has == false {
						// Is not a nested record set so should be added to fields object
						r.fields[k] = m
					}
				} else {
					r.fields[k] = v
				}
			}
		}
	}
}

// SObject returns attribute's Salesforce object name.
func (r *Record) SObject() string {
	return r.sobject
}

// URL returns the record attribute's URL.
func (r *Record) URL() string {
	return r.url
}

// FieldValue returns the field's value.  If there is no field
// for the field name, then false will be returned.
func (r *Record) FieldValue(field string) (interface{}, bool) {
	fieldChain := strings.Split(field, ".")
	var value interface{}
	var has bool
	value = r.fields
	for _, f := range fieldChain {
		if fieldMap, ok := value.(map[string]interface{}); ok {
			value, has = fieldMap[f]
			if has == false {
				return value, has
			}
		} else {
			return nil, false
		}
	}
	return value, has
}

// Fields returns the map of field name to value relationships.
func (r *Record) Fields() map[string]interface{} {

	fields := make(map[string]interface{})
	for k, v := range r.fields {
		fields[k] = v
	}
	return fields
}
