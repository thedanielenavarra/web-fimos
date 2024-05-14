package main

import (
	"fmt"
	"log"

	"github.com/godbus/dbus"
)

// Functions that interact with the firewalld d-bus interface to create new firewall rules to open or close ports
// to specific IPs

func addRule(sourceIP string, destinationPort int, networkInterface string, protocol string) error {
	conn, err := dbus.SystemBus()
	if err != nil {
		log.Fatal("Couln't create system dbus: ", err)
	}
	firewalld := conn.Object("org.fedoraproject.FirewallD1", "/org/fedoraproject/FirewallD1")

	call := firewalld.Call("org.fedoraproject.FirewallD1.zone.addRichRule", 0, networkInterface, fmt.Sprintf("rule family='ipv4' source address='%s' port protocol='%s' port='%d' accept", sourceIP, protocol, destinationPort), 0)

	if call.Err != nil {
		log.Fatal("Error during dbus call:", call.Err)
		return call.Err
	}
	return nil
}

func removeRule(sourceIP string, destinationPort int, networkInterface string, protocol string) error {
	conn, err := dbus.SystemBus()
	if err != nil {
		log.Fatal("Couln't create system dbus: ", err)
	}
	firewalld := conn.Object("org.fedoraproject.FirewallD1", "/org/fedoraproject/FirewallD1")
	fmt.Println("Removing rule for ", sourceIP, destinationPort, networkInterface, protocol)

	call := firewalld.Call("org.fedoraproject.FirewallD1.zone.removeRichRule", 0, networkInterface, fmt.Sprintf("rule family='ipv4' source address='%s' port protocol='%s' port='%d' accept", sourceIP, protocol, destinationPort))
	
	if call.Err != nil {
		log.Fatal("Error during dbus call:", call.Err)
		return call.Err
	}
	return nil
}