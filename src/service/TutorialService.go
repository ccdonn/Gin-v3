package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"../config"
	"../domain"
	"../request"
	"../utils"

	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic"
)

// FindTutorial : find tutorials
func FindTutorial(c *gin.Context) {
	var req request.TutorialRequest
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

	db := config.GetDBConn()

	rows, err := db.Query(sqlQuery, "%"+string(req.Query)+"%", "%"+string(req.Query)+"%")
	defer rows.Close()

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
	defer counts.Close()

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

	t := make([]*domain.Tutorial, req.PageSize)
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

		t[tindex] = &domain.Tutorial{
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

// GetTutorial : get single tutorial
func GetTutorial(c *gin.Context) {
	n, err := strconv.Atoi(c.Param("ID"))
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"status": "failure"})
		return
	}
	t, err := getTutorial(int32(n))
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

func getTutorial(id int32) (*domain.Tutorial, error) {
	db := config.GetDBConn()

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

	r := &domain.Tutorial{
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

// CreateTutorial : create a new tutorial
func CreateTutorial(c *gin.Context) {
	var body domain.Tutorial
	c.BindJSON(&body)

	if body.Title == "" || body.TitleImg == "" || body.Content == "" {
		c.AbortWithStatusJSON(400, gin.H{
			"status":       "failure",
			"errorCode":    10005000,
			"errorMessage": "param check fail",
		})
		log.Println("param check fail")
		return
	}

	uid := utils.ExtractAgentID(c.Request.Header.Get("Authorization"))
	_, err := FindAccount(uid)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"status":       "failure",
			"errorCode":    10005001,
			"errorMessage": "user not found",
		})
		log.Println(err)
		return
	}

	body.LastUpdateUser = uid

	if success := createTutorial(&body, c); success {
		c.JSON(200, gin.H{
			"status": "success",
		})
	}

	return
}

func createTutorial(t *domain.Tutorial, c *gin.Context) bool {
	db := config.GetDBConn()

	insert, err := db.Prepare("insert into tutorial (" + strings.Join(tutorialDBColumn, ",") + ") " +
		"values (null,?,?,?,now(),?,?,now()) ")
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"status":       "failure",
			"errorCode":    10003000,
			"errorMessage": "sql error",
		})
		log.Println(err)
		return false
	}

	if _, err = insert.Exec(t.Title, t.TitleImg, t.Content, 0, t.LastUpdateUser); err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"status":       "failure",
			"errorCode":    10003000,
			"errorMessage": "sql error",
		})
		log.Println(err)
		return false
	}

	return true
}

// UpdateTutorial : update exist tutorial
func UpdateTutorial(c *gin.Context) {
	n, err := strconv.Atoi(c.Param("ID"))
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"status": "failure"})
		return
	}

	var body domain.Tutorial
	c.BindJSON(&body)

	if n == 0 || body.Title == "" || body.TitleImg == "" || body.Content == "" {
		c.AbortWithStatusJSON(400, gin.H{
			"status":       "failure",
			"errorCode":    10005000,
			"errorMessage": "param check fail",
		})
		log.Println("param check fail")
		return
	}

	uid := utils.ExtractAgentID(c.Request.Header.Get("Authorization"))
	_, err = FindAccount(uid)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"status":       "failure",
			"errorCode":    10005001,
			"errorMessage": "user not found",
		})
		log.Println(err)
		return
	}

	t, err := getTutorial(int32(n))
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"status":       "failure",
			"errorCode":    10003111,
			"errorMessage": "tutorial not found",
		})
		log.Println(err)
		return
	}

	t.LastUpdateUser = uid
	t.Title = body.Title
	t.TitleImg = body.TitleImg
	t.Content = body.Content

	if ok, err := updateTutorial(t); ok {
		c.JSON(200, gin.H{
			"status": "success",
		})
	} else {
		c.AbortWithStatusJSON(500, gin.H{
			"status":       "failure",
			"errorCode":    10003000,
			"errorMessage": "sql error",
		})
		log.Println(err)
	}

	return
}

func updateTutorial(t *domain.Tutorial) (bool, error) {
	db := config.GetDBConn()

	update, err := db.Prepare("update tutorial set title=?, titleImg=?, content=?, last_update_user=?, last_update_time=now() where id = ?")
	if err != nil {
		log.Println(err)
		return false, err
	}

	_, err = update.Exec(t.Title, t.TitleImg, t.Content, t.LastUpdateUser, t.ID)
	if err != nil {
		log.Println(err)
		return false, err
	}

	return true, nil
}

// DeleteTutorial : delete exist tutorial
func DeleteTutorial(c *gin.Context) {
	n, err := strconv.Atoi(c.Param("ID"))
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"status": "failure"})
		return
	}

	uid := utils.ExtractAgentID(c.Request.Header.Get("Authorization"))
	_, err = FindAccount(uid)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"status":       "failure",
			"errorCode":    10005001,
			"errorMessage": "user not found",
		})
		log.Println(err)
		return
	}

	_, err = getTutorial(int32(n))
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"status":       "failure",
			"errorCode":    10003111,
			"errorMessage": "tutorial not found",
		})
		return
	}

	if ok, err := deleteTutorial(int32(n), uid); ok {
		c.JSON(200, gin.H{
			"status": "success",
		})
	} else {
		c.AbortWithStatusJSON(500, gin.H{
			"status":       "failure",
			"errorCode":    10003000,
			"errorMessage": "sql error",
		})
		log.Println(err)
	}

	return
}

func deleteTutorial(id, agentID int32) (bool, error) {
	db := config.GetDBConn()

	delete, err := db.Prepare("update tutorial set del = 1, last_update_user=?, last_update_time=now() where id = ?")
	if err != nil {
		log.Println(err)
		return false, err
	}

	_, err = delete.Exec(agentID, id)
	if err != nil {
		log.Println(err)
		return false, err
	}

	return true, nil
}

// SearchTutorial : search tutorials in elasticsearch
func SearchTutorial(c *gin.Context) {

	client, err := elastic.NewClient(elastic.SetURL("http://192.168.1.72:9200"))
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{
			"status":       "failure",
			"errorCode":    10009000,
			"errorMessage": "connection fail",
		})
		return
	}
	log.Println(client)

	exists, err := client.IndexExists("tutorial").Do(context.Background())
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{
			"status":       "failure",
			"errorCode":    10009001,
			"errorMessage": "execution fail",
		})
		return
	}
	if !exists {
		log.Println("index not exist")
		c.AbortWithStatusJSON(500, gin.H{
			"status":       "failure",
			"errorCode":    10009002,
			"errorMessage": "index not exist",
		})
		return
	}

	log.Println(exists)

	// get1, err := client.Get().Index("tutorial").Type("tutorial").Id("126").Do(context.Background())
	// if err != nil {
	// 	switch {
	// 	case elastic.IsNotFound(err):
	// 		panic(fmt.Sprintf("Document not found: %v", err))
	// 	case elastic.IsTimeout(err):
	// 		panic(fmt.Sprintf("Timeout retrieving document: %v", err))
	// 	case elastic.IsConnErr(err):
	// 		panic(fmt.Sprintf("Connection problem: %v", err))
	// 	default:
	// 		// Some other kind of error
	// 		panic(err)
	// 	}
	// }

	// var tes domain.TutorialES
	// json.Unmarshal(*get1.Source, &tes)

	log.Println("search part")

	matchPhaseQuery := elastic.NewMatchPhraseQuery("title", "教程")
	// termQuery := elastic.NewTermQuery("title", "教程")
	searchResult, err := client.Search().
		Index("tutorial").
		Query(matchPhaseQuery).
		Sort("createTime", false).
		From(0).
		Size(10).
		Do(context.Background())
	if err != nil {
		// Handle error
		log.Println(err)
		// panic(err)
	}

	// var ttyp domain.TutorialES
	// for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
	// 	tes := item.(domain.TutorialES)
	// }

	tutorialSlice := make([]domain.Tutorial, len(searchResult.Hits.Hits))
	index := 0
	if searchResult.Hits.TotalHits > 0 {
		var esr *domain.TutorialES
		for _, hit := range searchResult.Hits.Hits {
			err := json.Unmarshal(*hit.Source, &esr)
			if err != nil {
				// error
				log.Println(err)
			}
			tutorialSlice[index] = esr.ToTutorial()
			index++
		}
	}

	c.JSON(200, gin.H{
		"status": "success",
		"result": tutorialSlice,
		"total":  searchResult.Hits.TotalHits,
		"size":   index,
	})
	return
}
