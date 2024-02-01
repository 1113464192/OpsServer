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
		ak: configs.Conf.CloudSecretKey.TencentCloud.Ak,
		sk: configs.Conf.CloudSecretKey.TencentCloud.Sk,
	}
	return insTencentCloud
}
