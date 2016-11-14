package middleware

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/docker/distribution/context"
	"github.com/docker/distribution/registry/storage/cdn"
	storagedriver "github.com/docker/distribution/registry/storage/driver"
	storagemiddleware "github.com/docker/distribution/registry/storage/driver/middleware"
)

const EXPIRE_TIME = 20 * 60

type aliCDNStorageMiddleware struct {
	storagedriver.StorageDriver
	baseURL string
	key     string
}

func newAliCDNStorageMiddleware(storageDriver storagedriver.StorageDriver, options map[string]interface{}) (storagedriver.StorageDriver, error) {
	base, ok := options["baseurl"]
	if !ok {
		return nil, fmt.Errorf("no baseurl provided")
	}
	baseURL, ok := base.(string)
	if !ok {
		return nil, fmt.Errorf("baseurl must be a string")
	}

	baseURL = strings.TrimRight(baseURL, "/")

	argKey, ok := options["key"]
	if !ok {
		return nil, fmt.Errorf("no key provided")
	}
	key, ok := argKey.(string)
	if !ok {
		return nil, fmt.Errorf("key must be a string")
	}

	return &aliCDNStorageMiddleware{
		StorageDriver: storageDriver,
		baseURL:       baseURL,
		key:           key,
	}, nil
}

func md5sum(data []byte) string {
	digest := md5.Sum(data)
	return hex.EncodeToString(digest[:])
}

func (lh *aliCDNStorageMiddleware) sign(path string) string {
	now := time.Now().Unix()
	expire := time.Unix(now+EXPIRE_TIME, 0)
	formatedExpire := expire.Format("200601021504")
	s := fmt.Sprintf("%s%s%s", lh.key, formatedExpire, path)
	hashValue := md5sum([]byte(s))
	return fmt.Sprintf("/%s/%s%s", formatedExpire, hashValue, path)
}

func (lh *aliCDNStorageMiddleware) URLFor(ctx context.Context, path string, options map[string]interface{}) (string, error) {
	if lh.StorageDriver.Name() != "oss" {
		context.GetLogger(ctx).Warn("the alicdn middleware does not support this backend storage driver")
		return lh.StorageDriver.URLFor(ctx, path, options)
	}
	repo := ctx.Value("vars.name")

	if repo != nil && cdn.UseCDN(repo.(string)) {
		return lh.baseURL + lh.sign(path), nil
	}
	return lh.StorageDriver.URLFor(ctx, path, options)
}
func init() {
	storagemiddleware.Register("alicdn", storagemiddleware.InitFunc(newAliCDNStorageMiddleware))
}
