package models

import "bytes"

type InMemFile struct {
	FileName string
	Buffer   bytes.Buffer
}
