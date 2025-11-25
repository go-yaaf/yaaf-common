package utils

import (
	. "github.com/go-yaaf/yaaf-common/entity"
)

// region Factory method -----------------------------------------------------------------------------------------------

// JsonUtilsStruct provides utility functions for converting between structured entities and raw JSON maps.
// This can be useful for scenarios where you need to work with JSON data in a more dynamic, less-structured way.
type JsonUtilsStruct struct {
}

// JsonUtils is a factory method that returns a new instance of JsonUtilsStruct.
func JsonUtils() *JsonUtilsStruct {
	return &JsonUtilsStruct{}
}

// endregion

// region Public methods -----------------------------------------------------------------------------------------------

// ToJson converts a given entity (which should be a struct or a pointer to a struct) into a raw JSON map (map[string]any).
// This is done by first marshalling the entity to a JSON byte slice and then unmarshalling it into the map.
//
// Parameters:
//
//	entity: The entity to convert. It should be marshallable to JSON.
//
// Returns:
//
//	A map[string]any representing the entity as raw JSON.
//	An error if marshalling or unmarshalling fails.
func (t *JsonUtilsStruct) ToJson(entity any) (raw map[string]any, err error) {
	bytes, err := Marshal(entity)
	if err != nil {
		return nil, err
	}

	raw = make(map[string]any)
	if err = Unmarshal(bytes, &raw); err != nil {
		return nil, err
	}
	return raw, nil
}

// FromJson converts a raw JSON map (map[string]any) into a structured entity.
// This is achieved by first marshalling the map to a JSON byte slice and then unmarshalling it into a new entity
// created by the provided factory function.
//
// Parameters:
//
//	factory: A function that returns a new instance of the target entity (e.g., `func() Entity { return &MyEntity{} }`).
//	raw: The raw JSON map to convert.
//
// Returns:
//
//	The populated entity.
//	An error if marshalling or unmarshalling fails.
func (t *JsonUtilsStruct) FromJson(factory EntityFactory, raw map[string]any) (entity Entity, err error) {
	bytes, err := Marshal(raw)
	if err != nil {
		return nil, err
	}

	entity = factory()
	if err = Unmarshal(bytes, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

// endregion
