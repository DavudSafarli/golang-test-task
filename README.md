### Project

There are 3 services defined in docker-compose file next to rabbitmq and kafka:

- api: http server running on 8080
- reporting: http server running on 3000
- processor: MessageProcessor that consumes messages from rabbitmq and writes to redis


### How to run

`docker-compose up -d`


#### Send a message to api service
`curl -XPOST localhost:8080/message -d '{"message": "msg1", "sender": "davud", "receiver": "gamelight"}'`
`curl -XPOST localhost:8080/message -d '{"message": "msg2", "sender": "davud", "receiver": "gamelight"}'`


#### Read the messages from Reporting service in chronological descending order 
`curl -XGET localhost:3000/message/list?sender=davud&receiver=gamelight`
or follow [localhost:3000/message/list?sender=davud&receiver=gamelight](localhost:3000/message/list?sender=davud&receiver=gamelight)

you'll get a response similar:
```json
[
    {
        "message": "msg2",
        "sender": "davud",
        "receiver": "gamelight",
        "sent_at": "2023-08-20T10:40:20.965199051Z"
    },
    {
        "message": "msg1",
        "sender": "davud",
        "receiver": "gamelight",
        "sent_at": "2023-08-20T10:40:18.564155651Z"
    }
]
```


### Some decisions:
- RabbitMQ
I couldn't connect to `rabbitmq:3.7-management` with `amqp://user:password@rabbitmq:5672`. As soon as i noticed it didn't work, I followed some articles and changed it to `rabbitmq-3` and removed the custom username and password to make it easy to connect with `guest:guest@...` and to not lose time debugging rabbitmq.


- Env variables:
all connection strings are defined in the code with plain text because of time constraint

- http data binding
I used domain models to bind the request data for both POST /message and GET /message/list routes to save time.

