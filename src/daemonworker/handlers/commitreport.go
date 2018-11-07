package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/inspursoft/wand/src/daemonworker/dao"
	"github.com/inspursoft/wand/src/daemonworker/models"
	"github.com/inspursoft/wand/src/daemonworker/utils"
)

func (c *Handler) AddOrUpdateCommitReport(resp http.ResponseWriter, req *http.Request) {
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		rendStatus(resp, http.StatusInternalServerError, fmt.Sprintf("Failed to read from request body: %+v\n", err))
		return
	}
	var commitReport models.CommitReport
	err = json.Unmarshal(data, &commitReport)
	if err != nil {
		rendStatus(resp, http.StatusInternalServerError, fmt.Sprintf("Failed to unmarshal request body: %+v\n", err))
		return
	}
	err = dao.AddOrUpdateCommitReport(commitReport)
	if err != nil {
		rendStatus(resp, http.StatusInternalServerError, fmt.Sprintf("Failed to add or update commit report: %+v\n", err))
		return
	}
	c.Cache.Add(commitReport.ToCachedConfig())
}

func (c *Handler) ResolveCommitReport(resp http.ResponseWriter, req *http.Request) {
	commitID := req.FormValue("commit_id")
	if strings.TrimSpace(commitID) == "" {
		rendStatus(resp, http.StatusBadRequest, fmt.Sprintln("No commit ID provided."))
		return
	}
	config, found := c.Cache.Get(fmt.Sprintf("COMMITS_%s", commitID))
	var commitReport *models.CommitReport
	if !found {
		log.Printf("No found commit report with commit ID: %s, will retrieving from DB ...\n", commitID)
		commitReport = dao.GetCommitReport(commitID)
		if commitReport == nil {
			log.Printf("Initialized it with commit ID: %s, store into cache as no found from DB ...\n", commitID)
			commitReport = &models.CommitReport{CommitID: commitID, Report: ""}
		}
		c.Cache.Add(commitReport.ToCachedConfig())
	} else {
		commitReport = models.ToCommitReport(config)
	}
	if commitReport.CommitID == "" {
		return
	}
	if strings.Index(commitReport.Report, "|") == -1 {
		return
	}
	parts := strings.Split(commitReport.Report, "|")
	target := req.FormValue("target")
	if target == "report" {
		utils.DrawTag(resp, func(result string) string {
			if result == "pass" {
				return "correct.png"
			}
			return "wrong.png"
		}(parts[0]))
	} else {
		http.Redirect(resp, req, parts[1], http.StatusSeeOther)
	}
}
