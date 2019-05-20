package splunknozzle

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cloudfoundry-community/splunk-firehose-nozzle/events"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type Config struct {
	ApiEndpoint string `json:"api-endpoint"`
	User        string `json:"-"`
	Password    string `json:"-"`

	SplunkToken string `json:"-"`
	SplunkHost  string `json:"splunk-host"`
	SplunkIndex string `json:"splunk-index"`

	JobName  string `json:"job-name"`
	JobIndex string `json:"job-index"`
	JobHost  string `json:"job-host"`

	SkipSSLCF      bool          `json:"skip-ssl-cf"`
	SkipSSLSplunk  bool          `json:"skip-ssl-splunk"`
	SubscriptionID string        `json:"subscription-id"`
	KeepAlive      time.Duration `json:"keep-alive"`

	AddAppInfo         bool          `json:"add-app-info"`
	IgnoreMissingApps  bool          `json:"ignore-missing-apps"`
	MissingAppCacheTTL time.Duration `json:"missing-app-cache-ttl"`
	AppCacheTTL        time.Duration `json:"app-cache-ttl"`
	OrgSpaceCacheTTL   time.Duration `json:"org-space-cache-ttl"`
	AppLimits          int           `json:"app-limits"`

	BoltDBPath   string `json:"boltdb-path"`
	WantedEvents string `json:"wanted-events"`
	ExtraFields  string `json:"extra-fields"`

	FlushInterval time.Duration `json:"flush-interval"`
	QueueSize     int           `json:"queue-size"`
	BatchSize     int           `json:"batch-size"`
	Retries       int           `json:"retries"`
	HecWorkers    int           `json:"hec-workers"`
	SplunkVersion string        `json:"splunk-version"`

	Version string `json:"version"`
	Branch  string `json:"branch"`
	Commit  string `json:"commit"`
	BuildOS string `json:"buildos"`

	TraceLogging bool `json:"trace-logging"`
	Debug        bool `json:"debug"`
}

func NewConfigFromCmdFlags(version, branch, commit, buildos string) *Config {
	c := &Config{}

	c.Version = version
	c.Branch = branch
	c.Commit = commit
	c.BuildOS = buildos

	kingpin.Version(version)
	kingpin.Flag("api-endpoint", "API endpoint address").
		OverrideDefaultFromEnvar("API_ENDPOINT").Required().StringVar(&c.ApiEndpoint)
	kingpin.Flag("user", "Admin user.").
		OverrideDefaultFromEnvar("API_USER").Required().StringVar(&c.User)
	kingpin.Flag("password", "Admin password.").
		OverrideDefaultFromEnvar("API_PASSWORD").Required().StringVar(&c.Password)

	kingpin.Flag("splunk-host", "Splunk HTTP event collector host").
		OverrideDefaultFromEnvar("SPLUNK_HOST").Required().StringVar(&c.SplunkHost)
	kingpin.Flag("splunk-token", "Splunk HTTP event collector token").
		OverrideDefaultFromEnvar("SPLUNK_TOKEN").Required().StringVar(&c.SplunkToken)
	kingpin.Flag("splunk-index", "Splunk index").
		OverrideDefaultFromEnvar("SPLUNK_INDEX").Required().StringVar(&c.SplunkIndex)

	kingpin.Flag("job-name", "Job name to tag nozzle's own log events").
		OverrideDefaultFromEnvar("JOB_NAME").Default("splunk-nozzle").StringVar(&c.JobName)
	kingpin.Flag("job-index", "Job index to tag nozzle's own log events").
		OverrideDefaultFromEnvar("JOB_INDEX").Default("-1").StringVar(&c.JobIndex)
	kingpin.Flag("job-host", "Job host to tag nozzle's own log events").
		OverrideDefaultFromEnvar("JOB_HOST").Default("").StringVar(&c.JobHost)

	kingpin.Flag("skip-ssl-validation-cf", "Skip cert validation (for dev environments").
		OverrideDefaultFromEnvar("SKIP_SSL_VALIDATION_CF").Default("false").BoolVar(&c.SkipSSLCF)
	kingpin.Flag("skip-ssl-validation-splunk", "Skip cert validation (for dev environments").
		OverrideDefaultFromEnvar("SKIP_SSL_VALIDATION_SPLUNK").Default("false").BoolVar(&c.SkipSSLSplunk)
	kingpin.Flag("subscription-id", "Id for the subscription.").
		OverrideDefaultFromEnvar("FIREHOSE_SUBSCRIPTION_ID").Default("splunk-firehose").StringVar(&c.SubscriptionID)
	kingpin.Flag("firehose-keep-alive", "Keep Alive duration for the firehose consumer").
		OverrideDefaultFromEnvar("FIREHOSE_KEEP_ALIVE").Default("25s").DurationVar(&c.KeepAlive)

	kingpin.Flag("add-app-info", "Query API to fetch app details").
		OverrideDefaultFromEnvar("ADD_APP_INFO").Default("false").BoolVar(&c.AddAppInfo)
	kingpin.Flag("ignore-missing-app", "If app is missing, stop repeatedly querying app info from PCF").
		OverrideDefaultFromEnvar("IGNORE_MISSING_APP").Default("true").BoolVar(&c.IgnoreMissingApps)
	kingpin.Flag("missing-app-cache-invalidate-ttl", "How frequently the missing app info cache invalidates").
		OverrideDefaultFromEnvar("MISSING_APP_CACHE_INVALIDATE_TTL").Default("0s").DurationVar(&c.MissingAppCacheTTL)
	kingpin.Flag("app-cache-invalidate-ttl", "How frequently the app info local cache invalidates").
		OverrideDefaultFromEnvar("APP_CACHE_INVALIDATE_TTL").Default("0s").DurationVar(&c.AppCacheTTL)
	kingpin.Flag("org-space-cache-invalidate-ttl", "How frequently the org and space cache invalidates").
		OverrideDefaultFromEnvar("ORG_SPACE_CACHE_INVALIDATE_TTL").Default("72h").DurationVar(&c.OrgSpaceCacheTTL)
	kingpin.Flag("app-limits", "Restrict to APP_LIMITS most updated apps per request when populating the app metadata cache").
		OverrideDefaultFromEnvar("APP_LIMITS").Default("0").IntVar(&c.AppLimits)

	kingpin.Flag("boltdb-path", "Bolt Database path ").
		Default("cache.db").OverrideDefaultFromEnvar("BOLTDB_PATH").StringVar(&c.BoltDBPath)
	kingpin.Flag("events", fmt.Sprintf("Comma separated list of events you would like. Valid options are %s", events.AuthorizedEvents())).
		OverrideDefaultFromEnvar("EVENTS").Default("ValueMetric,CounterEvent,ContainerMetric").StringVar(&c.WantedEvents)
	kingpin.Flag("extra-fields", "Extra fields you want to annotate your events with, example: '--extra-fields=env:dev,something:other ").
		OverrideDefaultFromEnvar("EXTRA_FIELDS").Default("").StringVar(&c.ExtraFields)

	kingpin.Flag("flush-interval", "Every interval flushes to Splunk Http Event Collector server").
		OverrideDefaultFromEnvar("FLUSH_INTERVAL").Default("5s").DurationVar(&c.FlushInterval)
	kingpin.Flag("consumer-queue-size", "Consumer queue buffer size").
		OverrideDefaultFromEnvar("CONSUMER_QUEUE_SIZE").Default("10000").IntVar(&c.QueueSize)
	kingpin.Flag("hec-batch-size", "Batchsize of the events pushing to HEC").
		OverrideDefaultFromEnvar("HEC_BATCH_SIZE").Default("100").IntVar(&c.BatchSize)
	kingpin.Flag("hec-retries", "Number of retries before dropping events").
		OverrideDefaultFromEnvar("HEC_RETRIES").Default("5").IntVar(&c.Retries)
	kingpin.Flag("hec-workers", "How many workers (concurrency) when post data to HEC").
		OverrideDefaultFromEnvar("HEC_WORKERS").Default("8").IntVar(&c.HecWorkers)
	kingpin.Flag("splunk-version", "Splunk version will determine how metadata fields are ingested for HEC '--splunk-version=7.2	").
		OverrideDefaultFromEnvar("SPLUNK_VERSION").Default("7.2").StringVar(&c.SplunkVersion)

	kingpin.Flag("enable-event-tracing", "Enable event trace logging: Adds splunk trace logging fields to events. uuid, subscription-id, nozzle event counter").
		OverrideDefaultFromEnvar("ENABLE_EVENT_TRACING").Default("false").BoolVar(&c.TraceLogging)
	kingpin.Flag("debug", "Enable debug mode: forward to standard out instead of splunk").
		OverrideDefaultFromEnvar("DEBUG").Default("false").BoolVar(&c.Debug)

	kingpin.Parse()
	return c
}

func (c *Config) ToMap() map[string]interface{} {
	data, _ := json.Marshal(c)
	var r map[string]interface{}
	json.Unmarshal(data, &r)
	return r
}
