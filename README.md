# Swarm Service Handler

This is a Go script that allows you to filter, stop, start, and restart services in a Docker Swarm mode cluster.

## Prerequisites
- Docker installed and running in Swarm mode
- Go installed and configured

## Usage

1. Clone the repository
git clone https://github.com/klarekante/swarm-service-handler.git
2. Build the script
> go build -o swarm-service-handler
3. Run the script with command line arguments
> ./swarm-service-handler --help

You can pass command line arguments using the flag package, for example:
> swarm-service-handler --stop --services service1,service2
