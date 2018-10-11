package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/inspursoft/wand/src/daemonworker/dao"
	"github.com/inspursoft/wand/src/daemonworker/gogs"
	"github.com/inspursoft/wand/src/daemonworker/jenkins"
	"github.com/inspursoft/wand/src/daemonworker/models"
	"github.com/inspursoft/wand/src/daemonworker/utils"

	_ "github.com/mattn/go-sqlite3"
)

var kvmToolsPath = "/root/kvm"
var kvmRegistryPath = "/root/kvmregistry"
var uploadResourcePath = "/root/website"
var jsonHeader = http.Header{
	"content-type": []string{"application/json"},
}

func updateConfig(resp http.ResponseWriter, req *http.Request) {
	groupName := req.FormValue("group_name")
	username := req.FormValue("username")
	if strings.TrimSpace(groupName) == "" || strings.TrimSpace(username) == "" {
		log.Println("No group name or username provided.")
		return
	}
	var conf models.Config
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("Failed ot read from request body: %+v\n", err)
		return
	}
	err = json.Unmarshal(data, &conf)
	if err != nil {
		log.Printf("Failed to unmarshal request body: %+v\n", err)
		return
	}
	err = dao.NewBuildConfig(groupName, username).AddOrUpdateBuildConfig(conf.ConfigKey, conf.ConfigValue)
	if err != nil {
		log.Printf("Failed to add or update configure: %+v\n", err)
		return
	}
}

func fetchConfigs(resp http.ResponseWriter, req *http.Request) {
	repoName := req.FormValue("repo_name")
	username := req.FormValue("username")
	if strings.TrimSpace(repoName) == "" || strings.TrimSpace(username) == "" {
		log.Println("No repo name or username provided.")
		return
	}
	configs := dao.NewBuildConfig(repoName, username).GetBuildConfigs()
	if len(configs) > 0 {
		resp.Header().Set("content-disposition", "attachment;filename=env.cfg")
		w := bufio.NewWriter(resp)
		for _, c := range configs {
			w.WriteString(fmt.Sprintf("export %s=%s\n", c.ConfigKey, c.ConfigVal))
		}
		w.Flush()
	}
}

func uploadResource(resp http.ResponseWriter, req *http.Request) {
	fullName := req.FormValue("full_name")
	buildNumber := req.FormValue("build_number")
	if strings.TrimSpace(fullName) == "" || strings.TrimSpace(buildNumber) == "" {
		log.Println("No repo full name or build number provided.")
		return
	}
	f, fh, err := req.FormFile("upload")
	if err != nil {
		log.Printf("Failed to resolve uploaded file: %+v\n", err)
		return
	}
	uploadTargetPath := filepath.Join(uploadResourcePath, fullName, buildNumber)
	if _, err := os.Stat(uploadTargetPath); os.IsNotExist(err) {
		err = os.MkdirAll(uploadTargetPath, 0755)
		if err != nil {
			log.Printf("Failed to make dir for the target upload path: %+v\n", err)
			return
		}
	}
	targetUploadedFile := filepath.Join(uploadTargetPath, fh.Filename)
	targetFile, err := os.OpenFile(targetUploadedFile, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		log.Printf("Failed to create target uploaded: %+v\n", err)
		return
	}
	_, err = io.Copy(targetFile, f)
	if err != nil {
		log.Printf("Failed to write source to target: %+v\n", err)
		return
	}
	resp.Write([]byte(targetUploadedFile))
}

func interceptActionByURL(handler http.Handler, method string, urlList []string, action func(resp http.ResponseWriter, req *http.Request, body []byte)) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/config" {
			switch req.Method {
			case http.MethodGet:
				fetchConfigs(resp, req)
			case http.MethodPut:
				updateConfig(resp, req)
			default:
				log.Printf("Unsupported method for this request: %s\n", req.URL.Path)
			}
		} else if req.Method == http.MethodPost && req.URL.Path == "/upload" {
			uploadResource(resp, req)
		} else {
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
		}
	})
}

func main() {
	config, err := utils.LoadConfig("/root/config.ini")
	if err != nil {
		panic(err)
	}
	utils.ListConfig(config)

	gogitsBaseURL := fmt.Sprintf("http://%s:%s", config["gogits_host_ip"], config["gogits_host_port"])
	jenkinsBaseURL := fmt.Sprintf("http://%s:%s", config["jenkins_host_ip"], config["jenkins_host_port"])
	jenkinsNodeIP := config["jenkins_node_ip"]
	jenkinsNodeSSHPort := config["jenkins_node_ssh_port"]
	jenkinsNodeUsername := config["jenkins_node_username"]
	jenkinsNodePassword := config["jenkins_node_password"]
	kvmToolkitsPath := config["kvm_toolkits_path"]
	kvmRegistrySize := config["kvm_registry_size"]
	kvmRegistryPort := config["kvm_registry_port"]

	prepareKVMHost(jenkinsNodeIP, jenkinsNodeSSHPort, jenkinsNodeUsername, jenkinsNodePassword,
		kvmToolkitsPath, kvmRegistrySize, kvmRegistryPort)

	registryURL := fmt.Sprintf("http://%s:%s", jenkinsNodeIP, kvmRegistryPort)
	configURL := fmt.Sprintf("%s/config", gogitsBaseURL)
	jenkinsHandler := jenkins.NewJenkinsHandler(jenkinsBaseURL, registryURL, configURL)
	err = jenkinsHandler.CreateIgnitorJob()
	if err != nil {
		panic(fmt.Sprintf("Failed to create Jenkins ignitor job: %+v\n", err))
	}
	dao.InitDB()
	u, err := url.Parse("http://gogits:3000")
	if err != nil {
		panic(fmt.Sprintf("Failed to parse Gogits URL, error: %+v\n", err))
	}
	reverseProxy := httputil.NewSingleHostReverseProxy(u)
	chainedProxy := interceptActionByURL(reverseProxy, http.MethodPost, []string{"/user/sign_up"},
		func(resp http.ResponseWriter, req *http.Request, data []byte) {
			log.Printf("Intercepting user login: %s\n", req.URL.Path)
			form, _ := url.ParseQuery(string(data))
			username := form.Get("user_name")
			password := form.Get("password")
			time.Sleep(time.Second * 2)
			accessToken, err := gogs.NewGogsHandler(gogitsBaseURL, username, "").CreateAccessToken(password)
			if err != nil {
				log.Printf("Failed to get access token: %+v\n", err)
				return
			}
			dao.AddOrUpdateUserAccess(username, accessToken.Sha1, 0)
			log.Printf("Created access token with username: %s, access token: %s\n", username, accessToken.Sha1)
		})
	chainedProxy = interceptActionByURL(chainedProxy, http.MethodPost, []string{"/org/create"},
		func(resp http.ResponseWriter, req *http.Request, data []byte) {
			log.Printf("Intercepting org creation: %s\n", req.URL.Path)
			form, _ := url.ParseQuery(string(data))
			orgName := form.Get("org_name")
			dao.AddOrUpdateUserAccess(orgName, "", 1)
		})
	chainedProxy = interceptActionByURL(chainedProxy, http.MethodPost, []string{"/repo/create", "/repo/fork"},
		func(resp http.ResponseWriter, req *http.Request, data []byte) {
			log.Printf("Intercepting repo creation: %s\n", req.URL.Path)
			header := resp.Header()
			location := header.Get("Location")
			parts := strings.Split(location, "/")
			if len(parts) >= 3 {
				username := parts[1]
				repoName := parts[2]
				userAccess := dao.GetUserAccess(username)
				if userAccess != nil {
					log.Printf("Response location: %s, username: %s, repo name: %s, access token: %s, is org: %d\n",
						location, username, repoName, userAccess.AccessToken, userAccess.IsOrg)
					bc := dao.NewBuildConfig(repoName, username)
					if userAccess.IsOrg == 0 { //Org user does not have access token cannot create hook for repo.
						gogs.NewGogsHandler(gogitsBaseURL, username, userAccess.AccessToken).
							CreateHook(fmt.Sprintf("%s/receive-webhook", gogitsBaseURL), username, repoName)
						if strings.HasPrefix(req.URL.Path, "/repo/fork") {
							repoName = username + "_" + repoName
							bc.GroupName = repoName
						}
						bc.AddOrUpdateBuildConfig("group_name", repoName)
					} else if userAccess.IsOrg == 1 {
						bc.GroupName = repoName
						bc.AddOrUpdateBuildConfig("group_name", bc.GroupName)
					}

					form, _ := url.ParseQuery(string(data))
					affinity := form.Get("description")
					if strings.TrimSpace(affinity) == "" {
						affinity = "golang"
					}
					bc.AddOrUpdateBuildConfig("affinity", affinity)
					bc.AddOrUpdateBuildConfig("last_coverage", "-")
					jenkinsHandler.CreateJobWithParameter(repoName)
				}
			}
		})
	chainedProxy = interceptActionByURL(chainedProxy, http.MethodPost, []string{"/receive-webhook"},
		func(resp http.ResponseWriter, req *http.Request, data []byte) {
			resp.WriteHeader(http.StatusOK)
			event := req.Header.Get("X-Gogs-Event")
			log.Printf("Intercepting event: %s webhook ...\n", event)
			var payload models.CustomWebhookPayload
			payload.GogsURL = gogitsBaseURL
			payload.APIURL = fmt.Sprintf("%s/api/v1", gogitsBaseURL)
			payload.MasterURL = jenkinsBaseURL
			payload.NodeIP = jenkinsNodeIP
			payload.RegistryURL = registryURL
			trigger, err := payload.AdaptToCustom(event, data)
			if err != nil {
				log.Printf("Failed to convert payload with event: %s, error: %+v\n", event, err)
			}
			if trigger {
				dao.NewBuildConfig(payload.RepoName, payload.Username).SetPayload(&payload).Update()
				log.Printf("raw: %+v\n\ncustom: %+v\n", string(data), utils.PrettyPrintJSON(payload))
				utils.SimplePostRequestHandle(fmt.Sprintf("%s/generic-webhook-trigger/invoke", jenkinsBaseURL), jsonHeader, payload)
			}
		})
	http.ListenAndServe(":8088", chainedProxy)
}
