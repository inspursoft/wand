package proxies

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/inspursoft/wand/src/daemonworker/dao"
	"github.com/inspursoft/wand/src/daemonworker/gogs"
	"github.com/inspursoft/wand/src/daemonworker/jenkins"
)

func (p *ProxiedHandler) CreateOrForkRepo(resp http.ResponseWriter, req *http.Request, data []byte) {
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
				gogs.NewGogsHandler(p.GogitsBaseURL, username, userAccess.AccessToken).
					CreateHook(fmt.Sprintf("%s/receive-webhook", p.GogitsBaseURL), username, repoName)
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
			jenkins.NewJenkinsHandler(p.JenkinsMasterURL, p.KVMRegistryURL, p.ConfigURL).CreateJobWithParameter(repoName)
		}
	}
}
