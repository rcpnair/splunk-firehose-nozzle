[![CircleCI](https://circleci.com/gh/git-lfs/git-lfs.svg?style=shield&circle-token=856152c2b02bfd236f54d21e1f581f3e4ebf47ad)](https://circleci.com/gh/cloudfoundry-community/splunk-firehose-nozzle)
## Splunk Nozzle

Cloud Foundry Firehose-to-Splunk Nozzle

### Usage
Splunk nozzle is used to stream Cloud Foundry Firehose events to Splunk HTTP Event Collector. Using pre-defined Splunk sourcetypes, the nozzle automatically parses the events and enriches them with additional metadata before forwarding to Splunk. For detailed descriptions of each Firehose event type and their fields, refer to underlying [dropsonde protocol](https://github.com/cloudfoundry/dropsonde-protocol). Below is a mapping of each Firehose event type to its corresponding Splunk sourcetype. Refer to [Searching Events](#searching-events) for example Splunk searches.

| Firehose event type | Splunk sourcetype | Description
|---|---|---
| Error | `cf:error` | An Error event represents an error in the originating process
| HttpStartStop | `cf:httpstartstop` | An HttpStartStop event represents the whole lifecycle of an HTTP request
| LogMessage | `cf:logmessage` | A LogMessage contains a "log line" and associated metadata
| ContainerMetric | `cf:containermetric` | A ContainerMetric records resource usage of an app in a container
| CounterEvent | `cf:counterevent` | A CounterEvent represents the increment of a counter
| ValueMetric | `cf:valuemetric` | A ValueMetric indicates the value of a metric at an instant in time

In addition, logs from the nozzle itself are of sourcetype `cf:splunknozzle`.

### Setup

The Nozzle requires a user with the scope `doppler.firehose` and
`cloud_controller.admin_read_only` (the latter is only required if `ADD_APP_INFO` is true). If `cloud_controller.admin_read_only` is not
available in the system, switch to use `cloud_controller.admin`.

You can either
* Add the user manually using [uaac](https://github.com/cloudfoundry/cf-uaac)
* Add a new user to the deployment manifest; see [uaa.scim.users](https://github.com/cloudfoundry/uaa-release/blob/master/jobs/uaa/spec)

Manifest example:

```yaml
uaa:
  scim:
    users:
      - splunk-firehose|password123|cloud_controller.admin_read_only,doppler.firehose
```

`uaac` example:
```shell
uaac target https://uaa.[system domain url]
uaac token client get admin -s [admin client credentials secret]
uaac -t user add splunk-nozzle --password password123 --emails na
uaac -t member add cloud_controller.admin_read_only splunk-nozzle
uaac -t member add doppler.firehose splunk-nozzle
```

`cloud_controller.admin_read_only` will work for cf v241
or later. Earlier versions should use `cloud_controller.admin` instead.

- - - -
#### Environment Parameters
You can declare parameters by making a copy of the scripts/nozzle.sh.template.
* `DEBUG`: Enable debug mode (forward to standard out instead of Splunk).

__Cloud Foundry configuration parameters:__
* `API_ENDPOINT`: Cloud Foundry API endpoint address.
* `API_USER`: Cloud Foundry user name. (Must have scope described above)
* `API_PASSWORD`: Cloud Foundry user password.

__Splunk configuration parameters:__
* `SPLUNK_TOKEN`: [Splunk HTTP event collector token](http://docs.splunk.com/Documentation/Splunk/latest/Data/UsetheHTTPEventCollector/).
* `SPLUNK_HOST`: Splunk HTTP event collector host. example: https://example.cloud.splunk.com:8088
* `SPLUNK_INDEX`: The Splunk index events will be sent to. Warning: Setting an invalid index will cause events to be lost. This index must match one of the selected indexes for the Splunk HTTP event collector token used for the SPLUNK_TOKEN parameter.

__Advanced Configuration Features:__
* `JOB_NAME`: Tags nozzle log events with job name.
* `JOB_INDEX`: Tags nozzle log events with job index.
* `JOB_HOST`: Tags nozzle log events with job host.
* `SKIP_SSL_VALIDATION_CF`: Skips SSL certificate validation for connection to Cloud Foundry. Secure communications will not check SSL certificates against a trusted certificate authority.
This is recommended for dev environments only.
* `SKIP_SSL_VALIDATION_SPLUNK`: Skips SSL certificate validation for connection to Splunk. Secure communications will not check SSL certificates against a trusted certificate authority.
This is recommended for dev environments only.
* `FIREHOSE_SUBSCRIPTION_ID`: Tags nozzle events with a Firehose subscription id. See https://docs.pivotal.io/pivotalcf/1-11/loggregator/log-ops-guide.html.
* `FIREHOSE_KEEP_ALIVE`: Keep alive duration for the Firehose consumer.
* `ADD_APP_INFO`: Enriches raw data with app details.
* `IGNORE_MISSING_APP`: If the application is missing, then stop repeatedly querying application info from Cloud Foundry.
* `MISSING_APP_CACHE_INVALIDATE_TTL`:  How frequently the missing app info cache invalidates.
* `APP_CACHE_INVALIDATE_TTL`: How frequently the app info local cache invalidates.
* `APP_LIMITS`: Restrict to APP_LIMITS the most updated apps per request when populating the app metadata cache.
* `BOLTDB_PATH`: Bolt database path.
* `EVENTS`: A comma separated list of events to include. Possible values: ValueMetric,CounterEvent,Error,LogMessage,HttpStartStop,ContainerMetric
* `EXTRA_FIELDS`: Extra fields to annotate your events with (format is key:value,key:value).
* `FLUSH_INTERVAL`: Time interval for flushing queue to Splunk regardless of CONSUMER_QUEUE_SIZE. Protects against stale events in low throughput systems.
* `CONSUMER_QUEUE_SIZE`: Sets the internal consumer queue buffer size. Events will be pushed to Splunk after queue is full.
* `HEC_BATCH_SIZE`: Set the batch size for the events to push to HEC (Splunk HTTP Event Collector).
* `HEC_RETRIES`: Retry count for sending events to Splunk. After expiring, events will begin dropping causing data loss.
* `HEC_WORKERS`: Set the amount of Splunk HEC workers to increase concurrency while ingesting in Splunk.
* `SPLUNK_VERSION`: The Splunk version that determines how HEC ingests metadata fields. For example: 7.2
* `ENABLE_EVENT_TRACING`: Enables event trace logging. Splunk events will now contain a UUID, Splunk Nozzle Event Counts, and a Subscription-ID for Splunk correlation searches.

- - - -

### Push as an App to Cloud Foundry

[splunk-firehose-nozzle-release](https://github.com/cloudfoundry-community/splunk-firehose-nozzle-release)
packages this code into a [BOSH](https://bosh.io) release for deployment. The code could also be run on
Cloud Foundry as an application. See the **Setup** section for details
on making a user and credentials.

1. Download the latest release

    ```shell
    git clone https://github.com/cloudfoundry-community/splunk-firehose-nozzle.git
    cd splunk-firehose-nozzle
    ```

1. Authenticate to Cloud Foundry

    ```shell
    cf login -a https://api.[your cf system domain] -u [your id]
    ```

1. Copy the manifest template and fill in needed values (using the credentials created during setup)

    ```shell
    vim ci/nozzle_manifest.yml
    ```

1. Push the nozzle

    ```shell
    make deploy-nozzle
    ```

#### Dump application info to boltdb ####
If in production there are lots of PCF applications(say tens of thousands) and if the user would like to enrich
application logs by including application meta data,querying all application metadata information from PCF may take some time.
For example if we include, add app name, space ID, space name, org ID and org name to the events.
If there are multiple instances of Spunk nozzle deployed the situation will be even worse, since each of the Splunk nozzle(s) will query all applications meta data and
cache the meta data information to the local boltdb file. These queries will introduce load to the PCF system and could potentially take a long time to finish.
Users can run this tool to generate a copy of all application meta data and copy this to each Splunk nozzle deployment. Each Splunk nozzle can pick up the cache copy and update the cache file incrementally afterwards.

Example of how to run the dump application info tool:

```
$ cd tools/dump_app_info
$ go build dump_app_info.go
$ ./dump_app_info --skip-ssl-validation --api-endpoint=https://<your api endpoint> --user=<api endpoint login username> --password=<api endpoint login password>
```

After populating the application info cache file, user can copy to different Splunk nozzle deployments and start Splunk nozzle to pick up this cache file by
specifying correct "--boltdb-path" flag or "BOLTDB_PATH" environment variable.

###Per application index routing (deprecates instructions below)
in your app manifest provide an env var called SPLUNK_INDEX and assign it the index you would like to send the data to

```
applications:
- name: console
  memory: 256M
  disk_quota: 256M
  host: console
  timeout: 180
  buildpack: https://github.com/SUSE/stratos-buildpack
  health-check-type: port
  services:
  - splunk-index
  env:
    SPLUNK_INDEX: testing_index
```

#### Index routing
Index routing is a feature that can be used to send different Cloud Foundry logs to different indexes for better ACL and data retention control in Splunk.
Logs can be routed using fields such as app ID/name, space ID/name or org ID/name.
Users can configure the Splunk configuration files props.conf and transforms.conf on Splunk indexers or Splunk Heavy Forwarders if deployed.

The following is an example of how to route application ID `95930b4e-c16c-478e-8ded-5c6e9c5981f8` to a Splunk `prod` index.

$SPLUNK_HOME/etc/system/local/props.conf

```
[cf:logmessage]
TRANSFORMS-index_routing = route_data_to_index_by_field_cf_app_id
```


$SPLUNK_HOME/etc/system/local/transforms.conf

```
[route_data_to_index_by_field_cf_app_id]
REGEX = "(\w+)":"95930b4e-c16c-478e-8ded-5c6e9c5981f8"
DEST_KEY = _MetaData:Index
FORMAT = prod
```

Another example is routing application logs from any Cloud Foundry orgs whose names are prefixed with `sales` to a Splunk `sales` index.

$SPLUNK_HOME/etc/system/local/props.conf

```
[cf:logmessage]
TRANSFORMS-index_routing = route_data_to_index_by_field_cf_org_name
```


$SPLUNK_HOME/etc/system/local/transforms.conf

```
[route_data_to_index_by_field_cf_org_name]
REGEX = "cf_org_name":"(sales.*)"
DEST_KEY = _MetaData:Index
FORMAT = sales
```

#### Troubleshooting
In most cases you will only need to troubleshoot from Splunk which includes not only firehose data but also this Splunk nozzle internal logs.
However, if the nozzle is still not forwarding any data, a good place to start is to get the application internal logs directly:

```shell
cf logs splunk-firehose-nozzle
```

A common mis-configuration occurs when having invalid or unsigned certificate(s) for the Cloud Foundry API endpoint.
In the case for non-production environments, you can set `SKIP_SSL_VALIDATION` to `true` in manifest.yml before re-deploying the app.

#### Searching Events

Here are two short Splunk queries to start exploring some of the Cloud Foundry events in Splunk.

```
sourcetype="cf:valuemetric"
    | stats avg(value) by job_instance, name
```

```
sourcetype="cf:counterevent"
    | eval job_and_name=source+"-"+name
    | stats values(job_and_name)
```
### Development

#### Software Requirements

Make sure you have the following installed on your workstation:

| Software | Version
| --- | --- |
| go | go1.7.x
| glide | 0.12.x

Then install all dependent packages via [Glide](https://glide.sh/):

```
$ cd <REPO_ROOT_DIRECTORY>
$ make installdeps
```

#### Environment

For development against [bosh-lite](https://github.com/cloudfoundry/bosh-lite),
copy `tools/nozzle.sh.template` to `tools/nozzle.sh` and supply missing values:

```
$ cp script/dev.sh.template tools/nozzle.sh
$ chmod +x tools/nozzle.sh
```

Build project:

```
$ make VERSION=1.0
```

Run tests with [Ginkgo](http://onsi.github.io/ginkgo/)

```
$ ginkgo -r
```

Run all kinds of testing

```
$ make test # run all unittest
$ make race # test if there is race condition in the code
$ make vet  # examine GoLang code
$ make cov  # code coverage test and code coverage html report
```

Or run all testings: unit test, race condition test, code coverage etc
```
$ make testall
```

Run app

```
# this will run: go run main.go
$ ./tools/nozzle.sh
```

#### CI

https://concourse.cfplatformeng.com/teams/splunk/pipelines/splunk-firehose-tile-build
