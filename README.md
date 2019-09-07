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

