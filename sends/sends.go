/*
Benefits:
* pretty easy to implement
* fairly clear what's going on
* keeps a lid on memory usage, usage is fixed and calculable
Drawbacks:
  - fixed processor pool size may not meet your needs
  - you'll need to think about thread safety when modifying records
    in EventProcessor, or elsewhere.
  - retries not built in here.

Fixed processor size: this one's tricky. I think you could modify this to
support scaling up and down a worker pool with a little elbow grease.

For retries: in the past I've used PubSub which has a ack/nack mechanism
with a deadline, that allows you to extend the deadline for processing if
you need more time (I believe AWS has something similar for SQS). If you
succesfully process, you ack the message and it
drops off the queue, if you finish processing, but it's not successful, you
nack and let it be redelivered, and if you reach the deadline without doing
either, it's treated like a nack and the message is redelivered.
This style of retries complicates the above a little, because you need to
consider the message lifecycle in the processor. I think in your case
you only really need to consider ack/nack, but you won't need to worry about
deadline extension.

For ordering, it's tricky with messaging. Presumably your event processors have
access to the DB and can maybe check if related events must be processed first,
in which case, let the message drop back on the queue without processing.
This approach is fraught though, so maybe don't do that.

There are some libraries out there that might make some of the setup easier,
or give you other options, e.g. https://github.com/gocraft/work which bills itself
as Sidkiq for Go, or https://github.com/hibiken/asynq. I've used neither, and
can't speak to their quality, though I think Naimun is using asynq.

Personally, I'd start as small as possible and avoid the third party if possible
for a few reasons:
* understanding channels and goroutines will get you real far with Go,
  you should learn them, even if it slows down development time.
* these libraries have a lot of other dependencies
* It's easier to add these things later than it is to remove them.

That said, do you.

*/

package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Message struct {
	ID        int
	Timestamp time.Time
}

func ReceiveFromQueue(queue <-chan Message, eventIDs chan<- int) {
	// Assume that queue is provided by AWS libraries, but essentially
	// functions as a channel
	for msg := range queue {
		fmt.Printf("publishing event %d\n", msg.ID)
		eventIDs <- msg.ID
	}
}

func EventProcessor(eventIDs <-chan int) {
	for eventID := range eventIDs {
		fmt.Printf("processing event %d\n", eventID)
		// simulate work
		time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
	}
}

func main() {
	// this value is fixed at startup, one reason it might not be a good
	// choice is if you need to dynamically size up the number of workers
	// based on queue size, or some other factor.
	const maxWorkers = 3

	// This is a simulation for the external queue, probably SQS or something
	// similar. Avoid putting the message queue that kicks off processing in
	// memory, it's not durable enough to support that.
	queue := make(chan Message)
	// buffer this channel to prevent memory overflow
	// and apply backpressure on publisher, which can probably
	// outpace your processing speed (I assume)
	eventIDs := make(chan int, maxWorkers)

	for i := 0; i < 50; i++ {
		queue <- Message{i, time.Now()}
	}

	go ReceiveFromQueue(queue, eventIDs)

	for i := 0; i <= maxWorkers; i++ {
		go EventProcessor(eventIDs)
	}

	// simulate a server started on main thread
	time.Sleep(3 * time.Second)
}
