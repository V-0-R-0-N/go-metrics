package environ

import (
	"os"

	"github.com/V-0-R-0-N/go-metrics.git/internal/flags"
)

type allFlags interface {
	Set(string) error
}

func parseFlag(data allFlags, name string) error {
	if value, ok := os.LookupEnv(name); ok {
		if err := data.Set(value); err != nil {
			return err
		}
	}
	return nil
}

func Server(addr *flags.NetAddress, fileRestore *flags.FileRestore) error {

	if err := parseFlag(addr, "ADDRESS"); err != nil {
		return err
	}
	if err := parseFlag(&fileRestore.Interval, "STORE_INTERVAL"); err != nil {
		return err
	}
	if err := parseFlag(&fileRestore.Path, "FILE_STORAGE_PATH"); err != nil {
		return err
	}
	if err := parseFlag(&fileRestore.FileRestore, "RESTORE"); err != nil {
		return err
	}
	return nil
}

func Agent(addr *flags.NetAddress, poll *flags.Poll, report *flags.Report) error {

	err := parseFlag(addr, "ADDRESS")
	if err != nil {
		return err
	}
	if err = parseFlag(poll, "POLL_INTERVAL"); err != nil {
		return err
	}
	if err = parseFlag(report, "REPORT_INTERVAL"); err != nil {
		return err
	}

	return nil
}
