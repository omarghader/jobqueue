package shared

// Flags contains values of flags which are important in the whole API
type Flags struct {
	RavenDSN string

	BindAddress string
	BaseUrl     string

	APIVersion       string
	LogFormatterType string
	ForceColors      bool

	DatabaseDriver    string
	DatabaseHost      string
	DatabaseNamespace string
	DatabaseUser      string
	DatabasePassword  string

	AmqpURI        string
	WebhookURLs    string
	MemcachedHosts string
	RedisHost      string
}
