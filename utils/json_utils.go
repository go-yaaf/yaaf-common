package utils

import (
	. "github.com/go-yaaf/yaaf-common/entity"
)

// region Factory method -----------------------------------------------------------------------------------------------

// General File utils
type jsonUtils struct {
}

// JsonUtils factory method
func JsonUtils() *jsonUtils {
	return &jsonUtils{}
}

// endregion

// region Public methods -----------------------------------------------------------------------------------------------

// ToJson convert entity to raw json (map of string keys into values)
func (t *jsonUtils) ToJson(entity Entity) (raw map[string]any, err error) {

	// Convert entity to string
	bytes, err := Marshal(entity)
	if err != nil {
		return nil, err
	}
	raw = make(map[string]any)

	// Convert string to arbitrary json
	err = Unmarshal(bytes, &raw)
	if err != nil {
		return nil, err
	}

	return raw, nil
}

// FromJson convert raw json to entity
func (t *jsonUtils) FromJson(factory EntityFactory, raw map[string]any) (entity Entity, err error) {

	// Convert raw to string
	bytes, err := Marshal(raw)
	if err != nil {
		return nil, err
	}

	entity = factory()

	// Convert string to entity
	err = Unmarshal(bytes, entity)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

// endregion
