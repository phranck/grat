# Fonts

These files are served from this repository rather than fetched from a font
service, so a visitor's browser talks to nobody but the host serving the site.

| File | Family | Weight | Source |
| --- | --- | --- | --- |
| `barlow-400.woff2` | Barlow | 400 | Google Fonts, latin subset |
| `barlow-600.woff2` | Barlow | 600 | Google Fonts, latin subset |
| `barlow-condensed-600.woff2` | Barlow Condensed | 600 | Google Fonts, latin subset |
| `barlow-condensed-700.woff2` | Barlow Condensed | 700 | Google Fonts, latin subset |
| `fira-code-300-700.woff2` | Fira Code | 300 to 700, variable | Google Fonts, latin subset |

Both families are published under the SIL Open Font License, Version 1.1, which
allows redistribution as long as the licence travels with the files.
`Barlow-OFL.txt` and `FiraCode-OFL.txt` are those licences.

Only the latin subset is here, because the site is written in English. One
character it uses falls outside it, the check mark `U+2713` in the terminal
transcripts, and that has always come from a system font rather than from Fira
Code, whose Google subsets do not carry it either.
