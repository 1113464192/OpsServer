package tencentCloud

import (
	"encoding/json"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

type HostResponse struct {
	CloudHostResponse struct {
		InstanceSet []struct {
			RenewFlag            string `json:"RenewFlag"`
			Uuid                 string `json:"Uuid"`
			InstanceState        string `json:"InstanceState"`
			LatestOperationState string `json:"LatestOperationState"`
			LoginSettings        struct {
				Password       string   `json:"Password"`
				KeepImageLogin string   `json:"KeepImageLogin"`
				KeyIds         []string `json:"KeyIds"`
			} `json:"LoginSettings"`
			IPv6Addresses          []string `json:"IPv6Addresses"`
			DedicatedClusterId     string   `json:"DedicatedClusterId"`
			RestrictState          string   `json:"RestrictState"`
			ExpiredTime            string   `json:"ExpiredTime"`
			DisasterRecoverGroupId string   `json:"DisasterRecoverGroupId"`
			Memory                 int      `json:"Memory"`
			CreatedTime            string   `json:"CreatedTime"`
			CPU                    int      `json:"CPU"`
			RdmaIpAddresses        []string `json:"RdmaIpAddresses"`
			CamRoleName            string   `json:"CamRoleName"`
			PublicIpAddresses      []string `json:"PublicIpAddresses"`
			Tags                   []struct {
				Value string `json:"Value"`
				Key   string `json:"Key"`
			} `json:"Tags"`
			InstanceId         string `json:"InstanceId"`
			ImageId            string `json:"ImageId"`
			StopChargingMode   string `json:"StopChargingMode"`
			InstanceChargeType string `json:"InstanceChargeType"`
			InstanceType       string `json:"InstanceType"`
			SystemDisk         struct {
				DiskSize int    `json:"DiskSize"`
				CdcId    string `json:"CdcId"`
				DiskId   string `json:"DiskId"`
				DiskType string `json:"DiskType"`
			} `json:"SystemDisk"`
			Placement struct {
				HostId    string   `json:"HostId"`
				ProjectId int      `json:"ProjectId"`
				HostIds   []string `json:"HostIds"`
				Zone      string   `json:"Zone"`
			} `json:"Placement"`
			PrivateIpAddresses []string `json:"PrivateIpAddresses"`
			OsName             string   `json:"OsName"`
			SecurityGroupIds   []string `json:"SecurityGroupIds"`
			InstanceName       string   `json:"InstanceName"`
			DataDisks          []struct {
				DeleteWithInstance    bool        `json:"DeleteWithInstance"`
				Encrypt               bool        `json:"Encrypt"`
				CdcId                 string      `json:"CdcId"`
				DiskType              string      `json:"DiskType"`
				ThroughputPerformance int         `json:"ThroughputPerformance"`
				KmsKeyId              interface{} `json:"KmsKeyId"`
				DiskSize              int         `json:"DiskSize"`
				SnapshotId            interface{} `json:"SnapshotId"`
				DiskId                string      `json:"DiskId"`
			} `json:"DataDisks"`
			IsolatedSource      string `json:"IsolatedSource"`
			VirtualPrivateCloud struct {
				SubnetId           string   `json:"SubnetId"`
				AsVpcGateway       bool     `json:"AsVpcGateway"`
				Ipv6AddressCount   int      `json:"Ipv6AddressCount"`
				VpcId              string   `json:"VpcId"`
				PrivateIpAddresses []string `json:"PrivateIpAddresses"`
			} `json:"VirtualPrivateCloud"`
			LatestOperationRequestId string `json:"LatestOperationRequestId"`
			InternetAccessible       struct {
				PublicIpAssigned        bool   `json:"PublicIpAssigned"`
				InternetChargeType      string `json:"InternetChargeType"`
				InternetMaxBandwidthOut int    `json:"InternetMaxBandwidthOut"`
			} `json:"InternetAccessible"`
			HpcClusterId    string `json:"HpcClusterId"`
			LatestOperation string `json:"LatestOperation"`
		} `json:"InstanceSet"`
		TotalCount int    `json:"TotalCount"`
		RequestId  string `json:"RequestId"`
	} `json:"Response"`
}

func (s *TencentCloudService) GetCloudHostInfo(region string, publicIpv4 string, publicIpv6 string, offset int64, limit int64) (*HostResponse, error) {
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
	request := cvm.NewDescribeInstancesRequest()

	request.Filters = []*cvm.Filter{
		&cvm.Filter{
			Name:   common.StringPtr("instance-state"),
			Values: common.StringPtrs([]string{"RUNNING"}),
		},
	}
	// 都为空则全部返回
	if publicIpv4 != "" {
		request.Filters = append(request.Filters, &cvm.Filter{
			Name:   common.StringPtr("public-ip-address"),
			Values: common.StringPtrs([]string{publicIpv4}),
		})
	} else if publicIpv6 != "" {
		request.Filters = append(request.Filters, &cvm.Filter{
			Name:   common.StringPtr("ipv6-address"),
			Values: common.StringPtrs([]string{publicIpv6}),
		})
	}
	// 腾讯云默认都不为0的话，则offset=0，limit=20
	if offset != 0 && limit != 0 {
		request.Offset = common.Int64Ptr(offset)
		request.Limit = common.Int64Ptr(limit)
	} else if offset != 0 {
		request.Offset = common.Int64Ptr(offset)
	} else if limit != 0 {
		request.Limit = common.Int64Ptr(limit)
	}

	// 返回的resp是一个DescribeInstancesResponse的实例，与请求对象对应
	response, err := client.DescribeInstances(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return nil, fmt.Errorf("An API error has returned: %v", err)
	}
	if err != nil {
		return nil, fmt.Errorf("获取返回结果失败：%v", err)
	}
	var res HostResponse
	// 输出json格式的字符串回包
	if err = json.Unmarshal([]byte(response.ToJsonString()), &res); err != nil {
		return nil, fmt.Errorf("解析json失败：%v", err)
	}
	return &res, err
}
