package proxies

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type ProxiedHandler struct {
	GogitsBaseURL    string
	JenkinsMasterURL string
	JenkinsNodeIP    string
	KVMRegistryURL   string
	ConfigURL        string
}

func InterceptActionByURL(handler http.Handler, method string, urlList []string, action func(resp http.ResponseWriter, req *http.Request, body []byte)) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		data, _ := ioutil.ReadAll(req.Body)
		for _, urlStr := range urlList {
			if req.Method == method && strings.HasPrefix(req.URL.Path, urlStr) {
				defer func() {
					go func() {
						time.Sleep(time.Second * 2)
						action(resp, req, data)
					}()
				}()
			}
		}
		req.Body = ioutil.NopCloser(bytes.NewBuffer(data))
		handler.ServeHTTP(resp, req)
	})
}
