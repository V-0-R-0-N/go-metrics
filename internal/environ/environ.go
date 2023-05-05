package environ

import (
	fl "github.com/V-0-R-0-N/go-metrics.git/internal/flags"
	"os"
)

func parseAddr(addr *fl.NetAddress) error {
	address, ok := os.LookupEnv("ADDRESS")
	if ok {
		err := addr.Set(address)
		if err != nil {
			return err
		}
	}
	return nil
}

func parseAgentPollReport(poll *fl.Poll, report *fl.Report) error {

	if p, ok := os.LookupEnv("POLL_INTERVAL"); ok {
		err := poll.Set(p)
		if err != nil {
			return err
		}
	}
	if r, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
		err := report.Set(r)
		if err != nil {
			return err
		}
	}
	return nil
}

func Server(addr *fl.NetAddress) error {
	err := parseAddr(addr)
	if err != nil {
		return err
	}
	return nil
}

func Agent(addr *fl.NetAddress, poll *fl.Poll, report *fl.Report) error {
	err := parseAddr(addr)
	if err != nil {
		return err
	}
	err = parseAgentPollReport(poll, report)
	if err != nil {
		return err
	}
	return nil
}
