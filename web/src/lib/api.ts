import { supabase } from "./supabase";
import type { Project } from "../types";

const API_URL = "http://localhost:8080";

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