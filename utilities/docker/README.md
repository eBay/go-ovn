- Goal: Add ability to run the go-ovn tests inside of a docker container which in turn would make it easy to run on 
  non-linux systems that have docker installed. The docker image installs `go` binaries and sets the `GOPATH` environment
  variable.

- To build the docker image execute the below on a system where docker is installed

  ```docker build -t <image_name>:<tag_name> .```

- To execute the tests, 
  a)first, mount the path to the go-ovn repo when running the container as below:
  
  ```
  docker run -itd -v <PATH_TO_GO_OVN_REPO>:/go/src/github.com/eBay/go-ovn <image_name>:<tag_name>
  ```
  
  b)next, exec into shell of the running container ( `docker exec -it <running_container_name> bash`) and execute the
  below command to run the tests
  
  ```
  OVS_RUNDIR=/var/run/ovn OVN_NB_DB=unix:ovnnb_db.sock OVN_SB_DB=unix:ovnsb_db.sock go test -v ./...
  ```