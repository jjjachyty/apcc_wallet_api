package utils

import (
	"bytes"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Page struct {
	OrderBy     string
	Sort        string
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
	size, hasSize := c.GetQuery("size")
	page, hasPage := c.GetQuery("page")
	order, hasOrder := c.GetQuery("order")
	sort, hasSort := c.GetQuery("sort")

	pageSize, _ := strconv.ParseInt(size, 10, 64)
	currentPage, _ := strconv.ParseInt(page, 10, 64)
	if !hasSize {
		pageData.Page.PageSize = 15
	} else {
		pageData.Page.PageSize = int(pageSize)
	}
	if !hasPage {
		pageData.Page.CurrentPage = 1
	} else {
		pageData.Page.CurrentPage = int(currentPage)

	}
	if hasOrder && hasSort {

		order = Field2Cols(order)

		pageData.Page.OrderBy = order
		pageData.Page.Sort = sort
	}

	return pageData

}

func Field2Cols(filed string) string {
	byts := []byte(filed)
	col := bytes.NewBuffer(nil)
	for i, v := range byts {
		if (int(v) > 64 && int(v) < 91) && i > 0 { //大写字母
			col.WriteByte('_')
		}
		col.WriteByte(v)
	}
	return col.String()
}
