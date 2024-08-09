# Use the official Jenkins agent image as the base image
FROM jenkins/inbound-agent:latest

# Set environment variables for Go
ENV GO_VERSION=1.22.6
ENV GOLANG_DOWNLOAD_URL=https://golang.org/dl/go$GO_VERSION.linux-amd64.tar.gz
ENV GOLANG_DOWNLOAD_SHA256=999805bed7d9039ec3da1a53bfbcafc13e367da52aa823cb60b68ba22d44c616
# Switch to root user
USER root
# Install Go
RUN curl -fsSL "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz \
    && echo "$GOLANG_DOWNLOAD_SHA256  golang.tar.gz" | sha256sum -c - \
    && tar -C /usr/local -xzf golang.tar.gz \
    && rm golang.tar.gz

# Add Go binary to the PATH
ENV PATH="/usr/local/go/bin:${PATH}"

# Verify Go installation
RUN go version

# Set the Jenkins agent work directory
WORKDIR /home/jenkins/agent

# Start the Jenkins agent
ENTRYPOINT ["/usr/local/bin/jenkins-agent"]
