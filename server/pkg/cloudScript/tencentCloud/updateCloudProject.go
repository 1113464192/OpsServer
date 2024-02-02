package tencentCloud

import (
	"fmt"
	"fqhWeb/pkg/logger"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tag "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag/v20180813"
)

func (s *TencentCloudService) UpdateCloudProject(projectId uint64, projectName string, disable int64) (err error) {
	// 实例化一个认证对象，入参需要传入腾讯云账户 SecretId 和 SecretKey，此处还需注意密钥对的保密
	// 代码泄露可能会导致 SecretId 和 SecretKey 泄露，并威胁账号下所有资源的安全性。以下代码示例仅供参考，建议采用更安全的方式来使用密钥，请参见：https://cloud.tencent.com/document/product/1278/85305
	// 密钥可前往官网控制台 https://console.cloud.tencent.com/cam/capi 进行获取
	credential := common.NewCredential(
		s.ak,
		s.sk,
	)
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "tag.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := tag.NewClient(credential, "", cpf)

	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := tag.NewUpdateProjectRequest()

	request.ProjectId = common.Uint64Ptr(projectId)
	if projectName != "" {
		request.ProjectName = common.StringPtr(projectName)
	}
	if disable == 1 {
		request.Disable = common.Int64Ptr(TencentProjectEnable)
	} else if disable == 2 {
		request.Disable = common.Int64Ptr(TencentProjectDisable)
	}

	// 返回的resp是一个UpdateProjectResponse的实例，与请求对象对应
	response, err := client.UpdateProject(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return fmt.Errorf("an API error has returned: %v", err)
	}
	if err != nil {
		return err
	}
	logger.Log().Info("Project", "DisableProjectGroup", response)
	return nil
}
