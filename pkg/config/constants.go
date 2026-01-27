package config

const (
	// Default Configuration Values
	DefaultInputFile   = "protec.txt"
	DefaultOutputFile  = "found_secrets.txt"
	DefaultThreadCount = 100
	DefaultTimeoutSec  = 6
	DefaultIdleTimeout = 30
	DefaultDelayMs     = 0

	
	// Application Constants
	ProtectedPrefix       = "Protected S3 Bucket: "
	MinBucketNameLen      = 3
	ChannelBufferMulti    = 10
	UserAgent             = "S3Hunter/2.0"
	S3UrlRegex            = `http[s]?://([a-zA-Z0-9.-]+)\.s3\.amazonaws\.com`
	
	// WAF Evasion Constants
	JitterPercentage = 0.1
	JitterMultiplier = 2
)
