package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/davecgh/go-spew/spew"
	"github.com/dollarkillerx/erguotou"
	"log"
	"net/http"
	"strconv"
	"time"
)

/**
1. 创建自己的区块链
2. 理解HASH函数如何保存区块链的完整性
3. 如何创建并添加新的块
4. 通过浏览器来查看整个链
5. 所有其他关于区块链的基础知识
这里用到的WEB Framework `go get github.com/dollarkillerx/erguotou`
*/

type Block struct {
	Index     int    // 区块位置
	Timestamp string // 生成时间
	BPM       int    // 心率
	Hash      string
	PrevHash  string // 前一块的hash
}

// 获取当前区块的hash
func calculateHash(block Block) string {
	recode := string(block.Index) + block.Timestamp + string(block.BPM) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(recode))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// 更具前一个块生成一个新的块
func generateBlock(oldBlock Block, BPM int) (Block, error) {
	var newBlock Block

	now := time.Now()
	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = now.String()
	newBlock.BPM = BPM
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)
	return newBlock, nil
}

// 验证
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

// 如果当前链是旧的旧更新当前链
func replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}

var Blockchain []Block // 高效BUG 内存吃紧  推荐链表实现

func handleGetBlockchain(ctx *erguotou.Context) {
	bytes, err := json.MarshalIndent(Blockchain, "", " ")
	if err != nil {
		ctx.Json(http.StatusInternalServerError, err)
		return
	}
	ctx.Write(http.StatusOK, bytes)
}

type Message struct {
	BPM int
}

func handleWriteBlock(ctx *erguotou.Context) {
	dpmString, bool := ctx.PathValueString("dpm")
	if bool == false {
		ctx.Json(400, errors.New("DPM IS NULL"))
		return
	}
	atoi, err := strconv.Atoi(dpmString)
	if err != nil {
		ctx.Json(400, err)
		return
	}
	block, err := generateBlock(Blockchain[len(Blockchain)-1], atoi)
	if err != nil {
		ctx.Json(500, err)
		return
	}

	if isBlockValid(block, Blockchain[len(Blockchain)-1]) {
		newBlockchain := append(Blockchain, block)
		replaceChain(newBlockchain)
		spew.Dump(Blockchain)
	}

	ctx.Json(200, block)
}

// 上帝说要有光  创世块
func haveLight() {
	t := time.Now().String()
	genesisBlock := Block{0, t, 0, "", ""}
	genesisBlock.Hash = calculateHash(genesisBlock)
	spew.Dump(genesisBlock)
	Blockchain = append(Blockchain, genesisBlock)
}

func main() {
	haveLight()
	app := erguotou.New()

	app.Get("/", handleGetBlockchain)
	app.Get("/add/:dpm", handleWriteBlock)

	if err := app.Run(erguotou.SetHost(":8081"), erguotou.SetDebug(true)); err != nil {
		log.Fatalln(err)
	}
}
