install:
    @echo "Installing dependencies..."
    @cp .env.dist .env
    @go version
    @go mod download
    @docker-compose build

run:
    @echo "Running application..."
    @docker-compose up -d

stop:
    @echo "Stopping application..."
    @docker-compose stop

clean-logs:
    @rm -rf app/api/logs/*.log
    @rm -rf app/indexer/logs/*.log

nuke:
    @echo "Nuking application..."
    @docker-compose down
    @docker rmi -f tzkt-delegator-tools-indexer
    @docker rmi -f tzkt-delegator-tools-api