package service

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	. "../config"
	. "../domain"
	. "../request"

	"github.com/gin-gonic/gin"
)

func FindTutorial(c *gin.Context) {
	var req TutorialRequest
	c.Bind(&req)
	//	fmt.Println(req.Query)

	// build find tutorial sql
	sqlSelect := "select " + strings.Join(tutorialDBColumn, ",") + " from tutorial"
	sqlCount := "select count(*) from tutorial"
	sqlCondition := "where 1 = 1 and (title like ? or content like ?)"
	sqlOrder := "order by"
	sqlPagination := "limit " + strconv.Itoa(int((req.PageNum-1)*req.PageSize)) + ", " + strconv.Itoa(int(req.PageSize))

	switch req.Order {
	case "lastUpdateTime":
		sqlOrder += " " + "last_update_time"
	case "createTime":
		sqlOrder += " " + "create_time"
	default:
		sqlOrder += " " + "create_time"
	}

	switch req.Desc {
	case "asc":
		sqlOrder += " asc"
	case "desc":
		sqlOrder += " desc"
	default:
		sqlOrder += " desc"
	}

	sqlQuery := strings.Join([]string{sqlSelect, sqlCondition, sqlOrder, sqlPagination}, " ")
	sqlCountQuery := strings.Join([]string{sqlCount, sqlCondition}, " ")

	fmt.Println(sqlQuery)
	fmt.Println(sqlCountQuery)
	// execution sql
	var (
		id             int32
		title          string
		titleImg       string
		content        string
		createTime     *time.Time
		del            int8
		lastUpdateUser int32
		lastUpdateTime *time.Time
	)

	db := GetDBConn()
	rows, err := db.Query(sqlQuery, "%"+string(req.Query)+"%", "%"+string(req.Query)+"%")
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"status":       "failure",
			"errorCode":    10003000,
			"errorMessage": "sql error",
		})
		log.Println(err)
		return
	}

	counts, err := db.Query(sqlCountQuery, "%"+string(req.Query)+"%", "%"+string(req.Query)+"%")
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"status":       "failure",
			"errorCode":    10003000,
			"errorMessage": "sql error",
		})
		log.Println(err)
		return
	}

	total := 0
	counts.Next()
	err = counts.Scan(&total)
	if err != nil {
		c.JSON(500, gin.H{
			"status":       "failure",
			"errorCode":    10003001,
			"errorMessage": "scan error",
		})
		log.Println(err)
		return
	}

	t := make([]*Tutorial, req.PageSize)
	tindex := 0

	for rows.Next() {
		err = rows.Scan(&id, &title, &titleImg, &content, &createTime, &del, &lastUpdateUser, &lastUpdateTime)
		if err != nil {
			c.JSON(500, gin.H{
				"status":       "failure",
				"errorCode":    10003001,
				"errorMessage": "scan error",
			})
			log.Println(err)
			return
		}

		t[tindex] = &Tutorial{
			ID:             id,
			Title:          title,
			TitleImg:       titleImg,
			Content:        content,
			Del:            del,
			LastUpdateUser: lastUpdateUser,
		}

		if createTime != nil {
			t[tindex].CreateTimeValue = createTime.Unix()
		}

		if lastUpdateTime != nil {
			t[tindex].LastUpdateTimeValue = lastUpdateTime.Unix()
		}

		tindex++
	}

	c.JSON(200, gin.H{
		"status":    "success",
		"timestamp": time.Now(),
		"result":    t[:tindex],
		"size":      tindex,
		"total":     total,
	})
	return

}

func GetTutorial(c *gin.Context) {
	n, err := strconv.Atoi(c.Param("ID"))
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"status": "failure"})
		return
	}
	t, err := getTutorial(int32(n), c)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"status": "failure",
			"error":  10003111,
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"result": t,
	})
}

var tutorialDBColumn = []string{
	"id", "title", "titleImg", "content",
	"create_time", "del", "last_update_user", "last_update_time"}

func getTutorial(id int32, c *gin.Context) (*Tutorial, error) {
	db := GetDBConn()
	var (
		title          string
		titleImg       string
		content        string
		createTime     *time.Time
		del            int8
		lastUpdateUser int32
		lastUpdateTime *time.Time
	)

	row := db.QueryRow("select "+strings.Join(tutorialDBColumn, ",")+" from tutorial where id =? and del = 0", id)

	err := row.Scan(&id, &title, &titleImg, &content, &createTime, &del, &lastUpdateUser, &lastUpdateTime)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	r := &Tutorial{
		ID:             id,
		Title:          title,
		TitleImg:       titleImg,
		Content:        content,
		Del:            del,
		LastUpdateUser: lastUpdateUser,
	}

	if createTime != nil {
		r.CreateTimeValue = createTime.Unix()
	}

	if lastUpdateTime != nil {
		r.LastUpdateTimeValue = lastUpdateTime.Unix()
	}

	return r, nil
}

func CreateTutorial(c *gin.Context) {

}

func UpdateTutorial(c *gin.Context) {
	n, err := strconv.Atoi(c.Param("ID"))
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"status": "failure"})
		return
	}

	fmt.Println(n)
}

func DeleteTutorial(c *gin.Context) {
	n, err := strconv.Atoi(c.Param("ID"))
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"status": "failure"})
		return
	}

	fmt.Println(n)
}
