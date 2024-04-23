package nvdcache

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/na4ma4/config"
)

const (
	baseURL            = "https://nvd.nist.gov/feeds/json/cve/1.1/"
	defaultHTTPTimeout = 10 * time.Second
)

type NVD struct {
	cfg      config.Conf
	baseURL  *url.URL
	fileData string
	fileMeta string
}

func NewNVDByYear(cfg config.Conf, year int) *NVD {
	u, _ := url.Parse(baseURL)
	return &NVD{
		cfg:      cfg,
		baseURL:  u,
		fileData: fmt.Sprintf("nvdcve-1.1-%d.json.gz", year),
		fileMeta: fmt.Sprintf("nvdcve-1.1-%d.meta", year),
	}
}

func (n *NVD) fileNameMeta() string {
	return n.cfg.GetString("nvd.cache") + string(os.PathSeparator) + n.fileMeta
}

func (n *NVD) fileNameData() string {
	return n.cfg.GetString("nvd.cache") + string(os.PathSeparator) + n.fileData
}

func (n *NVD) CacheMeta() (*Metadata, error) {
	f, err := os.Open(n.fileNameMeta())
	if err != nil {
		return &Metadata{}, err
	}

	o := NewMetadataFromReader(f)
	if o == nil {
		return &Metadata{}, errors.New("unable to read cached metadata file")
	}

	return o, nil
}

func (n *NVD) LiveMeta(ctx context.Context) (*Metadata, error) {
	u := n.baseURL.ResolveReference(&url.URL{Path: n.fileMeta})
	log.Printf("URL: %s", u.String())
	r, err := n.getURL(ctx, u)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	md := NewMetadataFromReader(r.Body)
	return md, nil
}

func (n *NVD) getURL(ctx context.Context, u *url.URL) (*http.Response, error) {
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	c := &http.Client{
		Timeout: defaultHTTPTimeout,
	}

	return c.Do(r)
}

//nolint:nestif // TODO cleanup.
func (n *NVD) Download(ctx context.Context) error {
	var lmd *Metadata
	{
		var err error
		lmd, err = n.LiveMeta(ctx)
		if err != nil {
			return err
		}
	}

	var md *Metadata
	{
		var err error
		md, err = n.CacheMeta()
		if err != nil {
			return err
		}
	}

	if !md.Compare(lmd) {
		log.Println("Metadata does not match, need to download")

		u := n.baseURL.ResolveReference(&url.URL{Path: n.fileData})
		log.Printf("URL: %s", u.String())

		var r *http.Response
		var f *os.File
		{
			var err error
			r, err = n.getURL(ctx, u)
			if err != nil {
				return err
			}
			defer r.Body.Close()
		}

		{
			var err error
			f, err = os.Create(n.fileNameData())
			if err != nil {
				return err
			}
			defer f.Close()
		}

		if _, err := io.Copy(f, r.Body); err != nil {
			return err
		}

		{
			var err error
			f, err = os.Create(n.fileNameMeta())
			if err != nil {
				return err
			}
			defer f.Close()
		}

		if err := lmd.Write(f); err != nil {
			return err
		}
	} else {
		log.Println("Metadata matches, no need to download")
	}

	return nil
}
