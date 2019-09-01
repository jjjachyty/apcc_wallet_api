package commonCtr

import (
	"apcc_wallet_api/middlewares/jwt"
	"apcc_wallet_api/services/commonSrv"
	"apcc_wallet_api/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

type NewsController struct{}

var newsService commonSrv.NewsService

//AssetsLogs 获取我的转账记录
func (NewsController) NewsList(c *gin.Context) {
	var err error
	var news = new(commonSrv.News)

	var page = utils.GetPageData(c)
	err = newsService.GetNewsList(page, news)
	utils.Response(c, err, page)
}

func (NewsController) NewsDetail(c *gin.Context) {
	var err error
	uuid, hasid := c.GetQuery("uuid")
	var news = new(commonSrv.News)
	if hasid {
		news.UUID = uuid
		err = newsService.GetNewsDetail(news)
	}
	utils.Response(c, err, news)
}

func (NewsController) AddOrUpdateNews(c *gin.Context) {
	var err error
	var news = new(commonSrv.News)
	if err = c.BindJSON(news); err == nil {
		c := jwt.GetClaims(c)
		fmt.Println("\n\n", news.State)
		if news.UUID != "" {
			news.UpdateBy = c.UUID
			err = newsService.UpdateNews(news)
		} else {
			news.UUID = utils.GetUUID()
			news.CreateBy = c.UUID
			news.State = utils.STATE_ENABLE
			err = newsService.AddNews(news)

		}
	}
	utils.Response(c, err, news)
}

func (NewsController) RemoveNews(c *gin.Context) {
	var err error
	uuid, hasid := c.GetQuery("uuid")

	if hasid {
		err = newsService.DeleteNews(commonSrv.News{UUID: uuid})
	}
	utils.Response(c, err, nil)
}
