package uds

//import "fmt"
import "testing"
import "time"
import "os"

var conn *QueueConnection
var client *Client

func TestMain(m *testing.M) {
	conn = NewQueueConnection("test", 4095)
	client = NewClient(conn, 0.2)
	os.Exit(m.Run())
}

func TestTransferData(t *testing.T) {

	go func() {
		r := conn.touser_frame()
		eq(t, r, []byte{0x36, 0x22, 0x12, 0x34, 0x56})
		conn.fromuserm.Lock()
		conn.fromuser.PushBack([]byte{0x76, 0x22, 0x89, 0xab, 0xcd, 0xef})
		conn.fromuserm.Unlock()
	}()

	response := client.Transfer_data(0x22, []byte{0x12, 0x34, 0x56})
	eq(t, response.Service_data["sequence_number_echo"], 0x22)
	eq(t, response.Service_data["parameter_records"], []byte{0x89, 0xab, 0xcd, 0xef})

	time.Sleep(20 * time.Millisecond)
}
