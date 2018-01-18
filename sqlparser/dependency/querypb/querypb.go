// Code generated by "pkgimport -p querypb -i gopkg.in/go-grm/sqlparser.v1/dependency/querypb -o querypb.go"; DO NOT EDIT.
// Install by "go get -u -v gopkg.in/pkgimport.v1/cmd/pkgimport";
//go:generate pkgimport -p querypb -i gopkg.in/go-grm/sqlparser.v1/dependency/querypb -o querypb.go

package querypb

import (
	origin "gopkg.in/go-grm/sqlparser.v1/dependency/querypb"
)

type (
	// type

	/*
	 */
	MySqlFlag = origin.MySqlFlag

	/*
	 */
	Flag = origin.Flag

	/*
	 */
	Type = origin.Type

	/*
	 */
	TransactionState = origin.TransactionState

	/*
	 */
	ExecuteOptions_IncludedFields = origin.ExecuteOptions_IncludedFields

	/*
	 */
	ExecuteOptions_Workload = origin.ExecuteOptions_Workload

	/*
	 */
	ExecuteOptions_TransactionIsolation = origin.ExecuteOptions_TransactionIsolation

	/*
	 */
	StreamEvent_Statement_Category = origin.StreamEvent_Statement_Category

	/*
	 */
	SplitQueryRequest_Algorithm = origin.SplitQueryRequest_Algorithm

	/*
	 */
	Value = origin.Value

	/*
	 */
	BindVariable = origin.BindVariable

	/*
	 */
	BoundQuery = origin.BoundQuery
)

var (
	// function

	/*EnumName is a helper function to simplify printing protocol buffer enums
	by name.  Given an enum map and a value, it returns a useful string.

	*/
	EnumName = origin.EnumName
)

var (
	// value

	/*
	 */
	MySqlFlag_EMPTY = origin.MySqlFlag_EMPTY

	/*
	 */
	MySqlFlag_NOT_NULL_FLAG = origin.MySqlFlag_NOT_NULL_FLAG

	/*
	 */
	MySqlFlag_PRI_KEY_FLAG = origin.MySqlFlag_PRI_KEY_FLAG

	/*
	 */
	MySqlFlag_UNIQUE_KEY_FLAG = origin.MySqlFlag_UNIQUE_KEY_FLAG

	/*
	 */
	MySqlFlag_MULTIPLE_KEY_FLAG = origin.MySqlFlag_MULTIPLE_KEY_FLAG

	/*
	 */
	MySqlFlag_BLOB_FLAG = origin.MySqlFlag_BLOB_FLAG

	/*
	 */
	MySqlFlag_UNSIGNED_FLAG = origin.MySqlFlag_UNSIGNED_FLAG

	/*
	 */
	MySqlFlag_ZEROFILL_FLAG = origin.MySqlFlag_ZEROFILL_FLAG

	/*
	 */
	MySqlFlag_BINARY_FLAG = origin.MySqlFlag_BINARY_FLAG

	/*
	 */
	MySqlFlag_ENUM_FLAG = origin.MySqlFlag_ENUM_FLAG

	/*
	 */
	MySqlFlag_AUTO_INCREMENT_FLAG = origin.MySqlFlag_AUTO_INCREMENT_FLAG

	/*
	 */
	MySqlFlag_TIMESTAMP_FLAG = origin.MySqlFlag_TIMESTAMP_FLAG

	/*
	 */
	MySqlFlag_SET_FLAG = origin.MySqlFlag_SET_FLAG

	/*
	 */
	MySqlFlag_NO_DEFAULT_VALUE_FLAG = origin.MySqlFlag_NO_DEFAULT_VALUE_FLAG

	/*
	 */
	MySqlFlag_ON_UPDATE_NOW_FLAG = origin.MySqlFlag_ON_UPDATE_NOW_FLAG

	/*
	 */
	MySqlFlag_NUM_FLAG = origin.MySqlFlag_NUM_FLAG

	/*
	 */
	MySqlFlag_PART_KEY_FLAG = origin.MySqlFlag_PART_KEY_FLAG

	/*
	 */
	MySqlFlag_GROUP_FLAG = origin.MySqlFlag_GROUP_FLAG

	/*
	 */
	MySqlFlag_UNIQUE_FLAG = origin.MySqlFlag_UNIQUE_FLAG

	/*
	 */
	MySqlFlag_BINCMP_FLAG = origin.MySqlFlag_BINCMP_FLAG

	/*
	 */
	MySqlFlag_name = origin.MySqlFlag_name

	/*
	 */
	MySqlFlag_value = origin.MySqlFlag_value

	/*
	 */
	Flag_NONE = origin.Flag_NONE

	/*
	 */
	Flag_ISINTEGRAL = origin.Flag_ISINTEGRAL

	/*
	 */
	Flag_ISUNSIGNED = origin.Flag_ISUNSIGNED

	/*
	 */
	Flag_ISFLOAT = origin.Flag_ISFLOAT

	/*
	 */
	Flag_ISQUOTED = origin.Flag_ISQUOTED

	/*
	 */
	Flag_ISTEXT = origin.Flag_ISTEXT

	/*
	 */
	Flag_ISBINARY = origin.Flag_ISBINARY

	/*
	 */
	Flag_name = origin.Flag_name

	/*
	 */
	Flag_value = origin.Flag_value

	/*NULL_TYPE specifies a NULL type.

	 */
	Type_NULL_TYPE = origin.Type_NULL_TYPE

	/*INT8 specifies a TINYINT type.
	Properties: 1, IsNumber.

	*/
	Type_INT8 = origin.Type_INT8

	/*UINT8 specifies a TINYINT UNSIGNED type.
	Properties: 2, IsNumber, IsUnsigned.

	*/
	Type_UINT8 = origin.Type_UINT8

	/*INT16 specifies a SMALLINT type.
	Properties: 3, IsNumber.

	*/
	Type_INT16 = origin.Type_INT16

	/*UINT16 specifies a SMALLINT UNSIGNED type.
	Properties: 4, IsNumber, IsUnsigned.

	*/
	Type_UINT16 = origin.Type_UINT16

	/*INT24 specifies a MEDIUMINT type.
	Properties: 5, IsNumber.

	*/
	Type_INT24 = origin.Type_INT24

	/*UINT24 specifies a MEDIUMINT UNSIGNED type.
	Properties: 6, IsNumber, IsUnsigned.

	*/
	Type_UINT24 = origin.Type_UINT24

	/*INT32 specifies a INTEGER type.
	Properties: 7, IsNumber.

	*/
	Type_INT32 = origin.Type_INT32

	/*UINT32 specifies a INTEGER UNSIGNED type.
	Properties: 8, IsNumber, IsUnsigned.

	*/
	Type_UINT32 = origin.Type_UINT32

	/*INT64 specifies a BIGINT type.
	Properties: 9, IsNumber.

	*/
	Type_INT64 = origin.Type_INT64

	/*UINT64 specifies a BIGINT UNSIGNED type.
	Properties: 10, IsNumber, IsUnsigned.

	*/
	Type_UINT64 = origin.Type_UINT64

	/*FLOAT32 specifies a FLOAT type.
	Properties: 11, IsFloat.

	*/
	Type_FLOAT32 = origin.Type_FLOAT32

	/*FLOAT64 specifies a DOUBLE or REAL type.
	Properties: 12, IsFloat.

	*/
	Type_FLOAT64 = origin.Type_FLOAT64

	/*TIMESTAMP specifies a TIMESTAMP type.
	Properties: 13, IsQuoted.

	*/
	Type_TIMESTAMP = origin.Type_TIMESTAMP

	/*DATE specifies a DATE type.
	Properties: 14, IsQuoted.

	*/
	Type_DATE = origin.Type_DATE

	/*TIME specifies a TIME type.
	Properties: 15, IsQuoted.

	*/
	Type_TIME = origin.Type_TIME

	/*DATETIME specifies a DATETIME type.
	Properties: 16, IsQuoted.

	*/
	Type_DATETIME = origin.Type_DATETIME

	/*YEAR specifies a YEAR type.
	Properties: 17, IsNumber, IsUnsigned.

	*/
	Type_YEAR = origin.Type_YEAR

	/*DECIMAL specifies a DECIMAL or NUMERIC type.
	Properties: 18, None.

	*/
	Type_DECIMAL = origin.Type_DECIMAL

	/*TEXT specifies a TEXT type.
	Properties: 19, IsQuoted, IsText.

	*/
	Type_TEXT = origin.Type_TEXT

	/*BLOB specifies a BLOB type.
	Properties: 20, IsQuoted, IsBinary.

	*/
	Type_BLOB = origin.Type_BLOB

	/*VARCHAR specifies a VARCHAR type.
	Properties: 21, IsQuoted, IsText.

	*/
	Type_VARCHAR = origin.Type_VARCHAR

	/*VARBINARY specifies a VARBINARY type.
	Properties: 22, IsQuoted, IsBinary.

	*/
	Type_VARBINARY = origin.Type_VARBINARY

	/*CHAR specifies a CHAR type.
	Properties: 23, IsQuoted, IsText.

	*/
	Type_CHAR = origin.Type_CHAR

	/*BINARY specifies a BINARY type.
	Properties: 24, IsQuoted, IsBinary.

	*/
	Type_BINARY = origin.Type_BINARY

	/*BIT specifies a BIT type.
	Properties: 25, IsQuoted.

	*/
	Type_BIT = origin.Type_BIT

	/*ENUM specifies an ENUM type.
	Properties: 26, IsQuoted.

	*/
	Type_ENUM = origin.Type_ENUM

	/*SET specifies a SET type.
	Properties: 27, IsQuoted.

	*/
	Type_SET = origin.Type_SET

	/*TUPLE specifies a a tuple. This cannot
	be returned in a QueryResult, but it can
	be sent as a bind var.
	Properties: 28, None.

	*/
	Type_TUPLE = origin.Type_TUPLE

	/*GEOMETRY specifies a GEOMETRY type.
	Properties: 29, IsQuoted.

	*/
	Type_GEOMETRY = origin.Type_GEOMETRY

	/*JSON specifies a JSON type.
	Properties: 30, IsQuoted.

	*/
	Type_JSON = origin.Type_JSON

	/*EXPRESSION specifies a SQL expression.
	This type is for internal use only.
	Properties: 31, None.

	*/
	Type_EXPRESSION = origin.Type_EXPRESSION

	/*
	 */
	Type_name = origin.Type_name

	/*
	 */
	Type_value = origin.Type_value

	/*
	 */
	TransactionState_UNKNOWN = origin.TransactionState_UNKNOWN

	/*
	 */
	TransactionState_PREPARE = origin.TransactionState_PREPARE

	/*
	 */
	TransactionState_COMMIT = origin.TransactionState_COMMIT

	/*
	 */
	TransactionState_ROLLBACK = origin.TransactionState_ROLLBACK

	/*
	 */
	TransactionState_name = origin.TransactionState_name

	/*
	 */
	TransactionState_value = origin.TransactionState_value

	/*
	 */
	ExecuteOptions_TYPE_AND_NAME = origin.ExecuteOptions_TYPE_AND_NAME

	/*
	 */
	ExecuteOptions_TYPE_ONLY = origin.ExecuteOptions_TYPE_ONLY

	/*
	 */
	ExecuteOptions_ALL = origin.ExecuteOptions_ALL

	/*
	 */
	ExecuteOptions_IncludedFields_name = origin.ExecuteOptions_IncludedFields_name

	/*
	 */
	ExecuteOptions_IncludedFields_value = origin.ExecuteOptions_IncludedFields_value

	/*
	 */
	ExecuteOptions_UNSPECIFIED = origin.ExecuteOptions_UNSPECIFIED

	/*
	 */
	ExecuteOptions_OLTP = origin.ExecuteOptions_OLTP

	/*
	 */
	ExecuteOptions_OLAP = origin.ExecuteOptions_OLAP

	/*
	 */
	ExecuteOptions_DBA = origin.ExecuteOptions_DBA

	/*
	 */
	ExecuteOptions_Workload_name = origin.ExecuteOptions_Workload_name

	/*
	 */
	ExecuteOptions_Workload_value = origin.ExecuteOptions_Workload_value

	/*
	 */
	ExecuteOptions_DEFAULT = origin.ExecuteOptions_DEFAULT

	/*
	 */
	ExecuteOptions_REPEATABLE_READ = origin.ExecuteOptions_REPEATABLE_READ

	/*
	 */
	ExecuteOptions_READ_COMMITTED = origin.ExecuteOptions_READ_COMMITTED

	/*
	 */
	ExecuteOptions_READ_UNCOMMITTED = origin.ExecuteOptions_READ_UNCOMMITTED

	/*
	 */
	ExecuteOptions_SERIALIZABLE = origin.ExecuteOptions_SERIALIZABLE

	/*
	 */
	ExecuteOptions_TransactionIsolation_name = origin.ExecuteOptions_TransactionIsolation_name

	/*
	 */
	ExecuteOptions_TransactionIsolation_value = origin.ExecuteOptions_TransactionIsolation_value

	/*
	 */
	StreamEvent_Statement_Error = origin.StreamEvent_Statement_Error

	/*
	 */
	StreamEvent_Statement_DML = origin.StreamEvent_Statement_DML

	/*
	 */
	StreamEvent_Statement_DDL = origin.StreamEvent_Statement_DDL

	/*
	 */
	StreamEvent_Statement_Category_name = origin.StreamEvent_Statement_Category_name

	/*
	 */
	StreamEvent_Statement_Category_value = origin.StreamEvent_Statement_Category_value

	/*
	 */
	SplitQueryRequest_EQUAL_SPLITS = origin.SplitQueryRequest_EQUAL_SPLITS

	/*
	 */
	SplitQueryRequest_FULL_SCAN = origin.SplitQueryRequest_FULL_SCAN

	/*
	 */
	SplitQueryRequest_Algorithm_name = origin.SplitQueryRequest_Algorithm_name

	/*
	 */
	SplitQueryRequest_Algorithm_value = origin.SplitQueryRequest_Algorithm_value
)