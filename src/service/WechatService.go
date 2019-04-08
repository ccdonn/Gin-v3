package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"../config"
	"../constant"
	ApiErr "../error"
	"../request"
	"../vo"
	"github.com/gin-gonic/gin"
)

// GetWechat : get wechat account
func GetWechat(c *gin.Context) {
	var query request.WechatRequest
	c.Bind(&query)

	fmt.Println(query)

	response, err := http.Post("http://192.168.1.72:8083/simulator/wc/getPromoteWechat", "", nil)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{
			"status": "failure",
		})
		return
	}

	data, _ := ioutil.ReadAll(response.Body)
	var s vo.ExternalWechatMessage
	// fmt.Println(string(data))
	err = json.Unmarshal(data, &s)

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{
			"status": "failure",
		})
		return
	}

	// fmt.Println(s)

	if !s.Suc {
		// log.Println()
		c.AbortWithStatusJSON(500, gin.H{
			"status": "failure",
		})
		return
	}

	all := make(map[string][]string)
	for _, ewvo := range s.Data {
		if ewvo.Value != "" {
			all[ewvo.Channel] = strings.Split(ewvo.Value, "|")
		} else {
			all[ewvo.Channel] = []string{}
		}
	}

	// refactoring :: could it better ?
	r := make([]*vo.WechatVO, 1)
	if query.BrandKey == "" {
		if query.Wechat != "" {
			for k, vs := range all {
				for _, v := range vs {
					if v == query.Wechat {
						r[0] = &vo.WechatVO{
							BrandKey:      k,
							WechatAccount: vs,
						}
					}
				}
			}
		}
	} else if values, exist := all[query.BrandKey]; exist {
		if query.Wechat != "" {
			for _, v := range values {
				if v == query.Wechat {
					r[0] = &vo.WechatVO{
						BrandKey:      query.BrandKey,
						WechatAccount: values,
					}
				}
			}
		} else {
			r[0] = &vo.WechatVO{
				BrandKey:      query.BrandKey,
				WechatAccount: values,
			}
		}
	}

	if r[0] == nil {
		r = r[:0]
	}

	c.JSON(200, gin.H{
		"status":    "success",
		"timestamp": time.Now().Unix(),
		"result":    r,
	})
	return

}

// UpdateWechat : update wechat account
func UpdateWechat(c *gin.Context) {

}

// DeleteWechat : delete wechat account
func DeleteWechat(c *gin.Context) {

}

// GetWechatBrand : get wechat brand
func GetWechatBrand(c *gin.Context) {
	r, err := getWechatBrand()
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, ApiErr.ErrNotFound)
		return
	}

	c.JSON(200, gin.H{
		"status":    "success",
		"result":    r,
		"timestamp": time.Now(),
		"size":      len(r),
		"total":     len(r),
	})

	return

}

func getWechatBrand() ([]*vo.WechatBrandVO, error) {
	db := config.GetDBConn()

	row := db.QueryRow("select v from system_config where k = ?", "game.wechat")

	var value string

	err := row.Scan(&value)

	if err != nil {
		return nil, err
	}

	brandStrings := strings.Split(value, ",")
	count := strings.Count(value, ",") + 1
	wechatBrandVOs := make([]*vo.WechatBrandVO, count)
	index := 0
	for _, bStr := range brandStrings {
		brandInfo := strings.Split(bStr, ":")
		wechatBrandVOs[index] = &vo.WechatBrandVO{
			Key:  brandInfo[0],
			Name: brandInfo[1],
		}
		index++
	}

	return wechatBrandVOs, nil
}

func getWechatBrandMap() (map[string]string, error) {
	db := config.GetDBConn()

	row := db.QueryRow("select v from system_config where k = ?", "game.wechat")

	var value string

	err := row.Scan(&value)

	if err != nil {
		return nil, err
	}

	brandStrings := strings.Split(value, ",")
	count := strings.Count(value, ",") + 1
	wechatBrandVOMap := make(map[string]string, count)
	for _, bStr := range brandStrings {
		brandInfo := strings.Split(bStr, ":")
		wechatBrandVOMap[brandInfo[0]] = brandInfo[1]
	}

	return wechatBrandVOMap, nil
}

// CreateWechatBrand :
func CreateWechatBrand(c *gin.Context) {
	var body vo.WechatBrandVO
	c.Bind(&body)

	if body.Key == "" || body.Name == "" {
		// log.Println("")
		c.AbortWithStatusJSON(400, ApiErr.ErrRequestParam)
		return
	}

	if ok, err := createWechatBrand(&body); !ok {
		log.Println(err)
		c.AbortWithStatusJSON(500, ApiErr.ErrWechatBrandCreateFail)
	} else {
		c.JSON(200, gin.H{
			"status": "success",
		})
	}
	return
}

func createWechatBrand(wechatBrand *vo.WechatBrandVO) (bool, error) {

	wechatBrandMaps, err := getWechatBrandMap()
	if err != nil {
		log.Println(err)
		return false, err
	}

	if _, exist := wechatBrandMaps[wechatBrand.Key]; !exist {
		log.Println("key already exist")
		return false, nil // not exist
	}

	str := ""
	for k, v := range wechatBrandMaps {
		str += k + ":" + v + ","
	}

	str += wechatBrand.Key + ":" + wechatBrand.Name + ","
	// fmt.Println("debug:" + str)

	db := config.GetDBConn()
	create, err := db.Prepare("update system_config set v=?,update_time=now() where k=?")
	if err != nil {
		return false, err
	}
	_, err = create.Exec(str[:len(str)-1], constant.SysConfigWechatBrand)
	if err != nil {
		return false, err
	}

	return true, nil
}

// UpdateWechatBrand :
func UpdateWechatBrand(c *gin.Context) {

	n := c.Param("ID")
	var body vo.WechatBrandVO
	c.Bind(&body)
	body.Key = n

	if body.Key == "" || body.Name == "" {
		c.AbortWithStatusJSON(400, ApiErr.ErrRequestParam)
		return
	}

	if ok, err := updateWechatBrand(&body); !ok {
		log.Println(err)
		c.AbortWithStatusJSON(500, ApiErr.ErrWechatBrandUpdateFail)
	} else {
		c.JSON(200, gin.H{
			"status": "success",
		})
	}
	return
}

func updateWechatBrand(wechatBrand *vo.WechatBrandVO) (bool, error) {

	wechatBrandMaps, err := getWechatBrandMap()
	if err != nil {
		return false, err
	}

	if _, exist := wechatBrandMaps[wechatBrand.Key]; !exist {
		return false, nil
	}

	wechatBrandMaps[wechatBrand.Key] = wechatBrand.Name

	str := ""
	for k, v := range wechatBrandMaps {
		str += k + ":" + v + ","
	}

	db := config.GetDBConn()
	update, err := db.Prepare("update system_config set v=?,update_time=now() where k=?")
	if err != nil {
		return false, err
	}
	_, err = update.Exec(str[:len(str)-1], constant.SysConfigWechatBrand)
	if err != nil {
		return false, err
	}

	return true, nil
}

// DeleteWechatBrand :
func DeleteWechatBrand(c *gin.Context) {
	n := c.Param("ID")
	if n == "" {
		c.AbortWithStatusJSON(400, ApiErr.ErrRequestParam)
		return
	}

	if ok, err := deleteWechatBrand(n); !ok {
		log.Println(err)
		c.AbortWithStatusJSON(500, ApiErr.ErrWechatBrandDeleteFail)
	} else {
		c.JSON(200, gin.H{
			"status": "success",
		})
	}
	return
}

func deleteWechatBrand(key string) (bool, error) {
	wechatBrandMaps, err := getWechatBrandMap()
	if err != nil {
		return false, err
	}

	if _, exist := wechatBrandMaps[key]; !exist {
		return false, nil
	}

	delete(wechatBrandMaps, key)

	str := ""
	for k, v := range wechatBrandMaps {
		str += k + ":" + v + ","
	}

	db := config.GetDBConn()
	update, err := db.Prepare("update system_config set v=?,update_time=now() where k=?")
	if err != nil {
		return false, err
	}
	_, err = update.Exec(str[:len(str)-1], constant.SysConfigWechatBrand)
	if err != nil {
		return false, err
	}

	return true, nil
}
