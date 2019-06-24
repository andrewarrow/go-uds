package isotp

import "fmt"
import "testing"
import "os"

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestClient(t *testing.T) {
	conn := NewQueueConnection("test", 4095)
	client := NewClient(conn, 0.2)
	response := client.transfer_data(0x22, []byte{0x12, 0x34, 0x56})
	fmt.Println(response)
}
