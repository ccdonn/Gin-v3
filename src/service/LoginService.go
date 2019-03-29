package service

import (
	"regexp"
	"strings"

	"../config"
	"../constant"
	"../domain"
	"../request"
	"../utils"
	"github.com/gin-gonic/gin"
)

// Login : user login
func Login(c *gin.Context) {
	var req request.LoginRequest
	c.Bind(&req)

	if req.Username == "" || req.Password == "" {
		c.JSON(400, gin.H{
			"status":       "failure",
			"errorCode":    123,
			"errorMessage": "no input",
		})
		return
	}

	if ok, _ := regexp.MatchString(constant.UsernameRegex, req.Username); !ok {
		c.JSON(400, gin.H{
			"status":       "failure",
			"errorCode":    123,
			"errorMessage": "username ",
		})
	}

	// fmt.Println("user input", req)

	db := config.GetDBConn()
	rows, err := db.Query(
		strings.Join([]string{"select ac.agent_id, ac.username, ac.password from agent ag join account ac on ag.id = ac.agent_id",
			"where ac.username = \"" + req.Username + "\"",
			"and ac.is_del = 0 and ag.is_del = 0 and ac.ban = 0 and ag.ban = 0"}, " "))

	if err != nil {
		c.JSON(500, gin.H{
			"status":       "failure",
			"errorCode":    10001012,
			"errorMessage": "internal sql error",
		})
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
		c.JSON(500, gin.H{
			"status":       "failure",
			"errorCode":    10001013,
			"errorMessage": "internal scan error",
		})
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
		c.JSON(500, gin.H{
			"status":       "failure",
			"errorCode":    10001035,
			"errorMessage": "create token fail",
		})
	}

	if utils.VerifyPassword(req.Password, account.Password) {

		rc := config.GetRedisPool().Get()
		// rc.Do("SET", "kk", 0)
		if _, err := rc.Do("SETEX", account.AgentID, constant.TokenExpireTime, token); err != nil {
			c.JSON(500, gin.H{
				"status":       "failure",
				"errorCode":    10001025,
				"errorMessage": "login fail",
			})
		}

		defer rc.Close()

		c.JSON(200, gin.H{
			"status": "success",
			"result": token,
		})
	} else {
		c.JSON(400, gin.H{
			"status":       "failure",
			"errorCode":    10001033,
			"errorMessage": "login fail",
		})
	}

}
