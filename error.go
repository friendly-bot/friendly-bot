package main

type constErr string

const (
	UnknownDatastoreErr constErr = "unknown datastore type"
)

func (e constErr) Error() string {
	return string(e)
}
