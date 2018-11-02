package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/inspursoft/wand/src/daemonworker/dao"
	"github.com/inspursoft/wand/src/daemonworker/handlers"
	"github.com/inspursoft/wand/src/daemonworker/jenkins"
	"github.com/inspursoft/wand/src/daemonworker/models"
	"github.com/inspursoft/wand/src/daemonworker/proxies"
	"github.com/inspursoft/wand/src/daemonworker/utils"
	_ "github.com/mattn/go-sqlite3"
)

const (
	kvmToolsPath    = "/root/kvm"
	kvmRegistryPath = "/root/kvmregistry"
)

func prepareKVMHost(nodeIP, nodeSSHPort, username, password, kvmToolkitsPath, kvmRegistrySize, kvmRegistryPort string) error {
	sshPort, _ := strconv.Atoi(nodeSSHPort)
	sshHandler, err := utils.NewSecureShell(nodeIP, sshPort, username, password)
	if err != nil {
		return err
	}
	kvmToolsNodePath := filepath.Join(kvmToolkitsPath, "kvm")
	kvmRegistryNodePath := filepath.Join(kvmToolkitsPath, "kvmregistry")
	err = sshHandler.ExecuteCommand(fmt.Sprintf("mkdir -p %s %s", kvmToolsNodePath, kvmRegistryNodePath))
	if err != nil {
		return err
	}
	err = sshHandler.SecureCopy(kvmToolsPath, kvmToolsNodePath)
	if err != nil {
		return err
	}
	err = sshHandler.SecureCopy(kvmRegistryPath, kvmRegistryNodePath)
	if err != nil {
		return err
	}
	return sshHandler.ExecuteCommand(fmt.Sprintf(`
		cd %s && chmod +x kvmregistry && nohup ./kvmregistry -size %s -port %s > kvmregistry.out 2>&1 &`,
		kvmRegistryNodePath, kvmRegistrySize, kvmRegistryPort))
}

func main() {
	config, err := utils.LoadConfig("/root/config.ini")
	if err != nil {
		panic(err)
	}
	utils.ListConfig(config)

	gogitsBaseURL := fmt.Sprintf("http://%s:%s", config["gogits_host_ip"], config["gogits_host_port"])
	jenkinsMasterURL := fmt.Sprintf("http://%s:%s", config["jenkins_host_ip"], config["jenkins_host_port"])
	jenkinsNodeIP := config["jenkins_node_ip"]
	jenkinsNodeSSHPort := config["jenkins_node_ssh_port"]
	jenkinsNodeUsername := config["jenkins_node_username"]
	jenkinsNodePassword := config["jenkins_node_password"]
	kvmToolkitsPath := config["kvm_toolkits_path"]
	kvmRegistrySize := config["kvm_registry_size"]
	kvmRegistryPort := config["kvm_registry_port"]

	prepareKVMHost(jenkinsNodeIP, jenkinsNodeSSHPort, jenkinsNodeUsername, jenkinsNodePassword,
		kvmToolkitsPath, kvmRegistrySize, kvmRegistryPort)

	kvmRegistryURL := fmt.Sprintf("http://%s:%s", jenkinsNodeIP, kvmRegistryPort)
	configURL := fmt.Sprintf("%s/configs", gogitsBaseURL)
	jenkinsHandler := jenkins.NewJenkinsHandler(jenkinsMasterURL, kvmRegistryURL, configURL)
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

	proxiedHandler := &proxies.ProxiedHandler{
		GogitsBaseURL:    gogitsBaseURL,
		JenkinsMasterURL: jenkinsMasterURL,
		JenkinsNodeIP:    jenkinsNodeIP,
		KVMRegistryURL:   kvmRegistryURL,
		ConfigURL:        configURL,
	}
	chainedProxy := proxies.InterceptActionByURL(reverseProxy, http.MethodPost, []string{"/user/sign_up"}, proxiedHandler.SignUp)
	chainedProxy = proxies.InterceptActionByURL(chainedProxy, http.MethodPost, []string{"/org/create"}, proxiedHandler.CreateOrg)
	chainedProxy = proxies.InterceptActionByURL(chainedProxy, http.MethodPost, []string{"/repo/create", "/repo/fork"}, proxiedHandler.CreateOrForkRepo)
	chainedProxy = proxies.InterceptActionByURL(chainedProxy, http.MethodPost, []string{"/receive-webhook"}, proxiedHandler.Webhook)

	router := mux.NewRouter()

	handler := &handlers.Handler{Cache: models.NewCachedReport()}
	commitReportRouter := router.Path("/commit-report").Subrouter()
	commitReportRouter.Methods("GET").HandlerFunc(handler.ResolveCommitReport)
	commitReportRouter.Methods("POST").HandlerFunc(handler.AddOrUpdateCommitReport)

	configRouter := router.Path("/config").Subrouter()
	configRouter.Methods("GET").HandlerFunc(handler.GetConfig)
	configRouter.Methods("PUT").HandlerFunc(handler.AddOrUpdateConfig)

	router.Path("/configs").Methods("GET").HandlerFunc(handler.FetchConfigs)
	router.Path("/upload").Methods("POST").HandlerFunc(handler.UploadResource)
	router.Path("/icon").Methods("GET").HandlerFunc(handler.ResolveIcon)

	router.NotFoundHandler = chainedProxy
	log.Fatal(http.ListenAndServe(":8088", router))
}
