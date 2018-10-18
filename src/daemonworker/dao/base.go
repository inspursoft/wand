package dao

import (
	"database/sql"
	"fmt"
	"log"
)

func GetDBConn() (db *sql.DB) {
	db, err := sql.Open("sqlite3", "/data/storage.db")
	if err != nil {
		panic(fmt.Sprintf("Failed to connect db: %+v", err))
	}
	return
}

func InitDB() (err error) {
	_, err = GetDBConn().Exec(
		` create table if not exists build_config
		(group_name string, username string, config_key text, config_val text,
			primary key('group_name', 'username', 'config_key'));
		create table if not exists user_access
		 (id integer primary key autoincrement, username text unique, access_token text, is_org integer);
		create table if not exists commit_report
		 (commit_id string, report string, primary key('commit_id'));
		`)
	if err != nil {
		log.Printf("Failed to create table: %+v\n", err)
		return
	}
	return
}

func CleanUp() (err error) {
	_, err = GetDBConn().Exec(`
		drop table if exists user_access;
		drop table if exists build_config;`)
	if err != nil {
		log.Printf("Failed to drop table: %+v\n", err)
		return
	}
	return
}
