# kayron

This repository contains the operator responsible for our automated change
management. The relevant components here are described by the `server` and
`operator` packages. Kayron's sole responsibility is to keep our infrastructure
and services up to date, according to our [release artifacts] managed in Github.

![Kayron Overview](.github/assets/Kayron-Overview.svg)

### Server

The server is a simple HTTP backend exposing the Kayron's `/metrics` endopoint.

```
curl -s http://127.0.0.1:7777/metrics
```

```
# HELP go_gc_duration_seconds A summary of the wall-time pause (stop-the-world) duration in garbage collection cycles.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 0
go_gc_duration_seconds{quantile="0.25"} 0
go_gc_duration_seconds{quantile="0.5"} 0
```

### Worker

The worker is a [custom task engine] executing a directed acyclic graph of
worker handlers iteratively. New worker handlers can be added easily by
implementing the handler interface and registering the new handler in the
operator chain. Kayron's operator chain is executed by a sequential worker
engine.

```golang
type Interface interface {
	// Ensure is the minimal worker handler interface that all users have to
	// implement for their own business logic, regardless of the underlying worker
	// engine.
	//
	// Ensure executes the handler specific business logic in order to complete
	// the given task, if possible. Any error returned will be emitted using the
	// underlying logger interface, unless the injected metrics registry is
	// configured to filter the received error.
	Ensure() error
}
```

### Operator

Kayron is a change management controller implementing the [operator pattern], in
this particular case without the involvement of [Kubernetes]. The main goroutine
for the operator's reconciliation loop is the operator worker handler running a
sequence of steps according to their operator functions located in
[pkg/operator](./pkg/operator/). Secondary worker handlers may be executed
within their own isolated failure domain.

Operators try to continuously drive the current state of a system towards the
desired state of a system. In our case, the current state is represented by the
already deployed CloudFormation templates in AWS. The desired state is then
represented by the template changes introduced in any given environment. Those
changes my simply be a version change from `v1.8.2` to `v1.8.3` in case of a
service update, or a more complex CloudFormation template change that modifies
physical resources in AWS. Kayron’s job is to continuously check for any
detectable drift between current and desired state, and apply any changes made
upon detection.

For instance, if a new Release were to be created in the
[Specta](https://github.com/0xSplits/specta) repository, then Kayron will
generate all CloudFormation templates anew using the changed resource details in
the underlying cache channels, and reconcile the resulting CloudFormation
templates against AWS, if the specified policies permit the proposed update.
This process will then eventually result in a stack update in CloudFormation, so
that Kayron may also keep the respective deployment status up to date in Github.

![Operator Design](.github/assets/Operator-Design.svg)

### Triggers

1. The operator executes its graph of business logic sequentially on a time
   based schedule. E.g. the entire chain runs every 10 seconds and pulls all the
   data it needs from external APIs.
2. The operator can also executes its graph on demand in case a trigger event
   like a webhook notification may occurred. In this case the time based
   schedule described above resets so that the received webhook event creates a
   new 10 second interval.

### Usage

At its core, Kayron is a simple [Cobra] command line tool, providing e.g. the
daemon command to start the long running `server` and `worker` processes.

```
kayron -h
```

```
Golang based operator microservice.

Usage:
  kayron [flags]
  kayron [command]

Available Commands:
  daemon      Execute Kayron's long running process for running the operator.
  deploy      Manually trigger a CloudFormation stack update.
  lint        Validate the release configuration under the given path.
  version     Print the version information for this command line tool.

Flags:
  -h, --help   help for kayron

Use "kayron [command] --help" for more information about a command.
```

### Development

As a convention, Kayron's `.env` file should remain simple and generic. A
reasonable setting within that config file is e.g. `KAYRON_LOG_LEVEL`. Kayron
requires the standard AWS credentials format to be properly setup. More specific
organization and environment related settings must be provided separately. And
so running e.g. the Kayron daemon requires several additional environment
variables to be injected.

- `KAYRON_CLOUDFORMATION_STACK`, the CloudFormation stack name to reconcile, e.g. `server-test`.
- `KAYRON_ENVIRONMENT`, the environment Kayron is running in, one of `development` `testing` `staging` `production`.
- `KAYRON_GITHUB_TOKEN`, the Github token to use for fetching releases, requires the `repo` scope.
- `KAYRON_RELEASE_SOURCE`, the Github repository containing releases, e.g. https://github.com/0xSplits/releases.
- `KAYRON_S3_BUCKET`, the S3 bucket to upload CloudFormation templates, e.g. `splits-cf-templates`.

---

- `KAYRON_CLOUDFORMATION_PARAMETERS`, the optional CloudFormation parameter overwrites, in the format `key:value,foo:bar`.
- `KAYRON_CLOUDFORMATION_TAGS`, the optional CloudFormation tag overwrites, in the format `key:value,foo:bar`.

```
kayron daemon
```

```
{ "time":"2025-07-04 14:09:06", "level":"info", "message":"daemon is launching procs", "environment":"development", "caller":".../pkg/daemon/daemon.go:38" }
{ "time":"2025-07-04 14:09:06", "level":"info", "message":"server is accepting calls", "address":"127.0.0.1:7777",  "caller":".../pkg/server/server.go:95" }
{ "time":"2025-07-04 14:09:06", "level":"info", "message":"worker is executing tasks", "pipelines":"1",             "caller":".../pkg/worker/worker.go:110" }
```

There is an [integration-test](.github/workflows/integration-test.yaml) workflow
to verify several aspects of Kayron's various responsibilities. The test
verifies that all operator functions are free from race conditions, which is
critical, because several operator functions are running in parallel each and
every reconciliation loop. The test also verifies that the operator's
reconciliation loops are not unexpectedly cancelled due to internal data
inconsistencies when looking all kinds of current and desired state for the
internal artifact cache. The required AWS credentials need the
[ViewOnlyAccess], [AmazonEC2ContainerRegistryReadOnly] and
[ResourceGroupsandTagEditorReadOnlyAccess] permissions. The required Github auth
token needs the repository scope.

```
go test -tags=integration ./pkg/operator -v -race -run Test_Operator_Integration
```

```yaml
=== RUN   Test_Operator_Integration
{
    "time": "2025-09-17 15:31:59",
    "level": "debug",
    "message": "resetting operator cache",
    "caller": "/Users/xh3b4sd/project/0xSplits/kayron/pkg/cache/delete.go:9"
}
{
    "time": "2025-09-17 15:32:01",
    "level": "debug",
    "message": "resolved ref for github repository",
    "environment": "testing",
    "ref": "176ffefa272b210eb3be269887d665a682dbd548",
    "repository": "https://github.com/0xSplits/releases",
    "caller": "/Users/xh3b4sd/project/0xSplits/kayron/pkg/operator/release/ensure.go:66"
}
{
    "time": "2025-09-17 15:32:02",
    "level": "debug",
    "message": "caching release artifact",
    "deploy": "branch=preview",
    "github": "infrastructure",
    "preview": "false",
    "provider": "cloudformation",
    "caller": "/Users/xh3b4sd/project/0xSplits/kayron/pkg/cache/create.go:17"
}
{
    "time": "2025-09-17 15:32:02",
    "level": "debug",
    "message": "caching release artifact",
    "deploy": "release=v0.1.1",
    "docker": "splits-lite",
    "github": "splits-lite",
    "preview": "false",
    "caller": "/Users/xh3b4sd/project/0xSplits/kayron/pkg/cache/create.go:17"
}
{
    "time": "2025-09-17 15:32:02",
    "level": "debug",
    "message": "caching release artifact",
    "deploy": "branch=fancy-feature-branch",
    "docker": "splits-lite",
    "github": "splits-lite",
    "preview": "true",
    "caller": "/Users/xh3b4sd/project/0xSplits/kayron/pkg/cache/create.go:17"
}
{
    "time": "2025-09-17 15:32:02",
    "level": "debug",
    "message": "caching release artifact",
    "deploy": "branch=preview",
    "docker": "kayron",
    "github": "kayron",
    "preview": "false",
    "caller": "/Users/xh3b4sd/project/0xSplits/kayron/pkg/cache/create.go:17"
}
{
    "time": "2025-09-17 15:32:02",
    "level": "debug",
    "message": "caching release artifact",
    "deploy": "release=v0.2.2",
    "docker": "specta",
    "github": "specta",
    "preview": "false",
    "caller": "/Users/xh3b4sd/project/0xSplits/kayron/pkg/cache/create.go:17"
}
{
    "time": "2025-09-17 15:32:02",
    "level": "debug",
    "message": "instrumented worker handler",
    "handler": "release",
    "latency": "3.727451s",
    "success": "true",
    "caller": "/Users/xh3b4sd/go/pkg/mod/github.com/0x!splits/workit@v0.6.0/handler/metrics/ensure.go:55"
}
{
    "time": "2025-09-17 15:32:02",
    "level": "debug",
    "message": "caching current state",
    "current": "d5fa88afd502edc9052f89c956618b2cb567d984",
    "github": "infrastructure",
    "caller": "/Users/xh3b4sd/project/0xSplits/kayron/pkg/operator/template/ensure.go:35"
}
{
    "time": "2025-09-17 15:32:02",
    "level": "debug",
    "message": "instrumented worker handler",
    "handler": "template",
    "latency": "173.083µs",
    "success": "true",
    "caller": "/Users/xh3b4sd/go/pkg/mod/github.com/0x!splits/workit@v0.6.0/handler/metrics/ensure.go:55"
}
{
    "time": "2025-09-17 15:32:02",
    "level": "debug",
    "message": "caching desired state",
    "desired": "v0.1.1",
    "github": "splits-lite",
    "preview": "false",
    "caller": "/Users/xh3b4sd/project/0xSplits/kayron/pkg/operator/reference/ensure.go:45"
}
{
    "time": "2025-09-17 15:32:02",
    "level": "debug",
    "message": "caching desired state",
    "desired": "v0.2.2",
    "github": "specta",
    "preview": "false",
    "caller": "/Users/xh3b4sd/project/0xSplits/kayron/pkg/operator/reference/ensure.go:45"
}
{
    "time": "2025-09-17 15:32:03",
    "level": "debug",
    "message": "caching desired state",
    "desired": "d5fa88afd502edc9052f89c956618b2cb567d984",
    "github": "infrastructure",
    "preview": "false",
    "caller": "/Users/xh3b4sd/project/0xSplits/kayron/pkg/operator/reference/ensure.go:45"
}
{
    "time": "2025-09-17 15:32:03",
    "level": "debug",
    "message": "caching desired state",
    "desired": "eb9f56e195f25218499897a6c6d79c62261faba5",
    "github": "kayron",
    "preview": "false",
    "caller": "/Users/xh3b4sd/project/0xSplits/kayron/pkg/operator/reference/ensure.go:45"
}
{
    "time": "2025-09-17 15:32:03",
    "level": "debug",
    "message": "caching desired state",
    "desired": "e307253e62c2da119579e2505db0b6185642ab9b",
    "github": "splits-lite",
    "preview": "true",
    "caller": "/Users/xh3b4sd/project/0xSplits/kayron/pkg/operator/reference/ensure.go:45"
}
{
    "time": "2025-09-17 15:32:03",
    "level": "debug",
    "message": "instrumented worker handler",
    "handler": "reference",
    "latency": "367.2485ms",
    "success": "true",
    "caller": "/Users/xh3b4sd/go/pkg/mod/github.com/0x!splits/workit@v0.6.0/handler/metrics/ensure.go:55"
}
{
    "time": "2025-09-17 15:32:04",
    "level": "debug",
    "message": "caching current state",
    "current": "bc7891268e44f62e0aebbe339c0850b61d52c417",
    "docker": "splits-lite",
    "preview": "false",
    "caller": "/Users/xh3b4sd/project/0xSplits/kayron/pkg/operator/container/cache.go:10"
}
{
    "time": "2025-09-17 15:32:04",
    "level": "debug",
    "message": "caching current state",
    "current": "bc7891268e44f62e0aebbe339c0850b61d52c417",
    "docker": "splits-lite",
    "preview": "true",
    "caller": "/Users/xh3b4sd/project/0xSplits/kayron/pkg/operator/container/cache.go:10"
}
{
    "time": "2025-09-17 15:32:04",
    "level": "debug",
    "message": "caching current state",
    "current": "eb9f56e195f25218499897a6c6d79c62261faba5",
    "docker": "kayron",
    "preview": "false",
    "caller": "/Users/xh3b4sd/project/0xSplits/kayron/pkg/operator/container/cache.go:10"
}
{
    "time": "2025-09-17 15:32:04",
    "level": "debug",
    "message": "caching current state",
    "current": "v0.2.2",
    "docker": "specta",
    "preview": "false",
    "caller": "/Users/xh3b4sd/project/0xSplits/kayron/pkg/operator/container/cache.go:10"
}
{
    "time": "2025-09-17 15:32:04",
    "level": "debug",
    "message": "instrumented worker handler",
    "handler": "container",
    "latency": "1.705606292s",
    "success": "true",
    "caller": "/Users/xh3b4sd/go/pkg/mod/github.com/0x!splits/workit@v0.6.0/handler/metrics/ensure.go:55"
}
{
    "time": "2025-09-17 15:32:05",
    "level": "debug",
    "message": "executed image check",
    "exists": "true",
    "image": "splits-lite",
    "preview": "true",
    "tag": "e307253e62c2da119579e2505db0b6185642ab9b",
    "caller": "/Users/xh3b4sd/project/0xSplits/kayron/pkg/operator/registry/ensure.go:45"
}
{
    "time": "2025-09-17 15:32:05",
    "level": "debug",
    "message": "executed image check",
    "exists": "true",
    "image": "splits-lite",
    "preview": "false",
    "tag": "v0.1.1",
    "caller": "/Users/xh3b4sd/project/0xSplits/kayron/pkg/operator/registry/ensure.go:45"
}
{
    "time": "2025-09-17 15:32:05",
    "level": "debug",
    "message": "instrumented worker handler",
    "handler": "registry",
    "latency": "875.278125ms",
    "success": "true",
    "caller": "/Users/xh3b4sd/go/pkg/mod/github.com/0x!splits/workit@v0.6.0/handler/metrics/ensure.go:55"
}
{
    "time": "2025-09-17 15:32:05",
    "level": "info",
    "message": "continuing reconciliation loop",
    "reason": "detected state drift",
    "release": "splits-lite",
    "version": "v0.1.1",
    "caller": "/Users/xh3b4sd/project/0xSplits/kayron/pkg/operator/policy/ensure.go:68"
}
{
    "time": "2025-09-17 15:32:05",
    "level": "info",
    "message": "continuing reconciliation loop",
    "domain": "1d0fd508.lite.testing.splits.org",
    "reason": "detected state drift",
    "release": "splits-lite",
    "version": "e307253e62c2da119579e2505db0b6185642ab9b",
    "caller": "/Users/xh3b4sd/project/0xSplits/kayron/pkg/operator/policy/ensure.go:68"
}
{
    "time": "2025-09-17 15:32:05",
    "level": "debug",
    "message": "instrumented worker handler",
    "handler": "policy",
    "latency": "523.166µs",
    "success": "true",
    "caller": "/Users/xh3b4sd/go/pkg/mod/github.com/0x!splits/workit@v0.6.0/handler/metrics/ensure.go:55"
}
{
    "time": "2025-09-17 15:32:05",
    "level": "debug",
    "message": "resolved ref for github repository",
    "environment": "testing",
    "ref": "d5fa88afd502edc9052f89c956618b2cb567d984",
    "repository": "https://github.com/0xSplits/infrastructure",
    "caller": "/Users/xh3b4sd/project/0xSplits/kayron/pkg/operator/infrastructure/ensure.go:24"
}
{
    "time": "2025-09-17 15:32:05",
    "level": "debug",
    "message": "uploading cloudformation template",
    "bucket": "splits-cf-templates",
    "key": "testing/deployment/deployment.yaml",
    "caller": "/Users/xh3b4sd/project/0xSplits/kayron/pkg/operator/infrastructure/aws.go:26"
}
{
    "time": "2025-09-17 15:32:06",
    "level": "debug",
    "message": "uploading cloudformation template",
    "bucket": "splits-cf-templates",
    "key": "testing/discovery/discovery.yaml",
    "caller": "/Users/xh3b4sd/project/0xSplits/kayron/pkg/operator/infrastructure/aws.go:26"
}
{
    "time": "2025-09-17 15:32:06",
    "level": "debug",
    "message": "uploading cloudformation template",
    "bucket": "splits-cf-templates",
    "key": "testing/elasticache/elasticache.yaml",
    "caller": "/Users/xh3b4sd/project/0xSplits/kayron/pkg/operator/infrastructure/aws.go:26"
}
{
    "time": "2025-09-17 15:32:06",
    "level": "debug",
    "message": "uploading cloudformation template",
    "bucket": "splits-cf-templates",
    "key": "testing/fargate/fargate.yaml",
    "caller": "/Users/xh3b4sd/project/0xSplits/kayron/pkg/operator/infrastructure/aws.go:26"
}
{
    "time": "2025-09-17 15:32:06",
    "level": "debug",
    "message": "uploading cloudformation template",
    "bucket": "splits-cf-templates",
    "key": "testing/index.yaml",
    "caller": "/Users/xh3b4sd/project/0xSplits/kayron/pkg/operator/infrastructure/aws.go:26"
}
{
    "time": "2025-09-17 15:32:07",
    "level": "debug",
    "message": "uploading cloudformation template",
    "bucket": "splits-cf-templates",
    "key": "testing/kayron/kayron.yaml",
    "caller": "/Users/xh3b4sd/project/0xSplits/kayron/pkg/operator/infrastructure/aws.go:26"
}
{
    "time": "2025-09-17 15:32:07",
    "level": "debug",
    "message": "uploading cloudformation template",
    "bucket": "splits-cf-templates",
    "key": "testing/rds/rds.alarms.yaml",
    "caller": "/Users/xh3b4sd/project/0xSplits/kayron/pkg/operator/infrastructure/aws.go:26"
}
{
    "time": "2025-09-17 15:32:07",
    "level": "debug",
    "message": "uploading cloudformation template",
    "bucket": "splits-cf-templates",
    "key": "testing/rds/rds.yaml",
    "caller": "/Users/xh3b4sd/project/0xSplits/kayron/pkg/operator/infrastructure/aws.go:26"
}
{
    "time": "2025-09-17 15:32:07",
    "level": "debug",
    "message": "uploading cloudformation template",
    "bucket": "splits-cf-templates",
    "key": "testing/server/server.yaml",
    "caller": "/Users/xh3b4sd/project/0xSplits/kayron/pkg/operator/infrastructure/aws.go:26"
}
{
    "time": "2025-09-17 15:32:08",
    "level": "debug",
    "message": "uploading cloudformation template",
    "bucket": "splits-cf-templates",
    "key": "testing/specta/specta.yaml",
    "caller": "/Users/xh3b4sd/project/0xSplits/kayron/pkg/operator/infrastructure/aws.go:26"
}
{
    "time": "2025-09-17 15:32:08",
    "level": "debug",
    "message": "uploading cloudformation template",
    "bucket": "splits-cf-templates",
    "key": "testing/splits-lite/splits-lite.yaml",
    "caller": "/Users/xh3b4sd/project/0xSplits/kayron/pkg/operator/infrastructure/aws.go:26"
}
{
    "time": "2025-09-17 15:32:08",
    "level": "debug",
    "message": "uploading cloudformation template",
    "bucket": "splits-cf-templates",
    "key": "testing/telemetry/telemetry.yaml",
    "caller": "/Users/xh3b4sd/project/0xSplits/kayron/pkg/operator/infrastructure/aws.go:26"
}
{
    "time": "2025-09-17 15:32:08",
    "level": "debug",
    "message": "uploading cloudformation template",
    "bucket": "splits-cf-templates",
    "key": "testing/vpc/vpc.yaml",
    "caller": "/Users/xh3b4sd/project/0xSplits/kayron/pkg/operator/infrastructure/aws.go:26"
}
{
    "time": "2025-09-17 15:32:08",
    "level": "debug",
    "message": "instrumented worker handler",
    "handler": "infrastructure",
    "latency": "3.615857209s",
    "success": "true",
    "caller": "/Users/xh3b4sd/go/pkg/mod/github.com/0x!splits/workit@v0.6.0/handler/metrics/ensure.go:55"
}
{
    "time": "2025-09-17 15:32:08",
    "level": "info",
    "message": "updating cloudformation stack",
    "name": "server-test",
    "url": "https://splits-cf-templates.s3.us-west-2.amazonaws.com/testing/index.yaml",
    "caller": "/Users/xh3b4sd/project/0xSplits/kayron/pkg/operator/cloudformation/ensure.go:30"
}
{
    "time": "2025-09-17 15:32:08",
    "level": "debug",
    "message": "instrumented worker handler",
    "handler": "cloudformation",
    "latency": "190.125µs",
    "success": "true",
    "caller": "/Users/xh3b4sd/go/pkg/mod/github.com/0x!splits/workit@v0.6.0/handler/metrics/ensure.go:55"
}
--- PASS: Test_Operator_Integration (9.93s)
PASS
ok  	github.com/0xSplits/kayron/pkg/operator	11.224s
```

### Releases

In order to update the Docker image, prepare all desired changes within the
`main` branch and create a Github release for the desired Kayron version. The
release tag should be in [Semver Format]. Creating the Github release triggers
the responsible [Github Action] to build and push the Docker image to the
configured [Amazon ECR].

```
v0.1.11
```

The version command `kayron version` and the version endpoint `/version` provide
build specific version information about the build and runtime environment. A
live demo can be seen at https://kayron.testing.splits.org/version.

# Docker

Kayron's build artifact is a statically compiled binary running in a
[distroless] image for maximum security and minimum size. If you do not have Go
installed and just want to run Kayron locally in a Docker container, then use
the following commands.

```
docker build \
  --build-arg SHA="local-test-sha" \
  --build-arg TAG="local-test-tag" \
  -t kayron:local .
```

```
docker run \
  -e KAYRON_ENVIRONMENT=development \
  -p 7777:7777 \
  kayron:local \
  daemon
```

[Amazon ECR]: https://docs.aws.amazon.com/ecr
[AmazonEC2ContainerRegistryReadOnly]: https://docs.aws.amazon.com/aws-managed-policy/latest/reference/AmazonEC2ContainerRegistryReadOnly.html
[Cobra]: https://github.com/spf13/cobra
[custom task engine]: https://github.com/0xSplits/workit
[distroless]: https://github.com/GoogleContainerTools/distroless
[Github Action]: .github/workflows/docker-release.yaml
[operator pattern]: https://kubernetes.io/docs/concepts/extend-kubernetes/operator
[Kubernetes]: https://kubernetes.io/docs/concepts/overview
[release artifacts]: https://github.com/0xSplits/releases
[ResourceGroupsandTagEditorReadOnlyAccess]: https://docs.aws.amazon.com/aws-managed-policy/latest/reference/ResourceGroupsandTagEditorReadOnlyAccess.html
[Semver Format]: https://semver.org
[ViewOnlyAccess]: https://docs.aws.amazon.com/aws-managed-policy/latest/reference/ViewOnlyAccess.html
