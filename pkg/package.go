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
	"archive/tar"
	"fmt"
	"github.com/ulikunitz/xz"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
)

type Pask interface {
	Install(root string) error
	Run(root string, task string) error
}

type Package struct {
	Name     string `yaml:"name"`
	Version  string `yaml:"version"`
	Location string `yaml:"location"`
}

type Spec struct {
	Packages []Package `yaml:"packages"`
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
		return nil, fmt.Errorf("Unsupported archive type (not *.tar.xz): `%s`",
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
		return nil, fmt.Errorf("Unsupported URL type: `%s`", urlParts.Scheme)
	}

	return r, nil
}

func openTar(r io.Reader) (*tar.Reader, error) {
	// Decompress the stream
	if xzReader, err := xz.NewReader(r); err != nil {
		return nil, err
	} else {
		// Untar the stream
		tarReader := tar.NewReader(xzReader)
		return tarReader, nil
	}
}

func (p *Package) unpack(tarReader *tar.Reader, root string) error {
	findPaskFiles := regexp.MustCompile("^/*" + ControlDirectory + "/+")
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
			relativeDest := findPaskFiles.ReplaceAllString(header.Name, "")
			dest = path.Join(root, ControlDirectory, "packages", p.Name,
				p.Version, relativeDest)
		} else {
			dest = path.Join(root, header.Name)
		}

		log.Printf("Attempting to unpack file `%s`\n", dest)

		mode := header.FileInfo().Mode()
		if len(dest) == 0 {
			return fmt.Errorf("Empty destination file")
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if fi, err := os.Stat(dest); err != nil {
				if os.IsNotExist(err) {
					if err := os.MkdirAll(dest, mode); err != nil {
						return err
					}
				} else {
					return err
				}
			} else {
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
	return nil
}

// Accepts a UTF-8 encoded archive name.
// Name must be a valid URL. It currently supports URLS which point to files
// with the suffix ".tar.xz".
// This function will not override previously existing files. If it
// encounters one, it will stop processing and return an error.
// Runs the `postinst` task, if it exists, after installation is complete.
// TODO: add logging
func (p *Package) Install(root string) error {

	var tarReader *tar.Reader
	log.Printf("Opening archive at `%s`\n", p.Location)
	if rc, err := openArchive(p.Location); err != nil {
		return err
	} else {
		defer rc.Close()
		if r, err := openTar(rc); err != nil {
			return err
		} else {
			tarReader = r
		}
	}
	if err := p.unpack(tarReader, root); err != nil {
		return err
	}

	if err := p.Run(root, "pretmpl"); err != nil {
		return err
	}

	/*	if err := p.Template(); err != nil {
		return err
	} */

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

func (s *Spec) Run(root string, task string) error {
	for _, pkg := range s.Packages {
		if err := pkg.Run(root, task); err != nil {
			return err
		}
	}
	return nil
}

func (s *Spec) Install(root string) error {
	for _, pkg := range s.Packages {
		log.Printf("Installing `%s`\n", pkg)
		log.Printf("Installing archive `%s` from location `%s`\n",
			pkg.Name, pkg.Location)
		if err := pkg.Install(root); err != nil {
			return err
		}
	}
	return nil
}

func ReadSpec(sfile string) (*Spec, error) {
	var spec Spec
	if specData, err := ioutil.ReadFile(sfile); err != nil {
		return nil, fmt.Errorf("Couldn't read spec")
	} else {
		if err := yaml.Unmarshal(specData, &spec); err != nil {
			return nil, fmt.Errorf("Couldn't parse spec file: %v", err)
		}
		log.Printf("Spec input is as follows:\n\n%s\n", specData)
		log.Printf("Spec is as follows:\n\n%s\n", spec)
	}
	return &spec, nil
}
