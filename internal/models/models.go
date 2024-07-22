package models

import "time"

type Message struct {
	ID         int32
	Content    string
	Processed  bool
	Created_at time.Time
}

type SendMessageResponse struct {
	StatusCode int32
	Message    string
}

type GetMessagesResponse struct {
	Messages []string
}

type GetStatsResponse struct {
	ProcessedCount int32
}
