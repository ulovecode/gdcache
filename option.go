package gdcache

import "github.com/ulovecode/gdcache/log"

type Options struct {
	cacheTagName string
	log          log.Logger
	serializer   Serializer
}

type OptionsFunc func(o *Options)

func WithCacheTagName(cacheTagName string) OptionsFunc {
	return func(o *Options) {
		o.cacheTagName = cacheTagName
	}
}

func WithLogger(logger log.Logger) OptionsFunc {
	return func(o *Options) {
		o.log = logger
	}
}

func WithSerializer(serializer Serializer) OptionsFunc {
	return func(o *Options) {
		o.serializer = serializer
	}
}
