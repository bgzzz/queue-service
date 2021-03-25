# Simple queue implementation 

Simple application containing two elements queue tool
and reader-writer tool.

## Queue tool

Is a tool that implements queuing functionality on the basis of http server 

## Reader-writer, reader mode 

Reader is a tool that reads the file line by line and pushes the content to 
a predefined queue.

## Reader-writer, writer  mode 

Writer is a tool that reads the queue writes  the content to 
a predefined file. Writing terminates when there is a end of file signal received.
Write also terminates if queue is empty 

More about parameters [here](./reader-writer/main.go#L50)

## How to build 

```
make build
```

# TODO
- unit tests
- docker wrapping 
- docker-compose 