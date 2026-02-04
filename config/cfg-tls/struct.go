package cfg_tls

type TLS struct {
	MinVersion   TLSVersion    `json:"min_version" yaml:"min_version" ini:"min_version"`
	MaxVersion   TLSVersion    `json:"max_version" yaml:"max_version" ini:"max_version"`
	CipherSuites []CipherSuite `json:"cipher_suites" yaml:"cipher_suites" ini:"cipher_suites"`
}
