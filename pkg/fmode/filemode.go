package fmode

const (
	read       = 04
	write      = 02
	ex         = 01
	userShift  = 6
	groupShift = 3
	othShift   = 0
)

//nolint:stylecheck
const (
	USER_R   = read << userShift
	USER_W   = write << userShift
	USER_X   = ex << userShift
	USER_RW  = USER_R | USER_W
	USER_RWX = USER_RW | USER_X

	GROUP_R   = read << groupShift
	GROUP_W   = write << groupShift
	GROUP_X   = ex << groupShift
	GROUP_RW  = GROUP_R | GROUP_W
	GROUP_RWX = GROUP_RW | GROUP_X

	OTHER_R   = read << othShift
	OTHER_W   = write << othShift
	OTHER_X   = ex << othShift
	OTHER_RW  = OTHER_R | OTHER_W
	OTHER_RWX = OTHER_RW | OTHER_X

	ALL_R   = USER_R | GROUP_R | OTHER_R
	ALL_W   = USER_W | GROUP_W | OTHER_W
	ALL_X   = USER_X | GROUP_X | OTHER_X
	ALL_RW  = ALL_R | ALL_W
	ALL_RWX = ALL_RW | GROUP_X
)
