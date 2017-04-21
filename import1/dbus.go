/*
Copyright 2015 CoreOS Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Integration with the systemd-importd API - See https://www.freedesktop.org/wiki/Software/systemd/importd/
package import1

import (
	"fmt"
	"os"
	"strconv"

	"github.com/godbus/dbus"
)

const (
	dbusInterface = "org.freedesktop.import1.Manager"
	dbusPath      = "/org/freedesktop/import1"
)

type Transfer struct {
	ID   uint
	Path dbus.ObjectPath
}

type Conn struct {
	conn   *dbus.Conn
	object dbus.BusObject
}

func New() (*Conn, error) {
	c := new(Conn)
	if err := c.initConnection(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Conn) initConnection() error {
	var err error
	c.conn, err = dbus.SystemBusPrivate()
	if err != nil {
		return err
	}
	methods := []dbus.Auth{dbus.AuthExternal(strconv.Itoa(os.Getuid()))}
	err = c.conn.Auth(methods)
	if err != nil {
		c.conn.Close()
		return err
	}
	err = c.conn.Hello()
	if err != nil {
		c.conn.Close()
		return err
	}
	c.object = c.conn.Object("org.freedesktop.import1", dbus.ObjectPath(dbusPath))
	return nil
}

func (c *Conn) transferImage(method string, args ...interface{}) (*Transfer, error) {
	result := c.object.Call(fmt.Sprintf("%s.%s", dbusInterface, method), 0, args...)
	if result.Err != nil {
		return nil, result.Err
	}
	transferID, ok := result.Body[0].(uint)
	if !ok {
		return nil, fmt.Errorf("unable to convert dbus response '%v' to uint", result.Body[0])
	}
	transferPath, ok := result.Body[1].(dbus.ObjectPath)
	if !ok {
		return nil, fmt.Errorf("unable to convert dbus response '%v' to dbus.ObjectPath", result.Body[1])
	}
	return &Transfer{
		ID:   transferID,
		Path: transferPath,
	}, nil
}

func (c *Conn) ImportTar(fd int, localName string, force, readOnly bool) (*Transfer, error) {
	return c.importImage("ImportTar", fd, localName, force, readOnly)
}

func (c *Conn) ImportRaw(fd int, localName string, force, readOnly bool) (*Transfer, error) {
	return c.importImage("ImportRaw", fd, localName, force, readOnly)
}
