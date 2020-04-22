# Jobatator

Jobatator is a light alternative to RabbitMQ, if you need a way to connect workers together in order to work on some 
heavy jobs, this might be a solution for you. You can interact with jobatator using a TCP connexion.

## Configuration

You can find an example of configuration in `config.example.yml`. Access it [here](https://github.com/lefuturiste/jobatator/blob/master/config.example.yml).

**Warning: by default if you are not providing the port or host key in the config, the server will listen by default on host 0.0.0.0 and on port 8962.**

## Docker image

You can find the docker image details [here](https://hub.docker.com/repository/docker/lefuturiste/jobatator).

## Use with docker-compose

The most simple way to start a jobatator server using docker is by making use of a `docker-compose.yml` file.

You can find an example of a `docker-compose.yml` file to host jobatator [here](https://github.com/lefuturiste/jobatator/blob/master/docker-compose.yml).

In the same directory with the `docker-compose.yml` file type this command to start the server:

`docker-compose up -d`

### Use with docker cli

Not really recommanded, but you can do it:

```bash
# in the publish or -p flag, the right port is the container port whereas the left part is the machine port \
docker run \
    --name jobatator \
    -p 8962:8962 \
    -v /absolute/path/to/config.yml:/go/src/app/config.yml \
    lefuturiste/jobatator:latest
```

## Commands

Warning: For now the server is using "\r\n" but "\n" as end line!

The major commands:

- AUTH username password
- PING
- USE_GROUP group
- PUBLISH namespace queue job_type 'payload'
- SUBSCRIBE namespace queue
- JOB UPDATE job_id status # status can be 'done' 'in-progress' or 'errored'

## Ressources

### Group

- Slug: string

### User

- Username:     string
- Password:     string
- Groups:       string[]
- Addr:         string
- CurrentGroup: Group
- Conn:         net.Conn
- Status:       string

### Queue

- Group:   string
- Slug:    string
- Jobs:    Job[]
- Workers: User[]

### Job

- ID:                  string
- Slug:                string
- Type:                string
- Payload:             string
- Status:              string  ['pending', 'in-progress', 'errored', 'done']
- Attempts:            int
- StartedProcessingAt: timestamp
- EndProcessingAt:     timestamp