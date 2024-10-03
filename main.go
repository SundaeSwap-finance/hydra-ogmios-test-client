package main

import (
        "fmt"
        
        "github.com/SundaeSwap-finance/ogmigo/v6"
)

func main() {
        ogmigoClient := ogmigo.New(
                ogmigo.WithEndpoint("0.0.0.0:9001"),
        )
        var callback ogmigo.ChainSyncFunc = func(ctx context.Context, data []byte) error {
                fmt.Println("Received chainsync msg: %s", string(data))
	}
	chainSync, err := ogmigoClient.ChainSync(ctx, callback,
		ogmigo.WithReconnect(true),
	)
	if err != nil {
		return err
	}
	defer chainSync.Close()
}
