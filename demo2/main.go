package main

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
)

var Blockchain []Block

type Block struct {
	Index int `json:"index"`
	Timestamp string `json:"timestamp"`
	BPM int `json:"bpm"`
	Hash string `json:"hash"`
	PrevHash string `json:"prev_hash"`
}

func calculateHash(block Block) string {
	base := string(block.Index) + block.Timestamp + string(block.BPM) + block.Hash + block.PrevHash
	h := sha256.New()
	h.Write([]byte(base))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func generateBlock(oldBlock Block, BPM int) (Block, error) {
	var newBlock Block
	t := strconv.Itoa(int(time.Now().Unix()))
	newBlock.PrevHash = oldBlock.Hash
	newBlock.BPM = BPM
	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t
	newBlock.Hash = calculateHash(oldBlock)

	return newBlock, nil
}

func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index + 1 != newBlock.Index {
		return false
	}
	if oldBlock.Hash != newBlock.PrevHash {
		return  false
	}

	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

func replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}

func main() {

}
