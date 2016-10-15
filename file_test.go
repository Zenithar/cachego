package cachego

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type FileTestSuite struct {
	suite.Suite

	assert *assert.Assertions
	cache  Cache
}

func (s *FileTestSuite) SetupTest() {
	directory := "./cache-dir/"

	os.Mkdir(directory, 0777)

	s.cache = &File{
		dir: directory,
	}

	s.assert = assert.New(s.T())
}

func (s *FileTestSuite) TestSaveReturnFalseWhenThrowError() {
	cache := &File{}

	s.assert.False(cache.Save("foo", "bar", 0))
}

func (s *FileTestSuite) TestSave() {
	s.assert.True(s.cache.Save("", "bar", 0))
}

func (s *FileTestSuite) TestSaveWithCustomExtension() {
	directory := "./cache-dir/"

	os.Mkdir(directory, 0777)

	cache := &File{
		dir:       directory,
		extension: "data.cache",
	}

	s.assert.True(cache.Save("foo", "bar", 0))
}

func (s *FileTestSuite) TestFetchReturnFalseWhenThrowError() {
	key := "foo"
	value := "bar"

	s.cache.Save(key, value, 0)

	cache := &File{}

	result, status := cache.Fetch(key)

	s.assert.False(status)
	s.assert.Empty(result)
}

func (s *FileTestSuite) TestFetch() {
	key := "foo"
	value := "bar"

	s.cache.Save(key, value, 0)

	result, status := s.cache.Fetch(key)

	s.assert.True(status)
	s.assert.Equal(value, result)
}

func (s *FileTestSuite) TestContainsReturnFalseWhenThrowError() {
	cache := &File{}

	s.assert.False(cache.Contains("bar"))
}

func (s *FileTestSuite) TestContains() {
	s.cache.Save("foo", "bar", 0)

	s.assert.True(s.cache.Contains("foo"))
	s.assert.False(s.cache.Contains("bar"))
}

func (s *FileTestSuite) TestDeleteReturnFalseWhenThrowError() {
	cache := &File{}
	s.assert.False(cache.Delete("bar"))
}

func (s *FileTestSuite) TestDelete() {
	s.cache.Save("foo", "bar", 0)

	s.assert.True(s.cache.Delete("foo"))
	s.assert.False(s.cache.Contains("foo"))
	s.assert.False(s.cache.Delete("foo"))
}

func (s *FileTestSuite) TestFlushReturnFalseWhenThrowError() {
	cache := &File{}

	s.assert.False(cache.Flush())
}

func (s *FileTestSuite) TestFlush() {
	s.cache.Save("foo", "bar", 0)

	s.assert.True(s.cache.Flush())
	s.assert.False(s.cache.Contains("foo"))
}

func (s *FileTestSuite) TestFetchMultiReturnNoItemsWhenThrowError() {
	cache := &File{}
	result := cache.FetchMulti([]string{"foo"})

	s.assert.Len(result, 0)
}

func (s *FileTestSuite) TestFetchMulti() {
	s.cache.Save("foo", "bar", 0)
	s.cache.Save("john", "doe", 0)

	result := s.cache.FetchMulti([]string{"foo", "john"})

	s.assert.Len(result, 2)
}

func (s *FileTestSuite) TestFetchMultiWhenOnlyOneOfKeysExists() {
	s.cache.Save("foo", "bar", 0)

	result := s.cache.FetchMulti([]string{"foo", "alice"})

	s.assert.Len(result, 1)
}

func TestFileRunSuite(t *testing.T) {
	suite.Run(t, new(FileTestSuite))
}
