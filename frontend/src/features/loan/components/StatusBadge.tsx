import type { LoanStatus } from "@/features/loan/types";

export const StatusBadge = ({ status }: { status: LoanStatus }) => {
  const styles = {
    PENDING: "bg-yellow-100 text-yellow-800",
    APPROVED: "bg-blue-100 text-blue-800",
    REJECTED: "bg-red-100 text-red-800",
    PAID_OFF: "bg-green-100 text-green-800",
  };

  return (
    <span
      className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${styles[status]}`}
    >
      {status}
    </span>
  );
};
