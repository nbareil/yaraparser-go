# yaraparser-go


[![Coverage Status](https://coveralls.io/repos/github/nbareil/yaraparser-go/badge.svg?branch=master)](https://coveralls.io/github/nbareil/yaraparser-go?branch=master)
[![Build Status](https://travis-ci.org/nbareil/yaraparser-go.svg?branch=master)](https://travis-ci.org/nbareil/yaraparser-go#)


Scanning and parsing yara files without lex/yacc grammar.

⚠️  This project is just for fun, it is not (intented to be) finished. If you want something serious, use [Yara-Rules/yago](https://github.com/Yara-Rules/yago) from [@Xumeiquer](https://github.com/Xumeiquer)!


# Status

- [X] Support Meta section
- [ ] Support Strings section
  - [X] Parse regular strings
  - [X] Parse Hex format strings and patterns
  - [ ] Parse modifiers
  - [ ] Support regexp rules
    - [ ] Parse Regexp
- [ ] Support Condition section
  - [ ] Check boolean function
  - [ ] Check semantic
- [ ] Parse comments
  - [X] Skip them
  - [ ] Attach them to where they belong

That's where I ended up with a few hours of free time. I don't know if I will have more time like this in the future, so consider this status as the final one :)

