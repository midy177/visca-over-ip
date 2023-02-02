package visca

import (
	"testing"

	"github.com/matryer/is"
)

type payloadMarshalTest struct {
	name        string
	payload     *payload
	expected    []byte
	expectedErr error
}

var _payloadMarshalTests = []payloadMarshalTest{
	{
		name: "TiltUp",
		payload: &payload{
			Type:         _payloadTypeCommand,
			IsInquiry:    false,
			CategoryCode: _categoryCodePanTilter,
			Command:      _commandPanTiltDrive,
			Args: []byte{
				PanTiltSpeedMin,
				0x0e,
				PanDirectionStop,
				TiltDirectionUp,
			},
		},
		expected: []byte{0x01, 0x00, 0x00, 0x09, 0x00, 0x00, 0x00, 0x00, 0x81, 0x01, 0x06, 0x01, 0x00, 0x0e, 0x03, 0x01, 0xff},
	},
	{
		name: "TiltDown",
		payload: &payload{
			Type:         _payloadTypeCommand,
			IsInquiry:    false,
			CategoryCode: _categoryCodePanTilter,
			Command:      _commandPanTiltDrive,
			Args: []byte{
				PanTiltSpeedMin,
				0x0e,
				PanDirectionStop,
				TiltDirectionDown,
			},
		},
		expected: []byte{0x01, 0x00, 0x00, 0x09, 0x00, 0x00, 0x00, 0x00, 0x81, 0x01, 0x06, 0x01, 0x00, 0x0e, 0x03, 0x02, 0xff},
	},
	{
		name: "PanLeft",
		payload: &payload{
			Type:         _payloadTypeCommand,
			IsInquiry:    false,
			CategoryCode: _categoryCodePanTilter,
			Command:      _commandPanTiltDrive,
			Args: []byte{
				0x0b,
				PanTiltSpeedMin,
				PanDirectionLeft,
				TiltDirectionStop,
			},
		},
		expected: []byte{0x01, 0x00, 0x00, 0x09, 0x00, 0x00, 0x00, 0x00, 0x81, 0x01, 0x06, 0x01, 0x0b, 0x00, 0x01, 0x03, 0xff},
	},
	{
		name: "PanRight",
		payload: &payload{
			Type:         _payloadTypeCommand,
			IsInquiry:    false,
			CategoryCode: _categoryCodePanTilter,
			Command:      _commandPanTiltDrive,
			Args: []byte{
				0x0b,
				PanTiltSpeedMin,
				PanDirectionRight,
				TiltDirectionStop,
			},
		},
		expected: []byte{0x01, 0x00, 0x00, 0x09, 0x00, 0x00, 0x00, 0x00, 0x81, 0x01, 0x06, 0x01, 0x0b, 0x00, 0x02, 0x03, 0xff},
	},
	{
		name: "PanTiltStop",
		payload: &payload{
			Type:         _payloadTypeCommand,
			IsInquiry:    false,
			CategoryCode: _categoryCodePanTilter,
			Command:      _commandPanTiltDrive,
			Args: []byte{
				PanTiltSpeedMin,
				PanTiltSpeedMin,
				PanDirectionStop,
				TiltDirectionStop,
			},
		},
		expected: []byte{0x01, 0x00, 0x00, 0x09, 0x00, 0x00, 0x00, 0x00, 0x81, 0x01, 0x06, 0x01, 0x00, 0x00, 0x03, 0x03, 0xff},
	},
	{
		name: "ZoomTele",
		payload: &payload{
			Type:         _payloadTypeCommand,
			IsInquiry:    false,
			CategoryCode: _categoryCodeCamera1,
			Command:      _commandZoom,
			Args: []byte{
				_funcZoomTele,
			},
		},
		expected: []byte{0x01, 0x00, 0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x81, 0x01, 0x04, 0x07, 0x02, 0xff},
	},
	{
		name: "ZoomWide",
		payload: &payload{
			Type:         _payloadTypeCommand,
			IsInquiry:    false,
			CategoryCode: _categoryCodeCamera1,
			Command:      _commandZoom,
			Args: []byte{
				_funcZoomWide,
			},
		},
		expected: []byte{0x01, 0x00, 0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x81, 0x01, 0x04, 0x07, 0x03, 0xff},
	},
	{
		name: "ZoomStop",
		payload: &payload{
			Type:         _payloadTypeCommand,
			IsInquiry:    false,
			CategoryCode: _categoryCodeCamera1,
			Command:      _commandZoom,
			Args: []byte{
				_funcZoomStop,
			},
		},
		expected: []byte{0x01, 0x00, 0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x81, 0x01, 0x04, 0x07, 0x00, 0xff},
	},
	{
		name: "MemorySet0",
		payload: &payload{
			Type:         _payloadTypeCommand,
			IsInquiry:    false,
			CategoryCode: _categoryCodeCamera1,
			Command:      _commandMemory,
			Args: []byte{
				_funcMemorySet,
				0x00,
			},
		},
		expected: []byte{0x01, 0x00, 0x00, 0x07, 0x00, 0x00, 0x00, 0x00, 0x81, 0x01, 0x04, 0x03f, 0x01, 0x00, 0xff},
	},
	{
		name: "MemorySet2",
		payload: &payload{
			Type:         _payloadTypeCommand,
			IsInquiry:    false,
			CategoryCode: _categoryCodeCamera1,
			Command:      _commandMemory,
			Args: []byte{
				_funcMemorySet,
				0x02,
			},
		},
		expected: []byte{0x01, 0x00, 0x00, 0x07, 0x00, 0x00, 0x00, 0x00, 0x81, 0x01, 0x04, 0x03f, 0x01, 0x02, 0xff},
	},
	{
		name: "MemoryRecall1",
		payload: &payload{
			Type:         _payloadTypeCommand,
			IsInquiry:    false,
			CategoryCode: _categoryCodeCamera1,
			Command:      _commandMemory,
			Args: []byte{
				_funcMemoryRecall,
				0x01,
			},
		},
		expected: []byte{0x01, 0x00, 0x00, 0x07, 0x00, 0x00, 0x00, 0x00, 0x81, 0x01, 0x04, 0x03f, 0x02, 0x01, 0xff},
	},
	{
		name: "MemoryRecall3",
		payload: &payload{
			Type:         _payloadTypeCommand,
			IsInquiry:    false,
			CategoryCode: _categoryCodeCamera1,
			Command:      _commandMemory,
			Args: []byte{
				_funcMemoryRecall,
				0x03,
			},
		},
		expected: []byte{0x01, 0x00, 0x00, 0x07, 0x00, 0x00, 0x00, 0x00, 0x81, 0x01, 0x04, 0x03f, 0x02, 0x03, 0xff},
	},
}

func TestPayloadMarshal(t *testing.T) {
	for _, tt := range _payloadMarshalTests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			b, err := tt.payload.MarshalBinary()
			if tt.expectedErr == nil {
				is.NoErr(err)
			} else {
				is.True(err.Error() == tt.expectedErr.Error())
			}

			t.Logf("%s: %# x", tt.name, b)
			is.Equal(b, tt.expected)
		})
	}
}
