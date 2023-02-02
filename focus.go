package visca

import "context"

const (
	_commandFocusModel  command = 0x38
	_commandFocusAction command = 0x08
)

type focusMode byte

const (
	FocusAuto       focusMode = 0x02
	FocusManual     focusMode = 0x03
	FocusModeToggle focusMode = 0x10
)

const (
	_focusStop         = 0x00
	_focusFarStandard  = 0x02
	_focusNearStandard = 0x03
)

func (c *Camera) SetFocusMode(ctx context.Context, mode focusMode) error {
	return c.changeFocusMode(ctx, byte(mode))
}
func (c *Camera) changeFocusMode(ctx context.Context, args byte) error {
	payload := payload{
		Type:         _payloadTypeCommand,
		IsInquiry:    false,
		CategoryCode: _categoryCodeCamera1,
		Command:      _commandFocusModel,
		Args: []byte{
			args,
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

func (c *Camera) FocusStop(ctx context.Context) error {
	return c.doFocusAction(ctx, _focusStop)
}

// FocusFar speed=0x20(low) - 0x27(high),other value are standard speed.
func (c *Camera) FocusFar(ctx context.Context, speed byte) error {
	if speed < 0x20 && speed > 0x27 {
		speed = 0x02
	}
	return c.doFocusAction(ctx, speed)
}

// FocusNear speed=0x20(low) - 0x27(high),other value are standard speed.
func (c *Camera) FocusNear(ctx context.Context, speed byte) error {
	if speed < 0x30 && speed > 0x37 {
		speed = 0x03
	}
	return c.doFocusAction(ctx, speed)
}
func (c *Camera) doFocusAction(ctx context.Context, args byte) error {
	payload := payload{
		Type:         _payloadTypeCommand,
		IsInquiry:    false,
		CategoryCode: _categoryCodeCamera1,
		Command:      _commandFocusAction,
		Args: []byte{
			args,
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
