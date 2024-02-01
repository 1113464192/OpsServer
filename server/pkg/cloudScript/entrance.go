package cloudScript

import (
	"fmt"
	"fqhWeb/pkg/cloudScript/tencentCloud"
)

func UpdateCloudProjectEntry(cloudType string, projectId uint64, projectName string, disable int64) (err error) {
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

func CreateCloudProjectEntry(cloudType string, projectName string) (err error) {
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
