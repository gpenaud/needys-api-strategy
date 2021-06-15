# needys-api-strategy
An API micro-service for needys application to manage "need" objects.
On "POST" curl request, the api send a message through a rabbitmq broker
to trigger a pdf generation in another service: needys-output-producer

Please contact <guillaume.penaud@gmail.com> if you have questions !

##### manage the compose stack

```
### to start it
docker-compose up --detach
OR
make start

### to stop it
docker-compose down
OR
make stop
```

##### to start only a part of the stack
```
### only the api part
docker-compose up needys-api-strategy
OR
make api-only

### only the sidecars part
docker-compose up mariadb rabbitmq
OR
make sidecars-only
```

##### possible tests for the api
Thoses tests are covering every usage of the api needys-api-strategy. Use them
to validate both mysql and rabbitmq usage

```
### list database content
curl -X GET http://localhost:8010
OR
make test-list

### delete entries matching "testing-need" as name, then insert one into database
curl -X DELETE http://localhost:8010?name=testing-need
curl -d "name=testing-need&priority=high" -X POST http://localhost:8010
OR
make test-all

### only delete entries matching "testing-need" as name from database
curl -X DELETE http://localhost:8010?name=testing-need
OR
make test-delete

### only insert one entry "testing-need" in database
curl -d "name=testing-need&priority=high" -X POST http://localhost:8010
OR
test-insert
```

### Tricks for container debug

# CMD exec /bin/bash -c "trap : TERM INT; sleep infinity & wait"