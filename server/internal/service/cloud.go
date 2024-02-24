package service

import (
	"errors"
	"fmt"
	"fqhWeb/internal/model"
	"fqhWeb/pkg/api"
	"fqhWeb/pkg/cloudScript/tencentCloud"
	"fqhWeb/pkg/logger"
	"fqhWeb/pkg/util2"
)

type CloudService struct{}

var (
	insCloud = &CloudService{}
)

func Cloud() *CloudService {
	return insCloud
}

// 更改购买实例配置
func (s *CloudService) UpdateCloudInstanceConfig(param *api.UpdateCloudInstanceConfigReq) (res *model.CloudInstanceConfig, err error) {
	var config model.CloudInstanceConfig
	var count int64
	// 判断项目名是否已被使用
	if model.DB.Model(&config).Where("project_id = ? AND id != ?", param.ProjectId, param.Id).Count(&count); count > 0 {
		return &config, errors.New("项目名已被使用")
	}
	// ID查询
	if param.Id != 0 {
		if !util2.CheckIdExists(&config, param.Id) {
			return &config, errors.New("项目不存在")
		}

		if err := model.DB.Where("id = ?", param.Id).First(&config).Error; err != nil {
			return &config, errors.New("项目数据库查询失败")
		}

		config.Region = param.Region
		config.InstanceChargeType = param.InstanceChargeType
		config.Period = param.Period
		config.RenewFlag = param.RenewFlag
		config.ProjectId = param.ProjectId
		config.InstanceFamily = param.InstanceFamily
		config.CpuCores = param.CpuCores
		config.MemorySize = param.MemorySize
		config.Fpga = param.Fpga
		config.GpuCores = param.GpuCores
		config.ImageId = param.ImageId
		config.SystemDiskType = param.SystemDiskType
		config.SystemDiskSize = param.SystemDiskSize
		config.DataDiskType = param.DataDiskType
		config.DataDiskSize = param.DataDiskSize
		config.VpcId = param.VpcId
		config.SubnetId = param.SubnetId
		config.InternetChargeType = param.InternetChargeType
		config.InternetMaxBandwidthOut = param.InternetMaxBandwidthOut
		config.InstanceNamePrefix = param.InstanceNamePrefix
		config.SecurityGroupId = param.SecurityGroupId
		config.HostNamePrefix = param.HostNamePrefix

		if err = model.DB.Save(&config).Error; err != nil {
			return &config, fmt.Errorf("数据保存失败: %v", err)
		}
		return &config, err
	} else {
		config = model.CloudInstanceConfig{
			Region:                  param.Region,
			InstanceChargeType:      param.InstanceChargeType,
			Period:                  param.Period,
			RenewFlag:               param.RenewFlag,
			ProjectId:               param.ProjectId,
			InstanceFamily:          param.InstanceFamily,
			CpuCores:                param.CpuCores,
			MemorySize:              param.MemorySize,
			Fpga:                    param.Fpga,
			GpuCores:                param.GpuCores,
			ImageId:                 param.ImageId,
			SystemDiskType:          param.SystemDiskType,
			SystemDiskSize:          param.SystemDiskSize,
			DataDiskType:            param.DataDiskType,
			DataDiskSize:            param.DataDiskSize,
			VpcId:                   param.VpcId,
			SubnetId:                param.SubnetId,
			InternetChargeType:      param.InternetChargeType,
			InternetMaxBandwidthOut: param.InternetMaxBandwidthOut,
			InstanceNamePrefix:      param.InstanceNamePrefix,
			SecurityGroupId:         param.SecurityGroupId,
			HostNamePrefix:          param.HostNamePrefix,
		}

		if err = model.DB.Create(&config).Error; err != nil {
			logger.Log().Error("cloud", "创建实例对应配置失败", err)
			return &config, errors.New("创建实例对应配置失败")
		}
		return &config, err
	}
}

func (s *CloudService) GetCloudInstanceConfig(id uint) (*model.CloudInstanceConfig, error) {
	var (
		err     error
		project model.Project
	)
	if !util2.CheckIdExists(&project, id) {
		return nil, errors.New("项目不存在")
	}
	if err = model.DB.Preload("CloudInstanceConfig").Where("id = ?", id).First(&project).Error; err != nil {
		return nil, errors.New("项目数据库查询失败")
	}

	return &project.CloudInstanceConfig, err
}

func (s *CloudService) UpdateCloudProject(cloudType string, projectId uint64, projectName string, disable int64) (err error) {
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

func (s *CloudService) CreateCloudProject(cloudType string, projectName string) (err error) {
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

func (s *CloudService) GetCloudProjectId(cloudType string, projectName string) (uint64, error) {
	switch cloudType {
	case "腾讯云":
		return tencentCloud.TencentCloud().GetCloudProjectId(projectName)
	case "火山云":
		return 0, fmt.Errorf("%s 云商暂未加入平台，请通知运维加入", cloudType)
	default:
		return 0, fmt.Errorf("%s 云商字符串有误，通知运维检查", cloudType)
	}
}

func (s *CloudService) GetCloudInsInfo(cloudType string, region string, publicIpv4 string, publicIpv6 string, insName string, cloudPid string, offset int64, limit int64) (any, error) {
	switch cloudType {
	case "腾讯云":
		return tencentCloud.TencentCloud().GetCloudInsInfo(region, publicIpv4, publicIpv6, insName, cloudPid, offset, limit)
	case "火山云":
		return nil, fmt.Errorf("%s 云商暂未加入平台，请通知运维加入", cloudType)
	default:
		return nil, fmt.Errorf("%s 云商字符串有误，通知运维检查", cloudType)
	}
}

func (s *CloudService) CreateCloudInstance(cloudType string, region string, instanceChargeType string, period int64, renewFlag string, zone string, projectId int64,
	instanceType string, imageId string, systemDiskType string, systemDiskSize int64, dataDiskType string, dataDiskSize int64, vpcId string,
	subnetId string, internetChargeType string, internetMaxBandwidthOut int64, instanceName string, securityGroupId string, hostName string) error {
	switch cloudType {
	case "腾讯云":
		return tencentCloud.TencentCloud().CreateCloudInstance(region, instanceChargeType, period, renewFlag, zone, projectId, instanceType, imageId,
			systemDiskType, systemDiskSize, dataDiskType, dataDiskSize, vpcId, subnetId, internetChargeType, internetMaxBandwidthOut, instanceName,
			securityGroupId, hostName)
	case "火山云":
		return fmt.Errorf("%s 云商暂未加入平台，请通知运维加入", cloudType)
	default:
		return fmt.Errorf("%s 云商字符串有误，通知运维检查", cloudType)
	}
}

func (s *CloudService) GetCloudInstanceTypeList(cloudType string, region string, instanceFamily string, cpuCores int, memorySize int, fpga int, gpuCores int) (any, error) {
	switch cloudType {
	case "腾讯云":
		return tencentCloud.TencentCloud().GetCloudInstanceTypeList(region, instanceFamily, cpuCores, memorySize, fpga, gpuCores)
	case "火山云":
		return nil, fmt.Errorf("%s 云商暂未加入平台，请通知运维加入", cloudType)
	default:
		return nil, fmt.Errorf("%s 云商字符串有误，通知运维检查", cloudType)
	}
}

func (s *CloudService) GetAvailCloudZoneId(cloudType string, region string) ([]string, error) {
	switch cloudType {
	case "腾讯云":
		return tencentCloud.TencentCloud().GetAvailCloudZoneId(region)
	case "火山云":
		return nil, fmt.Errorf("%s 云商暂未加入平台，请通知运维加入", cloudType)
	default:
		return nil, fmt.Errorf("%s 云商字符串有误，通知运维检查", cloudType)
	}
}

func (s *CloudService) CreateCloudVpc(cloudType string, region string, vpcName string, cidrBlock string) error {
	switch cloudType {
	case "腾讯云":
		return tencentCloud.TencentCloud().CreateCloudVpc(region, vpcName, cidrBlock)
	case "火山云":
		return fmt.Errorf("%s 云商暂未加入平台，请通知运维加入", cloudType)
	default:
		return fmt.Errorf("%s 云商字符串有误，通知运维检查", cloudType)
	}
}

func (s *CloudService) GetCloudVpcId(cloudType string, region string, vpcName string) (string, error) {
	switch cloudType {
	case "腾讯云":
		return tencentCloud.TencentCloud().GetCloudVpcId(region, vpcName)
	case "火山云":
		return "", fmt.Errorf("%s 云商暂未加入平台，请通知运维加入", cloudType)
	default:
		return "", fmt.Errorf("%s 云商字符串有误，通知运维检查", cloudType)
	}
}

func (s *CloudService) CreateCloudSecurityGroup(cloudType string, region string, projectId string, groupName string, groupDescription string) error {
	switch cloudType {
	case "腾讯云":
		return tencentCloud.TencentCloud().CreateCloudSecurityGroup(region, projectId, groupName, groupDescription)
	case "火山云":
		return fmt.Errorf("%s 云商暂未加入平台，请通知运维加入", cloudType)
	default:
		return fmt.Errorf("%s 云商字符串有误，通知运维检查", cloudType)
	}
}

func (s *CloudService) CreateCloudVpcSubnet(cloudType string, region string, vpcId string, subnetName string, subnetCidrBlock string) error {
	switch cloudType {
	case "腾讯云":
		zone, err := tencentCloud.TencentCloud().GetAvailCloudZoneId(region)
		if err != nil {
			return fmt.Errorf("获取可用Zone失败: %v", err)
		}
		return tencentCloud.TencentCloud().CreateCloudVpcSubnet(region, vpcId, subnetName, subnetCidrBlock, zone[0])
	case "火山云":
		return fmt.Errorf("%s 云商暂未加入平台，请通知运维加入", cloudType)
	default:
		return fmt.Errorf("%s 云商字符串有误，通知运维检查", cloudType)
	}
}

func (s *CloudService) GetCloudSecurityGroupId(cloudType string, region string, securityGroupName string) (string, error) {
	switch cloudType {
	case "腾讯云":
		return tencentCloud.TencentCloud().GetCloudSecurityGroupId(region, securityGroupName)
	case "火山云":
		return "", fmt.Errorf("%s 云商暂未加入平台，请通知运维加入", cloudType)
	default:
		return "", fmt.Errorf("%s 云商字符串有误，通知运维检查", cloudType)
	}
}

func (s *CloudService) GetCloudVpcSubnetId(cloudType string, region string, subnetName string) (string, error) {
	switch cloudType {
	case "腾讯云":
		return tencentCloud.TencentCloud().GetCloudVpcSubnetId(region, subnetName)
	case "火山云":
		return "", fmt.Errorf("%s 云商暂未加入平台，请通知运维加入", cloudType)
	default:
		return "", fmt.Errorf("%s 云商字符串有误，通知运维检查", cloudType)
	}
}

func (s *CloudService) GetCloudHostInVpcSubnetSum(cloudType string, region string, subnetId string) (int64, error) {
	switch cloudType {
	case "腾讯云":
		return tencentCloud.TencentCloud().GetCloudHostInVpcSubnetSum(region, subnetId)
	case "火山云":
		return 0, fmt.Errorf("%s 云商暂未加入平台，请通知运维加入", cloudType)
	default:
		return 0, fmt.Errorf("%s 云商字符串有误，通知运维检查", cloudType)
	}
}

func (s *CloudService) ReturnCloudInstance(cloudType string, region string, instanceId string) error {
	switch cloudType {
	case "腾讯云":
		return tencentCloud.TencentCloud().ReturnCloudInstance(region, instanceId)
	case "火山云":
		return fmt.Errorf("%s 云商暂未加入平台，请通知运维加入", cloudType)
	default:
		return fmt.Errorf("%s 云商字符串有误，通知运维检查", cloudType)
	}
}

func (s *CloudService) GetCloudInsRenewPrice(cloudType string, hid uint, pid uint) (string, error) {
	var (
		err     error
		host    model.Host
		project model.Project
	)
	if err = model.DB.First(&host, hid).Error; err != nil {
		return "", fmt.Errorf("查询服务器失败: %v", err)
	}
	if err = model.DB.Preload("CloudInstanceConfig").First(&project, pid).Error; err != nil {
		return "", fmt.Errorf("查询项目失败: %v", err)
	}
	switch cloudType {
	case "腾讯云":
		var res *tencentCloud.HostResponse
		if res, err = tencentCloud.TencentCloud().GetCloudInsInfo(project.CloudInstanceConfig.Region, host.Ipv4.String, "", "", "", 1, 1); err != nil {
			return "", err
		}
		cost, err := tencentCloud.TencentCloud().GetCloudInsRenewPrice(res.CloudHostResponse.InstanceSet[0].InstanceId, &project.CloudInstanceConfig)
		if err != nil {
			return "", err
		}
		return cost, err
	case "火山云":
		return "", fmt.Errorf("%s 云商暂未加入平台，请通知运维加入", cloudType)
	default:
		return "", fmt.Errorf("%s 云商字符串有误，通知运维检查", cloudType)
	}
}
