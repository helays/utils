package cfg_tls

import "crypto/tls"

type TLSVersion string

const (
	TLSVersionTLSv1_0 TLSVersion = "TLSv1.0"
	TLSVersionTLSv1_1 TLSVersion = "TLSv1.1"
	TLSVersionTLSv1_2 TLSVersion = "TLSv1.2"
	TLSVersionTLSv1_3 TLSVersion = "TLSv1.3"
)

func (t TLSVersion) Version() uint16 {
	switch t {
	case TLSVersionTLSv1_0:
		return tls.VersionTLS10
	case TLSVersionTLSv1_1:
		return tls.VersionTLS11
	case TLSVersionTLSv1_2:
		return tls.VersionTLS12
	case TLSVersionTLSv1_3:
		return tls.VersionTLS13
	default:
		return tls.VersionTLS12
	}
}
