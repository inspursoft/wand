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
	groupName := req.FormValue("group_name")
	username := req.FormValue("username")
	if strings.TrimSpace(groupName) == "" || strings.TrimSpace(username) == "" {
		rendStatus(resp, http.StatusBadRequest, fmt.Sprintf("No group name or username provided.\n"))
		return
	}
	w := bufio.NewWriter(resp)
	defer w.Flush()
	isCached := req.FormValue("is_cached")
	if isCached == "true" {
		serialID := req.FormValue("serial_id")
		if strings.TrimSpace(serialID) == "" {
			rendStatus(resp, http.StatusBadRequest, fmt.Sprintf("No serial ID provided for retrieving cached data.\n"))
			return
		}
		resp.Header().Set("content-disposition", "attachment;filename=custom.cfg")
		for _, cachedConfig := range c.Cache.All() {
			if strings.HasPrefix(cachedConfig.Key, "COMMITS_") || strings.Index(cachedConfig.Key, fmt.Sprintf("%s_%s_%s", groupName, username, serialID)) == -1 {
				continue
			}
			w.WriteString(fmt.Sprintf("%s=%s\n", cachedConfig.Key[strings.LastIndex(cachedConfig.Key, "_")+1:], cachedConfig.Value))
		}
	} else {
		resp.Header().Set("content-disposition", "attachment;filename=env.cfg")
		for _, c := range dao.NewBuildConfig(groupName, username).GetBuildConfigs() {
			w.WriteString(fmt.Sprintf("%s=%s\n", c.ConfigKey, c.ConfigVal))
		}
	}
}

func (c *Handler) GetConfig(resp http.ResponseWriter, req *http.Request) {
	groupName := req.FormValue("group_name")
	username := req.FormValue("username")
	serialID := req.FormValue("serial_id")
	configKey := req.FormValue("config_key")
	if strings.TrimSpace(groupName) == "" || strings.TrimSpace(username) == "" || strings.TrimSpace(serialID) == "" || strings.TrimSpace(configKey) == "" {
		rendStatus(resp, http.StatusBadRequest, fmt.Sprintf("No group name, username, serial ID or config key provided.\n"))
		return
	}
	config := dao.NewBuildConfig(groupName, username).GetBuildConfigByKey(configKey)
	var value string
	if config == nil {
		cachedConfig, found := c.Cache.Get(fmt.Sprintf("%s_%s_%s_%s", groupName, username, serialID, configKey))
		if !found {
			rendStatus(resp, http.StatusNotFound,
				fmt.Sprintf("No config found with key: %s under groupName: %s and username: %s\n ", configKey, groupName, username))
			return
		}
		value = cachedConfig.Value
	} else {
		value = config.ConfigVal
	}
	resp.Write([]byte(value))
}

func (c *Handler) CacheConfig(resp http.ResponseWriter, req *http.Request) {
	groupName := req.FormValue("group_name")
	username := req.FormValue("username")
	serialID := req.FormValue("serial_id")
	if strings.TrimSpace(groupName) == "" || strings.TrimSpace(username) == "" || strings.TrimSpace(serialID) == "" {
		rendStatus(resp, http.StatusBadRequest, fmt.Sprintf("No group name, username or serial ID provided.\n"))
		return
	}

	var cachedConfig models.CachedConfig
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		rendStatus(resp, http.StatusInternalServerError, fmt.Sprintf("Failed to read from request body: %+v\n", err))
		return
	}
	err = json.Unmarshal(data, &cachedConfig)
	if err != nil {
		rendStatus(resp, http.StatusInternalServerError, fmt.Sprintf("Failed to unmarshal request body: %+v\n", err))
		return
	}
	cachedConfig.Key = fmt.Sprintf("%s_%s_%s_%s", groupName, username, serialID, cachedConfig.Key)
	c.Cache.Add(&cachedConfig)
}
