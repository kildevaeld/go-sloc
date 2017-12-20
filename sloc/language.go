package sloc

import "path"

type Commenter struct {
	LineComment  string
	StartComment string
	EndComment   string
	Nesting      bool
}

var (
	noComments     = Commenter{"\000", "\000", "\000", false}
	xmlComments    = Commenter{"\000", `<!--`, `-->`, false}
	cComments      = Commenter{`//`, `/*`, `*/`, false}
	cssComments    = Commenter{"\000", `/*`, `*/`, false}
	shComments     = Commenter{`#`, "\000", "\000", false}
	semiComments   = Commenter{`;`, "\000", "\000", false}
	hsComments     = Commenter{`--`, `{-`, `-}`, true}
	mlComments     = Commenter{`\000`, `(*`, `*)`, false}
	sqlComments    = Commenter{`--`, `/*`, `*/`, false}
	luaComments    = Commenter{`--`, `--[[`, `]]`, false}
	pyComments     = Commenter{`#`, `"""`, `"""`, false}
	matlabComments = Commenter{`%`, `%{`, `%}`, false}
	erlangComments = Commenter{`%`, "\000", "\000", false}
	rubyComments   = Commenter{`#`, "=begin", "=end", false}
	coffeeComments = Commenter{`#`, "###", "###", false}
	swiftComments  = Commenter{`//`, `/*`, `*/`, true}
	yamlComments   = Commenter{`#`, "\000", "\000", false}
	// TODO support POD and __END__
	perlComments = Commenter{`#`, "\000", "\000", false}
)

type Namer string

func (l Namer) Name() string { return string(l) }

type Matcher func(string) bool

func (m Matcher) Match(fname string) bool { return m(fname) }

func mExt(exts ...string) Matcher {
	return func(fname string) bool {
		for _, ext := range exts {
			if ext == path.Ext(fname) {
				return true
			}
		}
		return false
	}
}

func mName(names ...string) Matcher {
	return func(fname string) bool {
		for _, name := range names {
			if name == path.Base(fname) {
				return true
			}
		}
		return false
	}
}

type Language struct {
	Namer
	Matcher
	Commenter
}

// TODO work properly with unicode
func (l Language) Update(c []byte, s *Stats) {
	s.FileCount++

	commentLen := 0
	codeLen := 0
	inComment := 0 // this is an int for nesting
	inLComment := false
	lc := []byte(l.LineComment)
	sc := []byte(l.StartComment)
	ec := []byte(l.EndComment)
	lp, sp, ep := 0, 0, 0

	for _, b := range c {
		if b != byte(' ') && b != byte('\t') && b != byte('\n') {
			if !inLComment && inComment == 0 {
				codeLen++
			} else {
				commentLen++
			}
		}
		if inComment == 0 && b == lc[lp] {
			lp++
			if lp == len(lc) {
				if !inLComment {
					codeLen -= lp
				}
				inLComment = true
				lp = 0
			}
		} else {
			lp = 0
		}
		if !inLComment && b == sc[sp] {
			sp++
			if sp == len(sc) {
				if inComment == 0 {
					codeLen -= sp
				}
				inComment++
				if inComment > 1 && !l.Nesting {
					inComment = 1
				}
				sp = 0
			}
		} else {
			sp = 0
		}
		if !inLComment && inComment > 0 && b == ec[ep] {
			ep++
			if ep == len(ec) {
				if inComment > 0 {
					inComment--
				}
				if inComment == 0 {
					commentLen -= ep
				}
				ep = 0
			}
		} else {
			ep = 0
		}

		if b == byte('\n') {
			s.TotalLines++
			if commentLen > 0 {
				s.CommentLines++
			}
			if codeLen > 0 {
				s.CodeLines++
			}
			if commentLen == 0 && codeLen == 0 {
				s.BlankLines++
			}
			inLComment = false
			codeLen = 0
			commentLen = 0
			continue
		}
	}
}

var languages = []Language{
	Language{"Thrift", mExt(".thrift"), cComments},

	Language{"C", mExt(".c", ".h"), cComments},
	Language{"C++", mExt(".cc", ".cpp", ".cxx", ".hh", ".hpp", ".hxx"), cComments},
	Language{"Go", mExt(".go"), cComments},
	Language{"Rust", mExt(".rs", ".rc"), cComments},
	Language{"Scala", mExt(".scala"), cComments},
	Language{"Java", mExt(".java"), cComments},

	Language{"YACC", mExt(".y"), cComments},
	Language{"Lex", mExt(".l"), cComments},

	Language{"Lua", mExt(".lua"), luaComments},

	Language{"SQL", mExt(".sql"), sqlComments},

	Language{"Haskell", mExt(".hs", ".lhs"), hsComments},
	Language{"ML", mExt(".ml", ".mli"), mlComments},

	Language{"Perl", mExt(".pl", ".pm"), perlComments},
	Language{"PHP", mExt(".php"), cComments},

	Language{"Shell", mExt(".sh"), shComments},
	Language{"Bash", mExt(".bash"), shComments},
	Language{"R", mExt(".r", ".R"), shComments},
	Language{"Tcl", mExt(".tcl"), shComments},

	Language{"MATLAB", mExt(".m"), matlabComments},

	Language{"Ruby", mExt(".rb"), rubyComments},
	Language{"Python", mExt(".py"), pyComments},
	Language{"Assembly", mExt(".asm", ".s"), semiComments},
	Language{"Lisp", mExt(".lsp", ".lisp"), semiComments},
	Language{"Scheme", mExt(".scm", ".scheme"), semiComments},

	Language{"Make", mName("makefile", "Makefile", "MAKEFILE"), shComments},
	Language{"CMake", mName("CMakeLists.txt"), shComments},
	Language{"Jam", mName("Jamfile", "Jamrules"), shComments},

	Language{"Markdown", mExt(".md"), noComments},

	Language{"HAML", mExt(".haml"), noComments},
	Language{"SASS", mExt(".sass"), cssComments},
	Language{"SCSS", mExt(".scss"), cssComments},

	Language{"HTML", mExt(".htm", ".html", ".xhtml"), xmlComments},
	Language{"XML", mExt(".xml"), xmlComments},
	Language{"CSS", mExt(".css"), cssComments},
	Language{"JavaScript", mExt(".js", ".jsx"), cComments},
	Language{"Qml", mExt(".qml"), cComments},
	Language{"CoffeeScript", mExt(".coffee"), coffeeComments},
	Language{"Typescript", mExt(".ts", ".tsx"), cComments},

	Language{"Swift", mExt(".swift"), swiftComments},
	Language{"Erlang", mExt(".erl"), erlangComments},
	Language{"Yaml", mExt(".yml", ".yaml"), yamlComments},
}
