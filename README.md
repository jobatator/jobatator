# Jobatator

Jobatator is a light alternative to RabbitMQ, if you need a way to connect workers together in order to work on some 
heavy jobs, this might be a solution for you. You can interact with jobatator using a TCP connexion.

## Configuration

You can find an example of configuration in `config.example.yml`.

**Warning: by default if you are not providing the port or host key in the config, the server will listen by default on host 0.0.0.0 and on port 8962.**

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