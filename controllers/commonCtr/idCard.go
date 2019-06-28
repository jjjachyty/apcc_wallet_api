package commonCtr

import (
	"apcc_wallet_api/services/commonSrv"
	"apcc_wallet_api/utils"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

// 身份识别
func IDCardRecognition(c *gin.Context) {

	fh1, err1 := c.FormFile("card1")
	fh2, err2 := c.FormFile("card2")
	if err1 == nil && err2 == nil {
		file1, err1 := fh1.Open()

		file2, err2 := fh2.Open()
		if err1 == nil && err2 == nil {
			defer func() {

				file1.Close()
				file2.Close()
			}()
			if card, err := commonSrv.IDCadrPostRecognition(file1); err == nil && card.ErrorMessage == "" {
				if back, err := commonSrv.IDCadrPostRecognition(file2); err == nil && back.ErrorMessage == "" {
					card.Cards[0].IssuedBy = back.Cards[0].IssuedBy
					card.Cards[0].ValidDate = back.Cards[0].ValidDate
					fmt.Println(back.Cards[0])
					utils.Response(c, err, card)
					return
				}
				utils.Response(c, errors.New("身份证反面解析失败"), nil)
				return

			}
			utils.Response(c, errors.New("身份证正面解析失败"), nil)

		} else {
			utils.Response(c, errors.New("图片打开失败"), nil)
			return
		}

	}

	utils.Response(c, errors.New("身份证正反面缺失"), nil)
}
