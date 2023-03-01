// Common REST messages
//

package rest

import . "github.com/go-yaaf/yaaf-common/entity"

// region BaseRestResponse ---------------------------------------------------------------------------------------------

// BaseRestResponse is a common structure for all response types
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

// endregion

// region ActionResponse -----------------------------------------------------------------------------------------------

// ActionResponse message is returned for any action on entity with no return data (e.d. delete)
type ActionResponse struct {
	BaseRestResponse
	Key  string `json:"key,omitempty"`  // The entity key (Id)
	Data string `json:"data,omitempty"` // Additional data
}

// NewActionResponse factory method
func NewActionResponse(key, data string) (er *ActionResponse) {
	return &ActionResponse{Key: key, Data: data}
}

// NewActionResponseError with error
func NewActionResponseError(err error) (res *ActionResponse) {
	res = &ActionResponse{}
	res.SetError(err)
	return res
}

// endregion

// region EntityResponse -----------------------------------------------------------------------------------------------

// EntityResponse message is returned for any create/update action on entity
type EntityResponse struct {
	BaseRestResponse
	Entity Entity `json:"entity"` // The entity
}

// NewEntityResponse factory method
func NewEntityResponse(entity Entity) (er *EntityResponse) {
	return &EntityResponse{Entity: entity}
}

// NewEntityResponseError with error
func NewEntityResponseError(err error) (res *EntityResponse) {
	res = &EntityResponse{}
	res.SetError(err)
	return res
}

// endregion

// region EntityResponse -----------------------------------------------------------------------------------------------

// EntitiesResponse message is returned for any action returning multiple entities
type EntitiesResponse struct {
	BaseRestResponse
	Page  int      `json:"page"`  // Current page (Bulk) number
	Size  int      `json:"size"`  // Size of page (items in bulk)
	Pages int      `json:"pages"` // Total number of pages
	Total int      `json:"total"` // Total number of items in the query
	List  []Entity `json:"list"`  // List of objects in the current result set
}

// NewEntitiesResponse factory method
func NewEntitiesResponse(entities []Entity, page, size, total int) *EntitiesResponse {

	if size == 0 {
		size = 1
	}
	rem := 0
	if total%size > 0 {
		rem = 1
	}
	return &EntitiesResponse{
		Page:  page,
		Size:  size,
		Total: total,
		Pages: (total / size) + rem,
		List:  entities,
	}
}

// NewEntitiesResponseError with error
func NewEntitiesResponseError(err error) *EntitiesResponse {
	res := &EntitiesResponse{}
	res.SetError(err)
	return res
}

// endregion
