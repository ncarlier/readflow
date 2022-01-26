#########################################
# Build stage
#########################################
FROM golang:1.17 AS builder

# Repository location
ARG REPOSITORY=github.com/ncarlier

# Artifact name
ARG ARTIFACT=readflow

# Copy sources into the container
ADD . /go/src/$REPOSITORY/$ARTIFACT

# Set working directory
WORKDIR /go/src/$REPOSITORY/$ARTIFACT

# Build the binary
RUN make

#########################################
# Distribution stage
#########################################
FROM gcr.io/distroless/base-debian11

# Repository location
ARG REPOSITORY=github.com/ncarlier

# Artifact name
ARG ARTIFACT=readflow

# Install binary
COPY --from=builder /go/src/$REPOSITORY/$ARTIFACT/release/$ARTIFACT /usr/local/bin/$ARTIFACT

# Add configuration file
ADD ./pkg/config/readflow.toml /etc/readflow.toml

# Set configuration file
ENV READFLOW_CONFIG /etc/readflow.toml

# Exposed ports
EXPOSE 8080 9090

# Define entrypoint
ENTRYPOINT [ "readflow" ]
