package db

import "tagstore/middleware"

// i cheated

var DB db

type db struct {
	ManifestKV
	TagKV
}

func NewDB() middleware.Store {
	return db{
		make(ManifestKV),
		make(TagKV),
	}
}
