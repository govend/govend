package main

var stdpkgs = map[string][]pkg{
	"io": []pkg{
		pkg{importpath: "io", dir: "/go/src/io"},
	},
	"x86asm": []pkg{
		pkg{importpath: "cmd/internal/rsc.io/x86/x86asm", dir: "/go/src/cmd/internal/rsc.io/x86/x86asm"},
	},
	"hash": []pkg{
		pkg{importpath: "hash", dir: "/go/src/hash"},
	},
	"image": []pkg{
		pkg{importpath: "image", dir: "/go/src/image"},
	},
	"printer": []pkg{
		pkg{importpath: "go/printer", dir: "/go/src/go/printer"},
	},
	"tabwriter": []pkg{
		pkg{importpath: "text/tabwriter", dir: "/go/src/text/tabwriter"},
	},
	"tempfile": []pkg{
		pkg{importpath: "cmd/pprof/internal/tempfile", dir: "/go/src/cmd/pprof/internal/tempfile"},
	},
	"sha256": []pkg{
		pkg{importpath: "crypto/sha256", dir: "/go/src/crypto/sha256"},
	},
	"gob": []pkg{
		pkg{importpath: "encoding/gob", dir: "/go/src/encoding/gob"},
	},
	"binary": []pkg{
		pkg{importpath: "encoding/binary", dir: "/go/src/encoding/binary"},
	},
	"rpc": []pkg{
		pkg{importpath: "net/rpc", dir: "/go/src/net/rpc"},
	},
	"url": []pkg{
		pkg{importpath: "net/url", dir: "/go/src/net/url"},
	},
	"syntax": []pkg{
		pkg{importpath: "regexp/syntax", dir: "/go/src/regexp/syntax"},
	},
	"archive": []pkg{
		pkg{importpath: "archive", dir: "/go/src/archive"},
	},
	"objdump": []pkg{
		pkg{importpath: "cmd/objdump", dir: "/go/src/cmd/objdump"},
	},
	"pack": []pkg{
		pkg{importpath: "cmd/pack", dir: "/go/src/cmd/pack"},
	},
	"ast": []pkg{
		pkg{importpath: "go/ast", dir: "/go/src/go/ast"},
	},
	"jpeg": []pkg{
		pkg{importpath: "image/jpeg", dir: "/go/src/image/jpeg"},
	},
	"user": []pkg{
		pkg{importpath: "os/user", dir: "/go/src/os/user"},
	},
	"syslog": []pkg{
		pkg{importpath: "log/syslog", dir: "/go/src/log/syslog"},
	},
	"cgi": []pkg{
		pkg{importpath: "net/http/cgi", dir: "/go/src/net/http/cgi"},
	},
	"reflect": []pkg{
		pkg{importpath: "reflect", dir: "/go/src/reflect"},
	},
	"symbolz": []pkg{
		pkg{importpath: "cmd/pprof/internal/symbolz", dir: "/go/src/cmd/pprof/internal/symbolz"},
	},
	"fmt": []pkg{
		pkg{importpath: "fmt", dir: "/go/src/fmt"},
		pkg{importpath: "lib9/fmt", dir: "/go/src/lib9/fmt"},
	},
	"png": []pkg{
		pkg{importpath: "image/png", dir: "/go/src/image/png"},
	},
	"debug": []pkg{
		pkg{importpath: "debug", dir: "/go/src/debug"},
		pkg{importpath: "runtime/debug", dir: "/go/src/runtime/debug"},
	},
	"smtp": []pkg{
		pkg{importpath: "net/smtp", dir: "/go/src/net/smtp"},
	},
	"utf8": []pkg{
		pkg{importpath: "unicode/utf8", dir: "/go/src/unicode/utf8"},
	},
	"zlib": []pkg{
		pkg{importpath: "compress/zlib", dir: "/go/src/compress/zlib"},
	},
	"container": []pkg{
		pkg{importpath: "container", dir: "/go/src/container"},
	},
	"pkix": []pkg{
		pkg{importpath: "crypto/x509/pkix", dir: "/go/src/crypto/x509/pkix"},
	},
	"net": []pkg{
		pkg{importpath: "net", dir: "/go/src/net"},
	},
	"builtin": []pkg{
		pkg{importpath: "builtin", dir: "/go/src/builtin"},
	},
	"dist": []pkg{
		pkg{importpath: "cmd/dist", dir: "/go/src/cmd/dist"},
	},
	"profile": []pkg{
		pkg{importpath: "cmd/pprof/internal/profile", dir: "/go/src/cmd/pprof/internal/profile"},
	},
	"runtime": []pkg{
		pkg{importpath: "runtime", dir: "/go/src/runtime"},
	},
	"race": []pkg{
		pkg{importpath: "runtime/race", dir: "/go/src/runtime/race"},
	},
	"quick": []pkg{
		pkg{importpath: "testing/quick", dir: "/go/src/testing/quick"},
	},
	"xml": []pkg{
		pkg{importpath: "encoding/xml", dir: "/go/src/encoding/xml"},
	},
	"cookiejar": []pkg{
		pkg{importpath: "net/http/cookiejar", dir: "/go/src/net/http/cookiejar"},
	},
	"mail": []pkg{
		pkg{importpath: "net/mail", dir: "/go/src/net/mail"},
	},
	"atomic": []pkg{
		pkg{importpath: "sync/atomic", dir: "/go/src/sync/atomic"},
	},
	"text": []pkg{
		pkg{importpath: "text", dir: "/go/src/text"},
	},
	"commands": []pkg{
		pkg{importpath: "cmd/pprof/internal/commands", dir: "/go/src/cmd/pprof/internal/commands"},
	},
	"dwarf": []pkg{
		pkg{importpath: "debug/dwarf", dir: "/go/src/debug/dwarf"},
	},
	"lib9": []pkg{
		pkg{importpath: "lib9", dir: "/go/src/lib9"},
	},
	"rc4": []pkg{
		pkg{importpath: "crypto/rc4", dir: "/go/src/crypto/rc4"},
	},
	"sha1": []pkg{
		pkg{importpath: "crypto/sha1", dir: "/go/src/crypto/sha1"},
	},
	"macho": []pkg{
		pkg{importpath: "debug/macho", dir: "/go/src/debug/macho"},
	},
	"sort": []pkg{
		pkg{importpath: "sort", dir: "/go/src/sort"},
	},
	"unicode": []pkg{
		pkg{importpath: "unicode", dir: "/go/src/unicode"},
	},
	"cmd": []pkg{
		pkg{importpath: "cmd", dir: "/go/src/cmd"},
	},
	"list": []pkg{
		pkg{importpath: "container/list", dir: "/go/src/container/list"},
	},
	"rand": []pkg{
		pkg{importpath: "crypto/rand", dir: "/go/src/crypto/rand"},
		pkg{importpath: "math/rand", dir: "/go/src/math/rand"},
	},
	"hex": []pkg{
		pkg{importpath: "encoding/hex", dir: "/go/src/encoding/hex"},
	},
	"suffixarray": []pkg{
		pkg{importpath: "index/suffixarray", dir: "/go/src/index/suffixarray"},
	},
	"arm": []pkg{
		pkg{importpath: "cmd/internal/rsc.io/arm", dir: "/go/src/cmd/internal/rsc.io/arm"},
	},
	"nm": []pkg{
		pkg{importpath: "cmd/nm", dir: "/go/src/cmd/nm"},
	},
	"dsa": []pkg{
		pkg{importpath: "crypto/dsa", dir: "/go/src/crypto/dsa"},
	},
	"expvar": []pkg{
		pkg{importpath: "expvar", dir: "/go/src/expvar"},
	},
	"flag": []pkg{
		pkg{importpath: "flag", dir: "/go/src/flag"},
	},
	"cmplx": []pkg{
		pkg{importpath: "math/cmplx", dir: "/go/src/math/cmplx"},
	},
	"heap": []pkg{
		pkg{importpath: "container/heap", dir: "/go/src/container/heap"},
	},
	"sha512": []pkg{
		pkg{importpath: "crypto/sha512", dir: "/go/src/crypto/sha512"},
	},
	"json": []pkg{
		pkg{importpath: "encoding/json", dir: "/go/src/encoding/json"},
	},
	"ld": []pkg{
		pkg{importpath: "cmd/ld", dir: "/go/src/cmd/ld"},
	},
	"strconv": []pkg{
		pkg{importpath: "strconv", dir: "/go/src/strconv"},
	},
	"utf16": []pkg{
		pkg{importpath: "unicode/utf16", dir: "/go/src/unicode/utf16"},
	},
	"pe": []pkg{
		pkg{importpath: "debug/pe", dir: "/go/src/debug/pe"},
	},
	"base64": []pkg{
		pkg{importpath: "encoding/base64", dir: "/go/src/encoding/base64"},
	},
	"httputil": []pkg{
		pkg{importpath: "net/http/httputil", dir: "/go/src/net/http/httputil"},
	},
	"os": []pkg{
		pkg{importpath: "os", dir: "/go/src/os"},
	},
	"signal": []pkg{
		pkg{importpath: "os/signal", dir: "/go/src/os/signal"},
	},
	"elliptic": []pkg{
		pkg{importpath: "crypto/elliptic", dir: "/go/src/crypto/elliptic"},
	},
	"hmac": []pkg{
		pkg{importpath: "crypto/hmac", dir: "/go/src/crypto/hmac"},
	},
	"rsa": []pkg{
		pkg{importpath: "crypto/rsa", dir: "/go/src/crypto/rsa"},
	},
	"strings": []pkg{
		pkg{importpath: "strings", dir: "/go/src/strings"},
	},
	"zip": []pkg{
		pkg{importpath: "archive/zip", dir: "/go/src/archive/zip"},
	},
	"base32": []pkg{
		pkg{importpath: "encoding/base32", dir: "/go/src/encoding/base32"},
	},
	"textproto": []pkg{
		pkg{importpath: "net/textproto", dir: "/go/src/net/textproto"},
	},
	"subtle": []pkg{
		pkg{importpath: "crypto/subtle", dir: "/go/src/crypto/subtle"},
	},
	"errors": []pkg{
		pkg{importpath: "errors", dir: "/go/src/errors"},
	},
	"crc64": []pkg{
		pkg{importpath: "hash/crc64", dir: "/go/src/hash/crc64"},
	},
	"x509": []pkg{
		pkg{importpath: "crypto/x509", dir: "/go/src/crypto/x509"},
	},
	"gosym": []pkg{
		pkg{importpath: "debug/gosym", dir: "/go/src/debug/gosym"},
	},
	"format": []pkg{
		pkg{importpath: "go/format", dir: "/go/src/go/format"},
	},
	"jsonrpc": []pkg{
		pkg{importpath: "net/rpc/jsonrpc", dir: "/go/src/net/rpc/jsonrpc"},
	},
	"filepath": []pkg{
		pkg{importpath: "path/filepath", dir: "/go/src/path/filepath"},
	},
	"parse": []pkg{
		pkg{importpath: "text/template/parse", dir: "/go/src/text/template/parse"},
	},
	"api": []pkg{
		pkg{importpath: "cmd/api", dir: "/go/src/cmd/api"},
	},
	"driver": []pkg{
		pkg{importpath: "cmd/pprof/internal/driver", dir: "/go/src/cmd/pprof/internal/driver"},
		pkg{importpath: "database/sql/driver", dir: "/go/src/database/sql/driver"},
	},
	"big": []pkg{
		pkg{importpath: "math/big", dir: "/go/src/math/big"},
	},
	"index": []pkg{
		pkg{importpath: "index", dir: "/go/src/index"},
	},
	"syscall": []pkg{
		pkg{importpath: "internal/syscall", dir: "/go/src/internal/syscall"},
		pkg{importpath: "syscall", dir: "/go/src/syscall"},
	},
	"ecdsa": []pkg{
		pkg{importpath: "crypto/ecdsa", dir: "/go/src/crypto/ecdsa"},
	},
	"build": []pkg{
		pkg{importpath: "go/build", dir: "/go/src/go/build"},
	},
	"palette": []pkg{
		pkg{importpath: "image/color/palette", dir: "/go/src/image/color/palette"},
	},
	"pprof": []pkg{
		pkg{importpath: "cmd/pprof", dir: "/go/src/cmd/pprof"},
		pkg{importpath: "net/http/pprof", dir: "/go/src/net/http/pprof"},
		pkg{importpath: "runtime/pprof", dir: "/go/src/runtime/pprof"},
	},
	"ring": []pkg{
		pkg{importpath: "container/ring", dir: "/go/src/container/ring"},
	},
	"parser": []pkg{
		pkg{importpath: "go/parser", dir: "/go/src/go/parser"},
	},
	"sync": []pkg{
		pkg{importpath: "sync", dir: "/go/src/sync"},
	},
	"bufio": []pkg{
		pkg{importpath: "bufio", dir: "/go/src/bufio"},
	},
	"objfile": []pkg{
		pkg{importpath: "cmd/internal/objfile", dir: "/go/src/cmd/internal/objfile"},
	},
	"rsc.io": []pkg{
		pkg{importpath: "cmd/internal/rsc.io", dir: "/go/src/cmd/internal/rsc.io"},
	},
	"des": []pkg{
		pkg{importpath: "crypto/des", dir: "/go/src/crypto/des"},
	},
	"elf": []pkg{
		pkg{importpath: "debug/elf", dir: "/go/src/debug/elf"},
	},
	"template": []pkg{
		pkg{importpath: "html/template", dir: "/go/src/html/template"},
		pkg{importpath: "text/template", dir: "/go/src/text/template"},
	},
	"bytes": []pkg{
		pkg{importpath: "bytes", dir: "/go/src/bytes"},
	},
	"internal": []pkg{
		pkg{importpath: "cmd/internal", dir: "/go/src/cmd/internal"},
		pkg{importpath: "cmd/pprof/internal", dir: "/go/src/cmd/pprof/internal"},
		pkg{importpath: "internal", dir: "/go/src/internal"},
		pkg{importpath: "net/http/internal", dir: "/go/src/net/http/internal"},
	},
	"gzip": []pkg{
		pkg{importpath: "compress/gzip", dir: "/go/src/compress/gzip"},
	},
	"gif": []pkg{
		pkg{importpath: "image/gif", dir: "/go/src/image/gif"},
	},
	"ioutil": []pkg{
		pkg{importpath: "io/ioutil", dir: "/go/src/io/ioutil"},
	},
	"exec": []pkg{
		pkg{importpath: "os/exec", dir: "/go/src/os/exec"},
	},
	"testing": []pkg{
		pkg{importpath: "testing", dir: "/go/src/testing"},
	},
	"cc": []pkg{
		pkg{importpath: "cmd/cc", dir: "/go/src/cmd/cc"},
	},
	"compress": []pkg{
		pkg{importpath: "compress", dir: "/go/src/compress"},
	},
	"crypto": []pkg{
		pkg{importpath: "crypto", dir: "/go/src/crypto"},
	},
	"cipher": []pkg{
		pkg{importpath: "crypto/cipher", dir: "/go/src/crypto/cipher"},
	},
	"ascii85": []pkg{
		pkg{importpath: "encoding/ascii85", dir: "/go/src/encoding/ascii85"},
	},
	"token": []pkg{
		pkg{importpath: "go/token", dir: "/go/src/go/token"},
	},
	"http": []pkg{
		pkg{importpath: "net/http", dir: "/go/src/net/http"},
	},
	"httptest": []pkg{
		pkg{importpath: "net/http/httptest", dir: "/go/src/net/http/httptest"},
	},
	"addr2line": []pkg{
		pkg{importpath: "cmd/addr2line", dir: "/go/src/cmd/addr2line"},
	},
	"fix": []pkg{
		pkg{importpath: "cmd/fix", dir: "/go/src/cmd/fix"},
	},
	"flate": []pkg{
		pkg{importpath: "compress/flate", dir: "/go/src/compress/flate"},
	},
	"unsafe": []pkg{
		pkg{importpath: "unsafe", dir: "/go/src/unsafe"},
	},
	"tls": []pkg{
		pkg{importpath: "crypto/tls", dir: "/go/src/crypto/tls"},
	},
	"asn1": []pkg{
		pkg{importpath: "encoding/asn1", dir: "/go/src/encoding/asn1"},
	},
	"time": []pkg{
		pkg{importpath: "time", dir: "/go/src/time"},
	},
	"utf": []pkg{
		pkg{importpath: "lib9/utf", dir: "/go/src/lib9/utf"},
	},
	"math": []pkg{
		pkg{importpath: "math", dir: "/go/src/math"},
	},
	"mime": []pkg{
		pkg{importpath: "mime", dir: "/go/src/mime"},
	},
	"md5": []pkg{
		pkg{importpath: "crypto/md5", dir: "/go/src/crypto/md5"},
	},
	"crc32": []pkg{
		pkg{importpath: "hash/crc32", dir: "/go/src/hash/crc32"},
	},
	"html": []pkg{
		pkg{importpath: "html", dir: "/go/src/html"},
	},
	"scanner": []pkg{
		pkg{importpath: "go/scanner", dir: "/go/src/go/scanner"},
		pkg{importpath: "text/scanner", dir: "/go/src/text/scanner"},
	},
	"liblink": []pkg{
		pkg{importpath: "liblink", dir: "/go/src/liblink"},
	},
	"regexp": []pkg{
		pkg{importpath: "regexp", dir: "/go/src/regexp"},
	},
	"color": []pkg{
		pkg{importpath: "image/color", dir: "/go/src/image/color"},
	},
	"draw": []pkg{
		pkg{importpath: "image/draw", dir: "/go/src/image/draw"},
	},
	"fetch": []pkg{
		pkg{importpath: "cmd/pprof/internal/fetch", dir: "/go/src/cmd/pprof/internal/fetch"},
	},
	"aes": []pkg{
		pkg{importpath: "crypto/aes", dir: "/go/src/crypto/aes"},
	},
	"plan9obj": []pkg{
		pkg{importpath: "debug/plan9obj", dir: "/go/src/debug/plan9obj"},
	},
	"goobj": []pkg{
		pkg{importpath: "cmd/internal/goobj", dir: "/go/src/cmd/internal/goobj"},
	},
	"x86": []pkg{
		pkg{importpath: "cmd/internal/rsc.io/x86", dir: "/go/src/cmd/internal/rsc.io/x86"},
	},
	"bzip2": []pkg{
		pkg{importpath: "compress/bzip2", dir: "/go/src/compress/bzip2"},
	},
	"lzw": []pkg{
		pkg{importpath: "compress/lzw", dir: "/go/src/compress/lzw"},
	},
	"pem": []pkg{
		pkg{importpath: "encoding/pem", dir: "/go/src/encoding/pem"},
	},
	"tar": []pkg{
		pkg{importpath: "archive/tar", dir: "/go/src/archive/tar"},
	},
	"gc": []pkg{
		pkg{importpath: "cmd/gc", dir: "/go/src/cmd/gc"},
	},
	"go": []pkg{
		pkg{importpath: "cmd/go", dir: "/go/src/cmd/go"},
		pkg{importpath: "go", dir: "/go/src/go"},
	},
	"doc": []pkg{
		pkg{importpath: "go/doc", dir: "/go/src/go/doc"},
	},
	"adler32": []pkg{
		pkg{importpath: "hash/adler32", dir: "/go/src/hash/adler32"},
	},
	"log": []pkg{
		pkg{importpath: "log", dir: "/go/src/log"},
	},
	"database": []pkg{
		pkg{importpath: "database", dir: "/go/src/database"},
	},
	"symbolizer": []pkg{
		pkg{importpath: "cmd/pprof/internal/symbolizer", dir: "/go/src/cmd/pprof/internal/symbolizer"},
	},
	"csv": []pkg{
		pkg{importpath: "encoding/csv", dir: "/go/src/encoding/csv"},
	},
	"fcgi": []pkg{
		pkg{importpath: "net/http/fcgi", dir: "/go/src/net/http/fcgi"},
	},
	"encoding": []pkg{
		pkg{importpath: "encoding", dir: "/go/src/encoding"},
	},
	"fnv": []pkg{
		pkg{importpath: "hash/fnv", dir: "/go/src/hash/fnv"},
	},
	"libbio": []pkg{
		pkg{importpath: "libbio", dir: "/go/src/libbio"},
	},
	"iotest": []pkg{
		pkg{importpath: "testing/iotest", dir: "/go/src/testing/iotest"},
	},
	"gofmt": []pkg{
		pkg{importpath: "cmd/gofmt", dir: "/go/src/cmd/gofmt"},
	},
	"armasm": []pkg{
		pkg{importpath: "cmd/internal/rsc.io/arm/armasm", dir: "/go/src/cmd/internal/rsc.io/arm/armasm"},
	},
	"sql": []pkg{
		pkg{importpath: "database/sql", dir: "/go/src/database/sql"},
	},
	"svg": []pkg{
		pkg{importpath: "cmd/pprof/internal/svg", dir: "/go/src/cmd/pprof/internal/svg"},
	},
	"yacc": []pkg{
		pkg{importpath: "cmd/yacc", dir: "/go/src/cmd/yacc"},
	},
	"multipart": []pkg{
		pkg{importpath: "mime/multipart", dir: "/go/src/mime/multipart"},
	},
	"path": []pkg{
		pkg{importpath: "path", dir: "/go/src/path"},
	},
	"cgo": []pkg{
		pkg{importpath: "cmd/cgo", dir: "/go/src/cmd/cgo"},
		pkg{importpath: "runtime/cgo", dir: "/go/src/runtime/cgo"},
	},
	"plugin": []pkg{
		pkg{importpath: "cmd/pprof/internal/plugin", dir: "/go/src/cmd/pprof/internal/plugin"},
	},
	"report": []pkg{
		pkg{importpath: "cmd/pprof/internal/report", dir: "/go/src/cmd/pprof/internal/report"},
	},
}
