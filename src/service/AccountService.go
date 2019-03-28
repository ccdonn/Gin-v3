package service

import (
	"strconv"
	"strings"

	. "../config"
	. "../domain"
)

func FindAccount(uid int32) (*Account, error) {

	db := GetDBConn()
	rows, err := db.Query(
		strings.Join([]string{"select ac.agent_id, ac.username, ac.nickname, ac.password from agent ag join account ac on ag.id = ac.agent_id",
			"where ag.id = " + strconv.Itoa(int(uid)),
			"and ac.is_del = 0 and ag.is_del = 0 and ac.ban = 0 and ag.ban = 0"}, " "))

	if err != nil {
		return nil, err
	}

	// var account Account
	var (
		agentID  int32
		username string
		password string
		nickname string
		account  Account
	)

	rows.Next()
	if err = rows.Scan(&agentID, &username, &nickname, &password); err != nil {
		return nil, err
	}

	account = Account{
		AgentID:  agentID,
		Username: username,
		Nickname: nickname,
		Password: password,
	}
	return &account, nil
}
