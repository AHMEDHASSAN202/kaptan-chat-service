package localization

// ErrorCode is an enum representing all message IDs

const (
	Account_Email_is_unique_rules_validation    = "Account_Email_is_unique_rules_validation"    // Email is Unique
	User_Email_is_unique_rules_validation       = "User_Email_is_unique_rules_validation"       // Email is Unique
	App_type_is_unique_rules_validation         = "App_type_is_unique_rules_validation"         // Feature not supported
	Cuisine_id_is_exists_rules_validation       = "Cuisine_id_is_exists_rules_validation"       // Cuisine Id Not Exist
	E1000                                       = "E1000"                                       // Database layer error
	E1001                                       = "E1001"                                       // Invalid input
	E1002                                       = "E1002"                                       // Resource not found
	E1002Item                                   = "E1002Item"                                   // Menu item not found
	E1003                                       = "E1003"                                       // Operation not permitted
	E1004                                       = "E1004"                                       // Internal server error
	E1005                                       = "E1005"                                       // Service unavailable
	E1006                                       = "E1006"                                       // Unauthorized access
	E1007                                       = "E1007"                                       // Session expired
	E1008                                       = "E1008"                                       // Validation error
	E1009                                       = "E1009"                                       // Timeout
	E1010                                       = "E1010"                                       // Conflict
	E1011                                       = "E1011"                                       // Quota exceeded
	E1012                                       = "E1012"                                       // Feature not supported
	Invalid_mongo_ids_validation_rule           = "Invalid_mongo_ids_validation_rule"           // Invalid mongo ids validation rule
	Item_name_is_unique_rules_validation        = "Item_name_is_unique_rules_validation"        // Feature not supported
	Modifier_groups_ids_rules_validation        = "Modifier_groups_ids_rules_validation"        // Feature not supported
	Modifier_items_cant_contains_modifier_group = "Modifier_items_cant_contains_modifier_group" // Modifier items cant contains modifier group
	SKU_name_is_unique_rules_validation         = "SKU_name_is_unique_rules_validation"         // Feature not supported
	Timeformat                                  = "Timeformat"                                  // time format rule validation
	PhoneNumberValidator                        = "phonenumber_rule"                            // phone number rule validation
)
