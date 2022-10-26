# Quickstart

Build the image:
```
docker build -t priobike-biker-swarm-test-client --no-cache .
```

Deploy to stack:
```
docker stack deploy --compose-file docker-compose.yml biker-swarm
```

Check if it is running:
```
docker stack services biker-swarm
```

Remove the deployment:
```
docker stack rm biker-swarm
```