# Harvard-arts-reverse-proxy
A reverse proxy API to consume the Harvard Arts Museum API.

# Table of Contents
* [Features](https://github.com/Nedson202/Harvard-arts-reverse-proxy#features)
* [Technologies](https://github.com/Nedson202/Harvard-arts-reverse-proxy#technologies)
* [Installation Setup](https://github.com/Nedson202/Harvard-arts-reverse-proxy#installation-setup)
* [Testing](https://github.com/Nedson202/Harvard-arts-reverse-proxy#testing)
* [Usage](https://github.com/Nedson202/Harvard-arts-reverse-proxy#usage)
* [Language](https://github.com/Nedson202/Harvard-arts-reverse-proxy#language)
* [Dependencies](https://github.com/Nedson202/Harvard-arts-reverse-proxy#dependencies)
* [License](https://github.com/Nedson202/Harvard-arts-reverse-proxy#license)

## Features
* GET places information
* GET collection objects
* GET collection object
* GET publication object
* Full-text search support for fields like title, century, accessionmethod, period, technique, classification, department, culture, medium, verificationleveldescription

## Technologies
* Golang
* Redis
* Elasticsearch
* Docker
* Harvard Art Museum API

## Installation Setup

* **Clone repo:**

  Open **CMD(command prompt)** for windows users, or any other terminal you use.

  ```
    git clone https://github.com/Nedson202/Harvard-arts-reverse-proxy.git
  ```

* **Start app:**

  This codebase uses the compile-daemon module to watch for changes and trigger a restart.

  * Create a .env file in the root directory of the codebase
  * Copy the content of the .env.sample file and add their corresponding values appropriately.

  ```
    Change directory to cloned repo (Harvard-arts-reverse-proxy)

    $ cd Harvard-arts-reverse-proxy

    Run development server

    $ make start-dev
  ```

  If you want to use a dockerized version instead use the following:

  ```
    $ make start

    View docker logs

    $ make logs
  ```

  You can access the API via http://localhost:4000/api/v1/

## Testing
Golangs testing library *testing.T :) is in play here.

```
  $ go test --cover
```

## Project structure
```
├── Dockerfile
├── LICENSE
├── Procfile
├── README.md
├── bin
│   └── harvard-arts-reverse-proxy
├── docker-compose.yml
├── go.mod
├── go.sum
├── main.go
├── makefile
├── places.json
├── reverse_proxy
│   ├── app.go
│   ├── client.go
│   ├── collections.go
│   ├── config.go
│   ├── elasticsearch.go
│   ├── handler.go
│   ├── healthcheck.go
│   ├── healthcheck_test.go
│   ├── places.go
│   ├── places_test.go
│   ├── publications.go
│   ├── redis.go
│   ├── routes.go
│   └── types.go
└── vendor
    ├── github.com
    └── modules.txt
```

## Usage

| HTTP VERB | Description | Endpoints |
| --- | --- | --- |
| `GET` | Health check | /api/v1/health |
| `GET` | Retrieves collection objects | /api/v1/objects?size={size}&page={page} |
| `GET` | Retrieves a collection object | /api/v1/object/{objectId} |
| `GET` | Search for collection objects. See search fields below | /api/v1/search?query={searchText}&size={size}&from={from} |
| `GET` | Retrieves publications by year | /api/v1/publications?size={size}&page={page}&year={year} |
| `GET` | Retrieves place data including geolocation | /api/v1/places?placeId={placeId} |
| `GET` | Retrieves place IDs and parent location | /api/v1/places/id?size={size}&from={from} |

#### Search fields
* title, century, accessionmethod, period, technique, classification, department, culture, medium, verificationleveldescription

## Language
* Golang

## Dependencies
> Click [here](https://github.com/Nedson202/Harvard-arts-reverse-proxy/blob/develop/go.mod) to view all dependencies.

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please endeavour to update tests as appropriate.

## License

> You can check out the full license [here](https://github.com/Nedson202/Harvard-arts-reverse-proxy/blob/develop/LICENSE)

This project is licensed under the terms of the MIT license.

