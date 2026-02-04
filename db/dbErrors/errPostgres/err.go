package errPostgres

import "github.com/helays/utils/v2/db/dbErrors"

var PgErrorMap = map[string]dbErrors.DbError{
	// Class 00 — Successful Completion
	"00000": {Code: "00000", EN: "successful_completion", ZH: "成功完成", Class: "Successful Completion"},

	// Class 01 — Warning
	"01000": {Code: "01000", EN: "warning", ZH: "警告", Class: "Warning"},
	"0100C": {Code: "0100C", EN: "dynamic_result_sets_returned", ZH: "返回动态结果集", Class: "Warning"},
	"01008": {Code: "01008", EN: "implicit_zero_bit_padding", ZH: "隐式零位填充", Class: "Warning"},
	"01003": {Code: "01003", EN: "null_value_eliminated_in_set_function", ZH: "在集合函数中消除了空值", Class: "Warning"},
	"01007": {Code: "01007", EN: "privilege_not_granted", ZH: "权限未授予", Class: "Warning"},
	"01006": {Code: "01006", EN: "privilege_not_revoked", ZH: "权限未撤销", Class: "Warning"},
	"01004": {Code: "01004", EN: "string_data_right_truncation", ZH: "字符串数据右截断", Class: "Warning"},
	"01P01": {Code: "01P01", EN: "deprecated_feature", ZH: "已弃用的功能", Class: "Warning"},

	// Class 02 — No Data
	"02000": {Code: "02000", EN: "no_data", ZH: "无数据", Class: "No Data"},
	"02001": {Code: "02001", EN: "no_additional_dynamic_result_sets_returned", ZH: "未返回额外的动态结果集", Class: "No Data"},

	// Class 03 — SQL Statement Not Yet Complete
	"03000": {Code: "03000", EN: "sql_statement_not_yet_complete", ZH: "SQL 语句尚未完成", Class: "SQL Statement Not Yet Complete"},

	// Class 08 — Connection Exception
	"08000": {Code: "08000", EN: "connection_exception", ZH: "连接异常", Class: "Connection Exception"},
	"08003": {Code: "08003", EN: "connection_does_not_exist", ZH: "连接不存在", Class: "Connection Exception"},
	"08006": {Code: "08006", EN: "connection_failure", ZH: "连接失败", Class: "Connection Exception"},
	"08001": {Code: "08001", EN: "sqlclient_unable_to_establish_sqlconnection", ZH: "SQL 客户端无法建立连接", Class: "Connection Exception"},
	"08004": {Code: "08004", EN: "sqlserver_rejected_establishment_of_sqlconnection", ZH: "SQL 服务器拒绝建立连接", Class: "Connection Exception"},
	"08007": {Code: "08007", EN: "transaction_resolution_unknown", ZH: "事务解析未知", Class: "Connection Exception"},
	"08P01": {Code: "08P01", EN: "protocol_violation", ZH: "协议违规", Class: "Connection Exception"},

	// Class 09 — Triggered Action Exception
	"09000": {Code: "09000", EN: "triggered_action_exception", ZH: "触发动作异常", Class: "Triggered Action Exception"},

	// Class 0A — Feature Not Supported
	"0A000": {Code: "0A000", EN: "feature_not_supported", ZH: "不支持的功能", Class: "Feature Not Supported"},

	// Class 0B — Invalid Transaction Initiation
	"0B000": {Code: "0B000", EN: "invalid_transaction_initiation", ZH: "无效的事务启动", Class: "Invalid Transaction Initiation"},

	// Class 0F — Locator Exception
	"0F000": {Code: "0F000", EN: "locator_exception", ZH: "定位器异常", Class: "Locator Exception"},
	"0F001": {Code: "0F001", EN: "invalid_locator_specification", ZH: "无效的定位器规范", Class: "Locator Exception"},

	// Class 0L — Invalid Grantor
	"0L000": {Code: "0L000", EN: "invalid_grantor", ZH: "无效的授权者", Class: "Invalid Grantor"},
	"0LP01": {Code: "0LP01", EN: "invalid_grant_operation", ZH: "无效的授权操作", Class: "Invalid Grantor"},

	// Class 0P — Invalid Role Specification
	"0P000": {Code: "0P000", EN: "invalid_role_specification", ZH: "无效的角色规范", Class: "Invalid Role Specification"},

	// Class 0Z — Diagnostics Exception
	"0Z000": {Code: "0Z000", EN: "diagnostics_exception", ZH: "诊断异常", Class: "Diagnostics Exception"},
	"0Z002": {Code: "0Z002", EN: "stacked_diagnostics_accessed_without_active_handler", ZH: "在没有活动处理程序的情况下访问堆叠诊断", Class: "Diagnostics Exception"},

	// Class 20 — Case Not Found
	"20000": {Code: "20000", EN: "case_not_found", ZH: "未找到案例", Class: "Case Not Found"},

	// Class 21 — Cardinality Violation
	"21000": {Code: "21000", EN: "cardinality_violation", ZH: "基数违规", Class: "Cardinality Violation"},

	// Class 22 — Data Exception
	"22000": {Code: "22000", EN: "data_exception", ZH: "数据异常", Class: "Data Exception"},
	"2202E": {Code: "2202E", EN: "array_subscript_error", ZH: "数组下标错误", Class: "Data Exception"},
	"22021": {Code: "22021", EN: "character_not_in_repertoire", ZH: "字符不在字符集中", Class: "Data Exception"},
	"22008": {Code: "22008", EN: "datetime_field_overflow", ZH: "日期时间字段溢出", Class: "Data Exception"},
	"22012": {Code: "22012", EN: "division_by_zero", ZH: "除零错误", Class: "Data Exception"},
	"22005": {Code: "22005", EN: "error_in_assignment", ZH: "赋值错误", Class: "Data Exception"},
	"2200B": {Code: "2200B", EN: "escape_character_conflict", ZH: "转义字符冲突", Class: "Data Exception"},
	"22022": {Code: "22022", EN: "indicator_overflow", ZH: "指示器溢出", Class: "Data Exception"},
	"22015": {Code: "22015", EN: "interval_field_overflow", ZH: "间隔字段溢出", Class: "Data Exception"},
	"2201E": {Code: "2201E", EN: "invalid_argument_for_logarithm", ZH: "对数的无效参数", Class: "Data Exception"},
	"22014": {Code: "22014", EN: "invalid_argument_for_ntile_function", ZH: "NTILE 函数的无效参数", Class: "Data Exception"},
	"22016": {Code: "22016", EN: "invalid_argument_for_nth_value_function", ZH: "NTH_VALUE 函数的无效参数", Class: "Data Exception"},
	"2201F": {Code: "2201F", EN: "invalid_argument_for_power_function", ZH: "POWER 函数的无效参数", Class: "Data Exception"},
	"2201G": {Code: "2201G", EN: "invalid_argument_for_width_bucket_function", ZH: "WIDTH_BUCKET 函数的无效参数", Class: "Data Exception"},
	"22018": {Code: "22018", EN: "invalid_character_value_for_cast", ZH: "转换的无效字符值", Class: "Data Exception"},
	"22007": {Code: "22007", EN: "invalid_datetime_format", ZH: "无效的日期时间格式", Class: "Data Exception"},
	"22019": {Code: "22019", EN: "invalid_escape_character", ZH: "无效的转义字符", Class: "Data Exception"},
	"2200D": {Code: "2200D", EN: "invalid_escape_octet", ZH: "无效的转义八位字节", Class: "Data Exception"},
	"22025": {Code: "22025", EN: "invalid_escape_sequence", ZH: "无效的转义序列", Class: "Data Exception"},
	"22P06": {Code: "22P06", EN: "nonstandard_use_of_escape_character", ZH: "非标准使用转义字符", Class: "Data Exception"},
	"22010": {Code: "22010", EN: "invalid_indicator_parameter_value", ZH: "无效的指示器参数值", Class: "Data Exception"},
	"22023": {Code: "22023", EN: "invalid_parameter_value", ZH: "无效的参数值", Class: "Data Exception"},
	"22013": {Code: "22013", EN: "invalid_preceding_or_following_size", ZH: "无效的前导或跟随大小", Class: "Data Exception"},
	"2201B": {Code: "2201B", EN: "invalid_regular_expression", ZH: "无效的正则表达式", Class: "Data Exception"},
	"2201W": {Code: "2201W", EN: "invalid_row_count_in_limit_clause", ZH: "LIMIT 子句中的无效行数", Class: "Data Exception"},
	"2201X": {Code: "2201X", EN: "invalid_row_count_in_result_offset_clause", ZH: "结果偏移子句中的无效行数", Class: "Data Exception"},
	"2202H": {Code: "2202H", EN: "invalid_tablesample_argument", ZH: "无效的 TABLESAMPLE 参数", Class: "Data Exception"},
	"2202G": {Code: "2202G", EN: "invalid_tablesample_repeat", ZH: "无效的 TABLESAMPLE 重复", Class: "Data Exception"},
	"22009": {Code: "22009", EN: "invalid_time_zone_displacement_value", ZH: "无效的时区位移值", Class: "Data Exception"},
	"2200C": {Code: "2200C", EN: "invalid_use_of_escape_character", ZH: "无效的转义字符使用", Class: "Data Exception"},
	"2200G": {Code: "2200G", EN: "most_specific_type_mismatch", ZH: "最具体类型不匹配", Class: "Data Exception"},
	"22004": {Code: "22004", EN: "null_value_not_allowed", ZH: "不允许空值", Class: "Data Exception"},
	"22002": {Code: "22002", EN: "null_value_no_indicator_parameter", ZH: "空值无指示器参数", Class: "Data Exception"},
	"22003": {Code: "22003", EN: "numeric_value_out_of_range", ZH: "数值超出范围", Class: "Data Exception"},
	"2200H": {Code: "2200H", EN: "sequence_generator_limit_exceeded", ZH: "序列生成器超出限制", Class: "Data Exception"},
	"22026": {Code: "22026", EN: "string_data_length_mismatch", ZH: "字符串数据长度不匹配", Class: "Data Exception"},
	"22001": {Code: "22001", EN: "string_data_right_truncation", ZH: "字符串数据右截断", Class: "Data Exception"},
	"22011": {Code: "22011", EN: "substring_error", ZH: "子字符串错误", Class: "Data Exception"},
	"22027": {Code: "22027", EN: "trim_error", ZH: "修剪错误", Class: "Data Exception"},
	"22024": {Code: "22024", EN: "unterminated_c_string", ZH: "未终止的 C 字符串", Class: "Data Exception"},
	"2200F": {Code: "2200F", EN: "zero_length_character_string", ZH: "零长度字符串", Class: "Data Exception"},
	"22P01": {Code: "22P01", EN: "floating_point_exception", ZH: "浮点异常", Class: "Data Exception"},
	"22P02": {Code: "22P02", EN: "invalid_text_representation", ZH: "无效的文本表示", Class: "Data Exception"},
	"22P03": {Code: "22P03", EN: "invalid_binary_representation", ZH: "无效的二进制表示", Class: "Data Exception"},
	"22P04": {Code: "22P04", EN: "bad_copy_file_format", ZH: "错误的复制文件格式", Class: "Data Exception"},
	"22P05": {Code: "22P05", EN: "untranslatable_character", ZH: "不可翻译的字符", Class: "Data Exception"},
	"2200L": {Code: "2200L", EN: "not_an_xml_document", ZH: "不是 XML 文档", Class: "Data Exception"},
	"2200M": {Code: "2200M", EN: "invalid_xml_document", ZH: "无效的 XML 文档", Class: "Data Exception"},
	"2200N": {Code: "2200N", EN: "invalid_xml_content", ZH: "无效的 XML 内容", Class: "Data Exception"},
	"2200S": {Code: "2200S", EN: "invalid_xml_comment", ZH: "无效的 XML 注释", Class: "Data Exception"},
	"2200T": {Code: "2200T", EN: "invalid_xml_processing_instruction", ZH: "无效的 XML 处理指令", Class: "Data Exception"},
	"22030": {Code: "22030", EN: "duplicate_json_object_key_value", ZH: "重复的 JSON 对象键值", Class: "Data Exception"},
	"22031": {Code: "22031", EN: "invalid_argument_for_sql_json_datetime_function", ZH: "SQL/JSON 日期时间函数的无效参数", Class: "Data Exception"},
	"22032": {Code: "22032", EN: "invalid_json_text", ZH: "无效的 JSON 文本", Class: "Data Exception"},
	"22033": {Code: "22033", EN: "invalid_sql_json_subscript", ZH: "无效的 SQL/JSON 下标", Class: "Data Exception"},
	"22034": {Code: "22034", EN: "more_than_one_sql_json_item", ZH: "多个 SQL/JSON 项", Class: "Data Exception"},
	"22035": {Code: "22035", EN: "no_sql_json_item", ZH: "无 SQL/JSON 项", Class: "Data Exception"},
	"22036": {Code: "22036", EN: "non_numeric_sql_json_item", ZH: "非数字 SQL/JSON 项", Class: "Data Exception"},
	"22037": {Code: "22037", EN: "non_unique_keys_in_a_json_object", ZH: "JSON 对象中的非唯一键", Class: "Data Exception"},
	"22038": {Code: "22038", EN: "singleton_sql_json_item_required", ZH: "需要单个 SQL/JSON 项", Class: "Data Exception"},
	"22039": {Code: "22039", EN: "sql_json_array_not_found", ZH: "未找到 SQL/JSON 数组", Class: "Data Exception"},
	"2203A": {Code: "2203A", EN: "sql_json_member_not_found", ZH: "未找到 SQL/JSON 成员", Class: "Data Exception"},
	"2203B": {Code: "2203B", EN: "sql_json_number_not_found", ZH: "未找到 SQL/JSON 数字", Class: "Data Exception"},
	"2203C": {Code: "2203C", EN: "sql_json_object_not_found", ZH: "未找到 SQL/JSON 对象", Class: "Data Exception"},
	"2203D": {Code: "2203D", EN: "too_many_json_array_elements", ZH: "JSON 数组元素过多", Class: "Data Exception"},
	"2203E": {Code: "2203E", EN: "too_many_json_object_members", ZH: "JSON 对象成员过多", Class: "Data Exception"},
	"2203F": {Code: "2203F", EN: "sql_json_scalar_required", ZH: "需要 SQL/JSON 标量", Class: "Data Exception"},
	"2203G": {Code: "2203G", EN: "sql_json_item_cannot_be_cast_to_target_type", ZH: "SQL/JSON 项无法转换为目标类型", Class: "Data Exception"},

	// Class 23 — Integrity Constraint Violation
	"23000": {Code: "23000", EN: "integrity_constraint_violation", ZH: "完整性约束违规", Class: "Integrity Constraint Violation"},
	"23001": {Code: "23001", EN: "restrict_violation", ZH: "限制违规", Class: "Integrity Constraint Violation"},
	"23502": {Code: "23502", EN: "not_null_violation", ZH: "非空违规", Class: "Integrity Constraint Violation"},
	"23503": {Code: "23503", EN: "foreign_key_violation", ZH: "外键违规", Class: "Integrity Constraint Violation"},
	"23505": {Code: "23505", EN: "unique_violation", ZH: "唯一性违规", Class: "Integrity Constraint Violation"},
	"23514": {Code: "23514", EN: "check_violation", ZH: "检查违规", Class: "Integrity Constraint Violation"},
	"23P01": {Code: "23P01", EN: "exclusion_violation", ZH: "排除违规", Class: "Integrity Constraint Violation"},

	// Class 24 — Invalid Cursor State
	"24000": {Code: "24000", EN: "invalid_cursor_state", ZH: "无效的游标状态", Class: "Invalid Cursor State"},

	// Class 25 — Invalid Transaction State
	"25000": {Code: "25000", EN: "invalid_transaction_state", ZH: "无效的事务状态", Class: "Invalid Transaction State"},
	"25001": {Code: "25001", EN: "active_sql_transaction", ZH: "活动的 SQL 事务", Class: "Invalid Transaction State"},
	"25002": {Code: "25002", EN: "branch_transaction_already_active", ZH: "分支事务已激活", Class: "Invalid Transaction State"},
	"25008": {Code: "25008", EN: "held_cursor_requires_same_isolation_level", ZH: "持有的游标需要相同隔离级别", Class: "Invalid Transaction State"},
	"25003": {Code: "25003", EN: "inappropriate_access_mode_for_branch_transaction", ZH: "分支事务的访问模式不适当", Class: "Invalid Transaction State"},
	"25004": {Code: "25004", EN: "inappropriate_isolation_level_for_branch_transaction", ZH: "分支事务的隔离级别不适当", Class: "Invalid Transaction State"},
	"25005": {Code: "25005", EN: "no_active_sql_transaction_for_branch_transaction", ZH: "分支事务无活动的 SQL 事务", Class: "Invalid Transaction State"},
	"25006": {Code: "25006", EN: "read_only_sql_transaction", ZH: "只读 SQL 事务", Class: "Invalid Transaction State"},
	"25007": {Code: "25007", EN: "schema_and_data_statement_mixing_not_supported", ZH: "不支持模式和数据语句混合", Class: "Invalid Transaction State"},
	"25P01": {Code: "25P01", EN: "no_active_sql_transaction", ZH: "无活动的 SQL 事务", Class: "Invalid Transaction State"},
	"25P02": {Code: "25P02", EN: "in_failed_sql_transaction", ZH: "在失败的 SQL 事务中", Class: "Invalid Transaction State"},
	"25P03": {Code: "25P03", EN: "idle_in_transaction_session_timeout", ZH: "事务会话空闲超时", Class: "Invalid Transaction State"},
	"25P04": {Code: "25P04", EN: "transaction_timeout", ZH: "事务超时", Class: "Invalid Transaction State"},

	// Class 26 — Invalid SQL Statement Name
	"26000": {Code: "26000", EN: "invalid_sql_statement_name", ZH: "无效的 SQL 语句名称", Class: "Invalid SQL Statement Name"},

	// Class 27 — Triggered Data Change Violation
	"27000": {Code: "27000", EN: "triggered_data_change_violation", ZH: "触发数据更改违规", Class: "Triggered Data Change Violation"},

	// Class 28 — Invalid Authorization Specification
	"28000": {Code: "28000", EN: "invalid_authorization_specification", ZH: "无效的授权规范", Class: "Invalid Authorization Specification"},
	"28P01": {Code: "28P01", EN: "invalid_password", ZH: "无效的密码", Class: "Invalid Authorization Specification"},

	// Class 2B — Dependent Privilege Descriptors Still Exist
	"2B000": {Code: "2B000", EN: "dependent_privilege_descriptors_still_exist", ZH: "依赖的权限描述符仍然存在", Class: "Dependent Privilege Descriptors Still Exist"},
	"2BP01": {Code: "2BP01", EN: "dependent_objects_still_exist", ZH: "依赖的对象仍然存在", Class: "Dependent Privilege Descriptors Still Exist"},

	// Class 2D — Invalid Transaction Termination
	"2D000": {Code: "2D000", EN: "invalid_transaction_termination", ZH: "无效的事务终止", Class: "Invalid Transaction Termination"},

	// Class 2F — SQL Routine Exception
	"2F000": {Code: "2F000", EN: "sql_routine_exception", ZH: "SQL 例程异常", Class: "SQL Routine Exception"},
	"2F005": {Code: "2F005", EN: "function_executed_no_return_statement", ZH: "函数执行无返回语句", Class: "SQL Routine Exception"},
	"2F002": {Code: "2F002", EN: "modifying_sql_data_not_permitted", ZH: "不允许修改 SQL 数据", Class: "SQL Routine Exception"},
	"2F003": {Code: "2F003", EN: "prohibited_sql_statement_attempted", ZH: "尝试了禁止的 SQL 语句", Class: "SQL Routine Exception"},
	"2F004": {Code: "2F004", EN: "reading_sql_data_not_permitted", ZH: "不允许读取 SQL 数据", Class: "SQL Routine Exception"},

	// Class 34 — Invalid Cursor Name
	"34000": {Code: "34000", EN: "invalid_cursor_name", ZH: "无效的游标名称", Class: "Invalid Cursor Name"},

	// Class 38 — External Routine Exception
	"38000": {Code: "38000", EN: "external_routine_exception", ZH: "外部例程异常", Class: "External Routine Exception"},
	"38001": {Code: "38001", EN: "containing_sql_not_permitted", ZH: "不允许包含 SQL", Class: "External Routine Exception"},
	"38002": {Code: "38002", EN: "modifying_sql_data_not_permitted", ZH: "不允许修改 SQL 数据", Class: "External Routine Exception"},
	"38003": {Code: "38003", EN: "prohibited_sql_statement_attempted", ZH: "尝试了禁止的 SQL 语句", Class: "External Routine Exception"},
	"38004": {Code: "38004", EN: "reading_sql_data_not_permitted", ZH: "不允许读取 SQL 数据", Class: "External Routine Exception"},

	// Class 39 — External Routine Invocation Exception
	"39000": {Code: "39000", EN: "external_routine_invocation_exception", ZH: "外部例程调用异常", Class: "External Routine Invocation Exception"},
	"39001": {Code: "39001", EN: "invalid_sqlstate_returned", ZH: "返回了无效的 SQLSTATE", Class: "External Routine Invocation Exception"},
	"39004": {Code: "39004", EN: "null_value_not_allowed", ZH: "不允许空值", Class: "External Routine Invocation Exception"},
	"39P01": {Code: "39P01", EN: "trigger_protocol_violated", ZH: "违反触发器协议", Class: "External Routine Invocation Exception"},
	"39P02": {Code: "39P02", EN: "srf_protocol_violated", ZH: "违反 SRF 协议", Class: "External Routine Invocation Exception"},
	"39P03": {Code: "39P03", EN: "event_trigger_protocol_violated", ZH: "违反事件触发器协议", Class: "External Routine Invocation Exception"},

	// Class 3B — Savepoint Exception
	"3B000": {Code: "3B000", EN: "savepoint_exception", ZH: "保存点异常", Class: "Savepoint Exception"},
	"3B001": {Code: "3B001", EN: "invalid_savepoint_specification", ZH: "无效的保存点规范", Class: "Savepoint Exception"},

	// Class 3D — Invalid Catalog Name
	"3D000": {Code: "3D000", EN: "invalid_catalog_name", ZH: "无效的目录名称", Class: "Invalid Catalog Name"},

	// Class 3F — Invalid Schema Name
	"3F000": {Code: "3F000", EN: "invalid_schema_name", ZH: "无效的模式名称", Class: "Invalid Schema Name"},

	// Class 40 — Transaction Rollback
	"40000": {Code: "40000", EN: "transaction_rollback", ZH: "事务回滚", Class: "Transaction Rollback"},
	"40002": {Code: "40002", EN: "transaction_integrity_constraint_violation", ZH: "事务完整性约束违规", Class: "Transaction Rollback"},
	"40001": {Code: "40001", EN: "serialization_failure", ZH: "序列化失败", Class: "Transaction Rollback"},
	"40003": {Code: "40003", EN: "statement_completion_unknown", ZH: "语句完成未知", Class: "Transaction Rollback"},
	"40P01": {Code: "40P01", EN: "deadlock_detected", ZH: "检测到死锁", Class: "Transaction Rollback"},

	// Class 42 — Syntax Error or Access Rule Violation
	"42000": {Code: "42000", EN: "syntax_error_or_access_rule_violation", ZH: "语法错误或访问规则违规", Class: "Syntax Error or Access Rule Violation"},
	"42601": {Code: "42601", EN: "syntax_error", ZH: "语法错误", Class: "Syntax Error or Access Rule Violation"},
	"42501": {Code: "42501", EN: "insufficient_privilege", ZH: "权限不足", Class: "Syntax Error or Access Rule Violation"},
	"42846": {Code: "42846", EN: "cannot_coerce", ZH: "无法强制转换", Class: "Syntax Error or Access Rule Violation"},
	"42803": {Code: "42803", EN: "grouping_error", ZH: "分组错误", Class: "Syntax Error or Access Rule Violation"},
	"42P20": {Code: "42P20", EN: "windowing_error", ZH: "窗口错误", Class: "Syntax Error or Access Rule Violation"},
	"42P19": {Code: "42P19", EN: "invalid_recursion", ZH: "无效的递归", Class: "Syntax Error or Access Rule Violation"},
	"42830": {Code: "42830", EN: "invalid_foreign_key", ZH: "无效的外键", Class: "Syntax Error or Access Rule Violation"},
	"42602": {Code: "42602", EN: "invalid_name", ZH: "无效的名称", Class: "Syntax Error or Access Rule Violation"},
	"42622": {Code: "42622", EN: "name_too_long", ZH: "名称过长", Class: "Syntax Error or Access Rule Violation"},
	"42939": {Code: "42939", EN: "reserved_name", ZH: "保留名称", Class: "Syntax Error or Access Rule Violation"},
	"42804": {Code: "42804", EN: "datatype_mismatch", ZH: "数据类型不匹配", Class: "Syntax Error or Access Rule Violation"},
	"42P18": {Code: "42P18", EN: "indeterminate_datatype", ZH: "不确定的数据类型", Class: "Syntax Error or Access Rule Violation"},
	"42P21": {Code: "42P21", EN: "collation_mismatch", ZH: "排序规则不匹配", Class: "Syntax Error or Access Rule Violation"},
	"42P22": {Code: "42P22", EN: "indeterminate_collation", ZH: "不确定的排序规则", Class: "Syntax Error or Access Rule Violation"},
	"42809": {Code: "42809", EN: "wrong_object_type", ZH: "错误的对象类型", Class: "Syntax Error or Access Rule Violation"},
	"428C9": {Code: "428C9", EN: "generated_always", ZH: "始终生成", Class: "Syntax Error or Access Rule Violation"},
	"42703": {Code: "42703", EN: "undefined_column", ZH: "未定义的列", Class: "Syntax Error or Access Rule Violation"},
	"42883": {Code: "42883", EN: "undefined_function", ZH: "未定义的函数", Class: "Syntax Error or Access Rule Violation"},
	"42P01": {Code: "42P01", EN: "undefined_table", ZH: "未定义的表", Class: "Syntax Error or Access Rule Violation"},
	"42P02": {Code: "42P02", EN: "undefined_parameter", ZH: "未定义的参数", Class: "Syntax Error or Access Rule Violation"},
	"42704": {Code: "42704", EN: "undefined_object", ZH: "未定义的对象", Class: "Syntax Error or Access Rule Violation"},
	"42701": {Code: "42701", EN: "duplicate_column", ZH: "重复的列", Class: "Syntax Error or Access Rule Violation"},
	"42P03": {Code: "42P03", EN: "duplicate_cursor", ZH: "重复的游标", Class: "Syntax Error or Access Rule Violation"},
	"42P04": {Code: "42P04", EN: "duplicate_database", ZH: "重复的数据库", Class: "Syntax Error or Access Rule Violation"},
	"42723": {Code: "42723", EN: "duplicate_function", ZH: "重复的函数", Class: "Syntax Error or Access Rule Violation"},
	"42P05": {Code: "42P05", EN: "duplicate_prepared_statement", ZH: "重复的预处理语句", Class: "Syntax Error or Access Rule Violation"},
	"42P06": {Code: "42P06", EN: "duplicate_schema", ZH: "重复的模式", Class: "Syntax Error or Access Rule Violation"},
	"42P07": {Code: "42P07", EN: "duplicate_table", ZH: "重复的表", Class: "Syntax Error or Access Rule Violation"},
	"42712": {Code: "42712", EN: "duplicate_alias", ZH: "重复的别名", Class: "Syntax Error or Access Rule Violation"},
	"42710": {Code: "42710", EN: "duplicate_object", ZH: "重复的对象", Class: "Syntax Error or Access Rule Violation"},
	"42702": {Code: "42702", EN: "ambiguous_column", ZH: "歧义的列", Class: "Syntax Error or Access Rule Violation"},
	"42725": {Code: "42725", EN: "ambiguous_function", ZH: "歧义的函数", Class: "Syntax Error or Access Rule Violation"},
	"42P08": {Code: "42P08", EN: "ambiguous_parameter", ZH: "歧义的参数", Class: "Syntax Error or Access Rule Violation"},
	"42P09": {Code: "42P09", EN: "ambiguous_alias", ZH: "歧义的别名", Class: "Syntax Error or Access Rule Violation"},
	"42P10": {Code: "42P10", EN: "invalid_column_reference", ZH: "无效的列引用", Class: "Syntax Error or Access Rule Violation"},
	"42611": {Code: "42611", EN: "invalid_column_definition", ZH: "无效的列定义", Class: "Syntax Error or Access Rule Violation"},
	"42P11": {Code: "42P11", EN: "invalid_cursor_definition", ZH: "无效的游标定义", Class: "Syntax Error or Access Rule Violation"},
	"42P12": {Code: "42P12", EN: "invalid_database_definition", ZH: "无效的数据库定义", Class: "Syntax Error or Access Rule Violation"},
	"42P13": {Code: "42P13", EN: "invalid_function_definition", ZH: "无效的函数定义", Class: "Syntax Error or Access Rule Violation"},
	"42P14": {Code: "42P14", EN: "invalid_prepared_statement_definition", ZH: "无效的预处理语句定义", Class: "Syntax Error or Access Rule Violation"},
	"42P15": {Code: "42P15", EN: "invalid_schema_definition", ZH: "无效的模式定义", Class: "Syntax Error or Access Rule Violation"},
	"42P16": {Code: "42P16", EN: "invalid_table_definition", ZH: "无效的表定义", Class: "Syntax Error or Access Rule Violation"},
	"42P17": {Code: "42P17", EN: "invalid_object_definition", ZH: "无效的对象定义", Class: "Syntax Error or Access Rule Violation"},

	// Class 44 — WITH CHECK OPTION Violation
	"44000": {Code: "44000", EN: "with_check_option_violation", ZH: "WITH CHECK OPTION 违规", Class: "WITH CHECK OPTION Violation"},

	// Class 53 — Insufficient Resources
	"53000": {Code: "53000", EN: "insufficient_resources", ZH: "资源不足", Class: "Insufficient Resources"},
	"53100": {Code: "53100", EN: "disk_full", ZH: "磁盘已满", Class: "Insufficient Resources"},
	"53200": {Code: "53200", EN: "out_of_memory", ZH: "内存不足", Class: "Insufficient Resources"},
	"53300": {Code: "53300", EN: "too_many_connections", ZH: "连接过多", Class: "Insufficient Resources"},
	"53400": {Code: "53400", EN: "configuration_limit_exceeded", ZH: "超出配置限制", Class: "Insufficient Resources"},

	// Class 54 — Program Limit Exceeded
	"54000": {Code: "54000", EN: "program_limit_exceeded", ZH: "超出程序限制", Class: "Program Limit Exceeded"},
	"54001": {Code: "54001", EN: "statement_too_complex", ZH: "语句过于复杂", Class: "Program Limit Exceeded"},
	"54011": {Code: "54011", EN: "too_many_columns", ZH: "列过多", Class: "Program Limit Exceeded"},
	"54023": {Code: "54023", EN: "too_many_arguments", ZH: "参数过多", Class: "Program Limit Exceeded"},

	// Class 55 — Object Not In Prerequisite State
	"55000": {Code: "55000", EN: "object_not_in_prerequisite_state", ZH: "对象未处于先决状态", Class: "Object Not In Prerequisite State"},
	"55006": {Code: "55006", EN: "object_in_use", ZH: "对象正在使用中", Class: "Object Not In Prerequisite State"},
	"55P02": {Code: "55P02", EN: "cant_change_runtime_param", ZH: "无法更改运行时参数", Class: "Object Not In Prerequisite State"},
	"55P03": {Code: "55P03", EN: "lock_not_available", ZH: "锁不可用", Class: "Object Not In Prerequisite State"},
	"55P04": {Code: "55P04", EN: "unsafe_new_enum_value_usage", ZH: "不安全的枚举值使用", Class: "Object Not In Prerequisite State"},

	// Class 57 — Operator Intervention
	"57000": {Code: "57000", EN: "operator_intervention", ZH: "操作员干预", Class: "Operator Intervention"},
	"57014": {Code: "57014", EN: "query_canceled", ZH: "查询已取消", Class: "Operator Intervention"},
	"57P01": {Code: "57P01", EN: "admin_shutdown", ZH: "管理员关闭", Class: "Operator Intervention"},
	"57P02": {Code: "57P02", EN: "crash_shutdown", ZH: "崩溃关闭", Class: "Operator Intervention"},
	"57P03": {Code: "57P03", EN: "cannot_connect_now", ZH: "现在无法连接", Class: "Operator Intervention"},
	"57P04": {Code: "57P04", EN: "database_dropped", ZH: "数据库已删除", Class: "Operator Intervention"},
	"57P05": {Code: "57P05", EN: "idle_session_timeout", ZH: "空闲会话超时", Class: "Operator Intervention"},

	// Class 58 — System Error (errors external to PostgreSQL itself)
	"58000": {Code: "58000", EN: "system_error", ZH: "系统错误", Class: "System Error"},
	"58030": {Code: "58030", EN: "io_error", ZH: "I/O 错误", Class: "System Error"},
	"58P01": {Code: "58P01", EN: "undefined_file", ZH: "未定义的文件", Class: "System Error"},
	"58P02": {Code: "58P02", EN: "duplicate_file", ZH: "重复的文件", Class: "System Error"},

	// Class F0 — Configuration File Error
	"F0000": {Code: "F0000", EN: "config_file_error", ZH: "配置文件错误", Class: "Configuration File Error"},
	"F0001": {Code: "F0001", EN: "lock_file_exists", ZH: "锁定文件已存在", Class: "Configuration File Error"},

	// Class HV — Foreign Data Wrapper Error (SQL/MED)
	"HV000": {Code: "HV000", EN: "fdw_error", ZH: "外部数据包装器错误", Class: "Foreign Data Wrapper Error"},
	"HV005": {Code: "HV005", EN: "fdw_column_name_not_found", ZH: "未找到外部数据包装器列名", Class: "Foreign Data Wrapper Error"},
	"HV002": {Code: "HV002", EN: "fdw_dynamic_parameter_value_needed", ZH: "需要外部数据包装器动态参数值", Class: "Foreign Data Wrapper Error"},
	"HV010": {Code: "HV010", EN: "fdw_function_sequence_error", ZH: "外部数据包装器函数序列错误", Class: "Foreign Data Wrapper Error"},
	"HV021": {Code: "HV021", EN: "fdw_inconsistent_descriptor_information", ZH: "外部数据包装器描述符信息不一致", Class: "Foreign Data Wrapper Error"},
	"HV024": {Code: "HV024", EN: "fdw_invalid_attribute_value", ZH: "外部数据包装器无效的属性值", Class: "Foreign Data Wrapper Error"},
	"HV007": {Code: "HV007", EN: "fdw_invalid_column_name", ZH: "外部数据包装器无效的列名", Class: "Foreign Data Wrapper Error"},
	"HV008": {Code: "HV008", EN: "fdw_invalid_column_number", ZH: "外部数据包装器无效的列号", Class: "Foreign Data Wrapper Error"},
	"HV004": {Code: "HV004", EN: "fdw_invalid_data_type", ZH: "外部数据包装器无效的数据类型", Class: "Foreign Data Wrapper Error"},
	"HV006": {Code: "HV006", EN: "fdw_invalid_data_type_descriptors", ZH: "外部数据包装器无效的数据类型描述符", Class: "Foreign Data Wrapper Error"},
	"HV091": {Code: "HV091", EN: "fdw_invalid_descriptor_field_identifier", ZH: "外部数据包装器无效的描述符字段标识符", Class: "Foreign Data Wrapper Error"},
	"HV00B": {Code: "HV00B", EN: "fdw_invalid_handle", ZH: "外部数据包装器无效的句柄", Class: "Foreign Data Wrapper Error"},
	"HV00C": {Code: "HV00C", EN: "fdw_invalid_option_index", ZH: "外部数据包装器无效的选项索引", Class: "Foreign Data Wrapper Error"},
	"HV00D": {Code: "HV00D", EN: "fdw_invalid_option_name", ZH: "外部数据包装器无效的选项名称", Class: "Foreign Data Wrapper Error"},
	"HV090": {Code: "HV090", EN: "fdw_invalid_string_length_or_buffer_length", ZH: "外部数据包装器无效的字符串长度或缓冲区长度", Class: "Foreign Data Wrapper Error"},
	"HV00A": {Code: "HV00A", EN: "fdw_invalid_string_format", ZH: "外部数据包装器无效的字符串格式", Class: "Foreign Data Wrapper Error"},
	"HV009": {Code: "HV009", EN: "fdw_invalid_use_of_null_pointer", ZH: "外部数据包装器无效的空指针使用", Class: "Foreign Data Wrapper Error"},
	"HV014": {Code: "HV014", EN: "fdw_too_many_handles", ZH: "外部数据包装器句柄过多", Class: "Foreign Data Wrapper Error"},
	"HV001": {Code: "HV001", EN: "fdw_out_of_memory", ZH: "外部数据包装器内存不足", Class: "Foreign Data Wrapper Error"},
	"HV00P": {Code: "HV00P", EN: "fdw_no_schemas", ZH: "外部数据包装器无模式", Class: "Foreign Data Wrapper Error"},
	"HV00J": {Code: "HV00J", EN: "fdw_option_name_not_found", ZH: "未找到外部数据包装器选项名称", Class: "Foreign Data Wrapper Error"},
	"HV00K": {Code: "HV00K", EN: "fdw_reply_handle", ZH: "外部数据包装器回复句柄", Class: "Foreign Data Wrapper Error"},
	"HV00Q": {Code: "HV00Q", EN: "fdw_schema_not_found", ZH: "未找到外部数据包装器模式", Class: "Foreign Data Wrapper Error"},
	"HV00R": {Code: "HV00R", EN: "fdw_table_not_found", ZH: "未找到外部数据包装器表", Class: "Foreign Data Wrapper Error"},
	"HV00L": {Code: "HV00L", EN: "fdw_unable_to_create_execution", ZH: "外部数据包装器无法创建执行", Class: "Foreign Data Wrapper Error"},
	"HV00M": {Code: "HV00M", EN: "fdw_unable_to_create_reply", ZH: "外部数据包装器无法创建回复", Class: "Foreign Data Wrapper Error"},
	"HV00N": {Code: "HV00N", EN: "fdw_unable_to_establish_connection", ZH: "外部数据包装器无法建立连接", Class: "Foreign Data Wrapper Error"},

	// Class P0 — PL/pgSQL Error
	"P0000": {Code: "P0000", EN: "plpgsql_error", ZH: "PL/pgSQL 错误", Class: "PL/pgSQL Error"},
	"P0001": {Code: "P0001", EN: "raise_exception", ZH: "引发异常", Class: "PL/pgSQL Error"},
	"P0002": {Code: "P0002", EN: "no_data_found", ZH: "未找到数据", Class: "PL/pgSQL Error"},
	"P0003": {Code: "P0003", EN: "too_many_rows", ZH: "行数过多", Class: "PL/pgSQL Error"},
	"P0004": {Code: "P0004", EN: "assert_failure", ZH: "断言失败", Class: "PL/pgSQL Error"},

	// Class XX — Internal Error
	"XX000": {Code: "XX000", EN: "internal_error", ZH: "内部错误", Class: "Internal Error"},
	"XX001": {Code: "XX001", EN: "data_corrupted", ZH: "数据损坏", Class: "Internal Error"},
	"XX002": {Code: "XX002", EN: "index_corrupted", ZH: "索引损坏", Class: "Internal Error"},
}
