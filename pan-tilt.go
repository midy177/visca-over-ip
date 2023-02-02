package visca

import (
	"context"
	"errors"
)

const (
	TiltDirectionUp   = 0x01
	TiltDirectionDown = 0x02
	TiltDirectionStop = 0x03

	PanDirectionLeft  = 0x01
	PanDirectionRight = 0x02
	PanDirectionStop  = 0x03

	PanTiltSpeedMin = 0x00
	PanTiltSpeedMax = 0x18
)

var (
	ErrInvalidPanTiltSpeed     = errors.New("speed out of range")
	ErrInvalidPanTiltDirection = errors.New("invalid pan/tilt direction")
)

func (c *Camera) TiltUp(ctx context.Context, speed byte) error {
	return c.PanTiltDrive(ctx, PanDirectionStop, TiltDirectionUp, PanTiltSpeedMin, speed)
}

func (c *Camera) TiltDown(ctx context.Context, speed byte) error {
	return c.PanTiltDrive(ctx, PanDirectionStop, TiltDirectionDown, PanTiltSpeedMin, speed)
}

func (c *Camera) PanLeft(ctx context.Context, speed byte) error {
	return c.PanTiltDrive(ctx, PanDirectionLeft, TiltDirectionStop, speed, PanTiltSpeedMin)
}

func (c *Camera) PanRight(ctx context.Context, speed byte) error {
	return c.PanTiltDrive(ctx, PanDirectionRight, TiltDirectionStop, speed, PanTiltSpeedMin)
}

func (c *Camera) PanTiltStop(ctx context.Context) error {
	return c.PanTiltDrive(ctx, PanDirectionStop, TiltDirectionStop, PanTiltSpeedMin, PanTiltSpeedMin)
}

func (c *Camera) PanTiltDrive(ctx context.Context, panDir, tiltDir, panSpeed, tiltSpeed byte) error {
	switch {
	case panSpeed < PanTiltSpeedMin || panSpeed > PanTiltSpeedMax:
		return ErrInvalidPanTiltSpeed
	case tiltSpeed < PanTiltSpeedMin || tiltSpeed > PanTiltSpeedMax:
		return ErrInvalidPanTiltSpeed
	case panDir != PanDirectionLeft && panDir != PanDirectionRight && panDir != PanDirectionStop:
		return ErrInvalidPanTiltDirection
	case tiltDir != TiltDirectionUp && tiltDir != TiltDirectionDown && tiltDir != TiltDirectionStop:
		return ErrInvalidPanTiltDirection
	}

	payload := payload{
		Type:         _payloadTypeCommand,
		IsInquiry:    false,
		CategoryCode: _categoryCodePanTilter,
		Command:      _commandPanTiltDrive,
		Args: []byte{
			panSpeed,
			tiltSpeed,
			panDir,
			tiltDir,
		},
	}

	resp, err := c.sendPayload(ctx, payload)
	if err != nil {
		return err
	}

	if !resp.IsAck() {
		return resp.Error()
	}

	return nil
}
