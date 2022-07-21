// Copyright 2022 IBM Corp.
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
package archive

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"

	E "github.com/terraform-provider-hpcr/fp/either"
	F "github.com/terraform-provider-hpcr/fp/function"
	I "github.com/terraform-provider-hpcr/fp/identity"
	T "github.com/terraform-provider-hpcr/fp/tuple"
)

var (
	fileInfoHeaderE = E.Eitherize2(tar.FileInfoHeader)
	relE            = E.Eitherize2(filepath.Rel)
	openE           = E.Eitherize1(os.Open)
	copyE           = E.Eitherize2(io.Copy)
	skipDir         = E.Of[error, int64](-1)
)

func toReader[A io.Reader](a A) io.Reader {
	return a
}

func onClose[A io.Closer](a A) E.Either[error, any] {
	err := a.Close()
	return E.TryCatchError(func() (any, error) {
		return nil, err
	})
}

func onOpenFile(file string) func() E.Either[error, *os.File] {
	return func() E.Either[error, *os.File] {
		return F.Pipe2(
			file,
			filepath.Clean,
			openE,
		)
	}
}

// constructs a function that copies the content of a file into the writer
func copyFile(w io.Writer) func(string, os.FileInfo) E.Either[error, int64] {

	copyTo := F.Bind1st(copyE, w)

	return func(file string, fi os.FileInfo) E.Either[error, int64] {
		// do not copy for directories
		if fi.IsDir() {
			return skipDir
		}
		// copy the file content
		return F.Pipe1(
			E.WithResource[error, *os.File, int64](onOpenFile(file), onClose[*os.File]),
			I.Ap[func(*os.File) E.Either[error, int64], E.Either[error, int64]](F.Flow2(
				toReader[*os.File],
				copyTo,
			)),
		)
	}
}

func writeHeader(src string) func(*tar.Writer) func(file string, fi os.FileInfo) E.Either[error, *tar.Writer] {
	// callback to get the relative path
	rel := F.Bind1st(relE, src)

	fixHeader := func(file string) func(*tar.Header) E.Either[error, *tar.Header] {
		return func(hdr *tar.Header) E.Either[error, *tar.Header] {
			return F.Pipe3(
				file,
				rel,
				E.Map[error](filepath.ToSlash),
				E.Map[error](func(relName string) *tar.Header {
					hdr.Name = relName
					return hdr
				}),
			)
		}
	}

	return func(w *tar.Writer) func(file string, fi os.FileInfo) E.Either[error, *tar.Writer] {

		// callback to write the header
		writeHeader := func(hdr *tar.Header) E.Either[error, *tar.Writer] {
			return E.TryCatchError(func() (*tar.Writer, error) {
				return w, w.WriteHeader(hdr)
			})
		}

		return func(file string, fi os.FileInfo) E.Either[error, *tar.Writer] {
			return F.Pipe2(
				fileInfoHeaderE(fi, file),
				E.Chain(fixHeader(file)),
				E.Chain(writeHeader),
			)
		}
	}
}

type tarStreams = T.Tuple2[*gzip.Writer, *tar.Writer]

var (
	gzipStream = T.FirstOf2[*gzip.Writer, *tar.Writer]
	tarStream  = T.SecondOf2[*gzip.Writer, *tar.Writer]
)

func onCreateStreams(buf io.Writer) func() E.Either[error, tarStreams] {
	return func() E.Either[error, tarStreams] {
		gz := gzip.NewWriter(buf)
		tw := tar.NewWriter(gz)
		return E.Of[error](T.MakeTuple2(gz, tw))
	}
}

func onCloseStreams(streams tarStreams) E.Either[error, any] {
	tar := F.Pipe2(
		streams,
		tarStream,
		onClose[*tar.Writer],
	)
	gz := F.Pipe2(
		streams,
		gzipStream,
		onClose[*gzip.Writer],
	)
	return F.Pipe1(
		E.SequenceArray[error, any]()([]E.Either[error, any]{tar, gz}),
		E.MapTo[error, []any, any](nil),
	)
}

func TarFolder[W io.Writer](src string) func(W) E.Either[error, W] {

	writeRel := writeHeader(src)
	walk := F.Bind1st(filepath.Walk, src)

	return func(buf W) E.Either[error, W] {

		return E.WithResource[error, tarStreams, W](
			onCreateStreams(buf),
			onCloseStreams,
		)(func(streams tarStreams) E.Either[error, W] {
			// prepare some context
			tw := tarStream(streams)
			copy := copyFile(tw)
			header := writeRel(tw)

			walkFunc := func(file string, fi os.FileInfo, e error) error {
				// header
				return F.Pipe2(
					header(file, fi),
					E.Chain(func(_ *tar.Writer) E.Either[error, int64] {
						return copy(file, fi)
					}),
					E.ToError[int64],
				)
			}

			return E.TryCatchError(func() (W, error) {
				// walk through every file in the folder
				return buf, walk(walkFunc)

			})
		})
	}
}
