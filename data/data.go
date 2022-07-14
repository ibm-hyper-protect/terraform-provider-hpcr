package data

import (
	_ "embed"
)

//go:embed ibm-hyper-protect-container-runtime-1-0-s390x-1-encrypt.crt
var DefaultCertificate string
