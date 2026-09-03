package allocate

// Creds returns jailer uid and gid for id, offset from uidStart and gidStart.
func Creds(slot int, uidStart, gidStart int) (uid, gid int) {
	return uidStart + slot, gidStart + slot
}
