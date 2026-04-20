package chat

type UserMessage struct {
	ID      string   `json:"id"`
	Query   string   `json:"query"`
	History []string `json:"history"`
}
