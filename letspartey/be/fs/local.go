package fs

import (
	"errors"
	"os"
	gopath "path"
	"sync"
	"time"

	"iamlinix.com/partay/logger"
)

type CacheFile struct {
	Inst      *File
	Data      []byte
	Timestamp int64
}

type LocalFS struct {
	Cache         map[string]*CacheFile
	CachedSize    int64
	LastClearTime int64
	ClearInterval int64
	CacheMaxLife  int64
	MaxCacheSize  int64
	Mux           sync.Mutex
}

func FileSanityCheck(file *File) error {
	if file == nil {
		logger.Errorf("invalid file")
		return errors.New("Invalid file")
	}

	switch file.Fd.(type) {
	case *os.File:
		break
	default:
		logger.Errorf("incorrect file type")
		return errors.New("Incorrect file type")
	}

	return nil
}

func (self *LocalFS) ClearCache() {
	self.Mux.Lock()
	defer self.Mux.Unlock()

	now := time.Now().Unix()
	paths := make([]string, 0)
	for p, f := range self.Cache {
		if now-f.Timestamp >= self.CacheMaxLife {
			logger.Warnf("to delete cached file [%s]: %d", p, f.Timestamp)
			paths = append(paths, p)
		}
	}

	for _, p := range paths {
		delete(self.Cache, p)
	}
}

func (self *LocalFS) CacheFile(file *File, data []byte) {
	var cache *CacheFile
	var ok bool
	size := int64(len(data))
	if cache, ok = self.Cache[file.Path]; ok {
		if cache.Data != nil {
			logger.Warnf("file already cached and has valid data: %s", file.Path)
			return
		} else {
			logger.Warnf("file already cached but no data attached: %s", file.Path)
		}
	}

	if self.CachedSize+size < self.MaxCacheSize {
		if cache == nil {
			cache = &CacheFile{
				Inst:      file,
				Timestamp: time.Now().Unix(),
			}
			self.Cache[file.Path] = cache
		}
		cache.Data = data
		self.CachedSize += size
		logger.Infof("new file cached: %s", file.Path)
	} else {
		logger.Errorf("cannot cache any more files, quota is full: %d / %d", self.CachedSize, self.MaxCacheSize)
	}
}

func (self *LocalFS) RemoveCache(path string) {
	if _, ok := self.Cache[path]; ok {
		delete(self.Cache, path)
	}
}

func (self *LocalFS) Init(args map[string]interface{}) error {
	logger.Info("init local fs:", args)
	if interval, ok := args["clear_interval"]; ok {
		switch interval.(type) {
		case int:
			self.ClearInterval = int64(interval.(int)) * 1000
			break

		case int64:
			self.ClearInterval = interval.(int64) * 1000
			break

		default:
			return errors.New("Unsupported clear interval value")
		}
	}

	if maxCache, ok := args["max_cache"]; ok {
		switch maxCache.(type) {
		case int:
			self.MaxCacheSize = int64(maxCache.(int)) * 1024 * 1024
			break

		case int64:
			self.MaxCacheSize = maxCache.(int64) * 1024 * 1024
			break

		default:
			return errors.New("Unsupported clear max cache value")
		}
	}

	if cacheLife, ok := args["cache_life"]; ok {
		switch cacheLife.(type) {
		case int:
			self.CacheMaxLife = int64(cacheLife.(int)) * 1000
			break

		case int64:
			self.CacheMaxLife = cacheLife.(int64) * 1000
			break

		default:
			return errors.New("Unsupported cache lift value")
		}
	}

	return nil
}

func (self *LocalFS) Cleanup() {

}

func (self *LocalFS) Create(path string) (*File, error) {
	realPath := gopath.Join(_baseDir, path)
	fd, err := os.Create(realPath)
	if err != nil {
		logger.Errorf("error creating file %s: %v", path, err)
		return nil, err
	}

	file := &File{
		Path:     realPath,
		Fd:       fd,
		Position: 0,
	}
	return file, nil
}

func (self *LocalFS) Mkdir(path string) error {
	return os.MkdirAll(gopath.Join(_baseDir, path), os.ModePerm)
}

func (self *LocalFS) Open(path string) (*File, error) {
	now := time.Now().Unix()
	realPath := gopath.Join(_baseDir, path)
	if cache, ok := self.Cache[realPath]; ok {
		logger.Infof("local fs cache hit: %s", realPath)
		cache.Timestamp = now
		return cache.Inst, nil
	}

	if now-self.LastClearTime > self.ClearInterval {
		self.LastClearTime = now
		go self.ClearCache()
	}

	fd, err := os.Open(realPath)
	if err != nil {
		logger.Errorf("failed to open file %s: %v", realPath, err)
		return nil, err
	}

	file := &File{
		Path:     realPath,
		Fd:       fd,
		Position: 0,
	}
	return file, nil
}

func (self *LocalFS) Read(file *File, start, length int64) ([]byte, error) {
	if length <= 0 || start < 0 {
		logger.Errorf("invalid read length: %d", length)
		return nil, errors.New("Invalid read length")
	}

	if c, ok := self.Cache[file.Path]; ok {
		if c.Data != nil {
			logger.Infof("file read all cache hit: %s", file.Path)
			if int64(len(c.Data)) >= start+length {
				return c.Data[start : start+length], nil
			}
		}
	}

	err := FileSanityCheck(file)
	if err != nil {
		return nil, err
	}

	buffer := make([]byte, length)
	fd := file.Fd.(*os.File)
	n, err := fd.ReadAt(buffer, start)
	if err != nil {
		return nil, err
	}

	if int64(n) != length {
		logger.Warnf("read size mismatch, requested: %d, actually read: %d", length, n)
		newBuffer := make([]byte, n)
		for i := 0; i < n; i++ {
			newBuffer[i] = buffer[i]
		}
		buffer = newBuffer
	}

	file.Position += int64(n)

	return buffer, nil
}

func (self *LocalFS) ReadAll(file *File) ([]byte, error) {
	if c, ok := self.Cache[file.Path]; ok {
		if c.Data != nil {
			logger.Infof("file read all cache hit: %s", file.Path)
			return c.Data, nil
		}
	}

	err := FileSanityCheck(file)
	if err != nil {
		return nil, err
	}

	fd := file.Fd.(*os.File)
	stat, err := fd.Stat()
	if err != nil {
		logger.Errorf("error stating file %s: %v", file.Path, err)
		return nil, err
	}

	size := stat.Size()

	// do not use ReadAll if file size is larger than 1GB
	if size >= 1073741824 {
		logger.Errorf("file too big to ReadAll")
		return nil, errors.New("File too big")
	}

	buffer := make([]byte, size)
	n, err := fd.Read(buffer)
	if err != nil {
		logger.Errorf("error read file %s: %v", file.Path, err)
		return nil, err
	}

	if int64(n) != size {
		logger.Warnf("file read size mismatch, stat size: %d, read size: %d", size, n)
	}
	file.Position = size

	go self.CacheFile(file, buffer)

	return buffer, nil
}

func (*LocalFS) Stat(file *File) (*FileStat, error) {
	err := FileSanityCheck(file)
	if err != nil {
		return nil, err
	}

	fd := file.Fd.(*os.File)
	stat, err := fd.Stat()
	if err != nil {
		logger.Errorf("error stating file %s: %v", file.Path, err)
		return nil, err
	}

	fileType := REG_FILE
	if stat.IsDir() {
		fileType = DIR
	}

	return &FileStat{
		Path:  file.Path,
		Type:  fileType,
		Size:  stat.Size(),
		MTime: stat.ModTime(),
	}, nil
}

func (*LocalFS) Write(file *File, content []byte) error {
	err := FileSanityCheck(file)
	if err != nil {
		return err
	}

	fd := file.Fd.(*os.File)
	n, err := fd.Write(content)
	if err != nil {
		logger.Errorf("error writing file %s: %v", file.Path, err)
		return err
	}

	contentLen := len(content)
	if n != contentLen {
		logger.Warnf("file write size mismatch, content size: %d, written size: %d", contentLen, n)
	}

	file.Position += int64(n)

	return nil
}

func (*LocalFS) Flush(file *File) error {
	err := FileSanityCheck(file)
	if err != nil {
		return err
	}

	fd := file.Fd.(*os.File)
	err = fd.Sync()
	if err != nil {
		logger.Errorf("error flushing file %s: %v", file.Path, err)
	}

	return err
}

func (*LocalFS) Seek(file *File, offset int64, whence int) error {
	err := FileSanityCheck(file)
	if err != nil {
		return err
	}

	fd := file.Fd.(*os.File)
	ret, err := fd.Seek(offset, whence)
	if err != nil {
		logger.Errorf("error seeking file %s: %v", file.Path, err)
		return err
	}

	file.Position = ret
	return nil
}

func (*LocalFS) CloseFile(file *File) error {
	if file.Fd == nil {
		logger.Warnf("file fd is nil, may be cached: %s", file.Path)
		return nil
	}

	err := FileSanityCheck(file)
	if err != nil {
		return err
	}

	fd := file.Fd.(*os.File)
	err = fd.Close()
	if err != nil {
		logger.Errorf("error closing file %s: %v", file.Path, err)
		return err
	}
	file.Fd = nil

	return nil
}

func (self *LocalFS) Delete(file *File) error {
	err := FileSanityCheck(file)
	if err != nil {
		return err
	}

	go self.RemoveCache(file.Path)

	fd := file.Fd.(*os.File)
	stat, err := fd.Stat()
	if err != nil {
		logger.Errorf("error stating while deleting file %s: %v", file.Path, err)
		return err
	}

	if stat.IsDir() {
		return self.DeleteDir(file.Path)
	} else {
		return self.DeleteFile(file.Path)
	}
}

func (self *LocalFS) DeleteFile(path string) error {
	go self.RemoveCache(path)

	err := os.Remove(gopath.Join(_baseDir, path))
	if err != nil {
		logger.Errorf("error while deleting file: %v", err)
		return err
	}

	return nil
}

func (*LocalFS) DeleteDir(path string) error {
	err := os.RemoveAll(gopath.Join(_baseDir, path))
	if err != nil {
		logger.Errorf("error while deleting dir: %v", err)
		return err
	}

	return nil
}
