package service

import (
	"strings"

	"../config"
	"../constant"
	"../domain"
	ApiErr "../error"
	"../request"
	"../utils"
	"github.com/gin-gonic/gin"
)

// Login : user login
func Login(c *gin.Context) {
	var req request.LoginRequest
	c.Bind(&req)

	if req.Username == "" || req.Password == "" {
		c.AbortWithStatusJSON(400, ApiErr.ErrRequestParam)
		return
	}

	// do not have the varify username when login
	// if ok, _ := regexp.MatchString(constant.UsernameRegex, req.Username); !ok {
	// 	c.AbortWithStatusJSON(400, ApiErr.ErrRequestParam)
	// 	return
	// }

	// fmt.Println("user input", req)

	db := config.GetDBConn()

	rows, err := db.Query(
		strings.Join([]string{"select ac.agent_id, ac.username, ac.password from agent ag join account ac on ag.id = ac.agent_id",
			"where ac.username =?",
			"and ac.is_del = 0 and ag.is_del = 0 and ac.ban = 0 and ag.ban = 0"}, " "), req.Username)
	defer rows.Close()

	if err != nil {
		c.AbortWithStatusJSON(500, ApiErr.ErrSQLExec)
		return
	}

	// var account Account
	var (
		agentID  int32
		username string
		password string
		account  domain.Account
	)

	rows.Next()
	if err = rows.Scan(&agentID, &username, &password); err != nil {
		c.AbortWithStatusJSON(500, ApiErr.ErrSQLScan)
		return
	}

	account = domain.Account{
		AgentID:  agentID,
		Username: username,
		Password: password,
	}

	// fmt.Printf("%v's password is %v\n", account.Username, account.Password)

	token, err := utils.CreateToken(account.AgentID)
	if err != nil {
		c.AbortWithStatusJSON(500, ApiErr.ErrTokenGen)
		return
	}

	if utils.VerifyPassword(req.Password, account.Password) {

		rc := config.GetRedisPool().Get()
		defer rc.Close()

		if _, err := rc.Do("SETEX", account.AgentID, constant.TokenExpireTime, token); err != nil {
			c.AbortWithStatusJSON(500, ApiErr.ErrRedisExec)
			return
		}

		c.JSON(200, gin.H{
			"status": "success",
			"result": token,
		})
	} else {
		c.AbortWithStatusJSON(400, ApiErr.ErrLogin)
		return
	}

}
