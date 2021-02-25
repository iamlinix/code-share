package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"iamlinix.com/partay/db"
	"iamlinix.com/partay/logger"
	"iamlinix.com/partay/orm/models"
	"iamlinix.com/partay/web"
)

type SinglePost struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Creator string   `json:"creator"`
	Desc    string   `json:"desc"`
	Images  []string `json:"images"`
}

type CreatePostRsp struct {
	web.BaseResponse
	ID int `json:"id"`
}

type PostListRsp struct {
	web.BaseResponse
	Posts []*SinglePost `json:"posts"`
}

func HdlrCreatePost(c *gin.Context) {
	var post models.Post
	if err := UnmarshalRequestData(c, &post); err != nil {
		logger.Errorf("failed to parse request data: %#v", err)
		return
	}

	if _, err := db.Get().Execute("INSERT INTO posts (creator, `name`, `desc`) VALUES (?, ?, ?)",
		post.Creator, post.Name, post.Desc); err != nil {
		logger.Errorf("failed to create new post: %#v", err)
		c.JSON(http.StatusInternalServerError, &web.BaseResponse{
			Code:    web.ECGenUnknown,
			Message: web.EMGenUnknown,
		})
	} else {
		res, _ := db.Get().Execute("SELECT LAST_INSERT_ID()")
		c.JSON(http.StatusOK, &CreatePostRsp{
			web.BaseResponse{
				Code:    web.ECOK,
				Message: web.EMOK,
			},
			int(*res[0]["LAST_INSERT_ID()"].(*uint64)),
		})
	}
}

func HdlrListPost(c *gin.Context) {
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

	results, err := db.Get().Execute("SELECT a.id, a.name, a.creator, a.desc, b.url FROM "+
		"(SELECT * FROM posts ORDER BY ctime DESC LIMIT ? OFFSET ?) a "+
		"LEFT JOIN post_images b on a.id = b.post_id ORDER BY a.id DESC", pageSize, (page-1)*pageSize)
	if err != nil {
		logger.Errorf("failed to query db for activities: %#v", err)
		c.JSON(http.StatusInternalServerError, &web.BaseResponse{
			Code:    web.ECGenUnknown,
			Message: web.EMGenUnknown,
		})
	} else {
		var posts map[int32]*SinglePost = make(map[int32]*SinglePost)
		var resp PostListRsp
		var post *SinglePost
		var ok bool
		resp.Code = web.ECOK
		resp.Message = web.EMOK
		for _, r := range results {
			posID := *r["id"].(*int32)
			if post, ok = posts[posID]; !ok {
				post = &SinglePost{
					ID:      int(posID),
					Creator: string(*r["creator"].(*sql.RawBytes)),
					Desc:    string(*r["desc"].(*sql.RawBytes)),
					Name:    string(*r["name"].(*sql.RawBytes)),
				}
				posts[posID] = post
				resp.Posts = append(resp.Posts, post)
			}

			url := string(*r["url"].(*sql.RawBytes))
			if len(url) > 0 {
				post.Images = append(post.Images, url)
			}
		}
		c.JSON(http.StatusOK, &resp)
	}
}
