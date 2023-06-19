package pocket

import (
	"context"
	"testing"
	"time"

	"github.com/dan13ram/wpokt-backend/app"
	"github.com/dan13ram/wpokt-backend/models"
)

var mockResultTx ResultTx = ResultTx{
	Height: 1,
	Index:  1,
	Hash:   "test-hash",
	StdTx: StdTxParams{
		Memo: `{"address": "recipient-address", "chainId": "test-chain"}`,
		Msg: TxMsg{
			Type: "test-type",
			Value: TxMsgValue{
				FromAddress: "sender-address",
				ToAddress:   "recipient-address",
				Amount:      "1000000",
			},
		},
	},
}

// Mock implementation of the Client interface
type mockClient struct{}

func (c *mockClient) GetHeight() (*HeightResponse, error) {
	return &HeightResponse{Height: 10}, nil
}

func (c *mockClient) GetAccountTxsByHeight(height int64) ([]*ResultTx, error) {
	return []*ResultTx{&mockResultTx}, nil
}

func (c *mockClient) ValidateNetwork() {}

func (c *mockClient) GetBlock() (*BlockResponse, error) {
	return &BlockResponse{}, nil
}

// Helper function to set up a test instance of WPOKTMintMonitor
func setupMintMonitor() *WPOKTMintMonitor {
	app.Config.Pocket.MonitorIntervalSecs = 1
	app.Config.Pocket.ChainId = "test-chain"
	app.Config.Ethereum.ChainId = 0
	app.DB = &mockDB{}
	Client = &mockClient{}

	monitor := NewMintMonitor().(*WPOKTMintMonitor)
	monitor.stop = make(chan bool)

	return monitor
}

// Mock implementation of the database
type mockDB struct{}

// insertOne
func (db *mockDB) InsertOne(collectionName string, document interface{}) error {
	return nil
}

func (db *mockDB) Connect(ctx context.Context) error {
	return nil
}

func (db *mockDB) SetupIndexes() error {
	return nil
}

func (db *mockDB) Disconnect() error {
	return nil
}

func TestWPOKTMintMonitor_Stop(t *testing.T) {
	monitor := setupMintMonitor()

	// Call the Stop() method
	go func() {
		time.Sleep(2 * time.Second) // Wait for a few seconds before stopping
		monitor.Stop()
	}()

	// Start the monitor and wait for it to stop
	go monitor.Start()

	// If the monitor doesn't stop within a reasonable time, it means Stop() didn't work
	select {
	case <-monitor.stop:
		// Monitor stopped as expected
	case <-time.After(5 * time.Second):
		t.Errorf("Stop() did not stop the monitor")
	}
}

func TestWPOKTMintMonitor_UpdateCurrentHeight(t *testing.T) {
	monitor := setupMintMonitor()

	// Call the UpdateCurrentHeight() method
	monitor.UpdateCurrentHeight()

	// Check if the currentHeight field is updated correctly
	if monitor.currentHeight != 10 {
		t.Errorf("UpdateCurrentHeight() failed. Expected currentHeight: 10, got: %d", monitor.currentHeight)
	}
}

func TestWPOKTMintMonitor_HandleInvalidMint(t *testing.T) {
	monitor := setupMintMonitor()

	tx := &mockResultTx

	// Call the HandleInvalidMint() method
	result := monitor.HandleInvalidMint(tx)

	// Check if the method returns the expected result
	if result != true {
		t.Errorf("HandleInvalidMint() failed. Expected: true, got: %t", result)
	}
}

func TestWPOKTMintMonitor_HandleValidMint(t *testing.T) {
	monitor := setupMintMonitor()

	tx := &mockResultTx

	// Call the HandleValidMint() method
	result := monitor.HandleValidMint(tx, models.MintMemo{})

	// Check if the method returns the expected result
	if result != true {
		t.Errorf("HandleValidMint() failed. Expected: true, got: %t", result)
	}
}

func TestWPOKTMintMonitor_SyncTxs(t *testing.T) {
	monitor := setupMintMonitor()

	// Call the SyncTxs() method
	result := monitor.SyncTxs()

	// Check if the method returns the expected result
	if result != true {
		t.Errorf("SyncTxs() failed. Expected: true, got: %t", result)
	}
}

func TestWPOKTMintMonitor_Start(t *testing.T) {
	monitor := setupMintMonitor()

	// Call the Start() method
	go func() {
		time.Sleep(2 * time.Second) // Wait for a few seconds before stopping
		monitor.Stop()
	}()

	// Start the monitor and wait for it to stop
	go monitor.Start()

	// If the monitor doesn't stop within a reasonable time, it means Stop() didn't work
	select {
	case <-monitor.stop:
		// Monitor stopped as expected
	case <-time.After(5 * time.Second):
		t.Errorf("Start() did not stop the monitor")
	}
}

func TestNewMintMonitor(t *testing.T) {
	app.Config.Pocket.MonitorIntervalSecs = 1
	app.Config.Pocket.StartHeight = -1

	monitor := NewMintMonitor()

	// Check if the monitor instance is created correctly
	switch m := monitor.(type) {
	case *WPOKTMintMonitor:
		if m.monitorInterval != 1*time.Second || m.startHeight != 10 || m.currentHeight != 10 {
			t.Errorf("NewMintMonitor() failed. Unexpected monitor values")
		}
	default:
		t.Errorf("NewMintMonitor() failed. Unexpected monitor type")
	}
}
