create table jobs(
    id uuid primary key default gen_random_uuid(),
    project_id uuid not null references projects(id) on delete cascade,
    type text not null,
    status text not null default 'pending' check (status in ('pending','processing','done','failed')),
    payload jsonb not null,
    attempts int not null default 0,
    idempotency_key text not null unique,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);