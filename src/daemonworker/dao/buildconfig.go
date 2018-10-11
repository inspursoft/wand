package dao

import (
	"log"

	"github.com/inspursoft/wand/src/daemonworker/models"
)

type buildConfig struct {
	Payload   *models.CustomWebhookPayload
	GroupName string
	Username  string
	ConfigKey string
	ConfigVal string
}

func NewBuildConfig(groupName string, username string) *buildConfig {
	log.Printf("Get build configs with group name: %s, username: %s\n", groupName, username)
	return &buildConfig{GroupName: groupName, Username: username}
}

func (bc *buildConfig) SetPayload(payload *models.CustomWebhookPayload) *buildConfig {
	userAccess := GetUserAccess(bc.Username)
	if userAccess != nil {
		if userAccess.IsOrg == 0 {
			conf := bc.GetBuildConfigByKey("group_name")
			if conf.ConfigVal == "" {
				bc.GroupName = bc.Username + "_" + bc.GroupName
				payload.RepoName = bc.GroupName
			} else {
				bc.GroupName = conf.ConfigVal
			}
		}
		payload.GroupName = bc.GroupName
	} else {
		log.Printf("Failed to get user access with username: %s\n", bc.Username)
	}
	payload.Affinity = bc.GetBuildConfigByKey("affinity").ConfigVal
	payload.LastCoverage = bc.GetBuildConfigByKey("last_coverage").ConfigVal
	collaborator := GetUserAccess(payload.Collaborator)
	if collaborator != nil {
		payload.AccessToken = collaborator.AccessToken
	}
	bc.Payload = payload
	return bc
}

func (bc *buildConfig) AddOrUpdateBuildConfig(key, val string) (err error) {
	stmt, err := GetDBConn().Prepare(`insert or replace into build_config 
		(group_name, username, config_key, config_val) 
		values (?, ?, ?, ?);`)
	if err != nil {
		log.Printf("Failed to add or update table build_config: %+v\n", err)
		return
	}
	_, err = stmt.Exec(bc.GroupName, bc.Username, key, val)
	return
}

func (bc *buildConfig) GetBuildConfigs() (configs []buildConfig) {
	rows, err := GetDBConn().Query(
		`select config_key, config_val from build_config where group_name = ? and username = ?;`, bc.GroupName, bc.Username)
	if err != nil {
		log.Printf("Failed to query: %+v\n", err)
		return
	}
	log.Printf("Retrieve configure with group: %s and username: %s\n", bc.GroupName, bc.Username)
	configs = []buildConfig{}
	for rows.Next() {
		var key string
		var val string
		err = rows.Scan(&key, &val)
		if err != nil {
			log.Printf("Failed to gather row: %+v\n", err)
		}
		conf := buildConfig{
			ConfigKey: key,
			ConfigVal: val,
		}
		configs = append(configs, conf)
	}
	return
}

func (bc *buildConfig) GetBuildConfigByKey(key string) (config buildConfig) {
	rows, err := GetDBConn().Query(
		`select config_key, config_val from build_config where group_name = ? and username = ? and config_key = ?;`, bc.GroupName, bc.Username, key)
	if err != nil {
		log.Printf("Failed to query: %+v\n", err)
		return
	}
	log.Printf("Retrieve configure with group: %s and username: %s\n", bc.GroupName, bc.Username)
	for rows.Next() {
		var key string
		var val string
		err = rows.Scan(&key, &val)
		if err != nil {
			log.Printf("Failed to gather row: %+v\n", err)
		}
		config = buildConfig{
			ConfigKey: key,
			ConfigVal: val,
		}
	}
	return
}

func (bc *buildConfig) Update() {
	bc.AddOrUpdateBuildConfig("action", bc.Payload.Action)
	bc.AddOrUpdateBuildConfig("username", bc.Payload.Username)
	bc.AddOrUpdateBuildConfig("full_name", bc.Payload.FullName)
	bc.AddOrUpdateBuildConfig("gogs_url", bc.Payload.GogsURL)
	bc.AddOrUpdateBuildConfig("jenkins_master_url", bc.Payload.MasterURL)
	bc.AddOrUpdateBuildConfig("jenkins_node_ip", bc.Payload.NodeIP)
	bc.AddOrUpdateBuildConfig("kvm_registry_url", bc.Payload.RegistryURL)
	bc.AddOrUpdateBuildConfig("base_repo_name", bc.Payload.BaseRepo.RepoName)
	bc.AddOrUpdateBuildConfig("base_repo_branch", bc.Payload.BaseRepo.Branch)
	bc.AddOrUpdateBuildConfig("base_repo_clone_url", bc.Payload.BaseRepo.CloneURL)
	bc.AddOrUpdateBuildConfig("base_repo_html_url", bc.Payload.BaseRepo.HTMLURL)
	bc.AddOrUpdateBuildConfig("head_repo_name", bc.Payload.HeadRepo.RepoName)
	bc.AddOrUpdateBuildConfig("head_repo_branch", bc.Payload.HeadRepo.Branch)
	bc.AddOrUpdateBuildConfig("head_repo_clone_url", bc.Payload.HeadRepo.CloneURL)
	bc.AddOrUpdateBuildConfig("head_repo_html_url", bc.Payload.HeadRepo.HTMLURL)
	bc.AddOrUpdateBuildConfig("html_url", bc.Payload.HTMLURL)
	bc.AddOrUpdateBuildConfig("comment_url", bc.Payload.CommentURL)
	bc.AddOrUpdateBuildConfig("last_coverage", "\""+bc.Payload.LastCoverage+"\"")
	bc.AddOrUpdateBuildConfig("access_token", bc.Payload.AccessToken)

}
