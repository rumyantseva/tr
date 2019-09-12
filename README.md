# tr

A tool to find out the addresses of all intermediate hops between your machine and any target host.

### How to use?

Dockerized version:

```
docker build -t tr -f Dockerfile .
docker run tr -maxTTL 5 google.com
```

Compiled version:

```
GOOS=darwin GOARCH=amd64 make build
sudo ./bin/darwin-amd64/tr -maxTTL 5 google.com
```

Available flags:

```

```

### Task

- Find out the addresses of all intermediate hops between your machine and
any target host (for instance, www.google.com).
- Find the largest difference in response time between consecutive hops and output it separately.

### Idea and implementation details

#### Intermediate hops
In principle, multiple implementations are possible: we can use ICMP, UDP, TCP or other protocols.

For simplicity, let's choose ICMP. It won't let us specify port for the target but as
only host is mentioned in the original task, ICMP should be enough.

In implementation, we can send ICMP Echo Request and will wait for ICMP Echo Reply.
We start from TTL=1 and for every new request we increase TTL.
Until we reach the target, we expect to receive the Time Exceeded error instead of Echo Reply.

#### Response time

In fact, the network situation is changing quite often. So that, the response time is not a really
reliable metric in this particular case. Due to continuous changes in the network,
it's also quite expected that we can complete a round trip for a bigger TTL
faster than a round trip for a smaller TTL.


In this implementation, we do requests hops for every TTL value only once.
In principle, it might make sense to do such requests a few times to have a better representation for a median
response time.
A few attempts also might be useful if we haven't reached the desired hop from the first attempt.

However, as it was requested in the task, the largest difference in response time between consecutive hops is calculated.
As timeouts are potentially possible, we will calculate the largest difference only for hops where we don't have timeouts.

### Go-related Decisions

- The `pkg/trace` package is designed to be used as a library, it might be included into external projects.
- CLI-related libraries: the standard `flag` library is chosen to not to overcomplicate the tool.
- In addition to the standard library, `golang.org/x/net` is used to work with IPv4, IPv6 and ICMP.

### Dependency Management

Dependency management is provided with Go Modules.
[GoCenter](https://gocenter.io) is used as `GOPROXY` in `Makefile` and `Dockerfile`.

Typically, for production-readiness I still use vendor-mode in addition.
For this particular task, Go Modules + `GOPROXY` should be enough.

### Testing

Proper testing requires quite a lot of mock on network
(especially, if we want to test our tool against possible network states).
As time for the coding challenge is limited, unit / integration testing is our of scope now.

### Continuous Integration

I used this coding challenge as an opportunity to try Github Actions for the first time.
The implementation is very MVP'ish but it works :)

### Checks & tests

The idea is that the main CI-related checks are defined in `Dockerfile.test`.
During the `docker build` phase, linters ([golangci-lint](https://github.com/golangci/golangci-lint) configured
with `.golangci.yml`) and tests will be run.
In addition, compilation is also included into this phase, so it's checked that the binary is compilable and callable.

This dockerized approach will work with any CI tool (or even locally) if Docker is available there.

How to run: `docker build -f Dockerfile.test .`

The main Dockerfile to compile the

### Release

In addition, when the tag is pushed, Github Actions will call the Release phase of the pipeline.
The binaries for linux-amd64 and darwin-amd64 will be prepared and uploaded as
the assets of [latest release](https://github.com/rumyantseva/tr/releases).
