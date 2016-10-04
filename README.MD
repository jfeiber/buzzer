# Hack-A-Thing 10/4
Joon Cho: I built a comment box component in ReactJS using this tutorial: [ReactJs Tutorial](https://facebook.github.io/react/docs/tutorial.html).
I learned that ReactJS is extremely flexible, is probably the best library for creating front end objects, and should be the standard. It was hard to createClass something beautiful in CSS because I can't just modify the HTML - I have to work within ReactJS guidelines.


# Restaurant App

We are allowing restaurants to manage their clients better with a web app that displays analytics registered by proprietary restaurant pagers.


## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisities

1. [Docker](https://docs.docker.com/docker-for-mac/)
2.  You can install go locally if you really want but it's pretty annoying to get dependencies working properly.

### Installing

1. Clone this repo
2. cd into the repo folder and run: `docker build -t ra-image .`
3. Start a Postgres container: `docker run --name ra-postgres -e POSTGRES_PASSWORD=password -e POSTGRES_USER=ra -d postgres`
4. To start the docker container run: `docker run -it --rm --name ra-instance --link ra-postgres:postgres -p 3000:3000 -v $(pwd):/go/src/restaur-anteater -w /go/src/restaur-anteater ra-image`
5. Once the docker container has started, run `goose -env dev up` to setup the database tables.
6. Then run: `gin`. This will start the web app on port 3000 (go to localhost:3000 to access). Everytime you make changes to server.go it will automatically be recompiled in the docker container and the webapp will be restarted.

### Using psql to interact with the DB

Assumes you setup the Postgres container using the above instructions.

1. Run: `docker run -it --rm --link ra-postgres:postgres postgres psql -h postgres -U ra`

### Using migrations

1. Start the docker container: `docker run -it --rm --name ra-instance --link ra-postgres:postgres -p 3000:3000 -v $(pwd):/go/src/restaur-anteater -w /go/src/restaur-anteater ra-image`.
2. Run `goose -env dev up`. This will migrate the dev database (the Postgres db running in the Postgres docker container) to the latest migration. Migrations are in `db/migrations`. To migrate
the production database, run `goose -env production up`.

### Adding external go dependencies

1. Start the docker container: `docker run -it --rm --name ra-instance -p 3000:3000 -v $(pwd):/go/src/restaur-anteater -w /go/src/restaur-anteater ra-image`.
2. Run `godep restore`
3. Run `go get dependency` where `dependency` is something you want to install (like `github.com/tools/godep`)
4. Run `godep save`
5. Exit the docker container and check (with `git diff` or something) that the changes to Godeps/Godeps.json didn't do anything crazy.


## Authors

* **Kevin Neilson**
* **Michael Machlin**
* **Joon Cho**
* **Mickey Alperin**
* **Josh Feiber**
