package telegram

import (
	"encoding/json"
	"io"
	"m/lib/e"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

const (
	getUpdatesMethod  = "getUpdatesMethod"
	sendMessageMethod = "sendMessage"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

func New(host string, token string) Client {
	return Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) Updates(offset, limit int) ([]Update, error) {
	q := url.Values{}

	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, err
	}

	var res UpdatesResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (c *Client) SendMessage(chatId int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatId))
	q.Add("text", text)

	_, err := c.doRequest(sendMessageMethod, q)
	if err != nil {
		return e.Wrap("can't send message", err)
	}

	return nil
}

func (c *Client) doRequest(method string, query url.Values) (data []byte, err error) {
	const errMsg = "can not do request"
	defer func() { err = e.WrapIfErr(errMsg, err) }()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}