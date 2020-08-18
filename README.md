# SCOIR Technical Interview for Back-End Engineers
This repo contains an exercise intended for Back-End Engineers.

## Instructions
1. Fork this repo.
1. Using technology of your choice, complete [the assignment](./Assignment.md).
1. Update this README with
    * a `How-To` section containing any instructions needed to execute your program.
    * an `Assumptions` section containing documentation on any assumptions made while interpreting the requirements.
1. Before the deadline, submit a pull request with your solution.

## Expectations
1. Please take no more than 8 hours to work on this exercise. Complete as much as possible and then submit your solution.
1. This exercise is meant to showcase how you work. With consideration to the time limit, do your best to treat it like a production system.


## How-to
The most simple method of starting the app is using its `./bin/start-dev.sh` script. This will start a local db using docker-compose, build the apps binary, and run it.  
  
If you want to simply build and run the binary, that too should be fine, but without a localhost db or some database credentials set in the enviroment using the psql statandard envs the `-u` flag will not work. 

```sh
go build ./cmd/dirwatcher/
./dirwatcher <what ever flags you want>
```
Database envs
```sh
PGHOST
PGPORT
PGUSER
PGPASSWORD
PGDATABASE
```

### Customazation
The app requires 3 directories to work, these are created automatically for you if there is no configureation set, but you can also set the names of the dirs in two ways. 

```sh
INPUT_DIR
OUTPUT_DIR
ERROR_DIR
```