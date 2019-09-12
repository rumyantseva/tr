# tr

A tool to find out the addresses of all intermediate hops between your machine and any target host.

## How to use?

## Implementation notes



### Task

- Find out the addresses of all intermediate hops between your machine and
any target host (for instance, www.google.com).
- Find the largest difference in response time between consecutive hops and output it separately.

### Idea



In principle, multiple implementations are possible: we can use ICMP, UDP, TCP or other protocols.

For simplicity, let's choose ICMP. It won't let us specify port for the target but as
only host is mentioned in the original task, ICMP should be enough.


### Technical Decisions

- CLI-related libraries: the standard `flag` library is chosen to not to overcomplicate the tool.


### Dependency Management


### Testing


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

### Release

In addition, when the tag is pushed, Github Actions will call the Release phase of the pipeline.
The binaries for linux-amd64 and darwin-amd64 will be prepared and uploaded as
the [latest release](https://github.com/rumyantseva/tr/releases).
