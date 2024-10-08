package contentanalysis

import (
	"fmt"
	"io/fs"
	"net/http"
	FP "path/filepath"
	S "strings"

	"github.com/gabriel-vasile/mimetype"
	"golang.org/x/tools/godoc/util" // used once, L125

	L "github.com/fbaube/mlog"
	SU "github.com/fbaube/stringutils"
	XU "github.com/fbaube/xmlutils"
	FU "github.com/fbaube/fileutils"
)

// <!ELEMENT  map     (topicmeta?, (topicref | keydef)*)  >
// <!ELEMENT topicmeta (navtitle?, linktext?, data*) >

// NewPathAnalysis is called only by NewContentityRecord(..).
// It has very different handling for XML content versus non-XML content.
// Most of the function is making several checks for the presence of XML.
// When a file is identified as XML, we have much more info available,
// so processing becomes both simpler and more complicated.
//
// Binary content is tagged as such and no further examination is done.
// So, the basic top-level classificaton of content is:
//   - Binary
//   - XML (when a DOCTYPE is detected)
//   - Everything else (incl. plain text,
//     Markdown, and XML/HTML that lacks DOCTYPE)
//
// If the argument is "dirlike" (dir, symlink, etc.), 
// then NewPathAnalysis returns (nil, nil).
//
// If the first argument "sCont" (the content) is less than six bytes,
// return (nil, nil) to indicate that there is not enough content with
// which to do anything productive or informative. 
// .
func NewPathAnalysis(pFSI *FU.FSItem) (*PathAnalysis, error) {
     	// println("NewPathAnalysis: Entering!")
     	// If it's not a file, GTFO
	if pFSI.IsDirlike() {
	   // println("NewPathAnalysis: dirlike")
	   return nil, nil
	}
	if !pFSI.IsFile() {
	   // println("NewPathAnalysis: not a file")
	   return nil, nil
	}
	var sCont string
	if pFSI.TypedRaw != nil {
	   // println("NewPathAnalysis: dupe pFSI.LoadContents?")
	   }
	elc := pFSI.LoadContents()
	if elc != nil {
	   pFSI.SetError(fmt.Errorf("LoadContents: %w", elc))
	   return nil, &fs.PathError{ Op:"LoadContents",
	   	  Path:pFSI.FPs.CreationPath(), Err:elc }
	   }
	if pFSI.TypedRaw == nil {
	   // println("NewPathAnalysis: failed pFSI.LoadContents")
	   }
	if pFSI.TypedRaw != nil {
	   sCont = string(pFSI.TypedRaw.Raw)
	   } else {
	   // println("NewPathAnalysis: NO TypedRaw!")
	   }
	filext := FP.Ext(pFSI.FPs.AbsFP)

	// A trailing dot in the filename provides no filetype info
	filext = FP.Ext(filext)
	if filext == "." {
		filext = ""
	}
	// ===========================
	//  Handle pathological case:
	//  Too short or non-existent
	// ===========================
	if len(sCont) < 6 {
		if sCont == "" {
			L.L.Info("NewPathAnalysis: skipping zero-length content")
		} else {
			L.L.Warning("DNewPathnalysis: content too short (%d bytes)", len(sCont))
		}
		p := new(PathAnalysis)
		p.FileExt = filext
		// return nil, errors.New(fmt.Sprintf(
		//    "content is too short (%d bytes) to analyse", len(sCont)))
		return p, nil
	}
	L.L.Debug("NewPathAnalysis: filext<%s> len<%d> beg<%s>",
		filext, len(sCont), sCont[:5])

	// ========================
	//  Try a coupla shortcuts
	// ========================
	cheatYaml := S.HasPrefix(sCont, "---\n")
	// FIXME cheatHtml should also permit a preceding XML prolog
	cheatHtml := S.HasPrefix(sCont, "<DOCTYPE html")
	cheatXml := S.HasPrefix(sCont, "<?xml ") || S.HasPrefix(sCont, "<DOCTYPE ")
	cheatXml = cheatXml && !cheatHtml
	// =========================
	//  Content type detection
	//  using 3rd party library
	//  (this is AUTHORITATIVE!)
	// =========================
	var contypeData *mimetype.MIME
	var contype string // *mimetype.MIME
	// MIME type ?
	contypeData = mimetype.Detect([]byte(sCont))
	contype = contypeData.String()
	contype = S.TrimSuffix(contype, "; charset=utf-8")
	if S.Contains(contype, ";") {
		L.L.Warning("Content type from lib_3p " +
			"(still) has a semicolon: " + contype)
	}
	// ======================================
	//  Also do content type detection using
	//     HTTP stdlib (this is UNRELIABLE!)
	// ======================================
	var stdlib_contype string
	stdlib_contype = http.DetectContentType([]byte(sCont))
	stdlib_contype = S.TrimSuffix(stdlib_contype, "; charset=utf-8")
	if S.Contains(stdlib_contype, ";") {
		L.L.Warning("Content type from stdlib " +
			"(still) has a semicolon: " + stdlib_contype)
	}
	// ===========================
	//  Warn if they do not agree
	// ===========================
	if stdlib_contype != contype {
		L.L.Warning("NPA<%s>: MIME type: lib_3p<%s>OK?, stdlib<%s>Bad?",
			filext, contype, stdlib_contype)
	} else {
	  L.L.Info("NPA<%s>: MIME type: snift-as: %s", filext, contype)
	}
	// =====================================
	// INITIALIZE ANALYSIS RECORD:
	// pAnlRec is *xmlutils.PathAnalysis
	// is basically all of our analysis
	// results, including ContypingInfo.
	// Don't forget to set the content!
	// (omitting this caused a lot of bugs)
	// =====================================
	var pPA *PathAnalysis
	pPA = new(PathAnalysis)
	pPA.FileExt = filext
	pPA.MimeType = stdlib_contype // Junk this ?
	pPA.MimeTypeAsSnift = contype
	// ===========================
	//  Check for & handle BINARY
	// ===========================
	var isBinary bool        // Authoritative!
	var stdlib_isBinary bool // Unreliable!
	stdlib_isBinary = !util.IsText([]byte(sCont))
	isBinary = true
	for mime := contypeData; mime != nil; mime = mime.Parent() {
		if mime.Is("text/plain") {
			isBinary = false
		}
		// FIXME If "text/" here, is an error to sort out
	}
	// Warn if they disagree
	if stdlib_isBinary != isBinary {
		L.L.Warning("NPA<%s>: MIME type: stdlib err, " +
			"says is-binary: <%t>", filext, stdlib_isBinary)
	}
	if isBinary {
		if cheatYaml || cheatXml || cheatHtml {
			L.L.Error("NPA: both is-Binary & is-Yaml/Xml")
		}
		return pPA, pPA.DoAnalysis_bin()
	}
	// ======================================
	// We have text, but it might not be XML.
	// So process the MIME types returned by
	// the two libraries.
	// ======================================
	// hIsXml, hMsg := contypeIsXml("stdlib", stdlib_contype, filext)
	mIsXml, mMsg := contypeIsXml("lib_3p", contype, filext)
	if mIsXml { // || hIsXml {
		var mS string // hS
		// if !hIsXml {
		//	hS = "not "
		// }
		if !mIsXml {
			mS = "Not-"
		}
		L.L.Info("NPA: lib_3p (is-%sXML) %s", // \n\t\t stdlib (is-%sXML) %s",
			mS, mMsg) // , hS, hMsg)
	} else {
		L.L.Info("NPA: XML not detected by either MIME lib")
	}
	// ===================================
	//  MAIN XML PRELIMINARY ANALYSIS:
	//  Peek into file to look for root
	//  tag and other top-level XML stuff
	// ===================================
	var xmlParsingFailed bool
	var pPeek *XU.XmlPeek
	var e error
	// ==============================
	//  Peek for XML; this also sets
	//  KeyElms (Root,Meta,Text)
	// ==============================
	pPeek, e = XU.Peek_xml(sCont)

	// NOTE! An error from peeking might be caused
	// by, for example, applying XML parsing to a
	// Markdown file. So, an error is NOT fatal.
	if e != nil {
		L.L.Info("NPA: XML parsing got error: " + e.Error())
		xmlParsingFailed = true
	}
	// ===============================
	//  If it's DTD stuff, we're done
	// ===============================
	if pPeek.HasDTDstuff && SU.IsInSliceIgnoreCase(
		filext, XU.DTDtypeFileExtensions) {
		return pPA, pPA.DoAnalysis_sch()
	}
	// Check for pathological cases
	if xmlParsingFailed && mIsXml { // (hIsXml || mIsXml) {
		L.L.Panic("NPA: XML confusion (case #1) in DoAnalysis")
	}
	// Note that this next test dusnt always work for Markdown!
	// if (!xmlParsingFailed) && (! (hIsXml || mIsXml)) {
	//      L.L.Panic("XML confusion (case #2) in AnalyzeFile")
	// }
	var hasRootTag, gotSomeXml bool
	hasRootTag, _ = pPeek.ContentityBasics.CheckTopTags()
	gotSomeXml = hasRootTag || (pPeek.DoctypeRaw != "") ||
		(pPeek.PreambleRaw != "")
	// =============================
	//  If it's not XML, we're done
	// =============================
	if xmlParsingFailed || !gotSomeXml {
		if cheatXml {
			// L.L.Panic("(AF) both non-xml & xml")
			L.L.Panic(fmt.Sprintf("WHOOPS xmlParsingFailed<%t> gotSomeXml<%t> \n",
				xmlParsingFailed, gotSomeXml))
		}
		return pPA, pPA.DoAnalysis_txt(sCont)
	}
	// ===========================================
	//  It's XML, so crank thru it and we're done
	// ===========================================
	L.L.Debug("NPA passing to DoAnalysis_xml: Peek: %+v", *pPeek) 
	return pPA, pPA.DoAnalysis_xml(pPeek, sCont)
}
