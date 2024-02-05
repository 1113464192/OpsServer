package tencentCloud

import (
	"encoding/json"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

type InstanceConfigRes struct {
	InstanceType string `json:"InstanceType"`
	Zone         string `json:"Zone"`
}

func (s *TencentCloudService) GetCloudInstanceTypeList(region string, instanceFamily string, cpuCores int, memorySize int, fpga int, gpuCores int) ([]InstanceConfigRes, error) {
	type InstanceTypeConfig struct {
		CPU          int    `json:"CPU"`
		FPGA         int    `json:"FPGA"`
		GPU          int    `json:"GPU"`
		GpuCount     int    `json:"GpuCount"`
		InstanceType string `json:"InstanceType"`
		Memory       int    `json:"Memory"`
		Zone         string `json:"Zone"`
	}

	type Response struct {
		InstanceResponse struct {
			InstanceTypeConfigSet []InstanceTypeConfig `json:"InstanceTypeConfigSet"`
		} `json:"Response"`
	}

	// 实例化一个认证对象，入参需要传入腾讯云账户 SecretId 和 SecretKey，此处还需注意密钥对的保密
	// 代码泄露可能会导致 SecretId 和 SecretKey 泄露，并威胁账号下所有资源的安全性。以下代码示例仅供参考，建议采用更安全的方式来使用密钥，请参见：https://cloud.tencent.com/document/product/1278/85305
	// 密钥可前往官网控制台 https://console.cloud.tencent.com/cam/capi 进行获取
	credential := common.NewCredential(
		s.ak,
		s.sk,
	)
	completeRegion := RegionPrefix + region
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "cvm.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := cvm.NewClient(credential, completeRegion, cpf)

	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := cvm.NewDescribeInstanceTypeConfigsRequest()

	request.Filters = []*cvm.Filter{
		&cvm.Filter{
			Name:   common.StringPtr("instance-family"),
			Values: common.StringPtrs([]string{instanceFamily}),
		},
	}

	// 返回的resp是一个DescribeInstanceTypeConfigsResponse的实例，与请求对象对应
	response, err := client.DescribeInstanceTypeConfigs(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return nil, fmt.Errorf("an API error has returned: %v", err)
	}
	if err != nil {
		return nil, err
	}

	var res Response
	err = json.Unmarshal([]byte(response.ToJsonString()), &res)
	if err != nil {
		return nil, fmt.Errorf("json解析失败: %v", err)
	}

	var insList []InstanceConfigRes
	for _, config := range res.InstanceResponse.InstanceTypeConfigSet {
		if config.CPU == cpuCores && config.Memory == memorySize && config.FPGA == fpga && config.GPU == gpuCores {
			insList = append(insList, InstanceConfigRes{
				InstanceType: config.InstanceType,
				Zone:         config.Zone,
			})
		}
	}
	return insList, err
}
