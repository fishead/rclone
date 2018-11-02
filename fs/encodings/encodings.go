// +build !noencode

package encodings

import (
	"runtime"
	"strings"

	"github.com/rclone/rclone/lib/encoder"
)

// Base only encodes the zero byte and slash
const Base = encoder.MultiEncoder(
	encoder.EncodeZero |
		encoder.EncodeSlash)

// Display is the internal encoding for logging and output
const Display = encoder.Standard

// LocalUnix is the encoding used by the local backend for non windows platforms
const LocalUnix = Base

// LocalWindows is the encoding used by the local backend for windows platforms
//
// List of replaced characters:
//   < (less than)     -> '＜' // FULLWIDTH LESS-THAN SIGN
//   > (greater than)  -> '＞' // FULLWIDTH GREATER-THAN SIGN
//   : (colon)         -> '：' // FULLWIDTH COLON
//   " (double quote)  -> '＂' // FULLWIDTH QUOTATION MARK
//   \ (backslash)     -> '＼' // FULLWIDTH REVERSE SOLIDUS
//   | (vertical line) -> '｜' // FULLWIDTH VERTICAL LINE
//   ? (question mark) -> '？' // FULLWIDTH QUESTION MARK
//   * (asterisk)      -> '＊' // FULLWIDTH ASTERISK
//
// Additionally names can't end with a period (.) or space ( ).
// List of replaced characters:
//   . (period)        -> '．' // FULLWIDTH FULL STOP
//     (space)         -> '␠' // SYMBOL FOR SPACE
//
// Also encode invalid UTF-8 bytes as Go can't convert them to UTF-16.
//
// https://docs.microsoft.com/de-de/windows/desktop/FileIO/naming-a-file#naming-conventions
const LocalWindows = encoder.MultiEncoder(
	uint(Base) |
		encoder.EncodeWin |
		encoder.EncodeBackSlash |
		encoder.EncodeCtl |
		encoder.EncodeRightSpace |
		encoder.EncodeRightPeriod |
		encoder.EncodeInvalidUtf8)

// AmazonCloudDrive is the encoding used by the amazonclouddrive backend
//
// Encode invalid UTF-8 bytes as json doesn't handle them properly.
const AmazonCloudDrive = encoder.MultiEncoder(
	uint(Base) |
		encoder.EncodeInvalidUtf8)

// B2 is the encoding used by the b2 backend
//
// See: https://www.backblaze.com/b2/docs/files.html
// Encode invalid UTF-8 bytes as json doesn't handle them properly.
// FIXME: allow /, but not leading, trailing or double
const B2 = encoder.MultiEncoder(
	uint(Display) |
		encoder.EncodeBackSlash |
		encoder.EncodeInvalidUtf8)

// Box is the encoding used by the box backend
//
// From https://developer.box.com/docs/error-codes#section-400-bad-request :
// > Box only supports file or folder names that are 255 characters or less.
// > File names containing non-printable ascii, "/" or "\", names with leading
// > or trailing spaces, and the special names “.” and “..” are also unsupported.
//
// Testing revealed names with leading spaces work fine.
// Also encode invalid UTF-8 bytes as json doesn't handle them properly.
const Box = encoder.MultiEncoder(
	uint(Display) |
		encoder.EncodeBackSlash |
		encoder.EncodeRightSpace |
		encoder.EncodeInvalidUtf8)

// Drive is the encoding used by the drive backend
//
// Encode invalid UTF-8 bytes as json doesn't handle them.
// Don't encode / as it's a valid name character in drive.
const Drive = encoder.MultiEncoder(
	encoder.EncodeInvalidUtf8)

// Dropbox is the encoding used by the dropbox backend
//
// https://www.dropbox.com/help/syncing-uploads/files-not-syncing lists / and \
// as invalid characters.
// Testing revealed names with trailing spaces and the DEL character don't work.
// Also encode invalid UTF-8 bytes as json doesn't handle them properly.
const Dropbox = encoder.MultiEncoder(
	uint(Base) |
		encoder.EncodeBackSlash |
		encoder.EncodeDel |
		encoder.EncodeRightSpace |
		encoder.EncodeInvalidUtf8)

// GoogleCloudStorage is the encoding used by the googlecloudstorage backend
const GoogleCloudStorage = encoder.MultiEncoder(
	uint(Base) |
		//encoder.EncodeCrLF |
		encoder.EncodeInvalidUtf8)

// JottaCloud is the encoding used by the jottacloud backend
//
// Encode invalid UTF-8 bytes as xml doesn't handle them properly.
const JottaCloud = encoder.MultiEncoder(
	uint(Display) |
		encoder.EncodeInvalidUtf8)

// Koofr is the encoding used by the koofr backend
//
// Encode invalid UTF-8 bytes as json doesn't handle them properly.
const Koofr = encoder.MultiEncoder(
	uint(Display) |
		encoder.EncodeBackSlash |
		encoder.EncodeInvalidUtf8)

// Mega is the encoding used by the mega backend
//
// Encode invalid UTF-8 bytes as json doesn't handle them properly.
const Mega = encoder.MultiEncoder(
	uint(Base) |
		encoder.EncodeInvalidUtf8)

// OneDrive is the encoding used by the onedrive backend
//
// List of replaced characters:
//   < (less than)     -> '＜' // FULLWIDTH LESS-THAN SIGN
//   > (greater than)  -> '＞' // FULLWIDTH GREATER-THAN SIGN
//   : (colon)         -> '：' // FULLWIDTH COLON
//   " (double quote)  -> '＂' // FULLWIDTH QUOTATION MARK
//   \ (backslash)     -> '＼' // FULLWIDTH REVERSE SOLIDUS
//   | (vertical line) -> '｜' // FULLWIDTH VERTICAL LINE
//   ? (question mark) -> '？' // FULLWIDTH QUESTION MARK
//   * (asterisk)      -> '＊' // FULLWIDTH ASTERISK
//   # (number sign)  -> '＃'  // FULLWIDTH NUMBER SIGN
//   % (percent sign) -> '％'  // FULLWIDTH PERCENT SIGN
//
// Folder names cannot begin with a tilde ('~')
// List of replaced characters:
//   ~ (tilde)        -> '～'  // FULLWIDTH TILDE
//
// Additionally names can't begin with a space ( ) or end with a period (.) or space ( ).
// List of replaced characters:
//   . (period)        -> '．' // FULLWIDTH FULL STOP
//     (space)         -> '␠'  // SYMBOL FOR SPACE
//
// Also encode invalid UTF-8 bytes as json doesn't handle them.
//
// The OneDrive API documentation lists the set of reserved characters, but
// testing showed this list is incomplete. This are the differences:
//  - " (double quote) is rejected, but missing in the documentation
//  - space at the end of file and folder names is rejected, but missing in the documentation
//  - period at the end of file names is rejected, but missing in the documentation
//
// Adding these restrictions to the OneDrive API documentation yields exactly
// the same rules as the Windows naming conventions.
//
// https://docs.microsoft.com/en-us/onedrive/developer/rest-api/concepts/addressing-driveitems?view=odsp-graph-online#path-encoding
const OneDrive = encoder.MultiEncoder(
	uint(Display) |
		encoder.EncodeBackSlash |
		encoder.EncodeHashPercent |
		encoder.EncodeLeftSpace |
		encoder.EncodeLeftTilde |
		encoder.EncodeRightPeriod |
		encoder.EncodeRightSpace |
		encoder.EncodeWin |
		encoder.EncodeInvalidUtf8)

// OpenDrive is the encoding used by the opendrive backend
//
// List of replaced characters:
//   < (less than)     -> '＜' // FULLWIDTH LESS-THAN SIGN
//   > (greater than)  -> '＞' // FULLWIDTH GREATER-THAN SIGN
//   : (colon)         -> '：' // FULLWIDTH COLON
//   " (double quote)  -> '＂' // FULLWIDTH QUOTATION MARK
//   \ (backslash)     -> '＼' // FULLWIDTH REVERSE SOLIDUS
//   | (vertical line) -> '｜' // FULLWIDTH VERTICAL LINE
//   ? (question mark) -> '？' // FULLWIDTH QUESTION MARK
//   * (asterisk)      -> '＊' // FULLWIDTH ASTERISK
//
// Additionally names can't begin or end with a ASCII whitespace.
// List of replaced characters:
//     (space)           -> '␠'  // SYMBOL FOR SPACE
//     (horizontal tab)  -> '␉'  // SYMBOL FOR HORIZONTAL TABULATION
//     (line feed)       -> '␊'  // SYMBOL FOR LINE FEED
//     (vertical tab)    -> '␋'  // SYMBOL FOR VERTICAL TABULATION
//     (carriage return) -> '␍'  // SYMBOL FOR CARRIAGE RETURN
//
// Also encode invalid UTF-8 bytes as json doesn't handle them properly.
//
// https://www.opendrive.com/wp-content/uploads/guides/OpenDrive_API_guide.pdf
const OpenDrive = encoder.MultiEncoder(
	uint(Base) |
		encoder.EncodeWin |
		encoder.EncodeLeftCrLfHtVt |
		encoder.EncodeRightCrLfHtVt |
		encoder.EncodeBackSlash |
		encoder.EncodeLeftSpace |
		encoder.EncodeRightSpace |
		encoder.EncodeInvalidUtf8)

// Pcloud is the encoding used by the pcloud backend
//
// Encode invalid UTF-8 bytes as json doesn't handle them properly.
//
// TODO: Investigate Unicode simplification (＼ gets converted to \ server-side)
const Pcloud = encoder.MultiEncoder(
	uint(Base) |
		encoder.EncodeInvalidUtf8)

// ByName returns the encoder for a give backend name or nil
func ByName(name string) encoder.Encoder {
	switch strings.ToLower(name) {
	case "base":
		return Base
	case "display":
		return Display
	case "amazonclouddrive":
		return AmazonCloudDrive
	//case "azureblob":
	case "b2":
		return B2
	case "box":
		return Box
	//case "cache":
	case "drive":
		return Drive
	case "dropbox":
		return Dropbox
	//case "ftp":
	case "googlecloudstorage":
		return GoogleCloudStorage
	//case "http":
	//case "hubic":
	case "jottacloud":
		return JottaCloud
	case "koofr":
		return Koofr
	case "local":
		return Local()
	case "local-windows", "windows":
		return LocalWindows
	case "local-unix", "unix":
		return LocalUnix
	case "mega":
		return Mega
	case "onedrive":
		return OneDrive
	case "opendrive":
		return OpenDrive
	case "pcloud":
		return Pcloud
	//case "qingstor":
	//case "s3":
	//case "sftp":
	//case "swift":
	//case "webdav":
	//case "yandex":
	default:
		return nil
	}
}

// Local returns the local encoding for the current platform
func Local() encoder.MultiEncoder {
	if runtime.GOOS == "windows" {
		return LocalWindows
	}
	return LocalUnix
}

// Names returns the list of known encodings as accepted by ByName
func Names() []string {
	return []string{
		"base",
		"display",
		"amazonclouddrive",
		"b2",
		"box",
		"drive",
		"dropbox",
		"googlecloudstorage",
		"jottacloud",
		"koofr",
		"local-unix",
		"local-windows",
		"local",
		"mega",
		"onedrive",
		"opendrive",
		"pcloud",
	}
}