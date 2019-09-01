package commonSrv

import (
	"apcc_wallet_api/models"
	"apcc_wallet_api/utils"
	"time"
)

type NewsService struct{}
type News struct {
	UUID string `xorm:"varchar(36) notnull unique pk 'uuid'"`

	Title    string
	Banner   string
	Content  string
	State    int
	CreateAt time.Time `xorm:"created"`
	CreateBy string
	UpdateAt time.Time `xorm:"updated"`
	UpdateBy string
}

func (NewsService) GetNewsList(page *utils.PageData, condBean interface{}) error {
	return models.GetBeansPage(page, condBean)
}

func (NewsService) GetNewsDetail(news *News) error {
	return models.GetBean(news)
}

func (NewsService) AddNews(news *News) error {
	return models.Create(news)
}
func (NewsService) UpdateNews(news *News) error {
	return models.Update(news.UUID, news, "state", "title", "content", "banner")
}
func (NewsService) DeleteNews(news News) error {
	return models.Delete(news.UUID, news)
}
