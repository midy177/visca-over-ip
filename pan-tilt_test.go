package visca

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestDriveCircle(t *testing.T) {
	cam := New("10.13.34.8:52381", WithLogger(testLogger{t}))

	err := cam.TiltUp(context.Background(), 0x0e)
	if err != nil {
		t.Fatalf("failed to pan tilt up: %s\n", err)
	}

	time.Sleep(1500 * time.Millisecond)

	err = cam.PanRight(context.Background(), 0x0b)
	if err != nil {
		t.Fatalf("failed to pan right up: %s\n", err)
	}

	time.Sleep(1500 * time.Millisecond)

	err = cam.TiltDown(context.Background(), 0x0e)
	if err != nil {
		t.Fatalf("failed to pan right down: %s\n", err)
	}

	time.Sleep(1500 * time.Millisecond)

	err = cam.PanLeft(context.Background(), 0x0b)
	if err != nil {
		t.Fatalf("failed to pan right left: %s\n", err)
	}

	time.Sleep(1500 * time.Millisecond)

	err = cam.PanTiltStop(context.Background())
	if err != nil {
		t.Fatalf("failed to pan stop: %s\n", err)
	}
}

func TestDriveErrors(t *testing.T) {
	cam := New("")
	ctx := context.Background()

	if err := cam.TiltUp(ctx, 0x19); !errors.Is(err, ErrInvalidPanTiltSpeed) {
		t.Fatalf("expected %v, got %v", ErrInvalidPanTiltSpeed, err)
	}

	if err := cam.TiltDown(ctx, 0x22); !errors.Is(err, ErrInvalidPanTiltSpeed) {
		t.Fatalf("expected %v, got %v", ErrInvalidPanTiltSpeed, err)
	}

	if err := cam.PanLeft(ctx, 0x59); !errors.Is(err, ErrInvalidPanTiltSpeed) {
		t.Fatalf("expected %v, got %v", ErrInvalidPanTiltSpeed, err)
	}

	if err := cam.PanRight(ctx, 0xff); !errors.Is(err, ErrInvalidPanTiltSpeed) {
		t.Fatalf("expected %v, got %v", ErrInvalidPanTiltSpeed, err)
	}

	if err := cam.PanTiltDrive(ctx, PanDirectionLeft, 0x04, PanTiltSpeedMin, PanTiltSpeedMin); !errors.Is(err, ErrInvalidPanTiltDirection) {
		t.Fatalf("expected %v, got %v", ErrInvalidPanTiltDirection, err)
	}

	if err := cam.PanTiltDrive(ctx, 0x05, TiltDirectionDown, PanTiltSpeedMin, PanTiltSpeedMin); !errors.Is(err, ErrInvalidPanTiltDirection) {
		t.Fatalf("expected %v, got %v", ErrInvalidPanTiltDirection, err)
	}
}
