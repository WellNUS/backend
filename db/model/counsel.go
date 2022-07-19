package model

import (
	"time"
)

type CounselRequest struct {
	UserID 		int64		`json:"user_id"`
	Details 	string		`json:"details"`
	Topics		[]string	`json:"topics"`
	LastUpdated	time.Time 	`json:"last_updated"`
}

func (cr CounselRequest) HasTopic(topic string) bool {
	for _, t := range cr.Topics {
		if t == topic {
			return true
		}
	}
	return false
}