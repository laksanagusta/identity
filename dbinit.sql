CREATE TABLE IF NOT EXISTS organizations (
    uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    name VARCHAR(255) NULL,
    code VARCHAR(50) NULL UNIQUE,
    address TEXT NULL,
    type VARCHAR(50) NULL,
    parent_uuid UUID NULL REFERENCES organizations(uuid) ON DELETE SET NULL,
    level INT NULL,
    path TEXT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE

    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITHOUT TIME ZONE,
);

CREATE TABLE IF NOT EXISTS permissions (
    uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    name VARCHAR(100) NULL UNIQUE,
    action VARCHAR(50) NULL,
    resource VARCHAR(100) NULL,
    description TEXT NULL,

    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    
    CONSTRAINT unique_action_resource UNIQUE (action, resource)
);

CREATE TABLE IF NOT EXISTS roles (
    uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    name VARCHAR(100) NULL UNIQUE,
    description TEXT NULL,
    is_system BOOLEAN NOT NULL DEFAULT FALSE

    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITHOUT TIME ZONE,
);

CREATE TABLE IF NOT EXISTS users (
    uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    username VARCHAR(100) NULL UNIQUE,
    email VARCHAR(255) NULL UNIQUE,
    first_name VARCHAR(100) NULL,
    last_name VARCHAR(100) NULL,
    phone_number VARCHAR(20) NULL,
    password_hash VARCHAR(255) NULL,
    organization_uuid UUID NULL REFERENCES organizations(uuid) ON DELETE SET NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    last_login_at TIMESTAMP WITHOUT TIME ZONE NULL

    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITHOUT TIME ZONE,
);

CREATE TABLE IF NOT EXISTS role_permissions (
    uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    role_uuid UUID NOT NULL REFERENCES roles(uuid) ON DELETE CASCADE,
    permission_uuid UUID NOT NULL REFERENCES permissions(uuid) ON DELETE CASCADE,

    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    
    CONSTRAINT unique_role_permission UNIQUE (role_uuid, permission_uuid)
);

CREATE TABLE IF NOT EXISTS user_roles (
    uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    user_uuid UUID NOT NULL REFERENCES users(uuid) ON DELETE CASCADE,
    role_uuid UUID NOT NULL REFERENCES roles(uuid) ON DELETE CASCADE,

    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT unique_user_role UNIQUE (user_uuid, role_uuid)
);