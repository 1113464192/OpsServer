package tencentCloud

import (
	"fmt"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

func (s *TencentCloudService) CreateCloudHost(region string, instanceChargeType string, period string, renewFlag string, zone string, projectId string,
	instanceType string, imageId string, systemDiskType string, systemDiskSize string, dataDiskType string, dataDiskSize string, vpcId string,
	subnetId string, internetChargeType string, internetMaxBandwidthOut int, instanceName string, securityGroupId string, hostName string,
) error {
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
	request := cvm.NewRunInstancesRequest()

	request.InstanceChargeType = common.StringPtr("PREPAID")
	request.InstanceChargePrepaid = &cvm.InstanceChargePrepaid{
		Period:    common.Int64Ptr(1),
		RenewFlag: common.StringPtr("NOTIFY_AND_AUTO_RENEW"),
	}
	request.Placement = &cvm.Placement{
		Zone:      common.StringPtr("ap-guangzhou-01"),
		ProjectId: common.Int64Ptr(40),
	}
	request.InstanceType = common.StringPtr("asdsfhdsdsfds-asdas01")
	request.ImageId = common.StringPtr("img-487zeit5")
	request.SystemDisk = &cvm.SystemDisk{
		DiskType: common.StringPtr("CLOUD_BSSD"),
		DiskSize: common.Int64Ptr(40),
	}
	request.DataDisks = []*cvm.DataDisk{
		&cvm.DataDisk{
			DiskType:           common.StringPtr("CLOUD_BSSD"),
			DiskSize:           common.Int64Ptr(200),
			DeleteWithInstance: common.BoolPtr(true),
		},
	}
	request.VirtualPrivateCloud = &cvm.VirtualPrivateCloud{
		VpcId:        common.StringPtr("vpc-xxx"),
		SubnetId:     common.StringPtr("subnet-2ks"),
		AsVpcGateway: common.BoolPtr(true),
	}
	request.InternetAccessible = &cvm.InternetAccessible{
		InternetChargeType:      common.StringPtr("TRAFFIC_POSTPAID_BY_HOUR"),
		InternetMaxBandwidthOut: common.Int64Ptr(200),
		PublicIpAssigned:        common.BoolPtr(true),
	}
	request.InstanceName = common.StringPtr("xhxmj-01")
	request.SecurityGroupIds = common.StringPtrs([]string{"sg-hgdz3u17"})
	request.HostName = common.StringPtr("xhxmj-01")

	// 返回的resp是一个RunInstancesResponse的实例，与请求对象对应
	response, err := client.RunInstances(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return
	}
	if err != nil {
		panic(err)
	}
	// 输出json格式的字符串回包
	fmt.Printf("%s", response.ToJsonString())
}
