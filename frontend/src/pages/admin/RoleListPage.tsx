import { PageHeader } from "@/components/shared/PageHeader";
import { RoleList } from "@/features/role/components/RoleList";
import { ShieldAlert } from "lucide-react";

export default function RoleListPage() {
  return (
    <div className="space-y-6">
      <PageHeader
        title="Roles & Permissions"
        description="Configure access control and permissions for different user roles in the system."
        icon={ShieldAlert}
      />
      <RoleList />
    </div>
  );
}
