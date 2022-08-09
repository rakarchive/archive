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
	"io/fs"
	"os"
	"path/filepath"
)

func Zip(src string, dst io.Writer) error {
	writer := zip.NewWriter(dst)
	defer writer.Close()

	return filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		return addFile(writer, src, path, info)
	})
}

func addFile(writer *zip.Writer, src, path string, info fs.FileInfo) error {
	if info.IsDir() {
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	if header.Name, err = filepath.Rel(filepath.Dir(src), path); err != nil {
		return err
	}
	header.Method = zip.Deflate

	headerWriter, err := writer.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(headerWriter, file)
	return err
}
