package lib

import "encoding/json"

func Copy(dest any, src any) {
	data, _ := json.Marshal(src)
	_ = json.Unmarshal(data, dest)
}
