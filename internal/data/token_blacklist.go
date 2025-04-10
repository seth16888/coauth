package data

import "context"

var memStore = map[string]string{}

type TokenBlacklistData struct {
}

func NewTokenBlacklistData() *TokenBlacklistData {
	return &TokenBlacklistData{}
}

func (t *TokenBlacklistData) Push(ctx context.Context, tokenId string) error {
	key := "akbk:" + tokenId
	memStore[key] = tokenId
	return nil
}

func (t *TokenBlacklistData) Exists(ctx context.Context, tokenId string) (bool, error) {
	key := "akbk:" + tokenId
	_, ok := memStore[key]
	return ok, nil
}
