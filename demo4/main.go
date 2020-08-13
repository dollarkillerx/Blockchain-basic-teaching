package main

//我们将有一个基于GO的TCP服务器，其他节点（验证器）可以连接它。
//最新的区块链状态将周期性地广播到每个节点。
//每个节点将提出新的块。
//基于每个节点所持有的令牌数量，将随机选择一个节点（根据所持有的令牌数量加权）作为获胜者，并将其块添加到块链。

type Block struct {
	Index int `json:"index"`
	Timestamp string `json:"timestamp"`
	DPM int `json:"dpm"`
	Hash string `json:"hash"`
	PreHash string `json:"pre_hash"`
	Validator string `json:"validator"`  // 验证器
}



func main() {
	
}
