# Message Relay

Implementation of the [transactional outbox pattern](https://microservices.io/patterns/data/transactional-outbox.html) in Go.

## Architecture

```mermaid
flowchart LR
  subgraph a [Application]
    direction TB
    s[Service]
    subgraph db [Database]
        direction RL
        oDB[(Outbox Table)]
        sDB[(Service Table)]
        oDB -. transaction .- sDB
    end
        s --> db
  end

  subgraph mr [Message Relay]
    direction LR
  end

  subgraph mbr [Message Broker]
    direction LR
  end

  subgraph c [Consumers]
  direction TB
  cX[Consumer X]
  cY[Consumer Y]
  cZ[Consumer Z]
  end
  a <-->|reads outbox table| mr
  mr --> mbr
  mbr --> cX
  mbr --> cY
  mbr --> cZ
```

## Todo List

- tests
- metrics
- health checks
- add a retry mechanism in case of a failure
- add a way to configure the number of retries
