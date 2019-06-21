package utils

import (
	
	"strconv"

	"github.com/gin-gonic/gin"
)

type Page struct {
	TotalRows   int64 //总记录数
	PageSize    int   //页数量
	PageCount   int   //总页数
	CurrentPage int   //当前页数
}

type PageData struct {
	Page Page
	Rows interface{}
}

func GetPageData(c *gin.Context) *PageData {
	var pageData = new(PageData)
	pageSize, _ := strconv.ParseInt(c.GetHeader("X-PAGE-SIZE"), 10, 64)
	currentPage, _ := strconv.ParseInt(c.GetHeader("X-PAGE-CURRENT"), 10, 64)
	if 0 == pageSize {
		pageData.Page.PageSize = 15
	} else {
		pageData.Page.PageSize = int(pageSize)
	}
	if 0 == currentPage {
		pageData.Page.CurrentPage = 1
	} else {
		pageData.Page.CurrentPage = int(currentPage)

	}
	return pageData
}
