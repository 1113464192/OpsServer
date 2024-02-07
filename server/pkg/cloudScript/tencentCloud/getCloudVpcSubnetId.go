package tencentCloud

import (
	"encoding/json"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

func (s *TencentCloudService) GetCloudVpcSubnetId(region string, subnetName string) (string, error) {
	type Response struct {
		RequestId string `json:"RequestId"`
		SubnetSet []struct {
			SubnetId string `json:"SubnetId"`
		} `json:"SubnetSet"`
	}
	// 实例化一个认证对象，入参需要传入腾讯云账户 SecretId 和 SecretKey，此处还需注意密钥对的保密
	// 代码泄露可能会导致 SecretId 和 SecretKey 泄露，并威胁账号下所有资源的安全性。以下代码示例仅供参考，建议采用更安全的方式来使用密钥，请参见：https://cloud.tencent.com/document/product/1278/85305
	// 密钥可前往官网控制台 https://console.cloud.tencent.com/cam/capi 进行获取
	credential := common.NewCredential(
		"SecretId",
		"SecretKey",
	)
	//completeRegion := RegionPrefix + region
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "vpc.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := vpc.NewClient(credential, region, cpf)

	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := vpc.NewDescribeSubnetsRequest()

	request.Filters = []*vpc.Filter{
		&vpc.Filter{
			Name:   common.StringPtr("subnet-name"),
			Values: common.StringPtrs([]string{subnetName}),
		},
	}
	request.Offset = common.StringPtr("0")
	request.Limit = common.StringPtr("1")

	// 返回的resp是一个DescribeSubnetsResponse的实例，与请求对象对应
	response, err := client.DescribeSubnets(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return "", fmt.Errorf("an API error has returned: %v", err)
	}
	if err != nil {
		return "", err
	}
	var res Response
	err = json.Unmarshal([]byte(response.ToJsonString()), &res)
	if err != nil {
		return "", fmt.Errorf("json.Unmarshal error: %v", err)
	}

	return res.SubnetSet[0].SubnetId, err
}
