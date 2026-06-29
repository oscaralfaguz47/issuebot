
import { supabase } from "@/lib/supabase";
import { API_URL } from "@/lib/config"
import type { Project } from "../types";


export async function getProjects(): Promise<Project[]> {
  const { data } = await supabase.auth.getSession();
  const token = data.session?.access_token;

  const res = await fetch(`${API_URL}/projects`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  if (!res.ok) {
    throw new Error(`API error: ${res.status}`);
  }

  return res.json();
}

export async function createProject(orgId: string, name: string): Promise<void> {
  const { data } = await supabase.auth.getSession();
  const token = data.session?.access_token;

  const res = await fetch(`${API_URL}/projects`, {
    method: "POST",
    headers: {
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ org_id: orgId, name }),
  });

  if (!res.ok) throw new Error(`API error: ${res.status}`);
}