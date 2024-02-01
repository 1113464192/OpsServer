package crontab

import (
	"encoding/json"
	errs "errors"
	"fmt"
	"fqhWeb/configs"
	"fqhWeb/internal/model"
	"fqhWeb/pkg/logger"
	"fqhWeb/pkg/util"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tag "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag/v20180813"
)

type Project struct {
	ProjectId   int    `json:"projectId"`
	ProjectName string `json:"projectName"`
}

type Response struct {
	ResponseData struct {
		Total    int       `json:"total"`
		Projects []Project `json:"projects"`
	} `json:"response"`
}

// 检查机器和云平台的项目是否有出漏
func CronCheckCloudProject() {
	var (
		err                 error
		projects            []model.Project
		tencentProjectNames []string
		aliyunProjectNames  []string
		cloudProjectNames   []string
	)
	if err = model.DB.Where("cloud != ?", "").Find(&projects).Error; err != nil {
		logger.Log().Error("Mysql", "获取非空云项目失败", err)
		// 接入微信小程序之类的请求, 向运维发送获取非空云项目失败的问题
		fmt.Println("微信小程序=====向运维发送处理获取非空云项目的问题")
		return
	}

	// 腾讯云
	for _, i := range projects {
		if i.Cloud == "腾讯云" {
			tencentProjectNames = append(tencentProjectNames, i.Name)
		}
		if i.Cloud == "阿里云" {
			aliyunProjectNames = append(aliyunProjectNames, i.Name)
		}
	}
	if err = cronCheckTencentCloudProject(&tencentProjectNames, &cloudProjectNames); err != nil {
		logger.Log().Error("CronCheckCloudProject", "腾讯云", err)
		// 接入微信小程序之类的请求, 向运维发送获取腾讯云项目失败的问题
		fmt.Println("微信小程序=====向运维发送处理获取腾讯云项目失败的问题")
		return
	}

	// 阿里云...
}

func cronCheckTencentCloudProject(projectNames *[]string, cloudProjectNames *[]string) error {
	defer func() {
		*projectNames = nil
		*cloudProjectNames = nil
	}()
	// 实例化一个认证对象，入参需要传入腾讯云账户 SecretId 和 SecretKey，此处还需注意密钥对的保密
	// 代码泄露可能会导致 SecretId 和 SecretKey 泄露，并威胁账号下所有资源的安全性。以下代码示例仅供参考，建议采用更安全的方式来使用密钥，请参见：https://cloud.tencent.com/document/product/1278/85305
	// 密钥可前往官网控制台 https://console.cloud.tencent.com/cam/capi 进行获取
	credential := common.NewCredential(
		configs.Conf.CloudSecretKey.TencentCloud.Ak,
		configs.Conf.CloudSecretKey.TencentCloud.Sk,
	)
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "tag.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := tag.NewClient(credential, "", cpf)

	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := tag.NewDescribeProjectsRequest()

	request.AllList = common.Uint64Ptr(0)
	request.Limit = common.Uint64Ptr(1000)
	request.Offset = common.Uint64Ptr(0)

	// 返回的resp是一个DescribeProjectsResponse的实例，与请求对象对应
	response, err := client.DescribeProjects(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return errs.New("An API error has returned " + err.Error())
	}
	if err != nil {
		return errs.New("实例获取报错 " + err.Error())
	}
	// 输出json格式的字符串回包
	res := Response{}
	if err = json.Unmarshal([]byte(response.ToJsonString()), &res); err != nil {
		return errs.New("json解析报错 " + err.Error())
	}

	for _, i := range res.ResponseData.Projects {
		*cloudProjectNames = append(*cloudProjectNames, i.ProjectName)
	}

	if diffProjects := util.StringSliceDifference(*projectNames, *cloudProjectNames); len(diffProjects) > 0 {
		return errs.New("数据库腾讯云项目和腾讯云平台项目组不一致 " + err.Error())
	}
	return nil
}
