package contentanalysis

import (
	// "fmt"
	SU "github.com/fbaube/stringutils"
	XU "github.com/fbaube/xmlutils"
	L "github.com/fbaube/mlog"
	S "strings"
)

type Doctype string
type MimeType string

// PathAnalysis is the results of content analysis
// on the contents of a non-embedded [FSItem].
// .
type PathAnalysis struct { // this has has Raw
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
func (p PathAnalysis) IsXML() bool {
	s := p.MarkupTypeOfMType()
	return s == SU.MU_type_XML || s == SU.MU_type_HTML
}

// MarkupType returns an enum with values of SU.MU_type_*
func (p PathAnalysis) MarkupTypeOfMType() SU.MarkupType {
	// HTML is an exceptional case
	if S.HasPrefix(p.MType, "xml/html/") {
		return SU.MU_type_HTML
	}
	if S.HasPrefix(p.MimeType, "text/html") {
		return SU.MU_type_HTML
	}
	if S.HasPrefix(p.MimeType, "html/") {
		return SU.MU_type_HTML
	}
	if S.HasPrefix(p.MType, "html/") {
		return SU.MU_type_HTML
	}
	if S.HasPrefix(p.MType, "xml/") {
		return SU.MU_type_XML
	}
	if S.HasPrefix(p.MType, "text/") ||
		S.HasPrefix(p.MType, "txt/") ||
		S.HasPrefix(p.MType, "md/") ||
		S.HasPrefix(p.MType, "mkdn/") {
		return SU.MU_type_MKDN
	}
	if S.HasPrefix(p.MType, "bin/") {
		return SU.MU_type_BIN // opaque
	}
	if S.HasPrefix(p.MType, "dir") || S.HasPrefix(p.MType, "DIR") {
		return SU.MU_type_DIRLIKE 
	}
	if p.MType == "" { 
		return SU.MU_type_DIRLIKE 
	}
	// panic("UNK!")
	L.L.Error("ca.pa.MarkupTypeOfMType: failed (isDir?) on: %s", p.MType)
	return SU.MU_type_UNK
}
