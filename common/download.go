package common

import "github.com/ItsMeSamey/go_fuzzy"

type LanguageID string

const (
  LangAN    LanguageID = "an"    // Aragonese
  LangAR    LanguageID = "ar"    // Arabic
  LangAT    LanguageID = "at"    // Asturian (non standard?)
  LangBG    LanguageID = "bg"    // Bulgarian
  LangBR    LanguageID = "br"    // Breton
  LangCA    LanguageID = "ca"    // Catalan
  LangCS    LanguageID = "cs"    // Czech
  LangDA    LanguageID = "da"    // Danish
  LangDE    LanguageID = "de"    // German
  LangEL    LanguageID = "el"    // Greek
  LangEN    LanguageID = "en"    // English
  LangEO    LanguageID = "eo"    // Esperanto
  LangES    LanguageID = "es"    // Spanish
  LangET    LanguageID = "et"    // Estonian
  LangEU    LanguageID = "eu"    // Basque
  LangFA    LanguageID = "fa"    // Persian
  LangFI    LanguageID = "fi"    // Finnish
  LangFR    LanguageID = "fr"    // French
  LangGL    LanguageID = "gl"    // Galician
  LangHE    LanguageID = "he"    // Hebrew
  LangHI    LanguageID = "hi"    // Hindi
  LangHR    LanguageID = "hr"    // Croatian
  LangHU    LanguageID = "hu"    // Hungarian
  LangHY    LanguageID = "hy"    // Armenian
  LangID    LanguageID = "id"    // Indonesian
  LangIS    LanguageID = "is"    // Icelandic
  LangIT    LanguageID = "it"    // Italian
  LangJA    LanguageID = "ja"    // Japanese
  LangKA    LanguageID = "ka"    // Georgian
  LangKM    LanguageID = "km"    // Khmer
  LangKO    LanguageID = "ko"    // Korean
  LangMK    LanguageID = "mk"    // Macedonian
  LangMS    LanguageID = "ms"    // Malay
  LangNL    LanguageID = "nl"    // Dutch
  LangNO    LanguageID = "no"    // Norwegian
  LangOC    LanguageID = "oc"    // Occitan
  LangPTBR  LanguageID = "pt-br" // Portuguese, Brazilian
  LangPL    LanguageID = "pl"    // Polish
  LangPT    LanguageID = "pt"    // Portuguese
  LangRO    LanguageID = "ro"    // Romanian
  LangRU    LanguageID = "ru"    // Russian
  LangSI    LanguageID = "si"    // Sinhala
  LangSK    LanguageID = "sk"    // Slovak
  LangSL    LanguageID = "sl"    // Slovenian
  LangSQ    LanguageID = "sq"    // Albanian
  LangSR    LanguageID = "sr"    // Serbian
  LangSV    LanguageID = "sv"    // Swedish
  LangTH    LanguageID = "th"    // Thai
  LangTL    LanguageID = "tl"    // Tagalog
  LangTR    LanguageID = "tr"    // Turkish
  LangTT    LanguageID = "tt"    // Tatar
  LangUK    LanguageID = "uk"    // Ukrainian
  LangUZ    LanguageID = "uz"    // Uzbek
  LangVI    LanguageID = "vi"    // Vietnamese
  LangZH    LanguageID = "zh"    // Chinese
  LangZHTW  LanguageID = "zh-tw" // Chinese Traditional
)

var LanguageNameMap = map[LanguageID]string{
  LangAN:   "Aragonese",
  LangAR:   "Arabic",
  LangAT:   "Asturian",
  LangBG:   "Bulgarian",
  LangBR:   "Breton",
  LangCA:   "Catalan",
  LangCS:   "Czech",
  LangDA:   "Danish",
  LangDE:   "German",
  LangEL:   "Greek",
  LangEN:   "English",
  LangEO:   "Esperanto",
  LangES:   "Spanish",
  LangET:   "Estonian",
  LangEU:   "Basque",
  LangFA:   "Persian",
  LangFI:   "Finnish",
  LangFR:   "French",
  LangGL:   "Galician",
  LangHE:   "Hebrew",
  LangHI:   "Hindi",
  LangHR:   "Croatian",
  LangHU:   "Hungarian",
  LangHY:   "Armenian",
  LangID:   "Indonesian",
  LangIS:   "Icelandic",
  LangIT:   "Italian",
  LangJA:   "Japanese",
  LangKA:   "Georgian",
  LangKM:   "Khmer",
  LangKO:   "Korean",
  LangMK:   "Macedonian",
  LangMS:   "Malay",
  LangNL:   "Dutch",
  LangNO:   "Norwegian",
  LangOC:   "Occitan",
  LangPTBR: "Brazilian Portuguese",
  LangPL:   "Polish",
  LangPT:   "Portuguese",
  LangRO:   "Romanian",
  LangRU:   "Russian",
  LangSI:   "Sinhala",
  LangSK:   "Slovak",
  LangSL:   "Slovenian",
  LangSQ:   "Albanian",
  LangSR:   "Serbian",
  LangSV:   "Swedish",
  LangTH:   "Thai",
  LangTL:   "Tagalog",
  LangTR:   "Turkish",
  LangTT:   "Tatar",
  LangUK:   "Ukrainian",
  LangUZ:   "Uzbek",
  LangVI:   "Vietnamese",
  LangZH:   "Chinese",
  LangZHTW: "Traditional Chinese",
}

type SearchOptions struct {
  Language LanguageID
  Sorter   fuzzy.Sorter[float32, string, string]
}
type MovieListData struct {
  // The title of the movie
  Title   string

  // The Options when frtching the list of subtitles for this movie
  Options SearchOptions
}
type MovieListEntry interface {
  Data() *MovieListData
  ToSubtitleLinks() ([]SubtitleListEntry, error)
}

type SubtitleListData struct {
  Parent MovieListEntry
  // Name of the subtitle file
  Filename string
  // Language of the subtitle file
  Language string
  // Setting this to "" will force refetching when calling DownloadLink()
  Target string
}
type SubtitleListEntry interface {
  Data() *SubtitleListData
  // Weather the downloaded file is a zip file or not
  IsZip() bool
  // Returns the Download link from where we can fetch the subtitle file
  DownloadLink() (string, error)
}


type DownloadedSubtitleEntry struct {
  Subtitle []byte
  Filename  string
}
type DownloadedSubtitle struct {
  Parent SubtitleListEntry

  // May contain 0, or 1 or more subtitles
  Subtitles []DownloadedSubtitleEntry
}

