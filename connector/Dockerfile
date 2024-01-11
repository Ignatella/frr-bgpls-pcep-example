FROM golang:1.19

# Set destination for COPY
WORKDIR /connector

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY . ./
# Download Go modules
RUN go mod download

# Build
WORKDIR cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o bgpd

EXPOSE 1790

WORKDIR /connector

# Run
CMD ["./cmd/bgpd"]
