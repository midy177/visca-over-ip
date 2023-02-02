package visca

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type categoryCode byte

const (
	_categoryCodeInterface categoryCode = 0x00
	_categoryCodeCamera1   categoryCode = 0x04
	_categoryCodePanTilter categoryCode = 0x06
)

type command byte

const (
	_commandPanTiltDrive command = 0x01
	_commandMemory       command = 0x3f
	_commandZoom         command = 0x07
)

type payloadType [2]byte

var (
	_payloadTypeCommand = payloadType{0x01, 0x00}
	_payloadTypeInquiry = payloadType{0x01, 0x10}
	_payloadTypeReply   = payloadType{0x01, 0x11}
)

const _terminator byte = 0xff
const _headerLength = 8
const _minPayloadLength = _headerLength + 1 // +1 for the terminator

var (
	ErrArgsTooLong       = errors.New("payload arguments length must not exceed 0xffff")
	ErrMissingTerminator = errors.New("payload missing terminator")
)

var (
	ErrSyntaxError       = errors.New("command format valid or illegal command parameters")
	ErrCommandBufferFull = errors.New("unable to accept a command because 2 commands are currently being executed")
)

type payload struct {
	Type payloadType

	IsInquiry    bool
	CategoryCode categoryCode
	Command      command
	Args         []byte
}

func (p *payload) MarshalBinary() ([]byte, error) {
	length := uint16(len(p.Args) + 5)
	if length > 0xffff {
		return nil, ErrArgsTooLong
	}

	var buf []byte

	// bytes 0 and 1 are for the type
	// this function only works for typeCommand rn
	buf = append(buf, p.Type[0])
	buf = append(buf, p.Type[1])

	// bytes 2 and 3 are for the length
	lengthBuf := make([]byte, 2)
	binary.BigEndian.PutUint16(lengthBuf, length)
	buf = append(buf, lengthBuf...)

	// bytes 4-7 are for who knows what
	buf = append(buf, make([]byte, 4)...)

	// byte 8 is for the address (always 0x81 for IP)
	buf = append(buf, 0x81)

	// byte 9 is either command (0x01) or inquiry (0x09)
	if p.IsInquiry {
		buf = append(buf, 0x09)
	} else {
		buf = append(buf, 0x01)
	}

	// byte 10 is the category code
	buf = append(buf, byte(p.CategoryCode))

	// byte 11 is the command
	buf = append(buf, byte(p.Command))

	// byte 12+ is the command arguments
	buf = append(buf, p.Args...)

	// last byte is the terminator
	buf = append(buf, _terminator)

	return buf, nil
}

func (p *payload) UnmarshalBinary(data []byte) error {
	// Commenting out these lines due to having a similar health check 11 lines down
	//if len(data) < _minPayloadLength {
	//	return fmt.Errorf("data in bad format: %# x", data)
	//}

	// bytes 0 and 1 are the type
	p.Type[0] = data[0]
	p.Type[1] = data[1]

	// bytes 2 and 3 are the length
	length := binary.BigEndian.Uint16([]byte{data[2], data[3]})

	// bytes 4-7 are still a mystery...

	// Logic that checks that the data received the full visca package
	// This keeps fractured visca packets from getting processed
	// We found that AVer cameras do no send back the proper length bit but
	// other manufacturers do like BirdDog.  (AVer sends back a 1 for lenth)
	// This fix will still allow AVer to function while also allowing BirdDog to work
	if len(data) < _headerLength+int(length) {
		return fmt.Errorf("data in bad format: %# x", data)
	}

	switch p.Type {
	case _payloadTypeCommand:
		// byte 8 is just the address, which we don't use rn
		// byte 9 is either command or inquiry
		p.IsInquiry = data[9] == 0x09

		// byte 10 is the category code
		p.CategoryCode = categoryCode(data[10])

		// byte 11 is the command
		p.Command = command(data[11])

		// byte 12+ is the command args
		for i := 12; i < _headerLength+int(length); i++ {
			p.Args = append(p.Args, data[i])
		}
	case _payloadTypeReply:
		for i := _headerLength; i < _headerLength+int(length); i++ {
			p.Args = append(p.Args, data[i])
		}
	default:
		return fmt.Errorf("don't know how to parse payloadType: %# x", p.Type)
	}

	// the last byte is a terminator
	if data[len(data)-1] != _terminator {
		return ErrMissingTerminator
	}

	return nil
}

func (p *payload) IsAck() bool {
	return p.Type == _payloadTypeReply && len(p.Args) == 1 && p.Args[0] == 0x90
}

func (p *payload) Error() error {
	switch {
	case p.Type != _payloadTypeReply:
	case len(p.Args) != 3:
	case p.Args[0] != 0x90:
	case p.Args[1] == 0x60 && p.Args[2] == 0x02:
		return ErrSyntaxError
	case p.Args[1] == 0x60 && p.Args[2] == 0x03:
		return ErrCommandBufferFull
	default:
		return fmt.Errorf("unknown error: %# x", p.Args)
	}

	return nil

	// TODO need to check these errors, see if they ever show up
	//case p.Args[2] == 0x04:
	//case p.Args[2] == 0x05:
	//case p.Args[2] == 0x41:
}
