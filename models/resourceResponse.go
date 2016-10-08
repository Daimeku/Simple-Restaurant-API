package models

import ()

type FormattedResource struct {
	Type       string                 `json:type`
	Id         int                    `json:id`
	Attributes map[string]interface{} `json:attributes`
}

type ResourceResponse struct {
	Data [1]FormattedResource `json:data`
}
