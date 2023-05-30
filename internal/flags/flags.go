package flags

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const fileRestoreFlagTrue = "true"
const fileRestoreFlagFalse = "false"

type NetAddress struct {
	Host string
	Port int
}

type Poll struct {
	Interval time.Duration
}

type Report struct {
	Interval time.Duration
}

type FileRestore struct {
	Interval    fileInterval
	Path        filePath
	FileRestore restore
	Synchro     bool
	Restore     bool
	File        *os.File
}

type fileInterval struct {
	Data time.Duration
}

type filePath struct {
	Data string
}

type restore struct {
	Data bool
}

func (a *NetAddress) String() string {
	return fmt.Sprintf("%s:%s", a.Host, strconv.Itoa(a.Port))
}

func (a *NetAddress) Set(s string) error {
	hp := strings.Split(s, ":")
	if len(hp) != 2 {
		return fmt.Errorf("%s", "Need address in a form 'host:port' or ':port'")
	}
	port, err := strconv.Atoi(hp[1])
	if err != nil {
		return err
	}
	if port < 1025 || port > 65535 {
		return fmt.Errorf("%s", "Incorrect port assignment! Possible range is from 1025 to 65535")
	}
	hp[0] = strings.Trim(hp[0], " ")
	if len(hp[0]) > 0 {
		a.Host = hp[0]
	}

	hp[0] = strings.Trim(hp[0], " ")
	a.Port = port
	return nil
}

func (p *Poll) String() string {
	return fmt.Sprintf("%d", p.Interval/time.Second)
}

func (p *Poll) Set(s string) error {
	num, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	if num <= 0 {
		return fmt.Errorf("%s", "The value cannot be less than or equal to 0")
	}
	p.Interval = time.Duration(num) * time.Second
	return nil
}

func (r *Report) String() string {
	return fmt.Sprintf("%d", r.Interval/time.Second)
}

func (r *Report) Set(s string) error {
	num, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	if num <= 0 {
		return fmt.Errorf("%s", "The value cannot be less than or equal to 0")
	}
	r.Interval = time.Duration(num) * time.Second
	return nil
}

func (fi *fileInterval) String() string {
	return fmt.Sprintf("%d", fi.Data/time.Second)
}

func (fi *fileInterval) Set(s string) error {
	num, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	if num <= 0 {
		return fmt.Errorf("%s", "The value cannot be less than or equal to 0")
	}
	fi.Data = time.Duration(num) * time.Second
	return nil
}

func (fi *filePath) String() string {
	return fi.Data
}

func (fi *filePath) Set(s string) error {

	s = strings.TrimSpace(s)
	if len(s) != 0 {
		fi.Data = s
	}
	return nil
}

func (fi *restore) String() string {
	return fmt.Sprintf("%v", fi.Data)
}

func (fi *restore) Set(s string) error {

	if s == fileRestoreFlagTrue {
		fi.Data = true
	} else if s == fileRestoreFlagFalse {
		fi.Data = false
	} else {
		return errors.New("you can use only 'true' or 'false'")
	}
	return nil
}

func NewFileRestore() *FileRestore {
	return &FileRestore{
		Interval: fileInterval{
			Data: 300 * time.Second,
		},
		Path: filePath{
			Data: "/tmp/metrics-db.json",
		},
		FileRestore: restore{
			Data: true,
		},
	}
}

func Server(addr *NetAddress, file *FileRestore) {
	_ = flag.Value(addr)
	_ = flag.Value(&file.Interval)
	_ = flag.Value(&file.Path)
	_ = flag.Value(&file.FileRestore)
	// проверка реализации
	flag.Var(addr, "a", "Net address 'host:port' or ':port'")
	flag.Var(&file.Interval, "i", `time interval, in seconds,
after which the current server data is saved to disk
(by default 300 seconds, value 0 makes the recording synchronous)`)
	flag.Var(&file.Path, "f", `full name of the file where the current values are saved
(by default /tmp/metrics-db.json, an empty value disables the disk write function)`)
	flag.Var(&file.FileRestore, "r", `Boolean value (true/false) that determines whether or not
to load previously saved values from the specified file
when the server starts up (the default is true)`)
}

func Agent(addr *NetAddress, poll *Poll, report *Report) {
	_ = flag.Value(addr)
	_ = flag.Value(poll)
	_ = flag.Value(report)
	// проверка реализации
	flag.Var(addr, "a", "Net address 'host:port' or ':port'")
	flag.Var(poll, "p", "Poll interval in seconds")
	flag.Var(report, "r", "Report interval in seconds")
}
