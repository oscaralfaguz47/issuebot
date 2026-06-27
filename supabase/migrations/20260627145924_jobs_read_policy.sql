-- habilitar RLS en la tabla
alter table jobs enable row level security;

-- permitir lectura (select) a cualquier usuario autenticado
create policy "authenticated can read jobs"
on jobs
for select
to authenticated
using (true);