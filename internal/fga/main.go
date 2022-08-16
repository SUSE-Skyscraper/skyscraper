package fga

type Relation string
type Document string

const DefaultOrganizationID = "default"
const DocumentOrganization Document = "organization"
const DocumentOrganizationRelationAuditLogViewer Relation = "audit_logs_viewer"
const DocumentOrganizationRelationAPIKeysViewer Relation = "api_keys_viewer"
const DocumentOrganizationRelationAPIKeysEditor Relation = "api_keys_editor"
const DocumentOrganizationRelationStandardTagsViewer Relation = "standard_tags_viewer"
const DocumentOrganizationRelationStandardTagsEditor Relation = "standard_tags_editor"
const DocumentOrganizationRelationUsersViewer Relation = "users_viewer"
const DocumentOrganizationRelationCloudAccountsViewer Relation = "cloud_accounts_viewer"
const DocumentOrganizationRelationCloudTenantsViewer Relation = "cloud_tenants_viewer"
const DocumentOrganizationRelationOrganizationalUnitsViewer Relation = "organizational_units_viewer"
const DocumentOrganizationRelationOrganizationalUnitsEditor Relation = "organizational_units_editor"

const DocumentAccount Document = "account"
const DocumentAccountRelationViewer Relation = "viewer"
const DocumentAccountRelationEditor Relation = "editor"
