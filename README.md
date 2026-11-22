# heliosphere
Rate limiting for the system and it's applications

## Quickstart

1. Clone this repository

2. Run the following command

```
docker-compose up
```

## Endpoints

- GET http://localhost:8090/ratelimiter - retrieves all rate limiters
- GET http://localhost:8090/ratelimiter/{id} - retrieves a rate limiter by its id
- POST http://localhost:8090/ratelimiter - create a new rate limiter
- PUT http://localhost:8090/ratelimiter - update a rate limiter
- DELETE http://localhost:8090/ratelimiter/{id} - delete a rate limiter
- GET http://localhost:8090/ratelimiter/allow/{id} - checks if a rate limiter can allow request 
## Configuration

The gateway auto-creates a limiter (token bucket) the first time a
user/method key is seen. Its size comes from the rate-limiter service's
environment:

| Variable | Default | Meaning |
|---|---|---|
| `RL_DEFAULT_REQ_PER_SEC` | `1` | refill rate of an auto-created limiter |
| `RL_DEFAULT_BURST` | `1` | bucket size of an auto-created limiter |

Limiters are stored in Consul, so a changed default applies only to keys
created afterwards; delete the stored ones (`consul kv delete -recurse
rateLimiters/`) to pick up the new value.
