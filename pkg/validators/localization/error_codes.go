package localization

// ErrorCode is an enum representing all message IDs

const (
    AccountRequiredValidation = "AccountRequiredValidation" // AccountRequiredValidation
    Account_Email_is_unique_rules_validation = "Account_Email_is_unique_rules_validation" // Email is Unique
    App_type_is_unique_rules_validation = "App_type_is_unique_rules_validation" // Feature not supported
    Cuisine_id_is_exists_rules_validation = "Cuisine_id_is_exists_rules_validation" // Cuisine Id Not Exist
    E1000 = "E1000" // Database layer error
    E1001 = "E1001" // Invalid input
    E1002 = "E1002" // Resource not found
    E1002Item = "E1002Item" // Menu item not found
    E1003 = "E1003" // Operation not permitted
    E1004 = "E1004" // Internal server error
    E1005 = "E1005" // Service unavailable
    E1006 = "E1006" // Unauthorized access
    E1007 = "E1007" // Session expired
    E1008 = "E1008" // Validation error
    E1009 = "E1009" // Timeout
    E1010 = "E1010" // Conflict
    E1011 = "E1011" // Quota exceeded
    E1012 = "E1012" // Feature not supported
    E1013 = "E1013" // Invalid OTP
    E1014 = "E1014" // OTP has expired
    E1015 = "E1015" // Exceeded maximum OTP trials
    E1401 = "E1401" // Unauthorized access
    E1403 = "E1403" // Unauthenticated access
    Email_is_unique_rules_validation = "Email_is_unique_rules_validation" // Email is Unique
    ErrLoginBlocked = "ErrLoginBlocked" // Your Account Is disabled from admin
    ErrLoginEmail = "ErrLoginEmail" // Incorrect email
    ErrLoginInActive = "ErrLoginInActive" // Your account is currently inactive. Please contact our support team for assistance in reactivating your account.
    ErrLoginPassword = "ErrLoginPassword" // Incorrect password
    Invalid_mongo_ids_validation_rule = "Invalid_mongo_ids_validation_rule" // Invalid mongo ids validation rule
    Item_name_is_unique_rules_validation = "Item_name_is_unique_rules_validation" // Feature not supported
    JwtSigningError = "JwtSigningError" // Error generating signed string representation of JWT: there was a problem with the signing process
    JwtTokenExpiredError = "JwtTokenExpiredError" // Error: the JSON Web Token has expired and cannot be verified
    JwtTokenInvalidError = "JwtTokenInvalidError" // The JSON Web Token is invalid and could not be verified
    JwtTokenParsingError = "JwtTokenParsingError" // Error: the JSON Web Token could not be parsed
    Mobile_location_not_available_error = "Mobile_location_not_available_error" // mobile_location_not_available_error
    Modifier_groups_ids_rules_validation = "Modifier_groups_ids_rules_validation" // Feature not supported
    Modifier_items_cant_contains_modifier_group = "Modifier_items_cant_contains_modifier_group" // Modifier items cant contains modifier group
    Password_required_if_id_is_zero = "Password_required_if_id_is_zero" // password is a required field
    PhoneNumber_rule_validation = "PhoneNumber_rule_validation" // PhoneNumber is Wrong
    PreventDeleteRolesIdsValidation = "PreventDeleteRolesIdsValidation" // You are not allowed to delete main roles
    RoleExistsValidation = "RoleExistsValidation" // Role is not exists
    RoleHasAdminsValidation = "RoleHasAdminsValidation" // The role you are trying to delete is currently assigned to one or more admins. Please reassign these admins to a different role before attempting to delete this role.
    SKU_name_is_unique_rules_validation = "SKU_name_is_unique_rules_validation" // Feature not supported
    Timeformat = "Timeformat" // time format rule validation
)
