package domain

import (
	"fmt"
	"log"
	"time"
)

type TutorialES struct {
	ID       int32  `json:"id"`
	Title    string `json:"title"`
	TitleImg string `json:"titleImg"`
	Content  string `json:"content"`
	// CreateTimeValue     int64  `json:"createTime"`
	Del            int8  `json:"del"`
	LastUpdateUser int32 `json:"lastUpdateUser"`
	// LastUpdateTimeValue int64  `json:"lastUpdateTime"`

	// TopAndOrder    int

	CreateTime     string `json:"createTime"`
	LastUpdateTime string ` json:"lastUpdateTime"`
}

func (tes TutorialES) ToTutorial() Tutorial {
	t := Tutorial{
		ID:             tes.ID,
		Title:          tes.Title,
		TitleImg:       tes.TitleImg,
		Content:        tes.Content,
		Del:            tes.Del,
		LastUpdateUser: tes.LastUpdateUser,
	}

	// time should be
	layout := "2006-01-02T15:04:05.999Z"

	if tes.CreateTime != "" {
		fmt.Println(tes.CreateTime)
		createT, err := time.Parse(layout, tes.CreateTime)
		fmt.Println(createT)
		if err != nil {
			log.Println(err)
		}
		t.CreateTimeValue = createT.Unix()
	}

	if tes.LastUpdateTime != "" {
		fmt.Println(tes.LastUpdateTime)
		lastUpdateT, err := time.Parse(layout, tes.LastUpdateTime)
		fmt.Println(lastUpdateT)
		if err != nil {
			log.Println(err)
		}
		t.LastUpdateTimeValue = lastUpdateT.Unix()
	}

	return t
}
