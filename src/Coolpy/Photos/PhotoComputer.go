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

func parseRange(data string) int64 {
	stop := (int64)(0)
	part := 0
	for i := 0; i < len(data) && part < 2; i = i + 1 {
		if part == 0 {
			// part = 0 <=> equal isn't met.
			if data[i] == '=' {
				part = 1
			}
			continue
		}
		if part == 1 {
			// part = 1 <=> we've met the equal, parse beginning
			if data[i] == ',' || data[i] == '-' {
				part = 2 // part = 2 <=> OK DUDE.
			} else {
				if 48 <= data[i] && data[i] <= 57 {
					// If it's a digit ...
					// ... convert the char to integer and add it!
					stop = (stop * 10) + (((int64)(data[i])) - 48)
				} else {
					part = 2 // Parsing error! No error needed : 0 = from start.
				}
			}
		}
	}
	return stop
}
