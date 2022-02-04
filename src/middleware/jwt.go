package middleware

import (
	"NAS/src/sqlmap"
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

const TokenTimout = time.Hour

//const UpdateToken = time.Minute * 30
//const RefTokenTimeout = time.Hour * 2

var (
	ctx              = sqlmap.Ctx
	rdb              = sqlmap.RedisClient2
	db               = sqlmap.MyDb
	TokenExpired     = errors.New("token Is Expired")
	TokenNotValidYet = errors.New("token Not Active Yet")
	TokenMalformed   = errors.New("that's Not Even A Token")
	TokenInvalid     = errors.New("couldn't Handle This Token")
	key              = []byte("WebKey-MyTest-adminA1")
)

var pemKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIICXgIBAAKBgQDDcl4ekDzG9baDn47qB7EghIsSgfY6uXNHmiQ8h+OaDAcDa8Rc
UNtg772EXzOOxjaDFSV/JWMQFSVXRh2Q5xoEtzDx/rh3xd8Kc5T2GWH+ym1psk4s
RR7SFKzL82LzfXvTm+HF6fhHK9SY71sO9zcxhbAhm0RSKy23dhLj5woz+QIDAQAB
AoGBAJOaTcZbO+suKfZhi/bmdDiQoM8LYz+aSptqp68nGRZ/utQ0kQj+747Xv5K2
qyNKQmTglX7eZ1//+EFe7HlAbv6KkdOJ2EOAF5V8ySKytOBA7DPMKle/PMpjLUWE
ZfLpAHo5sC24AhISIf/CMHl7TTra8TPE/LMtbuGsHmr1r54BAkEA6MrL32/eZqxw
4Y+W37frLgQHeG7iPwHvGZofL8LPU8OdCjnnziMQKbEvBh6qp3pCtO+bDTYDaPXW
c2JdHKPRGQJBANbuddcO2R0uvlVgkCVmvPG40RvfNPQ0Ei8JPKZvg7ObauQ8O/Rr
stqZv0ZDtQ2MjMnM1QS1AEh99V+/+dVXdeECQQCqpg7xiiY0ifBtyT7GXSKPpvB6
/n3nxlkqIWr/LgWh1+HE31HoMJfmmDZqfAyJnPxNet/kvVWemahNCSxMlGHxAkBw
w4tv2YpvlSanBJKcDNr0t1J+nQzbUrZ3lxELAVbH1LKwLCoIgrjDmAaShtNm2GbF
OYJJhe0wG2WxZrddBxYBAkEAgHjmOPVcfEZwpIOfdSvIrCdeKzIIA+jS+lwp4a4C
+C9cVpRCEuLyVW1g/zDL0xQ617X92ChMJOfooJ7VVKt1gw==
-----END RSA PRIVATE KEY-----
`)

type (
	logs struct {
		Username string `json:"Username"`
		Password string `json:"Password"`
	}

	newPass struct {
		Username string `json:"Username"`
		Password string `json:"Password"`
		NewPass  string `json:"NewPass"`
	}

	sqlUser struct {
		ID       int64  `gorm:"primary_key"`
		Username string `json:"Username" gorm:"index"`
		Password string `json:"Password"`
		Sing     string `json:"Sing"`
	}

	mess struct {
		Status int    `json:"Status"`
		Msg    string `json:"Msg"`
		Data   interface{}
	}

	CustomClaims struct {
		Name string `json:"username"`
		Uid  int64  `json:"uid"`
		Ip   string `json:"ip"`
		Sing string `json:"sing"`
		Auth string `json:"auth"`
		// StandardClaims结构体实现了Claims接口(Valid()函数)
		jwt.StandardClaims
	}

	LoginResult struct {
		Token string `json:"token"`
		// 用户模型
		Name string `json:"name"`
	}
)

func (sqlUser) TableName() string {
	return "user"
}

func InitUser() {
	db.AutoMigrate(&sqlUser{})
}

func New(c *gin.Context) {
	var now logs
	err := c.BindJSON(&now)
	if err != nil {
		c.JSON(http.StatusBadRequest, mess{
			Status: -1,
			Msg:    "数据解析失败",
			Data:   nil,
		})
		return
	}
	if getUser(now.Username) {
		err = newUser(now.Username, md5V(now.Password))
		if err != nil {
			log.Println("New", err.Error())
			c.JSON(http.StatusBadRequest, mess{
				Status: -1,
				Msg:    "注册失败",
				Data:   nil,
			})
			return
		}
		c.JSON(http.StatusOK, mess{
			Status: 0,
			Msg:    "注册成功",
			Data:   nil,
		})
		return
	} else {
		c.JSON(http.StatusBadRequest, mess{
			Status: -1,
			Msg:    "用户已存在",
			Data:   nil,
		})
		return
	}

}

func Login(c *gin.Context) {
	var now logs
	err := c.BindJSON(&now)
	if err != nil {
		c.JSON(http.StatusBadRequest, mess{
			Status: -1,
			Msg:    "数据解析失败",
			Data:   nil,
		})
		return
	}
	ip := c.ClientIP()
	claim, bol := getPass(now.Username, md5V(now.Password))
	if bol {
		claims, bla := MyClaims(claim.Username, claim.ID, ip, sing(claim.Sing), TokenTimout)
		if bla {
			c.JSON(http.StatusOK, mess{
				Status: 0,
				Msg:    "ok",
				Data:   claims,
			})
		} else {
			c.JSON(http.StatusOK, mess{
				Status: 0,
				Msg:    "错误",
				Data:   claims,
			})
		}
	} else {
		c.JSON(http.StatusOK, mess{
			Status: -1,
			Msg:    "密码错误",
			Data:   nil,
		})
	}
}

func RestPass(c *gin.Context) {
	var now newPass
	err := c.BindJSON(&now)
	if err != nil {
		c.JSON(http.StatusBadRequest, mess{
			Status: -1,
			Msg:    "数据解析失败",
			Data:   nil,
		})
		return
	}
	token := c.Request.Header.Get("token")
	parserToken, _ := ParserToken(token)
	if parserToken.Name != now.Username {
		c.JSON(http.StatusUnauthorized, mess{
			Status: -1,
			Msg:    "权限不足",
			Data:   nil,
		})
		return
	}
	val := restPas(parserToken.Name, md5V(now.Password))
	if val {
		err = upUser(parserToken.Name, md5V(now.NewPass))
		if err != nil {
			log.Println("UpUser:", err.Error())
			c.JSON(http.StatusUnauthorized, mess{
				Status: -1,
				Msg:    "未知错误",
				Data:   nil,
			})
			return
		}
	} else {
		c.JSON(http.StatusUnauthorized, mess{
			Status: -1,
			Msg:    "用户名或密码错误",
			Data:   nil,
		})
		return
	}
	c.JSON(http.StatusOK, mess{
		Status: -1,
		Msg:    "修改成功",
		Data:   nil,
	})
}

// Logout 当使用Redis做验证的时候可以进行退出
func Logout(c *gin.Context) {
	token := c.Request.Header.Get("token")
	parserToken, _ := ParserToken(token)
	rdb.Del(ctx, parserToken.Name)
	c.JSON(http.StatusOK, mess{
		Status: 0,
		Msg:    "已退出",
		Data:   nil,
	})
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, mess{
				Status: -1,
				Msg:    "无权限访问",
				Data:   nil,
			})
			c.Abort()
			return
		}
		parserToken, err := ParserToken(token)
		if err != nil {
			// token过期
			if err == TokenExpired {
				c.JSON(http.StatusUnauthorized, mess{
					Status: -1,
					Msg:    "token授权已过期，请重新申请授权",
					Data:   nil,
				})
				c.Abort()
				return
			}
			// 其他错误
			c.JSON(http.StatusUnauthorized, mess{
				Status: -1,
				Msg:    err.Error(),
				Data:   nil,
			})
			c.Abort()
			return
		}
		if parserToken.Ip != c.ClientIP() || parserToken.Sing != getSing(parserToken.Name) {
			c.JSON(http.StatusUnauthorized, mess{
				Status: -1,
				Msg:    "验证token失败",
				Data:   nil,
			})
			c.Abort()
			return
		}
		c.Set("uid", parserToken.Uid)
		c.Next()
	}
}

//func JWTUpd() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		refToken := c.Request.Header.Get("refToken")
//		if refToken == "" {
//			getTime := c.GetInt64("ExpiresAt")
//			tokenTime := time.Unix(int64(getTime), 0).Unix()
//			newTime := time.Unix(time.Now().Unix(), 0).Unix()
//			remainSecond := time.Duration(newTime-tokenTime) * time.Second
//			log.Println(remainSecond.String())
//			if remainSecond < time.Duration(time.Minute*10) {
//				log.Println("????")
//				c.Set("TokenTimout", true)
//				c.Next()
//				return
//			}
//			c.Set("TokenTimout", false)
//			c.Next()
//			return
//		}
//		parserToken, err := ParserToken(refToken)
//		if err != nil {
//			// token过期
//			if err == TokenExpired {
//				c.JSON(http.StatusUnauthorized, mess{
//					Status: -1,
//					Msg:    "token授权已过期，请重新申请授权",
//					Data:   nil,
//				})
//				c.Abort()
//				return
//			}
//			// 其他错误
//			c.JSON(http.StatusUnauthorized, mess{
//				Status: -1,
//				Msg:    err.Error(),
//				Data:   nil,
//			})
//			c.Abort()
//			return
//		}
//		if parserToken.Ip != c.ClientIP() || parserToken.Sing != getSing(parserToken.Name) || parserToken.Auth != "update" {
//			c.JSON(http.StatusUnauthorized, mess{
//				Status: -1,
//				Msg:    "验证token失败",
//				Data:   nil,
//			})
//			c.Abort()
//			return
//		}
//		c.Next()
//	}
//}

func getUser(user string) bool {
	var users sqlUser
	//err := db.Model(&sqlUser{}).Where("Username = ?", user).First(&users).Error
	err := db.Table("user").First(&users, "Username = ?", user).Error
	if err != nil {
		log.Println("GetUser", err.Error())
	}
	if users.Username == user {
		return false
	}
	return true
}

func getPass(user, pass string) (sqlUser, bool) {
	var lo sqlUser
	db.Table("user").First(&lo, "Username = ?", user)
	if lo.Password == pass && lo.Username == user {
		//err := rdb.Ping(ctx).Err()
		//log.Println(err.Error())
		rdb.Set(ctx, lo.Username, sing(lo.Sing), TokenTimout)
		return lo, true
		//myClaims(lo)
	}
	return lo, false
}

func restPas(user, pass string) bool {
	var res sqlUser
	db.Table("user").First(&res, "Username = ?", user)
	if res.Password == pass {
		go func() {
			time.Sleep(time.Second * 2)
			rdb.Del(ctx, user)
		}()
		return true
	}
	return false
}

func newUser(user, pass string) error {
	var now sqlUser
	now = sqlUser{
		Username: user,
		Password: pass,
		Sing:     rsaSing(pass),
	}
	err := db.Table("user").Create(&now).Error
	return err
}

func upUser(user, pass string) error {
	var upd sqlUser
	upd = sqlUser{Sing: rsaSing(pass), Password: pass}
	err := db.Table("user").First(&upd, "Username = ?", user).Update(&upd).Error
	return err
}

//如果使用mysql则Token不会主动失效
func getSing(user string) string {
	val := rdb.Get(ctx, user).Val()
	return val
	//var sin sqlUser
	//db.Table("user").First(&sin, "Username = ?", user)
	//return sin.Sing
}

func MyClaims(username string, uid int64, ip string, sing string, time time.Duration) (interface{}, bool) {
	claims := CustomClaims{
		Name: username,
		Uid:  uid,
		Ip:   ip,
		Sing: sing,
		StandardClaims: jwt.StandardClaims{
			//NotBefore: int64(time.Now().Unix() - 1000), // 签名生效时间
			//ExpiresAt: int64(time.Now().Unix() + 3600), // 签名过期时间
			ExpiresAt: int64(time),
			Issuer:    "wang.com", // 签名颁发者
		},
	}

	token, err := createToken(claims)
	if err != nil {
		log.Println("Create:", err.Error())
		return nil, false
	}

	data := LoginResult{
		Name:  username,
		Token: token,
	}

	return data, true
}

//签发Token
func createToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signedString, err := token.SignedString(key)
	return signedString, err
}

// ParserToken Token解析
func ParserToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			// ValidationErrorMalformed是一个uint常量，表示token不可用
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
				// ValidationErrorExpired表示Token过期
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
				// ValidationErrorNotValidYet表示无效token
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

func rsaSing(str string) string {
	data := []byte(str)
	hashMd5 := md5.Sum(data)
	hashed := hashMd5[:]

	block, _ := pem.Decode(pemKey)
	if block == nil {
		panic(errors.New("private key error"))
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println("ParsePKCS1PrivateKey err", err)
		panic(err)
	}
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.MD5, hashed)
	return base64.StdEncoding.EncodeToString(signature)
}

func md5V(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func sing(str string) string {
	return str[60 : len(str)-40]
}
