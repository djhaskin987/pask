// Copyright © 2017 Daniel Jay Haskin <djhaskin987@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
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
	"re"
)

type Package struct {
	Name     string `yaml:"name"`
	Version  string `yaml:"version"`
	Location string `yaml:"location"`
}

// Name of the directory relative to the root where control files will
// be installed. Also the name of the directory relative to the root
// of package archives where control files will be looked for.
const ControlDirectory = "pask"

func openArchive(archive string) (io.ReadCloser, error) {

	var r io.ReadCloser
	// Verify that the archive is a valid URL, for "file" or "http(s)"
	urlParts, err := url.Parse(archive)
	if err != nil {
		return nil, err
	}

	// verify that the path ends in ".tar.xz", otherwise give error
	if !strings.HasSuffix(urlParts.Path, ".tar.xz") {
		return nil, fmt.Errorf("Unsupported archive type (not *.tar.xz): %s",
			path.Base(urlParts.Path))
	}

	switch urlParts.Scheme {
	case "file":
		if fileReader, err := os.Open(urlParts.Path); err != nil {
			return nil, err
		} else {
			r = fileReader
		}
	case "https":
		fallthrough
	case "http":
		if resp, err := http.Get(archive); err != nil {
			return nil, err
		} else {
			r = resp.Body
		}

	default:
		return nil, fmt.Errorf("Unsupported URL type: %s", urlParts.Scheme)
	}

	return r, nil
}

func untar(r io.Reader) (io.Reader, error) {
	// Decompress the stream
	if xzReader, err := xz.NewReader(r); err != nil {
		return r, err
	} else {
		// Untar the stream
		if tarReader, err := tar.NewReader(xzReader); err != nil {
			return xzReader, err
		} else {
			return tarReader, nil
		}
	}
}

func unpack(tarReader io.Reader, name string, version string, root string) error {
	findPaskFiles := re.MustCompile("^/*" + ControlDirectory + "/+")
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		} else if header == nil {
			continue
		}

		var dest string
		if results := findPaskFiles.FindStringSubmatch(header.Name); results != nil {
			dest = path.Join(root, ControlDirectory, "packages", name, version)
		} else {
			dest = path.Join(root, header.Name)
		}

		mode := header.FileInfo().Mode()
		if len(dest) == 0 {
			return fmt.Errorf("Empty destination file")
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if fi, err := os.Stat(dest); err == nil {
				if fi.IsDir() {
					os.Chmod(dest, mode)
				} else {
					if err := os.Remove(dest); err != nil {
						return err
					}
					if os.MkdirAll(dest, mode); err != nil {
						return err
					}
				}
			} else if os.IsNotExist(err) {
				if os.MkdirAll(dest, mode); err != nil {
					return err
				}
			} else if err != nil {
				return err
			}
		case tar.TypeReg:
			// Perhaps it is far cleaner and smoother in my
			// use case to blow away what's there (if anything)
			// and not apologize for it.
			if fi, err := os.Stat(dest); err == nil {
				if fi.IsDir() {
					os.RemoveAll(dest)
				} else {
					os.Remove(dest)
				}
			}

			if destWriter, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode); err != nil {
				return err
			} else {
				if _, err := io.Copy(destWriter, tarReader); err != nil {
					return err
				}
			}
		case tar.TypeSymlink:
			if err := os.Symlink(header.Linkname, dest); err != nil {
				return err
			}
		default:
			continue
		}
	}
}

// Accepts a UTF-8 encoded archive name.
// Name must be a valid URL. It currently supports URLS which point to files
// with the suffix ".tar.xz".
// This function will not override previously existing files. If it
// encounters one, it will stop processing and return an error.
// Runs the `postinst` task, if it exists, after installation is complete.
// TODO: add logging
func (p *Package) Install(root string) error {
	var tarReader io.Reader

	if rc, err := openArchive(p.Location); err != nil {
		return err
	} else {
		defer rc.Close()
		if r, err := untar(rc); err != nil {
			return err
		} else {
			tarReader = r
		}
	}
	if err := unpack(tarReader, name, version, root); err != nil {
		return err
	}

	// TODO: Templating?

	if err := p.Run(root, "postinst"); err != nil {
		return err
	}

	return nil
}

// Locates a file called
// <root>/<ControlDirectory>/packages/<name>/<version>/tasks/<task>. If it is
// executable, run it. If it is not, return an error. Return any errors along
// the way.
func (p *Package) Run(root string, task string) error {
	return nil
}
