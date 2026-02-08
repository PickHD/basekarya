import { Badge } from "@/components/ui/badge";

export const LeaveTypeBadge = ({ status }: { status: string }) => {
  const styles = {
    Sick: "bg-yellow-100 text-yellow-700 border-yellow-200",
    Annual: "bg-green-100 text-green-700 border-green-200",
    Unpaid: "bg-red-100 text-red-700 border-red-200",
  };
  return (
    <Badge
      variant="outline"
      className={styles[status as keyof typeof styles] || ""}
    >
      {status}
    </Badge>
  );
};
