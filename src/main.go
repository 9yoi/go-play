package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	fmt.Println("Hello World")

	// The method signature
	methodSig := "increment(uint256)"

	// Hash the method signature and take first 4 bytes as method id
	methodID := crypto.Keccak256([]byte(methodSig))[:4]

	fmt.Printf("Method ID: %x\n", methodID)

	// The value we want to pass to the function
	arg := big.NewInt(10) // Using big.Int to represent uint256

	// Pack the function arguments
	// Manually pack the argument to 32 bytes
	packedArg := make([]byte, 32)
	copy(packedArg[32-len(arg.Bytes()):], arg.Bytes())

	// The payload for the transaction is the function selector followed by the packed argument
	payload := append(methodID, packedArg...)

	fmt.Printf("Data: %x\n", payload)

	client, err := ethclient.Dial("wss://sepolia.infura.io/ws/v3/{key}")
	if err != nil {
		log.Fatal(err)
	}

	// The contract address
	contractAddress := common.HexToAddress("0x97e06bfb0b2EB74750f33d2af48aeB7E0656fa0E")

	// Sender's address
	fromAddress := common.HexToAddress("{address}")

	// Fetch the nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalf("Failed to fetch nonce: %v", err)
	}

	// Define gas price and gas limit
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Failed to suggest gas price: %v", err)
	}
	gasLimit := uint64(4700000) // Adjust as necessary

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &contractAddress,
		Value:    big.NewInt(0),
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     payload,
	})

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA("{address}")
	if err != nil {
		log.Fatalf("Failed to decode private key: %v", err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(signedTx)

	fmt.Printf("tx sent: %s", signedTx.Hash().Hex())

}

// func subscribeToLogs(client *ethclient.Client, query ethereum.FilterQuery, payload bytes) (ethereum.Subscription, error) {
// 	logs := make(chan types.Log)
// 	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
// 	if err != nil {
// 		return nil, err
// 	}

// 	go func() {
// 		for {
// 			select {
// 			case log := <-logs:
// 				fmt.Printf("Log received: %v\n", log)
// 				fmt.Println("Address: ", log.Address)
// 				fmt.Println("Topics:", log.Topics)
// 				fmt.Println("Data: ", log.Data)

// 			// payload

// 			case err := <-sub.Err():
// 				log.Printf("Error: %v\n", err)
// 				return
// 			}
// 		}
// 	}()

// 	return sub, nil
// }
