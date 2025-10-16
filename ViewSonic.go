package viewsonic

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/jpillora/backoff"
)

// ErrFunctionDisabled is returned when the projector indicates that a function is disabled (greyed out).
// This typically occurs when there are no source inputs to the projector, making certain functions
// like "Aspect Ratio" unavailable via OSD menu or remote control.
var ErrFunctionDisabled = fmt.Errorf("function is disabled (greyed out) on the projector")

// ViewSonic is a connection to the projector. It is designed to be thread-safe.
// It automatically handles reconnects in the background.
type ViewSonic struct {
	conn             net.Conn
	mutex            sync.Mutex
	cancelContext    context.CancelFunc
	triggerReconnect chan struct{}
}

// New generates a Connection and starts a background goroutine to maintain the connection.
func New(ip string) *ViewSonic {
	// Create custom dialer in order to set TCP KeepAlive
	var DefaultDialer = &net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 10 * time.Second,
	}

	// Create Context to cancel reconnect loop
	ctx, cancel := context.WithCancel(context.Background())

	c := &ViewSonic{
		cancelContext:    cancel,
		triggerReconnect: make(chan struct{}, 1),
	}

	// Initial connection attempt
	tmpConn, err := DefaultDialer.Dial("tcp", ip)
	if err != nil {
		log.Println(err)
		// Trigger immediate reconnect attempt in the background
		select {
		case c.triggerReconnect <- struct{}{}:
		default:
		}
	} else {
		c.conn = tmpConn
	}

	b := &backoff.Backoff{
		Min:    2 * time.Second,
		Max:    30 * time.Second,
		Factor: 1.19,
		Jitter: false,
	}

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("Stopping Reconnect loop")
				c.mutex.Lock()
				if c.conn != nil {
					c.conn.Close()
				}
				c.mutex.Unlock()
				return
			case <-c.triggerReconnect:
				c.mutex.Lock()
				if c.conn != nil {
					c.conn.Close()
					c.conn = nil
				}

				conn, err := DefaultDialer.Dial("tcp", ip)
				if err != nil {
					log.Println("reconnect error: ", err)
					c.mutex.Unlock()
					// Try again with backoff
					time.Sleep(b.Duration())
					select {
					case c.triggerReconnect <- struct{}{}:
					default:
					}
				} else {
					log.Println("reconnect success:", ip)
					b.Reset()
					c.conn = conn
					c.mutex.Unlock()
				}
			case <-ticker.C:
				// Health check: read power status
				c.mutex.Lock()
				if c.conn != nil {
					_, _, err := tx(c.conn, c.triggerReconnect, cmdRead, []byte{0x34, 0x00, 0x00, 0x11, 0x00})
					if err != nil {
						log.Printf("Health check failed for %v: %v", ip, err)
					}
				} else {
					c.mutex.Unlock()
					// If there's no connection, trigger a reconnect.
					select {
					case c.triggerReconnect <- struct{}{}:
					default:
					}
				}
				c.mutex.Unlock()
			}
		}
	}()

	return c
}

// Close closes the connection to the projector and stops the reconnect loop.
func (conn *ViewSonic) Close() {
	conn.cancelContext()
}

const (
	cmdError         = 0x00
	cmdWriteKey      = 0x02
	cmdWriteResponse = 0x03
	cmdReadResponse  = 0x05
	cmdWrite         = 0x06
	cmdRead          = 0x07
)

// tx handles the low-level transmission of a command and reception of a response.
// It is responsible for locking the connection, sending the packet, and parsing the response.
// If any network error occurs, it triggers a reconnect and returns the error.
// The caller is responsible for retrying the command if necessary.
func (conn *ViewSonic) tx(cmd1 uint8, data []byte) (uint8, []byte, error) {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()

	if conn.conn == nil {
		// If connection is not available, trigger a reconnect and return an error immediately.
		select {
		case conn.triggerReconnect <- struct{}{}:
		default:
		}
		return 0, nil, fmt.Errorf("connection not established")
	}

	return tx(conn.conn, conn.triggerReconnect, cmd1, data)
}

func tx(conn net.Conn, triggerReconnect chan struct{}, cmd1 uint8, data []byte) (uint8, []byte, error) {
	// Clear Buffer by setting a short deadline and reading whatever is there.
	conn.SetReadDeadline(time.Now().Add(1 * time.Millisecond))
	_, _ = io.Copy(io.Discard, conn)

	// Build Packet
	packet := make([]byte, 0, 4+len(data)+1)
	packet = append(packet, cmd1, 0x14, 0x00)
	packet = binary.LittleEndian.AppendUint16(packet, uint16(len(data)))
	packet = append(packet, data...)
	packet = append(packet, checkSum(packet[1:]))

	conn.SetWriteDeadline(time.Now().Add(2 * time.Second))
	_, err := conn.Write(packet)
	if err != nil {
		select {
		case triggerReconnect <- struct{}{}:
		default:
		}
		return 0, nil, fmt.Errorf("write error: %w", err)
	}

	// Read Response
	head := make([]byte, 5)

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, err := io.ReadFull(conn, head)
	if err != nil {
		log.Printf("Error reading Command Head: %x", head[:n])
		select {
		case triggerReconnect <- struct{}{}:
		default:
		}
		return 0, nil, err
	}

	dataLen := binary.LittleEndian.Uint16(head[3:5]) // Command Length

	rxData := make([]byte, dataLen+1) // +1 for Checksum
	conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	n, err = io.ReadFull(conn, rxData)
	if err != nil {
		log.Printf("Error reading Command Data: %x", rxData[:n])
		select {
		case triggerReconnect <- struct{}{}:
		default:
		}
		return 0, nil, err
	}

	// Verify Checksum
	if rxData[len(rxData)-1] != checkSum(head[1:], rxData[:len(rxData)-1]) {
		log.Printf("Error: invalid checksum: %x", rxData)
		select {
		case triggerReconnect <- struct{}{}:
		default:
		}
		return 0, nil, fmt.Errorf("invalid checksum")
	}

	return head[0], rxData[:len(rxData)-1], nil // return cmdType and data without checksum
}

// checkSum adds alle values toggetter and returns the sum
func checkSum(data ...[]byte) byte {
	cs := byte(0)
	for _, s := range data {
		for _, b := range s {
			cs += b
		}
	}
	return cs
}

// Write sends a write command to the projector.
// If the connection is down, it will return an error. The background process is responsible for reconnecting.
// The caller may choose to retry the command after a short delay.
func (conn *ViewSonic) Write(command uint16, value uint8) error {
	cmd1, data, err := conn.tx(cmdWrite, []byte{0x34, byte(command >> 8), byte(command), value})
	if err != nil {
		return err
	}

	if cmd1 == cmdError {
		return ErrFunctionDisabled
	}

	if cmd1 != cmdWriteResponse || len(data) != 0 {
		return fmt.Errorf("unexpected response command: 0x%02X, %x", cmd1, data)
	}

	return nil
}

// WriteKey is a specialized version of Write for Remote Key commands such as Menu, Enter, etc.
func (conn *ViewSonic) WriteKey(command uint16, value uint8) error {
	cmd1, data, err := conn.tx(cmdWriteKey, []byte{0x34, byte(command >> 8), byte(command), value})
	if err != nil {
		return err
	}

	if cmd1 == cmdError {
		return ErrFunctionDisabled
	}

	if cmd1 != cmdWriteResponse || len(data) != 0 {
		return fmt.Errorf("unexpected response command: 0x%02X, %x", cmd1, data)
	}

	return nil
}

// Read sends a read command to the projector and returns a single byte.
// If the connection is down, it will return an error. The background process is responsible for reconnecting.
// The caller may choose to retry the command after a short delay.
func (conn *ViewSonic) Read(command uint16) (uint8, error) {
	cmd1, data, err := conn.tx(cmdRead, []byte{0x34, 0x00, 0x00, byte(command >> 8), byte(command)})
	if err != nil {
		return 0, err
	}

	if cmd1 == cmdError {
		return 0, ErrFunctionDisabled
	}

	if cmd1 != cmdReadResponse || len(data) != 3 {
		return 0, fmt.Errorf("unexpected response command: 0x%02X, %x", cmd1, data)
	}

	return data[2], nil
}

// Read2Bytes sends a read command and returns two bytes as an int16.
// If the connection is down, it will return an error. The background process is responsible for reconnecting.
// The caller may choose to retry the command after a short delay.
func (conn *ViewSonic) Read2Bytes(command uint16) (int16, error) {
	cmd1, data, err := conn.tx(cmdRead, []byte{0x34, 0x00, 0x00, byte(command >> 8), byte(command)})
	if err != nil {
		return 0, err
	}

	if cmd1 == cmdError {
		return 0, ErrFunctionDisabled
	}

	if cmd1 != cmdReadResponse || len(data) != 4 {
		return 0, fmt.Errorf("unexpected response command: 0x%02X, %x", cmd1, data)
	}

	return int16(binary.LittleEndian.Uint16(data[2:4])), nil
}

// ReadNBytes sends a read command and returns N bytes of data.
// If the connection is down, it will return an error. The background process is responsible for reconnecting.
// The caller may choose to retry the command after a short delay.
func (conn *ViewSonic) ReadNBytes(command uint16) ([]byte, error) {
	cmd1, data, err := conn.tx(cmdRead, []byte{0x34, 0x00, 0x00, byte(command >> 8), byte(command)})
	if err != nil {
		return nil, err
	}

	if cmd1 == cmdError {
		return nil, ErrFunctionDisabled
	}

	if cmd1 != cmdReadResponse {
		return nil, fmt.Errorf("unexpected response command: 0x%02X, %x", cmd1, data)
	}

	return data, nil
}
