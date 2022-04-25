# Quickstart

Build the image:
```
docker build -t priobike-biker-swarm-test-client-staging --no-cache ./test-client
```

Deploy to stack:
```
docker stack deploy --compose-file docker-compose.yml biker-swarm-staging
```

Check if it is running:
```
docker stack services biker-swarm-staging
```

Remove the deployment:
```
docker stack rm biker-swarm-staging
```