package request

import (
	"bytes"
	"net/textproto"
)

type Files []File

type File struct {
	Filename string
	Size     int64
	Header   textproto.MIMEHeader
	Body     *bytes.Buffer
}
