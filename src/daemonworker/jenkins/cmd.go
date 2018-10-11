package jenkins

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/inspursoft/wand/src/daemonworker/utils"
)

type jenkinsHandler struct {
	jenkinsBaseURL string
	registryURL    string
	configURL      string
}

var maxRetryCount = 120
var seedIgnitorJobName = "base_ignitor"
var seedJobName = "base"

func NewJenkinsHandler(jenkinsBaseURL, registryURL, configURL string) *jenkinsHandler {
	pingURL := fmt.Sprintf("%s/job/%s", jenkinsBaseURL, seedJobName)
	for i := 0; i < maxRetryCount; i++ {
		log.Printf("Ping Jenkins server %d time(s)...\n", i+1)
		if i == maxRetryCount-1 {
			log.Println("Failed to ping Jenkins due to exceed max retry count.")
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
			log.Println("Successful connected to the Jenkins service.")
			break
		}
		time.Sleep(time.Second)
	}
	return &jenkinsHandler{
		jenkinsBaseURL: jenkinsBaseURL,
		registryURL:    registryURL,
		configURL:      configURL,
	}
}

func (j *jenkinsHandler) CreateIgnitorJob() error {
	return utils.SimpleGetRequestHandle(fmt.Sprintf("%s/job/%s/buildWithParameters?F00=%s&F01=%s", j.jenkinsBaseURL, seedIgnitorJobName, j.registryURL, j.configURL))
}

func (j *jenkinsHandler) CreateJobWithParameter(jobName string) error {
	return utils.SimpleGetRequestHandle(
		fmt.Sprintf("%s/job/%s/buildWithParameters?F00=%s&F01=%s&F02=%s&F03=%s", j.jenkinsBaseURL, seedJobName, jobName, j.configURL, j.jenkinsBaseURL, j.registryURL))
}
