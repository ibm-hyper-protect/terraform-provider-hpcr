package encrypt

import (
	"fmt"
	"testing"

	B "github.com/terraform-provider-hpcr/fp/bytes"
	E "github.com/terraform-provider-hpcr/fp/either"
	F "github.com/terraform-provider-hpcr/fp/function"
)

func TestVersion(t *testing.T) {

	res := openSSLVersion

	fmt.Println(res)
}

func TestRandomPassword(t *testing.T) {

	genPwd := RandomPassword(32)

	pwd := genPwd()

	fmt.Println(pwd)
}

func TestEncryptPassword(t *testing.T) {

	//	genPwd := RandomPassword(32)

}

func TestPrivateKey(t *testing.T) {
	privKey := PrivateKey()

	pubKey := F.Pipe2(
		privKey,
		E.Chain(PublicKey),
		E.Map[error](B.ToString),
	)

	fmt.Println(pubKey)
}
