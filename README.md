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
The most simple method of starting the app is using its `./bin/start-dev.sh` script. This will start a local db using docker-compose, build the apps binary, and run it. For  `./bin/start-dev.sh` to successfully run to will need 

```sh
docker-compose
flyway
```
  
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
env var keys for customizing the directories.
```sh
INPUT_DIR
OUTPUT_DIR
ERROR_DIR
```
And you can also use flags.  
```sh
Usage of ./dirwatcher:
  -errordir string
        the error dir, for csv files of errors caused the input file (default "./error-dir")
  -indir string
        the input dir or watch dir, for users to put their csv's in. (default "./input-dir")
  -outdir string
        the output dir, for for parsed user files (default "./out-dir")
  -u    pass this flag in to force parsing only unique file names
  ```
### -u flag
The `-u` flag is for forcing unique file names, so if a user adds a file into the input-dir, and then edits it and adds that file again, that file will not be processed again, and the file name should be updated for a reprocess. without this flag you can process the same file as many times as you'd like.

## Assumptions
I think 2 of the biggest assumtions I made are. 
1. this is a long running app, with no files in whatever dir the app is watching. 
"* once the application starts it watches _`input-directory`_ for any new files that need to be processed"  
2. to have the -u flag. The docs mention only processing "new" files that have not been marked as processed, but also mentiones in a file name colision to write over the last one. So because of that I made the unique file name optional. 

I also wasnt sure what kind of enviroment this should have lived in.. The psql db is probs overkill, but maybe not if this is living on like an ftp server where there can be tons and tons for files pouring in? 

Also wanted to consiter "outgrowning" a dir watcher, so created the parser as its own package to be able to use it on a server that is acceping files from users, or from events that pull the files down from some kind of document store etc.