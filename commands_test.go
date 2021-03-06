package main

import (
	"io/ioutil"
	"os/exec"
	"testing"

	drivers "github.com/docker/machine/drivers"
	"github.com/docker/machine/state"
)

type FakeDriver struct {
	MockState state.State
}

func (d *FakeDriver) DriverName() string {
	return "fakedriver"
}

func (d *FakeDriver) SetConfigFromFlags(flags drivers.DriverOptions) error {
	return nil
}

func (d *FakeDriver) GetURL() (string, error) {
	return "", nil
}

func (d *FakeDriver) GetIP() (string, error) {
	return "", nil
}

func (d *FakeDriver) GetState() (state.State, error) {
	return d.MockState, nil
}

func (d *FakeDriver) Create() error {
	return nil
}

func (d *FakeDriver) Remove() error {
	return nil
}

func (d *FakeDriver) Start() error {
	return nil
}

func (d *FakeDriver) Stop() error {
	return nil
}

func (d *FakeDriver) Restart() error {
	return nil
}

func (d *FakeDriver) Kill() error {
	return nil
}

func (d *FakeDriver) Upgrade() error {
	return nil
}

func (d *FakeDriver) StartDocker() error {
	return nil
}

func (d *FakeDriver) StopDocker() error {
	return nil
}

func (d *FakeDriver) GetDockerConfigDir() string {
	return ""
}

func (d *FakeDriver) GetSSHCommand(args ...string) (*exec.Cmd, error) {
	return &exec.Cmd{}, nil
}

func TestGetHostState(t *testing.T) {
	storePath, err := ioutil.TempDir("", ".docker")
	if err != nil {
		t.Fatal("Error creating tmp dir:", err)
	}
	hostListItems := make(chan hostListItem)
	store := NewStore(storePath, "", "")
	hosts := []Host{
		{
			Name:       "foo",
			DriverName: "fakedriver",
			Driver: &FakeDriver{
				MockState: state.Running,
			},
			storePath: storePath,
		},
		{
			Name:       "bar",
			DriverName: "fakedriver",
			Driver: &FakeDriver{
				MockState: state.Stopped,
			},
			storePath: storePath,
		},
		{
			Name:       "baz",
			DriverName: "fakedriver",
			Driver: &FakeDriver{
				MockState: state.Running,
			},
			storePath: storePath,
		},
	}
	expected := map[string]state.State{
		"foo": state.Running,
		"bar": state.Stopped,
		"baz": state.Running,
	}
	items := []hostListItem{}
	for _, host := range hosts {
		go getHostState(host, *store, hostListItems)
	}
	for i := 0; i < len(hosts); i++ {
		items = append(items, <-hostListItems)
	}
	for _, item := range items {
		if expected[item.Name] != item.State {
			t.Fatal("Expected state did not match for item", item)
		}
	}
}
