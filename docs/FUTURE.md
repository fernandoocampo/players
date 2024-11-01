# Future

How would you expand the solution if you would have spent more time on it?

## Opportunities for improvement

1. Add CI/CD pipelines to improve the feedback flow and discover hidden bugs. More quality gates, more velocity.
2. Improve continuous delivery by adding [semantic realease](https://github.com/semantic-release/semantic-release) steps based on [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/).
3. Add a commit linter inside the CI pipeline to ensure the semantic publishing process.
4. Add telemetry to have traces and metrics that will help maintainers identify bottlenecks, most demanded logic and performance.
5. Add authorization capabilities with OAuth2 to improve security and do so via a sidecar to avoid adding non-player-related logic to the service. e.g. oauth2-proxy project.
6. Improve data protection features to accomplish with international regulations. Encrypting data at database and transport layer.
7. Enable grpc ssl server to improve transport security.
8. Add a RESTful API to support clients that cannot make RPC calls.
9. Add Go vulnerability checker to verify that code and dependencies "are free" of security holes.
10. Add pipeline to version and generate new protobuffers for players grpc service.
11. Load tests and profiling to see performance and memory/cpu leaks.
12. Add more validations for value lengths.
13. Add goreleaser to build binaries.
14. If the service will run in a kubernetes cluster environment, helm charts could be included.
15. Use a real eventbus to send player events to well known topic systems. e.g. sns, kinesis, eventbridge, kafka, etc.
16. Get secret values from a secret manager.
17. Add benchmarks for critical tasks.
18. Add plantuml diagrams to show the design of the microservice.
19. Add test suite to unit tests to reduce boilerplate.
20. Add logic to disable and enable players.

## Technical Debts

1. I added an endpoint layer inside the player package to support multiple transports, http and grpc, but I decided to implement only one grpc server and connect it directly to the service layer. "if in the future the project decides it needs a RESTful api, It is advisable to use endpoint instead of service", Here I am breaking the YAGNI principle.