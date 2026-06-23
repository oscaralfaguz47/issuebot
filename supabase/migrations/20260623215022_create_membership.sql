create table memberships (
    id uuid primary key default gen_random_uuid(),
    org_id uuid not null references orgs(id) on delete cascade,
    user_id uuid not null,
    role text not null check (role in ('owner','member', 'viewer')),
    created_at timestamptz not null default now(),
    unique (org_id, user_id)
);