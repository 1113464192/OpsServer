package tencentCloud

import (
	"fmt"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

func (s *TencentCloudService) CreateCloudInstance(region string, instanceChargeType string, period int64, renewFlag string, zone string, projectId int64,
	instanceType string, imageId string, systemDiskType string, systemDiskSize int64, dataDiskType string, dataDiskSize int64, vpcId string,
	subnetId string, internetChargeType string, internetMaxBandwidthOut int64, instanceName string, securityGroupId string, hostName string,
) error {
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
	cpf.HttpProfile.Endpoint = "cvm.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := cvm.NewClient(credential, region, cpf)

	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := cvm.NewRunInstancesRequest()

	request.InstanceChargeType = common.StringPtr(instanceChargeType)
	request.InstanceChargePrepaid = &cvm.InstanceChargePrepaid{
		Period:    common.Int64Ptr(period),
		RenewFlag: common.StringPtr(renewFlag),
	}
	request.Placement = &cvm.Placement{
		Zone:      common.StringPtr(zone),
		ProjectId: common.Int64Ptr(projectId),
	}
	request.InstanceType = common.StringPtr(instanceType)
	request.ImageId = common.StringPtr(imageId)
	request.SystemDisk = &cvm.SystemDisk{
		DiskType: common.StringPtr(systemDiskType),
		DiskSize: common.Int64Ptr(systemDiskSize),
	}
	request.DataDisks = []*cvm.DataDisk{
		&cvm.DataDisk{
			DiskType:           common.StringPtr(dataDiskType),
			DiskSize:           common.Int64Ptr(dataDiskSize),
			DeleteWithInstance: common.BoolPtr(true),
		},
	}
	request.VirtualPrivateCloud = &cvm.VirtualPrivateCloud{
		VpcId:        common.StringPtr(vpcId),
		SubnetId:     common.StringPtr(subnetId),
		AsVpcGateway: common.BoolPtr(true),
	}
	request.InternetAccessible = &cvm.InternetAccessible{
		InternetChargeType:      common.StringPtr(internetChargeType),
		InternetMaxBandwidthOut: common.Int64Ptr(internetMaxBandwidthOut),
		PublicIpAssigned:        common.BoolPtr(true),
	}
	request.InstanceName = common.StringPtr(instanceName)
	request.SecurityGroupIds = common.StringPtrs([]string{securityGroupId})
	request.HostName = common.StringPtr(hostName)

	// 返回的resp是一个RunInstancesResponse的实例，与请求对象对应
	_, err := client.RunInstances(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return fmt.Errorf("an API error has returned: %v", err)
	}
	if err != nil {
		return err
	}
	return err
}
