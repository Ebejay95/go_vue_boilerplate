package models

type Notification struct {
	ID      int32  `json:"id"`
	Message string `json:"message"`
	Type    string `json:"type"`
}