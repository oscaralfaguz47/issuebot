import { useAuthStore } from "@/features/auth";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { createProject, getProjects } from "@/features/projects/api";
import LoadingScreen from "@/components/LoadingScreen";

function ProjectList() {
  // --- All hooks at the top, together, before any return ---
  const queryClient = useQueryClient();
  const memberships = useAuthStore((s) => s.memberships);

  const { data: projects, isPending, error } = useQuery({
    queryKey: ["projects"],
    queryFn: getProjects,
  });

  const createMutation = useMutation({
    mutationFn: (orgId: string) =>
      createProject(orgId, `Nuevo project ${Date.now()}`),
    onSuccess: () => {
      // Mark the cached projects as stale -> React Query refetches in background
      queryClient.invalidateQueries({ queryKey: ["projects"] });
    },
  });

  // --- Derived values and handlers (not hooks, can go here) ---
  const canCreate = memberships.some(
    (m) => m.role === "owner" || m.role === "member"
  );

  const handleCreate = () => {
    const orgId = memberships[0]?.org_id;
    if (!orgId) return;
    createMutation.mutate(orgId);
  };

  // --- Conditional returns, after all hooks ---
  if (isPending) return <LoadingScreen />;
  if (error) return <p>Error loading projects</p>;

  // --- Main render ---
  return (
    <div>
      <h2>Your projects ({projects.length})</h2>
      <ul>
        {projects.map((p) => (
          <li key={p.id}>{p.name}</li>
        ))}
      </ul>
      {canCreate && (
        <button onClick={handleCreate} disabled={createMutation.isPending}>
          {createMutation.isPending ? "Creando..." : "Crear project"}
        </button>
      )}
    </div>
  );
}

export default ProjectList;
