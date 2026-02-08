import { Badge } from "@/components/ui/badge";

export const LeaveStatusBadge = ({ status }: { status: string }) => {
  const styles = {
    PENDING: "bg-yellow-100 text-yellow-700 border-yellow-200",
    APPROVED: "bg-green-100 text-green-700 border-green-200",
    REJECTED: "bg-red-100 text-red-700 border-red-200",
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
