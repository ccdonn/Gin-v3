package service

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"../config"
	"../domain"
	"../request"
	"../utils"

	"github.com/gin-gonic/gin"
)

// FindSuggestion : find suggestions
func FindSuggestion(c *gin.Context) {
	var req request.SuggestionRequest
	c.Bind(&req)

	sqlSelect := "select id, agent_id, nickname, username, type, content, create_time, reply_content, reply_time, status from suggestion"
	sqlSelectCount := "select count(*) from suggestion"
	sqlOrder := "order by create_time desc"
	sqlPagination := "limit " + strconv.Itoa(int((req.PageNum-1)*req.PageSize)) + ", " + strconv.Itoa(int(req.PageSize))
	sqlCondition := "where 1 = 1"
	if req.Status > 0 {
		sqlCondition += " " + "and status = " + strconv.Itoa(int(req.Status))
	} else {
		sqlCondition += " " + "and status > 0"
	}

	if req.AgentID != 0 {
		sqlCondition += " " + "and agent_id = " + strconv.Itoa(int(req.AgentID))
	}

	if req.Type != 0 {
		sqlCondition += " " + "and type = " + strconv.Itoa(int(req.Type))
	}

	if req.StartTime != "" {
		if ok, _ := regexp.MatchString("([12]\\d{3}-(0[1-9]|1[0-2])-(0[1-9]|[12]\\d|3[01]))", req.StartTime); ok {
			sqlCondition += " " + "and date(create_time) >= \"" + req.StartTime + "\""
		}
	}

	if req.EndTime != "" {
		if ok, _ := regexp.MatchString("([12]\\d{3}-(0[1-9]|1[0-2])-(0[1-9]|[12]\\d|3[01]))", req.EndTime); ok {
			sqlCondition += " " + "and date(create_time) <= \"" + req.EndTime + "\""
		}
	}

	// sqlQuery := sqlSelect + sqlCondition + sqlOrder + sqlPagination
	sqlQuery := strings.Join([]string{sqlSelect, sqlCondition, sqlOrder, sqlPagination}, " ")
	sqlCountQuery := strings.Join([]string{sqlSelectCount, sqlCondition}, " ")
	// log.Println("SQL query:", sqlQuery)
	// log.Println("SQL countquery:", sqlCountQuery)

	db := config.GetDBConn()
	// rows, err := db.Query("select id, agent_id, nickname, username, type, content, create_time, reply_content, reply_time, status from suggestion where status > 0 order by create_time desc")
	rows, err := db.Query(sqlQuery)
	if err != nil {
		c.JSON(500, gin.H{
			"error":   "sql error",
			"message": err,
		})
		return
	}

	counts, err := db.Query(sqlCountQuery)
	if err != nil {
		c.JSON(500, gin.H{
			"error":   "sql error",
			"message": err,
		})
		return
	}

	total := 0
	counts.Next()
	err = counts.Scan(&total)
	if err != nil {
		c.JSON(500, gin.H{
			"error":   "scan error",
			"message": err,
		})
		return
	}

	var (
		id           int32
		agentID      int32
		nickname     string
		username     string
		Type         int32
		content      string
		createTime   *time.Time
		replyContent string
		replyTime    *time.Time
		status       int32
	)

	s := make([]*domain.Suggestion, req.PageSize)
	sindex := 0

	for rows.Next() {
		err = rows.Scan(&id, &agentID, &nickname, &username, &Type, &content, &createTime, &replyContent, &replyTime, &status)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   "scan error",
				"message": err,
			})
			return
		}

		sr := &domain.Suggestion{
			ID:           id,
			AgentID:      agentID,
			Nickname:     nickname,
			Username:     username,
			Type:         Type,
			Content:      content,
			ReplyContent: replyContent,
			Status:       status,
		}

		if createTime != nil {
			sr.CreateTimeValue = createTime.Unix()
		}

		if replyTime != nil {
			sr.ReplyTimeValue = replyTime.Unix()
		}

		s[sindex] = sr
		sindex++
	}

	c.JSON(200, gin.H{
		"status":    "success",
		"timestamp": time.Now(),
		"result":    s[:sindex],
		"size":      sindex,
		"total":     total,
	})
	return
}

// GetSuggestion : get suggestion
func GetSuggestion(c *gin.Context) {
	n, err := strconv.Atoi(c.Param("ID"))
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"status": "failure"})
		return
	}

	sr, err := getSuggestion(int32(n), c)
	if err != nil {
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"result": sr,
	})
	return
}

func getSuggestion(id int32, c *gin.Context) (*domain.Suggestion, error) {
	sqlSelect := "select id, agent_id, nickname, username, type, content, create_time, reply_content, reply_time, status from suggestion"
	sqlCondition := "where 1 = 1 and id = " + strconv.Itoa(int(id)) + " and status > 0"

	sqlQuery := strings.Join([]string{sqlSelect, sqlCondition}, " ")
	db := config.GetDBConn()
	row, err := db.Query(sqlQuery)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"status":       "failure",
			"errorCode":    10003000,
			"errorMessage": "sql error",
		})
		log.Println(err)
		return nil, err
	}

	var (
		agentID      int32
		nickname     string
		username     string
		Type         int32
		content      string
		createTime   *time.Time
		replyContent string
		replyTime    *time.Time
		status       int32
	)

	if !row.Next() {
		c.JSON(200, gin.H{
			"status": "success",
		})
		return nil, nil
	}

	err = row.Scan(&id, &agentID, &nickname, &username, &Type, &content, &createTime, &replyContent, &replyTime, &status)
	if err != nil {
		c.JSON(500, gin.H{
			"status":       "failure",
			"errorMessage": "scan error",
			"errorCode":    10003001,
		})
		log.Println(err)
		return nil, err
	}

	sr := &domain.Suggestion{
		ID:           id,
		AgentID:      agentID,
		Nickname:     nickname,
		Username:     username,
		Type:         Type,
		Content:      content,
		ReplyContent: replyContent,
		Status:       status,
	}

	if createTime != nil {
		sr.CreateTimeValue = createTime.Unix()
	}

	if replyTime != nil {
		sr.ReplyTimeValue = replyTime.Unix()
	}

	return sr, nil
}

// CreateSuggestion : create a new suggestion
func CreateSuggestion(c *gin.Context) {
	var body domain.Suggestion
	c.Bind(&body)

	if body.Type == 0 || body.Content == "" {
		c.AbortWithStatusJSON(400, gin.H{
			"status":       "failure",
			"errorCode":    10005000,
			"errorMessage": "param check fail",
		})
		return
	}

	uid := utils.ExtractAgentID(c.Request.Header.Get("Authorization"))
	account, err := FindAccount(uid)

	if err != nil {
		log.Println(err)
	}

	sr := &domain.Suggestion{
		AgentID:  account.AgentID,
		Username: account.Username,
		Nickname: account.Nickname,
		Content:  body.Content,
		Type:     body.Type,
	}

	create(sr)

	fmt.Println(body)
	return
}

func create(s *domain.Suggestion) {
	db := config.GetDBConn()

	insert, err := db.Prepare("insert into suggestion(" + strings.Join(SuggestionDBColumn, ",") + ") " +
		" values (?,?,?,?,?,?,?,now(),?,?,?,?)")
	if err != nil {
		log.Println(err)
	}

	// location, err := time.LoadLocation("Asia/Taipei")
	_, err = insert.Exec(nil, s.AgentID, s.Nickname, s.Username, s.Type, s.Content, "", "", "", nil, 1)
	if err != nil {
		log.Println(err)
	}

	return
}

// SuggestionDBColumn : columns of table Suggestion
var SuggestionDBColumn = []string{
	"id", "agent_id", "nickname", "username", "type",
	"content", "image", "create_time", "device_info", "reply_content",
	"reply_time", "status"}

// PartialUpdateSuggestion : update reply_content of Suggestion
func PartialUpdateSuggestion(c *gin.Context) {
	var body domain.Suggestion

	n, err := strconv.Atoi(c.Param("ID"))
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"status":       "failure",
			"errorCode":    10005000,
			"errorMessage": "param check fail",
		})
		return
	}

	sr, err := getSuggestion(int32(n), c)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"status":       "failure",
			"errorCode":    10002001,
			"errorMessage": "suggestion not found",
		})
		return
	}

	c.Bind(&body)
	if body.ReplyContent == "" {
		c.AbortWithStatusJSON(400, gin.H{
			"status":       "failure",
			"errorCode":    10005000,
			"errorMessage": "param check fail",
		})
		return
	}

	sr.ReplyContent = body.ReplyContent

	partialUpdateSuggestion(sr)

	return
}

func partialUpdateSuggestion(s *domain.Suggestion) bool {
	db := config.GetDBConn()

	update, err := db.Prepare("update suggestion set reply_content=?, reply_time=now(), status=2 where id = " + strconv.Itoa(int(s.ID)))
	if err != nil {
		log.Println(err)
		return false
	}

	// location, err := time.LoadLocation("Asia/Taipei")
	_, err = update.Exec(s.ReplyContent)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

/* not support delete operation */
// func DeleteSuggestion(c *gin.Context) {
// 	n, err := strconv.Atoi(c.Param("ID"))
// 	if err != nil {
// 		c.AbortWithStatusJSON(400, gin.H{"status": "failure"})
// 		return
// 	}

// 	fmt.Println(n)
// 	return
// }
