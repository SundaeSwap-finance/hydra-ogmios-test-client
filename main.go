package main

import (
        "context"
        "encoding/json"
        "fmt"
        "os"
        "os/signal"
        
        "github.com/SundaeSwap-finance/ogmigo/v6"
)

type Transaction struct {
        CborHex string
        TxId    string
}

func parseTx(filepath string) (*Transaction, error) {
        data, err := os.ReadFile(filepath)
        if err != nil {
                return nil, fmt.Errorf("failed to read file: %w", err)
        }
        var txEnvelope struct {
                Transaction Transaction
        }
        err = json.Unmarshal(data, &txEnvelope)
        if err != nil {
                return nil, fmt.Errorf("failed to decode tx envelope: %w", err)
        }
        return &txEnvelope.Transaction, nil
}

func submitTx(endpoint string, filepath string) {
        tx, err := parseTx(filepath)
        ctx := context.Background()
        ogmigoClient := ogmigo.New(
                ogmigo.WithEndpoint(endpoint),
        )
        resp, err := ogmigoClient.SubmitTx(ctx, tx.CborHex)
        if err != nil {
                fmt.Printf("failed to submit tx: %w\n", err)
                return
        }
        fmt.Printf("%s\n", resp.ID)
}

func doChainSync(endpoint string) {
        ctx := context.Background()
        ogmigoClient := ogmigo.New(
                ogmigo.WithEndpoint(endpoint),
        )
        var callback ogmigo.ChainSyncFunc = func(ctx context.Context, data []byte) error {
                
                fmt.Printf("Received chainsync msg: %s\n", string(data))
                return nil
	}
	chainSync, err := ogmigoClient.ChainSync(ctx, callback,
		ogmigo.WithReconnect(true),
	)
	if err != nil {
		fmt.Println("%w", err)
	}
	
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Kill, os.Interrupt)

	select {
	case <-chainSync.Done():
		fmt.Println("chainsync done")
	case <-ctx.Done():
		fmt.Println("context done")
	case <-stop:
		fmt.Println("caught SIGINT")
	}
        errs := chainSync.Close()
        fmt.Printf("chainsync errs: %v\n", errs)
}

func main() {
        if len(os.Args) < 2 {
                fmt.Printf("no action supplied\n")
                return
        }
        if len(os.Args) < 3 {
                fmt.Printf("no endpoint supplied\n")
                return
        }
        args := os.Args[1:]
        action := args[0]
        endpoint := args[1]
        if action == "chainsync" {
                doChainSync(endpoint)
        } else if action == "submit-tx" {
                if len(args) < 3 {
                         fmt.Printf("no filepath supplied\n")
                         return
                }
                filepath := args[2]
                submitTx(endpoint, filepath) 
        } else {
                fmt.Printf("unexpected action: %w\n", action)
        }
}
