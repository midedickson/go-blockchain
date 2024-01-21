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

type Block struct {
	Pos       int
	Data      BookCheckout
	TimeStamp string
	Hash      string
	PrevHash  string
}

func (b *Block) generateHash() {
	bytes, _ := json.Marshal(b.Data)
	data := string(rune(b.Pos)) + b.TimeStamp + string(bytes) + b.PrevHash

	hash := sha256.New()
	b.Hash = hex.EncodeToString(hash.Sum([]byte(data)))
}

type Book struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	PublishDate string `json:"publish_date"`
	ISBN        string `json:"isbn"`
}

type BookCheckout struct {
	BookID       string `json:"book_id"`
	User         string `json:"user"`
	CheckoutDate string `json:"checkout_date"`
	IsGenesis    bool   `json:"is_genesis"`
}

type BlockChain struct {
	blocks []*Block
}

func (bc *BlockChain) AddBlock(data BookCheckout) {
	lastBlock := bc.blocks[len(bc.blocks)-1]

	block := CreateBlock(lastBlock, data)

	if validBlock(block, lastBlock) {
		bc.blocks = append(bc.blocks, block)
	}
}

var MyBlockChain *BlockChain

func GenesisBlock() *Block {
	return CreateBlock(&Block{}, BookCheckout{IsGenesis: true})
}

func (b *Block) validateHash(hash string) bool {
	b.generateHash()
	return b.Hash != hash
}

func NewBlockChain() *BlockChain {
	return &BlockChain{[]*Block{GenesisBlock()}}
}

func CreateBlock(prevBlock *Block, data BookCheckout) *Block {
	block := &Block{}

	block.Pos = prevBlock.Pos + 1
	block.TimeStamp = time.Now().String()
	block.PrevHash = prevBlock.Hash
	block.generateHash()

	return block
}

func validBlock(block, prevBlock *Block) bool {
	if prevBlock.Hash != block.Hash {
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
		Data:    resp,
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
