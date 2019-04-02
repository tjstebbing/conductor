![Conductor Logo](/docs/logo.png)

Conductor does the least amount possible to solve the problem of 
starting up dependant services and shutting them down gracefully.

## Implementing a service:

A service implements the conductor.Service interface via a Run 
method which takes 3 channels:

* ***started chan bool*** - Send true when your service is running
* ***stopped chan bool*** - Send true when your service has stopped 
* ***stop chan context.Context*** - Select on this to be notified of shutdown, context with timeout.

```go

type MyService struct {}

func (t MyService) Run(started, stopped chan bool, stop chan context.Context) error {
    go func() {
      // Do some startup stuff then signal we're running
      started <- true
      select {
        case ctx := <- stop:
          // do some shutdown stuff then signal we're done
          stopped <- true
      }
    }()
    return nil
}
```


## Starting your services:

```go

c := conductor.New()
c.Service("Service 1", MyService{})
c.Service("Service 2", MyService{})
<- c.Start() // wait for graceful shutdown
```

## Stopping:

You can call `Conductor.Stop()` which will close the Shutdown channel returned by `Conductor.Start()`
once all services have closed or timeout is reached.

## Options!

Everyone loves options, these functions can all have their return values passed
to `Conductor.Run` to modify the behaviour of the Conductor.

* `HookSignals()` Will tell the conductor to listen for kill/stop signals and shutdown services.
* `StartupTimeout(time.Duration)` How long to wait for a service to start (Default 5s).
* `ShutdownTimeout(time.Duration)` How long to wait for a service to stop gracefully (Default 5s).
* `Noisy()` Log services starting and stopping via standard `log`. 

```go
c := conductor.New(
    conductor.HookSignals(),
    conductor.Noisy(),
    conductor.ShutdownTimeout(time.Duration(20 * time.Second)))
```
