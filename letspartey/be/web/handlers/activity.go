package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"iamlinix.com/partay/db"
	"iamlinix.com/partay/logger"
	"iamlinix.com/partay/orm/models"
	"iamlinix.com/partay/web"
)

type SingleActivity struct {
	ID        int        `json:"id"`
	Creator   string     `json:"creator"`
	Name      string     `json:"name"`
	StartTime *time.Time `json:"startTime"`
	Desc      string     `json:"desc"`
	Images    []string   `json:"images"`
}

type CreateActivityRsp struct {
	web.BaseResponse
	ID int `json:"id"`
}

type ActivityListRsp struct {
	web.BaseResponse
	Activities []*SingleActivity `json:"activities"`
}

func HdlrCreateActivity(c *gin.Context) {
	var activity models.Activity
	if err := UnmarshalRequestData(c, &activity); err != nil {
		logger.Errorf("failed to parse request data: %#v", err)
		return
	}

	if _, err := db.Get().Execute("INSERT INTO activities (creator, `name`, `desc`, start_time) VALUES (?, ?, ?, ?)",
		activity.Creator, activity.Name, activity.Desc, activity.StartTime); err != nil {
		logger.Errorf("failed to create new activity: %#v", err)
		c.JSON(http.StatusInternalServerError, &web.BaseResponse{
			Code:    web.ECGenUnknown,
			Message: web.EMGenUnknown,
		})
	} else {
		res, _ := db.Get().Execute("SELECT LAST_INSERT_ID()")
		c.JSON(http.StatusOK, &CreateActivityRsp{
			web.BaseResponse{
				Code:    web.ECOK,
				Message: web.EMOK,
			},
			int(*res[0]["LAST_INSERT_ID()"].(*uint64)),
		})
	}
}

func HdlrListActivity(c *gin.Context) {
	var p = c.Query("page")
	var ps = c.Query("querySize")
	var page int = 1
	var pageSize int = 20
	var err error

	if len(p) > 0 {
		page, err = strconv.Atoi(p)
		if err != nil {
			logger.Errorf("invalid page:", p)
			c.JSON(http.StatusBadRequest, &web.BaseResponse{
				Code:    web.ECGenIncorrectBody,
				Message: web.EMGenIncorrectBody,
			})
			return
		}
	}

	if len(ps) > 0 {
		pageSize, err = strconv.Atoi(p)
		if err != nil {
			logger.Errorf("invalid page size:", ps)
			c.JSON(http.StatusBadRequest, &web.BaseResponse{
				Code:    web.ECGenIncorrectBody,
				Message: web.EMGenIncorrectBody,
			})
			return
		}
	}

	results, err := db.Get().Execute("SELECT a.id, a.name, a.creator, a.start_time, a.desc, b.url FROM "+
		"(SELECT * FROM activities ORDER BY ctime DESC LIMIT ? OFFSET ?) a "+
		"LEFT JOIN activity_images b on a.id = b.activity_id ORDER BY a.id DESC", pageSize, (page-1)*pageSize)
	if err != nil {
		logger.Errorf("failed to query db for activities: %#v", err)
		c.JSON(http.StatusInternalServerError, &web.BaseResponse{
			Code:    web.ECGenUnknown,
			Message: web.EMGenUnknown,
		})
	} else {
		var acts map[int32]*SingleActivity = make(map[int32]*SingleActivity)
		var resp ActivityListRsp
		var act *SingleActivity
		var ok bool
		resp.Code = web.ECOK
		resp.Message = web.EMOK
		for _, r := range results {
			actID := *r["id"].(*int32)
			if act, ok = acts[actID]; !ok {
				act = &SingleActivity{
					ID:        int(actID),
					Creator:   string(*r["creator"].(*sql.RawBytes)),
					StartTime: &r["start_time"].(*mysql.NullTime).Time,
					Desc:      string(*r["desc"].(*sql.RawBytes)),
					Name:      string(*r["name"].(*sql.RawBytes)),
				}
				acts[actID] = act
				resp.Activities = append(resp.Activities, act)
			}

			url := string(*r["url"].(*sql.RawBytes))
			if len(url) > 0 {
				act.Images = append(act.Images, url)
			}
		}
		c.JSON(http.StatusOK, &resp)
	}
}
