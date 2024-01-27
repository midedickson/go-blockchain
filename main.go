package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/opensaucerer/barf"
)

// Block represenrs the individual block on the blockchain
// it holds the checkout data, position, timestamp, hash, and prevhash
type Block struct {
	Pos       int
	Data      BookCheckout
	TimeStamp string
	Hash      string
	PrevHash  string
}

func (b *Block) generateHash() {
	// method to generate block hash based on the checkout data sttacched to the block
	bytes, _ := json.Marshal(b.Data)
	data := string(rune(b.Pos)) + b.TimeStamp + string(bytes) + b.PrevHash

	hash := sha256.New()
	b.Hash = hex.EncodeToString(hash.Sum([]byte(data)))
}

func (b *Block) validateHash(hash string) bool {
	b.generateHash()
	return b.Hash == hash
}

// Book represents the item to be checked out
type Book struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	PublishDate string `json:"publish_date"`
	ISBN        string `json:"isbn"`
}

// BookCheckout represents the data of a Book to be checked out
type BookCheckout struct {
	BookID       string `json:"book_id"`
	User         string `json:"user"`
	CheckoutDate string `json:"checkout_date"`
	IsGenesis    bool   `json:"is_genesis"`
}

// BlockChain represents the core blockchain implementation
type BlockChain struct {
	blocks []*Block
}

func (bc *BlockChain) AddBlock(data BookCheckout) {
	// method to add a block to the array of block chains
	lastBlock := bc.blocks[len(bc.blocks)-1]

	block := CreateBlock(lastBlock, data)

	if validBlock(block, lastBlock) {
		bc.blocks = append(bc.blocks, block)
	}
}

var MyBlockChain *BlockChain

func CreateBlock(prevBlock *Block, data BookCheckout) *Block {
	// function to create a new block
	block := &Block{}

	if data.IsGenesis {
		block.Pos = 0
	} else {
		block.Pos = prevBlock.Pos + 1
	}

	block.TimeStamp = time.Now().String()
	block.Data = data
	block.PrevHash = prevBlock.Hash
	block.generateHash()

	return block
}

func validBlock(block, prevBlock *Block) bool {
	if prevBlock.Hash != block.PrevHash {
		return false
	}

	if !block.validateHash(block.Hash) {
		return false
	}

	if prevBlock.Pos+1 != block.Pos {
		return false
	}

	return true
}

func GenesisBlock() *Block {
	// creates the Genesis block in the block chain
	return CreateBlock(&Block{}, BookCheckout{IsGenesis: true})
}

func NewBlockChain() *BlockChain {
	return &BlockChain{[]*Block{GenesisBlock()}}
}

func newBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	err := barf.Request(r).Body().Format(&book)

	if err != nil {
		barf.Logger().Error(err.Error())
		barf.Response(w).Status(http.StatusInternalServerError).JSON(barf.Res{
			Status:  false,
			Data:    nil,
			Message: "Error creating book",
		})
	}

	h := md5.New()

	io.WriteString(h, book.ISBN+book.PublishDate)
	book.ID = hex.EncodeToString(h.Sum(nil))
	barf.Response(w).Status(http.StatusOK).JSON(barf.Res{
		Status:  true,
		Data:    book,
		Message: "New Book Created",
	})

}
func getBlockChain(w http.ResponseWriter, r *http.Request) {
	jbytes, err := json.MarshalIndent(MyBlockChain.blocks, "", " ")
	if err != nil {
		barf.Logger().Error(err.Error())
		barf.Response(w).Status(http.StatusInternalServerError).JSON(barf.Res{
			Status:  false,
			Data:    nil,
			Message: "Error getting blocks from chain",
		})
		return
	}

	barf.Response(w).Status(http.StatusOK).JSON(barf.Res{
		Status:  true,
		Data:    string(jbytes),
		Message: "Error getting blocks from chain",
	})

}
func writeBlock(w http.ResponseWriter, r *http.Request) {
	var bookCheckout BookCheckout
	err := barf.Request(r).Body().Format(&bookCheckout)

	if err != nil {
		barf.Logger().Error(err.Error())
		barf.Response(w).Status(http.StatusInternalServerError).JSON(barf.Res{
			Status:  false,
			Data:    nil,
			Message: "Error creating book checkout",
		})
	}

	MyBlockChain.AddBlock(bookCheckout)
	resp, err := json.MarshalIndent(bookCheckout, "", " ")
	if err != nil {
		barf.Logger().Error(err.Error())
		barf.Response(w).Status(http.StatusInternalServerError).JSON(barf.Res{
			Status:  false,
			Data:    nil,
			Message: "Error creating book checkout",
		})
	}

	barf.Response(w).Status(http.StatusOK).JSON(barf.Res{
		Status:  true,
		Data:    string(resp),
		Message: "New Block Created",
	})

}

func main() {

	MyBlockChain = NewBlockChain()

	barf.Get("/", getBlockChain)
	barf.Post("/", writeBlock)
	barf.Post("/new", newBook)
	go func() {
		for _, block := range MyBlockChain.blocks {
			fmt.Printf("Prev. Hash: %x\n", block.PrevHash)
			bytes, _ := json.MarshalIndent(block.Data, "", " ")
			fmt.Printf("Data: %v\n", string(bytes))
			fmt.Printf("Hash: %x\n", block.Hash)
			fmt.Println()
		}
	}()

	// start barf server
	if err := barf.Beck(); err != nil {
		barf.Logger().Error(err.Error())
	}

}
