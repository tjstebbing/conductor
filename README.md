![Conductor Logo](/docs/logo.png)

Conductor does the least amount possible to solve the problem of 
starting up dependant services and shutting them down gracefully.

## Implementing a service:

A service implements the conductor.Service interface via a Run 
method which takes 3 channels:

***started chan bool*** - Send true when your service is running
***stopped chan bool*** - Send true when your service has stopped 
***stop chan context.Context*** - Select on this to be notified of shutdown, context with timeout.

```go

type MyService struct {}

func (t MyService) Run(started, stopped chan bool, stop chan context.Context) {
  go func() {
    // Do some startup stuff then signal we're running
    started <- true
    select {
      case ctx := <- stop:
        // do some shutdown stuff then signal we're done
        stoopped <- true
    }
  }()

}```


## Starting your services:

```go

c := conductor.New()
c.Service(MyService{})
c.Service(MyService{})
c.Start()
```
