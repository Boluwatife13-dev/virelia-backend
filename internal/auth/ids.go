package auth

func (s *Service) HashNIN(nin string) string {
	return hashCode(nin, s.secret)
}

func (s *Service) NewNetworkID() string {
	return newNetworkID()
}
