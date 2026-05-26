package contentanalysis

import (
	"fmt"
	L "github.com/fbaube/mlog"
	SU "github.com/fbaube/stringutils"
	XU "github.com/fbaube/xmlutils"
	S "strings"
)

func (pCA *ContentAnalysis) DoAnalysis_xml(pXP *XU.XmlPeek, sCont string) error {
	var filext string
	filext = pCA.FileExt
	// ===============================
	//  Set bool variables, including
	//  supporting analysis by stdlib
	// ===============================
	gotRootElm, tagsMsg := (pXP.ContentityBasics.CheckTopTags())
	gotDoctype := (pXP.DoctypeRaw != "")
	gotPreambl := (pXP.PreambleRaw != "")
	L.L.Debug("DoAnalysis_xml: raw: Dctp<%s> Prmbl<%s>",
		pXP.DoctypeRaw, pXP.PreambleRaw)
	/* L.L.Warning("DoAnalysis_xml: raw: Dctp<%s> Prmbl<%s>",
		pXP.DoctypeRaw, pXP.PreambleRaw) */
	// gotSomeXml := (gotRootElm || gotDoctype || gotPreambl)
	// Write a progress string
	/* if true */
	{
		var sP, sD, sR, sDtd string
		if gotPreambl {
			sP = "<?xml..> "
		}
		if gotDoctype {
			sD = "<!DOCTYPE..> "
		}
		if gotRootElm {
			sR = "root<" + pXP.XmlRoot.TagName + "> "
		}
		if pXP.HasDTDstuff {
			sDtd = "<DTD stuff> "
		}
		L.L.Info("Is XML: found: %s%s%s%s", sP, sD, sR, sDtd)
		//fmt.Printf("Is XML: found: %s%s%s%s \n", sP, sD, sR, sDtd)
		if tagsMsg != "" {
			L.L.Info("DoAnalysisXml: " + tagsMsg)
		}
	}
	if !(gotRootElm || pXP.HasDTDstuff) {
		L.L.Warning("(An.X) XML file has no root tag (and is not DTD)")
	}
	var pParstPrmbl *XU.ParsedPreamble
	var e error
	if gotPreambl {
		L.L.Debug("(An.X) got: %s", pXP.PreambleRaw)
		pParstPrmbl, e = XU.ParsePreamble(pXP.PreambleRaw)
		if e != nil {
			// println("xm.peek: preamble failure in:", peek.RawPreamble)
			return fmt.Errorf("(An.X) preamble failure: %w", e)
		}
		// print("--> Parsed XML preamble: " + pParstPrmbl.String())
	}
	// ================================
	//  Time to do some heavy lifting.
	// ================================
	L.L.Info("(DoAnalysisXml) Now split the file")
	if sCont == "" { // pCA.PathProps.Raw == "" {
		L.L.Error("(CA) XML has nil Raw")
	}
	pCA.ContentityBasics = pXP.ContentityBasics
	// L.L.Warning("SKIPPING call to pCA.MakeXmlContentitySections")
	// pCA.MakeXmlContentitySections()
	/* more debugging
	fmt.Printf("--> meta pos<%d>len<%d> text pos<%d>len<%d> \n",
		pAnlRec.Meta.Beg.Pos, len(pAnlRec.MetaRaw()),
		pAnlRec.Text.Beg.Pos, len(pAnlRec.TextRaw()))
	if !peek.IsSplittable() {
		println("--> Can't split into meta and text")
	}
	*/
	// =================================
	//  If we have DOCTYPE,
	//  it is gospel (and we are done).
	// =================================
	if gotDoctype {
		// We are here if we got a DOCTYPE; we also have a file extension,
		// and we should have a root tag (or else the DOCTYPE makes no sense !)
		var pParstDoctp *XU.ParsedDoctype
		pParstDoctp, e = pCA.ContypingInfo.ParseDoctype(pXP.DoctypeRaw)
		if e != nil {
			L.L.Panic("FIXME:" + e.Error())
		}
		pCA.ParsedDoctype = pParstDoctp
		L.L.Debug("(CA) gotDT: MType: " + pCA.MType)
		L.L.Debug("(CA) gotDT: AnalysisRecord: " + pCA.String())
		// L.L.Debug("gotDT: DctpFlds: " + pParstDoctp.String())
		/* DBG
		L.L.Warning("====")
		L.L.Warning("Raw: %s", pXP.DoctypeRaw)
		L.L.Warning("MTp: " + pCA.MType)
		L.L.Warning("ARc: " + pCA.String())
		L.L.Warning("====")
		*/
		if pCA.MType == "" {
			L.L.Panic("(CA) no MType, L103")
		}
		if pCA.MType == "" {
			L.L.Panic("(CA) no MType, L106")
		}
		L.L.Okay("(CA) Success: got XML with DOCTYPE")
		// HACK ALERT
		if /* IS_MAP || */ S.HasSuffix(pCA.MType, "---") {
			rutag := S.ToLower(pXP.XmlRoot.TagName)
			if pCA.MType == "xml/map/---" {
				pCA.MType = "xml/map/" + rutag
				L.L.Okay("(CA) Patched MType to: " + pCA.MType)
			} else {
				panic("MType ending in \"---\" not fixed")
			}
		}
		return nil
	}
	// =====================
	//  No DOCTYPE. Bummer.
	// =====================
	if !gotRootElm {
		return fmt.Errorf("(CA) Got no XML root tag in file with ext <%s>", filext)
	}
	// ==========================================
	//  Let's at least try to set the MType.
	//  We have a root tag and a file extension.
	// ==========================================
	rutag := S.ToLower(pXP.XmlRoot.TagName)
	IS_MAP := ("map" == pXP.XmlRoot.TagName)
	L.L.Info("DoAnXMl: XML without DOCTYPE: <%s> root<%s> MType<%s> isMap<%>",
		filext, rutag, pCA.MType, IS_MAP)
	if pCA.MType == "" {
	   L.L.Info("XML without DOCTYPE has no MType assigned yet")
	   }
	// Do some easy cases
	if rutag == "html" && (filext == ".html" || filext == ".htm") {
		pCA.MType = "html/cnt/html5"
	} else if rutag == "html" && S.HasPrefix(filext, ".xht") {
		pCA.MType = "html/cnt/xhtml"
	} else if SU.IsInSliceIgnoreCase(rutag, XU.DITArootElms) &&
		SU.IsInSliceIgnoreCase(filext, XU.DITAtypeFileExtensions) {
		pCA.MType = "xml/cnt/" + rutag
		if S.HasSuffix(rutag, "map") && S.HasSuffix(filext, "map") {
			pCA.MType = "xml/map/" + rutag
		}
	}
	// pAnlRec.ContypingInfo = *pCntpg
	if pCA.MType == "-/-/-" {
		pCA.MType = "xml/???/" + rutag
	}
	// At this point, mt should be valid !
	L.L.Debug("(CA) Contyping: " + pCA.ContypingInfo.String())

	// Now we should fill in all the detail fields.
	pCA.XmlContype = "RootTagData"

	if pParstPrmbl != nil {
		pCA.ParsedPreamble = pParstPrmbl
	} else {
		// SKIP
		// pBA.XmlPreambleFields = XU.STD_PreambleFields
	}
	// pBA.DoctypeIsDeclared  =  true
	pCA.DitaFlavor = "TBS"
	pCA.DitaContype = "TBS"

	// L.L.Info("fu.af: MType<%s> xcntp<%s> ditaFlav<%s> ditaCntp<%s> DT<%s>",
	L.L.Debug("(CA) final: MType<%s> xcntp<%s> dita:TBS DcTpFlds<%s>",
		pCA.MType, pCA.XmlContype, // pAnlRec.XmlPreambleFields,
		// pAnlRec.DitaFlavor, pAnlRec.DitaContype,
		pCA.ParsedDoctype)
	// println("--> fu.af: MetaRaw:", pAnlRec.MetaRaw())
	// println("--> fu.af: TextRaw:", pAnlRec.TextRaw())

	// HACK!
	if pCA.MType == "" {
		switch pCA.ContypingInfo.MimeTypeAsSnift { // m_contype {
		case "image/svg+xml":
			pCA.MType = "xml/img/svg"
		case "image/xml":
			if pCA.FileExt == ".svg" {	
			   pCA.MType = "xml/img/svg"
			   }
		}
		if pCA.MType != "" {
			L.L.Info("Hacked MType from image/[svg+]xml to %s", pCA.MType)
		}
	}
	L.L.Okay("(CA) Success: got XML without DOCTYPE")
	return nil
}
