# Jobatator

## Commands

- PING
- AUTH username:password
- PUBLISH namespace queue job_type 'payload'
- SUBSCRIBE namespace queue
- JOB UPDATE job_id status # status can be 'done' 'processing' or 'failed'

CRUD on namespaces
- NAMESPACE LIST
- NAMESPACE CREATE
- NAMESPACE UPDATE
- NAMESPACE DELETE

CRUD on users
- USER LIST
- USER CREATE
- USER UPDATE
- USER DELETENew client: 127.0.0.1:57264
New cmd:  map[0:EVAL 1:-- Get all of the messages with an expired "score"...
local val = redis.call('zrangebyscore', KEYS[1], '-inf', ARGV[1])

-- If we have values in the array, we will remove them from the first queue
-- and add them onto the destination queue in chunks of 100, which moves
-- all of the appropriate messages onto the destination queue very safely.
if(next(val) ~= nil) then
    redis.call('zremrangebyrank', KEYS[1], 0, #val - 1)

    for i = 1, #val, 100 do
        redis.call('lpush', KEYS[2], unpack(val, i, math.min(i+99, #val)))
    end


Si on veut utiliser keyvaluer comme jobatator, alors il faut ajouter ces systèmes:
- on va alors avoir un serveur keyvaluer pour chaque namespace
- il faut un moyen afin de faire en sorte que le worker reçoit la notification qu'il y a un job 
- pouvoir dispatcher les jobs, et ça c'est compliqué

## Ressources

### queue

- slug: string
- jobs: collection

### job

- slug: string
- status: string  ['waiting', 'processing', 'done']