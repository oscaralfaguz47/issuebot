export interface Project {
  ID: string;
  OrgID: string;
  Name: string;
  GitHubRepo: string;
  GitHubInstallationID: string;
  CreatedAt: string;
}

export interface Membership {
  org_id: string;
  user_id: string;
  role: string;
}