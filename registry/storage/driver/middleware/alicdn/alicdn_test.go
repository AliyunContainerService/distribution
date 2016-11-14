package middleware

import (
	"fmt"
	"testing"
)

func TestSign(t *testing.T) {
	aliCDN := &aliCDNStorageMiddleware{
		key: "helloworld",
	}
	fmt.Println(aliCDN.sign("/hello"))
}
