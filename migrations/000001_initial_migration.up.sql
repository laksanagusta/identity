-- Database Schema DDL untuk PostgreSQL
-- ==================================================

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Organizations table dengan hierarchical support
CREATE TABLE organizations (
    uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    code VARCHAR(100) NOT NULL UNIQUE,
    address TEXT,
    type VARCHAR(100) NOT NULL,
    parent_uuid UUID REFERENCES organizations(uuid),
    level INTEGER DEFAULT 0,
    path TEXT NOT NULL, -- Materialized path: uuid.parent_uuid.grandparent_uuid
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by VARCHAR(255) NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_by VARCHAR(255) NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,
    deleted_by VARCHAR(255)
);

-- Permissions table
CREATE TABLE permissions (
    uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    action VARCHAR(100) NOT NULL,
    resource VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by VARCHAR(255) NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_by VARCHAR(255) NOT NULL
);

-- Roles table
CREATE TABLE roles (
    uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    is_system BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by VARCHAR(255) NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_by VARCHAR(255) NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,
    deleted_by VARCHAR(255)
);

-- Role permissions junction table
CREATE TABLE role_permissions (
    uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    role_uuid UUID NOT NULL REFERENCES roles(uuid) ON DELETE CASCADE,
    permission_uuid UUID NOT NULL REFERENCES permissions(uuid) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by VARCHAR(255) NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_by VARCHAR(255) NOT NULL,
    UNIQUE(role_uuid, permission_uuid)
);

-- Users table
CREATE TABLE users (
    uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(255) UNIQUE,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255),
    phone_number VARCHAR(20),
    password_hash VARCHAR(255) NOT NULL,
    organization_uuid UUID NOT NULL REFERENCES organizations(uuid),
    is_active BOOLEAN DEFAULT true,
    last_login_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by VARCHAR(255) NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_by VARCHAR(255) NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,
    deleted_by VARCHAR(255)
);

-- User roles junction table
CREATE TABLE user_roles (
    uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_uuid UUID NOT NULL REFERENCES users(uuid) ON DELETE CASCADE,
    role_uuid UUID NOT NULL REFERENCES roles(uuid) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by VARCHAR(255) NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_by VARCHAR(255) NOT NULL,
    UNIQUE(user_uuid, role_uuid)
);

-- Main Organization
INSERT INTO organizations (uuid, name, code, address, type, parent_uuid, level, path, is_active, created_by, updated_by)
VALUES
    ('11111111-1111-1111-1111-111111111111', 'PT. Utama Maju', 'UTAMA', 'Jl. Sudirman No. 123, Jakarta', 'Headquarters', NULL, 0, '11111111-1111-1111-1111-111111111111', TRUE, 'admin', 'admin'),
    ('22222222-2222-2222-2222-222222222222', 'Divisi Marketing', 'MKT', 'Jl. Thamrin No. 45, Jakarta', 'Division', '11111111-1111-1111-1111-111111111111', 1, '11111111-1111-1111-1111-111111111111.22222222-2222-2222-2222-222222222222', TRUE, 'admin', 'admin'),
    ('33333333-3333-3333-3333-333333333333', 'Divisi Keuangan', 'FIN', 'Jl. Gatot Subroto No. 67, Jakarta', 'Division', '11111111-1111-1111-1111-111111111111', 1, '11111111-1111-1111-1111-111111111111.33333333-3333-3333-3333-333333333333', TRUE, 'admin', 'admin');

-- Users
INSERT INTO users (uuid, username, email, first_name, last_name, password_hash, organization_uuid, created_by, updated_by)
VALUES
    ('a0a0a0a0-a0a0-a0a0-a0a0-a0a0a0a0a0a0', 'johndoe', 'john.doe@utama.com', 'John', 'Doe', 'hash-password-1', '22222222-2222-2222-2222-222222222222', 'admin', 'admin'),
    ('b1b1b1b1-b1b1-b1b1-b1b1-b1b1b1b1b1b1', 'janedoe', 'jane.doe@utama.com', 'Jane', 'Doe', 'hash-password-2', '33333333-3333-3333-3333-333333333333', 'admin', 'admin');

INSERT INTO permissions (uuid, name, action, resource, created_by, updated_by)
VALUES
    ('c1c1c1c1-c1c1-c1c1-c1c1-c1c1c1c1c1c1', 'Create Proposal', 'create', 'proposal', 'admin', 'admin'),
    ('d2d2d2d2-d2d2-d2d2-d2d2-d2d2d2d2d2d2', 'Approve Proposal', 'approve', 'proposal', 'admin', 'admin');
INSERT INTO roles (uuid, name, description, is_system, created_by, updated_by)
VALUES
    ('e1e1e1e1-e1e1-e1e1-e1e1-e1e1e1e1e1e1', 'Staff Marketing', 'Staff di divisi marketing', FALSE, 'admin', 'admin'),
    ('f2f2f2f2-f2f2-f2f2-f2f2-f2f2f2f2f2f2', 'Manajer Keuangan', 'Manajer di divisi keuangan', FALSE, 'admin', 'admin');

-- User Roles
INSERT INTO user_roles (user_uuid, role_uuid, created_by, updated_by)
VALUES
    ('a0a0a0a0-a0a0-a0a0-a0a0-a0a0a0a0a0a0', 'e1e1e1e1-e1e1-e1e1-e1e1-e1e1e1e1e1e1', 'admin', 'admin'),
    ('b1b1b1b1-b1b1-b1b1-b1b1-b1b1b1b1b1b1', 'f2f2f2f2-f2f2-f2f2-f2f2-f2f2f2f2f2f2', 'admin', 'admin');

-- Role Permissions
INSERT INTO role_permissions (role_uuid, permission_uuid, created_by, updated_by)
VALUES
    ('e1e1e1e1-e1e1-e1e1-e1e1-e1e1e1e1e1e1', 'c1c1c1c1-c1c1-c1c1-c1c1-c1c1c1c1c1c1', 'admin', 'admin'),
    ('f2f2f2f2-f2f2-f2f2-f2f2-f2f2f2f2f2f2', 'd2d2d2d2-d2d2-d2d2-d2d2-d2d2d2d2d2d2', 'admin', 'admin');