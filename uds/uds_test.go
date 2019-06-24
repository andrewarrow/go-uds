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
	fmt.Println(client)
}
