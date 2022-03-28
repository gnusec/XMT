//go:build !implant

package result

import (
	"io"
	"io/fs"
	"time"

	"github.com/iDigitalFlame/xmt/c2"
	"github.com/iDigitalFlame/xmt/cmd"
	"github.com/iDigitalFlame/xmt/com"
	"github.com/iDigitalFlame/xmt/data"
	"github.com/iDigitalFlame/xmt/device/regedit"
)

type fileInfo struct {
	mod  time.Time
	name string
	size int64
	mode fs.FileMode
}

func (f fileInfo) Size() int64 {
	return f.size
}
func (f fileInfo) IsDir() bool {
	return f.mode.IsDir()
}
func (f fileInfo) Name() string {
	return f.name
}
func (fileInfo) Sys() interface{} {
	return nil
}
func (f fileInfo) Mode() fs.FileMode {
	return f.mode
}
func (f fileInfo) ModTime() time.Time {
	return f.mod
}

// Pwd will parse the RvResult Packet from a MvPwd task.
//
// The return result is the current directory the client is located in.
//
// This function returns an error if any reading errors occur, the Packet is not
// in the expected format or the Packet is nil or empty.
func Pwd(n *com.Packet) (string, error) {
	if n == nil || n.Empty() {
		return "", c2.ErrMalformedPacket
	}
	return n.StringVal()
}

// Spawn will parse the RvResult Packet from a MvSpawn task.
//
// The return result is the new PID of the resulting Spawn operation.
//
// This function returns an error if any reading errors occur, the Packet is not
// in the expected format or the Packet is nil or empty.
func Spawn(n *com.Packet) (uint32, error) {
	if n == nil || n.Empty() {
		return 0, c2.ErrMalformedPacket
	}
	return n.Uint32()
}

// CheckDLL will parse the RvResult Packet from a TvCheckDLL task.
//
// The return result is true if the DLL provided is NOT hooked. A return value
// of false indicates that the DLL memory space differs from the on-disk value,
// which is an indicator of hooks.
//
// This function returns an error if any reading errors occur, the Packet is not
// in the expected format or the Packet is nil or empty.
func CheckDLL(n *com.Packet) (bool, error) {
	if n == nil || n.Empty() {
		return false, c2.ErrMalformedPacket
	}
	return n.Bool()
}

// Mounts will parse the RvResult Packet from a MvMounts task.
//
// The return result is a string list of all the exposed mount points on the
// client (drive letters on Windows).
//
// This function returns an error if any reading errors occur, the Packet is not
// in the expected format or the Packet is nil or empty.
func Mounts(n *com.Packet) ([]string, error) {
	if n == nil || n.Empty() {
		return nil, c2.ErrMalformedPacket
	}
	var s []string
	return s, data.ReadStringList(n, &s)
}

// Ls will parse the RvResult Packet from a MvList task.
//
// The return result is a slice of FileInfo interfaces that will return the
// data of the directory targed.
//
// This function returns an error if any reading errors occur, the Packet is not
// in the expected format or the Packet is nil or empty.
func Ls(n *com.Packet) ([]fs.FileInfo, error) {
	if n == nil || n.Empty() {
		return nil, c2.ErrMalformedPacket
	}
	c, err := n.Uint32()
	if err != nil || c == 0 {
		return nil, err
	}
	e := make([]fs.FileInfo, c)
	for i := range e {
		var v fileInfo
		if err = n.ReadString(&v.name); err != nil {
			return nil, err
		}
		if err = n.ReadUint32((*uint32)(&v.mode)); err != nil {
			return nil, err
		}
		if err = n.ReadInt64(&v.size); err != nil {
			return nil, err
		}
		t, err := n.Int64()
		if err != nil {
			return nil, err
		}
		v.mod = time.Unix(t, 0)
		e[i] = v
	}
	return e, nil
}

// Pull will parse the RvResult Packet from a TvPull task.
//
// The return result is the expended full file path on the host as a string, and
// the resulting count of bytes written to disk.
//
// This function returns an error if any reading errors occur, the Packet is not
// in the expected format or the Packet is nil or empty.
func Pull(n *com.Packet) (string, uint64, error) {
	return Upload(n)
}

// ScreenShot will parse the RvResult Packet from a TvScreenShot task.
//
// The return result is a Reader with the resulting screenshot data encoded as
// a png image inside. (This can be directly written to disk as a png file).
//
// This function returns an error if any reading errors occur, the Packet is not
// in the expected format or the Packet is nil or empty.
func ScreenShot(n *com.Packet) (io.Reader, error) {
	return ProcessDump(n)
}

// ProcessDump will parse the RvResult Packet from a TvProcDump task.
//
// The return result is a Reader with the resulting dump data inside.
//
// This function returns an error if any reading errors occur, the Packet is not
// in the expected format or the Packet is nil or empty.
func ProcessDump(n *com.Packet) (io.Reader, error) {
	if n == nil || n.Empty() {
		return nil, c2.ErrMalformedPacket
	}
	return n, nil
}

// Upload will parse the RvResult Packet from a TvUpload task.
//
// The return result is the expended full file path on the host as a string, and
// the resulting count of bytes written to disk.
//
// This function returns an error if any reading errors occur, the Packet is not
// in the expected format or the Packet is nil or empty.
func Upload(n *com.Packet) (string, uint64, error) {
	if n == nil || n.Empty() {
		return "", 0, c2.ErrMalformedPacket
	}
	var (
		s   string
		c   uint64
		err = n.ReadString(&s)
	)
	if err != nil {
		return "", 0, err
	}
	if err = n.ReadUint64(&c); err != nil {
		return "", 0, err
	}
	return s, c, nil
}

// DLL will parse the RvResult Packet from a TvDLL task.
//
// The return result is a handle to the memory location of the DLL (as a
// uintptr), the resulting PID of the DLL "host" and the exit code of the
// primary thread (if wait was specified, otherwise this is zero).
//
// This function returns an error if any reading errors occur, the Packet is not
// in the expected format or the Packet is nil or empty.
func DLL(n *com.Packet) (uintptr, uint32, int32, error) {
	return Assembly(n)
}

// ProcessList will parse the RvResult Packet from a TvProcList task.
//
// The return result is a slice of 'cmd.ProcessInfo' structs that will indicate
// the current processes running on the target device.
//
// This function returns an error if any reading errors occur, the Packet is not
// in the expected format or the Packet is nil or empty.
func ProcessList(n *com.Packet) ([]cmd.ProcessInfo, error) {
	if n == nil || n.Empty() {
		return nil, c2.ErrMalformedPacket
	}
	c, err := n.Uint32()
	if err != nil {
		return nil, err
	}
	e := make([]cmd.ProcessInfo, c)
	for i := range e {
		if err = n.ReadUint32(&e[i].PID); err != nil {
			return nil, err
		}
		if err = n.ReadUint32(&e[i].PPID); err != nil {
			return nil, err
		}
		if err = n.ReadString(&e[i].Name); err != nil {
			return nil, err
		}
	}
	return e, nil
}

// Registry will parse the RvResult Packet from a TvRegistry task.
//
// The return result is dependent on the resulting operation. If the result is
// from a 'RegLs' or 'RegGet' operation, this will return the resulting entries
// found (only one entry if this was a Get operation).
//
// The boolean value will return true if the result was a valid registry command
// that returns no output, such as a Set operation.
//
// This function returns an error if any reading errors occur, the Packet is not
// in the expected format or the Packet is nil or empty.
func Registry(n *com.Packet) ([]regedit.Entry, bool, error) {
	if n == nil || n.Empty() {
		return nil, false, c2.ErrMalformedPacket
	}
	o, err := n.Uint8()
	if err != nil {
		return nil, false, c2.ErrMalformedPacket
	}
	if o > 1 {
		return nil, o < 13, nil
	}
	var c uint32
	if o == 0 {
		if err = n.ReadUint32(&c); err != nil {
			return nil, false, err
		}
		if c == 0 {
			return nil, false, nil
		}
	} else {
		c = 1
	}
	r := make([]regedit.Entry, c)
	for i := range r {
		if err = r[i].UnmarshalStream(n); err != nil {
			return nil, false, err
		}
	}
	return r, true, nil
}

// Assembly will parse the RvResult Packet from a TvAssembly task.
//
// The return result is a handle to the memory location of the Assembly code (as
// a uintptr), the resulting PID of the Assembly "host" and the exit code of the
// primary thread (if wait was specified, otherwise this is zero).
//
// This function returns an error if any reading errors occur, the Packet is not
// in the expected format or the Packet is nil or empty.
func Assembly(n *com.Packet) (uintptr, uint32, int32, error) {
	if n == nil || n.Empty() {
		return 0, 0, 0, c2.ErrMalformedPacket
	}
	var (
		h   uint64
		p   uint32
		x   int32
		err = n.ReadUint64(&h)
	)
	if err != nil {
		return 0, 0, 0, err
	}
	if err = n.ReadUint32(&p); err != nil {
		return 0, 0, 0, err
	}
	if err = n.ReadInt32(&x); err != nil {
		return 0, 0, 0, err
	}
	return uintptr(h), p, x, nil
}

// Zombie will parse the RvResult Packet from a TvZombie task.
//
// The return result is the spawned PID of the new process and the resulting
// exit code and Stdout/Stderr data (if wait was specified, otherwise this the
// return code is zero and the reader will be empty).
//
// This function returns an error if any reading errors occur, the Packet is not
// in the expected format or the Packet is nil or empty.
func Zombie(n *com.Packet) (uint32, int32, io.Reader, error) {
	return Process(n)
}

// Process will parse the RvResult Packet from a TvExecute task.
//
// The return result is the spawned PID of the new process and the resulting
// exit code and Stdout/Stderr data (if wait was specified, otherwise this the
// return code is zero and the reader will be empty).
//
// This function returns an error if any reading errors occur, the Packet is not
// in the expected format or the Packet is nil or empty.
func Process(n *com.Packet) (uint32, int32, io.Reader, error) {
	if n == nil || n.Empty() {
		return 0, 0, nil, c2.ErrMalformedPacket
	}
	var (
		p   uint32
		x   int32
		err = n.ReadUint32(&p)
	)
	if err != nil {
		return 0, 0, nil, err
	}
	if err = n.ReadInt32(&x); err != nil {
		return 0, 0, nil, err
	}
	return p, x, n, nil
}

// PullExec will parse the RvResult Packet from a TvPullExecute task.
//
// The return result is the spawned PID of the new process and the resulting
// exit code and Stdout/Stderr data (if wait was specified, otherwise this the
// return code is zero and the reader will be empty).
//
// This function returns an error if any reading errors occur, the Packet is not
// in the expected format or the Packet is nil or empty.
func PullExec(n *com.Packet) (uint32, int32, io.Reader, error) {
	return Process(n)
}

// Download will parse the RvResult Packet from a TvDownload task.
//
// The return result is the expended full file path on the host as a string,
// a boolean representing if the path requested is a directory (true if the path
// is a directory, false otherwise), the size of the data in bytes (zero if the
// target is a directory) and a reader with the resulting file data (empty if
// the target is a directory).
//
// This function returns an error if any reading errors occur, the Packet is not
// in the expected format or the Packet is nil or empty.
func Download(n *com.Packet) (string, bool, uint64, io.Reader, error) {
	if n == nil || n.Empty() {
		return "", false, 0, nil, c2.ErrMalformedPacket
	}
	var (
		s   string
		d   bool
		c   uint64
		err = n.ReadString(&s)
	)
	if err != nil {
		return "", false, 0, nil, err
	}
	if err = n.ReadBool(&d); err != nil {
		return "", false, 0, nil, err
	}
	if err = n.ReadUint64(&c); err != nil {
		return "", false, 0, nil, err
	}
	return s, d, c, n, nil
}