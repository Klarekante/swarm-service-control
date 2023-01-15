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


```
usage: service-manager [-services string] [-restart] [-stop] [-backup] [-restore] [-all]

Service Manager is a script that allows you to perform various operations on Docker services deployed in a Swarm mode cluster.

optional arguments:
  -services string
        Comma separated names of services
  -restart
        Restart the specified services
  -stop
        Stop the specified services
  -backup
        Backup running services
  -restore
        Restore services from backup
  -all
        Operate on all running services
```
