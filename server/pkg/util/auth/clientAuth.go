package auth

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"fqhWeb/configs"
	"strings"
)

func CheckClientReqAuth(clientSign string, clientIp string) (err error) {
	if clientSign == "" {
		return errors.New("未从header获取签名到ClientAuthSign的值")
	} else if clientIp == "" {
		return errors.New("未从header获取签名到clientIp的值")
	}
	sign, err := Md5EncryptSign(clientIp, configs.Conf.SecurityVars.ClientReqMd5Key)
	if err != nil {
		return fmt.Errorf("sign生成报错: %v", err)
	}
	if sign != clientSign {
		return errors.New(`认证码错误，请确认认证码生成方式`)
	}
	return err
}

func Md5EncryptSign(clientIp string, md5Key string) (sign string, err error) {
	builder := strings.Builder{}
	builder.WriteString(md5Key)
	builder.WriteString(clientIp)
	md5Hash := md5.Sum([]byte(builder.String()))
	sign = hex.EncodeToString(md5Hash[:])
	return sign, err
}
