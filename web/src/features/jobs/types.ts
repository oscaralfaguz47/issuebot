export interface Job {
  id: string;
  type: string;
  status: string;
  attempts: number;
  created_at: string;
}