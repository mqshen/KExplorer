package explorer

import (
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/kversion"
	"kafkaexplorer/backend/types"
	"strings"
	"sync/atomic"
	"time"
)

type AdminClientService struct {
	adminClientCache map[string]*KafkaAdminClient
	clientIdSeq      atomic.Int32
}

func NewAdminClientService() (*AdminClientService, error) {
	return &AdminClientService{
		adminClientCache: make(map[string]*KafkaAdminClient),
	}, nil
}

func (a *AdminClientService) Get(c *types.Cluster) *KafkaAdminClient {
	adminClient, ok := a.adminClientCache[c.Name]
	if ok {
		return adminClient
	}
	adminClient = a.createAdminClient(c)
	a.adminClientCache[c.Name] = adminClient
	return adminClient
}

func (a *AdminClientService) createAdminClient(c *types.Cluster) *KafkaAdminClient {
	//clientId := a.generateClientId()
	//properties := &kafka.ConfigMap{
	//	"bootstrap.servers":  c.Bootstrap,
	//	"request.timeout.ms": c.ClientTimeout,
	//	"client.id":          clientId,
	//}
	bootstrapServers := strings.Split(c.Bootstrap, ",")
	adminClient, shutdownHook := newAdminClient(bootstrapServers)
	return &KafkaAdminClient{
		adminClient:  adminClient,
		shutdownHook: shutdownHook,
	}
}

func (a *AdminClientService) generateClientId() string {
	var seqNo int32
	a.clientIdSeq.Add(1)
	a.clientIdSeq.Store(seqNo)
	return "kafka-explorer-" + time.Nanosecond.String() + "-" + string(seqNo)
}

func newAdminClient(addrs []string) (*kadm.Client, func()) {

	client, err := kgo.NewClient(
		kgo.SeedBrokers(addrs...),

		// Do not try to send requests newer than 2.4.0 to avoid breaking changes in the request struct.
		// Sometimes there are breaking changes for newer versions where more properties are required to set.
		kgo.MaxVersions(kversion.V3_5_0()),
	)
	if err != nil {
		return nil, func() {}
	}

	adminClient := kadm.NewClient(client)

	return adminClient, func() { client.Close() }
}
