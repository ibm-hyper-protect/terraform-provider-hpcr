package archive

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	E "github.com/terraform-provider-hpcr/fp/either"
	F "github.com/terraform-provider-hpcr/fp/function"
)

func TestTgz(t *testing.T) {
	var body bytes.Buffer

	base64 := base64.NewEncoder(base64.StdEncoding, &body)

	resE := F.Pipe3(
		base64,
		TarFolder[io.WriteCloser]("../samples/nginx-golang"),
		E.Chain(onClose[io.WriteCloser]),
		E.MapTo[error, any](true),
	)

	assert.Equal(t, E.Of[error](true), resE)

	fmt.Println(body.String())
}
