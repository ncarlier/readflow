#########################################
# Build frontend stage
#########################################
FROM node:lts-alpine AS frontend-builder

# Setup env
RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app
ENV PATH /usr/src/app/node_modules/.bin:$PATH

# Install dependencies
COPY ui/package.json /usr/src/app/package.json
COPY ui/package-lock.json /usr/src/app/package-lock.json
RUN npm install --silent --legacy-peer-deps

# Build website
COPY ./ui /usr/src/app
RUN npm run build

#########################################
# Build backend stage
#########################################
FROM golang:1.19 AS backend-builder

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

# Install backend binary
COPY --from=backend-builder /go/src/$REPOSITORY/$ARTIFACT/release/$ARTIFACT /usr/local/bin/$ARTIFACT
# Install frontend assets
COPY --from=frontend-builder /usr/src/app/build /var/local/html

# Add configuration file
ADD ./internal/config/defaults.toml /etc/readflow.toml

# Set configuration file
ENV READFLOW_CONFIG /etc/readflow.toml

# Serve UI
ENV READFLOW_UI_DIRECTORY /var/local/html

# Exposed ports
EXPOSE 8080 9090

# Define entrypoint
ENTRYPOINT [ "readflow" ]

# Define command
CMD [ "serve" ]
