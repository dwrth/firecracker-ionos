package allocate

import "fmt"

type Names struct {
	Namespace string
	HostVeth  string
	NSVeth    string
	TAP       string
}

func InterfaceNames(id string) (Names, error) {
	n, err := parseID(id)
	if err != nil {
		return Names{}, err
	}
	return Names{
		Namespace: fmt.Sprintf("kn-%s", id),
		HostVeth:  fmt.Sprintf("kn%d-host", n),
		NSVeth:    fmt.Sprintf("kn%d-ns", n),
		TAP:       "tap0",
	}, nil
}
