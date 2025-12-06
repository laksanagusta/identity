package nullable

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	duplicateSpaceRegex = regexp.MustCompile(`\s+`)
)

func NewTime(value time.Time) NullTime {
	return NullTime{
		IsExists: true,
		Val:      &value,
	}
}

func NewNilTime() NullTime {
	return NullTime{
		IsExists: true,
		Val:      nil,
	}
}

func NewUUID(value uuid.UUID) uuid.NullUUID {
	return uuid.NullUUID{
		Valid: true,
		UUID:  value,
	}
}

func NewUUIDFromPtr(value *uuid.UUID) uuid.NullUUID {
	if value == nil {
		return uuid.NullUUID{}
	}

	return uuid.NullUUID{
		Valid: true,
		UUID:  *value,
	}
}

func NewString(value string) NullString {
	return NullString{
		IsExists: true,
		Val:      &value,
	}
}

func NewNilString() NullString {
	return NullString{
		IsExists: true,
	}
}

func NewInt64(value int64) NullInt64 {
	return NullInt64{
		IsExists: true,
		Val:      &value,
	}
}

func NewNilInt64() NullInt64 {
	return NullInt64{
		IsExists: true,
	}
}

func NewInt64FromPtr(value *int64) NullInt64 {
	if value == nil {
		return NullInt64{}
	}

	return NullInt64{
		IsExists: true,
		Val:      value,
	}
}

func NewInt32(value int32) NullInt32 {
	return NullInt32{
		IsExists: true,
		Val:      &value,
	}
}

func NewFloat64(value float64) NullFloat64 {
	return NullFloat64{
		IsExists: true,
		Val:      &value,
	}
}

func NewNilFloat64() NullFloat64 {
	return NullFloat64{
		IsExists: true,
	}
}

func NewFloat32(value float32) NullFloat32 {
	return NullFloat32{
		IsExists: true,
		Val:      &value,
	}
}

func NewNilFloat32() NullFloat32 {
	return NullFloat32{
		IsExists: true,
	}
}

func NewBool(value bool) NullBool {
	return NullBool{
		IsExists: true,
		Val:      value,
	}
}

func ToPointer[T comparable](value T) *T {
	return &value
}

type Nullable interface {
	Exists() bool
	Value() (driver.Value, error)
}

type NullString struct {
	Val      *string
	IsExists bool
}

func (n NullString) GetOrDefault(defaultValue ...string) string {
	if n.IsExists {
		return *n.Val
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return string("")
}

func (n NullString) IsNotEmpty() bool {
	return n.IsExists && *n.Val != ""
}

func (n *NullString) Scan(value any) error {
	var sqlValue sql.NullString
	err := sqlValue.Scan(value)
	if err != nil {
		return err
	}
	n.IsExists = sqlValue.Valid
	if n.IsExists {
		n.Val = &sqlValue.String
	}

	return nil
}

func (n NullString) Value() (driver.Value, error) {
	if !n.IsExists || n.Val == nil {
		return "", nil
	}
	return *n.Val, nil
}

func (n NullString) Exists() bool {
	return n.IsExists
}

func (n *NullString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.IsExists = true
		return nil
	}

	err := json.Unmarshal(data, &n.Val)
	if err != nil {
		return err
	}

	trimmedInput := duplicateSpaceRegex.ReplaceAllString(strings.TrimSpace(*n.Val), " ")
	n.Val = &trimmedInput
	n.IsExists = true
	return nil
}

func (n NullString) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.Val)
}

type NullFloat32 struct {
	Val      *float32
	IsExists bool
}

func (n NullFloat32) GetOrDefault(defaultValue ...float32) float32 {
	if n.IsExists {
		return *n.Val
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return float32(0.0)
}

func (n *NullFloat32) Scan(value any) error {
	var sqlValue sql.NullFloat64
	err := sqlValue.Scan(value)
	if err != nil {
		return err
	}

	n.IsExists = sqlValue.Valid
	if n.IsExists {
		float32Val := float32(sqlValue.Float64)
		n.Val = &float32Val
	}

	return nil
}

func (n NullFloat32) Value() (driver.Value, error) {
	if !n.IsExists || n.Val == nil {
		return nil, nil
	}
	return float64(*n.Val), nil
}

func (n NullFloat32) Exists() bool {
	return n.IsExists
}

func (n *NullFloat32) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.IsExists = true
		return nil
	}

	err := json.Unmarshal(data, &n.Val)
	if err != nil {
		return err
	}

	n.IsExists = true
	return nil
}

func (n NullFloat32) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.Val)
}

type NullFloat64 struct {
	Val      *float64
	IsExists bool
}

func (n NullFloat64) GetOrDefault(defaultValue ...float64) float64 {
	if n.IsExists {
		return *n.Val
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return float64(0.0)
}

func (n *NullFloat64) Scan(value any) error {
	var sqlValue sql.NullFloat64
	err := sqlValue.Scan(value)
	if err != nil {
		return err
	}

	n.IsExists = sqlValue.Valid
	if n.IsExists {
		n.Val = &sqlValue.Float64
	}

	return nil
}

func (n NullFloat64) Value() (driver.Value, error) {
	if !n.IsExists || n.Val == nil {
		return nil, nil
	}
	return *n.Val, nil
}

func (n NullFloat64) Exists() bool {
	return n.IsExists
}

func (n *NullFloat64) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.IsExists = true
		return nil
	}

	err := json.Unmarshal(data, &n.Val)
	if err != nil {
		return err
	}

	n.IsExists = true
	return nil
}

func (n NullFloat64) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.Val)
}

type NullInt32 struct {
	Val      *int32
	IsExists bool
}

func (n NullInt32) GetOrDefault(defaultValue ...int32) int32 {
	if n.IsExists {
		return *n.Val
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return int32(0)
}

func (n NullInt32) Exists() bool {
	return n.IsExists
}

func (n *NullInt32) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.IsExists = true
		return nil
	}
	err := json.Unmarshal(data, &n.Val)
	if err != nil {
		return err
	}
	n.IsExists = true
	return nil

}

func (n *NullInt32) Scan(value any) error {

	var sqlValue sql.NullInt32
	err := sqlValue.Scan(value)
	if err != nil {
		return err
	}
	n.IsExists = sqlValue.Valid
	if n.IsExists {
		n.Val = &sqlValue.Int32
	}

	return nil
}

func (n NullInt32) Value() (driver.Value, error) {
	if !n.IsExists || n.Val == nil {
		return nil, nil
	}
	return int64(*n.Val), nil
}

func (n NullInt32) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.Val)
}

type NullInt64 struct {
	Val      *int64
	IsExists bool
}

func (n NullInt64) GetOrDefault(defaultValue ...int64) int64 {
	if n.IsExists {
		return *n.Val
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return int64(0)
}

func (n *NullInt64) Scan(value any) error {

	var sqlValue sql.NullInt64
	err := sqlValue.Scan(value)
	if err != nil {
		return err
	}
	n.IsExists = sqlValue.Valid
	if n.IsExists {
		n.Val = &sqlValue.Int64
	}

	return nil
}

func (n NullInt64) Value() (driver.Value, error) {
	if !n.IsExists || n.Val == nil {
		return nil, nil
	}
	return *n.Val, nil
}

func (n NullInt64) Exists() bool {
	return n.IsExists
}
func (n *NullInt64) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.IsExists = true
		return nil
	}
	err := json.Unmarshal(data, &n.Val)
	if err != nil {
		return err
	}
	n.IsExists = true
	return nil

}

func (n NullInt64) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.Val)
}

type NullBool struct {
	Val      bool
	IsExists bool
}

func (n NullBool) GetOrDefault() bool {
	return n.Val
}

func (n *NullBool) Scan(value any) error {
	var sqlValue sql.NullBool
	err := sqlValue.Scan(value)
	if err != nil {
		return err
	}
	n.IsExists = sqlValue.Valid
	n.Val = sqlValue.Bool
	return nil
}

func (n NullBool) Value() (driver.Value, error) {
	if !n.IsExists {
		return false, nil
	}
	return n.Val, nil
}

func (n NullBool) Exists() bool {
	return n.IsExists
}

func (n *NullBool) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.IsExists = true
		return nil
	}
	err := json.Unmarshal(data, &n.Val)
	if err != nil {
		return err
	}
	n.IsExists = true
	return nil

}

func (n NullBool) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.Val)
}

type NullTime struct {
	Val      *time.Time
	IsExists bool
}

func (n NullTime) GetOrDefault(defaultValue ...time.Time) time.Time {
	if n.IsExists {
		return *n.Val
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return time.Time{}
}

func (t *NullTime) Exists() bool {
	return t.IsExists
}

func (t *NullTime) Scan(value any) error {
	var sqlTime sql.NullTime
	err := sqlTime.Scan(value)
	if err != nil {
		return err
	}

	t.IsExists = sqlTime.Valid
	if t.IsExists {
		t.Val = &sqlTime.Time
	}

	return nil
}

func (t NullTime) Value() (driver.Value, error) {
	if !t.IsExists || t.Val == nil {
		return nil, nil
	}

	return *t.Val, nil
}

func (t *NullTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		t.IsExists = true
		return nil
	}
	err := json.Unmarshal(data, &t.Val)
	if err != nil {
		return err
	}
	t.IsExists = true
	return nil

}
func (n NullTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.Val)
}
