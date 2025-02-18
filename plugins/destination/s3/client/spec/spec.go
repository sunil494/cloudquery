package spec

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/cloudquery/filetypes/v4"
	"github.com/cloudquery/plugin-sdk/v4/configtype"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const (
	varFormat = "{{FORMAT}}"
	varTable  = "{{TABLE}}"
	varUUID   = "{{UUID}}"
	varYear   = "{{YEAR}}"
	varMonth  = "{{MONTH}}"
	varDay    = "{{DAY}}"
	varHour   = "{{HOUR}}"
	varMinute = "{{MINUTE}}"
)

type Spec struct {
	filetypes.FileSpec

	// Bucket where to sync the files.
	Bucket string `json:"bucket,omitempty" jsonschema:"required,minLength=1"`

	// Region where bucket is located.
	Region string `json:"region,omitempty" jsonschema:"required,minLength=1"`

	//  Path to where the files will be uploaded in the above bucket. The path supports the following placeholder variables:
	//
	// - `{{TABLE}}` will be replaced with the table name
	// - `{{FORMAT}}` will be replaced with the file format, such as `csv`, `json` or `parquet`. If compression is enabled, the format will be `csv.gz`, `json.gz` etc.
	// - `{{UUID}}` will be replaced with a random UUID to uniquely identify each file
	// - `{{YEAR}}` will be replaced with the current year in `YYYY` format
	// - `{{MONTH}}` will be replaced with the current month in `MM` format
	// - `{{DAY}}` will be replaced with the current day in `DD` format
	// - `{{HOUR}}` will be replaced with the current hour in `HH` format
	// - `{{MINUTE}}` will be replaced with the current minute in `mm` format
	//
	// **Note** that timestamps are in `UTC` and will be the current time at the time the file is written, not when the sync started.
	Path string `json:"path,omitempty" jsonschema:"required,pattern=^[^/].*$"` // other cases (//, ./, ../) are covered in extended part

	// If set to `true`, the plugin will write to one file per table.
	// Otherwise, for every batch a new file will be created with a different `.<UUID>` suffix.
	NoRotate bool `json:"no_rotate,omitempty" jsonschema:"default=false"`

	// When `athena` is set to `true`, the S3 plugin will sanitize keys in JSON columns to be compatible with the Hive Metastore / Athena.
	// This allows tables to be created with a Glue Crawler and then queried via Athena, without changes to the table schema.
	Athena bool `json:"athena,omitempty" jsonschema:"default=false"`

	// Ensure write access to the given bucket and path by writing a test object on each sync.
	// If you are sure that the bucket and path are writable, you can set this to `false` to skip the test.
	TestWrite *bool `json:"test_write,omitempty" jsonschema:"default=true"`

	// Endpoint to use for S3 API calls. This is useful for S3-compatible storage services such as MinIO.
	// **Note**: if you want to use path-style addressing, i.e., `https://s3.amazonaws.com/BUCKET/KEY`, `use_path_style` should be enabled, too.
	Endpoint string `json:"endpoint,omitempty"  jsonschema:"default="`

	// SSE KMS Key ID appened to S3 API calls header. Used in conjuction with server_side_encryption.
	SSEKMSKeyId string `json:"sse_kms_key_id,omitempty" jsonschema:"default="`

	// Server Side Encryption header which declares encryption type in S3 API calls header: x-amz-server-side-encryption.
	ServerSideEncryption types.ServerSideEncryption `json:"server_side_encryption,omitempty" jsonschema:"enum=AES256,enum=aws:kms,enum=aws:kms:dsse"`

	// Allows to use path-style addressing in the `endpoint` option, i.e., `https://s3.amazonaws.com/BUCKET/KEY`.
	// By default, the S3 client will use virtual hosted bucket addressing when possible (`https://BUCKET.s3.amazonaws.com/KEY`).
	UsePathStyle bool `json:"use_path_style,omitempty" jsonschema:"default=false"`

	// Disable TLS verification for requests to your S3 endpoint.
	//
	// This option is intended to be used when using a custom endpoint using the `endpoint` option.
	EndpointSkipTLSVerify bool `json:"endpoint_skip_tls_verify,omitempty" jsonschema:"default=false"`

	// Maximum number of items that may be grouped together to be written in a single write.
	//
	// Defaults to `10000` unless `no_rotate` is `true` (will be `0` then).
	BatchSize *int64 `json:"batch_size" jsonschema:"minimum=1,default=10000"`

	// Maximum size of items that may be grouped together to be written in a single write.
	//
	// Defaults to `52428800` (50 MiB) unless `no_rotate` is `true` (will be `0` then).
	BatchSizeBytes *int64 `json:"batch_size_bytes" jsonschema:"minimum=1,default=52428800"`

	// Maximum interval between batch writes.
	//
	// Defaults to `30s` unless `no_rotate` is `true` (will be `0s` then).
	BatchTimeout *configtype.Duration `json:"batch_timeout" jsonschema:"default=30s"`
}

func (s *Spec) SetDefaults() {
	if !strings.Contains(s.Path, varTable) {
		// for backwards-compatibility, default to given path plus /{{TABLE}}.[format].{{UUID}} if
		// no {{TABLE}} value is found in the path string
		s.Path += fmt.Sprintf("/%s.%s", varTable, s.Format)
		if !s.NoRotate {
			s.Path += "." + varUUID
		}
	}
	if s.TestWrite == nil {
		b := true
		s.TestWrite = &b
	}
	if s.BatchSize == nil {
		if s.NoRotate {
			s.BatchSize = ptr(int64(0))
		} else {
			s.BatchSize = ptr(int64(10000))
		}
	}
	if s.BatchSizeBytes == nil {
		if s.NoRotate {
			s.BatchSizeBytes = ptr(int64(0))
		} else {
			s.BatchSizeBytes = ptr(int64(50 * 1024 * 1024)) // 50 MiB
		}
	}
	if s.BatchTimeout == nil {
		if s.NoRotate {
			d := configtype.NewDuration(0)
			s.BatchTimeout = &d
		} else {
			d := configtype.NewDuration(30 * time.Second)
			s.BatchTimeout = &d
		}
	}
}

func (s *Spec) Validate() error {
	if len(s.Bucket) == 0 {
		return fmt.Errorf("`bucket` is required")
	}
	if len(s.Region) == 0 {
		return fmt.Errorf("`region` is required")
	}

	if len(s.Path) == 0 {
		return fmt.Errorf("`path` is required")
	}
	if path.IsAbs(s.Path) {
		return fmt.Errorf("`path` should not start with a \"/\"")
	}
	if s.Path != path.Clean(s.Path) {
		return fmt.Errorf("`path` should not contain relative paths or duplicate slashes")
	}

	if s.NoRotate {
		if strings.Contains(s.Path, varUUID) {
			return fmt.Errorf("`path` should not contain %s when `no_rotate` = true", varUUID)
		}

		if (s.BatchSize != nil && *s.BatchSize > 0) || (s.BatchSizeBytes != nil && *s.BatchSizeBytes > 0) || (s.BatchTimeout != nil && s.BatchTimeout.Duration() > 0) {
			return fmt.Errorf("`no_rotate` cannot be used with non-zero `batch_size`, `batch_size_bytes` or `batch_timeout_ms`")
		}
	}

	if !strings.Contains(s.Path, varUUID) && s.batchingEnabled() {
		return fmt.Errorf("`path` should contain %s when using a non-zero `batch_size`, `batch_size_bytes` or `batch_timeout_ms`", varUUID)
	}

	// required for s.FileSpec.Validate call
	err := s.FileSpec.UnmarshalSpec()
	if err != nil {
		return err
	}
	s.FileSpec.SetDefaults()

	return s.FileSpec.Validate()
}

func (s *Spec) ReplacePathVariables(table string, fileIdentifier string, t time.Time) string {
	name := strings.ReplaceAll(s.Path, varTable, table)
	if strings.Contains(name, varFormat) {
		e := string(s.Format) + s.Compression.Extension()
		name = strings.ReplaceAll(name, varFormat, e)
	}
	name = strings.ReplaceAll(name, varUUID, fileIdentifier)
	name = strings.ReplaceAll(name, varYear, t.Format("2006"))
	name = strings.ReplaceAll(name, varMonth, t.Format("01"))
	name = strings.ReplaceAll(name, varDay, t.Format("02"))
	name = strings.ReplaceAll(name, varHour, t.Format("15"))
	name = strings.ReplaceAll(name, varMinute, t.Format("04"))
	return filepath.Clean(name)
}

func (s *Spec) PathContainsUUID() bool {
	return strings.Contains(s.Path, varUUID)
}

func (s *Spec) batchingEnabled() bool {
	if s.NoRotate {
		// if that's set we don't allow batching
		return false
	}

	return (s.BatchSize == nil || *s.BatchSize > 0) ||
		(s.BatchSizeBytes == nil || *s.BatchSizeBytes > 0) ||
		(s.BatchTimeout == nil || s.BatchTimeout.Duration() > 0)
}

func ptr[A any](a A) *A {
	return &a
}
