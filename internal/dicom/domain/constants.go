package domain

const (
	PDUAssociateRQ byte = 0x01
	PDUAssociateAC byte = 0x02
	PDUAssociateRJ byte = 0x03
	PDUDATATF      byte = 0x04
	PDUReleaseRQ   byte = 0x05
	PDUReleaseRP   byte = 0x06
	PDUAbort       byte = 0x07
)

var SupportedTransferSyntaxes = []string{"1.2.840.10008.1.2", "1.2.840.10008.1.2.1", "1.2.840.10008.1.2.2"}
var SupportedSOPClasses = []string{"1.2.840.10008.5.1.4.1.1.2", "1.2.840.10008.5.1.4.1.1.4", "1.2.840.10008.5.1.4.1.1.128"}
