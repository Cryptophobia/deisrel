package actions

import (
	"text/template"
)

const (
	// TODO: https://github.com/deis/deisrel/issues/11
	generateParamsTplStr = `#helm:generate helm template -o $HELM_GENERATE_DIR/manifests/deis-objectstorage-secret.yaml -d $HELM_GENERATE_FILE $HELM_GENERATE_DIR/tpl/deis-objectstorage-secret.yaml
#
# This is the main configuration file for Deis object storage. The values in
# this file are passed into the appropriate services so that they can configure
# themselves for persisting data in object storage.
#
# In general, all object storage credentials must be able to read and write to
# the container or bucket they are configured to use.
#
# When you change values in this file, make sure to re-run 'helm generate'
# on this chart.

# Set the storage backend
#
# Valid values are:
# - filesystem: Store persistent data on ephemeral disk
# - s3: Store persistent data in AWS S3 (configure in S3 section)
# - azure: Store persistent data in Azure's object storage
# - gcs: Store persistent data in Google Cloud Storage
# - minio: Store persistent data on in-cluster Minio server
storage = "minio"

[s3]
accesskey = "YOUR KEY HERE"
secretkey = "YOUR SECRET HERE"
# Any S3 region
region = "us-west-1"
# Your buckets.
registry_bucket = "your-registry-bucket-name"
database_bucket = "your-database-bucket-name"
builder_bucket = "your-builder-bucket-name"

[azure]
accountname = "YOUR ACCOUNT NAME"
accountkey = "YOUR ACCOUNT KEY"
registry_container = "your-registry-container-name"
database_container = "your-database-container-name"
builder_container = "your-builder-container-name"

[gcs]
# key_json is expanded into a JSON file on the remote server. It must be
# well-formatted JSON data.
key_json = '''Paste JSON data here.'''
registry_bucket = "your-registry-bucket-name"
database_bucket = "your-database-bucket-name"
builder_bucket = "your-builder-bucket-name"

[minio]
org = "{{.Minio.Org}}"
pullPolicy = "{{.Minio.PullPolicy}}"
dockerTag = "{{.Minio.Tag}}"

[builder]
org = "{{.Builder.Org}}"
pullPolicy = "{{.Builder.PullPolicy}}"
dockerTag = "{{.Builder.Tag}}"

[slugbuilder]
org = "{{.SlugBuilder.Org}}"
pullPolicy = "{{.SlugBuilder.PullPolicy}}"
dockerTag = "{{.SlugBuilder.Tag}}"

[dockerbuilder]
org = "{{.DockerBuilder.Org}}"
pullPolicy = "{{.DockerBuilder.PullPolicy}}"
dockerTag = "{{.DockerBuilder.Tag}}"

[controller]
org = "{{.Controller.Org}}"
pullPolicy = "{{.Controller.PullPolicy}}"
dockerTag = "{{.Controller.Tag}}"

[slugrunner]
org = "{{.SlugRunner.Org}}"
pullPolicy = "{{.SlugRunner.PullPolicy}}"
dockerTag = "{{.SlugRunner.Tag}}"

[database]
org = "{{.Database.Org}}"
pullPolicy = "{{.Database.PullPolicy}}"
dockerTag = "{{.Database.Tag}}"

[registry]
org = "{{.Registry.Org}}"
pullPolicy = "{{.Registry.PullPolicy}}"
dockerTag = "{{.Registry.Tag}}"

[workflowManager]
org = "{{.WorkflowManager.Org}}"
pullPolicy = "{{.WorkflowManager.PullPolicy}}"
dockerTag = "{{.WorkflowManager.Tag}}"

[logger]
org = "{{.Logger.Org}}"
pullPolicy = "{{.Logger.PullPolicy}}"
dockerTag = "{{.Logger.Tag}}"

[router]
org = "{{.Router.Org}}"
pullPolicy = "{{.Router.PullPolicy}}"
dockerTag = "{{.Router.Tag}}"

[fluentd]
org = "{{.FluentD.Org}}"
pullPolicy = "{{.FluentD.PullPolicy}}"
dockerTag = "{{.FluentD.Tag}}"

[grafana]
org = "{{.Grafana.Org}}"
pullPolicy = "{{.Grafana.PullPolicy}}"
dockerTag = "{{.Grafana.Tag}}"

[influxdb]
org = "{{.InfluxDB.Org}}"
pullPolicy = "{{.InfluxDB.PullPolicy}}"
dockerTag = "{{.InfluxDB.Tag}}"

[telegraf]
org = "{{.Telegraf.Org}}"
pullPolicy = "{{.Telegraf.PullPolicy}}"
dockerTag = "{{.Telegraf.Tag}}"

[stdoutmetrics]
org = "{{.StdoutMetrics.Org}}"
pullPolicy = "{{.StdoutMetrics.PullPolicy}}"
dockerTag = "{{.StdoutMetrics.Tag}}"
`
)

var (
	generateParamsTpl = template.Must(template.New("generateParamsTpl").Parse(generateParamsTplStr))
)
