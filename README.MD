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
3. To start the docker container run: `docker run -it --rm --name ra-instance -p 3000:3000 -v /local/path/to/repo:/go/src/restaur-anteater -w /go/src/restaur-anteater ra-image`. Replace /local/path/to/repo with the actual repo path.
4. Once the docker container has started, run: `gin`. This will start the web app on port 3000 (go to localhost:3000 to access). Everytime you make changes to server.go it will automatically be recompiled in the docker container and the webapp will be restarted.

### Adding external go dependencies

1. Start the docker contain: `docker run -it --rm --name ra-instance -p 3000:3000 -v /local/path/to/repo:/go/src/restaur-anteater -w /go/src/restaur-anteater ra-image`.
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
