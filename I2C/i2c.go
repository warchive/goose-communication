
/**i2c package to allow users to read from 
and write to slave I2C devices(arduinos)**/
package i2c

import (
		"golang.org/x/exp/io/i2c/driver"
)

const tenBitMask = 1 << 12

/** Device represents an I2C device(arduino for us).
Devices must be closed one they are no longer in use **/
type Device struct {
	conn driver.Conn
}

/** TenbBit marks an I2C address as a 10-bit address**/
func TenBit(addr int) int {
	return addr | tenBitMask
}

// Read reads len(buf) bytes from the device
func (d *Device) Read(buf []byte) error {
	return d.conn.Tx(nil,buf)
}

// ReadReg is Read from a register
func (d *Device) ReadReg(reg byte, buf []byte) error {
	return d.conn.Tx([]byte{reg}, buf)
}

/** Write writes the buffer to the device. If it is required
to write to a specific register, the register should be passed
as the first byte in the given buffer **/
func (d *Device) Write(buf []byte) (err error) {
	return d.conn.Tx(buf, nil)
}

// WriteReg is similar to Write but writes to a register.
func (d *Device) WriteReg(reg byte, buf []byte) (err error) {
	return d.conn.Tx(append( []byte{reg},buf...), nil)
}

// Close closes the devices and releases underlying sources.
func (d *Device) Close() error {
	return d.conn.Close()
}

/**Open opens a connection to an I2C device.
// All devices must be closed once they are no longer in use.
// For devices that use 10-bit I2C addresses, addr can be marked
as a 10-bit address with TenBit. **/
func Open(o driver.Opener, addr int) (*Device, error) {
	unmasked, tenbit := resolveAddr(addr)
	conn, err := o.Open(unmasked, tenbit)
	if err != nil {
		return nil, err
	}
	return &Device{conn: conn}, nil
}

// resolveAddr returns whether the addr is 10-bit masked or not.
// It also returns the unmasked address.
func resolveAddr(addr int) (unmasked int, tenbit bool) {
	return addr & (tenbitMask - 1), addr&tenbitMask == tenbitMask
}