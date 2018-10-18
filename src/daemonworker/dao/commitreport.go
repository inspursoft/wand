package dao

import (
	"log"

	"github.com/inspursoft/wand/src/daemonworker/models"
)

func AddOrUpdateCommitReport(commitReport models.CommitReport) (err error) {
	stmt, err := GetDBConn().Prepare(`insert or replace into commit_report (commit_id, report) values (?, ?);`)
	if err != nil {
		log.Printf("Failed to add or update commit report: %+v\n", err)
		return
	}
	_, err = stmt.Exec(commitReport.CommitID, commitReport.Report)
	return
}

func GetCommitReport(commitID string) (commitReport *models.CommitReport) {
	rows, err := GetDBConn().Query(`select commit_id, report from commit_report where commit_id = ?;`, commitID)
	for rows.Next() {
		var commitID string
		var report string
		err = rows.Scan(&commitID, &report)
		if err != nil {
			log.Printf("Failed to gather row: %+v\n", err)
		}
		commitReport = &models.CommitReport{
			CommitID: commitID,
			Report:   report,
		}
	}
	return
}
