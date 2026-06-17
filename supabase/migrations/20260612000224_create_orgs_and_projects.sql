create table orgs ( 
  id uuid primary key default gen_random_uuid(),
  name text not null,
  created_at timestamptz not null default now()
);

create table projects (
  id uuid primary key default gen_random_uuid(),
  org_id uuid not null references orgs(id) on delete cascade,
  name text not null,
  github_repo text not null,
  github_installation_id text not null,
  created_at timestamptz not null default now()
);
