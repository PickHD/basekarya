import { useParams, useNavigate } from "react-router-dom";
import { ArrowLeft, Users } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { ApplicantKanbanBoard } from "@/features/recruitment/components/ApplicantKanbanBoard";
import { useRequisitionDetail } from "@/features/recruitment/hooks/useRequisition";

export default function ApplicantBoardPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const requisitionId = Number(id);

  const { data: requisition, isLoading } = useRequisitionDetail(requisitionId || null);

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex flex-col sm:flex-row items-start sm:items-center gap-4">
        <Button
          variant="outline"
          size="sm"
          onClick={() => navigate("/admin/requisitions")}
          className="flex-shrink-0"
        >
          <ArrowLeft className="mr-2 h-4 w-4" />
          Back
        </Button>
        <div>
          <h2 className="text-3xl font-bold tracking-tight">
            {isLoading ? "Loading..." : requisition?.title ?? "Applicant Board"}
          </h2>
          <p className="text-slate-500 mt-1">
            {requisition
              ? `${requisition.department_name} · ${requisition.employment_type} · ${requisition.quantity} position(s)`
              : "Drag & drop applicants between stages"}
          </p>
        </div>
      </div>

      {/* Kanban Board Card */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Users className="h-5 w-5" />
            Applicant Tracking Board
          </CardTitle>
        </CardHeader>
        <CardContent>
          {requisitionId ? (
            <ApplicantKanbanBoard requisitionId={requisitionId} />
          ) : (
            <p className="text-slate-500 text-center py-8">Invalid requisition ID</p>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
