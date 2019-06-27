package jwt

import (
	"apcc_wallet_api/models/userMod"
	"apcc_wallet_api/utils"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// JWTAuth 中间件，检查token
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("authorization")
		if token == "" {
			c.JSON(http.StatusOK, gin.H{
				"Status":  false,
				"Message": "请求未携带token，无权限访问",
			})
			c.Abort()
			return
		}

		log.Print("get token: ", token)

		j := NewJWT()
		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == TokenExpired {
				c.JSON(http.StatusUnauthorized, gin.H{
					"Status":  false,
					"Message": "授权已过期",
				})
				c.Abort()
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"Status":  false,
				"Message": err.Error(),
			})
			c.Abort()
			return
		}
		// 继续交由下一个路由处理,并将解析出的信息传递下去
		c.Set("claims", claims)
	}
}

// JWT 签名结构
type JWT struct {
	SigningKey []byte
}

// 一些常量
var (
	TokenExpired     error  = errors.New("Token is expired")
	TokenNotValidYet error  = errors.New("Token not active yet")
	TokenMalformed   error  = errors.New("That's not even a token")
	TokenInvalid     error  = errors.New("Couldn't handle this token:")
	SignKey          string = "(GoFuckyourSelf~!)"
)

// 载荷，可以加一些自己需要的信息
type CustomClaims struct {
	UUID         string `json:"uuid"`
	Phone        string `json:"phone"`
	HasPayPasswd bool   `json:"hasPayPasswd"`
	jwt.StandardClaims
}

// 新建一个jwt实例
func NewJWT() *JWT {
	return &JWT{
		[]byte(GetSignKey()),
	}
}

// 获取signKey
func GetSignKey() string {
	return SignKey
}

// 这是SignKey
func SetSignKey(key string) string {
	SignKey = key
	return SignKey
}

// CreateToken 生成一个token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// 解析Tokne
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}

// 更新token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, _ := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})

	if claims, ok := token.Claims.(*CustomClaims); ok {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
}

// 生成令牌
func GenerateToken(user userMod.User) (string, error) {
	j := &JWT{
		[]byte(SignKey),
	}
	claims := CustomClaims{
		user.UUID,
		user.Phone,
		user.HasPayPasswd,
		jwt.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 1000), // 签名生效时间
			ExpiresAt: int64(time.Now().Unix() + 3600), // 过期时间 一小时
			Issuer:    "Janly",                         //签名的发行者
		},
	}
	return j.CreateToken(claims)
}

func GetClaims(c *gin.Context) *CustomClaims {
	return c.MustGet("claims").(*CustomClaims)
}

// GetDataByTime 一个需要token认证的测试接口
func GetDataByTime(c *gin.Context) {
	claims := c.MustGet("claims").(*CustomClaims)
	if claims != nil {
		c.JSON(http.StatusOK, gin.H{
			"Status":  true,
			"Message": "token有效",
			"Data":    claims,
		})
	}
}

// GetDataByTime
func RefreshToken(c *gin.Context) {
	var err = errors.New("authorization 不能为空")
	token := c.Request.Header.Get("authorization")
	fmt.Println("\n")
	fmt.Println(NewJWT().ParseToken(token))
	fmt.Println("\n")
	if token != "" {
		token, err = NewJWT().RefreshToken(token)
		fmt.Println("newtoken=", token, err)
	}

	utils.Response(c, err, gin.H{"Token": token})
}
