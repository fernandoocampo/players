FROM golang:1.23.2-alpine3.19 as builder

ARG appVersion
ARG commitHash
ARG buildDate

ENV VERSION=$appVersion
ENV COMMIT_HASH=$commitHash
ENV BUILD_DATE=$buildDate

WORKDIR /app
COPY . ./
RUN apk add --no-cache make

# Build
RUN go mod tidy
RUN make build-linux

# Run stage
FROM alpine:3.19.2
RUN apk --no-cache add ca-certificates
COPY --from=0 /app/bin/players-service-linux ./

EXPOSE 8080
EXPOSE 50051

CMD [ "/players-service-linux" ]
