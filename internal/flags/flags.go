package flags

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"
)

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

func (a NetAddress) String() string {
	return a.Host + ":" + strconv.Itoa(a.Port)
}

func (a *NetAddress) Set(s string) error {
	hp := strings.Split(s, ":")
	if len(hp) != 2 {
		return fmt.Errorf("%s\n", "Need address in a form 'host:port' or ':port'")
	}
	port, err := strconv.Atoi(hp[1])
	if err != nil {
		return err
	}
	if port < 1025 || port > 65535 {
		return fmt.Errorf("%s\n", "Incorrect port assignment! Possible range is from 1025 to 65535")
	}
	hp[0] = strings.Trim(hp[0], " ")
	if len(hp[0]) > 0 {
		a.Host = hp[0]
	}

	strings.Trim(hp[0], " ")
	a.Port = port
	return nil
}

func (p Poll) String() string {
	return fmt.Sprintf("%d\n", p.Interval/time.Second)
}

func (p *Poll) Set(s string) error {
	num, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	if num <= 0 {
		return fmt.Errorf("%s\n", "The value cannot be less than or equal to 0")
	}
	p.Interval = time.Duration(num)
	return nil
}

func (r Report) String() string {
	return fmt.Sprintf("%d\n", r.Interval/time.Second)
}

func (r *Report) Set(s string) error {
	num, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	if num <= 0 {
		return fmt.Errorf("%s\n", "The value cannot be less than or equal to 0")
	}
	r.Interval = time.Duration(num)
	return nil
}

func Server(addr *NetAddress) {
	_ = flag.Value(addr)
	// проверка реализации
	flag.Var(addr, "a", "Net address 'host:port' or ':port'")
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
