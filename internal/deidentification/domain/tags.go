package domain

const (
	PatientName        = "PatientName"
	PatientID          = "PatientID"
	PatientBirthDate   = "PatientBirthDate"
	PatientSex         = "PatientSex"
	AccessionNumber    = "AccessionNumber"
	StudyDate          = "StudyDate"
	StudyTime          = "StudyTime"
	StudyInstanceUID   = "StudyInstanceUID"
	SeriesInstanceUID  = "SeriesInstanceUID"
	SOPInstanceUID     = "SOPInstanceUID"
	InstitutionName    = "InstitutionName"
	ReferringPhysician = "ReferringPhysician"
)

var DefaultRemovedTags = []string{PatientName, PatientBirthDate, PatientSex, AccessionNumber, InstitutionName, ReferringPhysician}
