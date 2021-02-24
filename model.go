package goovn

import (
	"fmt"
	"reflect"
)

type TableName = string

// A Model is the base interface used to build Database Models. It is used
// to express how data from a specific Database Table shall be translated into structs
// A Model is a struct with at least one (most likely more) field tagged with the 'ovs' tag
// The value of 'ovs' field must be a valid column name in the OVS Database
// A field with the 'ovs' tag value '_uuid' is mandatory. The rest of the columns are optional
// The struct may also have non-tagged fields (which will be ignored by the API calls)
// The Model interface must be implemented by the pointer to such type
// Example:
//type MyLogicalRouter struct {
//	UUID          string            `ovs:"_uuid"`
//	Name          string            `ovs:"name"`
//	ExternalIDs   map[string]string `ovs:"external_ids"`
//	LoadBalancers []string          `ovs:"load_balancer"`
//}
//
//func (lr *MyLogicalRouter) Table() TableName {
//	return "Logical_Router"
//}
type Model interface {
	// Table returns the name of the Table this model represents
	Table() TableName
}

// BaseModel is a base structure that can be embedded to build models
type BaseModel struct {
	UUID string `ovs:"_uuid"`
}

// DBModel is a Database model
type DBModel struct {
	types map[TableName]reflect.Type
}

// newModel returns a new instance of a model from a specific TableName
func (db DBModel) newModel(table TableName) (Model, error) {
	mtype, ok := db.types[table]
	if !ok {
		return nil, fmt.Errorf("Table %s not found in Database Model", string(table))
	}
	model := reflect.New(mtype.Elem())
	return model.Interface().(Model), nil
}

// GetTypes returns the DBModel Types
func (db DBModel) GetTypes() map[TableName]reflect.Type {
	return db.types
}

// GetType returns the DBModel Types
func NewDBModel(models []Model) (*DBModel, error) {
	types := make(map[TableName]reflect.Type, len(models))
	for _, model := range models {
		if reflect.TypeOf(model).Kind() != reflect.Ptr {
			return nil, fmt.Errorf("Model is expected to be a pointer")
		}
		uField := reflect.Indirect(reflect.ValueOf(model)).FieldByName("UUID")
		if !uField.IsValid() || uField.Type().Kind() != reflect.String {
			return nil, fmt.Errorf("Model is expected to have a string field called UUID")

		}
		types[model.Table()] = reflect.TypeOf(model)
	}
	return &DBModel{
		types: types,
	}, nil
}
