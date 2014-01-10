// starting import: "github.com/rocky/go-fish"
package repl

import (
	"bufio"
	"bytes"
	"code.google.com/p/go-columnize"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"github.com/0xfaded/eval"
	"github.com/mgutz/ansi"
	"go/ast"
	"go/parser"
	"go/scanner"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"regexp/syntax"
	"runtime"
	"runtime/pprof"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
	"text/tabwriter"
	"time"
	"unicode"
	"unicode/utf8"
)

type pkgType map[string] eval.Pkg

// EvalEnvironment adds to eval.Pkg those packages included
// with import "github.com/rocky/go-fish".

func EvalEnvironment(pkgs pkgType) {
	var consts map[string] reflect.Value
	var vars   map[string] reflect.Value
	var types  map[string] reflect.Type
	var funcs  map[string] reflect.Value

	consts = make(map[string] reflect.Value)
	consts["MaxScanTokenSize"] = reflect.ValueOf(bufio.MaxScanTokenSize)

	funcs = make(map[string] reflect.Value)
	funcs["NewReaderSize"] = reflect.ValueOf(bufio.NewReaderSize)
	funcs["NewReader"] = reflect.ValueOf(bufio.NewReader)
	funcs["NewWriterSize"] = reflect.ValueOf(bufio.NewWriterSize)
	funcs["NewWriter"] = reflect.ValueOf(bufio.NewWriter)
	funcs["NewReadWriter"] = reflect.ValueOf(bufio.NewReadWriter)
	funcs["NewScanner"] = reflect.ValueOf(bufio.NewScanner)
	funcs["ScanBytes"] = reflect.ValueOf(bufio.ScanBytes)
	funcs["ScanRunes"] = reflect.ValueOf(bufio.ScanRunes)
	funcs["ScanLines"] = reflect.ValueOf(bufio.ScanLines)
	funcs["ScanWords"] = reflect.ValueOf(bufio.ScanWords)

	types = make(map[string] reflect.Type)
	types["Reader"] = reflect.TypeOf(*new(bufio.Reader))
	types["Writer"] = reflect.TypeOf(*new(bufio.Writer))
	types["ReadWriter"] = reflect.TypeOf(*new(bufio.ReadWriter))
	types["Scanner"] = reflect.TypeOf(*new(bufio.Scanner))
	types["SplitFunc"] = reflect.TypeOf(*new(bufio.SplitFunc))

	vars = make(map[string] reflect.Value)
	vars["ErrInvalidUnreadByte"] = reflect.ValueOf(&bufio.ErrInvalidUnreadByte)
	vars["ErrInvalidUnreadRune"] = reflect.ValueOf(&bufio.ErrInvalidUnreadRune)
	vars["ErrBufferFull"] = reflect.ValueOf(&bufio.ErrBufferFull)
	vars["ErrNegativeCount"] = reflect.ValueOf(&bufio.ErrNegativeCount)
	vars["ErrTooLong"] = reflect.ValueOf(&bufio.ErrTooLong)
	vars["ErrNegativeAdvance"] = reflect.ValueOf(&bufio.ErrNegativeAdvance)
	vars["ErrAdvanceTooFar"] = reflect.ValueOf(&bufio.ErrAdvanceTooFar)
	pkgs["bufio"] = &eval.Env {
		Name: "bufio",
		Consts: consts,
		Funcs:  funcs,
		Types:  types,
		Vars:   vars,
		Pkgs:   pkgs,
		Path:   "bufio",
	}
	consts = make(map[string] reflect.Value)
	consts["MinRead"] = reflect.ValueOf(bytes.MinRead)

	funcs = make(map[string] reflect.Value)
	funcs["NewBuffer"] = reflect.ValueOf(bytes.NewBuffer)
	funcs["NewBufferString"] = reflect.ValueOf(bytes.NewBufferString)
	funcs["Count"] = reflect.ValueOf(bytes.Count)
	funcs["Contains"] = reflect.ValueOf(bytes.Contains)
	funcs["Index"] = reflect.ValueOf(bytes.Index)
	funcs["LastIndex"] = reflect.ValueOf(bytes.LastIndex)
	funcs["IndexRune"] = reflect.ValueOf(bytes.IndexRune)
	funcs["IndexAny"] = reflect.ValueOf(bytes.IndexAny)
	funcs["LastIndexAny"] = reflect.ValueOf(bytes.LastIndexAny)
	funcs["SplitN"] = reflect.ValueOf(bytes.SplitN)
	funcs["SplitAfterN"] = reflect.ValueOf(bytes.SplitAfterN)
	funcs["Split"] = reflect.ValueOf(bytes.Split)
	funcs["SplitAfter"] = reflect.ValueOf(bytes.SplitAfter)
	funcs["Fields"] = reflect.ValueOf(bytes.Fields)
	funcs["FieldsFunc"] = reflect.ValueOf(bytes.FieldsFunc)
	funcs["Join"] = reflect.ValueOf(bytes.Join)
	funcs["HasPrefix"] = reflect.ValueOf(bytes.HasPrefix)
	funcs["HasSuffix"] = reflect.ValueOf(bytes.HasSuffix)
	funcs["Map"] = reflect.ValueOf(bytes.Map)
	funcs["Repeat"] = reflect.ValueOf(bytes.Repeat)
	funcs["ToUpper"] = reflect.ValueOf(bytes.ToUpper)
	funcs["ToLower"] = reflect.ValueOf(bytes.ToLower)
	funcs["ToTitle"] = reflect.ValueOf(bytes.ToTitle)
	funcs["ToUpperSpecial"] = reflect.ValueOf(bytes.ToUpperSpecial)
	funcs["ToLowerSpecial"] = reflect.ValueOf(bytes.ToLowerSpecial)
	funcs["ToTitleSpecial"] = reflect.ValueOf(bytes.ToTitleSpecial)
	funcs["Title"] = reflect.ValueOf(bytes.Title)
	funcs["TrimLeftFunc"] = reflect.ValueOf(bytes.TrimLeftFunc)
	funcs["TrimRightFunc"] = reflect.ValueOf(bytes.TrimRightFunc)
	funcs["TrimFunc"] = reflect.ValueOf(bytes.TrimFunc)
	funcs["TrimPrefix"] = reflect.ValueOf(bytes.TrimPrefix)
	funcs["TrimSuffix"] = reflect.ValueOf(bytes.TrimSuffix)
	funcs["IndexFunc"] = reflect.ValueOf(bytes.IndexFunc)
	funcs["LastIndexFunc"] = reflect.ValueOf(bytes.LastIndexFunc)
	funcs["Trim"] = reflect.ValueOf(bytes.Trim)
	funcs["TrimLeft"] = reflect.ValueOf(bytes.TrimLeft)
	funcs["TrimRight"] = reflect.ValueOf(bytes.TrimRight)
	funcs["TrimSpace"] = reflect.ValueOf(bytes.TrimSpace)
	funcs["Runes"] = reflect.ValueOf(bytes.Runes)
	funcs["Replace"] = reflect.ValueOf(bytes.Replace)
	funcs["EqualFold"] = reflect.ValueOf(bytes.EqualFold)
	funcs["IndexByte"] = reflect.ValueOf(bytes.IndexByte)
	funcs["Equal"] = reflect.ValueOf(bytes.Equal)
	funcs["Compare"] = reflect.ValueOf(bytes.Compare)
	funcs["NewReader"] = reflect.ValueOf(bytes.NewReader)

	types = make(map[string] reflect.Type)
	types["Buffer"] = reflect.TypeOf(*new(bytes.Buffer))
	types["Reader"] = reflect.TypeOf(*new(bytes.Reader))

	vars = make(map[string] reflect.Value)
	vars["ErrTooLarge"] = reflect.ValueOf(&bytes.ErrTooLarge)
	pkgs["bytes"] = &eval.Env {
		Name: "bytes",
		Consts: consts,
		Funcs:  funcs,
		Types:  types,
		Vars:   vars,
		Pkgs:   pkgs,
		Path:   "bytes",
	}
	consts = make(map[string] reflect.Value)
	consts["VERSION"] = reflect.ValueOf(columnize.VERSION)

	funcs = make(map[string] reflect.Value)
	funcs["DefaultOptions"] = reflect.ValueOf(columnize.DefaultOptions)
	funcs["SetOptions"] = reflect.ValueOf(columnize.SetOptions)
	funcs["CellSize"] = reflect.ValueOf(columnize.CellSize)
	funcs["ToStringSliceFromIndexable"] = reflect.ValueOf(columnize.ToStringSliceFromIndexable)
	funcs["ToStringSlice"] = reflect.ValueOf(columnize.ToStringSlice)
	funcs["Columnize"] = reflect.ValueOf(columnize.Columnize)
	funcs["ColumnizeS"] = reflect.ValueOf(columnize.ColumnizeS)

	types = make(map[string] reflect.Type)
	types["Opts_t"] = reflect.TypeOf(*new(columnize.Opts_t))
	types["KeyValuePair_t"] = reflect.TypeOf(*new(columnize.KeyValuePair_t))

	vars = make(map[string] reflect.Value)
	pkgs["columnize"] = &eval.Env {
		Name: "columnize",
		Consts: consts,
		Funcs:  funcs,
		Types:  types,
		Vars:   vars,
		Pkgs:   pkgs,
		Path:   "code.google.com/p/go-columnize",
	}
	consts = make(map[string] reflect.Value)
	consts["MaxVarintLen16"] = reflect.ValueOf(binary.MaxVarintLen16)
	consts["MaxVarintLen32"] = reflect.ValueOf(binary.MaxVarintLen32)
	consts["MaxVarintLen64"] = reflect.ValueOf(binary.MaxVarintLen64)

	funcs = make(map[string] reflect.Value)
	funcs["Read"] = reflect.ValueOf(binary.Read)
	funcs["Write"] = reflect.ValueOf(binary.Write)
	funcs["Size"] = reflect.ValueOf(binary.Size)
	funcs["PutUvarint"] = reflect.ValueOf(binary.PutUvarint)
	funcs["Uvarint"] = reflect.ValueOf(binary.Uvarint)
	funcs["PutVarint"] = reflect.ValueOf(binary.PutVarint)
	funcs["Varint"] = reflect.ValueOf(binary.Varint)
	funcs["ReadUvarint"] = reflect.ValueOf(binary.ReadUvarint)
	funcs["ReadVarint"] = reflect.ValueOf(binary.ReadVarint)

	types = make(map[string] reflect.Type)
	types["ByteOrder"] = reflect.TypeOf(*new(binary.ByteOrder))

	vars = make(map[string] reflect.Value)
	vars["LittleEndian"] = reflect.ValueOf(&binary.LittleEndian)
	vars["BigEndian"] = reflect.ValueOf(&binary.BigEndian)
	pkgs["binary"] = &eval.Env {
		Name: "binary",
		Consts: consts,
		Funcs:  funcs,
		Types:  types,
		Vars:   vars,
		Pkgs:   pkgs,
		Path:   "encoding/binary",
	}
	consts = make(map[string] reflect.Value)

	funcs = make(map[string] reflect.Value)
	funcs["New"] = reflect.ValueOf(errors.New)

	types = make(map[string] reflect.Type)

	vars = make(map[string] reflect.Value)
	pkgs["errors"] = &eval.Env {
		Name: "errors",
		Consts: consts,
		Funcs:  funcs,
		Types:  types,
		Vars:   vars,
		Pkgs:   pkgs,
		Path:   "errors",
	}
	consts = make(map[string] reflect.Value)
	consts["ContinueOnError"] = reflect.ValueOf(flag.ContinueOnError)
	consts["ExitOnError"] = reflect.ValueOf(flag.ExitOnError)
	consts["PanicOnError"] = reflect.ValueOf(flag.PanicOnError)

	funcs = make(map[string] reflect.Value)
	funcs["VisitAll"] = reflect.ValueOf(flag.VisitAll)
	funcs["Visit"] = reflect.ValueOf(flag.Visit)
	funcs["Lookup"] = reflect.ValueOf(flag.Lookup)
	funcs["Set"] = reflect.ValueOf(flag.Set)
	funcs["PrintDefaults"] = reflect.ValueOf(flag.PrintDefaults)
	funcs["NFlag"] = reflect.ValueOf(flag.NFlag)
	funcs["Arg"] = reflect.ValueOf(flag.Arg)
	funcs["NArg"] = reflect.ValueOf(flag.NArg)
	funcs["Args"] = reflect.ValueOf(flag.Args)
	funcs["BoolVar"] = reflect.ValueOf(flag.BoolVar)
	funcs["Bool"] = reflect.ValueOf(flag.Bool)
	funcs["IntVar"] = reflect.ValueOf(flag.IntVar)
	funcs["Int"] = reflect.ValueOf(flag.Int)
	funcs["Int64Var"] = reflect.ValueOf(flag.Int64Var)
	funcs["Int64"] = reflect.ValueOf(flag.Int64)
	funcs["UintVar"] = reflect.ValueOf(flag.UintVar)
	funcs["Uint"] = reflect.ValueOf(flag.Uint)
	funcs["Uint64Var"] = reflect.ValueOf(flag.Uint64Var)
	funcs["Uint64"] = reflect.ValueOf(flag.Uint64)
	funcs["StringVar"] = reflect.ValueOf(flag.StringVar)
	funcs["String"] = reflect.ValueOf(flag.String)
	funcs["Float64Var"] = reflect.ValueOf(flag.Float64Var)
	funcs["Float64"] = reflect.ValueOf(flag.Float64)
	funcs["DurationVar"] = reflect.ValueOf(flag.DurationVar)
	funcs["Duration"] = reflect.ValueOf(flag.Duration)
	funcs["Var"] = reflect.ValueOf(flag.Var)
	funcs["Parse"] = reflect.ValueOf(flag.Parse)
	funcs["Parsed"] = reflect.ValueOf(flag.Parsed)
	funcs["NewFlagSet"] = reflect.ValueOf(flag.NewFlagSet)

	types = make(map[string] reflect.Type)
	types["Value"] = reflect.TypeOf(*new(flag.Value))
	types["Getter"] = reflect.TypeOf(*new(flag.Getter))
	types["ErrorHandling"] = reflect.TypeOf(*new(flag.ErrorHandling))
	types["FlagSet"] = reflect.TypeOf(*new(flag.FlagSet))
	types["Flag"] = reflect.TypeOf(*new(flag.Flag))

	vars = make(map[string] reflect.Value)
	vars["ErrHelp"] = reflect.ValueOf(&flag.ErrHelp)
	vars["Usage"] = reflect.ValueOf(&flag.Usage)
	vars["CommandLine"] = reflect.ValueOf(&flag.CommandLine)
	pkgs["flag"] = &eval.Env {
		Name: "flag",
		Consts: consts,
		Funcs:  funcs,
		Types:  types,
		Vars:   vars,
		Pkgs:   pkgs,
		Path:   "flag",
	}
	consts = make(map[string] reflect.Value)

	funcs = make(map[string] reflect.Value)
	funcs["Fprintf"] = reflect.ValueOf(fmt.Fprintf)
	funcs["Printf"] = reflect.ValueOf(fmt.Printf)
	funcs["Sprintf"] = reflect.ValueOf(fmt.Sprintf)
	funcs["Errorf"] = reflect.ValueOf(fmt.Errorf)
	funcs["Fprint"] = reflect.ValueOf(fmt.Fprint)
	funcs["Print"] = reflect.ValueOf(fmt.Print)
	funcs["Sprint"] = reflect.ValueOf(fmt.Sprint)
	funcs["Fprintln"] = reflect.ValueOf(fmt.Fprintln)
	funcs["Println"] = reflect.ValueOf(fmt.Println)
	funcs["Sprintln"] = reflect.ValueOf(fmt.Sprintln)
	funcs["Scan"] = reflect.ValueOf(fmt.Scan)
	funcs["Scanln"] = reflect.ValueOf(fmt.Scanln)
	funcs["Scanf"] = reflect.ValueOf(fmt.Scanf)
	funcs["Sscan"] = reflect.ValueOf(fmt.Sscan)
	funcs["Sscanln"] = reflect.ValueOf(fmt.Sscanln)
	funcs["Sscanf"] = reflect.ValueOf(fmt.Sscanf)
	funcs["Fscan"] = reflect.ValueOf(fmt.Fscan)
	funcs["Fscanln"] = reflect.ValueOf(fmt.Fscanln)
	funcs["Fscanf"] = reflect.ValueOf(fmt.Fscanf)

	types = make(map[string] reflect.Type)
	types["State"] = reflect.TypeOf(*new(fmt.State))
	types["Formatter"] = reflect.TypeOf(*new(fmt.Formatter))
	types["Stringer"] = reflect.TypeOf(*new(fmt.Stringer))
	types["GoStringer"] = reflect.TypeOf(*new(fmt.GoStringer))
	types["ScanState"] = reflect.TypeOf(*new(fmt.ScanState))
	types["Scanner"] = reflect.TypeOf(*new(fmt.Scanner))

	vars = make(map[string] reflect.Value)
	pkgs["fmt"] = &eval.Env {
		Name: "fmt",
		Consts: consts,
		Funcs:  funcs,
		Types:  types,
		Vars:   vars,
		Pkgs:   pkgs,
		Path:   "fmt",
	}
	consts = make(map[string] reflect.Value)

	funcs = make(map[string] reflect.Value)
	funcs["CheckExpr"] = reflect.ValueOf(eval.CheckExpr)
	funcs["NewConstInteger"] = reflect.ValueOf(eval.NewConstInteger)
	funcs["NewConstFloat"] = reflect.ValueOf(eval.NewConstFloat)
	funcs["NewConstImag"] = reflect.ValueOf(eval.NewConstImag)
	funcs["NewConstRune"] = reflect.ValueOf(eval.NewConstRune)
	funcs["NewConstInt64"] = reflect.ValueOf(eval.NewConstInt64)
	funcs["NewConstUint64"] = reflect.ValueOf(eval.NewConstUint64)
	funcs["NewConstFloat64"] = reflect.ValueOf(eval.NewConstFloat64)
	funcs["NewConstComplex128"] = reflect.ValueOf(eval.NewConstComplex128)
	funcs["EvalExpr"] = reflect.ValueOf(eval.EvalExpr)
	funcs["DerefValue"] = reflect.ValueOf(eval.DerefValue)
	funcs["EvalIdentExpr"] = reflect.ValueOf(eval.EvalIdentExpr)
	funcs["SetEvalIdentExprCallback"] = reflect.ValueOf(eval.SetEvalIdentExprCallback)
	funcs["GetEvalIdentExprCallback"] = reflect.ValueOf(eval.GetEvalIdentExprCallback)
	funcs["CannotIndex"] = reflect.ValueOf(eval.CannotIndex)
	funcs["InspectPtr"] = reflect.ValueOf(eval.InspectPtr)
	funcs["Inspect"] = reflect.ValueOf(eval.Inspect)
	funcs["EvalSelectorExpr"] = reflect.ValueOf(eval.EvalSelectorExpr)
	funcs["SetEvalSelectorExprCallback"] = reflect.ValueOf(eval.SetEvalSelectorExprCallback)
	funcs["GetEvalSelectorExprCallback"] = reflect.ValueOf(eval.GetEvalSelectorExprCallback)
	funcs["SetUserConversion"] = reflect.ValueOf(eval.SetUserConversion)
	funcs["GetUserConversion"] = reflect.ValueOf(eval.GetUserConversion)
	funcs["FormatErrorPos"] = reflect.ValueOf(eval.FormatErrorPos)

	types = make(map[string] reflect.Type)
	types["Expr"] = reflect.TypeOf(*new(eval.Expr))
	types["BadExpr"] = reflect.TypeOf(*new(eval.BadExpr))
	types["Ident"] = reflect.TypeOf(*new(eval.Ident))
	types["Ellipsis"] = reflect.TypeOf(*new(eval.Ellipsis))
	types["BasicLit"] = reflect.TypeOf(*new(eval.BasicLit))
	types["FuncLit"] = reflect.TypeOf(*new(eval.FuncLit))
	types["CompositeLit"] = reflect.TypeOf(*new(eval.CompositeLit))
	types["ParenExpr"] = reflect.TypeOf(*new(eval.ParenExpr))
	types["SelectorExpr"] = reflect.TypeOf(*new(eval.SelectorExpr))
	types["IndexExpr"] = reflect.TypeOf(*new(eval.IndexExpr))
	types["SliceExpr"] = reflect.TypeOf(*new(eval.SliceExpr))
	types["TypeAssertExpr"] = reflect.TypeOf(*new(eval.TypeAssertExpr))
	types["CallExpr"] = reflect.TypeOf(*new(eval.CallExpr))
	types["StarExpr"] = reflect.TypeOf(*new(eval.StarExpr))
	types["UnaryExpr"] = reflect.TypeOf(*new(eval.UnaryExpr))
	types["BinaryExpr"] = reflect.TypeOf(*new(eval.BinaryExpr))
	types["KeyValueExpr"] = reflect.TypeOf(*new(eval.KeyValueExpr))
	types["ArrayType"] = reflect.TypeOf(*new(eval.ArrayType))
	types["StructType"] = reflect.TypeOf(*new(eval.StructType))
	types["FuncType"] = reflect.TypeOf(*new(eval.FuncType))
	types["InterfaceType"] = reflect.TypeOf(*new(eval.InterfaceType))
	types["MapType"] = reflect.TypeOf(*new(eval.MapType))
	types["ChanType"] = reflect.TypeOf(*new(eval.ChanType))
	types["BigComplex"] = reflect.TypeOf(*new(eval.BigComplex))
	types["ConstNumber"] = reflect.TypeOf(*new(eval.ConstNumber))
	types["ConstType"] = reflect.TypeOf(*new(eval.ConstType))
	types["ConstIntType"] = reflect.TypeOf(*new(eval.ConstIntType))
	types["ConstRuneType"] = reflect.TypeOf(*new(eval.ConstRuneType))
	types["ConstFloatType"] = reflect.TypeOf(*new(eval.ConstFloatType))
	types["ConstComplexType"] = reflect.TypeOf(*new(eval.ConstComplexType))
	types["ConstStringType"] = reflect.TypeOf(*new(eval.ConstStringType))
	types["ConstNilType"] = reflect.TypeOf(*new(eval.ConstNilType))
	types["ConstBoolType"] = reflect.TypeOf(*new(eval.ConstBoolType))
	types["Ctx"] = reflect.TypeOf(*new(eval.Ctx))
	types["Pkg"] = reflect.TypeOf(*new(eval.Pkg))
	types["Env"] = reflect.TypeOf(*new(eval.Env))
	types["ErrBadBasicLit"] = reflect.TypeOf(*new(eval.ErrBadBasicLit))
	types["ErrInvalidOperand"] = reflect.TypeOf(*new(eval.ErrInvalidOperand))
	types["ErrInvalidIndirect"] = reflect.TypeOf(*new(eval.ErrInvalidIndirect))
	types["ErrMismatchedTypes"] = reflect.TypeOf(*new(eval.ErrMismatchedTypes))
	types["ErrInvalidOperands"] = reflect.TypeOf(*new(eval.ErrInvalidOperands))
	types["ErrBadFunArgument"] = reflect.TypeOf(*new(eval.ErrBadFunArgument))
	types["ErrBadComplexArguments"] = reflect.TypeOf(*new(eval.ErrBadComplexArguments))
	types["ErrBadBuiltinArgument"] = reflect.TypeOf(*new(eval.ErrBadBuiltinArgument))
	types["ErrWrongNumberOfArgsOld"] = reflect.TypeOf(*new(eval.ErrWrongNumberOfArgsOld))
	types["ErrWrongNumberOfArgs"] = reflect.TypeOf(*new(eval.ErrWrongNumberOfArgs))
	types["ErrMissingValue"] = reflect.TypeOf(*new(eval.ErrMissingValue))
	types["ErrMultiInSingleContext"] = reflect.TypeOf(*new(eval.ErrMultiInSingleContext))
	types["ErrArrayIndexOutOfBounds"] = reflect.TypeOf(*new(eval.ErrArrayIndexOutOfBounds))
	types["ErrInvalidIndexOperation"] = reflect.TypeOf(*new(eval.ErrInvalidIndexOperation))
	types["ErrInvalidIndex"] = reflect.TypeOf(*new(eval.ErrInvalidIndex))
	types["ErrDivideByZero"] = reflect.TypeOf(*new(eval.ErrDivideByZero))
	types["ErrInvalidBinaryOperation"] = reflect.TypeOf(*new(eval.ErrInvalidBinaryOperation))
	types["ErrInvalidUnaryOperation"] = reflect.TypeOf(*new(eval.ErrInvalidUnaryOperation))
	types["ErrBadConversion"] = reflect.TypeOf(*new(eval.ErrBadConversion))
	types["ErrBadConstConversion"] = reflect.TypeOf(*new(eval.ErrBadConstConversion))
	types["ErrTruncatedConstant"] = reflect.TypeOf(*new(eval.ErrTruncatedConstant))
	types["ErrOverflowedConstant"] = reflect.TypeOf(*new(eval.ErrOverflowedConstant))
	types["ErrUntypedNil"] = reflect.TypeOf(*new(eval.ErrUntypedNil))
	types["ErrorContext"] = reflect.TypeOf(*new(eval.ErrorContext))
	types["EvalIdentExprFunc"] = reflect.TypeOf(*new(eval.EvalIdentExprFunc))
	types["Rune"] = reflect.TypeOf(*new(eval.Rune))
	types["EvalSelectorExprFunc"] = reflect.TypeOf(*new(eval.EvalSelectorExprFunc))
	types["UntypedNil"] = reflect.TypeOf(*new(eval.UntypedNil))
	types["UserConvertFunc"] = reflect.TypeOf(*new(eval.UserConvertFunc))

	vars = make(map[string] reflect.Value)
	vars["ConstInt"] = reflect.ValueOf(&eval.ConstInt)
	vars["ConstRune"] = reflect.ValueOf(&eval.ConstRune)
	vars["ConstFloat"] = reflect.ValueOf(&eval.ConstFloat)
	vars["ConstComplex"] = reflect.ValueOf(&eval.ConstComplex)
	vars["ConstString"] = reflect.ValueOf(&eval.ConstString)
	vars["ConstNil"] = reflect.ValueOf(&eval.ConstNil)
	vars["ConstBool"] = reflect.ValueOf(&eval.ConstBool)
	vars["ErrArrayKey"] = reflect.ValueOf(&eval.ErrArrayKey)
	vars["RuneType"] = reflect.ValueOf(&eval.RuneType)
	pkgs["eval"] = &eval.Env {
		Name: "eval",
		Consts: consts,
		Funcs:  funcs,
		Types:  types,
		Vars:   vars,
		Pkgs:   pkgs,
		Path:   "github.com/0xfaded/eval",
	}
	consts = make(map[string] reflect.Value)
	consts["Reset"] = reflect.ValueOf(ansi.Reset)

	funcs = make(map[string] reflect.Value)
	funcs["ColorCode"] = reflect.ValueOf(ansi.ColorCode)
	funcs["Color"] = reflect.ValueOf(ansi.Color)
	funcs["ColorFunc"] = reflect.ValueOf(ansi.ColorFunc)
	funcs["DisableColors"] = reflect.ValueOf(ansi.DisableColors)

	types = make(map[string] reflect.Type)

	vars = make(map[string] reflect.Value)
	pkgs["ansi"] = &eval.Env {
		Name: "ansi",
		Consts: consts,
		Funcs:  funcs,
		Types:  types,
		Vars:   vars,
		Pkgs:   pkgs,
		Path:   "github.com/mgutz/ansi",
	}
	consts = make(map[string] reflect.Value)

	funcs = make(map[string] reflect.Value)
	funcs["AddAlias"] = reflect.ValueOf(AddAlias)
	funcs["AddToCategory"] = reflect.ValueOf(AddToCategory)
	funcs["LookupCmd"] = reflect.ValueOf(LookupCmd)
	funcs["Errmsg"] = reflect.ValueOf(Errmsg)
	funcs["MsgNoCr"] = reflect.ValueOf(MsgNoCr)
	funcs["Msg"] = reflect.ValueOf(Msg)
	funcs["Section"] = reflect.ValueOf(Section)
	funcs["PrintSorted"] = reflect.ValueOf(PrintSorted)
	funcs["HistoryFile"] = reflect.ValueOf(HistoryFile)
	funcs["SimpleReadLine"] = reflect.ValueOf(SimpleReadLine)
	funcs["SimpleInspect"] = reflect.ValueOf(SimpleInspect)
	funcs["MakeEvalEnv"] = reflect.ValueOf(MakeEvalEnv)
	funcs["REPL"] = reflect.ValueOf(REPL)
	funcs["EvalEnvironment"] = reflect.ValueOf(EvalEnvironment)
	funcs["ArgCountOK"] = reflect.ValueOf(ArgCountOK)
	funcs["GetInt"] = reflect.ValueOf(GetInt)
	funcs["GetUInt"] = reflect.ValueOf(GetUInt)

	types = make(map[string] reflect.Type)
	types["CmdFunc"] = reflect.TypeOf(*new(CmdFunc))
	types["CmdInfo"] = reflect.TypeOf(*new(CmdInfo))
	types["ReadLineFnType"] = reflect.TypeOf(*new(ReadLineFnType))
	types["InspectFnType"] = reflect.TypeOf(*new(InspectFnType))
	types["NumError"] = reflect.TypeOf(*new(NumError))

	vars = make(map[string] reflect.Value)
	vars["Cmds"] = reflect.ValueOf(&Cmds)
	vars["Aliases"] = reflect.ValueOf(&Aliases)
	vars["Categories"] = reflect.ValueOf(&Categories)
	vars["CmdLine"] = reflect.ValueOf(&CmdLine)
	vars["Highlight"] = reflect.ValueOf(&Highlight)
	vars["Maxwidth"] = reflect.ValueOf(&Maxwidth)
	vars["GOFISH_RESTART_CMD"] = reflect.ValueOf(&GOFISH_RESTART_CMD)
	vars["Input"] = reflect.ValueOf(&Input)
	vars["LeaveREPL"] = reflect.ValueOf(&LeaveREPL)
	vars["ExitCode"] = reflect.ValueOf(&ExitCode)
	vars["Env"] = reflect.ValueOf(&Env)
	pkgs["repl"] = &eval.Env {
		Name: "repl",
		Consts: consts,
		Funcs:  funcs,
		Types:  types,
		Vars:   vars,
		Pkgs:   pkgs,
		Path:   "github.com/rocky/go-fish",
	}
	consts = make(map[string] reflect.Value)
	consts["SEND"] = reflect.ValueOf(ast.SEND)
	consts["RECV"] = reflect.ValueOf(ast.RECV)
	consts["FilterFuncDuplicates"] = reflect.ValueOf(ast.FilterFuncDuplicates)
	consts["FilterUnassociatedComments"] = reflect.ValueOf(ast.FilterUnassociatedComments)
	consts["FilterImportDuplicates"] = reflect.ValueOf(ast.FilterImportDuplicates)
	consts["Bad"] = reflect.ValueOf(ast.Bad)
	consts["Pkg"] = reflect.ValueOf(ast.Pkg)
	consts["Con"] = reflect.ValueOf(ast.Con)
	consts["Typ"] = reflect.ValueOf(ast.Typ)
	consts["Var"] = reflect.ValueOf(ast.Var)
	consts["Fun"] = reflect.ValueOf(ast.Fun)
	consts["Lbl"] = reflect.ValueOf(ast.Lbl)

	funcs = make(map[string] reflect.Value)
	funcs["NewIdent"] = reflect.ValueOf(ast.NewIdent)
	funcs["IsExported"] = reflect.ValueOf(ast.IsExported)
	funcs["NewCommentMap"] = reflect.ValueOf(ast.NewCommentMap)
	funcs["FileExports"] = reflect.ValueOf(ast.FileExports)
	funcs["PackageExports"] = reflect.ValueOf(ast.PackageExports)
	funcs["FilterDecl"] = reflect.ValueOf(ast.FilterDecl)
	funcs["FilterFile"] = reflect.ValueOf(ast.FilterFile)
	funcs["FilterPackage"] = reflect.ValueOf(ast.FilterPackage)
	funcs["MergePackageFiles"] = reflect.ValueOf(ast.MergePackageFiles)
	funcs["SortImports"] = reflect.ValueOf(ast.SortImports)
	funcs["NotNilFilter"] = reflect.ValueOf(ast.NotNilFilter)
	funcs["Fprint"] = reflect.ValueOf(ast.Fprint)
	funcs["Print"] = reflect.ValueOf(ast.Print)
	funcs["NewPackage"] = reflect.ValueOf(ast.NewPackage)
	funcs["NewScope"] = reflect.ValueOf(ast.NewScope)
	funcs["NewObj"] = reflect.ValueOf(ast.NewObj)
	funcs["Walk"] = reflect.ValueOf(ast.Walk)
	funcs["Inspect"] = reflect.ValueOf(ast.Inspect)

	types = make(map[string] reflect.Type)
	types["Node"] = reflect.TypeOf(*new(ast.Node))
	types["Expr"] = reflect.TypeOf(*new(ast.Expr))
	types["Stmt"] = reflect.TypeOf(*new(ast.Stmt))
	types["Decl"] = reflect.TypeOf(*new(ast.Decl))
	types["Comment"] = reflect.TypeOf(*new(ast.Comment))
	types["CommentGroup"] = reflect.TypeOf(*new(ast.CommentGroup))
	types["Field"] = reflect.TypeOf(*new(ast.Field))
	types["FieldList"] = reflect.TypeOf(*new(ast.FieldList))
	types["BadExpr"] = reflect.TypeOf(*new(ast.BadExpr))
	types["Ident"] = reflect.TypeOf(*new(ast.Ident))
	types["Ellipsis"] = reflect.TypeOf(*new(ast.Ellipsis))
	types["BasicLit"] = reflect.TypeOf(*new(ast.BasicLit))
	types["FuncLit"] = reflect.TypeOf(*new(ast.FuncLit))
	types["CompositeLit"] = reflect.TypeOf(*new(ast.CompositeLit))
	types["ParenExpr"] = reflect.TypeOf(*new(ast.ParenExpr))
	types["SelectorExpr"] = reflect.TypeOf(*new(ast.SelectorExpr))
	types["IndexExpr"] = reflect.TypeOf(*new(ast.IndexExpr))
	types["SliceExpr"] = reflect.TypeOf(*new(ast.SliceExpr))
	types["TypeAssertExpr"] = reflect.TypeOf(*new(ast.TypeAssertExpr))
	types["CallExpr"] = reflect.TypeOf(*new(ast.CallExpr))
	types["StarExpr"] = reflect.TypeOf(*new(ast.StarExpr))
	types["UnaryExpr"] = reflect.TypeOf(*new(ast.UnaryExpr))
	types["BinaryExpr"] = reflect.TypeOf(*new(ast.BinaryExpr))
	types["KeyValueExpr"] = reflect.TypeOf(*new(ast.KeyValueExpr))
	types["ChanDir"] = reflect.TypeOf(*new(ast.ChanDir))
	types["ArrayType"] = reflect.TypeOf(*new(ast.ArrayType))
	types["StructType"] = reflect.TypeOf(*new(ast.StructType))
	types["FuncType"] = reflect.TypeOf(*new(ast.FuncType))
	types["InterfaceType"] = reflect.TypeOf(*new(ast.InterfaceType))
	types["MapType"] = reflect.TypeOf(*new(ast.MapType))
	types["ChanType"] = reflect.TypeOf(*new(ast.ChanType))
	types["BadStmt"] = reflect.TypeOf(*new(ast.BadStmt))
	types["DeclStmt"] = reflect.TypeOf(*new(ast.DeclStmt))
	types["EmptyStmt"] = reflect.TypeOf(*new(ast.EmptyStmt))
	types["LabeledStmt"] = reflect.TypeOf(*new(ast.LabeledStmt))
	types["ExprStmt"] = reflect.TypeOf(*new(ast.ExprStmt))
	types["SendStmt"] = reflect.TypeOf(*new(ast.SendStmt))
	types["IncDecStmt"] = reflect.TypeOf(*new(ast.IncDecStmt))
	types["AssignStmt"] = reflect.TypeOf(*new(ast.AssignStmt))
	types["GoStmt"] = reflect.TypeOf(*new(ast.GoStmt))
	types["DeferStmt"] = reflect.TypeOf(*new(ast.DeferStmt))
	types["ReturnStmt"] = reflect.TypeOf(*new(ast.ReturnStmt))
	types["BranchStmt"] = reflect.TypeOf(*new(ast.BranchStmt))
	types["BlockStmt"] = reflect.TypeOf(*new(ast.BlockStmt))
	types["IfStmt"] = reflect.TypeOf(*new(ast.IfStmt))
	types["CaseClause"] = reflect.TypeOf(*new(ast.CaseClause))
	types["SwitchStmt"] = reflect.TypeOf(*new(ast.SwitchStmt))
	types["TypeSwitchStmt"] = reflect.TypeOf(*new(ast.TypeSwitchStmt))
	types["CommClause"] = reflect.TypeOf(*new(ast.CommClause))
	types["SelectStmt"] = reflect.TypeOf(*new(ast.SelectStmt))
	types["ForStmt"] = reflect.TypeOf(*new(ast.ForStmt))
	types["RangeStmt"] = reflect.TypeOf(*new(ast.RangeStmt))
	types["Spec"] = reflect.TypeOf(*new(ast.Spec))
	types["ImportSpec"] = reflect.TypeOf(*new(ast.ImportSpec))
	types["ValueSpec"] = reflect.TypeOf(*new(ast.ValueSpec))
	types["TypeSpec"] = reflect.TypeOf(*new(ast.TypeSpec))
	types["BadDecl"] = reflect.TypeOf(*new(ast.BadDecl))
	types["GenDecl"] = reflect.TypeOf(*new(ast.GenDecl))
	types["FuncDecl"] = reflect.TypeOf(*new(ast.FuncDecl))
	types["File"] = reflect.TypeOf(*new(ast.File))
	types["Package"] = reflect.TypeOf(*new(ast.Package))
	types["CommentMap"] = reflect.TypeOf(*new(ast.CommentMap))
	types["Filter"] = reflect.TypeOf(*new(ast.Filter))
	types["MergeMode"] = reflect.TypeOf(*new(ast.MergeMode))
	types["FieldFilter"] = reflect.TypeOf(*new(ast.FieldFilter))
	types["Importer"] = reflect.TypeOf(*new(ast.Importer))
	types["Scope"] = reflect.TypeOf(*new(ast.Scope))
	types["Object"] = reflect.TypeOf(*new(ast.Object))
	types["ObjKind"] = reflect.TypeOf(*new(ast.ObjKind))
	types["Visitor"] = reflect.TypeOf(*new(ast.Visitor))

	vars = make(map[string] reflect.Value)
	pkgs["ast"] = &eval.Env {
		Name: "ast",
		Consts: consts,
		Funcs:  funcs,
		Types:  types,
		Vars:   vars,
		Pkgs:   pkgs,
		Path:   "go/ast",
	}
	consts = make(map[string] reflect.Value)
	consts["PackageClauseOnly"] = reflect.ValueOf(parser.PackageClauseOnly)
	consts["ImportsOnly"] = reflect.ValueOf(parser.ImportsOnly)
	consts["ParseComments"] = reflect.ValueOf(parser.ParseComments)
	consts["Trace"] = reflect.ValueOf(parser.Trace)
	consts["DeclarationErrors"] = reflect.ValueOf(parser.DeclarationErrors)
	consts["SpuriousErrors"] = reflect.ValueOf(parser.SpuriousErrors)
	consts["AllErrors"] = reflect.ValueOf(parser.AllErrors)

	funcs = make(map[string] reflect.Value)
	funcs["ParseFile"] = reflect.ValueOf(parser.ParseFile)
	funcs["ParseDir"] = reflect.ValueOf(parser.ParseDir)
	funcs["ParseExpr"] = reflect.ValueOf(parser.ParseExpr)

	types = make(map[string] reflect.Type)
	types["Mode"] = reflect.TypeOf(*new(parser.Mode))

	vars = make(map[string] reflect.Value)
	pkgs["parser"] = &eval.Env {
		Name: "parser",
		Consts: consts,
		Funcs:  funcs,
		Types:  types,
		Vars:   vars,
		Pkgs:   pkgs,
		Path:   "go/parser",
	}
	consts = make(map[string] reflect.Value)
	consts["ScanComments"] = reflect.ValueOf(scanner.ScanComments)

	funcs = make(map[string] reflect.Value)
	funcs["PrintError"] = reflect.ValueOf(scanner.PrintError)

	types = make(map[string] reflect.Type)
	types["Error"] = reflect.TypeOf(*new(scanner.Error))
	types["ErrorList"] = reflect.TypeOf(*new(scanner.ErrorList))
	types["ErrorHandler"] = reflect.TypeOf(*new(scanner.ErrorHandler))
	types["Scanner"] = reflect.TypeOf(*new(scanner.Scanner))
	types["Mode"] = reflect.TypeOf(*new(scanner.Mode))

	vars = make(map[string] reflect.Value)
	pkgs["scanner"] = &eval.Env {
		Name: "scanner",
		Consts: consts,
		Funcs:  funcs,
		Types:  types,
		Vars:   vars,
		Pkgs:   pkgs,
		Path:   "go/scanner",
	}
	consts = make(map[string] reflect.Value)
	consts["NoPos"] = reflect.ValueOf(token.NoPos)
	consts["ILLEGAL"] = reflect.ValueOf(token.ILLEGAL)
	consts["EOF"] = reflect.ValueOf(token.EOF)
	consts["COMMENT"] = reflect.ValueOf(token.COMMENT)
	consts["IDENT"] = reflect.ValueOf(token.IDENT)
	consts["INT"] = reflect.ValueOf(token.INT)
	consts["FLOAT"] = reflect.ValueOf(token.FLOAT)
	consts["IMAG"] = reflect.ValueOf(token.IMAG)
	consts["CHAR"] = reflect.ValueOf(token.CHAR)
	consts["STRING"] = reflect.ValueOf(token.STRING)
	consts["ADD"] = reflect.ValueOf(token.ADD)
	consts["SUB"] = reflect.ValueOf(token.SUB)
	consts["MUL"] = reflect.ValueOf(token.MUL)
	consts["QUO"] = reflect.ValueOf(token.QUO)
	consts["REM"] = reflect.ValueOf(token.REM)
	consts["AND"] = reflect.ValueOf(token.AND)
	consts["OR"] = reflect.ValueOf(token.OR)
	consts["XOR"] = reflect.ValueOf(token.XOR)
	consts["SHL"] = reflect.ValueOf(token.SHL)
	consts["SHR"] = reflect.ValueOf(token.SHR)
	consts["AND_NOT"] = reflect.ValueOf(token.AND_NOT)
	consts["ADD_ASSIGN"] = reflect.ValueOf(token.ADD_ASSIGN)
	consts["SUB_ASSIGN"] = reflect.ValueOf(token.SUB_ASSIGN)
	consts["MUL_ASSIGN"] = reflect.ValueOf(token.MUL_ASSIGN)
	consts["QUO_ASSIGN"] = reflect.ValueOf(token.QUO_ASSIGN)
	consts["REM_ASSIGN"] = reflect.ValueOf(token.REM_ASSIGN)
	consts["AND_ASSIGN"] = reflect.ValueOf(token.AND_ASSIGN)
	consts["OR_ASSIGN"] = reflect.ValueOf(token.OR_ASSIGN)
	consts["XOR_ASSIGN"] = reflect.ValueOf(token.XOR_ASSIGN)
	consts["SHL_ASSIGN"] = reflect.ValueOf(token.SHL_ASSIGN)
	consts["SHR_ASSIGN"] = reflect.ValueOf(token.SHR_ASSIGN)
	consts["AND_NOT_ASSIGN"] = reflect.ValueOf(token.AND_NOT_ASSIGN)
	consts["LAND"] = reflect.ValueOf(token.LAND)
	consts["LOR"] = reflect.ValueOf(token.LOR)
	consts["ARROW"] = reflect.ValueOf(token.ARROW)
	consts["INC"] = reflect.ValueOf(token.INC)
	consts["DEC"] = reflect.ValueOf(token.DEC)
	consts["EQL"] = reflect.ValueOf(token.EQL)
	consts["LSS"] = reflect.ValueOf(token.LSS)
	consts["GTR"] = reflect.ValueOf(token.GTR)
	consts["ASSIGN"] = reflect.ValueOf(token.ASSIGN)
	consts["NOT"] = reflect.ValueOf(token.NOT)
	consts["NEQ"] = reflect.ValueOf(token.NEQ)
	consts["LEQ"] = reflect.ValueOf(token.LEQ)
	consts["GEQ"] = reflect.ValueOf(token.GEQ)
	consts["DEFINE"] = reflect.ValueOf(token.DEFINE)
	consts["ELLIPSIS"] = reflect.ValueOf(token.ELLIPSIS)
	consts["LPAREN"] = reflect.ValueOf(token.LPAREN)
	consts["LBRACK"] = reflect.ValueOf(token.LBRACK)
	consts["LBRACE"] = reflect.ValueOf(token.LBRACE)
	consts["COMMA"] = reflect.ValueOf(token.COMMA)
	consts["PERIOD"] = reflect.ValueOf(token.PERIOD)
	consts["RPAREN"] = reflect.ValueOf(token.RPAREN)
	consts["RBRACK"] = reflect.ValueOf(token.RBRACK)
	consts["RBRACE"] = reflect.ValueOf(token.RBRACE)
	consts["SEMICOLON"] = reflect.ValueOf(token.SEMICOLON)
	consts["COLON"] = reflect.ValueOf(token.COLON)
	consts["BREAK"] = reflect.ValueOf(token.BREAK)
	consts["CASE"] = reflect.ValueOf(token.CASE)
	consts["CHAN"] = reflect.ValueOf(token.CHAN)
	consts["CONST"] = reflect.ValueOf(token.CONST)
	consts["CONTINUE"] = reflect.ValueOf(token.CONTINUE)
	consts["DEFAULT"] = reflect.ValueOf(token.DEFAULT)
	consts["DEFER"] = reflect.ValueOf(token.DEFER)
	consts["ELSE"] = reflect.ValueOf(token.ELSE)
	consts["FALLTHROUGH"] = reflect.ValueOf(token.FALLTHROUGH)
	consts["FOR"] = reflect.ValueOf(token.FOR)
	consts["FUNC"] = reflect.ValueOf(token.FUNC)
	consts["GO"] = reflect.ValueOf(token.GO)
	consts["GOTO"] = reflect.ValueOf(token.GOTO)
	consts["IF"] = reflect.ValueOf(token.IF)
	consts["IMPORT"] = reflect.ValueOf(token.IMPORT)
	consts["INTERFACE"] = reflect.ValueOf(token.INTERFACE)
	consts["MAP"] = reflect.ValueOf(token.MAP)
	consts["PACKAGE"] = reflect.ValueOf(token.PACKAGE)
	consts["RANGE"] = reflect.ValueOf(token.RANGE)
	consts["RETURN"] = reflect.ValueOf(token.RETURN)
	consts["SELECT"] = reflect.ValueOf(token.SELECT)
	consts["STRUCT"] = reflect.ValueOf(token.STRUCT)
	consts["SWITCH"] = reflect.ValueOf(token.SWITCH)
	consts["TYPE"] = reflect.ValueOf(token.TYPE)
	consts["VAR"] = reflect.ValueOf(token.VAR)
	consts["LowestPrec"] = reflect.ValueOf(token.LowestPrec)
	consts["UnaryPrec"] = reflect.ValueOf(token.UnaryPrec)
	consts["HighestPrec"] = reflect.ValueOf(token.HighestPrec)

	funcs = make(map[string] reflect.Value)
	funcs["NewFileSet"] = reflect.ValueOf(token.NewFileSet)
	funcs["Lookup"] = reflect.ValueOf(token.Lookup)

	types = make(map[string] reflect.Type)
	types["Position"] = reflect.TypeOf(*new(token.Position))
	types["Pos"] = reflect.TypeOf(*new(token.Pos))
	types["File"] = reflect.TypeOf(*new(token.File))
	types["FileSet"] = reflect.TypeOf(*new(token.FileSet))
	types["Token"] = reflect.TypeOf(*new(token.Token))

	vars = make(map[string] reflect.Value)
	pkgs["token"] = &eval.Env {
		Name: "token",
		Consts: consts,
		Funcs:  funcs,
		Types:  types,
		Vars:   vars,
		Pkgs:   pkgs,
		Path:   "go/token",
	}
	consts = make(map[string] reflect.Value)

	funcs = make(map[string] reflect.Value)
	funcs["WriteString"] = reflect.ValueOf(io.WriteString)
	funcs["ReadAtLeast"] = reflect.ValueOf(io.ReadAtLeast)
	funcs["ReadFull"] = reflect.ValueOf(io.ReadFull)
	funcs["CopyN"] = reflect.ValueOf(io.CopyN)
	funcs["Copy"] = reflect.ValueOf(io.Copy)
	funcs["LimitReader"] = reflect.ValueOf(io.LimitReader)
	funcs["NewSectionReader"] = reflect.ValueOf(io.NewSectionReader)
	funcs["TeeReader"] = reflect.ValueOf(io.TeeReader)
	funcs["MultiReader"] = reflect.ValueOf(io.MultiReader)
	funcs["MultiWriter"] = reflect.ValueOf(io.MultiWriter)
	funcs["Pipe"] = reflect.ValueOf(io.Pipe)

	types = make(map[string] reflect.Type)
	types["Reader"] = reflect.TypeOf(*new(io.Reader))
	types["Writer"] = reflect.TypeOf(*new(io.Writer))
	types["Closer"] = reflect.TypeOf(*new(io.Closer))
	types["Seeker"] = reflect.TypeOf(*new(io.Seeker))
	types["ReadWriter"] = reflect.TypeOf(*new(io.ReadWriter))
	types["ReadCloser"] = reflect.TypeOf(*new(io.ReadCloser))
	types["WriteCloser"] = reflect.TypeOf(*new(io.WriteCloser))
	types["ReadWriteCloser"] = reflect.TypeOf(*new(io.ReadWriteCloser))
	types["ReadSeeker"] = reflect.TypeOf(*new(io.ReadSeeker))
	types["WriteSeeker"] = reflect.TypeOf(*new(io.WriteSeeker))
	types["ReadWriteSeeker"] = reflect.TypeOf(*new(io.ReadWriteSeeker))
	types["ReaderFrom"] = reflect.TypeOf(*new(io.ReaderFrom))
	types["WriterTo"] = reflect.TypeOf(*new(io.WriterTo))
	types["ReaderAt"] = reflect.TypeOf(*new(io.ReaderAt))
	types["WriterAt"] = reflect.TypeOf(*new(io.WriterAt))
	types["ByteReader"] = reflect.TypeOf(*new(io.ByteReader))
	types["ByteScanner"] = reflect.TypeOf(*new(io.ByteScanner))
	types["ByteWriter"] = reflect.TypeOf(*new(io.ByteWriter))
	types["RuneReader"] = reflect.TypeOf(*new(io.RuneReader))
	types["RuneScanner"] = reflect.TypeOf(*new(io.RuneScanner))
	types["LimitedReader"] = reflect.TypeOf(*new(io.LimitedReader))
	types["SectionReader"] = reflect.TypeOf(*new(io.SectionReader))
	types["PipeReader"] = reflect.TypeOf(*new(io.PipeReader))
	types["PipeWriter"] = reflect.TypeOf(*new(io.PipeWriter))

	vars = make(map[string] reflect.Value)
	vars["ErrShortWrite"] = reflect.ValueOf(&io.ErrShortWrite)
	vars["ErrShortBuffer"] = reflect.ValueOf(&io.ErrShortBuffer)
	vars["EOF"] = reflect.ValueOf(&io.EOF)
	vars["ErrUnexpectedEOF"] = reflect.ValueOf(&io.ErrUnexpectedEOF)
	vars["ErrNoProgress"] = reflect.ValueOf(&io.ErrNoProgress)
	vars["ErrClosedPipe"] = reflect.ValueOf(&io.ErrClosedPipe)
	pkgs["io"] = &eval.Env {
		Name: "io",
		Consts: consts,
		Funcs:  funcs,
		Types:  types,
		Vars:   vars,
		Pkgs:   pkgs,
		Path:   "io",
	}
	consts = make(map[string] reflect.Value)

	funcs = make(map[string] reflect.Value)
	funcs["ReadAll"] = reflect.ValueOf(ioutil.ReadAll)
	funcs["ReadFile"] = reflect.ValueOf(ioutil.ReadFile)
	funcs["WriteFile"] = reflect.ValueOf(ioutil.WriteFile)
	funcs["ReadDir"] = reflect.ValueOf(ioutil.ReadDir)
	funcs["NopCloser"] = reflect.ValueOf(ioutil.NopCloser)
	funcs["TempFile"] = reflect.ValueOf(ioutil.TempFile)
	funcs["TempDir"] = reflect.ValueOf(ioutil.TempDir)

	types = make(map[string] reflect.Type)

	vars = make(map[string] reflect.Value)
	vars["Discard"] = reflect.ValueOf(&ioutil.Discard)
	pkgs["ioutil"] = &eval.Env {
		Name: "ioutil",
		Consts: consts,
		Funcs:  funcs,
		Types:  types,
		Vars:   vars,
		Pkgs:   pkgs,
		Path:   "io/ioutil",
	}
	consts = make(map[string] reflect.Value)
	consts["Ldate"] = reflect.ValueOf(log.Ldate)
	consts["Ltime"] = reflect.ValueOf(log.Ltime)
	consts["Lmicroseconds"] = reflect.ValueOf(log.Lmicroseconds)
	consts["Llongfile"] = reflect.ValueOf(log.Llongfile)
	consts["Lshortfile"] = reflect.ValueOf(log.Lshortfile)
	consts["LstdFlags"] = reflect.ValueOf(log.LstdFlags)

	funcs = make(map[string] reflect.Value)
	funcs["New"] = reflect.ValueOf(log.New)
	funcs["SetOutput"] = reflect.ValueOf(log.SetOutput)
	funcs["Flags"] = reflect.ValueOf(log.Flags)
	funcs["SetFlags"] = reflect.ValueOf(log.SetFlags)
	funcs["Prefix"] = reflect.ValueOf(log.Prefix)
	funcs["SetPrefix"] = reflect.ValueOf(log.SetPrefix)
	funcs["Print"] = reflect.ValueOf(log.Print)
	funcs["Printf"] = reflect.ValueOf(log.Printf)
	funcs["Println"] = reflect.ValueOf(log.Println)
	funcs["Fatal"] = reflect.ValueOf(log.Fatal)
	funcs["Fatalf"] = reflect.ValueOf(log.Fatalf)
	funcs["Fatalln"] = reflect.ValueOf(log.Fatalln)
	funcs["Panic"] = reflect.ValueOf(log.Panic)
	funcs["Panicf"] = reflect.ValueOf(log.Panicf)
	funcs["Panicln"] = reflect.ValueOf(log.Panicln)

	types = make(map[string] reflect.Type)
	types["Logger"] = reflect.TypeOf(*new(log.Logger))

	vars = make(map[string] reflect.Value)
	pkgs["log"] = &eval.Env {
		Name: "log",
		Consts: consts,
		Funcs:  funcs,
		Types:  types,
		Vars:   vars,
		Pkgs:   pkgs,
		Path:   "log",
	}
	consts = make(map[string] reflect.Value)
	consts["E"] = reflect.ValueOf(math.E)
	consts["Pi"] = reflect.ValueOf(math.Pi)
	consts["Phi"] = reflect.ValueOf(math.Phi)
	consts["Sqrt2"] = reflect.ValueOf(math.Sqrt2)
	consts["SqrtE"] = reflect.ValueOf(math.SqrtE)
	consts["SqrtPi"] = reflect.ValueOf(math.SqrtPi)
	consts["SqrtPhi"] = reflect.ValueOf(math.SqrtPhi)
	consts["Ln2"] = reflect.ValueOf(math.Ln2)
	consts["Log2E"] = reflect.ValueOf(math.Log2E)
	consts["Ln10"] = reflect.ValueOf(math.Ln10)
	consts["Log10E"] = reflect.ValueOf(math.Log10E)
	consts["MaxFloat32"] = reflect.ValueOf(math.MaxFloat32)
	consts["SmallestNonzeroFloat32"] = reflect.ValueOf(math.SmallestNonzeroFloat32)
	consts["MaxFloat64"] = reflect.ValueOf(math.MaxFloat64)
	consts["SmallestNonzeroFloat64"] = reflect.ValueOf(math.SmallestNonzeroFloat64)
	consts["MaxInt8"] = reflect.ValueOf(math.MaxInt8)
	consts["MinInt8"] = reflect.ValueOf(math.MinInt8)
	consts["MaxInt16"] = reflect.ValueOf(math.MaxInt16)
	consts["MinInt16"] = reflect.ValueOf(math.MinInt16)
	consts["MaxInt32"] = reflect.ValueOf(math.MaxInt32)
	consts["MinInt32"] = reflect.ValueOf(math.MinInt32)
	consts["MaxInt64"] = reflect.ValueOf(int64(math.MaxInt64))
	consts["MinInt64"] = reflect.ValueOf(int64(math.MinInt64))
	consts["MaxUint8"] = reflect.ValueOf(uint8(math.MaxUint8))
	consts["MaxUint16"] = reflect.ValueOf(uint16(math.MaxUint16))
	consts["MaxUint32"] = reflect.ValueOf(uint32(math.MaxUint32))
	consts["MaxUint64"] = reflect.ValueOf(uint64(math.MaxUint64))

	funcs = make(map[string] reflect.Value)
	funcs["Abs"] = reflect.ValueOf(math.Abs)
	funcs["Acosh"] = reflect.ValueOf(math.Acosh)
	funcs["Asin"] = reflect.ValueOf(math.Asin)
	funcs["Acos"] = reflect.ValueOf(math.Acos)
	funcs["Asinh"] = reflect.ValueOf(math.Asinh)
	funcs["Atan"] = reflect.ValueOf(math.Atan)
	funcs["Atan2"] = reflect.ValueOf(math.Atan2)
	funcs["Atanh"] = reflect.ValueOf(math.Atanh)
	funcs["Inf"] = reflect.ValueOf(math.Inf)
	funcs["NaN"] = reflect.ValueOf(math.NaN)
	funcs["IsNaN"] = reflect.ValueOf(math.IsNaN)
	funcs["IsInf"] = reflect.ValueOf(math.IsInf)
	funcs["Cbrt"] = reflect.ValueOf(math.Cbrt)
	funcs["Copysign"] = reflect.ValueOf(math.Copysign)
	funcs["Dim"] = reflect.ValueOf(math.Dim)
	funcs["Max"] = reflect.ValueOf(math.Max)
	funcs["Min"] = reflect.ValueOf(math.Min)
	funcs["Erf"] = reflect.ValueOf(math.Erf)
	funcs["Erfc"] = reflect.ValueOf(math.Erfc)
	funcs["Exp"] = reflect.ValueOf(math.Exp)
	funcs["Exp2"] = reflect.ValueOf(math.Exp2)
	funcs["Expm1"] = reflect.ValueOf(math.Expm1)
	funcs["Floor"] = reflect.ValueOf(math.Floor)
	funcs["Ceil"] = reflect.ValueOf(math.Ceil)
	funcs["Trunc"] = reflect.ValueOf(math.Trunc)
	funcs["Frexp"] = reflect.ValueOf(math.Frexp)
	funcs["Gamma"] = reflect.ValueOf(math.Gamma)
	funcs["Hypot"] = reflect.ValueOf(math.Hypot)
	funcs["J0"] = reflect.ValueOf(math.J0)
	funcs["Y0"] = reflect.ValueOf(math.Y0)
	funcs["J1"] = reflect.ValueOf(math.J1)
	funcs["Y1"] = reflect.ValueOf(math.Y1)
	funcs["Jn"] = reflect.ValueOf(math.Jn)
	funcs["Yn"] = reflect.ValueOf(math.Yn)
	funcs["Ldexp"] = reflect.ValueOf(math.Ldexp)
	funcs["Lgamma"] = reflect.ValueOf(math.Lgamma)
	funcs["Log"] = reflect.ValueOf(math.Log)
	funcs["Log10"] = reflect.ValueOf(math.Log10)
	funcs["Log2"] = reflect.ValueOf(math.Log2)
	funcs["Log1p"] = reflect.ValueOf(math.Log1p)
	funcs["Logb"] = reflect.ValueOf(math.Logb)
	funcs["Ilogb"] = reflect.ValueOf(math.Ilogb)
	funcs["Mod"] = reflect.ValueOf(math.Mod)
	funcs["Modf"] = reflect.ValueOf(math.Modf)
	funcs["Nextafter"] = reflect.ValueOf(math.Nextafter)
	funcs["Pow"] = reflect.ValueOf(math.Pow)
	funcs["Pow10"] = reflect.ValueOf(math.Pow10)
	funcs["Remainder"] = reflect.ValueOf(math.Remainder)
	funcs["Signbit"] = reflect.ValueOf(math.Signbit)
	funcs["Cos"] = reflect.ValueOf(math.Cos)
	funcs["Sin"] = reflect.ValueOf(math.Sin)
	funcs["Sincos"] = reflect.ValueOf(math.Sincos)
	funcs["Sinh"] = reflect.ValueOf(math.Sinh)
	funcs["Cosh"] = reflect.ValueOf(math.Cosh)
	funcs["Sqrt"] = reflect.ValueOf(math.Sqrt)
	funcs["Tan"] = reflect.ValueOf(math.Tan)
	funcs["Tanh"] = reflect.ValueOf(math.Tanh)
	funcs["Float32bits"] = reflect.ValueOf(math.Float32bits)
	funcs["Float32frombits"] = reflect.ValueOf(math.Float32frombits)
	funcs["Float64bits"] = reflect.ValueOf(math.Float64bits)
	funcs["Float64frombits"] = reflect.ValueOf(math.Float64frombits)

	types = make(map[string] reflect.Type)

	vars = make(map[string] reflect.Value)
	pkgs["math"] = &eval.Env {
		Name: "math",
		Consts: consts,
		Funcs:  funcs,
		Types:  types,
		Vars:   vars,
		Pkgs:   pkgs,
		Path:   "math",
	}
	consts = make(map[string] reflect.Value)
	consts["MaxBase"] = reflect.ValueOf(big.MaxBase)

	funcs = make(map[string] reflect.Value)
	funcs["NewInt"] = reflect.ValueOf(big.NewInt)
	funcs["NewRat"] = reflect.ValueOf(big.NewRat)

	types = make(map[string] reflect.Type)
	types["Word"] = reflect.TypeOf(*new(big.Word))
	types["Int"] = reflect.TypeOf(*new(big.Int))
	types["Rat"] = reflect.TypeOf(*new(big.Rat))

	vars = make(map[string] reflect.Value)
	pkgs["big"] = &eval.Env {
		Name: "big",
		Consts: consts,
		Funcs:  funcs,
		Types:  types,
		Vars:   vars,
		Pkgs:   pkgs,
		Path:   "math/big",
	}
	consts = make(map[string] reflect.Value)

	funcs = make(map[string] reflect.Value)
	funcs["NewSource"] = reflect.ValueOf(rand.NewSource)
	funcs["New"] = reflect.ValueOf(rand.New)
	funcs["Seed"] = reflect.ValueOf(rand.Seed)
	funcs["Int63"] = reflect.ValueOf(rand.Int63)
	funcs["Uint32"] = reflect.ValueOf(rand.Uint32)
	funcs["Int31"] = reflect.ValueOf(rand.Int31)
	funcs["Int"] = reflect.ValueOf(rand.Int)
	funcs["Int63n"] = reflect.ValueOf(rand.Int63n)
	funcs["Int31n"] = reflect.ValueOf(rand.Int31n)
	funcs["Intn"] = reflect.ValueOf(rand.Intn)
	funcs["Float64"] = reflect.ValueOf(rand.Float64)
	funcs["Float32"] = reflect.ValueOf(rand.Float32)
	funcs["Perm"] = reflect.ValueOf(rand.Perm)
	funcs["NormFloat64"] = reflect.ValueOf(rand.NormFloat64)
	funcs["ExpFloat64"] = reflect.ValueOf(rand.ExpFloat64)
	funcs["NewZipf"] = reflect.ValueOf(rand.NewZipf)

	types = make(map[string] reflect.Type)
	types["Source"] = reflect.TypeOf(*new(rand.Source))
	types["Rand"] = reflect.TypeOf(*new(rand.Rand))
	types["Zipf"] = reflect.TypeOf(*new(rand.Zipf))

	vars = make(map[string] reflect.Value)
	pkgs["rand"] = &eval.Env {
		Name: "rand",
		Consts: consts,
		Funcs:  funcs,
		Types:  types,
		Vars:   vars,
		Pkgs:   pkgs,
		Path:   "math/rand",
	}
	consts = make(map[string] reflect.Value)
	consts["O_RDONLY"] = reflect.ValueOf(os.O_RDONLY)
	consts["O_WRONLY"] = reflect.ValueOf(os.O_WRONLY)
	consts["O_RDWR"] = reflect.ValueOf(os.O_RDWR)
	consts["O_APPEND"] = reflect.ValueOf(os.O_APPEND)
	consts["O_CREATE"] = reflect.ValueOf(os.O_CREATE)
	consts["O_EXCL"] = reflect.ValueOf(os.O_EXCL)
	consts["O_SYNC"] = reflect.ValueOf(os.O_SYNC)
	consts["O_TRUNC"] = reflect.ValueOf(os.O_TRUNC)
	consts["SEEK_SET"] = reflect.ValueOf(os.SEEK_SET)
	consts["SEEK_CUR"] = reflect.ValueOf(os.SEEK_CUR)
	consts["SEEK_END"] = reflect.ValueOf(os.SEEK_END)
	consts["DevNull"] = reflect.ValueOf(os.DevNull)
	consts["PathSeparator"] = reflect.ValueOf(os.PathSeparator)
	consts["PathListSeparator"] = reflect.ValueOf(os.PathListSeparator)
	consts["ModeDir"] = reflect.ValueOf(os.ModeDir)
	consts["ModeAppend"] = reflect.ValueOf(os.ModeAppend)
	consts["ModeExclusive"] = reflect.ValueOf(os.ModeExclusive)
	consts["ModeTemporary"] = reflect.ValueOf(os.ModeTemporary)
	consts["ModeSymlink"] = reflect.ValueOf(os.ModeSymlink)
	consts["ModeDevice"] = reflect.ValueOf(os.ModeDevice)
	consts["ModeNamedPipe"] = reflect.ValueOf(os.ModeNamedPipe)
	consts["ModeSocket"] = reflect.ValueOf(os.ModeSocket)
	consts["ModeSetuid"] = reflect.ValueOf(os.ModeSetuid)
	consts["ModeSetgid"] = reflect.ValueOf(os.ModeSetgid)
	consts["ModeCharDevice"] = reflect.ValueOf(os.ModeCharDevice)
	consts["ModeSticky"] = reflect.ValueOf(os.ModeSticky)
	consts["ModeType"] = reflect.ValueOf(os.ModeType)
	consts["ModePerm"] = reflect.ValueOf(os.ModePerm)

	funcs = make(map[string] reflect.Value)
	funcs["FindProcess"] = reflect.ValueOf(os.FindProcess)
	funcs["StartProcess"] = reflect.ValueOf(os.StartProcess)
	funcs["Hostname"] = reflect.ValueOf(os.Hostname)
	funcs["Expand"] = reflect.ValueOf(os.Expand)
	funcs["ExpandEnv"] = reflect.ValueOf(os.ExpandEnv)
	funcs["Getenv"] = reflect.ValueOf(os.Getenv)
	funcs["Setenv"] = reflect.ValueOf(os.Setenv)
	funcs["Clearenv"] = reflect.ValueOf(os.Clearenv)
	funcs["Environ"] = reflect.ValueOf(os.Environ)
	funcs["NewSyscallError"] = reflect.ValueOf(os.NewSyscallError)
	funcs["IsExist"] = reflect.ValueOf(os.IsExist)
	funcs["IsNotExist"] = reflect.ValueOf(os.IsNotExist)
	funcs["IsPermission"] = reflect.ValueOf(os.IsPermission)
	funcs["Getpid"] = reflect.ValueOf(os.Getpid)
	funcs["Getppid"] = reflect.ValueOf(os.Getppid)
	funcs["Mkdir"] = reflect.ValueOf(os.Mkdir)
	funcs["Chdir"] = reflect.ValueOf(os.Chdir)
	funcs["Open"] = reflect.ValueOf(os.Open)
	funcs["Create"] = reflect.ValueOf(os.Create)
	funcs["Link"] = reflect.ValueOf(os.Link)
	funcs["Symlink"] = reflect.ValueOf(os.Symlink)
	funcs["Readlink"] = reflect.ValueOf(os.Readlink)
	funcs["Rename"] = reflect.ValueOf(os.Rename)
	funcs["Chmod"] = reflect.ValueOf(os.Chmod)
	funcs["Chown"] = reflect.ValueOf(os.Chown)
	funcs["Lchown"] = reflect.ValueOf(os.Lchown)
	funcs["Chtimes"] = reflect.ValueOf(os.Chtimes)
	funcs["NewFile"] = reflect.ValueOf(os.NewFile)
	funcs["OpenFile"] = reflect.ValueOf(os.OpenFile)
	funcs["Stat"] = reflect.ValueOf(os.Stat)
	funcs["Lstat"] = reflect.ValueOf(os.Lstat)
	funcs["Truncate"] = reflect.ValueOf(os.Truncate)
	funcs["Remove"] = reflect.ValueOf(os.Remove)
	funcs["TempDir"] = reflect.ValueOf(os.TempDir)
	funcs["Getwd"] = reflect.ValueOf(os.Getwd)
	funcs["MkdirAll"] = reflect.ValueOf(os.MkdirAll)
	funcs["RemoveAll"] = reflect.ValueOf(os.RemoveAll)
	funcs["IsPathSeparator"] = reflect.ValueOf(os.IsPathSeparator)
	funcs["Pipe"] = reflect.ValueOf(os.Pipe)
	funcs["Getuid"] = reflect.ValueOf(os.Getuid)
	funcs["Geteuid"] = reflect.ValueOf(os.Geteuid)
	funcs["Getgid"] = reflect.ValueOf(os.Getgid)
	funcs["Getegid"] = reflect.ValueOf(os.Getegid)
	funcs["Getgroups"] = reflect.ValueOf(os.Getgroups)
	funcs["Exit"] = reflect.ValueOf(os.Exit)
	funcs["Getpagesize"] = reflect.ValueOf(os.Getpagesize)
	funcs["SameFile"] = reflect.ValueOf(os.SameFile)

	types = make(map[string] reflect.Type)
	types["PathError"] = reflect.TypeOf(*new(os.PathError))
	types["SyscallError"] = reflect.TypeOf(*new(os.SyscallError))
	types["Process"] = reflect.TypeOf(*new(os.Process))
	types["ProcAttr"] = reflect.TypeOf(*new(os.ProcAttr))
	types["Signal"] = reflect.TypeOf(*new(os.Signal))
	types["ProcessState"] = reflect.TypeOf(*new(os.ProcessState))
	types["LinkError"] = reflect.TypeOf(*new(os.LinkError))
	types["File"] = reflect.TypeOf(*new(os.File))
	types["FileInfo"] = reflect.TypeOf(*new(os.FileInfo))
	types["FileMode"] = reflect.TypeOf(*new(os.FileMode))

	vars = make(map[string] reflect.Value)
	vars["ErrInvalid"] = reflect.ValueOf(&os.ErrInvalid)
	vars["ErrPermission"] = reflect.ValueOf(&os.ErrPermission)
	vars["ErrExist"] = reflect.ValueOf(&os.ErrExist)
	vars["ErrNotExist"] = reflect.ValueOf(&os.ErrNotExist)
	vars["Interrupt"] = reflect.ValueOf(&os.Interrupt)
	vars["Kill"] = reflect.ValueOf(&os.Kill)
	vars["Stdin"] = reflect.ValueOf(&os.Stdin)
	vars["Stdout"] = reflect.ValueOf(&os.Stdout)
	vars["Stderr"] = reflect.ValueOf(&os.Stderr)
	vars["Args"] = reflect.ValueOf(&os.Args)
	pkgs["os"] = &eval.Env {
		Name: "os",
		Consts: consts,
		Funcs:  funcs,
		Types:  types,
		Vars:   vars,
		Pkgs:   pkgs,
		Path:   "os",
	}
	consts = make(map[string] reflect.Value)

	funcs = make(map[string] reflect.Value)
	funcs["Command"] = reflect.ValueOf(exec.Command)
	funcs["LookPath"] = reflect.ValueOf(exec.LookPath)

	types = make(map[string] reflect.Type)
	types["Error"] = reflect.TypeOf(*new(exec.Error))
	types["Cmd"] = reflect.TypeOf(*new(exec.Cmd))
	types["ExitError"] = reflect.TypeOf(*new(exec.ExitError))

	vars = make(map[string] reflect.Value)
	vars["ErrNotFound"] = reflect.ValueOf(&exec.ErrNotFound)
	pkgs["exec"] = &eval.Env {
		Name: "exec",
		Consts: consts,
		Funcs:  funcs,
		Types:  types,
		Vars:   vars,
		Pkgs:   pkgs,
		Path:   "os/exec",
	}
	consts = make(map[string] reflect.Value)
	consts["Separator"] = reflect.ValueOf(filepath.Separator)
	consts["ListSeparator"] = reflect.ValueOf(filepath.ListSeparator)

	funcs = make(map[string] reflect.Value)
	funcs["Match"] = reflect.ValueOf(filepath.Match)
	funcs["Glob"] = reflect.ValueOf(filepath.Glob)
	funcs["Clean"] = reflect.ValueOf(filepath.Clean)
	funcs["ToSlash"] = reflect.ValueOf(filepath.ToSlash)
	funcs["FromSlash"] = reflect.ValueOf(filepath.FromSlash)
	funcs["SplitList"] = reflect.ValueOf(filepath.SplitList)
	funcs["Split"] = reflect.ValueOf(filepath.Split)
	funcs["Join"] = reflect.ValueOf(filepath.Join)
	funcs["Ext"] = reflect.ValueOf(filepath.Ext)
	funcs["EvalSymlinks"] = reflect.ValueOf(filepath.EvalSymlinks)
	funcs["Abs"] = reflect.ValueOf(filepath.Abs)
	funcs["Rel"] = reflect.ValueOf(filepath.Rel)
	funcs["Walk"] = reflect.ValueOf(filepath.Walk)
	funcs["Base"] = reflect.ValueOf(filepath.Base)
	funcs["Dir"] = reflect.ValueOf(filepath.Dir)
	funcs["VolumeName"] = reflect.ValueOf(filepath.VolumeName)
	funcs["IsAbs"] = reflect.ValueOf(filepath.IsAbs)
	funcs["HasPrefix"] = reflect.ValueOf(filepath.HasPrefix)

	types = make(map[string] reflect.Type)
	types["WalkFunc"] = reflect.TypeOf(*new(filepath.WalkFunc))

	vars = make(map[string] reflect.Value)
	vars["ErrBadPattern"] = reflect.ValueOf(&filepath.ErrBadPattern)
	vars["SkipDir"] = reflect.ValueOf(&filepath.SkipDir)
	pkgs["filepath"] = &eval.Env {
		Name: "filepath",
		Consts: consts,
		Funcs:  funcs,
		Types:  types,
		Vars:   vars,
		Pkgs:   pkgs,
		Path:   "path/filepath",
	}
	consts = make(map[string] reflect.Value)
	consts["Invalid"] = reflect.ValueOf(reflect.Invalid)
	consts["Bool"] = reflect.ValueOf(reflect.Bool)
	consts["Int"] = reflect.ValueOf(reflect.Int)
	consts["Int8"] = reflect.ValueOf(reflect.Int8)
	consts["Int16"] = reflect.ValueOf(reflect.Int16)
	consts["Int32"] = reflect.ValueOf(reflect.Int32)
	consts["Int64"] = reflect.ValueOf(reflect.Int64)
	consts["Uint"] = reflect.ValueOf(reflect.Uint)
	consts["Uint8"] = reflect.ValueOf(reflect.Uint8)
	consts["Uint16"] = reflect.ValueOf(reflect.Uint16)
	consts["Uint32"] = reflect.ValueOf(reflect.Uint32)
	consts["Uint64"] = reflect.ValueOf(reflect.Uint64)
	consts["Uintptr"] = reflect.ValueOf(reflect.Uintptr)
	consts["Float32"] = reflect.ValueOf(reflect.Float32)
	consts["Float64"] = reflect.ValueOf(reflect.Float64)
	consts["Complex64"] = reflect.ValueOf(reflect.Complex64)
	consts["Complex128"] = reflect.ValueOf(reflect.Complex128)
	consts["Array"] = reflect.ValueOf(reflect.Array)
	consts["Chan"] = reflect.ValueOf(reflect.Chan)
	consts["Func"] = reflect.ValueOf(reflect.Func)
	consts["Interface"] = reflect.ValueOf(reflect.Interface)
	consts["Map"] = reflect.ValueOf(reflect.Map)
	consts["Ptr"] = reflect.ValueOf(reflect.Ptr)
	consts["Slice"] = reflect.ValueOf(reflect.Slice)
	consts["String"] = reflect.ValueOf(reflect.String)
	consts["Struct"] = reflect.ValueOf(reflect.Struct)
	consts["UnsafePointer"] = reflect.ValueOf(reflect.UnsafePointer)
	consts["RecvDir"] = reflect.ValueOf(reflect.RecvDir)
	consts["SendDir"] = reflect.ValueOf(reflect.SendDir)
	consts["BothDir"] = reflect.ValueOf(reflect.BothDir)
	consts["SelectSend"] = reflect.ValueOf(reflect.SelectSend)
	consts["SelectRecv"] = reflect.ValueOf(reflect.SelectRecv)
	consts["SelectDefault"] = reflect.ValueOf(reflect.SelectDefault)

	funcs = make(map[string] reflect.Value)
	funcs["DeepEqual"] = reflect.ValueOf(reflect.DeepEqual)
	funcs["MakeFunc"] = reflect.ValueOf(reflect.MakeFunc)
	funcs["TypeOf"] = reflect.ValueOf(reflect.TypeOf)
	funcs["PtrTo"] = reflect.ValueOf(reflect.PtrTo)
	funcs["ChanOf"] = reflect.ValueOf(reflect.ChanOf)
	funcs["MapOf"] = reflect.ValueOf(reflect.MapOf)
	funcs["SliceOf"] = reflect.ValueOf(reflect.SliceOf)
	funcs["Append"] = reflect.ValueOf(reflect.Append)
	funcs["AppendSlice"] = reflect.ValueOf(reflect.AppendSlice)
	funcs["Copy"] = reflect.ValueOf(reflect.Copy)
	funcs["Select"] = reflect.ValueOf(reflect.Select)
	funcs["MakeSlice"] = reflect.ValueOf(reflect.MakeSlice)
	funcs["MakeChan"] = reflect.ValueOf(reflect.MakeChan)
	funcs["MakeMap"] = reflect.ValueOf(reflect.MakeMap)
	funcs["Indirect"] = reflect.ValueOf(reflect.Indirect)
	funcs["ValueOf"] = reflect.ValueOf(reflect.ValueOf)
	funcs["Zero"] = reflect.ValueOf(reflect.Zero)
	funcs["New"] = reflect.ValueOf(reflect.New)
	funcs["NewAt"] = reflect.ValueOf(reflect.NewAt)

	types = make(map[string] reflect.Type)
	types["Type"] = reflect.TypeOf(*new(reflect.Type))
	types["Kind"] = reflect.TypeOf(*new(reflect.Kind))
	types["ChanDir"] = reflect.TypeOf(*new(reflect.ChanDir))
	types["Method"] = reflect.TypeOf(*new(reflect.Method))
	types["StructField"] = reflect.TypeOf(*new(reflect.StructField))
	types["StructTag"] = reflect.TypeOf(*new(reflect.StructTag))
	types["Value"] = reflect.TypeOf(*new(reflect.Value))
	types["ValueError"] = reflect.TypeOf(*new(reflect.ValueError))
	types["StringHeader"] = reflect.TypeOf(*new(reflect.StringHeader))
	types["SliceHeader"] = reflect.TypeOf(*new(reflect.SliceHeader))
	types["SelectDir"] = reflect.TypeOf(*new(reflect.SelectDir))
	types["SelectCase"] = reflect.TypeOf(*new(reflect.SelectCase))

	vars = make(map[string] reflect.Value)
	pkgs["reflect"] = &eval.Env {
		Name: "reflect",
		Consts: consts,
		Funcs:  funcs,
		Types:  types,
		Vars:   vars,
		Pkgs:   pkgs,
		Path:   "reflect",
	}
	consts = make(map[string] reflect.Value)

	funcs = make(map[string] reflect.Value)
	funcs["Compile"] = reflect.ValueOf(regexp.Compile)
	funcs["CompilePOSIX"] = reflect.ValueOf(regexp.CompilePOSIX)
	funcs["MustCompile"] = reflect.ValueOf(regexp.MustCompile)
	funcs["MustCompilePOSIX"] = reflect.ValueOf(regexp.MustCompilePOSIX)
	funcs["MatchReader"] = reflect.ValueOf(regexp.MatchReader)
	funcs["MatchString"] = reflect.ValueOf(regexp.MatchString)
	funcs["Match"] = reflect.ValueOf(regexp.Match)
	funcs["QuoteMeta"] = reflect.ValueOf(regexp.QuoteMeta)

	types = make(map[string] reflect.Type)
	types["Regexp"] = reflect.TypeOf(*new(regexp.Regexp))

	vars = make(map[string] reflect.Value)
	pkgs["regexp"] = &eval.Env {
		Name: "regexp",
		Consts: consts,
		Funcs:  funcs,
		Types:  types,
		Vars:   vars,
		Pkgs:   pkgs,
		Path:   "regexp",
	}
	consts = make(map[string] reflect.Value)
	consts["ErrInternalError"] = reflect.ValueOf(syntax.ErrInternalError)
	consts["ErrInvalidCharClass"] = reflect.ValueOf(syntax.ErrInvalidCharClass)
	consts["ErrInvalidCharRange"] = reflect.ValueOf(syntax.ErrInvalidCharRange)
	consts["ErrInvalidEscape"] = reflect.ValueOf(syntax.ErrInvalidEscape)
	consts["ErrInvalidNamedCapture"] = reflect.ValueOf(syntax.ErrInvalidNamedCapture)
	consts["ErrInvalidPerlOp"] = reflect.ValueOf(syntax.ErrInvalidPerlOp)
	consts["ErrInvalidRepeatOp"] = reflect.ValueOf(syntax.ErrInvalidRepeatOp)
	consts["ErrInvalidRepeatSize"] = reflect.ValueOf(syntax.ErrInvalidRepeatSize)
	consts["ErrInvalidUTF8"] = reflect.ValueOf(syntax.ErrInvalidUTF8)
	consts["ErrMissingBracket"] = reflect.ValueOf(syntax.ErrMissingBracket)
	consts["ErrMissingParen"] = reflect.ValueOf(syntax.ErrMissingParen)
	consts["ErrMissingRepeatArgument"] = reflect.ValueOf(syntax.ErrMissingRepeatArgument)
	consts["ErrTrailingBackslash"] = reflect.ValueOf(syntax.ErrTrailingBackslash)
	consts["ErrUnexpectedParen"] = reflect.ValueOf(syntax.ErrUnexpectedParen)
	consts["FoldCase"] = reflect.ValueOf(syntax.FoldCase)
	consts["Literal"] = reflect.ValueOf(syntax.Literal)
	consts["ClassNL"] = reflect.ValueOf(syntax.ClassNL)
	consts["DotNL"] = reflect.ValueOf(syntax.DotNL)
	consts["OneLine"] = reflect.ValueOf(syntax.OneLine)
	consts["NonGreedy"] = reflect.ValueOf(syntax.NonGreedy)
	consts["PerlX"] = reflect.ValueOf(syntax.PerlX)
	consts["UnicodeGroups"] = reflect.ValueOf(syntax.UnicodeGroups)
	consts["WasDollar"] = reflect.ValueOf(syntax.WasDollar)
	consts["Simple"] = reflect.ValueOf(syntax.Simple)
	consts["MatchNL"] = reflect.ValueOf(syntax.MatchNL)
	consts["Perl"] = reflect.ValueOf(syntax.Perl)
	consts["POSIX"] = reflect.ValueOf(syntax.POSIX)
	consts["InstAlt"] = reflect.ValueOf(syntax.InstAlt)
	consts["InstAltMatch"] = reflect.ValueOf(syntax.InstAltMatch)
	consts["InstCapture"] = reflect.ValueOf(syntax.InstCapture)
	consts["InstEmptyWidth"] = reflect.ValueOf(syntax.InstEmptyWidth)
	consts["InstMatch"] = reflect.ValueOf(syntax.InstMatch)
	consts["InstFail"] = reflect.ValueOf(syntax.InstFail)
	consts["InstNop"] = reflect.ValueOf(syntax.InstNop)
	consts["InstRune"] = reflect.ValueOf(syntax.InstRune)
	consts["InstRune1"] = reflect.ValueOf(syntax.InstRune1)
	consts["InstRuneAny"] = reflect.ValueOf(syntax.InstRuneAny)
	consts["InstRuneAnyNotNL"] = reflect.ValueOf(syntax.InstRuneAnyNotNL)
	consts["EmptyBeginLine"] = reflect.ValueOf(syntax.EmptyBeginLine)
	consts["EmptyEndLine"] = reflect.ValueOf(syntax.EmptyEndLine)
	consts["EmptyBeginText"] = reflect.ValueOf(syntax.EmptyBeginText)
	consts["EmptyEndText"] = reflect.ValueOf(syntax.EmptyEndText)
	consts["EmptyWordBoundary"] = reflect.ValueOf(syntax.EmptyWordBoundary)
	consts["EmptyNoWordBoundary"] = reflect.ValueOf(syntax.EmptyNoWordBoundary)
	consts["OpNoMatch"] = reflect.ValueOf(syntax.OpNoMatch)
	consts["OpEmptyMatch"] = reflect.ValueOf(syntax.OpEmptyMatch)
	consts["OpLiteral"] = reflect.ValueOf(syntax.OpLiteral)
	consts["OpCharClass"] = reflect.ValueOf(syntax.OpCharClass)
	consts["OpAnyCharNotNL"] = reflect.ValueOf(syntax.OpAnyCharNotNL)
	consts["OpAnyChar"] = reflect.ValueOf(syntax.OpAnyChar)
	consts["OpBeginLine"] = reflect.ValueOf(syntax.OpBeginLine)
	consts["OpEndLine"] = reflect.ValueOf(syntax.OpEndLine)
	consts["OpBeginText"] = reflect.ValueOf(syntax.OpBeginText)
	consts["OpEndText"] = reflect.ValueOf(syntax.OpEndText)
	consts["OpWordBoundary"] = reflect.ValueOf(syntax.OpWordBoundary)
	consts["OpNoWordBoundary"] = reflect.ValueOf(syntax.OpNoWordBoundary)
	consts["OpCapture"] = reflect.ValueOf(syntax.OpCapture)
	consts["OpStar"] = reflect.ValueOf(syntax.OpStar)
	consts["OpPlus"] = reflect.ValueOf(syntax.OpPlus)
	consts["OpQuest"] = reflect.ValueOf(syntax.OpQuest)
	consts["OpRepeat"] = reflect.ValueOf(syntax.OpRepeat)
	consts["OpConcat"] = reflect.ValueOf(syntax.OpConcat)
	consts["OpAlternate"] = reflect.ValueOf(syntax.OpAlternate)

	funcs = make(map[string] reflect.Value)
	funcs["Compile"] = reflect.ValueOf(syntax.Compile)
	funcs["Parse"] = reflect.ValueOf(syntax.Parse)
	funcs["EmptyOpContext"] = reflect.ValueOf(syntax.EmptyOpContext)
	funcs["IsWordChar"] = reflect.ValueOf(syntax.IsWordChar)

	types = make(map[string] reflect.Type)
	types["Error"] = reflect.TypeOf(*new(syntax.Error))
	types["ErrorCode"] = reflect.TypeOf(*new(syntax.ErrorCode))
	types["Flags"] = reflect.TypeOf(*new(syntax.Flags))
	types["Prog"] = reflect.TypeOf(*new(syntax.Prog))
	types["InstOp"] = reflect.TypeOf(*new(syntax.InstOp))
	types["EmptyOp"] = reflect.TypeOf(*new(syntax.EmptyOp))
	types["Inst"] = reflect.TypeOf(*new(syntax.Inst))
	types["Regexp"] = reflect.TypeOf(*new(syntax.Regexp))
	types["Op"] = reflect.TypeOf(*new(syntax.Op))

	vars = make(map[string] reflect.Value)
	pkgs["syntax"] = &eval.Env {
		Name: "syntax",
		Consts: consts,
		Funcs:  funcs,
		Types:  types,
		Vars:   vars,
		Pkgs:   pkgs,
		Path:   "regexp/syntax",
	}
	consts = make(map[string] reflect.Value)
	consts["Compiler"] = reflect.ValueOf(runtime.Compiler)
	consts["GOOS"] = reflect.ValueOf(runtime.GOOS)
	consts["GOARCH"] = reflect.ValueOf(runtime.GOARCH)

	funcs = make(map[string] reflect.Value)
	funcs["Breakpoint"] = reflect.ValueOf(runtime.Breakpoint)
	funcs["LockOSThread"] = reflect.ValueOf(runtime.LockOSThread)
	funcs["UnlockOSThread"] = reflect.ValueOf(runtime.UnlockOSThread)
	funcs["GOMAXPROCS"] = reflect.ValueOf(runtime.GOMAXPROCS)
	funcs["NumCPU"] = reflect.ValueOf(runtime.NumCPU)
	funcs["NumCgoCall"] = reflect.ValueOf(runtime.NumCgoCall)
	funcs["NumGoroutine"] = reflect.ValueOf(runtime.NumGoroutine)
	funcs["MemProfile"] = reflect.ValueOf(runtime.MemProfile)
	funcs["ThreadCreateProfile"] = reflect.ValueOf(runtime.ThreadCreateProfile)
	funcs["GoroutineProfile"] = reflect.ValueOf(runtime.GoroutineProfile)
	funcs["CPUProfile"] = reflect.ValueOf(runtime.CPUProfile)
	funcs["SetCPUProfileRate"] = reflect.ValueOf(runtime.SetCPUProfileRate)
	funcs["SetBlockProfileRate"] = reflect.ValueOf(runtime.SetBlockProfileRate)
	funcs["BlockProfile"] = reflect.ValueOf(runtime.BlockProfile)
	funcs["Stack"] = reflect.ValueOf(runtime.Stack)
	funcs["Gosched"] = reflect.ValueOf(runtime.Gosched)
	funcs["Goexit"] = reflect.ValueOf(runtime.Goexit)
	funcs["Caller"] = reflect.ValueOf(runtime.Caller)
	funcs["Callers"] = reflect.ValueOf(runtime.Callers)
	funcs["FuncForPC"] = reflect.ValueOf(runtime.FuncForPC)
	funcs["SetFinalizer"] = reflect.ValueOf(runtime.SetFinalizer)
	funcs["GOROOT"] = reflect.ValueOf(runtime.GOROOT)
	funcs["Version"] = reflect.ValueOf(runtime.Version)
	funcs["ReadMemStats"] = reflect.ValueOf(runtime.ReadMemStats)
	funcs["GC"] = reflect.ValueOf(runtime.GC)

	types = make(map[string] reflect.Type)
	types["MemProfileRecord"] = reflect.TypeOf(*new(runtime.MemProfileRecord))
	types["StackRecord"] = reflect.TypeOf(*new(runtime.StackRecord))
	types["BlockProfileRecord"] = reflect.TypeOf(*new(runtime.BlockProfileRecord))
	types["Error"] = reflect.TypeOf(*new(runtime.Error))
	types["TypeAssertionError"] = reflect.TypeOf(*new(runtime.TypeAssertionError))
	types["Func"] = reflect.TypeOf(*new(runtime.Func))
	types["MemStats"] = reflect.TypeOf(*new(runtime.MemStats))

	vars = make(map[string] reflect.Value)
	vars["MemProfileRate"] = reflect.ValueOf(&runtime.MemProfileRate)
	pkgs["runtime"] = &eval.Env {
		Name: "runtime",
		Consts: consts,
		Funcs:  funcs,
		Types:  types,
		Vars:   vars,
		Pkgs:   pkgs,
		Path:   "runtime",
	}
	consts = make(map[string] reflect.Value)

	funcs = make(map[string] reflect.Value)
	funcs["NewProfile"] = reflect.ValueOf(pprof.NewProfile)
	funcs["Lookup"] = reflect.ValueOf(pprof.Lookup)
	funcs["Profiles"] = reflect.ValueOf(pprof.Profiles)
	funcs["WriteHeapProfile"] = reflect.ValueOf(pprof.WriteHeapProfile)
	funcs["StartCPUProfile"] = reflect.ValueOf(pprof.StartCPUProfile)
	funcs["StopCPUProfile"] = reflect.ValueOf(pprof.StopCPUProfile)

	types = make(map[string] reflect.Type)
	types["Profile"] = reflect.TypeOf(*new(pprof.Profile))

	vars = make(map[string] reflect.Value)
	pkgs["pprof"] = &eval.Env {
		Name: "pprof",
		Consts: consts,
		Funcs:  funcs,
		Types:  types,
		Vars:   vars,
		Pkgs:   pkgs,
		Path:   "runtime/pprof",
	}
	consts = make(map[string] reflect.Value)

	funcs = make(map[string] reflect.Value)
	funcs["Search"] = reflect.ValueOf(sort.Search)
	funcs["SearchInts"] = reflect.ValueOf(sort.SearchInts)
	funcs["SearchFloat64s"] = reflect.ValueOf(sort.SearchFloat64s)
	funcs["SearchStrings"] = reflect.ValueOf(sort.SearchStrings)
	funcs["Sort"] = reflect.ValueOf(sort.Sort)
	funcs["Reverse"] = reflect.ValueOf(sort.Reverse)
	funcs["IsSorted"] = reflect.ValueOf(sort.IsSorted)
	funcs["Ints"] = reflect.ValueOf(sort.Ints)
	funcs["Float64s"] = reflect.ValueOf(sort.Float64s)
	funcs["Strings"] = reflect.ValueOf(sort.Strings)
	funcs["IntsAreSorted"] = reflect.ValueOf(sort.IntsAreSorted)
	funcs["Float64sAreSorted"] = reflect.ValueOf(sort.Float64sAreSorted)
	funcs["StringsAreSorted"] = reflect.ValueOf(sort.StringsAreSorted)
	funcs["Stable"] = reflect.ValueOf(sort.Stable)

	types = make(map[string] reflect.Type)
	types["Interface"] = reflect.TypeOf(*new(sort.Interface))
	types["IntSlice"] = reflect.TypeOf(*new(sort.IntSlice))
	types["Float64Slice"] = reflect.TypeOf(*new(sort.Float64Slice))
	types["StringSlice"] = reflect.TypeOf(*new(sort.StringSlice))

	vars = make(map[string] reflect.Value)
	pkgs["sort"] = &eval.Env {
		Name: "sort",
		Consts: consts,
		Funcs:  funcs,
		Types:  types,
		Vars:   vars,
		Pkgs:   pkgs,
		Path:   "sort",
	}
	consts = make(map[string] reflect.Value)
	consts["IntSize"] = reflect.ValueOf(strconv.IntSize)

	funcs = make(map[string] reflect.Value)
	funcs["ParseBool"] = reflect.ValueOf(strconv.ParseBool)
	funcs["FormatBool"] = reflect.ValueOf(strconv.FormatBool)
	funcs["AppendBool"] = reflect.ValueOf(strconv.AppendBool)
	funcs["ParseFloat"] = reflect.ValueOf(strconv.ParseFloat)
	funcs["ParseUint"] = reflect.ValueOf(strconv.ParseUint)
	funcs["ParseInt"] = reflect.ValueOf(strconv.ParseInt)
	funcs["Atoi"] = reflect.ValueOf(strconv.Atoi)
	funcs["FormatFloat"] = reflect.ValueOf(strconv.FormatFloat)
	funcs["AppendFloat"] = reflect.ValueOf(strconv.AppendFloat)
	funcs["FormatUint"] = reflect.ValueOf(strconv.FormatUint)
	funcs["FormatInt"] = reflect.ValueOf(strconv.FormatInt)
	funcs["Itoa"] = reflect.ValueOf(strconv.Itoa)
	funcs["AppendInt"] = reflect.ValueOf(strconv.AppendInt)
	funcs["AppendUint"] = reflect.ValueOf(strconv.AppendUint)
	funcs["Quote"] = reflect.ValueOf(strconv.Quote)
	funcs["AppendQuote"] = reflect.ValueOf(strconv.AppendQuote)
	funcs["QuoteToASCII"] = reflect.ValueOf(strconv.QuoteToASCII)
	funcs["AppendQuoteToASCII"] = reflect.ValueOf(strconv.AppendQuoteToASCII)
	funcs["QuoteRune"] = reflect.ValueOf(strconv.QuoteRune)
	funcs["AppendQuoteRune"] = reflect.ValueOf(strconv.AppendQuoteRune)
	funcs["QuoteRuneToASCII"] = reflect.ValueOf(strconv.QuoteRuneToASCII)
	funcs["AppendQuoteRuneToASCII"] = reflect.ValueOf(strconv.AppendQuoteRuneToASCII)
	funcs["CanBackquote"] = reflect.ValueOf(strconv.CanBackquote)
	funcs["UnquoteChar"] = reflect.ValueOf(strconv.UnquoteChar)
	funcs["Unquote"] = reflect.ValueOf(strconv.Unquote)
	funcs["IsPrint"] = reflect.ValueOf(strconv.IsPrint)

	types = make(map[string] reflect.Type)
	types["NumError"] = reflect.TypeOf(*new(strconv.NumError))

	vars = make(map[string] reflect.Value)
	vars["ErrRange"] = reflect.ValueOf(&strconv.ErrRange)
	vars["ErrSyntax"] = reflect.ValueOf(&strconv.ErrSyntax)
	pkgs["strconv"] = &eval.Env {
		Name: "strconv",
		Consts: consts,
		Funcs:  funcs,
		Types:  types,
		Vars:   vars,
		Pkgs:   pkgs,
		Path:   "strconv",
	}
	consts = make(map[string] reflect.Value)

	funcs = make(map[string] reflect.Value)
	funcs["NewReader"] = reflect.ValueOf(strings.NewReader)
	funcs["NewReplacer"] = reflect.ValueOf(strings.NewReplacer)
	funcs["Count"] = reflect.ValueOf(strings.Count)
	funcs["Contains"] = reflect.ValueOf(strings.Contains)
	funcs["ContainsAny"] = reflect.ValueOf(strings.ContainsAny)
	funcs["ContainsRune"] = reflect.ValueOf(strings.ContainsRune)
	funcs["Index"] = reflect.ValueOf(strings.Index)
	funcs["LastIndex"] = reflect.ValueOf(strings.LastIndex)
	funcs["IndexRune"] = reflect.ValueOf(strings.IndexRune)
	funcs["IndexAny"] = reflect.ValueOf(strings.IndexAny)
	funcs["LastIndexAny"] = reflect.ValueOf(strings.LastIndexAny)
	funcs["SplitN"] = reflect.ValueOf(strings.SplitN)
	funcs["SplitAfterN"] = reflect.ValueOf(strings.SplitAfterN)
	funcs["Split"] = reflect.ValueOf(strings.Split)
	funcs["SplitAfter"] = reflect.ValueOf(strings.SplitAfter)
	funcs["Fields"] = reflect.ValueOf(strings.Fields)
	funcs["FieldsFunc"] = reflect.ValueOf(strings.FieldsFunc)
	funcs["Join"] = reflect.ValueOf(strings.Join)
	funcs["HasPrefix"] = reflect.ValueOf(strings.HasPrefix)
	funcs["HasSuffix"] = reflect.ValueOf(strings.HasSuffix)
	funcs["Map"] = reflect.ValueOf(strings.Map)
	funcs["Repeat"] = reflect.ValueOf(strings.Repeat)
	funcs["ToUpper"] = reflect.ValueOf(strings.ToUpper)
	funcs["ToLower"] = reflect.ValueOf(strings.ToLower)
	funcs["ToTitle"] = reflect.ValueOf(strings.ToTitle)
	funcs["ToUpperSpecial"] = reflect.ValueOf(strings.ToUpperSpecial)
	funcs["ToLowerSpecial"] = reflect.ValueOf(strings.ToLowerSpecial)
	funcs["ToTitleSpecial"] = reflect.ValueOf(strings.ToTitleSpecial)
	funcs["Title"] = reflect.ValueOf(strings.Title)
	funcs["TrimLeftFunc"] = reflect.ValueOf(strings.TrimLeftFunc)
	funcs["TrimRightFunc"] = reflect.ValueOf(strings.TrimRightFunc)
	funcs["TrimFunc"] = reflect.ValueOf(strings.TrimFunc)
	funcs["IndexFunc"] = reflect.ValueOf(strings.IndexFunc)
	funcs["LastIndexFunc"] = reflect.ValueOf(strings.LastIndexFunc)
	funcs["Trim"] = reflect.ValueOf(strings.Trim)
	funcs["TrimLeft"] = reflect.ValueOf(strings.TrimLeft)
	funcs["TrimRight"] = reflect.ValueOf(strings.TrimRight)
	funcs["TrimSpace"] = reflect.ValueOf(strings.TrimSpace)
	funcs["TrimPrefix"] = reflect.ValueOf(strings.TrimPrefix)
	funcs["TrimSuffix"] = reflect.ValueOf(strings.TrimSuffix)
	funcs["Replace"] = reflect.ValueOf(strings.Replace)
	funcs["EqualFold"] = reflect.ValueOf(strings.EqualFold)
	funcs["IndexByte"] = reflect.ValueOf(strings.IndexByte)

	types = make(map[string] reflect.Type)
	types["Reader"] = reflect.TypeOf(*new(strings.Reader))
	types["Replacer"] = reflect.TypeOf(*new(strings.Replacer))

	vars = make(map[string] reflect.Value)
	pkgs["strings"] = &eval.Env {
		Name: "strings",
		Consts: consts,
		Funcs:  funcs,
		Types:  types,
		Vars:   vars,
		Pkgs:   pkgs,
		Path:   "strings",
	}
	consts = make(map[string] reflect.Value)

	funcs = make(map[string] reflect.Value)
	funcs["NewCond"] = reflect.ValueOf(sync.NewCond)

	types = make(map[string] reflect.Type)
	types["Cond"] = reflect.TypeOf(*new(sync.Cond))
	types["Mutex"] = reflect.TypeOf(*new(sync.Mutex))
	types["Locker"] = reflect.TypeOf(*new(sync.Locker))
	types["Once"] = reflect.TypeOf(*new(sync.Once))
	types["RWMutex"] = reflect.TypeOf(*new(sync.RWMutex))
	types["WaitGroup"] = reflect.TypeOf(*new(sync.WaitGroup))

	vars = make(map[string] reflect.Value)
	pkgs["sync"] = &eval.Env {
		Name: "sync",
		Consts: consts,
		Funcs:  funcs,
		Types:  types,
		Vars:   vars,
		Pkgs:   pkgs,
		Path:   "sync",
	}
	consts = make(map[string] reflect.Value)

	funcs = make(map[string] reflect.Value)
	funcs["SwapInt32"] = reflect.ValueOf(atomic.SwapInt32)
	funcs["SwapInt64"] = reflect.ValueOf(atomic.SwapInt64)
	funcs["SwapUint32"] = reflect.ValueOf(atomic.SwapUint32)
	funcs["SwapUint64"] = reflect.ValueOf(atomic.SwapUint64)
	funcs["SwapUintptr"] = reflect.ValueOf(atomic.SwapUintptr)
	funcs["SwapPointer"] = reflect.ValueOf(atomic.SwapPointer)
	funcs["CompareAndSwapInt32"] = reflect.ValueOf(atomic.CompareAndSwapInt32)
	funcs["CompareAndSwapInt64"] = reflect.ValueOf(atomic.CompareAndSwapInt64)
	funcs["CompareAndSwapUint32"] = reflect.ValueOf(atomic.CompareAndSwapUint32)
	funcs["CompareAndSwapUint64"] = reflect.ValueOf(atomic.CompareAndSwapUint64)
	funcs["CompareAndSwapUintptr"] = reflect.ValueOf(atomic.CompareAndSwapUintptr)
	funcs["CompareAndSwapPointer"] = reflect.ValueOf(atomic.CompareAndSwapPointer)
	funcs["AddInt32"] = reflect.ValueOf(atomic.AddInt32)
	funcs["AddUint32"] = reflect.ValueOf(atomic.AddUint32)
	funcs["AddInt64"] = reflect.ValueOf(atomic.AddInt64)
	funcs["AddUint64"] = reflect.ValueOf(atomic.AddUint64)
	funcs["AddUintptr"] = reflect.ValueOf(atomic.AddUintptr)
	funcs["LoadInt32"] = reflect.ValueOf(atomic.LoadInt32)
	funcs["LoadInt64"] = reflect.ValueOf(atomic.LoadInt64)
	funcs["LoadUint32"] = reflect.ValueOf(atomic.LoadUint32)
	funcs["LoadUint64"] = reflect.ValueOf(atomic.LoadUint64)
	funcs["LoadUintptr"] = reflect.ValueOf(atomic.LoadUintptr)
	funcs["LoadPointer"] = reflect.ValueOf(atomic.LoadPointer)
	funcs["StoreInt32"] = reflect.ValueOf(atomic.StoreInt32)
	funcs["StoreInt64"] = reflect.ValueOf(atomic.StoreInt64)
	funcs["StoreUint32"] = reflect.ValueOf(atomic.StoreUint32)
	funcs["StoreUint64"] = reflect.ValueOf(atomic.StoreUint64)
	funcs["StoreUintptr"] = reflect.ValueOf(atomic.StoreUintptr)
	funcs["StorePointer"] = reflect.ValueOf(atomic.StorePointer)

	types = make(map[string] reflect.Type)

	vars = make(map[string] reflect.Value)
	pkgs["atomic"] = &eval.Env {
		Name: "atomic",
		Consts: consts,
		Funcs:  funcs,
		Types:  types,
		Vars:   vars,
		Pkgs:   pkgs,
		Path:   "sync/atomic",
	}
	consts = make(map[string] reflect.Value)
	//syscall constants excluded

	funcs = make(map[string] reflect.Value)
	funcs["Getenv"] = reflect.ValueOf(syscall.Getenv)
	funcs["Setenv"] = reflect.ValueOf(syscall.Setenv)
	funcs["Clearenv"] = reflect.ValueOf(syscall.Clearenv)
	funcs["Environ"] = reflect.ValueOf(syscall.Environ)
	funcs["CloseOnExec"] = reflect.ValueOf(syscall.CloseOnExec)
	funcs["SetNonblock"] = reflect.ValueOf(syscall.SetNonblock)
	funcs["StartProcess"] = reflect.ValueOf(syscall.StartProcess)
	funcs["Exec"] = reflect.ValueOf(syscall.Exec)
	funcs["StringByteSlice"] = reflect.ValueOf(syscall.StringByteSlice)
	funcs["ByteSliceFromString"] = reflect.ValueOf(syscall.ByteSliceFromString)
	funcs["StringBytePtr"] = reflect.ValueOf(syscall.StringBytePtr)
	funcs["BytePtrFromString"] = reflect.ValueOf(syscall.BytePtrFromString)
	funcs["Open"] = reflect.ValueOf(syscall.Open)
	funcs["Pipe"] = reflect.ValueOf(syscall.Pipe)
	funcs["Utimes"] = reflect.ValueOf(syscall.Utimes)
	funcs["UtimesNano"] = reflect.ValueOf(syscall.UtimesNano)
	funcs["Getwd"] = reflect.ValueOf(syscall.Getwd)
	funcs["Getgroups"] = reflect.ValueOf(syscall.Getgroups)
	funcs["Accept"] = reflect.ValueOf(syscall.Accept)
	funcs["Getsockname"] = reflect.ValueOf(syscall.Getsockname)
	funcs["Getpagesize"] = reflect.ValueOf(syscall.Getpagesize)
	funcs["TimespecToNsec"] = reflect.ValueOf(syscall.TimespecToNsec)
	funcs["NsecToTimespec"] = reflect.ValueOf(syscall.NsecToTimespec)
	funcs["NsecToTimeval"] = reflect.ValueOf(syscall.NsecToTimeval)
	funcs["Seek"] = reflect.ValueOf(syscall.Seek)
	funcs["Listen"] = reflect.ValueOf(syscall.Listen)
	funcs["Shutdown"] = reflect.ValueOf(syscall.Shutdown)
	funcs["Syscall"] = reflect.ValueOf(syscall.Syscall)
	funcs["Syscall6"] = reflect.ValueOf(syscall.Syscall6)
	funcs["Read"] = reflect.ValueOf(syscall.Read)
	funcs["Write"] = reflect.ValueOf(syscall.Write)
	funcs["Bind"] = reflect.ValueOf(syscall.Bind)
	funcs["Connect"] = reflect.ValueOf(syscall.Connect)
	funcs["Getpeername"] = reflect.ValueOf(syscall.Getpeername)
	funcs["GetsockoptInt"] = reflect.ValueOf(syscall.GetsockoptInt)
	funcs["Recvfrom"] = reflect.ValueOf(syscall.Recvfrom)
	funcs["Sendto"] = reflect.ValueOf(syscall.Sendto)
	funcs["SetsockoptInt"] = reflect.ValueOf(syscall.SetsockoptInt)
	funcs["SetsockoptInet4Addr"] = reflect.ValueOf(syscall.SetsockoptInet4Addr)
	funcs["SetsockoptIPMreq"] = reflect.ValueOf(syscall.SetsockoptIPMreq)
	funcs["SetsockoptIPv6Mreq"] = reflect.ValueOf(syscall.SetsockoptIPv6Mreq)
	funcs["SetsockoptLinger"] = reflect.ValueOf(syscall.SetsockoptLinger)
	funcs["SetsockoptTimeval"] = reflect.ValueOf(syscall.SetsockoptTimeval)
	funcs["Socket"] = reflect.ValueOf(syscall.Socket)
	funcs["Chdir"] = reflect.ValueOf(syscall.Chdir)
	funcs["Chmod"] = reflect.ValueOf(syscall.Chmod)
	funcs["Close"] = reflect.ValueOf(syscall.Close)
	funcs["Exit"] = reflect.ValueOf(syscall.Exit)
	funcs["Fchdir"] = reflect.ValueOf(syscall.Fchdir)
	funcs["Fchmod"] = reflect.ValueOf(syscall.Fchmod)
	funcs["Fsync"] = reflect.ValueOf(syscall.Fsync)
	funcs["Getpid"] = reflect.ValueOf(syscall.Getpid)
	funcs["Getppid"] = reflect.ValueOf(syscall.Getppid)
	funcs["Link"] = reflect.ValueOf(syscall.Link)
	funcs["Mkdir"] = reflect.ValueOf(syscall.Mkdir)
	funcs["Readlink"] = reflect.ValueOf(syscall.Readlink)
	funcs["Rename"] = reflect.ValueOf(syscall.Rename)
	funcs["Rmdir"] = reflect.ValueOf(syscall.Rmdir)
	funcs["Symlink"] = reflect.ValueOf(syscall.Symlink)
	funcs["Unlink"] = reflect.ValueOf(syscall.Unlink)
	funcs["Chown"] = reflect.ValueOf(syscall.Chown)
	funcs["Fchown"] = reflect.ValueOf(syscall.Fchown)
	funcs["Ftruncate"] = reflect.ValueOf(syscall.Ftruncate)
	funcs["Getegid"] = reflect.ValueOf(syscall.Getegid)
	funcs["Geteuid"] = reflect.ValueOf(syscall.Geteuid)
	funcs["Getgid"] = reflect.ValueOf(syscall.Getgid)
	funcs["Getuid"] = reflect.ValueOf(syscall.Getuid)
	funcs["Lchown"] = reflect.ValueOf(syscall.Lchown)

    addArchSyscallFuncs(funcs)

	types = make(map[string] reflect.Type)
	types["SysProcAttr"] = reflect.TypeOf(*new(syscall.SysProcAttr))
	types["ProcAttr"] = reflect.TypeOf(*new(syscall.ProcAttr))
	types["WaitStatus"] = reflect.TypeOf(*new(syscall.WaitStatus))
	types["Errno"] = reflect.TypeOf(*new(syscall.Errno))
	types["Signal"] = reflect.TypeOf(*new(syscall.Signal))
	types["Sockaddr"] = reflect.TypeOf(*new(syscall.Sockaddr))
	types["SockaddrInet4"] = reflect.TypeOf(*new(syscall.SockaddrInet4))
	types["SockaddrInet6"] = reflect.TypeOf(*new(syscall.SockaddrInet6))
	types["SockaddrUnix"] = reflect.TypeOf(*new(syscall.SockaddrUnix))
	types["Timespec"] = reflect.TypeOf(*new(syscall.Timespec))
	types["Timeval"] = reflect.TypeOf(*new(syscall.Timeval))
	types["Rusage"] = reflect.TypeOf(*new(syscall.Rusage))
	types["RawSockaddrInet4"] = reflect.TypeOf(*new(syscall.RawSockaddrInet4))
	types["RawSockaddrInet6"] = reflect.TypeOf(*new(syscall.RawSockaddrInet6))
	types["RawSockaddrAny"] = reflect.TypeOf(*new(syscall.RawSockaddrAny))
	types["RawSockaddr"] = reflect.TypeOf(*new(syscall.RawSockaddr))
	types["Linger"] = reflect.TypeOf(*new(syscall.Linger))
	types["IPMreq"] = reflect.TypeOf(*new(syscall.IPMreq))
	types["IPv6Mreq"] = reflect.TypeOf(*new(syscall.IPv6Mreq))

    addArchSyscallTypes(types)

	vars = make(map[string] reflect.Value)
	vars["ForkLock"] = reflect.ValueOf(&syscall.ForkLock)
	vars["Stdin"] = reflect.ValueOf(&syscall.Stdin)
	vars["Stdout"] = reflect.ValueOf(&syscall.Stdout)
	vars["Stderr"] = reflect.ValueOf(&syscall.Stderr)
	vars["SocketDisableIPv6"] = reflect.ValueOf(&syscall.SocketDisableIPv6)
	pkgs["syscall"] = &eval.Env {
		Name: "syscall",
		Consts: consts,
		Funcs:  funcs,
		Types:  types,
		Vars:   vars,
		Pkgs:   pkgs,
		Path:   "syscall",
	}
	consts = make(map[string] reflect.Value)

	funcs = make(map[string] reflect.Value)
	funcs["AllocsPerRun"] = reflect.ValueOf(testing.AllocsPerRun)
	funcs["RunBenchmarks"] = reflect.ValueOf(testing.RunBenchmarks)
	funcs["Benchmark"] = reflect.ValueOf(testing.Benchmark)
	funcs["RegisterCover"] = reflect.ValueOf(testing.RegisterCover)
	funcs["RunExamples"] = reflect.ValueOf(testing.RunExamples)
	funcs["Short"] = reflect.ValueOf(testing.Short)
	funcs["Verbose"] = reflect.ValueOf(testing.Verbose)
	funcs["Main"] = reflect.ValueOf(testing.Main)
	funcs["RunTests"] = reflect.ValueOf(testing.RunTests)

	types = make(map[string] reflect.Type)
	types["InternalBenchmark"] = reflect.TypeOf(*new(testing.InternalBenchmark))
	types["B"] = reflect.TypeOf(*new(testing.B))
	types["BenchmarkResult"] = reflect.TypeOf(*new(testing.BenchmarkResult))
	types["CoverBlock"] = reflect.TypeOf(*new(testing.CoverBlock))
	types["Cover"] = reflect.TypeOf(*new(testing.Cover))
	types["InternalExample"] = reflect.TypeOf(*new(testing.InternalExample))
	types["TB"] = reflect.TypeOf(*new(testing.TB))
	types["T"] = reflect.TypeOf(*new(testing.T))
	types["InternalTest"] = reflect.TypeOf(*new(testing.InternalTest))

	vars = make(map[string] reflect.Value)
	pkgs["testing"] = &eval.Env {
		Name: "testing",
		Consts: consts,
		Funcs:  funcs,
		Types:  types,
		Vars:   vars,
		Pkgs:   pkgs,
		Path:   "testing",
	}
	consts = make(map[string] reflect.Value)
	consts["FilterHTML"] = reflect.ValueOf(tabwriter.FilterHTML)
	consts["StripEscape"] = reflect.ValueOf(tabwriter.StripEscape)
	consts["AlignRight"] = reflect.ValueOf(tabwriter.AlignRight)
	consts["DiscardEmptyColumns"] = reflect.ValueOf(tabwriter.DiscardEmptyColumns)
	consts["TabIndent"] = reflect.ValueOf(tabwriter.TabIndent)
	consts["Debug"] = reflect.ValueOf(tabwriter.Debug)
	consts["Escape"] = reflect.ValueOf(tabwriter.Escape)

	funcs = make(map[string] reflect.Value)
	funcs["NewWriter"] = reflect.ValueOf(tabwriter.NewWriter)

	types = make(map[string] reflect.Type)
	types["Writer"] = reflect.TypeOf(*new(tabwriter.Writer))

	vars = make(map[string] reflect.Value)
	pkgs["tabwriter"] = &eval.Env {
		Name: "tabwriter",
		Consts: consts,
		Funcs:  funcs,
		Types:  types,
		Vars:   vars,
		Pkgs:   pkgs,
		Path:   "text/tabwriter",
	}
	consts = make(map[string] reflect.Value)
	consts["ANSIC"] = reflect.ValueOf(time.ANSIC)
	consts["UnixDate"] = reflect.ValueOf(time.UnixDate)
	consts["RubyDate"] = reflect.ValueOf(time.RubyDate)
	consts["RFC822"] = reflect.ValueOf(time.RFC822)
	consts["RFC822Z"] = reflect.ValueOf(time.RFC822Z)
	consts["RFC850"] = reflect.ValueOf(time.RFC850)
	consts["RFC1123"] = reflect.ValueOf(time.RFC1123)
	consts["RFC1123Z"] = reflect.ValueOf(time.RFC1123Z)
	consts["RFC3339"] = reflect.ValueOf(time.RFC3339)
	consts["RFC3339Nano"] = reflect.ValueOf(time.RFC3339Nano)
	consts["Kitchen"] = reflect.ValueOf(time.Kitchen)
	consts["Stamp"] = reflect.ValueOf(time.Stamp)
	consts["StampMilli"] = reflect.ValueOf(time.StampMilli)
	consts["StampMicro"] = reflect.ValueOf(time.StampMicro)
	consts["StampNano"] = reflect.ValueOf(time.StampNano)
	consts["January"] = reflect.ValueOf(time.January)
	consts["February"] = reflect.ValueOf(time.February)
	consts["March"] = reflect.ValueOf(time.March)
	consts["April"] = reflect.ValueOf(time.April)
	consts["May"] = reflect.ValueOf(time.May)
	consts["June"] = reflect.ValueOf(time.June)
	consts["July"] = reflect.ValueOf(time.July)
	consts["August"] = reflect.ValueOf(time.August)
	consts["September"] = reflect.ValueOf(time.September)
	consts["October"] = reflect.ValueOf(time.October)
	consts["November"] = reflect.ValueOf(time.November)
	consts["December"] = reflect.ValueOf(time.December)
	consts["Sunday"] = reflect.ValueOf(time.Sunday)
	consts["Monday"] = reflect.ValueOf(time.Monday)
	consts["Tuesday"] = reflect.ValueOf(time.Tuesday)
	consts["Wednesday"] = reflect.ValueOf(time.Wednesday)
	consts["Thursday"] = reflect.ValueOf(time.Thursday)
	consts["Friday"] = reflect.ValueOf(time.Friday)
	consts["Saturday"] = reflect.ValueOf(time.Saturday)
	consts["Nanosecond"] = reflect.ValueOf(time.Nanosecond)
	consts["Microsecond"] = reflect.ValueOf(time.Microsecond)
	consts["Millisecond"] = reflect.ValueOf(time.Millisecond)
	consts["Second"] = reflect.ValueOf(time.Second)
	consts["Minute"] = reflect.ValueOf(time.Minute)
	consts["Hour"] = reflect.ValueOf(time.Hour)

	funcs = make(map[string] reflect.Value)
	funcs["Parse"] = reflect.ValueOf(time.Parse)
	funcs["ParseInLocation"] = reflect.ValueOf(time.ParseInLocation)
	funcs["ParseDuration"] = reflect.ValueOf(time.ParseDuration)
	funcs["Sleep"] = reflect.ValueOf(time.Sleep)
	funcs["NewTimer"] = reflect.ValueOf(time.NewTimer)
	funcs["After"] = reflect.ValueOf(time.After)
	funcs["AfterFunc"] = reflect.ValueOf(time.AfterFunc)
	funcs["NewTicker"] = reflect.ValueOf(time.NewTicker)
	funcs["Tick"] = reflect.ValueOf(time.Tick)
	funcs["Since"] = reflect.ValueOf(time.Since)
	funcs["Now"] = reflect.ValueOf(time.Now)
	funcs["Unix"] = reflect.ValueOf(time.Unix)
	funcs["Date"] = reflect.ValueOf(time.Date)
	funcs["FixedZone"] = reflect.ValueOf(time.FixedZone)
	funcs["LoadLocation"] = reflect.ValueOf(time.LoadLocation)

	types = make(map[string] reflect.Type)
	types["ParseError"] = reflect.TypeOf(*new(time.ParseError))
	types["Timer"] = reflect.TypeOf(*new(time.Timer))
	types["Ticker"] = reflect.TypeOf(*new(time.Ticker))
	types["Time"] = reflect.TypeOf(*new(time.Time))
	types["Month"] = reflect.TypeOf(*new(time.Month))
	types["Weekday"] = reflect.TypeOf(*new(time.Weekday))
	types["Duration"] = reflect.TypeOf(*new(time.Duration))
	types["Location"] = reflect.TypeOf(*new(time.Location))

	vars = make(map[string] reflect.Value)
	vars["UTC"] = reflect.ValueOf(&time.UTC)
	vars["Local"] = reflect.ValueOf(&time.Local)
	pkgs["time"] = &eval.Env {
		Name: "time",
		Consts: consts,
		Funcs:  funcs,
		Types:  types,
		Vars:   vars,
		Pkgs:   pkgs,
		Path:   "time",
	}
	consts = make(map[string] reflect.Value)
	consts["MaxRune"] = reflect.ValueOf(unicode.MaxRune)
	consts["ReplacementChar"] = reflect.ValueOf(unicode.ReplacementChar)
	consts["MaxASCII"] = reflect.ValueOf(unicode.MaxASCII)
	consts["MaxLatin1"] = reflect.ValueOf(unicode.MaxLatin1)
	consts["UpperCase"] = reflect.ValueOf(unicode.UpperCase)
	consts["LowerCase"] = reflect.ValueOf(unicode.LowerCase)
	consts["TitleCase"] = reflect.ValueOf(unicode.TitleCase)
	consts["MaxCase"] = reflect.ValueOf(unicode.MaxCase)
	consts["UpperLower"] = reflect.ValueOf(unicode.UpperLower)
	consts["Version"] = reflect.ValueOf(unicode.Version)

	funcs = make(map[string] reflect.Value)
	funcs["IsDigit"] = reflect.ValueOf(unicode.IsDigit)
	funcs["IsGraphic"] = reflect.ValueOf(unicode.IsGraphic)
	funcs["IsPrint"] = reflect.ValueOf(unicode.IsPrint)
	funcs["IsOneOf"] = reflect.ValueOf(unicode.IsOneOf)
	funcs["In"] = reflect.ValueOf(unicode.In)
	funcs["IsControl"] = reflect.ValueOf(unicode.IsControl)
	funcs["IsLetter"] = reflect.ValueOf(unicode.IsLetter)
	funcs["IsMark"] = reflect.ValueOf(unicode.IsMark)
	funcs["IsNumber"] = reflect.ValueOf(unicode.IsNumber)
	funcs["IsPunct"] = reflect.ValueOf(unicode.IsPunct)
	funcs["IsSpace"] = reflect.ValueOf(unicode.IsSpace)
	funcs["IsSymbol"] = reflect.ValueOf(unicode.IsSymbol)
	funcs["Is"] = reflect.ValueOf(unicode.Is)
	funcs["IsUpper"] = reflect.ValueOf(unicode.IsUpper)
	funcs["IsLower"] = reflect.ValueOf(unicode.IsLower)
	funcs["IsTitle"] = reflect.ValueOf(unicode.IsTitle)
	funcs["To"] = reflect.ValueOf(unicode.To)
	funcs["ToUpper"] = reflect.ValueOf(unicode.ToUpper)
	funcs["ToLower"] = reflect.ValueOf(unicode.ToLower)
	funcs["ToTitle"] = reflect.ValueOf(unicode.ToTitle)
	funcs["SimpleFold"] = reflect.ValueOf(unicode.SimpleFold)

	types = make(map[string] reflect.Type)
	types["RangeTable"] = reflect.TypeOf(*new(unicode.RangeTable))
	types["Range16"] = reflect.TypeOf(*new(unicode.Range16))
	types["Range32"] = reflect.TypeOf(*new(unicode.Range32))
	types["CaseRange"] = reflect.TypeOf(*new(unicode.CaseRange))
	types["SpecialCase"] = reflect.TypeOf(*new(unicode.SpecialCase))

	vars = make(map[string] reflect.Value)
	vars["TurkishCase"] = reflect.ValueOf(&unicode.TurkishCase)
	vars["AzeriCase"] = reflect.ValueOf(&unicode.AzeriCase)
	vars["GraphicRanges"] = reflect.ValueOf(&unicode.GraphicRanges)
	vars["PrintRanges"] = reflect.ValueOf(&unicode.PrintRanges)
	vars["Categories"] = reflect.ValueOf(&unicode.Categories)
	vars["Cc"] = reflect.ValueOf(&unicode.Cc)
	vars["Cf"] = reflect.ValueOf(&unicode.Cf)
	vars["Co"] = reflect.ValueOf(&unicode.Co)
	vars["Cs"] = reflect.ValueOf(&unicode.Cs)
	vars["Digit"] = reflect.ValueOf(&unicode.Digit)
	vars["Nd"] = reflect.ValueOf(&unicode.Nd)
	vars["Letter"] = reflect.ValueOf(&unicode.Letter)
	vars["L"] = reflect.ValueOf(&unicode.L)
	vars["Lm"] = reflect.ValueOf(&unicode.Lm)
	vars["Lo"] = reflect.ValueOf(&unicode.Lo)
	vars["Lower"] = reflect.ValueOf(&unicode.Lower)
	vars["Ll"] = reflect.ValueOf(&unicode.Ll)
	vars["Mark"] = reflect.ValueOf(&unicode.Mark)
	vars["M"] = reflect.ValueOf(&unicode.M)
	vars["Mc"] = reflect.ValueOf(&unicode.Mc)
	vars["Me"] = reflect.ValueOf(&unicode.Me)
	vars["Mn"] = reflect.ValueOf(&unicode.Mn)
	vars["Nl"] = reflect.ValueOf(&unicode.Nl)
	vars["No"] = reflect.ValueOf(&unicode.No)
	vars["Number"] = reflect.ValueOf(&unicode.Number)
	vars["N"] = reflect.ValueOf(&unicode.N)
	vars["Other"] = reflect.ValueOf(&unicode.Other)
	vars["C"] = reflect.ValueOf(&unicode.C)
	vars["Pc"] = reflect.ValueOf(&unicode.Pc)
	vars["Pd"] = reflect.ValueOf(&unicode.Pd)
	vars["Pe"] = reflect.ValueOf(&unicode.Pe)
	vars["Pf"] = reflect.ValueOf(&unicode.Pf)
	vars["Pi"] = reflect.ValueOf(&unicode.Pi)
	vars["Po"] = reflect.ValueOf(&unicode.Po)
	vars["Ps"] = reflect.ValueOf(&unicode.Ps)
	vars["Punct"] = reflect.ValueOf(&unicode.Punct)
	vars["P"] = reflect.ValueOf(&unicode.P)
	vars["Sc"] = reflect.ValueOf(&unicode.Sc)
	vars["Sk"] = reflect.ValueOf(&unicode.Sk)
	vars["Sm"] = reflect.ValueOf(&unicode.Sm)
	vars["So"] = reflect.ValueOf(&unicode.So)
	vars["Space"] = reflect.ValueOf(&unicode.Space)
	vars["Z"] = reflect.ValueOf(&unicode.Z)
	vars["Symbol"] = reflect.ValueOf(&unicode.Symbol)
	vars["S"] = reflect.ValueOf(&unicode.S)
	vars["Title"] = reflect.ValueOf(&unicode.Title)
	vars["Lt"] = reflect.ValueOf(&unicode.Lt)
	vars["Upper"] = reflect.ValueOf(&unicode.Upper)
	vars["Lu"] = reflect.ValueOf(&unicode.Lu)
	vars["Zl"] = reflect.ValueOf(&unicode.Zl)
	vars["Zp"] = reflect.ValueOf(&unicode.Zp)
	vars["Zs"] = reflect.ValueOf(&unicode.Zs)
	vars["Scripts"] = reflect.ValueOf(&unicode.Scripts)
	vars["Arabic"] = reflect.ValueOf(&unicode.Arabic)
	vars["Armenian"] = reflect.ValueOf(&unicode.Armenian)
	vars["Avestan"] = reflect.ValueOf(&unicode.Avestan)
	vars["Balinese"] = reflect.ValueOf(&unicode.Balinese)
	vars["Bamum"] = reflect.ValueOf(&unicode.Bamum)
	vars["Batak"] = reflect.ValueOf(&unicode.Batak)
	vars["Bengali"] = reflect.ValueOf(&unicode.Bengali)
	vars["Bopomofo"] = reflect.ValueOf(&unicode.Bopomofo)
	vars["Brahmi"] = reflect.ValueOf(&unicode.Brahmi)
	vars["Braille"] = reflect.ValueOf(&unicode.Braille)
	vars["Buginese"] = reflect.ValueOf(&unicode.Buginese)
	vars["Buhid"] = reflect.ValueOf(&unicode.Buhid)
	vars["Canadian_Aboriginal"] = reflect.ValueOf(&unicode.Canadian_Aboriginal)
	vars["Carian"] = reflect.ValueOf(&unicode.Carian)
	vars["Chakma"] = reflect.ValueOf(&unicode.Chakma)
	vars["Cham"] = reflect.ValueOf(&unicode.Cham)
	vars["Cherokee"] = reflect.ValueOf(&unicode.Cherokee)
	vars["Common"] = reflect.ValueOf(&unicode.Common)
	vars["Coptic"] = reflect.ValueOf(&unicode.Coptic)
	vars["Cuneiform"] = reflect.ValueOf(&unicode.Cuneiform)
	vars["Cypriot"] = reflect.ValueOf(&unicode.Cypriot)
	vars["Cyrillic"] = reflect.ValueOf(&unicode.Cyrillic)
	vars["Deseret"] = reflect.ValueOf(&unicode.Deseret)
	vars["Devanagari"] = reflect.ValueOf(&unicode.Devanagari)
	vars["Egyptian_Hieroglyphs"] = reflect.ValueOf(&unicode.Egyptian_Hieroglyphs)
	vars["Ethiopic"] = reflect.ValueOf(&unicode.Ethiopic)
	vars["Georgian"] = reflect.ValueOf(&unicode.Georgian)
	vars["Glagolitic"] = reflect.ValueOf(&unicode.Glagolitic)
	vars["Gothic"] = reflect.ValueOf(&unicode.Gothic)
	vars["Greek"] = reflect.ValueOf(&unicode.Greek)
	vars["Gujarati"] = reflect.ValueOf(&unicode.Gujarati)
	vars["Gurmukhi"] = reflect.ValueOf(&unicode.Gurmukhi)
	vars["Han"] = reflect.ValueOf(&unicode.Han)
	vars["Hangul"] = reflect.ValueOf(&unicode.Hangul)
	vars["Hanunoo"] = reflect.ValueOf(&unicode.Hanunoo)
	vars["Hebrew"] = reflect.ValueOf(&unicode.Hebrew)
	vars["Hiragana"] = reflect.ValueOf(&unicode.Hiragana)
	vars["Imperial_Aramaic"] = reflect.ValueOf(&unicode.Imperial_Aramaic)
	vars["Inherited"] = reflect.ValueOf(&unicode.Inherited)
	vars["Inscriptional_Pahlavi"] = reflect.ValueOf(&unicode.Inscriptional_Pahlavi)
	vars["Inscriptional_Parthian"] = reflect.ValueOf(&unicode.Inscriptional_Parthian)
	vars["Javanese"] = reflect.ValueOf(&unicode.Javanese)
	vars["Kaithi"] = reflect.ValueOf(&unicode.Kaithi)
	vars["Kannada"] = reflect.ValueOf(&unicode.Kannada)
	vars["Katakana"] = reflect.ValueOf(&unicode.Katakana)
	vars["Kayah_Li"] = reflect.ValueOf(&unicode.Kayah_Li)
	vars["Kharoshthi"] = reflect.ValueOf(&unicode.Kharoshthi)
	vars["Khmer"] = reflect.ValueOf(&unicode.Khmer)
	vars["Lao"] = reflect.ValueOf(&unicode.Lao)
	vars["Latin"] = reflect.ValueOf(&unicode.Latin)
	vars["Lepcha"] = reflect.ValueOf(&unicode.Lepcha)
	vars["Limbu"] = reflect.ValueOf(&unicode.Limbu)
	vars["Linear_B"] = reflect.ValueOf(&unicode.Linear_B)
	vars["Lisu"] = reflect.ValueOf(&unicode.Lisu)
	vars["Lycian"] = reflect.ValueOf(&unicode.Lycian)
	vars["Lydian"] = reflect.ValueOf(&unicode.Lydian)
	vars["Malayalam"] = reflect.ValueOf(&unicode.Malayalam)
	vars["Mandaic"] = reflect.ValueOf(&unicode.Mandaic)
	vars["Meetei_Mayek"] = reflect.ValueOf(&unicode.Meetei_Mayek)
	vars["Meroitic_Cursive"] = reflect.ValueOf(&unicode.Meroitic_Cursive)
	vars["Meroitic_Hieroglyphs"] = reflect.ValueOf(&unicode.Meroitic_Hieroglyphs)
	vars["Miao"] = reflect.ValueOf(&unicode.Miao)
	vars["Mongolian"] = reflect.ValueOf(&unicode.Mongolian)
	vars["Myanmar"] = reflect.ValueOf(&unicode.Myanmar)
	vars["New_Tai_Lue"] = reflect.ValueOf(&unicode.New_Tai_Lue)
	vars["Nko"] = reflect.ValueOf(&unicode.Nko)
	vars["Ogham"] = reflect.ValueOf(&unicode.Ogham)
	vars["Ol_Chiki"] = reflect.ValueOf(&unicode.Ol_Chiki)
	vars["Old_Italic"] = reflect.ValueOf(&unicode.Old_Italic)
	vars["Old_Persian"] = reflect.ValueOf(&unicode.Old_Persian)
	vars["Old_South_Arabian"] = reflect.ValueOf(&unicode.Old_South_Arabian)
	vars["Old_Turkic"] = reflect.ValueOf(&unicode.Old_Turkic)
	vars["Oriya"] = reflect.ValueOf(&unicode.Oriya)
	vars["Osmanya"] = reflect.ValueOf(&unicode.Osmanya)
	vars["Phags_Pa"] = reflect.ValueOf(&unicode.Phags_Pa)
	vars["Phoenician"] = reflect.ValueOf(&unicode.Phoenician)
	vars["Rejang"] = reflect.ValueOf(&unicode.Rejang)
	vars["Runic"] = reflect.ValueOf(&unicode.Runic)
	vars["Samaritan"] = reflect.ValueOf(&unicode.Samaritan)
	vars["Saurashtra"] = reflect.ValueOf(&unicode.Saurashtra)
	vars["Sharada"] = reflect.ValueOf(&unicode.Sharada)
	vars["Shavian"] = reflect.ValueOf(&unicode.Shavian)
	vars["Sinhala"] = reflect.ValueOf(&unicode.Sinhala)
	vars["Sora_Sompeng"] = reflect.ValueOf(&unicode.Sora_Sompeng)
	vars["Sundanese"] = reflect.ValueOf(&unicode.Sundanese)
	vars["Syloti_Nagri"] = reflect.ValueOf(&unicode.Syloti_Nagri)
	vars["Syriac"] = reflect.ValueOf(&unicode.Syriac)
	vars["Tagalog"] = reflect.ValueOf(&unicode.Tagalog)
	vars["Tagbanwa"] = reflect.ValueOf(&unicode.Tagbanwa)
	vars["Tai_Le"] = reflect.ValueOf(&unicode.Tai_Le)
	vars["Tai_Tham"] = reflect.ValueOf(&unicode.Tai_Tham)
	vars["Tai_Viet"] = reflect.ValueOf(&unicode.Tai_Viet)
	vars["Takri"] = reflect.ValueOf(&unicode.Takri)
	vars["Tamil"] = reflect.ValueOf(&unicode.Tamil)
	vars["Telugu"] = reflect.ValueOf(&unicode.Telugu)
	vars["Thaana"] = reflect.ValueOf(&unicode.Thaana)
	vars["Thai"] = reflect.ValueOf(&unicode.Thai)
	vars["Tibetan"] = reflect.ValueOf(&unicode.Tibetan)
	vars["Tifinagh"] = reflect.ValueOf(&unicode.Tifinagh)
	vars["Ugaritic"] = reflect.ValueOf(&unicode.Ugaritic)
	vars["Vai"] = reflect.ValueOf(&unicode.Vai)
	vars["Yi"] = reflect.ValueOf(&unicode.Yi)
	vars["Properties"] = reflect.ValueOf(&unicode.Properties)
	vars["ASCII_Hex_Digit"] = reflect.ValueOf(&unicode.ASCII_Hex_Digit)
	vars["Bidi_Control"] = reflect.ValueOf(&unicode.Bidi_Control)
	vars["Dash"] = reflect.ValueOf(&unicode.Dash)
	vars["Deprecated"] = reflect.ValueOf(&unicode.Deprecated)
	vars["Diacritic"] = reflect.ValueOf(&unicode.Diacritic)
	vars["Extender"] = reflect.ValueOf(&unicode.Extender)
	vars["Hex_Digit"] = reflect.ValueOf(&unicode.Hex_Digit)
	vars["Hyphen"] = reflect.ValueOf(&unicode.Hyphen)
	vars["IDS_Binary_Operator"] = reflect.ValueOf(&unicode.IDS_Binary_Operator)
	vars["IDS_Trinary_Operator"] = reflect.ValueOf(&unicode.IDS_Trinary_Operator)
	vars["Ideographic"] = reflect.ValueOf(&unicode.Ideographic)
	vars["Join_Control"] = reflect.ValueOf(&unicode.Join_Control)
	vars["Logical_Order_Exception"] = reflect.ValueOf(&unicode.Logical_Order_Exception)
	vars["Noncharacter_Code_Point"] = reflect.ValueOf(&unicode.Noncharacter_Code_Point)
	vars["Other_Alphabetic"] = reflect.ValueOf(&unicode.Other_Alphabetic)
	vars["Other_Default_Ignorable_Code_Point"] = reflect.ValueOf(&unicode.Other_Default_Ignorable_Code_Point)
	vars["Other_Grapheme_Extend"] = reflect.ValueOf(&unicode.Other_Grapheme_Extend)
	vars["Other_ID_Continue"] = reflect.ValueOf(&unicode.Other_ID_Continue)
	vars["Other_ID_Start"] = reflect.ValueOf(&unicode.Other_ID_Start)
	vars["Other_Lowercase"] = reflect.ValueOf(&unicode.Other_Lowercase)
	vars["Other_Math"] = reflect.ValueOf(&unicode.Other_Math)
	vars["Other_Uppercase"] = reflect.ValueOf(&unicode.Other_Uppercase)
	vars["Pattern_Syntax"] = reflect.ValueOf(&unicode.Pattern_Syntax)
	vars["Pattern_White_Space"] = reflect.ValueOf(&unicode.Pattern_White_Space)
	vars["Quotation_Mark"] = reflect.ValueOf(&unicode.Quotation_Mark)
	vars["Radical"] = reflect.ValueOf(&unicode.Radical)
	vars["STerm"] = reflect.ValueOf(&unicode.STerm)
	vars["Soft_Dotted"] = reflect.ValueOf(&unicode.Soft_Dotted)
	vars["Terminal_Punctuation"] = reflect.ValueOf(&unicode.Terminal_Punctuation)
	vars["Unified_Ideograph"] = reflect.ValueOf(&unicode.Unified_Ideograph)
	vars["Variation_Selector"] = reflect.ValueOf(&unicode.Variation_Selector)
	vars["White_Space"] = reflect.ValueOf(&unicode.White_Space)
	vars["CaseRanges"] = reflect.ValueOf(&unicode.CaseRanges)
	vars["FoldCategory"] = reflect.ValueOf(&unicode.FoldCategory)
	vars["FoldScript"] = reflect.ValueOf(&unicode.FoldScript)
	pkgs["unicode"] = &eval.Env {
		Name: "unicode",
		Consts: consts,
		Funcs:  funcs,
		Types:  types,
		Vars:   vars,
		Pkgs:   pkgs,
		Path:   "unicode",
	}
	consts = make(map[string] reflect.Value)
	consts["RuneError"] = reflect.ValueOf(utf8.RuneError)
	consts["RuneSelf"] = reflect.ValueOf(utf8.RuneSelf)
	consts["MaxRune"] = reflect.ValueOf(utf8.MaxRune)
	consts["UTFMax"] = reflect.ValueOf(utf8.UTFMax)

	funcs = make(map[string] reflect.Value)
	funcs["FullRune"] = reflect.ValueOf(utf8.FullRune)
	funcs["FullRuneInString"] = reflect.ValueOf(utf8.FullRuneInString)
	funcs["DecodeRune"] = reflect.ValueOf(utf8.DecodeRune)
	funcs["DecodeRuneInString"] = reflect.ValueOf(utf8.DecodeRuneInString)
	funcs["DecodeLastRune"] = reflect.ValueOf(utf8.DecodeLastRune)
	funcs["DecodeLastRuneInString"] = reflect.ValueOf(utf8.DecodeLastRuneInString)
	funcs["RuneLen"] = reflect.ValueOf(utf8.RuneLen)
	funcs["EncodeRune"] = reflect.ValueOf(utf8.EncodeRune)
	funcs["RuneCount"] = reflect.ValueOf(utf8.RuneCount)
	funcs["RuneCountInString"] = reflect.ValueOf(utf8.RuneCountInString)
	funcs["RuneStart"] = reflect.ValueOf(utf8.RuneStart)
	funcs["Valid"] = reflect.ValueOf(utf8.Valid)
	funcs["ValidString"] = reflect.ValueOf(utf8.ValidString)
	funcs["ValidRune"] = reflect.ValueOf(utf8.ValidRune)

	types = make(map[string] reflect.Type)

	vars = make(map[string] reflect.Value)
	pkgs["utf8"] = &eval.Env {
		Name: "utf8",
		Consts: consts,
		Funcs:  funcs,
		Types:  types,
		Vars:   vars,
		Pkgs:   pkgs,
		Path:   "unicode/utf8",
	}
}
