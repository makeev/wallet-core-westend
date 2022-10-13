package main

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"tw/core"
	"tw/protos/polkadot"
)

func main() {
	fmt.Println("==> calling wallet core from go")

	// можем создать кошелек из mnemonic или seed, но только по схеме ed25519(polkadot по умолчанию sr25519)
	mn := "excess shift gym east salmon ghost enroll winner document library pen minimum"

	fmt.Println("==> mnemonic is valid: ", core.IsMnemonicValid(mn))

	coin := core.CoinTypeWestend

	// coin details
	fmt.Println(coin.GetName())
	fmt.Println(coin.Decimals())

	w, _ := core.CreateWalletWithMnemonic(mn, coin)
	printWallet(w)

	fromAddress := w.Address
	fmt.Println("From: ", fromAddress)

	// будет разный, для разных сетей, можно получить через rpc call
	genesisHash, _ := hex.DecodeString("e143f23803ac50e8f6f8e62695d1ce9e4e1d68aa36c1cd2cfd15340213f3423e")
	fmt.Println("Genesis hash:", genesisHash)

	// ключ кошелька, которым будем подписывать
	priKeyByte, _ := hex.DecodeString(w.PriKey)
	fmt.Println("Private key:", priKeyByte)

	// последний блок тоже достается через rpc call
	blockHash, _ := hex.DecodeString("90d81701b3c3d38eb9aa54e2f6fa167e2e7fc21e154cf91fabdcca7bf58164e2")
	toAddress := "5FPJXVP6C8h53epLDm74PxCNQTHSiCXsgqFEbt5Z1eFikX9j"

	input := polkadot.SigningInput{
		BlockHash:          blockHash,
		GenesisHash:        genesisHash,
		Nonce:              3,
		SpecVersion:        9300,
		TransactionVersion: 13,
		Tip:                big.NewInt(100).Bytes(),
		Era: &polkadot.Era{
			BlockNumber: 12861317,
			Period:      64,
		},
		PrivateKey: priKeyByte,
		Network:    polkadot.Network_WESTEND,
		MessageOneof: &polkadot.SigningInput_BalanceCall{
			BalanceCall: &polkadot.Balance{
				MessageOneof: &polkadot.Balance_Transfer_{
					Transfer: &polkadot.Balance_Transfer{
						ToAddress: toAddress,
						Value:     big.NewInt(10000).Bytes(),
					},
				},
			},
		},
	}

	fmt.Println(input.String())

	var output polkadot.SigningOutput
	output_err := core.CreateSignedTx(&input, coin, &output)
	if output_err != nil {
		panic(output_err)
	}

	fmt.Println("")
	// тут бует tx-hash, используя c++ и Extrinsic.h можно отдельно генерить
	// payload и подпись, либо написать функцию чтобы бробросить через cgo
	fmt.Println(hex.EncodeToString(output.GetEncoded()))
}

func printWallet(w *core.Wallet) {
	fmt.Printf("%s wallet: \n", w.CoinType.GetName())
	fmt.Printf("\t address: %s \n", w.Address)
	fmt.Printf("\t pri key: %s \n", w.PriKey)
	fmt.Printf("\t pub key: %s \n", w.PubKey)
	fmt.Println("")
}
