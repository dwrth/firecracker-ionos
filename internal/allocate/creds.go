package allocate

// Creds returns jailer uid and gid for id, offset from uidStart and gidStart.
func Creds(id string, uidStart, gidStart int) (uid, gid int, err error) {
	n, err := parseSlotKey(id)
	if err != nil {
		return 0, 0, err
	}
	return uidStart + n, gidStart + n, nil
}
