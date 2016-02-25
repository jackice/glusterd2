// Package context is the runtime context of GlusterD
//
// Any package that needs access to the GlusterD runtime context just needs to
// import this package.
package context

import (
	"os"
	"sync"

	"github.com/gluster/glusterd2/config"
	"github.com/gluster/glusterd2/rest"
	"github.com/gluster/glusterd2/transaction"
	"github.com/gluster/glusterd2/utils"

	log "github.com/Sirupsen/logrus"
	etcdclient "github.com/coreos/etcd/client"
	"github.com/pborman/uuid"
)

// Various version constants that will be used by GD2
const (
	MaxOpVersion    = 40000
	APIVersion      = 1
	GlusterdVersion = "4.0-dev"
)

// Any object that is a part of the GlusterD context and needs to be available
// to other packages should be declared here as exported global variables
var (
	MyUUID         uuid.UUID
	Rest           *rest.GDRest
	TxnFw          *transaction.GDTxnFw
	OpVersion      int
	EtcdProcessCtx *os.Process
	EtcdClient     etcdclient.Client
	HostIP         string
)

var (
	initOnce sync.Once
)

func initOpVersion() {
	//TODO : Need cluster awareness and then decide the op-version
	OpVersion = MaxOpVersion
}

//initETCDClient will initialize etcd client that will be use during member add/remove in the cluster
func initETCDClient() error {
	c, err := etcdclient.New(etcdclient.Config{Endpoints: []string{"http://" + HostIP + ":2379"}})
	if err != nil {
		log.WithField("err", err).Error("Failed to create etcd client")
		return err
	}
	EtcdClient = c

	return nil
}

func doInit() {
	log.Debug("Initializing GlusterD context")

	utils.InitDir(config.LocalStateDir)

	initMyUUID()
	initOpVersion()

	Rest = rest.New()

	initStore()

	// Initializing etcd client
	err := initETCDClient()
	if err != nil {
		log.WithField("err", err).Error("Failed to initialize etcd client")
		return
	}

	log.Debug("Initialized GlusterD context")
}

// Init initializes the GlusterD context. This should be called once before doing anything else.
func Init() {
	initOnce.Do(doInit)
}

// GetEtcdMemberAPI() will return etcd MemberAPI interface
func GetEtcdMemberAPI() etcdclient.MembersAPI {
	var c etcdclient.Client
	return etcdclient.NewMembersAPI(c)
}

// AssignEtcdProcessCtx () is to assign the etcd ctx in context.EtcdCtx
func AssignEtcdProcessCtx(ctx *os.Process) {
	EtcdProcessCtx = ctx
}

// SetLocalHostIP() function will set LocalIP address
func SetLocalHostIP() {
	hostIP, err := utils.GetLocalIP()
	if err != nil {
		log.Fatal("Could not able to get IP address")
	}

	HostIP = hostIP
}
