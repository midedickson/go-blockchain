package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/opensaucerer/barf"
)

type Block struct {
	Pos       int
	Data      BookCheckout
	TimeStamp string
	Hash      string
	PrevHash  string
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

func newBook(w http.ResponseWriter, r *http.Request) {
	var book *Book
	err := barf.Request(r).Body().Format(&book)

	if err != nil {
		barf.Logger().Error(err.Error())
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
func getBlockChain(w http.ResponseWriter, r *http.Request) {}
func writeBlock(w http.ResponseWriter, r *http.Request)    {}

func main() {

	barf.Get("/", getBlockChain)
	barf.Post("/", writeBlock)
	barf.Post("/new", newBook)

	// start barf server
	if err := barf.Beck(); err != nil {
		barf.Logger().Error(err.Error())
	}
}
