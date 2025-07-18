package utils

import (
	"context"
)

type UIDError struct {
}

func (err UIDError) Error() string {
	return "Couldn't retrieve UID"
}

func GetUID(ctx context.Context) (string, error) {
	uid, ok := ctx.Value("uid").(string)
	if !ok {
		return "", UIDError{}
	}
	return uid, nil
}
