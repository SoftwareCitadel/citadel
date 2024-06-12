package dockerDriver

import "net"

func (d *DockerDriver) setIPs() error {
	interfaces, err := net.Interfaces()
	if err != nil {
		return err
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return err
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			if ip4 := ip.To4(); ip4 != nil {
				if d.ipv4 == "" {
					d.ipv4 = ip4.String()
				}
			} else if ip16 := ip.To16(); ip16 != nil {
				if d.ipv6 == "" {
					d.ipv6 = ip16.String()
				}
			}
		}
	}

	return nil
}
