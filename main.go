package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/coinbase/rosetta-sdk-go/asserter"
	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/marpme/digibyte-rosetta-node/client"
	"github.com/marpme/digibyte-rosetta-node/configuration"
	"github.com/marpme/digibyte-rosetta-node/services"
)

const (
	ConfigPath = "ROSETTA_DGB_CONFIG_PATH"
)

func NewBlockchainRouter(client client.DigibyteClient) http.Handler {
	assert, err := asserter.NewServer([]*types.NetworkIdentifier{
		{
			Blockchain:           client.GetConfig().NetworkIdentifier.Blockchain,
			Network:              client.GetConfig().NetworkIdentifier.Network,
			SubNetworkIdentifier: nil,
		},
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to create asserter: %v\n", err)
		os.Exit(1)
	}

	networkAPIController := server.NewNetworkAPIController(services.NewNetworkAPIService(client), assert)
	return server.NewRouter(networkAPIController)
}

func main() {
	configPath := os.Getenv(ConfigPath)
	if configPath == "" {
		configPath = "config.yaml"
	}
	cfg, err := configuration.New(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to parse config: %v\n", err)
		os.Exit(1)
	}

	client, err := client.NewDigibyteClient(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to prepare IoTex gRPC client: %v\n", err)
		os.Exit(1)
	}

	router := NewBlockchainRouter(client)
	fmt.Println("Listening on ", "0.0.0.0:"+cfg.Server.Port)
	err = http.ListenAndServe("0.0.0.0:"+cfg.Server.Port, router)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Digibyte Rosetta Gateway server exited suddenly: %v\n", err)
		os.Exit(1)
	}
}
