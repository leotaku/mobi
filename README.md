# Go MOBI

[![Github Actions CI](https://github.com/leotaku/mobi/workflows/check/badge.svg)](https://github.com/leotaku/mobi/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/leotaku/mobi)](https://goreportcard.com/report/github.com/leotaku/mobi)
[![Go Reference](https://pkg.go.dev/badge/github.com/leotaku/mobi.svg)](https://pkg.go.dev/github.com/leotaku/mobi)

This package implements facilities to create KF8-formatted MOBI and AZW3 books.
We also export the raw PalmDB writer and various PalmDoc, MOBI and KF8 components as subpackages, which can be used to implement other formats that build on these standards.

## Known issues

+ Chapters are supported but subchapters are not
+ Old readers without KF8 are not supported (Kindle 1, 2 and DX)
+ Books without any text content are always malformed
+ Errors during template expansion result in a panic

## References

+ MobileRead Wiki
  + [MOBI format](https://wiki.mobileread.com/wiki/MOBI)
  + [PDB format](https://wiki.mobileread.com/wiki/PDB)
  + [KF8 format](https://wiki.mobileread.com/wiki/KF8)
+ Calibre source code for [MOBI support](https://github.com/kovidgoyal/calibre/tree/master/src/calibre/ebooks/mobi)
+ Vladimir Konovalov's [Golang MOBI writer](https://github.com/766b/mobi)
+ Library of Congress on the [MOBI format](https://www.loc.gov/preservation/digital/formats/fdd/fdd000472.shtml)

Huge thanks to the authors of these resources.

## License

[MIT](./LICENSE) © Leo Gaskin 2020-2024
