package commonSrv

import (
	"apcc_wallet_api/utils"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

const (
	REGISTER_TPID = "a061eb95f24441f99ea611d0bb3ca4ca" //注册验证码模板ID
)

//必填,请参考"开发准备"获取如下数据,替换为实际值
var realURL = "https://api.rtc.huaweicloud.com:10443/sms/batchSendSms/v1" //APP接入地址+接口访问URI
var appKey = "b9JGuVoNnD529tLfu7Ui3x48LsHq"                               //APP_Key
var appSecret = "HBkCrHvT71sguLt2nLtF9GQZ73F0"                            //APP_Secret
var sender = "8819081219913"                                              //国内短信签名通道号或国际/港澳台短信通道号
var signature = "register"                                                //签名名称

type SMSService struct{}

func (SMSService) SendSMS(phone string) error {
	rand := utils.CreateRandomNumber(4)
	if sms, err := utils.Get(phone); err.Error() == "redis: nil" && sms == "" {

		if err = utils.Set(phone, rand, time.Second*60); err == nil {
			return send(phone, REGISTER_TPID, []string{rand})
		}

	}
	return errors.New("1分钟内只允许发送一次短信")

}

func (SMSService) VerificationSMS(phone, sms string) error {
	if smsRedis, err := utils.Get(phone); err == nil && sms == smsRedis {
		return nil
	}
	return errors.New("短信校验失败")
}

//发送短息
func send(to string, templateID string, templateParas []string) error {
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	},
	}
	v := url.Values{}
	templateParasByts, _ := json.Marshal(templateParas)
	v.Set("from", sender)
	v.Add("to", to)
	v.Add("templateId", templateID)
	v.Add("templateParas", string(templateParasByts))
	v.Add("statusCallback", "")
	v.Add("signature", signature)

	var err error
	var resp *http.Response
	var req *http.Request
	var body []byte
	if req, err = http.NewRequest("POST", realURL, strings.NewReader(v.Encode())); err == nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Authorization", "WSSE realm=\"SDP\",profile=\"UsernameToken\",type=\"Appkey\"")
		xwsse := buildWsseHeader(appKey, appSecret)
		fmt.Println("xwsse", xwsse)
		req.Header.Set("X-WSSE", xwsse)

		if resp, err = client.Do(req); err == nil {
			defer resp.Body.Close()
			if body, err = ioutil.ReadAll(resp.Body); err == nil {
				fmt.Println(string(body))
				code := gjson.Get(string(body), "code").String()
				if code == "000000" {
					return nil
				}
				return fmt.Errorf("发送短信消息失败 code=%s", code)
			}

		}
	}
	return err

}

func buildWsseHeader(appKey, appSecret string) string {

	var created = time.Now().Local().UTC().Format("2006-01-02T15:04:05Z")
	var nonce = strings.ReplaceAll(utils.GetUUID(), "-", "")
	hash := sha256.New()
	hash.Write([]byte(nonce + created + appSecret))

	return fmt.Sprintf("UsernameToken Username=%q,PasswordDigest=%q,Nonce=%q,Created=%q", appKey, base64.StdEncoding.EncodeToString(hash.Sum(nil)), nonce, created)
}
