package tencentCloud

func (s *TencentCloudService) CreateCloudHost(region string) {
	zoneId, err := s.GetAvailCloudZoneId(region)
	if err != nil {
		return err
	}
}
