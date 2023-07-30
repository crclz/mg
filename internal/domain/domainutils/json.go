package domainutils

import "encoding/json"

func ToJson(object any) string {
	var bytes, err = json.Marshal(object)
	if err != nil {
		panic(err)
	}

	return string(bytes)
}
