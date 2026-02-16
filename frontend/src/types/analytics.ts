export interface ProjectActivity {
  project_id: number;
  project_name: string;
  completed_tasks_count: number;
}

export interface MemberWorkload {
  user_id: number;
  full_name: string;
  active_tasks_count: number;
}

export interface DailyThroughput {
  date: string;
  count: number;
}

export interface StaleProject {
  project_id: number;
  project_name: string;
  stale_tasks_count: number;
}

export interface PriorityDistribution {
  priority: string;
  count: number;
}

export interface OverdueTask {
  task_id: number;
  title: string;
  assignee_name: string;
  due_date: string;
  priority: string;
}

export interface UserPerformanceMetric {
  user_id: number;
  full_name: string;
  completed_tasks_count: number;
  on_time_rate: number;
  avg_lead_time: number;
  active_tasks_count: number;
  overdue_tasks_count: number;
}

export interface AnalyticsOverview {
  top_projects: ProjectActivity[];
  member_workload: MemberWorkload[];
  throughput: DailyThroughput[];
  stale_projects: StaleProject[];

  // Advanced
  average_lead_time: number;
  on_time_completion_rate: number;
  priority_distribution: PriorityDistribution[];
  overdue_tasks: OverdueTask[];

  // Team Insights
  user_performance: UserPerformanceMetric[];
}
