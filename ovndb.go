package libovndb

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/socketplane/libovsdb"
)

const (
	insert string = "insert"
	mutate string = "mutate"
	del    string = "delete"
	list   string = "select"
	update string = "update"
)

const (
	NBDB string = "OVN_Northbound"
)

const (
	LSWITCH     string = "Logical_Switch"
	LPORT       string = "Logical_Switch_Port"
	ACLS         string = "ACL"
	Address_Set string = "Address_Set"
)

const (
	UNIX string = "unix"
	TCP  string = "tcp"
)

const (
	//random seed.
	MAX_TRANSACTION = 1000
)

type ovnDBClient struct {
	socket   string
	server   string
	port     int
	protocol string
	dbclient *libovsdb.OvsdbClient
}

type ovnDBImp struct {
	client     *ovnDBClient
	cache      map[string]map[string]libovsdb.Row
	cachemutex sync.Mutex
	tranmutex  sync.Mutex
}

type OVNDB struct {
	imp *ovnDBImp
}

var once sync.Once
var ovnDBApi OVNDBApi

func GetInstance(socketfile string, protocol string, server string, port int) OVNDBApi {
	once.Do(func() {
		var dbapi *OVNDB
		var err error
		if protocol == UNIX {
			dbapi, err = newNBCtlBySocket(socketfile)
		} else if protocol == TCP {
			dbapi, err = newNBCtlByServer(server, port)
		} else {
			err = errors.New(fmt.Sprintf("The protocol [%s] is not supported", protocol))
		}

		if err != nil {
			panic(fmt.Sprint("Library libovndb initilizing failed", err))
			os.Exit(1)
		}
		ovnDBApi = dbapi
	})
	return ovnDBApi
}
