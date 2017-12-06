package manifest

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

var fixture = `
language = "foo"
status = ""

[documentation]
  description = "foo"

[runtime]
  os = "alpine"
  native_version = ["42"]
  go_version = "1.8"
`[1:]

func TestEncode(t *testing.T) {
	m := &Manifest{}
	m.Language = "foo"
	m.Documentation.Description = "foo"
	m.Runtime.OS = Alpine
	m.Runtime.GoVersion = "1.8"
	m.Runtime.NativeVersion = []string{"42"}

	buf := bytes.NewBuffer(nil)
	err := m.Encode(buf)
	assert.Nil(t, err)

	assert.Equal(t, fixture, buf.String())
}

func TestDecode(t *testing.T) {
	m := &Manifest{}

	buf := bytes.NewBufferString(fixture)
	err := m.Decode(buf)
	assert.Nil(t, err)

	assert.Equal(t, "foo", m.Language)
	assert.Equal(t, Alpine, m.Runtime.OS)
}
