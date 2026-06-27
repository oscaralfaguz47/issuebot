import { useState, useEffect } from "react";
import { createProject, getProjects } from "../lib/api";
import type { Project } from "../types";
import { useAuthStore } from "../store/authStore";


// dentro del componente Dashboard:
function Dashboard() {
  const [projects, setProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const memberships = useAuthStore((s) => s.memberships);
  const canCreate = memberships.some((m) => m.role === "owner" || m.role === "member");

const handleCreate = async () => {
  const orgId = memberships[0]?.org_id;
  if (!orgId) return;

  try {
    await createProject(orgId, `Nuevo project ${Date.now()}`);
    const updated = await getProjects(); // recargar la lista
    setProjects(updated);
  } catch (err) {
    console.error("error creando project:", err);
  }
};

  useEffect(() => {
    getProjects()
      .then((data) => setProjects(data))
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <p>Loading projects...</p>;
  if (error) return <p>Error: {error}</p>;

  return (
    <div>
      <h1>Dashboard</h1>
      <h2>Your projects ({projects.length})</h2>
      <ul>
        {projects.map((p) => (
          <li key={p.id}>{p.name}</li>
        ))}
      </ul>
      {canCreate && (
        <button onClick={handleCreate}>Crear project</button>
      )}
    </div>
  );
}

export default Dashboard;