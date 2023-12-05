package jwt

import (
	"errors"
	"fmt"
	"fqhWeb/configs"
	"fqhWeb/internal/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// CustomSecret 用于加盐的字符串
var CustomSecret = []byte(configs.Conf.SecurityVars.TokenKey)

type CustomClaims struct {
	// 可根据需要自行添加字段
	User                 model.User
	jwt.RegisteredClaims // 内嵌标准的声明
}

// GenToken 生成JWT
func GenToken(user model.User) (string, error) {
	duration, err := time.ParseDuration(configs.Conf.SecurityVars.TokenExpireDuration)
	if err != nil {
		return "", fmt.Errorf("生成Token过期时间失败: %v", err)
	}
	// 创建一个我们自己的声明
	claims := CustomClaims{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			Issuer:    "fqh", // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(CustomSecret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*CustomClaims, error) {
	// 解析token
	// 如果是自定义Claim结构体则需要使用 ParseWithClaims 方法
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i any, err error) {
		// 直接使用标准的Claim则可以直接使用Parse方法
		return CustomSecret, nil
	})
	if err != nil {
		return nil, err
	}
	// 对token对象中的Claim进行类型断言
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
