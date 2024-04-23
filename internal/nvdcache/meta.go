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

func parseLine(o *Metadata, key, value string) {
	switch key {
	case "lastModifiedDate":
		if ts, err := time.Parse("2006-01-02T15:04:05-07:00", value); err == nil {
			o.LastModifiedDate = ts
		}
	case "size":
		if i, err := strconv.ParseInt(value, 10, 64); err == nil {
			o.Size = i
		}
	case "zipSize":
		if i, err := strconv.ParseInt(value, 10, 64); err == nil {
			o.ZipSize = i
		}
	case "gzSize":
		if i, err := strconv.ParseInt(value, 10, 64); err == nil {
			o.GZSize = i
		}
	case "sha256":
		o.SHA256 = value
	}
}

func NewMetadataFromReader(in io.Reader) *Metadata {
	o := &Metadata{}
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, in); err != nil {
		return nil
	}
	o.data = buf.Bytes()
	scanner := bufio.NewScanner(buf)

	for scanner.Scan() {
		sp := strings.SplitN(scanner.Text(), ":", 2) //nolint:gomnd // key-value.
		if len(sp) == 2 {                            //nolint:gomnd // key-value.
			parseLine(o, sp[0], sp[1])
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
