package models

import ()

type ErrorResponse struct {
	Status  int    `json:status`
	Title   string `json:title`
	Code    string `json:code`
	Details string `json:details`
}
