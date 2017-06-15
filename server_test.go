package worq

import (
	"bufio"
	"encoding/json"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func runServer(runner func()) {
	s := NewServer()
	go s.Start()
	runner()
}

func TestServerStart(t *testing.T) {
	t.Parallel()
	runServer(func() {
		conn, err := net.DialTimeout("tcp", "localhost:7419", 1*time.Second)
		assert.NoError(t, err)
		buf := bufio.NewReader(conn)

		conn.Write([]byte("AHOY pwd:123456 other:thing\n"))
		result, err := buf.ReadString('\n')
		assert.NoError(t, err)
		assert.Equal(t, "OK\n", result)

		conn.Write([]byte("CMD foo\n"))
		result, err = buf.ReadString('\n')
		assert.NoError(t, err)
		assert.Equal(t, "ERR unknown command\n", result)

		conn.Write([]byte("PUSH {\"jid\":\"12345678901234567890abcd\",\"class\":\"Thing\",\"args\":[123],\"queue\":\"default\"}\n"))
		result, err = buf.ReadString('\n')
		assert.NoError(t, err)
		assert.Equal(t, "12345678901234567890abcd\n", result)

		conn.Write([]byte("POP default some other\n"))
		result, err = buf.ReadString('\n')
		assert.NoError(t, err)

		hash := make(map[string]interface{})
		err = json.Unmarshal([]byte(result), &hash)
		assert.NoError(t, err)
		//fmt.Println(hash)
		assert.Equal(t, "12345678901234567890abcd", hash["jid"])
		//assert.Equal(t, "{\"jid\":\"12345678901234567890abcd\",\"class\":\"Thing\",\"args\":[123],\"queue\":\"default\"}\n", result)

		conn.Write([]byte("END\n"))
		//result, err = buf.ReadString('\n')
		//assert.NoError(t, err)
		//assert.Equal(t, "OK\n", result)

		conn.Close()
	})

}
