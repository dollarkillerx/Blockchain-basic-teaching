package main

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

/**
	1. 创建自己的区块链
	2. 理解HASH函数如何保存区块链的完整性
	3. 如何创建并添加新的块
	4. 通过浏览器来查看整个链
	5. 所有其他关于区块链的基础知识
 */

type Block struct {
	Index int          // 区块位置
	Timestamp string   // 生成时间
	BPM int            // 心率
	Hash string
	PrevHash string    // 前一块的hash
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
	return newBlock,nil
}

// 验证
func isBlockValid(newBlock,oldBlock Block) bool {
	if oldBlock.Index + 1 != newBlock.Index {
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

var Blockchain []Block   // 高效BUG 内存吃紧  推荐链表实现

func main() {

}
