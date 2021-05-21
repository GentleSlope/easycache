package util

import (
	"easycache/pkg/setting"
	"fmt"
)

// Setup Initialize the util
func Setup() {
	jwtSecret = []byte(setting.AppSetting.JwtSecret)
}

//Interface2String ...
func Interface2String(inter interface{}) string {

	switch inter.(type) {
	case string:
		return fmt.Sprintf("%s", inter.(string))
	case int:
		return fmt.Sprintf("%d", inter.(int))
	case int64:
		return fmt.Sprintf("%d", inter.(int64))
	default:
		return ""
	}
}
