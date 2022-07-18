package model

import (
	"time"
)

type CounselRequest struct {
	UserID 		int64		`json:"user_id"`
	Details 	string		`json:"details"`
	Topics		[]string	`json:"topics"`
	TimeAdded	time.Time 	`json:"time_added"`
}

func (cr CounselRequest) HasTopic(topic string) bool {
	for _, t := range cr.Topics {
		if t == topic {
			return true
		}
	}
	return false
}