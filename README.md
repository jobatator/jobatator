# Jobatator

Jobatator is a light alternative to RabbitMQ, if you need a way to connect workers together in order to work on some heavy jobs, this might be a solution for you. You can interact with jobatator using a TCP connexion.

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