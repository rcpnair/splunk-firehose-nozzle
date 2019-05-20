package splunknozzle_test

import (
	"os"
	"time"

	. "github.com/cloudfoundry-community/splunk-firehose-nozzle/splunknozzle"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {

	Context("Env config parsing", func() {
		var (
			version = "1.0"
			branch  = "develop"
			commit  = "08a9e9bd557ca9038e9b391d9a77d47aa56210a3"
			buildos = "Linux"
		)

		BeforeEach(func() {
			// FIX "nozzle.test: error: unknown short flag '-t', try --help" error when coverage
			os.Args = os.Args[:1]
			os.Clearenv()

			os.Setenv("API_ENDPOINT", "api.bosh-lite.com")
			os.Setenv("API_USER", "admin")
			os.Setenv("API_PASSWORD", "abc123")

			os.Setenv("SPLUNK_TOKEN", "sometoken")
			os.Setenv("SPLUNK_HOST", "splunk.example.com")
			os.Setenv("SPLUNK_INDEX", "splunk_index")
		})

		It("parses config from environment", func() {
			os.Setenv("JOB_NAME", "my-job")
			os.Setenv("JOB_INDEX", "2")
			os.Setenv("JOB_HOST", "nozzle.example.com")

			os.Setenv("SKIP_SSL_VALIDATION_CF", "true")
			os.Setenv("SKIP_SSL_VALIDATION_SPLUNK", "true")

			os.Setenv("FIREHOSE_SUBSCRIPTION_ID", "my-nozzle")
			os.Setenv("FIREHOSE_KEEP_ALIVE", "42s")

			os.Setenv("ADD_APP_INFO", "true")
			os.Setenv("IGNORE_MISSING_APP", "true")
			os.Setenv("MISSING_APP_CACHE_INVALIDATE_TTL", "100s")
			os.Setenv("APP_CACHE_INVALIDATE_TTL", "100s")
			os.Setenv("APP_LIMITS", "2000")

			os.Setenv("BOLTDB_PATH", "foo.db")
			os.Setenv("EVENTS", "LogMessage")
			os.Setenv("EXTRA_FIELDS", "foo:bar")

			os.Setenv("FLUSH_INTERVAL", "43s")
			os.Setenv("CONSUMER_QUEUE_SIZE", "15000")
			os.Setenv("HEC_RETRIES", "10")
			os.Setenv("HEC_WORKERS", "5")
			os.Setenv("SPLUNK_VERSION", "6.2")

			os.Setenv("ENABLE_EVENT_TRACING", "true")
			os.Setenv("DEBUG", "true")

			c := NewConfigFromCmdFlags(version, branch, commit, buildos)

			Expect(c.ApiEndpoint).To(Equal("api.bosh-lite.com"))
			Expect(c.User).To(Equal("admin"))
			Expect(c.Password).To(Equal("abc123"))

			Expect(c.SplunkHost).To(Equal("splunk.example.com"))
			Expect(c.SplunkToken).To(Equal("sometoken"))
			Expect(c.SplunkIndex).To(Equal("splunk_index"))

			Expect(c.JobName).To(Equal("my-job"))
			Expect(c.JobIndex).To(Equal("2"))
			Expect(c.JobHost).To(Equal("nozzle.example.com"))

			Expect(c.SkipSSLCF).To(BeTrue())
			Expect(c.SkipSSLSplunk).To(BeTrue())

			Expect(c.SubscriptionID).To(Equal("my-nozzle"))
			Expect(c.KeepAlive).To(Equal(42 * time.Second))

			Expect(c.AddAppInfo).To(BeTrue())
			Expect(c.IgnoreMissingApps).To(BeTrue())
			Expect(c.MissingAppCacheTTL).To(Equal(100 * time.Second))
			Expect(c.AppCacheTTL).To(Equal(100 * time.Second))
			Expect(c.AppLimits).To(Equal(2000))

			Expect(c.BoltDBPath).To(Equal("foo.db"))
			Expect(c.WantedEvents).To(Equal("LogMessage"))
			Expect(c.ExtraFields).To(Equal("foo:bar"))

			Expect(c.FlushInterval).To(Equal(43 * time.Second))
			Expect(c.QueueSize).To(Equal(15000))
			Expect(c.BatchSize).To(Equal(100))
			Expect(c.Retries).To(Equal(10))
			Expect(c.HecWorkers).To(Equal(5))
			Expect(c.SplunkVersion).To(Equal("6.2"))

			Expect(c.Version).To(Equal(version))
			Expect(c.Branch).To(Equal(branch))
			Expect(c.Commit).To(Equal(commit))
			Expect(c.BuildOS).To(Equal(buildos))

			Expect(c.TraceLogging).To(BeTrue())
			Expect(c.Debug).To(BeTrue())
		})

		It("check defaults", func() {
			c := NewConfigFromCmdFlags(version, branch, commit, buildos)

			Expect(c.JobName).To(Equal("splunk-nozzle"))
			Expect(c.JobIndex).To(Equal("-1"))
			Expect(c.JobHost).To(Equal(""))

			Expect(c.SkipSSLCF).To(BeFalse())
			Expect(c.SkipSSLCF).To(BeFalse())
			Expect(c.SubscriptionID).To(Equal("splunk-firehose"))
			Expect(c.KeepAlive).To(Equal(25 * time.Second))

			Expect(c.AddAppInfo).To(BeFalse())
			Expect(c.IgnoreMissingApps).To(BeTrue())
			Expect(c.MissingAppCacheTTL).To(Equal(0 * time.Second))
			Expect(c.AppCacheTTL).To(Equal(0 * time.Second))
			Expect(c.AppLimits).To(Equal(0))

			Expect(c.BoltDBPath).To(Equal("cache.db"))
			Expect(c.WantedEvents).To(Equal("ValueMetric,CounterEvent,ContainerMetric"))
			Expect(c.ExtraFields).To(Equal(""))

			Expect(c.FlushInterval).To(Equal(5 * time.Second))
			Expect(c.QueueSize).To(Equal(10000))
			Expect(c.BatchSize).To(Equal(100))
			Expect(c.Retries).To(Equal(5))
			Expect(c.HecWorkers).To(Equal(8))
			Expect(c.SplunkVersion).To(Equal("7.2"))

			Expect(c.TraceLogging).To(BeFalse())
			Expect(c.Debug).To(BeFalse())
		})
	})

	Context("Flags config parsing", func() {
		var (
			version = "1.0"
			branch  = "develop"
			commit  = "08a9e9bd557ca9038e9b391d9a77d47aa56210a3"
			buildos = "Linux"
		)

		BeforeEach(func() {
			os.Clearenv()
			// FIX "nozzle.test: error: unknown short flag '-t', try --help" error when coverage
			args := []string{
				"splunk-firehose-nozzle",
				"--api-endpoint=api.bosh-lite.comc",
				"--user=adminc",
				"--password=abc123c",
				"--splunk-host=splunk.example.comc",
				"--splunk-token=sometokenc",
				"--splunk-index=splunk_indexc",
				"--job-name=my-jobc",
				"--job-index=3",
				"--job-host=nozzle.example.comc",
				"--skip-ssl-validation-cf",
				"--skip-ssl-validation-splunk",
				"--subscription-id=my-nozzlec",
				"--firehose-keep-alive=24s",
				"--add-app-info",
				"--ignore-missing-app",
				"--missing-app-cache-invalidate-ttl=15s",
				"--app-cache-invalidate-ttl=15s",
				"--app-limits=35",
				"--boltdb-path=foo.dbc",
				"--events=LogMessagec",
				"--extra-fields=foo:barc",
				"--flush-interval=34s",
				"--consumer-queue-size=2323",
				"--hec-batch-size=1234",
				"--hec-retries=9",
				"--hec-workers=16",
				"--splunk-version=5.2",
				"--enable-event-tracing",
				"--debug",
			}
			os.Args = args
		})

		It("parses config from cli flags", func() {
			c := NewConfigFromCmdFlags(version, branch, commit, buildos)

			Expect(c.ApiEndpoint).To(Equal("api.bosh-lite.comc"))
			Expect(c.User).To(Equal("adminc"))
			Expect(c.Password).To(Equal("abc123c"))

			Expect(c.SplunkHost).To(Equal("splunk.example.comc"))
			Expect(c.SplunkToken).To(Equal("sometokenc"))
			Expect(c.SplunkIndex).To(Equal("splunk_indexc"))

			Expect(c.JobName).To(Equal("my-jobc"))
			Expect(c.JobIndex).To(Equal("3"))
			Expect(c.JobHost).To(Equal("nozzle.example.comc"))

			Expect(c.SkipSSLCF).To(BeTrue())
			Expect(c.SkipSSLSplunk).To(BeTrue())
			Expect(c.SubscriptionID).To(Equal("my-nozzlec"))
			Expect(c.KeepAlive).To(Equal(24 * time.Second))

			Expect(c.AddAppInfo).To(BeTrue())
			Expect(c.IgnoreMissingApps).To(BeTrue())
			Expect(c.MissingAppCacheTTL).To(Equal(15 * time.Second))
			Expect(c.AppCacheTTL).To(Equal(15 * time.Second))
			Expect(c.AppLimits).To(Equal(35))

			Expect(c.BoltDBPath).To(Equal("foo.dbc"))
			Expect(c.WantedEvents).To(Equal("LogMessagec"))
			Expect(c.ExtraFields).To(Equal("foo:barc"))

			Expect(c.FlushInterval).To(Equal(34 * time.Second))
			Expect(c.QueueSize).To(Equal(2323))
			Expect(c.BatchSize).To(Equal(1234))
			Expect(c.Retries).To(Equal(9))
			Expect(c.HecWorkers).To(Equal(16))
			Expect(c.SplunkVersion).To(Equal("5.2"))

			Expect(c.Debug).To(BeTrue())
			Expect(c.TraceLogging).To(BeTrue())

			Expect(c.Version).To(Equal(version))
			Expect(c.Branch).To(Equal(branch))
			Expect(c.Commit).To(Equal(commit))
			Expect(c.BuildOS).To(Equal(buildos))

		})
	})
})
