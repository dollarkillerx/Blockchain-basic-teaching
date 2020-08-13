package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/dollarkillerx/erguotou"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

const difficulty = 10

type Block struct {
	Index      int    `json:"index"`
	Timestamp  string `json:"timestamp"`
	BPM        int    `json:"bpm"`
	Hash       string `json:"hash"`
	PrevHash   string `json:"prev_hash"`
	Difficulty int    `json:"difficulty"`
	Nonce      string `json:"nonce"`
}

var Blockchain []Block

func calculateHash(block Block) string {
	base := string(block.Index) + block.Timestamp + string(block.BPM) + block.Hash + block.PrevHash + block.Nonce
	h := sha256.New()
	h.Write([]byte(base))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

//func generateBlock(oldBlock Block, BPM int) (Block, error) {
//	var newBlock Block
//	t := int(time.Now().Unix())
//	newBlock.PrevHash = oldBlock.Hash
//	newBlock.BPM = BPM
//	newBlock.Index = oldBlock.Index + 1
//	newBlock.Timestamp = strconv.Itoa(t)
//	newBlock.Hash = calculateHash(oldBlock)
//
//	return newBlock, nil
//}

func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}
	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

type Message struct {
	BPM int `json:"BPM"`
}

var mutex = &sync.Mutex{}

// 创世纪
func era() {
	t := time.Now().Unix()
	genesisBlock := Block{Timestamp: strconv.Itoa(int(t)), Index: 0}
	genesisBlock.Hash = calculateHash(genesisBlock)
	Blockchain = append(Blockchain, genesisBlock)
}

// 商品的价值来自与稀缺
func main() {
	era()
	web()
}



func web() {
	app := erguotou.New()

	app.Get("/", blockchainGet)
	app.Post("/", blockchainWrite)

	if err := app.Run(erguotou.SetHost("0.0.0.0:8082")); err != nil {
		log.Fatalln(err)
	}
}


func blockchainGet(ctx *erguotou.Context) {
	ctx.Json(200, Blockchain)
}


func blockchainWrite(ctx *erguotou.Context) {
	var m Message
	if err := ctx.BindJson(&m); err != nil {
		ctx.Json(400, err)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()
	newBlock, err := generateBlock(Blockchain[len(Blockchain)-1], m.BPM)
	if err != nil {
		ctx.Json(500, err)
		return
	}

	if isBlockValid(newBlock, Blockchain[len(Blockchain) -1]) {
		Blockchain = append(Blockchain, newBlock)
		spew.Dump(Blockchain)
	}

	ctx.Json(200, newBlock)
}

// 工作量证明
func isHashValid(hash string, difficulty int) bool {
	repeat := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, repeat)
}

func generateBlock(oldBlock Block, BPM int) (Block, error) {
	var newBlock Block

	t := int(time.Now().Unix())

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = strconv.Itoa(t)
	newBlock.BPM = BPM
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Difficulty = difficulty

	for i:=0;;i++ {
		fmt.Println(i)
		hex := fmt.Sprintf("%x", i)
		newBlock.Nonce = hex
		if !isHashValid(calculateHash(newBlock), newBlock.Difficulty) {  // 证明工作量
			fmt.Println(calculateHash(newBlock), "do more work!")
			continue
		}else {
			fmt.Println(calculateHash(newBlock), "Work Done!")
			newBlock.Hash = calculateHash(newBlock)
			break
		}
	}

	return newBlock, nil
}