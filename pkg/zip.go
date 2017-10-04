package pkg

import (
	"archive/zip"
	"errors"
	"fmt"
	"github.com/ulikunitz/xz"
	"io"
	"net"
	"net/url"
	"path"
)

// Accepts a UTF-8 encoded archive name.
// Name must be a valid URL. It currently supports URLS which point to files
// with the suffix ".tar.xz".
// This function will not override previously existing files. If it
// encounters one, it will stop processing and return an error.
func InstallMe(archive string, root string) error {

	// Verify that the archive is a valid URL, for "file" or "http(s)"
	urlParts, err := url.Parse(archive)
	if err != nil {
		return err
	}

	// verify that the path ends in ".tar.xz", otherwise give error
	if !strings.HasSuffix(urlParts.Path, ".tar.xz") {
		return fmt.Errorf("Unsupported archive type (not *.tar.xz): %s",
			path.Base(urlParts.Path))
	}

	var r io.Reader

	switch urlParts.Scheme {
	case "file":
		fileReader, err := os.Open(urlParts.Path)
		if err != nil {
			return err
		}
		defer fileReader.Close()
		r = fileReader
	case "https":
		fallthrough
	case "http":
		resp, err := http.Get(archive)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		r = resp.Body

	default:
		return fmt.Errorf("Unsupported URL type: %s", urlParts.Scheme)
	}

	// Decompress the stream
	xzReader, err := xz.NewReader(r)
	if err != nil {
		return err
	}

	// Untar the stream
	tarReader, err := tar.NewReader(xzReader)
	if err != nil {
		return err
	}

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		dest := path.Join(root, header.Name)
		if header.FileInfo().IsDir() {
			os.MkdirAll(dest, header.Mode)
		} else {
			// Below, I use os.Create, which truncates
			// any existing file.
			//
			// I could have checked for existence first,
			// but by then I would have already unpacked lots
			// of archive by then.
			// Perhaps it is far cleaner and smoother in my
			// use case to blow away what's there (if anything)
			// and not apologize for it.
			destWriter, err := os.Create(dest)

			if err != nil {
				return err
			}
			_, err := io.Copy(destWriter, tarReader)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
