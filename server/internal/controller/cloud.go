package controller

import (
	"fqhWeb/internal/service"
	"fqhWeb/pkg/api"
	"fqhWeb/pkg/logger"
	"github.com/gin-gonic/gin"
	"strconv"
)

// CreateCloudVpc
// @Tags 云平台相关
// @title 创建项目VPC
// @description 返回创建是否成功
// @Summary 创建项目VPC
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data formData api.CreateCloudVpcReq true "按对应云的API文档填写,VpcName与项目同名"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/cloud/vpc [post]
func CreateCloudVpc(c *gin.Context) {
	var vpcReq api.CreateCloudVpcReq
	if err := c.ShouldBind(&vpcReq); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	if err := service.Cloud().CreateCloudVpc(vpcReq.CloudType, vpcReq.Region, vpcReq.VpcName, vpcReq.CidrBlock); err != nil {
		logger.Log().Error("Cloud", "新增VPC失败", err)
		c.JSON(500, api.Err("新增VPC失败", err))
		return
	}
	c.JSON(200, api.Response{
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}

// CreateCloudVpcSubnet
// @Tags 云平台相关
// @title 创建项目VPCSubnet
// @description 返回创建是否成功
// @Summary 创建项目VPCSubnet
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data formData api.CreateCloudVpcSubnetReq true "按对应云的API文档填写，Name一般为ProjectName-01以此类推"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/cloud/vpc-subnet [post]
func CreateCloudVpcSubnet(c *gin.Context) {
	var subnetReq api.CreateCloudVpcSubnetReq
	if err := c.ShouldBind(&subnetReq); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	vpcId, err := service.Cloud().GetCloudVpcId(subnetReq.CloudType, subnetReq.Region, subnetReq.VpcName)
	if err != nil {
		logger.Log().Error("Cloud", "获取VPCId失败", err)
		c.JSON(500, api.Err("获取VPCId失败", err))
		return
	}
	if err = service.Cloud().CreateCloudVpcSubnet(subnetReq.CloudType, subnetReq.Region, vpcId, subnetReq.SubnetName, subnetReq.SubnetCidrBlock); err != nil {
		logger.Log().Error("Cloud", "新增VPCSubnet失败", err)
		c.JSON(500, api.Err("新增VPCSubnet失败", err))
		return
	}
	c.JSON(200, api.Response{
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}

// CreateCloudSecurityGroup
// @Tags 云平台相关
// @title 创建项目安全组
// @description 返回创建是否成功
// @Summary 创建项目安全组
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data formData api.CreateCloudSecurityGroupReq true "按对应云的API文档填写，Name(一般与项目名同名)"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/cloud/security-group [post]
func CreateCloudSecurityGroup(c *gin.Context) {
	var securityGroupReq api.CreateCloudSecurityGroupReq
	if err := c.ShouldBind(&securityGroupReq); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	projectId, err := service.Cloud().GetCloudProjectId(securityGroupReq.CloudType, securityGroupReq.ProjectName)
	if err != nil {
		logger.Log().Error("Cloud", "获取ProjectId失败", err)
		c.JSON(500, api.Err("获取ProjectId失败", err))
		return
	}
	// 让projectId变成字符串
	pidStr := strconv.FormatUint(projectId, 10)

	if err = service.Cloud().CreateCloudSecurityGroup(securityGroupReq.CloudType, securityGroupReq.Region, pidStr, securityGroupReq.GroupName, securityGroupReq.GroupDescription); err != nil {
		logger.Log().Error("Cloud", "创建安全组失败", err)
		c.JSON(500, api.Err("创建安全组失败", err))
		return
	}
	c.JSON(200, api.Response{
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}

// GetCloudHostInVpcSubnetSum
// @Tags 云平台相关
// @title 项目机器占指定VPCSubnet总数
// @description 返回子网包含IP总数(<=255)
// @Summary 项目机器占指定VPCSubnet总数
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data query api.GetCloudHostInVpcSubnetSumReq true "按对应云的API文档填写"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/cloud/host-in-subnet-sum [get]
func GetCloudHostInVpcSubnetSum(c *gin.Context) {
	var getSumReq api.GetCloudHostInVpcSubnetSumReq
	if err := c.ShouldBind(&getSumReq); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	subnetId, err := service.Cloud().GetCloudVpcSubnetId(getSumReq.CloudType, getSumReq.Region, getSumReq.SubnetName)
	if err != nil {
		logger.Log().Error("Cloud", "获取SubnetId失败", err)
		c.JSON(500, api.Err("获取SubnetId失败", err))
		return
	}
	sum, err := service.Cloud().GetCloudHostInVpcSubnetSum(getSumReq.CloudType, getSumReq.Region, subnetId)
	if err != nil {
		logger.Log().Error("Cloud", "获取总数失败", err)
		c.JSON(500, api.Err("获取总数失败", err))
		return
	}
	c.JSON(200, api.Response{
		Data: map[string]int64{
			"sum": sum,
		},
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}

// UpdateCloudInstanceConfig
// @Tags 云平台相关
// @title 创建/更新云项目的实例配置
// @description 返回新的项目指定创建实例的配置Json
// @Summary 创建/更新云项目的实例配置
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data formData api.GetCloudHostInVpcSubnetSumReq true "按对应云的API文档填写"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/cloud/instance-config [put]
func UpdateCloudInstanceConfig(c *gin.Context) {
	var configReq api.UpdateCloudInstanceConfigReq
	if err := c.ShouldBind(&configReq); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	res, err := service.Cloud().UpdateCloudInstanceConfig(&configReq)
	if err != nil {
		logger.Log().Error("Cloud", "创建/更新云项目的实例配置失败", err)
		c.JSON(500, api.Err("创建/更新云项目的实例配置失败", err))
		return
	}
	c.JSON(200, api.Response{
		Data: res,
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}

// GetCloudInstanceConfig
// @Tags 云平台相关
// @title 获取云项目的实例配置
// @description 返回指定云项目的实例配置，P.S: 输入项目ID(不是云项目ID)
// @Summary 获取云项目的实例配置
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data query api.IdReq true "输入项目ID(不是云项目ID)"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/cloud/instance-config [get]
func GetCloudInstanceConfig(c *gin.Context) {
	var param api.IdReq
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	res, err := service.Cloud().GetCloudInstanceConfig(param.Id)
	if err != nil {
		logger.Log().Error("Cloud", "获取云项目的实例配置失败", err)
		c.JSON(500, api.Err("获取云项目的实例配置失败", err))
		return
	}
	c.JSON(200, api.Response{
		Data: res,
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}

// GetCloudInsInfo
// @Tags 云平台相关
// @title 获取云实例的详细信息
// @description 返回云实例的详细信息，P.S: 可选输入云项目ID(不是项目ID)
// @Summary 获取云实例的详细信息
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data query api.GetCloudInsInfoReq true "按对应云的API文档填写"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/cloud/instance [get]
func GetCloudInsInfo(c *gin.Context) {
	var param api.GetCloudInsInfoReq
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	res, err := service.Cloud().GetCloudInsInfo(param.CloudType, param.Region, param.PublicIpv4, param.PublicIpv6, param.InsName, param.CloudPid, param.Offset, param.Limit)
	if err != nil {
		logger.Log().Error("Cloud", "获取云实例信息失败", err)
		c.JSON(500, api.Err("获取云实例信息失败", err))
		return
	}
	c.JSON(200, api.Response{
		Data: res,
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}

// GetCloudVpcId
// @Tags 云平台相关
// @title 获取VpcId
// @description 返回VpcId
// @Summary 获取VpcId
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data query api.GetCloudVpcIdReq true "按对应云的API文档填写，VpcName(一般与项目名同名)"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/cloud/vpc-id [get]
func GetCloudVpcId(c *gin.Context) {
	var vpcIdReq api.GetCloudVpcIdReq
	if err := c.ShouldBind(&vpcIdReq); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	vpcId, err := service.Cloud().GetCloudVpcId(vpcIdReq.CloudType, vpcIdReq.Region, vpcIdReq.VpcName)
	if err != nil {
		logger.Log().Error("Cloud", "获取VpcId失败", err)
		c.JSON(500, api.Err("获取VpcId失败", err))
		return
	}
	c.JSON(200, api.Response{
		Data: map[string]string{
			"vpcId": vpcId,
		},
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}

// GetCloudVpcSubnetId
// @Tags 云平台相关
// @title 获取VpcSubnetId
// @description 返回VpcSubnetId
// @Summary 获取VpcSubnetId
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data query api.GetCloudVpcSubnetIdReq true "按对应云的API文档填写，Name一般为ProjectName-01/VpcName-01以此类推"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/cloud/vpc-subnet-id [get]
func GetCloudVpcSubnetId(c *gin.Context) {
	var subnetIdReq api.GetCloudVpcSubnetIdReq
	if err := c.ShouldBind(&subnetIdReq); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	subnetId, err := service.Cloud().GetCloudVpcSubnetId(subnetIdReq.CloudType, subnetIdReq.Region, subnetIdReq.SubnetName)
	if err != nil {
		logger.Log().Error("Cloud", "获取VpcSubnetId失败", err)
		c.JSON(500, api.Err("获取VpcSubnetId失败", err))
		return
	}
	c.JSON(200, api.Response{
		Data: map[string]string{
			"subnetId": subnetId,
		},
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}

// GetCloudSecurityGroupId
// @Tags 云平台相关
// @title 获取安全组ID
// @description 返回安全组ID
// @Summary 获取安全组ID
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data query api.GetCloudSecurityGroupIdReq true "按对应云的API文档填写，Name一般为项目名"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/cloud/security-group-id [get]
func GetCloudSecurityGroupId(c *gin.Context) {
	var idReq api.GetCloudSecurityGroupIdReq
	if err := c.ShouldBind(&idReq); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	id, err := service.Cloud().GetCloudVpcSubnetId(idReq.CloudType, idReq.Region, idReq.SecurityGroupName)
	if err != nil {
		logger.Log().Error("Cloud", "获取安全组ID失败", err)
		c.JSON(500, api.Err("获取安全组ID失败", err))
		return
	}
	c.JSON(200, api.Response{
		Data: map[string]string{
			"id": id,
		},
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}

// GetCloudProjectId
// @Tags 云平台相关
// @title 获取ProjectId
// @description 返回ProjectId
// @Summary 获取ProjectId
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data query api.GetCloudProjectIdReq true "按对应云的API文档填写"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/cloud/project-id [get]
func GetCloudProjectId(c *gin.Context) {
	var idReq api.GetCloudProjectIdReq
	if err := c.ShouldBind(&idReq); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	id, err := service.Cloud().GetCloudProjectId(idReq.CloudType, idReq.ProjectName)
	if err != nil {
		logger.Log().Error("Cloud", "获取项目云ID失败", err)
		c.JSON(500, api.Err("获取项目云ID失败", err))
		return
	}
	c.JSON(200, api.Response{
		Data: map[string]uint64{
			"id": id,
		},
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}

// GetCloudInstanceTypeList
// @Tags 云平台相关
// @title 获取instance-type
// @description 返回instance-type及对应zone
// @Summary 获取instance-type
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data query api.GetCloudInstanceTypeListReq true "按对应云的API文档填写"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/cloud/instance-type [get]
func GetCloudInstanceTypeList(c *gin.Context) {
	var param api.GetCloudInstanceTypeListReq
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	res, err := service.Cloud().GetCloudInstanceTypeList(param.CloudType, param.Region, param.InstanceFamily, param.CpuCores, param.MemorySize, param.Fpga, param.GpuCores)
	if err != nil {
		logger.Log().Error("Cloud", "获取instance-type失败", err)
		c.JSON(500, api.Err("获取instance-type失败", err))
		return
	}
	c.JSON(200, api.Response{
		Data: res,
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}

// CreateCloudInstance
// @Tags 云平台相关
// @title 创建云服务器
// @description 返回是否成功创建，P.S: 一般不需要使用，由装服判断自动创建就好
// @Summary 创建云服务器
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data formData api.CreateCloudInstanceReq true "通过GetCloudInstanceConfig接口和GetCloudInstanceTypeList接口获取后输入"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/cloud/create-instance [post]
func CreateCloudInstance(c *gin.Context) {
	var param api.CreateCloudInstanceReq
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	err := service.Cloud().CreateCloudInstance(param.CloudType, param.Region, param.InstanceChargeType, param.Period, param.RenewFlag, param.Zone, param.ProjectId, param.InstanceType, param.ImageId, param.SystemDiskType, param.SystemDiskSize, param.DataDiskType, param.DataDiskSize, param.VpcId, param.SubnetId, param.InternetChargeType, param.InternetMaxBandwidthOut, param.InstanceName, param.SecurityGroupId, param.HostName)
	if err != nil {
		logger.Log().Error("Cloud", "创建云instance失败", err)
		c.JSON(500, api.Err("创建云instance失败", err))
		return
	}
	c.JSON(200, api.Response{
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}

// ReturnCloudInstance
// @Tags 云平台相关
// @title 退还云服务器
// @description 返回是否成功退还
// @Summary 退还云服务器
// @Produce  application/json
// @Param Authorization header string true "格式为：Bearer 用户令牌"
// @Param data formData api.ReturnCloudInstanceReq true "通过GetCloudInsInfo接口获取InsId后输入"
// @Success 200 {object} api.Response "{"data":{},"meta":{msg":"Success"}}"
// @Failure 401 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 403 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Failure 500 {object} api.Response "{"data":{}, "meta":{"msg":"错误信息", "error":"错误格式输出(如存在)"}}"
// @Router /api/v1/cloud/return-instance [post]
func ReturnCloudInstance(c *gin.Context) {
	var param api.ReturnCloudInstanceReq
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(500, api.ErrorResponse(err))
		return
	}
	err := service.Cloud().ReturnCloudInstance(param.CloudType, param.Region, param.InstanceId)
	if err != nil {
		logger.Log().Error("Cloud", "退还instance失败", err)
		c.JSON(500, api.Err("退还instance失败", err))
		return
	}
	c.JSON(200, api.Response{
		Meta: api.Meta{
			Msg: "Success",
		},
	})
}
