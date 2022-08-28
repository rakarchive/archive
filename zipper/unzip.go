// Copyright © 2022 Rak Laptudirm <rak@laptudirm.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package zipper

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func Unzip(src io.ReaderAt, size int64, dst string) error {
	r, err := zip.NewReader(src, size)
	if err != nil {
		return err
	}

	for _, f := range r.File {
		path := filepath.Join(dst, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			return err
		}

		file, err := os.Create(path)
		if err != nil {
			return err
		}

		unzipped, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(file, unzipped)
		if err != nil {
			return err
		}

		file.Close()
		unzipped.Close()
	}

	return nil
}
