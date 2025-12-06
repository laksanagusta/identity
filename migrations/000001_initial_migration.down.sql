-- Drop all tables in reverse dependency order

DROP TABLE IF EXISTS proposal_document_uploads;
DROP TABLE IF EXISTS proposal_level_organization_users;
DROP TABLE IF EXISTS proposal_level_organizations;
DROP TABLE IF EXISTS proposal_levels;
DROP TABLE IF EXISTS proposal_history;
DROP TABLE IF EXISTS proposal_comments;
DROP TABLE IF EXISTS proposals;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS document_upload_rules;
DROP TABLE IF EXISTS workflow_level_organization_users;
DROP TABLE IF EXISTS workflow_level_organizations;
DROP TABLE IF EXISTS workflow_level_conditions;
DROP TABLE IF EXISTS workflow_levels;
DROP TABLE IF EXISTS workflows;
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS organizations;

-- Optionally drop the UUID extension
DROP EXTENSION IF EXISTS "uuid-ossp";
