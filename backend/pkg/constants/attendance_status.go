package constants

type AttendanceStatus string

const (
	AttendanceStatusPresent AttendanceStatus = "PRESENT"
	AttendanceStatusLate    AttendanceStatus = "LATE"
	AttendanceStatusExcused AttendanceStatus = "EXCUSED"
	AttendanceStatusAbsent  AttendanceStatus = "ABSENT"
	AttendanceStatusSick    AttendanceStatus = "SICK"
)
