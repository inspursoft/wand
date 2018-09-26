package gogs

import (
	"devopssuite/src/daemonworker/utils"
	"fmt"
	"log"
	"net/http"
	"time"
)

var maxRetryCount = 20

func NewGogsHandler(gogitsBaseURL, username, token string) *gogsHandler {
	pingURL := fmt.Sprintf("%s", gogitsBaseURL)
	for i := 0; i < maxRetryCount; i++ {
		log.Printf("Ping Gogits server %d time(s)...\n", i+1)
		if i == maxRetryCount-1 {
			log.Println("Failed to ping Gogits due to exceed max retry count.")
			break
		}
		err := utils.RequestHandle(http.MethodGet, pingURL, nil, nil,
			func(req *http.Request, resp *http.Response) error {
				if resp.StatusCode <= 400 {
					return nil
				}
				return fmt.Errorf("Requested URL %s with unexpected response: %d", pingURL, resp.StatusCode)
			})
		if err == nil {
			log.Println("Successful connected to the Gogits service.")
			break
		}
		time.Sleep(time.Second)
	}
	return &gogsHandler{
		baseURL:  gogitsBaseURL,
		username: username,
		token:    token,
	}
}

func (g *gogsHandler) CreateAccessToken(password string) (*AccessToken, error) {
	opt := createAccessTokenOption{Name: "ACCESS-TOKEN"}
	var token AccessToken
	log.Println("Requesting Gogits API of create access token ...")
	err := utils.RequestHandle(http.MethodPost, fmt.Sprintf("%s/api/v1/users/%s/tokens", g.baseURL, g.username), func(req *http.Request) error {
		req.Header = http.Header{
			"content-type":  []string{"application/json"},
			"Authorization": []string{"Basic " + utils.BasicAuthEncode(g.username, password)},
		}
		return nil
	}, &opt, func(req *http.Request, resp *http.Response) error {
		return utils.UnmarshalToJSON(resp.Body, &token)
	})
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (g *gogsHandler) getAccessHeader() http.Header {
	return http.Header{
		"content-type":  []string{"application/json"},
		"Authorization": []string{"token " + g.token},
	}
}

func (g *gogsHandler) CreateHook(triggerURL string, ownerName string, repoName string) error {

	config := make(map[string]string)
	config["url"] = triggerURL
	config["content_type"] = "json"

	opt := createHookOption{
		Type:   "gogs",
		Config: config,
		Events: []string{"push", "pull_request"},
		Active: true,
	}
	log.Println("Requesting Gogits API of create hook ...")
	return utils.SimplePostRequestHandle(fmt.Sprintf("%s/api/v1/repos/%s/%s/hooks", g.baseURL, ownerName, repoName), g.getAccessHeader(), &opt)
}
