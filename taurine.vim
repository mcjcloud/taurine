" Vim syntax file
" Language: Celestia Star Catalogs
" Maintainer: Kevin Lauder
" Latest Revision: 26 April 2008

if exists("b:current_syntax")
  finish
endif

" Integer with - + or nothing in front
syn match numberRe '\d\+'
syn match numberRe '[-+]\d\+'

" Floating point number with decimal no E or e 
syn match numberRe '[-+]\d\+\.\d*'

" Floating point like number with E and no decimal point (+,-)
syn match numberRe '[-+]\=\d[[:digit:]]*[eE][\-+]\=\d\+'
syn match numberRe '\d[[:digit:]]*[eE][\-+]\=\d\+'

" Floating point like number with E and decimal point (+,-)
syn match numberRe '[-+]\=\d[[:digit:]]*\.\d*[eE][\-+]\=\d\+'
syn match numberRe '\d[[:digit:]]*\.\d*[eE][\-+]\=\d\+'

" comment
syn match commentRe '\/\/.*$'

syn keyword basicLanguageKeywords var while if else etch len
syn keyword basicTypes func num str bool

syn region stringRe start='"' end='"' contained

let b:current_syntax = "taurine"

hi def link commentRe Comment
hi def link basicLanguageKeywords Constant
hi def link basicTypes Type
hi def link stringRe Constant
hi def link numberRe Constant

