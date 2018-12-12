package eth

import (
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/rpc"
)

func TestSwagger(t *testing.T) {
	var s = &Ethereum{}

	apis := []rpc.API{}
	// Append all the local APIs and return
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

	for _, api := range apis {
		sw, err := api.Swagger()
		if err != nil {
			t.Fatal(err)
		}
		b, err := json.MarshalIndent(sw, "", "    ")
		if err != nil {
			t.Fatal(err)
		}
		t.Log(string(b))
	}
}
