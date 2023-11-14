package utils

import (
	"encoding/json"
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

	"github.com/Knetic/govaluate"
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

// 是否包含
func IsSliceContain(slice interface{}, value interface{}) bool {
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

// 环境变量提取整数
func GetEnvInt(key string, fallback int) int {
	ret := fallback
	value, exists := os.LookupEnv(key)
	if !exists {
		return ret
	}
	if t, err := strconv.Atoi(value); err != nil { //nolint:gosec
		return ret
	} else {
		ret = t
	}
	return ret
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

func IntSliceToStringSlice(intSlice []int) []string {
	stringSlice := make([]string, len(intSlice))
	for i, v := range intSlice {
		stringSlice[i] = strconv.Itoa(v)
	}
	return stringSlice
}

func Float64SliceToStringSlice(floatSlice []float64) []string {
	stringSlice := make([]string, len(floatSlice))
	for i, v := range floatSlice {
		stringSlice[i] = strconv.FormatFloat(v, 'f', -1, 64)
	}
	return stringSlice
}

// 有最大对应取最大，否则只取[0]
func SplitStringMap(originalMap map[string][]string) []map[string]string {
	maxLength := 0
	for _, values := range originalMap {
		if len(values) > maxLength {
			maxLength = len(values)
		}
	}

	// 创建一个切片用于存储拆分后的map
	splitMaps := make([]map[string]string, maxLength)

	// 遍历原始map
	for key, values := range originalMap {
		for i := 0; i < maxLength; i++ {
			// 如果值的长度大于i，则将值拆分到对应的map中；否则将空字符串放入map中
			if maxLength == len(values) {
				if splitMaps[i] == nil {
					splitMaps[i] = make(map[string]string)
				}
				splitMaps[i][key] = values[i]
			} else {
				if splitMaps[i] == nil {
					splitMaps[i] = make(map[string]string)
				}
				splitMaps[i][key] = values[0]
			}
		}
	}

	return splitMaps
}

// 用flag map类型, 做表达式中flag字符串的变量替换，生成结果为float64 slice类型
func GenerateExprResult(rule map[int]string, flag any) ([]float64, error) {
	var resultList []float64
	for _, rule := range rule {
		// 判断规则是否规范
		if !strings.Contains(rule, "flag") {
			return nil, errors.New(rule + " 不包含 flag 字符串")
		}
		// 创建规则表达式
		expr, err := govaluate.NewEvaluableExpression(rule)
		if err != nil {
			return nil, fmt.Errorf("创建表达式解析器报错: %v", err)
		}
		vars := map[string]any{
			"flag": flag,
		}
		// 获取出float64
		num, err := expr.Evaluate(vars)
		if err != nil {
			return nil, fmt.Errorf("表达式计算报错: %v", err)
		}
		resultList = append(resultList, num.(float64))
	}
	return resultList, nil
}

func ConvertToJson(params []string) (res string, err error) {
	var extraByte []byte
	var extra = make(map[int]string)
	if len(params) > 0 {
		for i, v := range params {
			extra[i] = v
		}
	}
	extraByte, err = json.Marshal(extra)
	if err != nil {
		return "", err
	}
	return string(extraByte), err
}

// 传x=y切片
func ConvertToJsonPair(params []string) (res string, err error) {
	data := make(map[string][]string)
	for _, param := range params {
		pair := strings.SplitN(param, "=", -1)
		if len(pair) != 2 {
			return "", fmt.Errorf("invalid key-value pair: %s", param)
		}
		key := pair[0]
		value := pair[1]
		data[key] = append(data[key], value)
	}
	var jsonData []byte
	jsonData, err = json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("%s: %v", "转换json报错", err)
	}
	return string(jsonData), err
}

func DeleteUintSlice(s []uint, elem uint) []uint {
	j := 0
	for _, v := range s {
		if v != elem {
			s[j] = v
			j++
		}
	}
	// 如果直接使用 return s，那么返回的切片 s 将包含原始切片中的所有元素，包括指定元素和非指定元素。这是因为切片是引用类型，返回的切片与传入的切片共享相同的底层数组。在这种情况下，虽然在循环中将指定元素跳过并将非指定元素移动到切片的前面，但没有对底层数组进行修改。因此，返回的切片 s 仍然包含了原始切片中的所有元素。
	return s[:j]
}

func DeleteAnySlice(s any, elem any) (any, error) {
	sliceValue := reflect.ValueOf(s)
	if sliceValue.Kind() != reflect.Slice {
		return s, errors.New("传入的首位参数, 类型不是slice")
	}
	j := 0
	for i := 0; i < sliceValue.Len(); i++ {
		v := sliceValue.Index(i).Interface()
		if v != elem {
			sliceValue.Index(j).Set(sliceValue.Index(i))
			j++
		}
	}
	return sliceValue.Slice(0, j).Interface(), nil
}
