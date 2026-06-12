import type { AssetStatus, AssetCondition, AssetAssignmentStatus } from "@/features/asset/types";

export const AssetStatusBadge = ({ status }: { status: AssetStatus }) => {
  const styles: Record<AssetStatus, string> = {
    AVAILABLE: "bg-green-100 text-green-800",
    ASSIGNED: "bg-blue-100 text-blue-800",
    MAINTENANCE: "bg-yellow-100 text-yellow-800",
    RETIRED: "bg-gray-100 text-gray-800",
  };

  return (
    <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${styles[status]}`}>
      {status}
    </span>
  );
};

export const AssetConditionBadge = ({ condition }: { condition: AssetCondition }) => {
  const styles: Record<AssetCondition, string> = {
    GOOD: "bg-green-100 text-green-800",
    FAIR: "bg-yellow-100 text-yellow-800",
    DAMAGED: "bg-red-100 text-red-800",
    LOST: "bg-red-200 text-red-900",
  };

  return (
    <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${styles[condition]}`}>
      {condition}
    </span>
  );
};

export const AssetAssignmentStatusBadge = ({ status }: { status: AssetAssignmentStatus }) => {
  const styles: Record<AssetAssignmentStatus, string> = {
    PENDING: "bg-yellow-100 text-yellow-800",
    ACTIVE: "bg-blue-100 text-blue-800",
    RETURNED: "bg-green-100 text-green-800",
    REJECTED: "bg-red-100 text-red-800",
  };

  return (
    <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${styles[status]}`}>
      {status}
    </span>
  );
};
