import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Users, UserCheck, Clock, UserX, Loader2 } from "lucide-react";
import type { DashboardStats } from "@/features/admin/types";

interface StatsCardsProps {
  data?: DashboardStats;
  isLoading: boolean;
}

export function StatsCards({ data, isLoading }: StatsCardsProps) {
  if (isLoading) {
    return (
      <div className="flex gap-4">
        <Loader2 className="animate-spin" /> Loading stats...
      </div>
    );
  }

  const cards = [
    {
      title: "Total Employees",
      value: data?.total_employees || 0,
      icon: Users,
      desc: "Registered active employees",
      color: "text-blue-600",
      bg: "bg-blue-100",
    },
    {
      title: "Present Today",
      value: data?.present_today || 0,
      icon: UserCheck,
      desc: "Checked in today",
      color: "text-green-600",
      bg: "bg-green-100",
    },
    {
      title: "Late Today",
      value: data?.late_today || 0,
      icon: Clock,
      desc: "Arrived after shift start",
      color: "text-orange-600",
      bg: "bg-orange-100",
    },
    {
      title: "Not Present",
      value: data?.absent_today || 0,
      icon: UserX,
      desc: "Haven't checked in yet",
      color: "text-red-600",
      bg: "bg-red-100",
    },
  ];

  return (
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      {cards.map((item) => (
        <Card key={item.title}>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">{item.title}</CardTitle>
            <div className={`p-2 rounded-full ${item.bg}`}>
              <item.icon className={`h-4 w-4 ${item.color}`} />
            </div>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{item.value}</div>
            <p className="text-xs text-muted-foreground">{item.desc}</p>
          </CardContent>
        </Card>
      ))}
    </div>
  );
}
