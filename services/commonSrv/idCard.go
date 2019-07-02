package commonSrv

import (
	"apcc_wallet_api/models/userMod"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

type response struct {
	ImageId      string
	RequestId    string
	Cards        []userMod.IdCard
	TimeUsed     int
	ErrorMessage string `json:"error_message"`
}

var idcard_url = "https://api-cn.faceplusplus.com/cardpp/v1/ocridcard"
var api_key = "HwdG52YWaf8tgX-FNwyHY68jKcROhjx7"
var api_secret = "o00516a7fjOAP3CoofuzSGy60jxPXJcf"

func IDCadrPostRecognition(imageFile multipart.File) (*response, error) {
	var ret = new(response)
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	bodyWriter.WriteField("api_key", api_key)
	bodyWriter.WriteField("api_secret", api_secret)
	//关键的一步操作
	imageFileWriter, err := bodyWriter.CreateFormFile("image_file", "front")
	if err != nil {
		fmt.Println("error writing to buffer")
		return nil, err
	}

	//iocopy
	_, err = io.Copy(imageFileWriter, imageFile)
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
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(respBody, ret)
	return ret, err
}
