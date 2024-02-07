package tencentCloud

import (
	"encoding/json"
	"fmt"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

func (s *TencentCloudService) GetCloudHostInVpcSubnetSum(region string, subnetId string) (int64, error) {
	type Response struct {
		RequestId             string `json:"RequestId"`
		ResourceStatisticsSet []struct {
			Ip                        int64 `json:"Ip"`
			ResourceStatisticsItemSet []struct {
				ResourceCount int    `json:"ResourceCount"`
				ResourceName  string `json:"ResourceName"`
				ResourceType  string `json:"ResourceType"`
			} `json:"ResourceStatisticsItemSet"`
			SubnetId string `json:"SubnetId"`
			VpcId    string `json:"VpcId"`
		} `json:"ResourceStatisticsSet"`
	}
	// 实例化一个认证对象，入参需要传入腾讯云账户 SecretId 和 SecretKey，此处还需注意密钥对的保密
	// 代码泄露可能会导致 SecretId 和 SecretKey 泄露，并威胁账号下所有资源的安全性。以下代码示例仅供参考，建议采用更安全的方式来使用密钥，请参见：https://cloud.tencent.com/document/product/1278/85305
	// 密钥可前往官网控制台 https://console.cloud.tencent.com/cam/capi 进行获取
	credential := common.NewCredential(
		s.ak,
		s.sk,
	)
	//completeRegion := RegionPrefix + region
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "vpc.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := vpc.NewClient(credential, region, cpf)

	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := vpc.NewDescribeSubnetResourceDashboardRequest()

	request.SubnetIds = common.StringPtrs([]string{subnetId})

	// 返回的resp是一个DescribeSubnetResourceDashboardResponse的实例，与请求对象对应
	response, err := client.DescribeSubnetResourceDashboard(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return 0, fmt.Errorf("an API error has returned: %v", err)
	}
	if err != nil {
		return 0, err
	}

	var res Response
	if err = json.Unmarshal([]byte(response.ToJsonString()), &res); err != nil {
		return 0, fmt.Errorf("json.Unmarshal error: %v", err)
	}
	return res.ResourceStatisticsSet[0].Ip, err
}
