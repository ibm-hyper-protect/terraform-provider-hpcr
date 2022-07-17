//
// Licensed Materials - Property of IBM
//
// 5737-I09
//
// Copyright IBM Corp. 2022 All Rights Reserved.
// US Government Users Restricted Rights - Use, duplication or
// disclosure restricted by GSA ADP Schedule Contract with IBM Corp
//
package data

import (
	_ "embed"
)

//go:embed ibm-hyper-protect-container-runtime-1-0-s390x-1-encrypt.crt
var DefaultCertificate string
