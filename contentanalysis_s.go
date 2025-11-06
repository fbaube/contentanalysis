package contentanalysis

import (
	S "strings"
)

func (pCA *ContentAnalysis) String() string {
	var sb S.Builder
	var sPDT string
	if pCA.ParsedDoctype != nil {
		sPDT = pCA.ParsedDoctype.String()
	}
	sb.WriteString("ContentAnalysis: ")
	sb.WriteString("CntpgInfo: \n\t" + pCA.ContypingInfo.String() + "\n\t")
	sb.WriteString("XmlCntp<" + pCA.XmlContype + "> ")
	sb.WriteString("XmlDctp<" + sPDT + "> ")
	return sb.String()
}
