import { Badge } from "@/components/ui/badge";
import type { RequisitionPriority } from "../types";

const priorityConfig: Record<RequisitionPriority, { label: string; className: string }> = {
  LOW: {
    label: "Low",
    className: "bg-slate-100 text-slate-600 border-slate-200",
  },
  MEDIUM: {
    label: "Medium",
    className: "bg-blue-50 text-blue-700 border-blue-200",
  },
  HIGH: {
    label: "High",
    className: "bg-amber-50 text-amber-700 border-amber-200",
  },
  URGENT: {
    label: "Urgent",
    className: "bg-red-50 text-red-700 border-red-200",
  },
};

export function PriorityBadge({ priority }: { priority: RequisitionPriority }) {
  const config = priorityConfig[priority] ?? priorityConfig.MEDIUM;
  return (
    <Badge variant="outline" className={config.className}>
      {config.label}
    </Badge>
  );
}
