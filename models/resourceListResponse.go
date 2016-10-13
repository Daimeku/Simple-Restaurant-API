package models

import ()

type ResourceListResponse struct {
	Data []ResourceResponse `json:data`
}
