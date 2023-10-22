package utils

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func GetRootPath() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	// return dir
	return strings.Replace(dir, "\\", "/", -1)
}

func IsDir(path string) bool {
	dirStat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return dirStat.IsDir()
}

func tmpLogWrite(msg string) bool {
	filePath := GetRootPath() + "/logs/tmp.log"

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("临时文件打开失败")
		return false
	}

	defer file.Close()
	// 创建一个写入器用作追加
	writer := io.MultiWriter(file)
	if _, err := io.WriteString(writer, msg+"\n"); err != nil {
		return false
	}
	return true
}

func CommonLog(service string, msg string) bool {
	var dirPath, file string
	if service == "" {
		dirPath = GetRootPath() + "/logs" + "/common"
		file = dirPath + "/" + "common" + time.Now().Format("01") + ".log"
	} else {
		dirPath = GetRootPath() + "/logs/" + service
		file = dirPath + "/" + service + time.Now().Format("01") + ".log"
	}

	if !IsDir(dirPath) {
		if err := os.Mkdir(dirPath, 0775); err != nil {
			tmpBool := tmpLogWrite(time.Now().Local().Format("2006-01-02 15:04:05") + "mkdir failed！ " + err.Error())
			if !tmpBool {
				panic(fmt.Errorf("临时日志文件写入失败"))
			}
		}
	}
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		tmpBool := tmpLogWrite(time.Now().Local().Format("2006-01-02 15:04:05") + "打开日志文件失败！ " + err.Error())
		if !tmpBool {
			panic(fmt.Errorf("临时日志文件写入失败"))
		}
	}
	log.SetOutput(logFile)
	log.SetPrefix("[" + service + "]" + "[" + time.Now().Local().Format("2006-01-02 15:04:05") + "] ")
	log.Println(msg)
	return true
}

func IsContain(slice interface{}, value interface{}) bool {
	sliceValue := reflect.ValueOf(slice)
	if sliceValue.Kind() != reflect.Slice {
		return false
	}

	for i := 0; i < sliceValue.Len(); i++ {
		item := sliceValue.Index(i).Interface()
		if reflect.DeepEqual(item, value) {
			return true
		}
	}

	return false
}

// 匹配手机号
func CheckMobile(phone string) bool {
	reg := `^1(3\d{2}|4[14-9]\d|5([0-35689]\d|7[1-79])|66\d|7[2-35-8]\d|8\d{2}|9[13589]\d)\d{7}$`
	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(phone)
}

// 匹配电子邮箱
func CheckEmail(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

func Pointer[T any](in T) (out *T) {
	return &in
}

// RandStringRunes 返回随机字符串
func RandStringRunes(n int) string {
	var letterRunes = []byte("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// 返回IP类型是IPV4orIPV6	1:Ipv4 2:Ipv6
func ReturnIpType(ipStr string) (uint8, error) {
	ip := net.ParseIP(ipStr)
	if ip != nil {
		if ip.To4() != nil {
			return 1, nil
		} else if ip.To16() != nil {
			return 2, nil
		} else {
			return 255, errors.New("IP类型错误")
		}
	}
	return 255, errors.New("未收到IP")
}

// string转换uint
func StringToUint(idStr *string) (id uint, err error) {
	var oldId uint64
	oldId, err = strconv.ParseUint(*idStr, 10, 0)
	if err != nil {
		return 0, errors.New("uint类型转换失败")
	}
	id = uint(oldId)
	return id, err
}
