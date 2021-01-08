package nvdcache

import (
	"bufio"
	"bytes"
	"io"
	"strconv"
	"strings"
	"time"
)

type Metadata struct {
	data             []byte
	LastModifiedDate time.Time
	Size             int64
	ZipSize          int64
	GZSize           int64
	SHA256           string
}

func NewMetadataFromReader(in io.Reader) *Metadata {
	o := &Metadata{}
	buf := bytes.NewBuffer(nil)
	_, err := io.Copy(buf, in)
	if err != nil {
		return nil
	}
	o.data = buf.Bytes()
	scanner := bufio.NewScanner(buf)

	for scanner.Scan() {
		sp := strings.SplitN(scanner.Text(), ":", 2)
		if len(sp) == 2 {
			switch sp[0] {
			case "lastModifiedDate":
				if ts, err := time.Parse("2006-01-02T15:04:05-07:00", sp[1]); err == nil {
					o.LastModifiedDate = ts
				}
			case "size":
				if i, err := strconv.ParseInt(sp[1], 10, 64); err == nil {
					o.Size = i
				}
			case "zipSize":
				if i, err := strconv.ParseInt(sp[1], 10, 64); err == nil {
					o.ZipSize = i
				}
			case "gzSize":
				if i, err := strconv.ParseInt(sp[1], 10, 64); err == nil {
					o.GZSize = i
				}
			case "sha256":
				o.SHA256 = sp[1]
			}
		}
	}

	return o
}

func (md *Metadata) Compare(o *Metadata) bool {
	return strings.EqualFold(md.SHA256, o.SHA256)
}

func (md *Metadata) Write(out io.Writer) error {
	buf := bytes.NewBuffer(nil)
	if _, err := buf.Write(md.data); err != nil {
		return err
	}
	if _, err := io.Copy(out, buf); err != nil {
		return err
	}

	return nil
}
