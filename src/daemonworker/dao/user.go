package dao

import (
	"log"

	"github.com/inspursoft/wand/src/daemonworker/models"
)

func AddOrUpdateUserAccess(username, token string, isOrg int) (err error) {
	stmt, err := GetDBConn().Prepare(
		`insert or ignore into user_access 
			(username, access_token, is_org) values (?, ?, ?);`)
	if err != nil {
		log.Printf("Failed to add or update user_access: %+v\n", err)
		return
	}
	_, err = stmt.Exec(username, token, isOrg)
	return
}

func GetUserAccess(username string) (userAccess *models.UserAccess) {
	rows, err := GetDBConn().Query(`select id, username, access_token, is_org from user_access where username = ?`, username)
	if err != nil {
		log.Printf("Failed to query: %+v\n", err)
		return
	}
	for rows.Next() {
		var ID int64
		var username string
		var accessToken string
		var isOrg int
		err = rows.Scan(&ID, &username, &accessToken, &isOrg)
		if err != nil {
			log.Printf("Failed to gather row: %+v\n", err)
		}
		userAccess = &models.UserAccess{
			ID:          ID,
			Username:    username,
			AccessToken: accessToken,
			IsOrg:       isOrg,
		}
	}
	return
}
