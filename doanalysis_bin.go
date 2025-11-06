package contentanalysis

import (
	L "github.com/fbaube/mlog"
	S "strings"
)

// DoAnalysis_bin doesn't do any further processing for binary, cos we
// basically trust that the sniffed MIME type is sufficient, and return.
// .
func (pCA *ContentAnalysis) DoAnalysis_bin() error {
	// pAnlRec.MimeType = m_contype
	pCA.MType = "bin/"
	m_contype := pCA.ContypingInfo.MimeTypeAsSnift
	if S.HasPrefix(m_contype, "image/") {
		det := S.TrimPrefix(m_contype, "image/")
		pCA.MType += "img/" // append, not replace 
		hasEPS := S.Contains(m_contype, "eps")
		hasTXT := S.Contains(m_contype, "text") ||
			S.Contains(m_contype, "txt")
		if hasTXT || hasEPS {
			// TODO
			L.L.Warning("(CA) EPS/TXT confusion for MIME type: " + m_contype)
			pCA.MType = "txt/img/?eps"
		} else {
			pCA.MType += det
			if !S.EqualFold(pCA.FileExt, "."+det) {
				L.L.Warning("Image: detMime<%s> filext<%s>", det, pCA.FileExt)
			}
		}
	} else {
		L.L.Warning("Image problem: mime<%s> filext<%s>",
			pCA.MimeTypeAsSnift, pCA.FileExt)
	}
	// L.L.Okay("(CA) Success: detected BINARY")
	return nil
}
