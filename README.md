<p align="center"><img width="200" src="https://avatars0.githubusercontent.com/u/65870207" alt="Jobatator logo"></a></p>

<p align="center">
  <a href="https://github.com/jobatator/jobatator/actions"><img src="https://github.com/jobatator/jobatator/workflows/Continuous%20integration/badge.svg" alt="Build Status"></a>
  <a href="https://discord.gg/9M4vVsX"><img src="https://img.shields.io/badge/chat-on%20discord-7289da.svg?sanitize=true" alt="Chat"></a>
  <br>
</p>

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

- AUTH {username} {password}
- USE_GROUP {group}
- PUBLISH {queue_slug} {job_type} 'payload'
- SUBSCRIBE {queue_slug}
- UPDATE_JOB {job_id} {job_status} # status can be 'done', 'in-progress' or 'errored'

You can find description of all the commands of the jobatator server [here](https://github.com/lefuturiste/jobatator/blob/master/pkg/commands/commands.go)

## Internal Entities/Ressources

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

## Roadmap/Todolist

- Being able to mock commands and not necessarly use the socket interface (can be usefull for the http gateway interface but also for unit tests)
    - Should we use a CmdOutput struct type?
    - Should we use a (string, error) return tuple type?
- User that have * groups can access all groups and are considered as an administrator
- Get stat about a queue (number of job of certain state and type)
    - By job.type and By all:
        - How many jobs are pending?
        - How many jobs are in-progress
        - How many jobs are errored?
        - How many jobs are done?
- Refactor all the data management part and come up with a mini librairy to use a relational database in memory
    - Organization by tables, fields, relation ship
    - Only code for now the hasMany() and belongsTo() relation ship
- Simple web interface which is using the gateway HTTP interface to admistrate
   - What UI/JS Framework we want to use?
   - Should we separate this web interface in a different repository? 
   - Log in with username and password
   - Show some stats about the server
   - Easily switch between groups
   - See the queues in that group
   - See the jobs/recurrent jobs in the queue
   - See workers in a queue
   - Kick out a worker
   - Delete a job
   - Delete a queue