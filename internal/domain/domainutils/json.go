package domainutils

import (
	"github.com/bytedance/sonic"
)

func ToJson(object any) string {
	var bytes, err = sonic.Marshal(object)
	if err != nil {
		panic(err)
	}

	return string(bytes)
}
