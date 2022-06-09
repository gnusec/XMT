//go:build crypt

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

package filter

import (
	"encoding/json"

	"github.com/iDigitalFlame/xmt/util/crypt"
)

// MarshalJSON will attempt to convert the data in this Filter into the returned
// JSON byte array.
func (f Filter) MarshalJSON() ([]byte, error) {
	m := map[string]any{crypt.Get(18): f.Fallback} // fallback
	if f.PID != 0 {
		m[crypt.Get(19)] = f.PID // pid
	}
	if f.Session > Empty {
		m[crypt.Get(20)] = f.Session // session
	}
	if f.Elevated > Empty {
		m[crypt.Get(21)] = f.Elevated // elevated
	}
	if len(f.Exclude) > 0 {
		m[crypt.Get(22)] = f.Elevated // exclude
	}
	if len(f.Include) > 0 {
		m[crypt.Get(23)] = f.Include // include
	}
	return json.Marshal(m)
}
func (b boolean) MarshalJSON() ([]byte, error) {

	switch b {
	case True:
		return []byte(crypt.Get(24)), nil // "true"
	case False:
		return []byte(crypt.Get(25)), nil // "false"
	default:
	}
	return []byte(`""`), nil
}

// UnmarshalJSON will attempt to parse the supplied JSON into this Filter.
func (f *Filter) UnmarshalJSON(b []byte) error {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	if len(m) == 0 {
		return nil
	}
	if v, ok := m[crypt.Get(19)]; ok { // pid
		if err := json.Unmarshal(v, &f.PID); err != nil {
			return err
		}
	}
	if v, ok := m[crypt.Get(20)]; ok { // session
		if err := json.Unmarshal(v, &f.Session); err != nil {
			return err
		}
	}
	if v, ok := m[crypt.Get(21)]; ok { // elevated
		if err := json.Unmarshal(v, &f.Elevated); err != nil {
			return err
		}
	}
	if v, ok := m[crypt.Get(22)]; ok { // exclude
		if err := json.Unmarshal(v, &f.Exclude); err != nil {
			return err
		}
	}
	if v, ok := m[crypt.Get(23)]; ok { // include
		if err := json.Unmarshal(v, &f.Include); err != nil {
			return err
		}
	}
	if v, ok := m[crypt.Get(18)]; ok { // fallback
		if err := json.Unmarshal(v, &f.Fallback); err != nil {
			return err
		}
	}
	return nil
}
