package commonSrv

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

type response struct {
	ImageId   string
	RequestId string
	Cards     []card
}

type card struct {
	Gender       string
	Name         string
	IdCardNumber string
	Birthday     string
	Race         string
	Address      string
	// Legality": {
	// 	"Edited": 0.001,
	// 	"Photocopy": 0.0,
	// 	"ID Photo": 0.502,
	// 	"Screen": 0.496,
	// 	"Temporary ID Photo": 0.0
	// },
	Type      string
	Side      string
	IssuedBy  string
	ValidDate string
}

var idcard_url = "https://api-cn.faceplusplus.com/cardpp/v1/ocridcard"
var api_key = "HwdG52YWaf8tgX-FNwyHY68jKcROhjx7"
var api_secret = "o00516a7fjOAP3CoofuzSGy60jxPXJcf"

func IDCadrPostRecognition(imageFile multipart.File) (*card, error) {
	var cd = new(card)
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	bodyWriter.WriteField("api_key", api_key)
	bodyWriter.WriteField("api_secret", api_secret)
	//关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile("image_file", "image_file")
	if err != nil {
		fmt.Println("error writing to buffer")
		return nil, err
	}

	//iocopy
	_, err = io.Copy(fileWriter, imageFile)
	if err != nil {
		return nil, err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(idcard_url, contentType, bodyBuf)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(resp.Status)
	fmt.Println(string(resp_body))
	err = json.Unmarshal(resp_body, cd)
	return cd, err
}
