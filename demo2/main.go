package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

var Blockchain []Block

type Block struct {
	Index     int    `json:"index"`
	Timestamp int    `json:"timestamp"`
	BPM       int    `json:"bpm"`
	Hash      string `json:"hash"`
	PrevHash  string `json:"prev_hash"`
}

func calculateHash(block Block) string {
	base := string(block.Index) + strconv.Itoa(block.Timestamp) + string(block.BPM) + block.Hash + block.PrevHash
	h := sha256.New()
	h.Write([]byte(base))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func generateBlock(oldBlock Block, BPM int) (Block, error) {
	var newBlock Block
	t := int(time.Now().Unix())
	newBlock.PrevHash = oldBlock.Hash
	newBlock.BPM = BPM
	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t
	newBlock.Hash = calculateHash(oldBlock)

	return newBlock, nil
}

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

func replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}

var bcServer chan []Block

// 创世纪
func era() {
	t := time.Now().Unix()
	genesisBlock := Block{Timestamp: int(t), Index: 0}
	genesisBlock.Hash = calculateHash(genesisBlock)
	Blockchain = append(Blockchain, genesisBlock)
}

func main() {
	era()

	listen, err := net.Listen("tcp", "8085")
	if err != nil {
		log.Fatalln(err)
	}
	defer listen.Close()

	for {
		accept, err := listen.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		go handleConn(accept)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	io.WriteString(conn, "Enter a new BPM: ")
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		bpm, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Printf("%v not a number: %v \n", scanner.Text(), err)
			break
		}

		newBlock, err := generateBlock(Blockchain[len(Blockchain)-1], bpm)
		if err != nil {
			log.Println(err)
			continue
		}

		if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
			newBlockchain := append(Blockchain, newBlock)
			replaceChain(newBlockchain)
		}

		bcServer <- Blockchain
		io.WriteString(conn, "\n Enter a new BPM: ")
	}
}

func broadcast(conn net.Conn) {
loop:
	for {
		select {
		case block, exit := <- bcServer:
			if !exit {
				break loop
			}

			marshal, err := json.Marshal(block)
			if err != nil {
				log.Println(err)
				continue
			}
			conn.Write(marshal)
			spew.Dump(marshal)
		}
	}
}


