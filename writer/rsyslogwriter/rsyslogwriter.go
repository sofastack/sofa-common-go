// nolint
// Copyright 20xx The Alipay Authors.
//
// @authors[0]: bingwu.ybw(bingwu.ybw@antfin.com|detailyang@gmail.com)
// @authors[1]: robotx(robotx@antfin.com)
//
// *Legal Disclaimer*
// Within this source code, the comments in Chinese shall be the original, governing version. Any comment in other languages are for reference only. In the event of any conflict between the Chinese language version comments and other language version comments, the Chinese language version shall prevail.
// *法律免责声明*
// 关于代码注释部分，中文注释为官方版本，其它语言注释仅做参考。中文注释可能与其它语言注释存在不一致，当中文注释与其它语言注释存在不一致时，请以中文注释为准。
//
//

package rsyslogwriter

// rsyslogwriter package implements rfc5424

// SYSLOG-MSG      = HEADER SP STRUCTURED-DATA [SP MSG]

// HEADER          = PRI VERSION SP TIMESTAMP SP HOSTNAME
// 				  SP APP-NAME SP PROCID SP MSGID
// PRI             = "<" PRIVAL ">"
// PRIVAL          = 1*3DIGIT ; range 0 .. 191
// VERSION         = NONZERO-DIGIT 0*2DIGIT
// HOSTNAME        = NILVALUE / 1*255PRINTUSASCII

// APP-NAME        = NILVALUE / 1*48PRINTUSASCII
// PROCID          = NILVALUE / 1*128PRINTUSASCII
// MSGID           = NILVALUE / 1*32PRINTUSASCII

// TIMESTAMP       = NILVALUE / FULL-DATE "T" FULL-TIME
// FULL-DATE       = DATE-FULLYEAR "-" DATE-MONTH "-" DATE-MDAY
// DATE-FULLYEAR   = 4DIGIT
// DATE-MONTH      = 2DIGIT  ; 01-12
// DATE-MDAY       = 2DIGIT  ; 01-28, 01-29, 01-30, 01-31 based on
// 						  ; month/year
// FULL-TIME       = PARTIAL-TIME TIME-OFFSET
// PARTIAL-TIME    = TIME-HOUR ":" TIME-MINUTE ":" TIME-SECOND
// 				  [TIME-SECFRAC]
// TIME-HOUR       = 2DIGIT  ; 00-23
// TIME-MINUTE     = 2DIGIT  ; 00-59
// TIME-SECOND     = 2DIGIT  ; 00-59
// TIME-SECFRAC    = "." 1*6DIGIT
// TIME-OFFSET     = "Z" / TIME-NUMOFFSET
// TIME-NUMOFFSET  = ("+" / "-") TIME-HOUR ":" TIME-MINUTE

// STRUCTURED-DATA = NILVALUE / 1*SD-ELEMENT
// SD-ELEMENT      = "[" SD-ID *(SP SD-PARAM) "]"
// SD-PARAM        = PARAM-NAME "=" %d34 PARAM-VALUE %d34
// SD-ID           = SD-NAME
// PARAM-NAME      = SD-NAME
// PARAM-VALUE     = UTF-8-STRING ; characters '"', '\' and
// 							   ; ']' MUST be escaped.
// SD-NAME         = 1*32PRINTUSASCII
// 				  ; except '=', SP, ']', %d34 (")

// MSG             = MSG-ANY / MSG-UTF8
// MSG-ANY         = *OCTET ; not starting with BOM
// MSG-UTF8        = BOM UTF-8-STRING
// BOM             = %xEF.BB.BF

import (
	"bytes"
	"errors"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var hostname string

func init() {
	hostname, _ = os.Hostname()
}

type Severity uint8

const (
	EMERG   Severity = 0 // system is unusable
	ALERT   Severity = 1 // action must be taken immediately
	CRIT    Severity = 2 // critical conditions
	ERR     Severity = 3 // error conditions
	WARNING Severity = 4 // warning conditions
	NOTICE  Severity = 5 // normal but significant condition
	INFO    Severity = 6 // informational
	DEBUG   Severity = 7 // debug-level messages
)

func ParseSeverity(s string) (Severity, error) {
	switch strings.ToUpper(s) {
	case "EMERG":
		return EMERG, nil

	case "ALERT":
		return ALERT, nil

	case "CRIT":
		return CRIT, nil

	case "ERR":
		return ERR, nil

	case "WARNING":
		return WARNING, nil

	case "NOTICE":
		return NOTICE, nil

	case "INFO":
		return INFO, nil

	case "DEBUG":
		return DEBUG, nil
	default:
		return 0, errors.New("unknown severity")
	}
}

type Facility uint8

const (
	KERN     Facility = 0  // kernel messages
	USER     Facility = 1  // random user-level messages
	MAIL     Facility = 2  // mail system
	DAEMON   Facility = 3  // system daemons
	AUTH     Facility = 4  // security/authorization messages
	SYSLOG   Facility = 5  // messages generated internally by syslogd
	LPR      Facility = 6  // line printer subsystem
	NEWS     Facility = 7  // network news subsystem
	UUCP     Facility = 8  // UUCP subsystem
	CRON     Facility = 9  // clock daemon
	AUTHPRIV Facility = 10 // security/authorization messages (private)
	FTP      Facility = 11 // FTP daemon
	LOCAL0   Facility = 16 // reserved for local use
	LOCAL1   Facility = 17 // reserved for local use
	LOCAL2   Facility = 18 // reserved for local use
	LOCAL3   Facility = 19 // reserved for local use
	LOCAL4   Facility = 20 // reserved for local use
	LOCAL5   Facility = 21 // reserved for local use
	LOCAL6   Facility = 22 // reserved for local use
	LOCAL7   Facility = 23 // reserved for local use
)

func ParseFacility(s string) (Facility, error) {
	switch strings.ToUpper(s) {
	case "KERN":
		return KERN, nil
	case "USER":
		return USER, nil
	case "MAIL":
		return MAIL, nil
	case "DAEMON":
		return DAEMON, nil
	case "AUTH":
		return AUTH, nil
	case "SYSLOG":
		return SYSLOG, nil
	case "LPR":
		return LPR, nil
	case "NEWS":
		return NEWS, nil
	case "UUCP":
		return UUCP, nil
	case "CRON":
		return CRON, nil
	case "AUTHPRIV":
		return AUTHPRIV, nil
	case "FTP":
		return FTP, nil
	case "LOCAL0":
		return LOCAL0, nil
	case "LOCAL1":
		return LOCAL1, nil
	case "LOCAL2":
		return LOCAL2, nil
	case "LOCAL3":
		return LOCAL3, nil
	case "LOCAL4":
		return LOCAL4, nil
	case "LOCAL5":
		return LOCAL5, nil
	case "LOCAL6":
		return LOCAL6, nil
	case "LOCAL7":
		return LOCAL7, nil
	default:
		return 0, errors.New("unknown facility")
	}
}

type RsyslogWriter struct {
	option *Option
	pri    string
	pid    string
	buffer bytes.Buffer
	conn   *net.UDPConn
}

type Option struct {
	server   string
	hostname string
	appname  string
	severity Severity
	facility Facility
}

func NewOption() *Option {
	return &Option{
		hostname: hostname,
	}
}

func (o *Option) SetServer(s string) *Option     { o.server = s; return o }
func (o *Option) SetHostname(s string) *Option   { o.hostname = s; return o }
func (o *Option) SetAppname(s string) *Option    { o.appname = s; return o }
func (o *Option) SetSeverity(s Severity) *Option { o.severity = s; return o }
func (o *Option) SetFacility(s Facility) *Option { o.facility = s; return o }

func New(o *Option) (*RsyslogWriter, error) {
	dstAddr, err := net.ResolveUDPAddr("udp4", o.server)
	if err != nil {
		return nil, err
	}

	pri := strconv.Itoa(int(o.facility*8) + int(o.severity))

	// Let the kernel choose a source port
	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}

	// Allocate a socket and set the src and dst address
	conn, err := net.DialUDP("udp4", srcAddr, dstAddr)
	if err != nil {
		return nil, err
	}

	return &RsyslogWriter{
		pri:    pri,
		pid:    strconv.Itoa(os.Getpid()),
		option: o,
		conn:   conn,
	}, nil
}

// nolint
func (rs *RsyslogWriter) Write(p []byte) (int, error) {
	// RFC5424
	// <165>1 2003-08-24T05:14:15.000003-07:00 192.0.2.1 appname - - It's time to make the do-nuts
	b := rs.buffer

	b.Reset()
	b.WriteString("<")
	b.WriteString(rs.pri)
	b.WriteString(">1 ")
	b.WriteString(time.Now().Format(time.RFC3339))
	b.WriteString(" ")
	b.WriteString(rs.option.hostname)
	b.WriteString(" ")
	b.WriteString(rs.option.appname)
	b.WriteString(" ")
	b.WriteString(rs.pid)
	b.WriteString(" -") // No msg id
	b.WriteString(" ")
	b.Write(p)

	return rs.conn.Write(b.Bytes())
}

func (rs *RsyslogWriter) Close() error {
	return rs.conn.Close()
}
