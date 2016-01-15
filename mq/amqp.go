package mq

import (
    "fmt"
    "log"
    "strings"
    "encoding/json"

    "github.com/streadway/amqp"

    "github.com/ugarcia/go_test_common/models"
    "github.com/ugarcia/go_test_common/util"
)

// Main structure
type AMQP struct {
    Connection *amqp.Connection
    Channel *amqp.Channel
    Queues map[string]*amqp.Queue
}

/**
 * Initializes the connection, channel and queues for the struct
 */
func (mq *AMQP) Init(address string) (*AMQP) {

    // Holder var for errors
    var err error

    // Create connection
    mq.Connection, err = amqp.Dial(address)
    util.FailOnError(err, "Failed to connect to RabbitMQ address: " + address)

    // Create Channel
    mq.Channel, err = mq.Connection.Channel()
    util.FailOnError(err, "Failed to open RabbitMQ channel")

    // Set quality into channel (not sure if it even works ...)
    err = mq.Channel.Qos(
        1,     // prefetch count
        0,     // prefetch size
        false, // global
    )
    util.FailOnError(err, "Failed to set QoS on channel")

    // Init queues map
    mq.Queues = make(map[string]*amqp.Queue)

    return mq
}

/**
 * Closes the connection and the channel
 */
func (mq *AMQP) Close() (*AMQP) {

    // Just close channel and connection
    mq.Channel.Close()
    mq.Connection.Close()

    return mq
}

/**
 * Register exchange dispatchers inside a channel
 */
func (mq *AMQP) RegisterExchanges(exchanges []string) (*AMQP) {

    // Loop over exchanges
    for _, exchange := range exchanges {

        // Declare exchange dispatcher
        err := mq.Channel.ExchangeDeclare(
            exchange,   // name
            "topic", // type
            true,     // durable
            false,    // auto-deleted
            false,    // internal
            false,    // no-wait
            nil,      // arguments
        )
        util.FailOnError(err, "Failed to declare into channel the exchange " + exchange)
    }

    return mq
}

/**
 * Register Queues
 */
func (mq *AMQP) RegisterQueues(queues []string) (*AMQP) {

    // Loop over exchanges
    for _, queue := range queues {

        // Create MCP Inbound Queue
        q, err := mq.Channel.QueueDeclare(
            queue, // name
            true, // durable
            false, // delete when unused
            false, // exclusive
            false, // no-wait
            nil, // arguments
        )
        util.FailOnError(err, "Failed to declare queue " + queue)

        // Assign to struct map (hope it does not complain if it exists ...)
        mq.Queues[queue] = &q
    }

    return mq
}

/**
 * Binds Queues to Exchanges under a routing key
 */
func (mq *AMQP) BindQueuesToExchange(queues []string, exchange string, routing string) (*AMQP) {

    // Loop over exchanges
    for _, queue := range queues {

        // Bind MCP Queue to MCP exchange thru channel
        err := mq.Channel.QueueBind(
            queue, // queue name
            routing,     // routing key
            exchange, // exchange
            false,
            nil,
        )
        util.FailOnError(err, "Failed to bind queue " + queue + " to exchange " + exchange + " with routing " + routing)
    }

    return mq
}

/**
 * Starts consuming messages into a queue
 */
func (mq *AMQP) Consume(queue string, fn func(models.QueueMessage, amqp.Delivery)) {

    // Consume queue messages
    deliveries, err := mq.Channel.Consume(
        queue, // queue
        "", // consumer
        false, // auto-ack
        false, // exclusive
        false, // no-local
        false, // no-wait
        nil, // args
    )
    util.FailOnError(err, "Failed to register consumer for queue " + queue)

    // Loop forever for messages, and call handler for each
    forever := make(chan bool)
    go func() {

        // Loop over deliveries
        for d := range deliveries {

            // Log message
            log.Printf("Received a message: %s", d.Body)

            // Parse received message
            msg := models.QueueMessage{}
            if err := json.Unmarshal([]byte(d.Body), &msg); err != nil {
                fmt.Println(err.Error())
                continue
            }

            // Call handler function
            go fn(msg, d)
        }
    }()
    log.Printf(" [*] Waiting for messages on queue %s. To exit press CTRL+C\n", queue)
    <-forever
}

/**
 * Sends a message to a exchange channel with route
 */
func (mq *AMQP) SendMessage(inMsg models.QueueMessage) (*AMQP) {

    // Encode message
    msg, err := json.Marshal(inMsg)
    if err != nil {
        fmt.Println(err.Error())
        return mq
    }

    // Parse exchange from original sender
    exchange := strings.Split(inMsg.Receiver, ".")[0]

    // Try to publish message to relevant exchange channel with proper routing
    err = mq.Channel.Publish(
        exchange,     // exchange
        inMsg.Receiver, // routing key
        false,  // mandatory
        false,  // immediate
        amqp.Publishing{
            DeliveryMode: amqp.Persistent,
            ContentType: "application/json",
            Body:        []byte(msg),
        })
    util.FailOnError(err, "Failed to publish a message")
    log.Printf(" [x] Sent %s", string(msg))

    return mq
}
