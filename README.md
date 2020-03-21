# gateway-producer
My gateway-producer

this code uses dependencies of:
 * streadway/amqp
 * julienschmidt/httprouter
 
This project goes with another project called consumer (found in another repo),
Both uses as POC to rabbitmq producer consumer, and working with redis.

- The producer is a server that every route/ entrypoint send a message to some queue.
- The consumer is forever waiting for messages on specific queue, 
  and activating functions by info sent in the message.

To make this project work locally run 2 docker commands to start redis & rabbit before starting the producer or the consumer:
1) docker run -d -p 15672:15672 -p 5672:5672 -p 5671:5671 --hostname my-rabbitmq --name my-rabbitmq-container rabbitmq
2) docker run --name my-redis-container -p 6379:6379 -d redis

Enjoy and feel free to entr any new changs or comment about any wrong idea here :)
