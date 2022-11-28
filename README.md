### Start

# build and push docker image
1. docker build -t docker-first .
2. docker image tag docker-first imigaka/docker-first:latest
3. docker push imigaka/docker-first:latest

# init docker swarm
1. docker swarm init
# start docker service
1. docker stack deploy -c docker-compose.yml docker-first
# check go service logs
1. docker service logs -f docker-first_go-service

## Kill
1. docker stack rm docker-first
2. docker swarm leave --force

## Set name
In docker-compose.yml in service go-service in environment set NAME to you name