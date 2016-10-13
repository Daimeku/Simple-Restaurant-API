package models

import ()

type ResourceListResponse struct {
	Data []FormattedResource `json:data`
}
