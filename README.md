# Quickstart

Build the image:
```
docker build -t bikenow.vkw.tu-dresden.de/priobike-biker-swarm:main --no-cache .
```

Deploy to stack:
```
docker stack deploy \
    --prune \
    --with-registry-auth \
    --compose-file docker-compose.yml \
    biker-swarm
```

Check if it is running:
```
docker stack services biker-swarm
```

Remove the deployment:
```
docker stack rm biker-swarm
```