package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// 定义过期时间

const TokenExpireDuration = time.Hour * 24

// CustomSecret 用于加盐的字符串

var MySecret = []byte("redditdemosecret")

// 定义自己的数据结构
// MyClaims 自定义声明类型 并内嵌jwt.RegisteredClaims
// jwt包自带的jwt.RegisteredClaims只包含了官方字段
// 假设我们这里需要额外记录一个username字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中

type MyClaims struct {
	// 可根据需要自行添加字段
	UserID             int64  `json:"user_id"`
	Username           string `json:"username"`
	jwt.StandardClaims        // 内嵌标准的声明
}

//// GenToken 生成JWT
//func GenToken(userID int64, username string) (string, error) {
//	// 创建一个我们自己的声明的数据结构
//	claims := MyClaims{
//		userID,
//		username, // 自定义字段
//		jwt.RegisteredClaims{
//			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)),
//			Issuer:    "redditdemo", // 签发人
//		},
//	}
//	// 使用指定的签名方法创建签名对象
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//	// 使用指定的secret签名并获得完整的编码后的字符串token
//	return token.SignedString(MySecret)
//}

// GenToken 生成JWT
func GenToken(userID int64, username string) (string, error) {
	// 创建一个我们自己的声明的数据
	c := MyClaims{
		userID,
		"username", // 自定义字段
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), // 过期时间
			Issuer:    "redditdemo",                               // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(MySecret)
}

//// ParseToken 解析JWT
//func ParseToken(tokenString string) (*MyClaims, error) {
//	// 解析token
//	// 如果是自定义Claim结构体则需要使用 ParseWithClaims 方法
//	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
//		// 直接使用标准的Claim则可以直接使用Parse方法
//		//token, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
//		return MySecret, nil
//	})
//	if err != nil {
//		return nil, err
//	}
//	// 对token对象中的Claim进行类型断言
//	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid { // 校验token
//		return claims, nil
//	}
//	return nil, errors.New("invalid token")
//}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token
	var mc = new(MyClaims)
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (i interface{}, err error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if token.Valid { // 校验token
		return mc, nil
	}
	return nil, errors.New("invalid token")
}
