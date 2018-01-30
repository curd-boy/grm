




{{define "GetSchema"}}

-- @Type Select
-- @Comm "获取 当前数据库名"
-- @Resp TableSchema         string "获取 当前数据库名"

SELECT DATABASE() AS table_schema;

{{end}}


{{define "GetTable"}}

-- @Type  Select                []
-- @Comm  "获取 指定库的所有表"
-- @Count
-- @Req   TableSchema           string 数据库名
-- @Resp  TableName             string 所有表

SELECT 
    table_name
FROM
    information_schema.columns
WHERE
    table_schema = :TableSchema
GROUP BY table_name;

{{end}}


{{define "GetColumn"}}

-- @Type Select                  []
-- @Comm "获取 指定库表的结构体"
-- @Req  TableSchema             string 数据库名
-- @Req  TableName               string 表名
-- @Resp TableCatalog            string 表目录
-- @Resp TableSchema             string 数据库名
-- @Resp TableName               string 表名
-- @Resp ColumnName              string 字段名
-- @Resp OrdinalPosition         int    字段所在位置
-- @Resp ColumnDefault           string 字段默认值
-- @Resp IsNullable              string 能是空值
-- @Resp DataType                string 数据类型
-- @Resp CharacterMaximumLength  int    字符串最大长度
-- @Resp CharacterOctetLength    int    字符串位长度
-- @Resp NumericPrecision        int    数字精度
-- @Resp NumericScale            string 数字比例
-- @Resp DatetimePrecision       int    时间精度
-- @Resp CharacterSetName        string 字符编码
-- @Resp CollationName           string 字符编码
-- @Resp ColumnType              string 列类型
-- @Resp ColumnKey               string 列键
-- @Resp Extra                   string 额外属性
-- @Resp Privileges              string 特供
-- @Resp ColumnComment           string 注释

SELECT 
    table_catalog,
    table_schema,
    table_name,
    column_name,
    ordinal_position,
    column_default,
    is_nullable,
    data_type,
    character_maximum_length,
    character_octet_length,
    numeric_precision,
    numeric_scale,
    datetime_precision,
    character_set_name,
    collation_name,
    column_type,
    column_key,
    extra,
    privileges,
    column_comment
FROM
    information_schema.columns
WHERE
    table_schema = :TableSchema
        AND table_name = :TableName
ORDER BY ordinal_position;
{{end}}


{{define "GetCreateTable"}}
-- @Type Select
-- @Comm "mysql 获取 指定库表的结构体"
-- @Req  TableSchema                   string 数据库名
-- @Req  TableName                     string 表名
-- @Resp SqlCreateTable                string 创建表的sql

SELECT 
    CONCAT(
            'CREATE TABLE IF NOT EXISTS ',
            '`',
            table_name,
            '` (',
            GROUP_CONCAT('\n    ',
                CONCAT('`', column_name, '`'),
                ' ',
                column_type,
                ' ',
                IF(is_nullable = 'NO', 'NOT NULL', 'NULL'),
                IF(column_default IS NULL OR column_default = '', '', CONCAT(' DEFAULT \'', column_default, '\'')),
                IF(extra IS NULL OR extra = '', '', CONCAT(' ', extra)),
                IF(column_comment IS NULL OR column_comment = '', '', CONCAT(' COMMENT \'', column_comment, '\''))
            ),
            ',',
            (SELECT 
                GROUP_CONCAT('\n    ', key_column.key_column)
            FROM
                (SELECT 
                    CONCAT(
                        IF(constraint_name = 'PRIMARY', 'PRIMARY KEY', CONCAT('UNIQUE KEY `', constraint_name, '`')), 
                        ' (', 
                        GROUP_CONCAT(' `', column_name, '`'), 
                        ' )'
                    ) key_column
                FROM
                    information_schema.key_column_usage
                WHERE
                    table_name = 'a'
                GROUP BY constraint_name) AS key_column
            ),
            '\n)',
            (SELECT 
                CONCAT(
                    ' ENGINE=',
                    engine,
                    ' DEFAULT CHARSET=',
                    (SELECT 
                        character_set_name
                    FROM
                        information_schema.collations
                    WHERE
                        collation_name = table_collation
                    ),
                    ' COMMENT=\'',
                    table_comment,
                    '\'' 
                )
            FROM
                information_schema.tables
            WHERE
                table_schema = :TableSchema AND table_name = :TableName
            ),
            ';'
    ) AS sql_create_table
FROM
    information_schema.columns
WHERE
    table_schema = :TableSchema AND table_name = :TableName;
{{end}}

