package visca

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/byuoitav/connpool"
	"log"
	"net"
	"time"
	"unsafe"
)

type Device struct {
	HardwareAddr net.HardwareAddr
	Model        string
	SoftVersion  string
	IP           net.IP
	MASK         net.IPMask
	Gateway      net.IP
	Name         string
	Write        bool
}

func (d *Device) New(opts ...Option) *Camera {
	options := options{
		ttl:    _defaultTTL,
		delay:  _defaultDelay,
		dialer: net.Dialer{},
	}

	for _, o := range opts {
		o.apply(&options)
	}
	cam := &Camera{
		address: bytes2Str(d.IP.To4()) + ":52381",
		pool: &connpool.Pool{
			TTL:    options.ttl,
			Delay:  options.delay,
			Logger: options.logger,
		},
		logger: options.logger,
		dialer: options.dialer,
	}

	cam.pool.NewConnection = func(ctx context.Context) (net.Conn, error) {
		return cam.dialer.DialContext(ctx, "udp", cam.address)
	}

	return cam
}
func Discover() []Device {
	bCast, err := GetBroadcast("以太网")
	if err != nil {
		log.Fatalln(err)
	}
	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	dstAddr := &net.UDPAddr{IP: bCast, Port: 52380}
	conn, err := net.ListenUDP("udp", srcAddr)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	dsc := []byte{0x02, 0x45, 0x4E, 0x51, 0x3A, 0x6E, 0x65, 0x74, 0x77, 0x6F, 0x72, 0x6B, 0xFF, 0x03}
	_, err = conn.WriteToUDP(dsc, dstAddr)
	if err != nil {
		fmt.Println(err)
	}
	tk := time.NewTimer(time.Second * 5)
	var devices []Device
forSel:
	for {
		select {
		case <-tk.C:
			break forSel
		default:
			data := make([]byte, 1024)
			n, _, err := conn.ReadFrom(data)
			if err != nil {
				//fmt.Println(err)
				//opError, ok := err.(*net.OpError)
				//if ok {
				//	log.Println(opError.Err.Error())
				//}
				break forSel
			}
			device := ParseDevice(data[1 : n-2])
			if device != nil {
				devices = append(devices, *device)
			}
		}
	}
	return devices
}
func ParseDevice(pStr []byte) *Device {
	var device = new(Device)
	arr := bytes.Split(pStr, []byte{255})
	for _, v := range arr {
		macIndex := bytes.Index(v, []byte("MAC:"))
		if macIndex != -1 {
			if hw, err := net.ParseMAC(bytes2Str(v[macIndex+4:])); err == nil {
				device.HardwareAddr = hw
			}
			continue
		}
		modelIndex := bytes.Index(v, []byte("MODEL:"))
		if modelIndex != -1 {
			device.Model = bytes2Str(v[modelIndex+6:])
			continue
		}
		softVersionIndex := bytes.Index(v, []byte("SOFTVERSION:"))
		if softVersionIndex != -1 {
			device.SoftVersion = bytes2Str(v[softVersionIndex+12:])
			continue
		}
		ipIndex := bytes.Index(v, []byte("IPADR:"))
		if ipIndex != -1 {
			device.IP = net.ParseIP(bytes2Str(v[ipIndex+6:]))
			continue
		}
		maskIndex := bytes.Index(v, []byte("MASK:"))
		if maskIndex != -1 {
			device.MASK = v[maskIndex+5:]
			continue
		}
		gatewayIndex := bytes.Index(v, []byte("GATEWAY:"))
		if gatewayIndex != -1 {
			device.Gateway = net.ParseIP(bytes2Str(v[gatewayIndex+8:]))
			continue
		}
		nameIndex := bytes.Index(v, []byte("NAME:"))
		if nameIndex != -1 {
			device.Name = bytes2Str(v[nameIndex+5:])
			continue
		}
		writeIndex := bytes.Index(v, []byte("WRITE:"))
		if writeIndex != -1 {
			if bytes.Contains(v, []byte("on")) {
				device.Write = true
			}
			continue
		}
	}
	return device
}
func bytes2Str(slice []byte) string {
	return *(*string)(unsafe.Pointer(&slice))
}
func GetBroadcast(name string) ([]byte, error) {
	ifn, err := GetInterfaceIpv4Net(name)
	if err != nil {
		return nil, err
	}
	//mask := net.CIDRMask(20, 32)
	mask := ifn.Mask
	//ip := net.IP([]byte{140, 45, 32, 0})
	ip := ifn.IP.To4()
	broadcast := net.IP(make([]byte, 4))
	for i := range ip {
		broadcast[i] = ip[i] | ^mask[i]
	}
	return broadcast, err
}

func GetInterfaceIpv4Net(interfaceName string) (*net.IPNet, error) {
	var (
		ief     *net.Interface
		addrs   []net.Addr
		ipv4Net *net.IPNet
		ok      bool
		err     error
	)
	if ief, err = net.InterfaceByName(interfaceName); err != nil { // get interface
		return nil, err
	}
	if addrs, err = ief.Addrs(); err != nil { // get addresses
		return nil, err
	}
	for _, v := range addrs { // get ipv4Net address

		if ipv4Net, ok = v.(*net.IPNet); ok {
			if ipv4Net.IP.To4() != nil {
				break
			}
		}
	}
	if ipv4Net == nil {
		return nil, errors.New(fmt.Sprintf("interface %s don't have an ipv4Net address\n", interfaceName))
	}
	return ipv4Net, nil
}
