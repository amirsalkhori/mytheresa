package derrors

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/pkg/errors"
)

type DomainError interface {
	error
	Code() ErrorCode
	Checksums() string
}

func NewDomainErrorWithMessage(code ErrorCode, msg string) DomainError {
	return baseError{
		code:    code,
		message: msg,
	}
}

func NewDomainError(code ErrorCode) DomainError {
	return baseError{
		code: code,
	}
}

type ErrorCode int32

type baseError struct {
	code    ErrorCode
	message string
}

func (b baseError) Error() string {
	return fmt.Sprintf("appError(%d)", b.code)
}

func (b baseError) Code() ErrorCode {
	return b.code
}

func (b baseError) Checksums() string {
	hash := md5.Sum([]byte(fmt.Sprintf("%d-%s", b.code, b.message)))
	return hex.EncodeToString(hash[:])
}

func Cause(err error) error {
	root := err

	var domainError DomainError
	for {
		if errors.As(root, &domainError) {
			return domainError
		}
		if err = errors.Unwrap(root); err == nil {
			return root
		}

		root = err
	}
}
