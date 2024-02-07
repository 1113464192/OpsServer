package api

type CreateCloudVpcReq struct {
	CloudType string `json:"cloud_type" form:"cloud_type"  binding:"required"` // 如: 腾讯云、阿里云
	Region    string `json:"region" form:"region"  binding:"required"`
	VpcName   string `json:"vpc_name" form:"vpc_name"  binding:"required"` // VpcName与项目同名
	CidrBlock string `json:"cidr_block" form:"cidr_block"  binding:"required"`
}

type CreateCloudVpcSubnetReq struct {
	CloudType       string `json:"cloud_type" form:"cloud_type"  binding:"required"` // 如: 腾讯云、阿里云
	Region          string `json:"region" form:"region"  binding:"required"`
	VpcName         string `json:"vpc_name" form:"vpc_name"  binding:"required"`
	SubnetName      string `json:"subnet_name" form:"subnet_name"  binding:"required"`
	SubnetCidrBlock string `json:"subnet_cidr_block" form:"subnet_cidr_block"  binding:"required"`
}

type CreateCloudSecurityGroupReq struct {
	CloudType        string `json:"cloud_type" form:"cloud_type"  binding:"required"` // 如: 腾讯云、阿里云
	Region           string `json:"region" form:"region"  binding:"required"`
	ProjectName      string `json:"project_name" form:"project_name"  binding:"required"`
	GroupName        string `json:"group_name" form:"group_name"  binding:"required"`
	GroupDescription string `json:"group_description" form:"group_description" binding:"required"`
}

type GetCloudHostInVpcSubnetSumReq struct {
	CloudType  string `json:"cloud_type" form:"cloud_type"  binding:"required"` // 如: 腾讯云、阿里云
	Region     string `json:"region" form:"region"  binding:"required"`
	SubnetName string `json:"subnet_name" form:"subnet_name"  binding:"required"`
}

type UpdateCloudInstanceConfigReq struct {
	Id                      uint   `json:"id" form:"id"`
	Region                  string `json:"region" form:"region"  binding:"required"`                                         // 如：ap-guangzhou
	InstanceChargeType      string `json:"instance_charge_type" form:"instance_charge_type"  binding:"required"`             // 如：PREPAID
	Period                  int64  `json:"period" form:"period"  binding:"required"`                                         // 如：1
	RenewFlag               string `json:"renew_flag" form:"renew_flag"  binding:"required"`                                 // 如：NOTIFY_AND_AUTO_RENEW
	ProjectId               int64  `json:"project_id" form:"project_id"  binding:"required"`                                 // 通过GetCloudProjectId接口，提供projectName(一般为项目名)获取值，值如：1308247
	InstanceFamily          string `json:"instance_family" form:"instance_family"  binding:"required"`                       // 如：SA5
	CpuCores                int    `json:"cpu_cores" form:"cpu_cores"  binding:"required"`                                   // 如：4
	MemorySize              int    `json:"memory_size" form:"memory_size"  binding:"required"`                               // 如：16，单位GB
	Fpga                    int    `json:"fpga" form:"fpga"  binding:"required"`                                             // 如：0
	GpuCores                int    `json:"gpu_cores" form:"gpu_cores"  binding:"required"`                                   // 如：0
	ImageId                 string `json:"image_id" form:"image_id"  binding:"required"`                                     // 如：img-8toqc6s3
	SystemDiskType          string `json:"system_disk_type" form:"system_disk_type"  binding:"required"`                     // 如：CLOUD_BASIC
	SystemDiskSize          int64  `json:"system_disk_size" form:"system_disk_size"  binding:"required"`                     // 如：40
	DataDiskType            string `json:"data_disk_type" form:"data_disk_type"  binding:"required"`                         // 如：CLOUD_BASIC
	DataDiskSize            int64  `json:"data_disk_size" form:"data_disk_size"  binding:"required"`                         // 如：200
	VpcId                   string `json:"vpc_id" form:"vpc_id"  binding:"required"`                                         // 通过GetCloudVpcId接口，提供vpcname(一般为项目名)获取值，值如：vpc-0t8v2z9w
	SubnetId                string `json:"subnet_id" form:"subnet_id"  binding:"required"`                                   // 通过GetCloudVpcSubnetId接口，提供通过subnetname(一般为项目名)获取值，值如：subnet-0t8v2z9w
	InternetChargeType      string `json:"internet_charge_type" form:"internet_charge_type"  binding:"required"`             // 如：TRAFFIC_POSTPAID_BY_HOUR
	InternetMaxBandwidthOut int64  `json:"internet_max_bandwidth_out" form:"internet_max_bandwidth_out"  binding:"required"` // 如：200
	InstanceNamePrefix      string `json:"instance_name_prefix" form:"instance_name_prefix"  binding:"required"`             // 如：xhxmj
	SecurityGroupId         string `json:"security_group_id" form:"security_group_id"  binding:"required"`                   // 通过GetCloudSecurityGroupId接口，提供SecurityGroupName(一般为项目名)获取值，值如：sg-0t8v2z9w
	HostNamePrefix          string `json:"host_name_prefix" form:"host_name_prefix"  binding:"required"`                     // 如：xhxmj
}

type GetCloudVpcIdReq struct {
	CloudType string `json:"cloud_type" form:"cloud_type"  binding:"required"` // 如: 腾讯云、阿里云
	Region    string `json:"region" form:"region"  binding:"required"`
	VpcName   string `json:"vpc_name" form:"vpc_name"  binding:"required"` // VpcName与项目同名
}

type GetCloudVpcSubnetIdReq struct {
	CloudType  string `json:"cloud_type" form:"cloud_type"  binding:"required"` // 如: 腾讯云、阿里云
	Region     string `json:"region" form:"region"  binding:"required"`
	SubnetName string `json:"subnet_name" form:"subnet_name"  binding:"required"`
}

type GetCloudSecurityGroupIdReq struct {
	CloudType         string `json:"cloud_type" form:"cloud_type"  binding:"required"` // 如: 腾讯云、阿里云
	Region            string `json:"region" form:"region"  binding:"required"`
	SecurityGroupName string `json:"security_group_name" form:"security_group_name"  binding:"required"`
}

type GetCloudProjectIdReq struct {
	CloudType   string `json:"cloud_type" form:"cloud_type"  binding:"required"` // 如: 腾讯云、阿里云
	ProjectName string `json:"project_name" form:"project_name"  binding:"required"`
}

type GetCloudInstanceTypeListReq struct {
	CloudType      string `json:"cloud_type" form:"cloud_type"  binding:"required"`
	Region         string `json:"region" form:"region"  binding:"required"`
	InstanceFamily string `json:"instance_family" form:"instance_family"  binding:"required"`
	CpuCores       int    `json:"cpu_cores" form:"cpu_cores"  binding:"required"`
	MemorySize     int    `json:"memory_size" form:"memory_size"  binding:"required"`
	Fpga           int    `json:"fpga" form:"fpga"  binding:"required"`
	GpuCores       int    `json:"gpu_cores" form:"gpu_cores"  binding:"required"`
}

type CreateCloudInstanceReq struct {
	CloudType               string `json:"cloud_type" form:"cloud_type"  binding:"required"` // 如: 腾讯云、阿里云
	Region                  string `json:"region" form:"region"  binding:"required"`
	InstanceChargeType      string `json:"instance_charge_type" form:"instance_charge_type"  binding:"required"`
	Period                  int64  `json:"period" form:"period"  binding:"required"`
	RenewFlag               string `json:"renew_flag" form:"renew_flag"  binding:"required"`
	Zone                    string `json:"zone" form:"zone"  binding:"required"`
	ProjectId               int64  `json:"project_id" form:"project_id"  binding:"required"`
	InstanceType            string `json:"instance_type" form:"instance_type"  binding:"required"`
	ImageId                 string `json:"image_id" form:"image_id"  binding:"required"`
	SystemDiskType          string `json:"system_disk_type" form:"system_disk_type"  binding:"required"`
	SystemDiskSize          int64  `json:"system_disk_size" form:"system_disk_size"  binding:"required"`
	DataDiskType            string `json:"data_disk_type" form:"data_disk_type"  binding:"required"`
	DataDiskSize            int64  `json:"data_disk_size" form:"data_disk_size"  binding:"required"`
	VpcId                   string `json:"vpc_id" form:"vpc_id"  binding:"required"`
	SubnetId                string `json:"subnet_id" form:"subnet_id"  binding:"required"`
	InternetChargeType      string `json:"internet_charge_type" form:"internet_charge_type"  binding:"required"`
	InternetMaxBandwidthOut int64  `json:"internet_max_bandwidth_out" form:"internet_max_bandwidth_out"  binding:"required"`
	InstanceName            string `json:"instance_name" form:"instance_name"  binding:"required"`
	SecurityGroupId         string `json:"security_group_id" form:"security_group_id"  binding:"required"`
	HostName                string `json:"host_name" form:"host_name"  binding:"required"`
}

type ReturnCloudInstanceReq struct {
	CloudType  string `json:"cloud_type" form:"cloud_type"  binding:"required"` // 如: 腾讯云、阿里云
	Region     string `json:"region" form:"region"  binding:"required"`
	InstanceId string `json:"instance_id" form:"instance_id"  binding:"required"`
}
