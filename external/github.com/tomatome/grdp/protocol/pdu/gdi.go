package pdu

/* Binary Raster Operations (ROP2) */
const (
	GDI_R2_BLACK       = 0x01
	GDI_R2_NOTMERGEPEN = 0x02
	GDI_R2_MASKNOTPEN  = 0x03
	GDI_R2_NOTCOPYPEN  = 0x04
	GDI_R2_MASKPENNOT  = 0x05
	GDI_R2_NOT         = 0x06
	GDI_R2_XORPEN      = 0x07
	GDI_R2_NOTMASKPEN  = 0x08
	GDI_R2_MASKPEN     = 0x09
	GDI_R2_NOTXORPEN   = 0x0A
	GDI_R2_NOP         = 0x0B
	GDI_R2_MERGENOTPEN = 0x0C
	GDI_R2_COPYPEN     = 0x0D
	GDI_R2_MERGEPENNOT = 0x0E
	GDI_R2_MERGEPEN    = 0x0F
	GDI_R2_WHITE       = 0x10
)

/* Ternary Raster Operations (ROP3) */
const (
	GDI_BLACKNESS   = 0x00000042
	GDI_DPSoon      = 0x00010289
	GDI_DPSona      = 0x00020C89
	GDI_PSon        = 0x000300AA
	GDI_SDPona      = 0x00040C88
	GDI_DPon        = 0x000500A9
	GDI_PDSxnon     = 0x00060865
	GDI_PDSaon      = 0x000702C5
	GDI_SDPnaa      = 0x00080F08
	GDI_PDSxon      = 0x00090245
	GDI_DPna        = 0x000A0329
	GDI_PSDnaon     = 0x000B0B2A
	GDI_SPna        = 0x000C0324
	GDI_PDSnaon     = 0x000D0B25
	GDI_PDSonon     = 0x000E08A5
	GDI_Pn          = 0x000F0001
	GDI_PDSona      = 0x00100C85
	GDI_NOTSRCERASE = 0x001100A6
	GDI_SDPxnon     = 0x00120868
	GDI_SDPaon      = 0x001302C8
	GDI_DPSxnon     = 0x00140869
	GDI_DPSaon      = 0x001502C9
	GDI_PSDPSanaxx  = 0x00165CCA
	GDI_SSPxDSxaxn  = 0x00171D54
	GDI_SPxPDxa     = 0x00180D59
	GDI_SDPSanaxn   = 0x00191CC8
	GDI_PDSPaox     = 0x001A06C5
	GDI_SDPSxaxn    = 0x001B0768
	GDI_PSDPaox     = 0x001C06CA
	GDI_DSPDxaxn    = 0x001D0766
	GDI_PDSox       = 0x001E01A5
	GDI_PDSoan      = 0x001F0385
	GDI_DPSnaa      = 0x00200F09
	GDI_SDPxon      = 0x00210248
	GDI_DSna        = 0x00220326
	GDI_SPDnaon     = 0x00230B24
	GDI_SPxDSxa     = 0x00240D55
	GDI_PDSPanaxn   = 0x00251CC5
	GDI_SDPSaox     = 0x002606C8
	GDI_SDPSxnox    = 0x00271868
	GDI_DPSxa       = 0x00280369
	GDI_PSDPSaoxxn  = 0x002916CA
	GDI_DPSana      = 0x002A0CC9
	GDI_SSPxPDxaxn  = 0x002B1D58
	GDI_SPDSoax     = 0x002C0784
	GDI_PSDnox      = 0x002D060A
	GDI_PSDPxox     = 0x002E064A
	GDI_PSDnoan     = 0x002F0E2A
	GDI_PSna        = 0x0030032A
	GDI_SDPnaon     = 0x00310B28
	GDI_SDPSoox     = 0x00320688
	GDI_NOTSRCCOPY  = 0x00330008
	GDI_SPDSaox     = 0x003406C4
	GDI_SPDSxnox    = 0x00351864
	GDI_SDPox       = 0x003601A8
	GDI_SDPoan      = 0x00370388
	GDI_PSDPoax     = 0x0038078A
	GDI_SPDnox      = 0x00390604
	GDI_SPDSxox     = 0x003A0644
	GDI_SPDnoan     = 0x003B0E24
	GDI_PSx         = 0x003C004A
	GDI_SPDSonox    = 0x003D18A4
	GDI_SPDSnaox    = 0x003E1B24
	GDI_PSan        = 0x003F00EA
	GDI_PSDnaa      = 0x00400F0A
	GDI_DPSxon      = 0x00410249
	GDI_SDxPDxa     = 0x00420D5D
	GDI_SPDSanaxn   = 0x00431CC4
	GDI_SRCERASE    = 0x00440328
	GDI_DPSnaon     = 0x00450B29
	GDI_DSPDaox     = 0x004606C6
	GDI_PSDPxaxn    = 0x0047076A
	GDI_SDPxa       = 0x00480368
	GDI_PDSPDaoxxn  = 0x004916C5
	GDI_DPSDoax     = 0x004A0789
	GDI_PDSnox      = 0x004B0605
	GDI_SDPana      = 0x004C0CC8
	GDI_SSPxDSxoxn  = 0x004D1954
	GDI_PDSPxox     = 0x004E0645
	GDI_PDSnoan     = 0x004F0E25
	GDI_PDna        = 0x00500325
	GDI_DSPnaon     = 0x00510B26
	GDI_DPSDaox     = 0x005206C9
	GDI_SPDSxaxn    = 0x00530764
	GDI_DPSonon     = 0x005408A9
	GDI_DSTINVERT   = 0x00550009
	GDI_DPSox       = 0x005601A9
	GDI_DPSoan      = 0x00570389
	GDI_PDSPoax     = 0x00580785
	GDI_DPSnox      = 0x00590609
	GDI_PATINVERT   = 0x005A0049
	GDI_DPSDonox    = 0x005B18A9
	GDI_DPSDxox     = 0x005C0649
	GDI_DPSnoan     = 0x005D0E29
	GDI_DPSDnaox    = 0x005E1B29
	GDI_DPan        = 0x005F00E9
	GDI_PDSxa       = 0x00600365
	GDI_DSPDSaoxxn  = 0x006116C6
	GDI_DSPDoax     = 0x00620786
	GDI_SDPnox      = 0x00630608
	GDI_SDPSoax     = 0x00640788
	GDI_DSPnox      = 0x00650606
	GDI_SRCINVERT   = 0x00660046
	GDI_SDPSonox    = 0x006718A8
	GDI_DSPDSonoxxn = 0x006858A6
	GDI_PDSxxn      = 0x00690145
	GDI_DPSax       = 0x006A01E9
	GDI_PSDPSoaxxn  = 0x006B178A
	GDI_SDPax       = 0x006C01E8
	GDI_PDSPDoaxxn  = 0x006D1785
	GDI_SDPSnoax    = 0x006E1E28
	GDI_PDSxnan     = 0x006F0C65
	GDI_PDSana      = 0x00700CC5
	GDI_SSDxPDxaxn  = 0x00711D5C
	GDI_SDPSxox     = 0x00720648
	GDI_SDPnoan     = 0x00730E28
	GDI_DSPDxox     = 0x00740646
	GDI_DSPnoan     = 0x00750E26
	GDI_SDPSnaox    = 0x00761B28
	GDI_DSan        = 0x007700E6
	GDI_PDSax       = 0x007801E5
	GDI_DSPDSoaxxn  = 0x00791786
	GDI_DPSDnoax    = 0x007A1E29
	GDI_SDPxnan     = 0x007B0C68
	GDI_SPDSnoax    = 0x007C1E24
	GDI_DPSxnan     = 0x007D0C69
	GDI_SPxDSxo     = 0x007E0955
	GDI_DPSaan      = 0x007F03C9
	GDI_DPSaa       = 0x008003E9
	GDI_SPxDSxon    = 0x00810975
	GDI_DPSxna      = 0x00820C49
	GDI_SPDSnoaxn   = 0x00831E04
	GDI_SDPxna      = 0x00840C48
	GDI_PDSPnoaxn   = 0x00851E05
	GDI_DSPDSoaxx   = 0x008617A6
	GDI_PDSaxn      = 0x008701C5
	GDI_SRCAND      = 0x008800C6
	GDI_SDPSnaoxn   = 0x00891B08
	GDI_DSPnoa      = 0x008A0E06
	GDI_DSPDxoxn    = 0x008B0666
	GDI_SDPnoa      = 0x008C0E08
	GDI_SDPSxoxn    = 0x008D0668
	GDI_SSDxPDxax   = 0x008E1D7C
	GDI_PDSanan     = 0x008F0CE5
	GDI_PDSxna      = 0x00900C45
	GDI_SDPSnoaxn   = 0x00911E08
	GDI_DPSDPoaxx   = 0x009217A9
	GDI_SPDaxn      = 0x009301C4
	GDI_PSDPSoaxx   = 0x009417AA
	GDI_DPSaxn      = 0x009501C9
	GDI_DPSxx       = 0x00960169
	GDI_PSDPSonoxx  = 0x0097588A
	GDI_SDPSonoxn   = 0x00981888
	GDI_DSxn        = 0x00990066
	GDI_DPSnax      = 0x009A0709
	GDI_SDPSoaxn    = 0x009B07A8
	GDI_SPDnax      = 0x009C0704
	GDI_DSPDoaxn    = 0x009D07A6
	GDI_DSPDSaoxx   = 0x009E16E6
	GDI_PDSxan      = 0x009F0345
	GDI_DPa         = 0x00A000C9
	GDI_PDSPnaoxn   = 0x00A11B05
	GDI_DPSnoa      = 0x00A20E09
	GDI_DPSDxoxn    = 0x00A30669
	GDI_PDSPonoxn   = 0x00A41885
	GDI_PDxn        = 0x00A50065
	GDI_DSPnax      = 0x00A60706
	GDI_PDSPoaxn    = 0x00A707A5
	GDI_DPSoa       = 0x00A803A9
	GDI_DPSoxn      = 0x00A90189
	GDI_DSTCOPY     = 0x00AA0029
	GDI_DPSono      = 0x00AB0889
	GDI_SPDSxax     = 0x00AC0744
	GDI_DPSDaoxn    = 0x00AD06E9
	GDI_DSPnao      = 0x00AE0B06
	GDI_DPno        = 0x00AF0229
	GDI_PDSnoa      = 0x00B00E05
	GDI_PDSPxoxn    = 0x00B10665
	GDI_SSPxDSxox   = 0x00B21974
	GDI_SDPanan     = 0x00B30CE8
	GDI_PSDnax      = 0x00B4070A
	GDI_DPSDoaxn    = 0x00B507A9
	GDI_DPSDPaoxx   = 0x00B616E9
	GDI_SDPxan      = 0x00B70348
	GDI_PSDPxax     = 0x00B8074A
	GDI_DSPDaoxn    = 0x00B906E6
	GDI_DPSnao      = 0x00BA0B09
	GDI_MERGEPAINT  = 0x00BB0226
	GDI_SPDSanax    = 0x00BC1CE4
	GDI_SDxPDxan    = 0x00BD0D7D
	GDI_DPSxo       = 0x00BE0269
	GDI_DPSano      = 0x00BF08C9
	GDI_MERGECOPY   = 0x00C000CA
	GDI_SPDSnaoxn   = 0x00C11B04
	GDI_SPDSonoxn   = 0x00C21884
	GDI_PSxn        = 0x00C3006A
	GDI_SPDnoa      = 0x00C40E04
	GDI_SPDSxoxn    = 0x00C50664
	GDI_SDPnax      = 0x00C60708
	GDI_PSDPoaxn    = 0x00C707AA
	GDI_SDPoa       = 0x00C803A8
	GDI_SPDoxn      = 0x00C90184
	GDI_DPSDxax     = 0x00CA0749
	GDI_SPDSaoxn    = 0x00CB06E4
	GDI_SRCCOPY     = 0x00CC0020
	GDI_SDPono      = 0x00CD0888
	GDI_SDPnao      = 0x00CE0B08
	GDI_SPno        = 0x00CF0224
	GDI_PSDnoa      = 0x00D00E0A
	GDI_PSDPxoxn    = 0x00D1066A
	GDI_PDSnax      = 0x00D20705
	GDI_SPDSoaxn    = 0x00D307A4
	GDI_SSPxPDxax   = 0x00D41D78
	GDI_DPSanan     = 0x00D50CE9
	GDI_PSDPSaoxx   = 0x00D616EA
	GDI_DPSxan      = 0x00D70349
	GDI_PDSPxax     = 0x00D80745
	GDI_SDPSaoxn    = 0x00D906E8
	GDI_DPSDanax    = 0x00DA1CE9
	GDI_SPxDSxan    = 0x00DB0D75
	GDI_SPDnao      = 0x00DC0B04
	GDI_SDno        = 0x00DD0228
	GDI_SDPxo       = 0x00DE0268
	GDI_SDPano      = 0x00DF08C8
	GDI_PDSoa       = 0x00E003A5
	GDI_PDSoxn      = 0x00E10185
	GDI_DSPDxax     = 0x00E20746
	GDI_PSDPaoxn    = 0x00E306EA
	GDI_SDPSxax     = 0x00E40748
	GDI_PDSPaoxn    = 0x00E506E5
	GDI_SDPSanax    = 0x00E61CE8
	GDI_SPxPDxan    = 0x00E70D79
	GDI_SSPxDSxax   = 0x00E81D74
	GDI_DSPDSanaxxn = 0x00E95CE6
	GDI_DPSao       = 0x00EA02E9
	GDI_DPSxno      = 0x00EB0849
	GDI_SDPao       = 0x00EC02E8
	GDI_SDPxno      = 0x00ED0848
	GDI_SRCPAINT    = 0x00EE0086
	GDI_SDPnoo      = 0x00EF0A08
	GDI_PATCOPY     = 0x00F00021
	GDI_PDSono      = 0x00F10885
	GDI_PDSnao      = 0x00F20B05
	GDI_PSno        = 0x00F3022A
	GDI_PSDnao      = 0x00F40B0A
	GDI_PDno        = 0x00F50225
	GDI_PDSxo       = 0x00F60265
	GDI_PDSano      = 0x00F708C5
	GDI_PDSao       = 0x00F802E5
	GDI_PDSxno      = 0x00F90845
	GDI_DPo         = 0x00FA0089
	GDI_PATPAINT    = 0x00FB0A09
	GDI_PSo         = 0x00FC008A
	GDI_PSDnoo      = 0x00FD0A0A
	GDI_DPSoo       = 0x00FE02A9
	GDI_WHITENESS   = 0x00FF0062
	GDI_GLYPH_ORDER = 0xFFFFFFFF
)

var Rop3CodeTable = map[int]string{
	GDI_BLACKNESS:   "0",
	GDI_DPSoon:      "DPSoon",
	GDI_DPSona:      "DPSona",
	GDI_PSon:        "PSon",
	GDI_SDPona:      "SDPona",
	GDI_DPon:        "DPon",
	GDI_PDSxnon:     "PDSxnon",
	GDI_PDSaon:      "PDSaon",
	GDI_SDPnaa:      "SDPnaa",
	GDI_PDSxon:      "PDSxon",
	GDI_DPna:        "DPna",
	GDI_PSDnaon:     "PSDnaon",
	GDI_SPna:        "SPna",
	GDI_PDSnaon:     "PDSnaon",
	GDI_PDSonon:     "PDSonon",
	GDI_Pn:          "Pn",
	GDI_PDSona:      "PDSona",
	GDI_NOTSRCERASE: "DSon",
	GDI_SDPxnon:     "SDPxnon",
	GDI_SDPaon:      "SDPaon",
	GDI_DPSxnon:     "DPSxnon",
	GDI_DPSaon:      "DPSaon",
	GDI_PSDPSanaxx:  "PSDPSanaxx",
	GDI_SSPxDSxaxn:  "SSPxDSxaxn",
	GDI_SPxPDxa:     "SPxPDxa",
	GDI_SDPSanaxn:   "SDPSanaxn",
	GDI_PDSPaox:     "PDSPaox",
	GDI_SDPSxaxn:    "SDPSxaxn",
	GDI_PSDPaox:     "PSDPaox",
	GDI_DSPDxaxn:    "DSPDxaxn",
	GDI_PDSox:       "PDSox",
	GDI_PDSoan:      "PDSoan",
	GDI_DPSnaa:      "DPSnaa",
	GDI_SDPxon:      "SDPxon",
	GDI_DSna:        "DSna",
	GDI_SPDnaon:     "SPDnaon",
	GDI_SPxDSxa:     "SPxDSxa",
	GDI_PDSPanaxn:   "PDSPanaxn",
	GDI_SDPSaox:     "SDPSaox",
	GDI_SDPSxnox:    "SDPSxnox",
	GDI_DPSxa:       "DPSxa",
	GDI_PSDPSaoxxn:  "PSDPSaoxxn",
	GDI_DPSana:      "DPSana",
	GDI_SSPxPDxaxn:  "SSPxPDxaxn",
	GDI_SPDSoax:     "SPDSoax",
	GDI_PSDnox:      "PSDnox",
	GDI_PSDPxox:     "PSDPxox",
	GDI_PSDnoan:     "PSDnoan",
	GDI_PSna:        "PSna",
	GDI_SDPnaon:     "SDPnaon",
	GDI_SDPSoox:     "SDPSoox",
	GDI_NOTSRCCOPY:  "Sn",
	GDI_SPDSaox:     "SPDSaox",
	GDI_SPDSxnox:    "SPDSxnox",
	GDI_SDPox:       "SDPox",
	GDI_SDPoan:      "SDPoan",
	GDI_PSDPoax:     "PSDPoax",
	GDI_SPDnox:      "SPDnox",
	GDI_SPDSxox:     "SPDSxox",
	GDI_SPDnoan:     "SPDnoan",
	GDI_PSx:         "PSx",
	GDI_SPDSonox:    "SPDSonox",
	GDI_SPDSnaox:    "SPDSnaox",
	GDI_PSan:        "PSan",
	GDI_PSDnaa:      "PSDnaa",
	GDI_DPSxon:      "DPSxon",
	GDI_SDxPDxa:     "SDxPDxa",
	GDI_SPDSanaxn:   "SPDSanaxn",
	GDI_SRCERASE:    "SDna",
	GDI_DPSnaon:     "DPSnaon",
	GDI_DSPDaox:     "DSPDaox",
	GDI_PSDPxaxn:    "PSDPxaxn",
	GDI_SDPxa:       "SDPxa",
	GDI_PDSPDaoxxn:  "PDSPDaoxxn",
	GDI_DPSDoax:     "DPSDoax",
	GDI_PDSnox:      "PDSnox",
	GDI_SDPana:      "SDPana",
	GDI_SSPxDSxoxn:  "SSPxDSxoxn",
	GDI_PDSPxox:     "PDSPxox",
	GDI_PDSnoan:     "PDSnoan",
	GDI_PDna:        "PDna",
	GDI_DSPnaon:     "DSPnaon",
	GDI_DPSDaox:     "DPSDaox",
	GDI_SPDSxaxn:    "SPDSxaxn",
	GDI_DPSonon:     "DPSonon",
	GDI_DSTINVERT:   "Dn",
	GDI_DPSox:       "DPSox",
	GDI_DPSoan:      "DPSoan",
	GDI_PDSPoax:     "PDSPoax",
	GDI_DPSnox:      "DPSnox",
	GDI_PATINVERT:   "DPx",
	GDI_DPSDonox:    "DPSDonox",
	GDI_DPSDxox:     "DPSDxox",
	GDI_DPSnoan:     "DPSnoan",
	GDI_DPSDnaox:    "DPSDnaox",
	GDI_DPan:        "DPan",
	GDI_PDSxa:       "PDSxa",
	GDI_DSPDSaoxxn:  "DSPDSaoxxn",
	GDI_DSPDoax:     "DSPDoax",
	GDI_SDPnox:      "SDPnox",
	GDI_SDPSoax:     "SDPSoax",
	GDI_DSPnox:      "DSPnox",
	GDI_SRCINVERT:   "DSx",
	GDI_SDPSonox:    "SDPSonox",
	GDI_DSPDSonoxxn: "DSPDSonoxxn",
	GDI_PDSxxn:      "PDSxxn",
	GDI_DPSax:       "DPSax",
	GDI_PSDPSoaxxn:  "PSDPSoaxxn",
	GDI_SDPax:       "SDPax",
	GDI_PDSPDoaxxn:  "PDSPDoaxxn",
	GDI_SDPSnoax:    "SDPSnoax",
	GDI_PDSxnan:     "PDSxnan",
	GDI_PDSana:      "PDSana",
	GDI_SSDxPDxaxn:  "SSDxPDxaxn",
	GDI_SDPSxox:     "SDPSxox",
	GDI_SDPnoan:     "SDPnoan",
	GDI_DSPDxox:     "DSPDxox",
	GDI_DSPnoan:     "DSPnoan",
	GDI_SDPSnaox:    "SDPSnaox",
	GDI_DSan:        "DSan",
	GDI_PDSax:       "PDSax",
	GDI_DSPDSoaxxn:  "DSPDSoaxxn",
	GDI_DPSDnoax:    "DPSDnoax",
	GDI_SDPxnan:     "SDPxnan",
	GDI_SPDSnoax:    "SPDSnoax",
	GDI_DPSxnan:     "DPSxnan",
	GDI_SPxDSxo:     "SPxDSxo",
	GDI_DPSaan:      "DPSaan",
	GDI_DPSaa:       "DPSaa",
	GDI_SPxDSxon:    "SPxDSxon",
	GDI_DPSxna:      "DPSxna",
	GDI_SPDSnoaxn:   "SPDSnoaxn",
	GDI_SDPxna:      "SDPxna",
	GDI_PDSPnoaxn:   "PDSPnoaxn",
	GDI_DSPDSoaxx:   "DSPDSoaxx",
	GDI_PDSaxn:      "PDSaxn",
	GDI_SRCAND:      "DSa",
	GDI_SDPSnaoxn:   "SDPSnaoxn",
	GDI_DSPnoa:      "DSPnoa",
	GDI_DSPDxoxn:    "DSPDxoxn",
	GDI_SDPnoa:      "SDPnoa",
	GDI_SDPSxoxn:    "SDPSxoxn",
	GDI_SSDxPDxax:   "SSDxPDxax",
	GDI_PDSanan:     "PDSanan",
	GDI_PDSxna:      "PDSxna",
	GDI_SDPSnoaxn:   "SDPSnoaxn",
	GDI_DPSDPoaxx:   "DPSDPoaxx",
	GDI_SPDaxn:      "SPDaxn",
	GDI_PSDPSoaxx:   "PSDPSoaxx",
	GDI_DPSaxn:      "DPSaxn",
	GDI_DPSxx:       "DPSxx",
	GDI_PSDPSonoxx:  "PSDPSonoxx",
	GDI_SDPSonoxn:   "SDPSonoxn",
	GDI_DSxn:        "DSxn",
	GDI_DPSnax:      "DPSnax",
	GDI_SDPSoaxn:    "SDPSoaxn",
	GDI_SPDnax:      "SPDnax",
	GDI_DSPDoaxn:    "DSPDoaxn",
	GDI_DSPDSaoxx:   "DSPDSaoxx",
	GDI_PDSxan:      "PDSxan",
	GDI_DPa:         "DPa",
	GDI_PDSPnaoxn:   "PDSPnaoxn",
	GDI_DPSnoa:      "DPSnoa",
	GDI_DPSDxoxn:    "DPSDxoxn",
	GDI_PDSPonoxn:   "PDSPonoxn",
	GDI_PDxn:        "PDxn",
	GDI_DSPnax:      "DSPnax",
	GDI_PDSPoaxn:    "PDSPoaxn",
	GDI_DPSoa:       "DPSoa",
	GDI_DPSoxn:      "DPSoxn",
	GDI_DSTCOPY:     "D",
	GDI_DPSono:      "DPSono",
	GDI_SPDSxax:     "SPDSxax",
	GDI_DPSDaoxn:    "DPSDaoxn",
	GDI_DSPnao:      "DSPnao",
	GDI_DPno:        "DPno",
	GDI_PDSnoa:      "PDSnoa",
	GDI_PDSPxoxn:    "PDSPxoxn",
	GDI_SSPxDSxox:   "SSPxDSxox",
	GDI_SDPanan:     "SDPanan",
	GDI_PSDnax:      "PSDnax",
	GDI_DPSDoaxn:    "DPSDoaxn",
	GDI_DPSDPaoxx:   "DPSDPaoxx",
	GDI_SDPxan:      "SDPxan",
	GDI_PSDPxax:     "PSDPxax",
	GDI_DSPDaoxn:    "DSPDaoxn",
	GDI_DPSnao:      "DPSnao",
	GDI_MERGEPAINT:  "DSno",
	GDI_SPDSanax:    "SPDSanax",
	GDI_SDxPDxan:    "SDxPDxan",
	GDI_DPSxo:       "DPSxo",
	GDI_DPSano:      "DPSano",
	GDI_MERGECOPY:   "PSa",
	GDI_SPDSnaoxn:   "SPDSnaoxn",
	GDI_SPDSonoxn:   "SPDSonoxn",
	GDI_PSxn:        "PSxn",
	GDI_SPDnoa:      "SPDnoa",
	GDI_SPDSxoxn:    "SPDSxoxn",
	GDI_SDPnax:      "SDPnax",
	GDI_PSDPoaxn:    "PSDPoaxn",
	GDI_SDPoa:       "SDPoa",
	GDI_SPDoxn:      "SPDoxn",
	GDI_DPSDxax:     "DPSDxax",
	GDI_SPDSaoxn:    "SPDSaoxn",
	GDI_SRCCOPY:     "S",
	GDI_SDPono:      "SDPono",
	GDI_SDPnao:      "SDPnao",
	GDI_SPno:        "SPno",
	GDI_PSDnoa:      "PSDnoa",
	GDI_PSDPxoxn:    "PSDPxoxn",
	GDI_PDSnax:      "PDSnax",
	GDI_SPDSoaxn:    "SPDSoaxn",
	GDI_SSPxPDxax:   "SSPxPDxax",
	GDI_DPSanan:     "DPSanan",
	GDI_PSDPSaoxx:   "PSDPSaoxx",
	GDI_DPSxan:      "DPSxan",
	GDI_PDSPxax:     "PDSPxax",
	GDI_SDPSaoxn:    "SDPSaoxn",
	GDI_DPSDanax:    "DPSDanax",
	GDI_SPxDSxan:    "SPxDSxan",
	GDI_SPDnao:      "SPDnao",
	GDI_SDno:        "SDno",
	GDI_SDPxo:       "SDPxo",
	GDI_SDPano:      "SDPano",
	GDI_PDSoa:       "PDSoa",
	GDI_PDSoxn:      "PDSoxn",
	GDI_DSPDxax:     "DSPDxax",
	GDI_PSDPaoxn:    "PSDPaoxn",
	GDI_SDPSxax:     "SDPSxax",
	GDI_PDSPaoxn:    "PDSPaoxn",
	GDI_SDPSanax:    "SDPSanax",
	GDI_SPxPDxan:    "SPxPDxan",
	GDI_SSPxDSxax:   "SSPxDSxax",
	GDI_DSPDSanaxxn: "DSPDSanaxxn",
	GDI_DPSao:       "DPSao",
	GDI_DPSxno:      "DPSxno",
	GDI_SDPao:       "SDPao",
	GDI_SDPxno:      "SDPxno",
	GDI_SRCPAINT:    "DSo",
	GDI_SDPnoo:      "SDPnoo",
	GDI_PATCOPY:     "P",
	GDI_PDSono:      "PDSono",
	GDI_PDSnao:      "PDSnao",
	GDI_PSno:        "PSno",
	GDI_PSDnao:      "PSDnao",
	GDI_PDno:        "PDno",
	GDI_PDSxo:       "PDSxo",
	GDI_PDSano:      "PDSano",
	GDI_PDSao:       "PDSao",
	GDI_PDSxno:      "PDSxno",
	GDI_DPo:         "DPo",
	GDI_PATPAINT:    "DPSnoo",
	GDI_PSo:         "PSo",
	GDI_PSDnoo:      "PSDnoo",
	GDI_DPSoo:       "DPSoo",
	GDI_WHITENESS:   "1",
}

var rop2Table = [16]int{
	GDI_R2_BLACK,
	GDI_R2_NOTMERGEPEN,
	GDI_R2_MASKNOTPEN,
	GDI_R2_NOTCOPYPEN,
	GDI_R2_MASKPENNOT,
	GDI_R2_NOT,
	GDI_R2_XORPEN,
	GDI_R2_NOTMASKPEN,
	GDI_R2_MASKPEN,
	GDI_R2_NOTXORPEN,
	GDI_R2_NOP,
	GDI_R2_MERGENOTPEN,
	GDI_R2_COPYPEN,
	GDI_R2_MERGEPENNOT,
	GDI_R2_MERGEPEN,
	GDI_R2_WHITE,
}
