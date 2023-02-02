package visca

import (
	"context"
	"fmt"
)

const (
	_funcZoomStop = 0x00
	_funcZoomTele = 0x02
	_funcZoomWide = 0x03
)

func (c *Camera) ZoomStop(ctx context.Context) error {
	return c.zoom(ctx, _funcZoomStop)
}

func (c *Camera) ZoomTele(ctx context.Context) error {
	return c.zoom(ctx, _funcZoomTele)
}

func (c *Camera) ZoomWide(ctx context.Context) error {
	return c.zoom(ctx, _funcZoomWide)
}

func (c *Camera) zoom(ctx context.Context, speedDir byte) error {
	if speedDir != _funcZoomStop {
		if err := c.ZoomStop(ctx); err != nil {
			return fmt.Errorf("unable to stop zoom: %w", err)
		}
	}

	payload := payload{
		Type:         _payloadTypeCommand,
		IsInquiry:    false,
		CategoryCode: _categoryCodeCamera1,
		Command:      _commandZoom,
		Args: []byte{
			speedDir,
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
