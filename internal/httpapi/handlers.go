package httpapi

import (
	"sendtokindle/internal/storage"
	"sendtokindle/internal/web"
)

type Handlers struct {
	Store    *storage.Store
	Renderer *web.Renderer
	KindleURL string
	StoreRoot string
}
