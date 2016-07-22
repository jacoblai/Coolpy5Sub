package Photos

import (
	"bytes"
)

var mimeTable = map[string]string{
	"\xff\xd8\xff":      "image/jpeg",
	"\x89PNG\r\n\x1a\n": "image/png",
	"GIF87a":            "image/gif",
	"GIF89a":            "image/gif",
}

func mimeComput(incipit []byte) string {
	for magic, mime := range mimeTable {
		if bytes.HasPrefix(incipit, []byte(magic)) {
			return mime
		}
	}
	return ""
}
