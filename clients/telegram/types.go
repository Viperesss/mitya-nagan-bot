package telegram

type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	UpdateId int              `json:"update_id"`
	Message  *IncomingMessage `json:"message"`
	// указатель здесь используют как признак:
	// Есть объект? -> *IncomingMessage
	// Нет объекта? -> nil
}

type IncomingMessage struct {
	Text string `json:"text"`
	From From   `json:"from"`
	Chat Chat   `json:"chat"`
}

type From struct {
	Username string `json:"username"`
}

type Chat struct {
	ID int `json:"id"`
}
