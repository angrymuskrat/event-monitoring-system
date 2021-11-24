# Microservice architecture

## Components

microservices: 
- data-storage - database access
- event-detection - event detection algorithm
- backend - communication with frontend
- coordinator - base pipeline
- insta-crawler - crawler for Instagram

data: 
- proto -GRPS data structures

### GRPC data types
GRPC data structures for communications between microservices.

### Data storage

Is the link between databases and services.

### Event detection 

The service is an implementation of the event detection algorithm.
<br>[Original algorithm](https://dl.acm.org/doi/10.1145/3282866.3282867).

### Coordinator

Automates the event search and monitoring pipelines, i.e. calls methods of different services 
according to a given scenario to collect data, process it, search for events and monitor it.

Due to the state of the crawler at the moment is only useful for restarting event search for 2018-2020

### Backend 

Is the link between the application core (modules, microservices, etc.) and the frontend.
<br>[Backend API](backend/README.md)

### Insta-crawler

Not working, dead component.