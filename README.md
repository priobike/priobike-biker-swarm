# priobike-biker-swarm

This service can be used to simulate the behavior of a client against the backend. When scaling horizontally, the service can be used to make load tests.

Using in PrioBike, this service mimics the behavior of the PrioBike app. We use this to test our backend.

It depends on the backend services it is running against. If you want to report/save the results of the biker swarm, the [priobike-biker-swarm-monitor](https://github.com/priobike/priobike-biker-swarm-monitor) needs to run as well.

[Learn more about PrioBike](https://github.com/priobike)

We are using this as part of as Docker Swarm stack to scale the service horizontally across multiple nodes: https://github.com/priobike/priobike-biker-swarm-deployment

## Quickstart

The easiest way to run this is by using the included Docker Compose file.

```bash
docker-compose up
```

Depending on the hardware resources of the machine, you are using, you might want to adjust the number of replicas in the `docker-compose.yml` file.

There are the following environment variables that can be set:

- `TIMEOUT`: Sets after how many seconds a request should timeout.
- `DEPLOYMENT`: Sets against which backend deployment the service should run. The deployments are configured in the code.
- `REPORT_RESULTS`: Sets whether the results should be reported to the backend. This is useful when you want to analyze how long the requests took overall and what services timed out or failed. To report results, the [priobike-biker-swarm-monitor](https://github.com/priobike/priobike-biker-swarm-monitor) needs to run.

## CLI

Build the Go binary:

```bash
go build -o main .
```

Run the binary:

```bash
./main
```

## Contributing

We highly encourage you to open an issue or a pull request. You can also use our repository freely with the `MIT` license. 

Every service runs through testing before it is deployed in our release setup. Read more in our [PrioBike deployment readme](https://github.com/priobike/.github/blob/main/wiki/deployment.md) to understand how specific branches/tags are deployed.

## Anything unclear?

Help us improve this documentation. If you have any problems or unclarities, feel free to open an issue.
