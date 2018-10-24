package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/inspursoft/wand/src/daemonworker/dao"
	"github.com/inspursoft/wand/src/daemonworker/models"
)

func (c *Handler) AddOrUpdateConfig(resp http.ResponseWriter, req *http.Request) {
	groupName := req.FormValue("group_name")
	username := req.FormValue("username")
	if strings.TrimSpace(groupName) == "" || strings.TrimSpace(username) == "" {
		rendStatus(resp, http.StatusBadRequest, fmt.Sprintf("No group name or username provided.\n"))
		return
	}
	var conf models.Config
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		rendStatus(resp, http.StatusInternalServerError, fmt.Sprintf("Failed to read from request body: %+v\n", err))
		return
	}
	err = json.Unmarshal(data, &conf)
	if err != nil {
		rendStatus(resp, http.StatusInternalServerError, fmt.Sprintf("Failed to unmarshal request body: %+v\n", err))
		return
	}
	err = dao.NewBuildConfig(groupName, username).AddOrUpdateBuildConfig(conf.ConfigKey, conf.ConfigValue)
	if err != nil {
		rendStatus(resp, http.StatusInternalServerError, fmt.Sprintf("Failed to add or update configure: %+v\n", err))
		return
	}
}

func (c *Handler) FetchConfigs(resp http.ResponseWriter, req *http.Request) {
	repoName := req.FormValue("repo_name")
	username := req.FormValue("username")
	if strings.TrimSpace(repoName) == "" || strings.TrimSpace(username) == "" {
		rendStatus(resp, http.StatusBadRequest, fmt.Sprintf("No repo name or username provided.\n"))
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

func (c *Handler) GetConfig(resp http.ResponseWriter, req *http.Request) {
	groupName := req.FormValue("group_name")
	username := req.FormValue("username")
	configKey := req.FormValue("config_key")
	if strings.TrimSpace(groupName) == "" || strings.TrimSpace(username) == "" || strings.TrimSpace(configKey) == "" {
		rendStatus(resp, http.StatusBadRequest, fmt.Sprintf("No group name, username or config key provided.\n"))
		return
	}
	config := dao.NewBuildConfig(groupName, username).GetBuildConfigByKey(configKey)
	if config == nil {
		rendStatus(resp, http.StatusNotFound, fmt.Sprintf("No config found with key: %s\n", configKey))
		return
	}
	resp.Write([]byte(config.ConfigVal))
}
