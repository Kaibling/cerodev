package service

import "fmt"

func HandleError[T any](d any, err error, errMsg string) (T, error) { //nolint:ireturn
	if err != nil {
		return d.(T), fmt.Errorf(errMsg+": %w", err) //nolint:err113,forcetypeassert
	}

	return d.(T), nil //nolint:forcetypeassert
}
