//go:build plan9 && crypt

// Copyright (C) 2020 - 2022 iDigitalFlame
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//

package local

import (
	"os"

	"github.com/iDigitalFlame/xmt/util/crypt"
)

func sysID() []byte {
	if b, err := os.ReadFile(crypt.Get(87)); err == nil { // /etc/hostid
		return b
	}
	o, _ := output(crypt.Get(88)) // kenv -q smbios.system.uuid
	return o
}
func version() string {
	var (
		ok      bool
		b, n, v string
	)
	if m := release(); len(m) > 0 {
		b = m[crypt.Get(8)]                // ID
		if n, ok = m[crypt.Get(80)]; !ok { // PRETTY_NAME
			n = m[crypt.Get(81)] // NAME
		}
		if v, ok = m[crypt.Get(82)]; !ok { // VERSION_ID
			v = m[crypt.Get(83)] // VERSION
		}
	}
	if len(v) == 0 {
		v = crypt.Get(93) // plan9
	}
	switch {
	case len(n) == 0 && len(b) == 0 && len(v) == 0:
		return crypt.Get(93) // plan9
	case len(n) == 0 && len(b) > 0 && len(v) > 0:
		return crypt.Get(93) + " (" + v + ", " + b + ")" // plan9
	case len(n) == 0 && len(b) == 0 && len(v) > 0:
		return crypt.Get(93) + " (" + v + ")" // plan9
	case len(n) == 0 && len(b) > 0 && len(v) == 0:
		return crypt.Get(93) + " (" + b + ")" // plan9
	case len(n) > 0 && len(b) > 0 && len(v) > 0:
		return n + " (" + v + ", " + b + ")"
	case len(n) > 0 && len(b) == 0 && len(v) > 0:
		return n + " (" + v + ")"
	case len(n) > 0 && len(b) > 0 && len(v) == 0:
		return n + " (" + b + ")"
	}
	return crypt.Get(93) // plan9
}