package gdcache

// Options Optional
type Options struct {
	cacheTagName string
	log          Logger
	serializer   Serializer
	serviceName  string
}

// OptionsFunc Alternative method
type OptionsFunc func(o *Options)

// WithServiceName Service Name
func WithServiceName(serviceName string) OptionsFunc {
	return func(o *Options) {
		o.serviceName = serviceName
	}
}

// WithCacheTagName Specify cache tag name
func WithCacheTagName(cacheTagName string) OptionsFunc {
	return func(o *Options) {
		o.cacheTagName = cacheTagName
	}
}

// WithLogger Specify log
func WithLogger(logger Logger) OptionsFunc {
	return func(o *Options) {
		o.log = logger
	}
}

// WithSerializer Specify the serialization method
func WithSerializer(serializer Serializer) OptionsFunc {
	return func(o *Options) {
		o.serializer = serializer
	}
}
