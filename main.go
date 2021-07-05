package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	heroku "github.com/heroku/heroku-go/v5"
)

const userAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36"

type AuthTokenResponse struct {
	Token   string `json:"token"`
	Success bool   `json:"success"`
}

type SessionResponse struct {
	APIHost       string `json:"api_host"`
	Session       string `json:"session"`
	Success       bool   `json:"success"`
	UID           int    `json:"uid"`
	URL           string `json:"url"`
	WebsocketHost string `json:"websocket_host"`
	WebsocketPath string `json:"websocket_path"`
}

type AuthRequest struct {
	Cookie string `json:"cookie"`
	Method string `json:"_method"`
	ReqID  int    `json:"_reqid"`
}

type AuthResponse struct {
	Success bool `json:"success"`
}

func httpRequest(method, url string, form *url.Values, headers map[string]string) ([]byte, error) {
	var err error
	var req *http.Request
	if form != nil {
		req, err = http.NewRequest(method, url, strings.NewReader(form.Encode()))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(
			k,
			v,
		)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	ret, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func wsClient(host, path string) (*websocket.Conn, error) {
	var hs http.Header = make(http.Header)
	hs.Set("Host", host)
	hs.Set("Origin", "https://www.irccloud.com")
	hs.Set("User-Agent", userAgent)

	u := url.URL{Scheme: "wss", Host: host, Path: path}
	client, _, err := websocket.DefaultDialer.Dial(u.String(), hs)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func getAuthToken() (*AuthTokenResponse, error) {
	resp, err := httpRequest(
		"POST",
		"https://www.irccloud.com/chat/auth-formtoken",
		nil,
		map[string]string{
			"User-Agent": userAgent,
		},
	)
	if err != nil {
		return nil, err
	}

	var r AuthTokenResponse
	err = json.Unmarshal(resp, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func getSession(email, password, token string) (*SessionResponse, error) {
	form := url.Values{}
	form.Add("email", email)
	form.Add("password", password)
	form.Add("token", token)

	b, err := httpRequest(
		"POST",
		"https://www.irccloud.com/chat/login",
		&form,
		map[string]string{
			"Content-Type":     "application/x-www-form-urlencoded",
			"User-Agent":       userAgent,
			"X-Auth-FormToken": token,
		},
	)
	if err != nil {
		return nil, err
	}

	var r SessionResponse
	err = json.Unmarshal(b, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func authWebsocket(session, host, path string) (bool, error) {
	client, err := wsClient(host, path)
	if err != nil {
		return false, err
	}
	defer client.Close()

	areq, _ := json.Marshal(AuthRequest{session, "auth", 1})
	err = client.WriteMessage(websocket.TextMessage, areq)
	if err != nil {
		return false, err
	}

	_, ares, err := client.ReadMessage()
	if err != nil {
		return false, err
	}
	var r AuthResponse
	err = json.Unmarshal(ares, &r)
	if err != nil {
		return false, err
	}

	err = client.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
	)
	if err != nil {
		return false, err
	}

	return r.Success, nil
}

func keepAlive(email, password string) error {
	token, err := getAuthToken()
	if err != nil {
		return err
	}
	if !token.Success {
		return fmt.Errorf("get auth token failed")
	}

	session, err := getSession(
		email,
		password,
		token.Token,
	)
	if err != nil {
		return err
	}
	if !session.Success {
		return fmt.Errorf("get session failed, check email and password")
	}

	success, err := authWebsocket(
		session.Session,
		session.WebsocketHost,
		session.WebsocketPath+"?exclude_archives=1",
	)
	if err != nil {
		return err
	}
	if !success {
		return fmt.Errorf("auth websocket request failed")
	}

	return nil
}

func heroku_shutdown(heroku_app_name, heroku_key string) error {
	heroku.DefaultTransport.BearerToken = heroku_key

	h := heroku.NewService(heroku.DefaultClient)

	qdes := 0
	hdes := heroku.FormationUpdateOpts{Quantity: &qdes}

	_, err := h.FormationUpdate(context.Background(), heroku_app_name, "worker", hdes)
	if err != nil {
		fmt.Println("error " + heroku_app_name + " " + err.Error())
	}
	time.Sleep(1 * time.Second)
	return nil
}

func die(message string) {
	fmt.Fprintf(os.Stderr, message+"\n")
	os.Exit(1)
}

func main() {
	email := os.Getenv("IRCCLOUD_USERNAME")
	password := os.Getenv("IRCCLOUD_PASSWORD")
	heroku_app_name := os.Getenv("heroku-app-name")
	heroku_key := os.Getenv("heroku-key")

	if email == "" {
		die("IRCCLOUD_USERNAME is required")
	}
	if password == "" {
		die("IRCCLOUD_PASSWORD is required")
	}
	if heroku_app_name == "" {
		die("heroku-app-name is required")
	}
	if heroku_key == "" {
		die("heroku-key is required")
	}

	err := keepAlive(email, password)
	if err != nil {
		die(err.Error())
	}

	errh := heroku_shutdown(heroku_app_name, heroku_key)
	if errh != nil {
		die(errh.Error())
	}

}
