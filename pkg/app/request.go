package app

import (
	"github.com/astaxie/beego/validation"

	"easycache/pkg/logger"
)

// MarkErrors logs error logs
func MarkErrors(errors []*validation.Error) {
	for _, err := range errors {
		logger.Info(err.Key, err.Message)
	}

	return
}
