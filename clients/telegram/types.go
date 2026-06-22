// Package telegram contains Telegram Bot API models and client logic.
package telegram

// UpdateResponse represents a Telegram getUpdates response.
type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

// Update represents a Telegram update.
type Update struct {
	UpdateId int              `json:"update_id"`
	Message  *IncomingMessage `json:"message"`
	// указатель здесь используют как признак:
	// Есть объект? -> *IncomingMessage
	// Нет объекта? -> nil
}

// IncomingMessage repressents a Telegram message.
type IncomingMessage struct {
	Text string `json:"text"`
	From From   `json:"from"`
	Chat Chat   `json:"chat"`
}

// From contains imforamation about the message sender.
type From struct {
	Username string `json:"username"`
}

// From contains information about the chat.
type Chat struct {
	ID int `json:"id"`
}
