package allocate

func Creds(id string, uidStart, gidStart int) (uid, gid int, err error) {
	n, err := parseID(id)
	if err != nil {
		return 0, 0, err
	}
	return uidStart + n, gidStart + n, nil
}
