package tencentCloud

import "fqhWeb/configs"

type TencentCloudService struct {
	ak string
	sk string
}

var (
	insTencentCloud = &TencentCloudService{}
)

func TencentCloud() *TencentCloudService {
	insTencentCloud = &TencentCloudService{
		ak: configs.Conf.Cloud.TencentCloud.Ak,
		sk: configs.Conf.Cloud.TencentCloud.Sk,
	}
	return insTencentCloud
}
