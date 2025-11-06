package contentanalysis

import (
	L "github.com/fbaube/mlog"
	S "strings"
)

// DoAnalysis_sch will handle DTDs and related files,
// and the code is mostly written but not yet integrated,
// so this func doesn't really worry about it yet.
// .
func (pCA *ContentAnalysis) DoAnalysis_sch() error {
	// L.L.Okay("(CA) Success: DTD-type content detected (filext<%s>)", filext)
	pCA.MimeType = "application/xml-dtd"
	pCA.MType = "xml/sch/" + S.ToLower(S.TrimPrefix(pCA.FileExt, "."))
	L.L.Warning("(CA) DTD stuff: should allocate and fill fields")
	return nil
}
