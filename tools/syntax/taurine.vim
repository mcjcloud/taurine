" Vim syntax file
" Language: Taurine
" Maintainer: Brayden Cloud
" Latest Revision: 19 December 2021

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

syn match closeParenthesisRe ')'

" keywords
syn keyword basicLanguageKeywords var while if else etch len func return
syn keyword basicTypes num str bool obj arr func contained nextgroup=closeParenthesisRe skipwhite

syn region stringRe start='"' end='"' 
syn region typeRe start='(' end=')' transparent contains=basicTypes,func

hi def link commentRe Comment
hi def link basicTypes Type
hi def link basicLanguageKeywords Preproc
hi def link typeRe Type
hi def link stringRe Constant
hi def link numberRe Constant

let b:current_syntax = 'taurine'

