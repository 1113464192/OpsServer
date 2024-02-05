package cloudScript

import (
	"fmt"
	"fqhWeb/pkg/cloudScript/tencentCloud"
)

func UpdateCloudProject(cloudType string, projectId uint64, projectName string, disable int64) (err error) {
	switch cloudType {
	case "腾讯云":
		if err = tencentCloud.TencentCloud().UpdateCloudProject(projectId, projectName, disable); err != nil {
			return err
		}

	case "火山云":
		return fmt.Errorf("%s 云商暂未加入平台，请通知运维加入", cloudType)
	default:
		return fmt.Errorf("%s 云商字符串有误，通知运维检查", cloudType)
	}
	return err
}

func CreateCloudProject(cloudType string, projectName string) (err error) {
	switch cloudType {
	case "腾讯云":
		if err = tencentCloud.TencentCloud().CreateCloudProject(projectName); err != nil {
			return err
		}
	case "火山云":
		return fmt.Errorf("%s 云商暂未加入平台，请通知运维加入", cloudType)
	default:
		return fmt.Errorf("%s 云商字符串有误，通知运维检查", cloudType)
	}
	return err
}

func GetCloudProjectId(cloudType string, projectName string) (uint64, error) {
	switch cloudType {
	case "腾讯云":
		return tencentCloud.TencentCloud().GetCloudProjectId(projectName)
	case "火山云":
		return 0, fmt.Errorf("%s 云商暂未加入平台，请通知运维加入", cloudType)
	default:
		return 0, fmt.Errorf("%s 云商字符串有误，通知运维检查", cloudType)
	}
}

func GetCloudHostInfo(cloudType string, region string, publicIpv4 string, publicIpv6 string, offset int64, limit int64) (any, error) {
	switch cloudType {
	case "腾讯云":
		return tencentCloud.TencentCloud().GetCloudHostInfo(region, publicIpv4, publicIpv6, offset, limit)
	case "火山云":
		return nil, fmt.Errorf("%s 云商暂未加入平台，请通知运维加入", cloudType)
	default:
		return nil, fmt.Errorf("%s 云商字符串有误，通知运维检查", cloudType)
	}
}

func CreateCloudHost(cloudType string, region string, hostName string, projectId uint64, publicIpv4 string, publicIpv6 string) (err error) {
	switch cloudType {
	case "腾讯云":
		if err = tencentCloud.TencentCloud().CreateCloudHost(region, hostName, projectId, publicIpv4, publicIpv6); err != nil {
			return err
		}
	case "火山云":
		return fmt.Errorf("%s 云商暂未加入平台，请通知运维加入", cloudType)
	default:
		return fmt.Errorf("%s 云商字符串有误，通知运维检查", cloudType)
	}
	return err
}

func GetAvailCloudZoneId(cloudType string, region string) ([]string, error) {
	switch cloudType {
	case "腾讯云":
		return tencentCloud.TencentCloud().GetAvailCloudZoneId(region)
	case "火山云":
		return nil, fmt.Errorf("%s 云商暂未加入平台，请通知运维加入", cloudType)
	default:
		return nil, fmt.Errorf("%s 云商字符串有误，通知运维检查", cloudType)
	}
}

func GetCloudInstanceTypeList(cloudType string, region string, instanceFamily string, cpuCores int, memorySize int, fpga int, gpuCores int) (any, error) {
	switch cloudType {
	case "腾讯云":
		return tencentCloud.TencentCloud().GetCloudInstanceTypeList(region, instanceFamily, cpuCores, memorySize, fpga, gpuCores)
	case "火山云":
		return nil, fmt.Errorf("%s 云商暂未加入平台，请通知运维加入", cloudType)
	default:
		return nil, fmt.Errorf("%s 云商字符串有误，通知运维检查", cloudType)
	}
}

func GetAvailCloudZoneI(cloudType string, region string) ([]string, error) {
	switch cloudType {
	case "腾讯云":
		return tencentCloud.TencentCloud().GetAvailCloudZoneId(region)
	case "火山云":
		return nil, fmt.Errorf("%s 云商暂未加入平台，请通知运维加入", cloudType)
	default:
		return nil, fmt.Errorf("%s 云商字符串有误，通知运维检查", cloudType)
	}
}

func CreateCloudVpc(cloudType string, region string, vpcName string, cidrBlock string) error {
	switch cloudType {
	case "腾讯云":
		return tencentCloud.TencentCloud().CreateCloudVpc(region, vpcName, cidrBlock)
	case "火山云":
		return fmt.Errorf("%s 云商暂未加入平台，请通知运维加入", cloudType)
	default:
		return fmt.Errorf("%s 云商字符串有误，通知运维检查", cloudType)
	}
}

func GetCloudVpcId(cloudType string, region string, vpcName string) (string, error) {
	switch cloudType {
	case "腾讯云":
		return tencentCloud.TencentCloud().GetCloudVpcId(region, vpcName)
	case "火山云":
		return "", fmt.Errorf("%s 云商暂未加入平台，请通知运维加入", cloudType)
	default:
		return "", fmt.Errorf("%s 云商字符串有误，通知运维检查", cloudType)
	}
}

func CreateCloudSecurityGroup(cloudType string, region string, projectId string, groupName string, groupDescription string) error {
	switch cloudType {
	case "腾讯云":
		return tencentCloud.TencentCloud().CreateCloudSecurityGroup(region, projectId, groupName, groupDescription)
	case "火山云":
		return fmt.Errorf("%s 云商暂未加入平台，请通知运维加入", cloudType)
	default:
		return fmt.Errorf("%s 云商字符串有误，通知运维检查", cloudType)
	}
}

func CreateCloudVpcSubnet(cloudType string, region string, vpcId string, subnetName string, subnetCidrBlock string, zone string) error {
	switch cloudType {
	case "腾讯云":
		return tencentCloud.TencentCloud().CreateCloudVpcSubnet(region, vpcId, subnetName, subnetCidrBlock, zone)
	case "火山云":
		return fmt.Errorf("%s 云商暂未加入平台，请通知运维加入", cloudType)
	default:
		return fmt.Errorf("%s 云商字符串有误，通知运维检查", cloudType)
	}
}

func GetCloudSecurityGroupId(cloudType string, region string, securityGroupName string) (string, error) {
	switch cloudType {
	case "腾讯云":
		return tencentCloud.TencentCloud().GetCloudSecurityGroupId(region, securityGroupName)
	case "火山云":
		return "", fmt.Errorf("%s 云商暂未加入平台，请通知运维加入", cloudType)
	default:
		return "", fmt.Errorf("%s 云商字符串有误，通知运维检查", cloudType)
	}
}

func GetCloudVpcSubnetId(cloudType string, region string, subnetName string) (string, error) {
	switch cloudType {
	case "腾讯云":
		return tencentCloud.TencentCloud().GetCloudVpcSubnetId(region, subnetName)
	case "火山云":
		return "", fmt.Errorf("%s 云商暂未加入平台，请通知运维加入", cloudType)
	default:
		return "", fmt.Errorf("%s 云商字符串有误，通知运维检查", cloudType)
	}
}

func GetCloudHostInVpcSubnetSum(cloudType string, region string, subnetId string) (int64, error) {
	switch cloudType {
	case "腾讯云":
		return tencentCloud.TencentCloud().GetCloudHostInVpcSubnetSum(region, subnetId)
	case "火山云":
		return 0, fmt.Errorf("%s 云商暂未加入平台，请通知运维加入", cloudType)
	default:
		return 0, fmt.Errorf("%s 云商字符串有误，通知运维检查", cloudType)
	}
}

func ReturnCloudHost(cloudType string, region string, instanceId string) error {
	switch cloudType {
	case "腾讯云":
		return tencentCloud.TencentCloud().ReturnCloudHost(region, instanceId)
	case "火山云":
		return fmt.Errorf("%s 云商暂未加入平台，请通知运维加入", cloudType)
	default:
		return fmt.Errorf("%s 云商字符串有误，通知运维检查", cloudType)
	}
}
