// Package telegram provides s client for the Telegram Bot API.
package telegram

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"the-mitya-nagan-bot/lib/e"
	"time"
)

// с клиентом работаем через указатель чтобы не копировать объект
// Clients communicates with the Telegram Bot API.
type Client struct {
	host     string      // хост API сервиса телеграмма
	basePath string      // префикс с которого начинаются все запросы
	client   http.Client // чтобы не созхдавать для каждого запроса отдельно
}

// клиент знает:
// - куда отправлять запросы (host);
// - какой токен использовать (basePath);
// - каким HTTP-клиентом пользоваться.

// tg-bot.com/bot<token>

// New cerates a new Telegram Bot API client.
func New(host, token string) (*Client, error) {
	proxyURL, err := url.Parse("http://127.0.0.1:10801")
	if err != nil {
		return nil, err
	}

	httpClient := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
		Timeout: 30 * time.Second,
	}

	return &Client{
		host:     host,
		basePath: "bot" + token,
		client:   httpClient, // объект, который умеет общаться по HTTP
	}, nil
}

// offset - с какого обновления начинать выдачу, начиная с какого update_id получать обновления
// limit - сколько максимум обновлений вернуть

// формируем запрос в API тг на получение сообщений
// Updates retrieves updates from Telegarm.
func (c *Client) Updates(offset, limit int) ([]Update, error) {
	q := url.Values{}                     // структура для хранения GET-параметров URL, контейнер для GET-параметров запроса
	q.Add("offset", strconv.Itoa(offset)) // .Add добавляет параметр в будущий URL-запрос
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest("getUpdates", q)
	if err != nil {
		return nil, e.Wrap("updates Updates fail:", err)
	}

	var res UpdatesResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, e.Wrap("updates Updates fail:", err)
	}

	return res.Result, err
}

// формируем запрос в API тг на отправку сообщений
// SendMessage sends a text message to the specified chat.
func (c *Client) SendMessage(chatId int, text string) error {
	q := url.Values{} // подготавливаем параметры запроса
	q.Add("chat_id", strconv.Itoa(chatId))
	q.Add("text", text)

	_, err := c.doRequest("sendMessage", q)
	if err != nil {
		return e.Wrap("Can't send message", err)
	}
	return nil
}

// doRequest делает HTTP-запрос
func (c *Client) doRequest(method string, query url.Values) ([]byte, error) {
	// собираем адрес по которому сделаем запрос
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	// "Создай объект HTTP-запроса типа GET по адресу Telegram API. Тело запроса не нужно."
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, e.Wrap("updates doRequest fail:", err)
	}

	req.URL.RawQuery = query.Encode() // Добавляем параметры в запрос

	resp, err := c.client.Do(req) // запрос уходит в Telegram
	if err != nil {
		return nil, e.Wrap("updates doRequest fail:", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, e.Wrap("updates doRequest fail:", err)
	}

	return body, nil
}
