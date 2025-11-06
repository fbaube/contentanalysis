package contentanalysis

import (
	// "fmt"
	SU "github.com/fbaube/stringutils"
	XU "github.com/fbaube/xmlutils"
	L "github.com/fbaube/mlog"
	S "strings"
)

type Doctype  string
type MimeType string

// ContentAnalysis is the results of content analysis
// on the contents of a non-embedded [FSItem].
// .
type ContentAnalysis struct { // this has has Raw
	// ContypingInfo is simple fields:
	// FileExt MType MimeType's
	XU.ContypingInfo
	// ContentityBasics does NOT include Raw
	// (the entire input content)
	XU.ContentityBasics 
	// KeyElms is: (Root,Meta,Text)ElmExtent
	// KeyElmsWithRanges
	// ContentitySections is: Text_raw, Meta_raw, MetaFormat; MetaProps SU.PropSet
	// ContentityRawSections
	// XmlInfo is: XmlPreambleFields, XmlDoctype, XmlDoctypeFields, ENTITY stuff
	// ** XmlInfo **
	// XmlContype is an enum: "Unknown", "DTD", "DTDmod", "DTDent",
	// "RootTagData", "RootTagMixedContent", "MultipleRootTags", "INVALID"}
	XmlContype string
	// XmlPreambleFields is nil if no preamble - it can always
	// default to xmlutils.STD_PreambleFields (from stdlib)
	*XU.ParsedPreamble
	// XmlDoctypeFields is a ptr - nil if ContypingInfo.Doctype
	// is "", i.e. if there is no DOCTYPE declaration
	*XU.ParsedDoctype
	// DitaInfo
	DitaFlavor  string
	DitaContype string
}

// IsXML is true for all XML, including all HTML.
func (pCA ContentAnalysis) IsXML() bool {
	s := pCA.RawType()
	return s == SU.Raw_type_XML || s == SU.Raw_type_HTML
}

// MarkupType returns an enum with values of SU.Raw_type_*
func (pCA ContentAnalysis) RawType() SU.Raw_type {
	// ======
     	//  HTML
	// ======
	// xml/HTML is an exceptional case
	if S.HasPrefix(pCA.MType, "xml/html/") {
		return SU.Raw_type_HTML
	}
	if S.HasPrefix(pCA.MimeType, "text/html") {
		return SU.Raw_type_HTML
	}
	if S.HasPrefix(pCA.MimeType, "html/") {
		return SU.Raw_type_HTML
	}
	if S.HasPrefix(pCA.MType, "html/") {
		return SU.Raw_type_HTML
	}
	// ======
     	//  XML
	// ======
	if S.HasPrefix(pCA.MType, "xml/") {
		return SU.Raw_type_XML
	}
	// ======
     	//  TEXT
	//  (MD)
	// ======
	if S.HasPrefix(pCA.MType, "text/") ||
		S.HasPrefix(pCA.MType, "txt/") ||
		S.HasPrefix(pCA.MType, "md/") ||
		S.HasPrefix(pCA.MType, "mkdn/") {
		return SU.Raw_type_MKDN
	}
	// ======
     	//  Misc.
	// ======
	if S.HasPrefix(pCA.MType, "bin/") {
		return SU.Raw_type_BIN // opaque
	}
	if S.HasPrefix(pCA.MType, "dir") || S.HasPrefix(pCA.MType, "DIR") {
		return SU.Raw_type_DIRLIKE 
	}
	// This is an unfortunate hack 
	if pCA.MType == "" { 
		return SU.Raw_type_DIRLIKE 
	}
	L.L.Error("CntAns.PathAns.RawType: failed on: <%s|%s>",
		pCA.MType, pCA.MimeType)
	return "" // SU.Raw_type_UNK (or OTHER?) 
}
