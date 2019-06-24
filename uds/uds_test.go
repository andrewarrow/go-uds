package uds

//import "fmt"
import "testing"
import "time"
import "os"

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestTransferData(t *testing.T) {
	conn := NewQueueConnection("test", 4095)
	client := NewClient(conn, 0.2)

	go func() {
		conn.other()
		conn.fromuserm.Lock()
		conn.fromuser.PushBack([]byte{0x76, 0x22, 0x89, 0xab, 0xcd, 0xef})
		conn.fromuserm.Unlock()
	}()

	response := client.transfer_data(0x22, []byte{0x12, 0x34, 0x56})
	eq(t, response.service_data["sequence_number_echo"], 0x22, "TestTransferData")
	eq(t, response.service_data["parameter_records"], []byte{0x89, 0xab, 0xcd, 0xef}, "TestTransferData")

	time.Sleep(20 * time.Millisecond)
}
