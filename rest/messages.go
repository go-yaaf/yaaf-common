// Common REST messages
//

package rest

import . "github.com/go-yaaf/yaaf-common/entity"

// region BaseRestResponse ---------------------------------------------------------------------------------------------

// BaseRestResponse is a common structure for all response types
// @Data
type BaseRestResponse struct {
	Code  int    `json:"code"`            // Error code (0 for success)
	Error string `json:"error,omitempty"` // Error message
}

// SetError sets error message
func (res *BaseRestResponse) SetError(err error) {
	if err != nil {
		res.Code = -1
		res.Error = err.Error()
	}
}

// SetErrorWithCode sets an error message with customized code (must be negative)
func (res *BaseRestResponse) SetErrorWithCode(code int, err error) {
	if err != nil {
		res.Code = code
		res.Error = err.Error()
	}
}

// NewErrorResponse response with error
func NewErrorResponse(err error) (res *BaseRestResponse) {
	res = &BaseRestResponse{}
	res.SetError(err)
	return res
}

// NewErrorResponseWithCode response with error and customized code
func NewErrorResponseWithCode(code int, err error) (res *BaseRestResponse) {
	res = &BaseRestResponse{}
	res.SetErrorWithCode(code, err)
	return res
}

// endregion

// region ActionResponse -----------------------------------------------------------------------------------------------

// ActionResponse message is returned for any action on entity with no return data (e.d. delete)
// @Data
type ActionResponse struct {
	BaseRestResponse
	Key  string `json:"key,omitempty"`  // The entity key (Id)
	Data string `json:"data,omitempty"` // Additional data
}

// NewActionResponse factory method
func NewActionResponse(key, data string) (er *ActionResponse) {
	return &ActionResponse{Key: key, Data: data}
}

// endregion

// region EntityResponse -----------------------------------------------------------------------------------------------

// EntityResponse message is returned for any create/update action on entity
// @Data
type EntityResponse[T Entity] struct {
	BaseRestResponse
	Entity T `json:"entity"` // The entity
}

// NewEntityResponse factory method
func NewEntityResponse[T Entity](entity T) (er *EntityResponse[T]) {
	return &EntityResponse[T]{Entity: entity}
}

// endregion

// region EntitiesResponse ---------------------------------------------------------------------------------------------

// EntitiesResponse message is returned for any action returning multiple entities
// @Data
type EntitiesResponse[T Entity] struct {
	BaseRestResponse
	Page  int `json:"page"`  // Current page (Bulk) number
	Size  int `json:"size"`  // Size of page (items in bulk)
	Pages int `json:"pages"` // Total number of pages
	Total int `json:"total"` // Total number of items in the query
	List  []T `json:"list"`  // List of objects in the current result set
}

// NewEntitiesResponse factory method
func NewEntitiesResponse[T Entity](entities []T, page, size, total int) *EntitiesResponse[T] {

	if size == 0 {
		size = 1
	}
	rem := 0
	if total%size > 0 {
		rem = 1
	}
	return &EntitiesResponse[T]{
		Page:  page,
		Size:  size,
		Total: total,
		Pages: (total / size) + rem,
		List:  entities,
	}
}

// endregion

// region EntityResponse -----------------------------------------------------------------------------------------------

// EntityRequest message is returned for any create/update action on entity
// @Data
type EntityRequest[T Entity] struct {
	Entity T `json:"entity"` // The entity
}

// NewEntityRequest factory method
func NewEntityRequest[T Entity](entity T) (er *EntityRequest[T]) {
	return &EntityRequest[T]{Entity: entity}
}

// endregion

// region EntitiesRequest ---------------------------------------------------------------------------------------------

// EntitiesRequest message is returned for any action returning multiple entities
// @Data
type EntitiesRequest[T Entity] struct {
	List []T `json:"list"` // List of objects in the current result set
}

// NewEntitiesRequest factory method
func NewEntitiesRequest[T Entity](entities []T, page, size, total int) *EntitiesRequest[T] {
	return &EntitiesRequest[T]{
		List: entities,
	}
}

// endregion
