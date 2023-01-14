# event driven architecture

EDA introduces new components, such as producer, consumer, event handler, command handler etc.

Some are categorised as primary ports, meaning they are responsible for messages coming in. 
- consumers
- event handler
- command handler

Secondary ports interacts with the message queue to publish messages.
- producers
- background workers

Background workers are there for pooling. 

There's concept of event/command processors too [^1].




[^1]: https://docs.axoniq.io/reference-guide/axon-framework/events/event-processors




