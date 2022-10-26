/*
 *     Copyright 2022 The Dragonfly Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package storage

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/stretchr/testify/assert"
)

func TestStorage_New(t *testing.T) {
	tests := []struct {
		name    string
		baseDir string
		options []Option
		expect  func(t *testing.T, s Storage, err error)
	}{
		{
			name:    "new storage",
			baseDir: os.TempDir(),
			options: []Option{},
			expect: func(t *testing.T, s Storage, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(reflect.TypeOf(s).Elem().Name(), "storage")
				assert.Equal(s.(*storage).maxSize, int64(DefaultMaxSize*megabyte))
				assert.Equal(s.(*storage).maxBackups, DefaultMaxBackups)
				assert.Equal(s.(*storage).bufferSize, DefaultBufferSize)
				assert.Equal(cap(s.(*storage).buffer), DefaultBufferSize)
				assert.Equal(len(s.(*storage).buffer), 0)
				assert.Equal(s.(*storage).count, int64(0))

				if err := s.Clear(); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name:    "new storage with maxSize",
			baseDir: os.TempDir(),
			options: []Option{WithMaxSize(1)},
			expect: func(t *testing.T, s Storage, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(reflect.TypeOf(s).Elem().Name(), "storage")
				assert.Equal(s.(*storage).maxSize, int64(1*megabyte))
				assert.Equal(s.(*storage).maxBackups, DefaultMaxBackups)
				assert.Equal(s.(*storage).bufferSize, DefaultBufferSize)
				assert.Equal(cap(s.(*storage).buffer), DefaultBufferSize)
				assert.Equal(len(s.(*storage).buffer), 0)

				if err := s.Clear(); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name:    "new storage with maxBackups",
			baseDir: os.TempDir(),
			options: []Option{WithMaxBackups(1)},
			expect: func(t *testing.T, s Storage, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(reflect.TypeOf(s).Elem().Name(), "storage")
				assert.Equal(s.(*storage).maxSize, int64(DefaultMaxSize*megabyte))
				assert.Equal(s.(*storage).maxBackups, 1)
				assert.Equal(s.(*storage).bufferSize, DefaultBufferSize)
				assert.Equal(cap(s.(*storage).buffer), DefaultBufferSize)
				assert.Equal(len(s.(*storage).buffer), 0)

				if err := s.Clear(); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name:    "new storage with bufferSize",
			baseDir: os.TempDir(),
			options: []Option{WithBufferSize(1)},
			expect: func(t *testing.T, s Storage, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(reflect.TypeOf(s).Elem().Name(), "storage")
				assert.Equal(s.(*storage).maxSize, int64(DefaultMaxSize*megabyte))
				assert.Equal(s.(*storage).maxBackups, DefaultMaxBackups)
				assert.Equal(s.(*storage).bufferSize, 1)
				assert.Equal(cap(s.(*storage).buffer), 1)
				assert.Equal(len(s.(*storage).buffer), 0)
				assert.Equal(len(s.(*storage).buffer), 0)

				if err := s.Clear(); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name:    "new storage failed",
			baseDir: "/foo",
			options: []Option{WithMaxBackups(100)},
			expect: func(t *testing.T, s Storage, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s, err := New(tc.baseDir, tc.options...)
			tc.expect(t, s, err)
		})
	}
}

func TestStorage_Create(t *testing.T) {
	tests := []struct {
		name    string
		baseDir string
		options []Option
		mock    func(s Storage)
		expect  func(t *testing.T, s Storage, baseDir string)
	}{
		{
			name:    "create record",
			baseDir: os.TempDir(),
			options: []Option{WithBufferSize(1)},
			mock:    func(s Storage) {},
			expect: func(t *testing.T, s Storage, baseDir string) {
				assert := assert.New(t)
				err := s.Create(Record{})
				assert.NoError(err)
				assert.Equal(s.(*storage).count, int64(0))
			},
		},
		{
			name:    "create record without buffer",
			baseDir: os.TempDir(),
			options: []Option{WithBufferSize(0)},
			mock: func(s Storage) {
			},
			expect: func(t *testing.T, s Storage, baseDir string) {
				assert := assert.New(t)
				err := s.Create(Record{})
				assert.NoError(err)
				assert.Equal(s.(*storage).count, int64(1))
			},
		},
		{
			name:    "write record to file",
			baseDir: os.TempDir(),
			options: []Option{WithBufferSize(1)},
			mock: func(s Storage) {
			},
			expect: func(t *testing.T, s Storage, baseDir string) {
				assert := assert.New(t)
				err := s.Create(Record{})
				assert.NoError(err)
				err = s.Create(Record{})
				assert.NoError(err)
				assert.Equal(s.(*storage).count, int64(1))
			},
		},
		{
			name:    "open file failed",
			baseDir: os.TempDir(),
			options: []Option{WithBufferSize(0)},
			mock: func(s Storage) {
				s.(*storage).baseDir = "foo"
			},
			expect: func(t *testing.T, s Storage, baseDir string) {
				assert := assert.New(t)
				err := s.Create(Record{})
				assert.Error(err)
				s.(*storage).baseDir = baseDir
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s, err := New(tc.baseDir, tc.options...)
			if err != nil {
				t.Fatal(err)
			}

			tc.mock(s)
			tc.expect(t, s, tc.baseDir)
			if err := s.Clear(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestStorage_List(t *testing.T) {
	tests := []struct {
		name    string
		baseDir string
		options []Option
		record  Record
		mock    func(t *testing.T, s Storage, baseDir string, record Record)
		expect  func(t *testing.T, s Storage, baseDir string, record Record)
	}{
		{
			name:    "empty csv file given",
			baseDir: os.TempDir(),
			options: []Option{},
			mock:    func(t *testing.T, s Storage, baseDir string, record Record) {},
			expect: func(t *testing.T, s Storage, baseDir string, record Record) {
				assert := assert.New(t)
				_, err := s.List()
				assert.Error(err)
			},
		},
		{
			name:    "get file infos failed",
			baseDir: os.TempDir(),
			options: []Option{},
			mock: func(t *testing.T, s Storage, baseDir string, record Record) {
				s.(*storage).baseDir = "bae"
			},
			expect: func(t *testing.T, s Storage, baseDir string, record Record) {
				assert := assert.New(t)
				_, err := s.List()
				assert.Error(err)
				s.(*storage).baseDir = baseDir
			},
		},
		{
			name:    "open file failed",
			baseDir: os.TempDir(),
			options: []Option{},
			mock: func(t *testing.T, s Storage, baseDir string, record Record) {
				file, err := os.OpenFile(filepath.Join(baseDir, "record-test.csv"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0300)
				if err != nil {
					t.Fatal(err)
				}
				file.Close()
			},
			expect: func(t *testing.T, s Storage, baseDir string, record Record) {
				assert := assert.New(t)
				_, err := s.List()
				assert.Error(err)
			},
		},
		{
			name:    "list records of a file",
			baseDir: os.TempDir(),
			options: []Option{WithBufferSize(1)},
			record: Record{
				IP:             rand.Intn(100)%2 + 1,
				HostName:       rand.Intn(100)%2 + 1,
				Tag:            rand.Intn(100)%2 + 1,
				Rate:           float64(rand.Intn(300) + 10),
				ParentPiece:    float64(rand.Intn(240) + 14),
				SecurityDomain: rand.Intn(100)%2 + 1,
				IDC:            rand.Intn(100)%2 + 1,
				NetTopology:    rand.Intn(100)%2 + 1,
				Location:       rand.Intn(100)%2 + 1,
				UploadRate:     float64(rand.Intn(550) + 3),
				CreateAt:       time.Now().Unix()/7200 + rand.Int63n(10),
				UpdateAt:       time.Now().Unix()/7200 + rand.Int63n(10),
				ParentCreateAt: time.Now().Unix()/7200 + rand.Int63n(10),
				ParentUpdateAt: time.Now().Unix()/7200 + rand.Int63n(10),
			},
			mock: func(t *testing.T, s Storage, baseDir string, record Record) {
				if err := s.Create(record); err != nil {
					t.Fatal(err)
				}
			},
			expect: func(t *testing.T, s Storage, baseDir string, record Record) {
				assert := assert.New(t)
				_, err := s.List()
				assert.Error(err)

				if err := s.Create(record); err != nil {
					t.Fatal(err)
				}
				records, err := s.List()
				assert.NoError(err)
				assert.Equal(len(records), 1)
				assert.EqualValues(records[0], record)
			},
		},
		{
			name:    "list records of multi files",
			baseDir: os.TempDir(),
			options: []Option{WithBufferSize(1)},
			record:  Record{},
			mock: func(t *testing.T, s Storage, baseDir string, record Record) {
				file, err := os.OpenFile(filepath.Join(baseDir, "record-test.csv"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
				if err != nil {
					t.Fatal(err)
				}
				defer file.Close()

				if err := gocsv.MarshalWithoutHeaders([]Record{{Tag: 2}}, file); err != nil {
					t.Fatal(err)
				}

				if err := s.Create(Record{Tag: 1}); err != nil {
					t.Fatal(err)
				}

				if err := s.Create(Record{Tag: 3}); err != nil {
					t.Fatal(err)
				}
			},
			expect: func(t *testing.T, s Storage, baseDir string, record Record) {
				assert := assert.New(t)
				records, err := s.List()
				assert.NoError(err)
				assert.Equal(len(records), 2)
				assert.Equal(records[0].Tag, 2)
				assert.Equal(records[1].Tag, 1)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s, err := New(tc.baseDir, tc.options...)
			if err != nil {
				t.Fatal(err)
			}

			tc.mock(t, s, tc.baseDir, tc.record)
			tc.expect(t, s, tc.baseDir, tc.record)
			if err := s.Clear(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestStorage_Open(t *testing.T) {
	tests := []struct {
		name    string
		baseDir string
		options []Option
		record  Record
		mock    func(t *testing.T, s Storage, baseDir string, record Record)
		expect  func(t *testing.T, s Storage, baseDir string, record Record)
	}{
		{
			name:    "open storage withempty csv file given",
			baseDir: os.TempDir(),
			options: []Option{},
			mock:    func(t *testing.T, s Storage, baseDir string, record Record) {},
			expect: func(t *testing.T, s Storage, baseDir string, record Record) {
				assert := assert.New(t)
				_, err := s.Open()
				assert.NoError(err)
			},
		},
		{
			name:    "open file infos failed",
			baseDir: os.TempDir(),
			options: []Option{},
			mock: func(t *testing.T, s Storage, baseDir string, record Record) {
				s.(*storage).baseDir = "bas"
			},
			expect: func(t *testing.T, s Storage, baseDir string, record Record) {
				assert := assert.New(t)
				_, err := s.Open()
				assert.Error(err)
				s.(*storage).baseDir = baseDir
			},
		},
		{
			name:    "open storage with records of a file",
			baseDir: os.TempDir(),
			options: []Option{WithBufferSize(1)},
			record: Record{
				IP:             rand.Intn(100)%2 + 1,
				HostName:       rand.Intn(100)%2 + 1,
				Tag:            rand.Intn(100)%2 + 1,
				Rate:           float64(rand.Intn(300) + 10),
				ParentPiece:    float64(rand.Intn(240) + 14),
				SecurityDomain: rand.Intn(100)%2 + 1,
				IDC:            rand.Intn(100)%2 + 1,
				NetTopology:    rand.Intn(100)%2 + 1,
				Location:       rand.Intn(100)%2 + 1,
				UploadRate:     float64(rand.Intn(550) + 3),
				CreateAt:       time.Now().Unix()/7200 + rand.Int63n(10),
				UpdateAt:       time.Now().Unix()/7200 + rand.Int63n(10),
				ParentCreateAt: time.Now().Unix()/7200 + rand.Int63n(10),
				ParentUpdateAt: time.Now().Unix()/7200 + rand.Int63n(10),
			},
			mock: func(t *testing.T, s Storage, baseDir string, record Record) {
				if err := s.Create(record); err != nil {
					t.Fatal(err)
				}
			},
			expect: func(t *testing.T, s Storage, baseDir string, record Record) {
				assert := assert.New(t)
				_, err := s.Open()
				assert.NoError(err)

				if err := s.Create(record); err != nil {
					t.Fatal(err)
				}

				readCloser, err := s.Open()
				assert.NoError(err)

				var records []Record
				err = gocsv.UnmarshalWithoutHeaders(readCloser, &records)
				assert.NoError(err)
				assert.Equal(len(records), 1)
				assert.EqualValues(records[0], record)
			},
		},
		{
			name:    "open storage with records of multi files",
			baseDir: os.TempDir(),
			options: []Option{WithBufferSize(1)},
			record:  Record{},
			mock: func(t *testing.T, s Storage, baseDir string, record Record) {
				file, err := os.OpenFile(filepath.Join(baseDir, "record-test.csv"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
				if err != nil {
					t.Fatal(err)
				}
				defer file.Close()

				if err := gocsv.MarshalWithoutHeaders([]Record{{Tag: 2}}, file); err != nil {
					t.Fatal(err)
				}

				if err := s.Create(Record{Tag: 1}); err != nil {
					t.Fatal(err)
				}

				if err := s.Create(Record{Tag: 3}); err != nil {
					t.Fatal(err)
				}
			},
			expect: func(t *testing.T, s Storage, baseDir string, record Record) {
				assert := assert.New(t)
				readCloser, err := s.Open()
				assert.NoError(err)

				var records []Record
				err = gocsv.UnmarshalWithoutHeaders(readCloser, &records)
				assert.NoError(err)
				assert.Equal(len(records), 2)
				assert.Equal(records[0].Tag, 2)
				assert.Equal(records[1].Tag, 1)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s, err := New(tc.baseDir, tc.options...)
			if err != nil {
				t.Fatal(err)
			}

			tc.mock(t, s, tc.baseDir, tc.record)
			tc.expect(t, s, tc.baseDir, tc.record)
			if err := s.Clear(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestStorage_Clear(t *testing.T) {
	tests := []struct {
		name    string
		baseDir string
		options []Option
		mock    func(s Storage)
		expect  func(t *testing.T, s Storage, baseDir string)
	}{
		{
			name:    "clear file",
			baseDir: os.TempDir(),
			options: []Option{},
			mock:    func(s Storage) {},
			expect: func(t *testing.T, s Storage, baseDir string) {
				assert := assert.New(t)
				assert.NoError(s.Clear())
				fileInfos, err := ioutil.ReadDir(filepath.Join(baseDir))
				assert.NoError(err)

				var backups []fs.FileInfo
				regexp := regexp.MustCompile(RecordFilePrefix)
				for _, fileInfo := range fileInfos {
					if !fileInfo.IsDir() && regexp.MatchString(fileInfo.Name()) {
						backups = append(backups, fileInfo)
					}
				}
				assert.Equal(len(backups), 0)
			},
		},
		{
			name:    "open file failed",
			baseDir: os.TempDir(),
			options: []Option{},
			mock: func(s Storage) {
				s.(*storage).baseDir = "baz"
			},
			expect: func(t *testing.T, s Storage, baseDir string) {
				assert := assert.New(t)
				assert.Error(s.Clear())

				s.(*storage).baseDir = baseDir
				assert.NoError(s.Clear())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s, err := New(tc.baseDir, tc.options...)
			if err != nil {
				t.Fatal(err)
			}

			tc.mock(s)
			tc.expect(t, s, tc.baseDir)
		})
	}
}

func TestStorage_create(t *testing.T) {
	tests := []struct {
		name    string
		baseDir string
		options []Option
		mock    func(s Storage)
		expect  func(t *testing.T, s Storage, baseDir string)
	}{
		{
			name:    "create record",
			baseDir: os.TempDir(),
			options: []Option{},
			mock:    func(s Storage) {},
			expect: func(t *testing.T, s Storage, baseDir string) {
				assert := assert.New(t)
				err := s.(*storage).create(Record{})
				assert.NoError(err)
			},
		},
		{
			name:    "open file failed",
			baseDir: os.TempDir(),
			options: []Option{},
			mock: func(s Storage) {
				s.(*storage).baseDir = "foo"
			},
			expect: func(t *testing.T, s Storage, baseDir string) {
				assert := assert.New(t)
				err := s.(*storage).create(Record{})
				assert.Error(err)
				s.(*storage).baseDir = baseDir
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s, err := New(tc.baseDir, tc.options...)
			if err != nil {
				t.Fatal(err)
			}

			tc.mock(s)
			tc.expect(t, s, tc.baseDir)
			if err := s.Clear(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestStorage_openFile(t *testing.T) {
	tests := []struct {
		name    string
		baseDir string
		options []Option
		mock    func(t *testing.T, s Storage)
		expect  func(t *testing.T, s Storage, baseDir string)
	}{
		{
			name:    "open file failed",
			baseDir: os.TempDir(),
			options: []Option{},
			mock: func(t *testing.T, s Storage) {
				s.(*storage).baseDir = "bat"
			},
			expect: func(t *testing.T, s Storage, baseDir string) {
				assert := assert.New(t)
				_, err := s.(*storage).openFile()
				assert.Error(err)
				s.(*storage).baseDir = baseDir
			},
		},
		{
			name:    "open new record file",
			baseDir: os.TempDir(),
			options: []Option{WithMaxSize(0), WithBufferSize(1)},
			mock: func(t *testing.T, s Storage) {
				if err := s.Create(Record{Tag: 1}); err != nil {
					t.Fatal(err)
				}

				if err := s.Create(Record{Tag: 2}); err != nil {
					t.Fatal(err)
				}
			},
			expect: func(t *testing.T, s Storage, baseDir string) {
				assert := assert.New(t)
				file, err := s.(*storage).openFile()
				assert.NoError(err)
				assert.Equal(file.Name(), filepath.Join(baseDir, fmt.Sprintf("%s.%s", RecordFilePrefix, RecordFileExt)))
				file.Close()
			},
		},
		{
			name:    "remove record file",
			baseDir: os.TempDir(),
			options: []Option{WithMaxSize(0), WithMaxBackups(1), WithBufferSize(1)},
			mock: func(t *testing.T, s Storage) {
				if err := s.Create(Record{Tag: 1}); err != nil {
					t.Fatal(err)
				}
			},
			expect: func(t *testing.T, s Storage, baseDir string) {
				assert := assert.New(t)
				file, err := s.(*storage).openFile()
				assert.NoError(err)
				assert.Equal(file.Name(), filepath.Join(baseDir, fmt.Sprintf("%s.%s", RecordFilePrefix, RecordFileExt)))
				file.Close()
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s, err := New(tc.baseDir, tc.options...)
			if err != nil {
				t.Fatal(err)
			}

			tc.mock(t, s)
			tc.expect(t, s, tc.baseDir)
			if err := s.Clear(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestStorage_backupFilename(t *testing.T) {
	baseDir := os.TempDir()
	s, err := New(baseDir)
	if err != nil {
		t.Fatal(err)
	}

	filename := s.(*storage).backupFilename()
	regexp := regexp.MustCompile(fmt.Sprintf("%s-.*.%s$", RecordFilePrefix, RecordFileExt))
	assert := assert.New(t)
	assert.True(regexp.MatchString(filename))

	if err := s.Clear(); err != nil {
		t.Fatal(err)
	}
}

func TestStorage_backups(t *testing.T) {
	tests := []struct {
		name    string
		baseDir string
		options []Option
		mock    func(t *testing.T, s Storage)
		expect  func(t *testing.T, s Storage, baseDir string)
	}{
		{
			name:    "open file failed",
			baseDir: os.TempDir(),
			options: []Option{},
			mock: func(t *testing.T, s Storage) {
				s.(*storage).baseDir = "bar"
			},
			expect: func(t *testing.T, s Storage, baseDir string) {
				assert := assert.New(t)
				_, err := s.(*storage).backups()
				assert.Error(err)
				s.(*storage).baseDir = baseDir
				if err := s.Clear(); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name:    "not found record file",
			baseDir: os.TempDir(),
			options: []Option{},
			mock:    func(t *testing.T, s Storage) {},
			expect: func(t *testing.T, s Storage, baseDir string) {
				assert := assert.New(t)
				if err := s.Clear(); err != nil {
					t.Fatal(err)
				}

				_, err := s.(*storage).backups()
				assert.Error(err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s, err := New(tc.baseDir, tc.options...)
			if err != nil {
				t.Fatal(err)
			}

			tc.mock(t, s)
			tc.expect(t, s, tc.baseDir)
		})
	}
}
