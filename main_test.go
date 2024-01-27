package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateBlock(t *testing.T) {
	prevBlock := &Block{
		Pos:       1,
		TimeStamp: "2022-01-01",
		Hash:      "prevHash",
	}

	data := BookCheckout{
		BookID:       "123",
		User:         "John Doe",
		CheckoutDate: "2022-01-02",
		IsGenesis:    false,
	}

	block := CreateBlock(prevBlock, data)

	// Add your assertions based on the expected behavior of CreateBlock
	if block.Pos != prevBlock.Pos+1 {
		t.Errorf("Expected Pos to be %d, got %d", prevBlock.Pos+1, block.Pos)
	}

}

func TestAddBlock(t *testing.T) {
	// Initialize a sample blockchain
	blockchain := NewBlockChain()

	// Create a sample book checkout data
	bookData := BookCheckout{
		BookID:       "456",
		User:         "Jane Doe",
		CheckoutDate: "2022-01-03",
		IsGenesis:    false,
	}

	// Add a block to the blockchain
	blockchain.AddBlock(bookData)

	// Add your assertions based on the expected behavior of AddBlock
	if len(blockchain.blocks) != 2 {
		t.Errorf("Expected blockchain length to be 2, got %d", len(blockchain.blocks))
	}

}

func TestHTTPHandlers(t *testing.T) {
	// Create a sample HTTP request for the newBook handler
	reqBody := `{"isbn": "123456789", "publish_date": "2022-01-01", "title": "Sample Book", "author": "John Doe"}`
	req, err := http.NewRequest("POST", "/new", strings.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	// Create a sample response recorder
	rr := httptest.NewRecorder()

	// Call the newBook handler
	newBook(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

}

func TestGenesisBlock(t *testing.T) {
	// Create a Genesis block
	genesisBlock := GenesisBlock()

	// Ensure that the Genesis block has the expected properties
	if genesisBlock.Pos != 0 {
		t.Errorf("Expected Genesis block Pos to be 1, got %d", genesisBlock.Pos)
	}

}

func TestHashGeneration(t *testing.T) {
	// Create a sample block
	block := &Block{
		Pos:       1,
		TimeStamp: "2022-01-01",
		PrevHash:  "prevHash",
		Data: BookCheckout{
			BookID:       "123",
			User:         "John Doe",
			CheckoutDate: "2022-01-02",
			IsGenesis:    false,
		},
	}

	// Generate the hash for the block
	block.generateHash()

	// Validate that the block's hash is not empty
	if block.Hash == "" {
		t.Error("Expected non-empty hash, got an empty string")
	}

}

func TestBlockchainStartsWithGenesisBlock(t *testing.T) {
	// Create a new blockchain
	blockchain := NewBlockChain()

	// Ensure that the blockchain starts with the Genesis block
	if len(blockchain.blocks) != 1 || blockchain.blocks[0].Pos != 0 {
		t.Error("Blockchain should start with the Genesis block")
	}

}
