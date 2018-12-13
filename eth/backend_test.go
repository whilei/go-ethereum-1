package eth

import (
	"testing"

	"github.com/ethereum/go-ethereum/rpc"
)

// func TestSwagger2(t *testing.T) {
// }

func TestSwagger(t *testing.T) {
	var s = &Ethereum{}

	apis := []rpc.API{}
	// Append all the local APIs and return
	// NOTE: commented ones, in this test, are commented because they demand set up of additional
	// strucs and configs and things that I just don't have the patience to do right now
	apis = append(apis, []rpc.API{
		{
			Namespace: "eth",
			Version:   "1.0",
			Service:   NewPublicEthereumAPI(s),
			Public:    true,
		}, {
			Namespace: "eth",
			Version:   "1.0",
			Service:   NewPublicMinerAPI(s),
			Public:    true,
			// }, {
			// 	// 	Namespace: "eth",
			// 	// 	Version:   "1.0",
			// 	// 	Service:   downloader.NewPublicDownloaderAPI(s.protocolManager.downloader, s.eventMux),
			// 	// 	Public:    true,
		}, {
			Namespace: "miner",
			Version:   "1.0",
			Service:   NewPrivateMinerAPI(s),
			Public:    false,
			// }, {
			// 	Namespace: "eth",
			// 	Version:   "1.0",
			// 	Service:   filters.NewPublicFilterAPI(s.APIBackend, false),
			// 	Public:    true,
		}, {
			Namespace: "admin",
			Version:   "1.0",
			Service:   NewPrivateAdminAPI(s),
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPublicDebugAPI(s),
			Public:    true,
			// }, {
			// 	Namespace: "debug",
			// 	Version:   "1.0",
			// 	Service:   NewPrivateDebugAPI(s.chainConfig, s),
			// }, {
			// 	Namespace: "net",
			// 	Version:   "1.0",
			// 	Service:   s.netRPCService,
			// 	Public:    true,
		},
	}...)

	services := []interface{}{}
	for _, s := range apis {
		services = append(services, s)
	}

	metaAPI := rpc.API{
		Namespace: "introspect",
		Version:   "1.0",
		Public:    true,
	}
	introAPI := rpc.NewPublicIntrospectAPI(services)
	introAPI.API = metaAPI // if this even works... i'll gonna go drink a beer
	introAPI.Service = introAPI

	// apis = append(apis, metaAPI)

	// for _, api := range apis {
	b, err := introAPI.Swagger()
	if err != nil {
		t.Fatal(err)
	}
	// b, err := json.MarshalIndent(sw, "", "    ")
	// if err != nil {
	// t.Fatal(err)
	// }
	t.Log(string(b))
	// }
}
