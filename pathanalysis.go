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
	s := p.RawType()
	return s == SU.Raw_type_XML || s == SU.Raw_type_HTML
}

// MarkupType returns an enum with values of SU.Raw_type_*
func (p PathAnalysis) RawType() SU.Raw_type {
	// ======
     	//  HTML
	// ======
	// xml/HTML is an exceptional case
	if S.HasPrefix(p.MType, "xml/html/") {
		return SU.Raw_type_HTML
	}
	if S.HasPrefix(p.MimeType, "text/html") {
		return SU.Raw_type_HTML
	}
	if S.HasPrefix(p.MimeType, "html/") {
		return SU.Raw_type_HTML
	}
	if S.HasPrefix(p.MType, "html/") {
		return SU.Raw_type_HTML
	}
	// ======
     	//  XML
	// ======
	if S.HasPrefix(p.MType, "xml/") {
		return SU.Raw_type_XML
	}
	// ======
     	//  TEXT
	//  (MD)
	// ======
	if S.HasPrefix(p.MType, "text/") ||
		S.HasPrefix(p.MType, "txt/") ||
		S.HasPrefix(p.MType, "md/") ||
		S.HasPrefix(p.MType, "mkdn/") {
		return SU.Raw_type_MKDN
	}
	// ======
     	//  Misc.
	// ======
	if S.HasPrefix(p.MType, "bin/") {
		return SU.Raw_type_BIN // opaque
	}
	if S.HasPrefix(p.MType, "dir") || S.HasPrefix(p.MType, "DIR") {
		return SU.Raw_type_DIRLIKE 
	}
	// This is an unfortunate hack 
	if p.MType == "" { 
		return SU.Raw_type_DIRLIKE 
	}
	L.L.Error("CntAns.PathAns.RawType: failed on: <%s|%s>",
		p.MType, p.MimeType)
	return "" // SU.Raw_type_UNK (or OTHER?) 
}
