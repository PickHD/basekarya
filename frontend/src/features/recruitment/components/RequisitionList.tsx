import { format } from "date-fns";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import {
  MoreHorizontal,
  Eye,
  Loader2,
  Users,
  Briefcase,
  Calendar,
  Building2,
} from "lucide-react";
import { RequisitionStatusBadge } from "./RequisitionStatusBadge";
import { PriorityBadge } from "./PriorityBadge";
import type { JobRequisition } from "../types";

interface Props {
  data: JobRequisition[];
  isLoading: boolean;
  onView: (requisition: JobRequisition) => void;
}

export function RequisitionList({ data, isLoading, onView }: Props) {
  return (
    <>
      {/* Desktop Table */}
      <div className="hidden md:block">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Job Title</TableHead>
              <TableHead>Department</TableHead>
              <TableHead>Type</TableHead>
              <TableHead>Priority</TableHead>
              <TableHead>Qty</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Requester</TableHead>
              <TableHead>Created</TableHead>
              <TableHead className="text-right">Action</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoading ? (
              <TableRow>
                <TableCell colSpan={9} className="text-center py-8">
                  <Loader2 className="h-5 w-5 animate-spin mx-auto text-slate-400" />
                </TableCell>
              </TableRow>
            ) : data.length === 0 ? (
              <TableRow>
                <TableCell colSpan={9} className="text-center py-8 text-slate-500">
                  No requisitions found
                </TableCell>
              </TableRow>
            ) : (
              data.map((req) => (
                <TableRow key={req.id}>
                  <TableCell className="font-medium max-w-[200px] truncate">{req.title}</TableCell>
                  <TableCell className="text-slate-600 text-sm">{req.department_name}</TableCell>
                  <TableCell>
                    <span className="text-xs font-medium px-2 py-0.5 bg-slate-100 text-slate-700 rounded">
                      {req.employment_type}
                    </span>
                  </TableCell>
                  <TableCell>
                    <PriorityBadge priority={req.priority} />
                  </TableCell>
                  <TableCell className="text-center">{req.quantity}</TableCell>
                  <TableCell>
                    <RequisitionStatusBadge status={req.status} />
                  </TableCell>
                  <TableCell className="text-slate-600 text-sm">{req.requester_name || "-"}</TableCell>
                  <TableCell className="text-slate-500 text-sm">
                    {format(new Date(req.created_at), "dd MMM yyyy")}
                  </TableCell>
                  <TableCell className="text-right">
                    <DropdownMenu>
                      <DropdownMenuTrigger asChild>
                        <Button variant="ghost" className="h-8 w-8 p-0">
                          <MoreHorizontal className="h-4 w-4" />
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end">
                        <DropdownMenuItem onClick={() => onView(req)}>
                          <Eye className="mr-2 h-4 w-4" />
                          View Detail
                        </DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>

      {/* Mobile Cards */}
      {isLoading ? (
        <div className="flex justify-center py-8 md:hidden">
          <Loader2 className="h-5 w-5 animate-spin text-slate-400" />
        </div>
      ) : data.length === 0 ? (
        <div className="text-center py-8 text-slate-500 md:hidden">No requisitions found</div>
      ) : (
        <div className="grid grid-cols-1 gap-3 md:hidden">
          {data.map((req) => (
            <Card key={req.id} className="cursor-pointer hover:shadow-md transition-shadow" onClick={() => onView(req)}>
              <CardContent className="p-4 space-y-3">
                <div className="flex items-start justify-between gap-2">
                  <p className="font-semibold text-sm leading-tight">{req.title}</p>
                  <RequisitionStatusBadge status={req.status} />
                </div>
                <div className="flex flex-wrap gap-3 text-xs text-slate-600">
                  <span className="flex items-center gap-1">
                    <Building2 className="h-3 w-3" />
                    {req.department_name}
                  </span>
                  <span className="flex items-center gap-1">
                    <Briefcase className="h-3 w-3" />
                    {req.employment_type}
                  </span>
                  <span className="flex items-center gap-1">
                    <Users className="h-3 w-3" />
                    {req.quantity} pax
                  </span>
                  <span className="flex items-center gap-1">
                    <Calendar className="h-3 w-3" />
                    {format(new Date(req.created_at), "dd MMM yyyy")}
                  </span>
                </div>
                <div className="flex items-center gap-2">
                  <PriorityBadge priority={req.priority} />
                  <span className="text-xs text-slate-500">by {req.requester_name || "-"}</span>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </>
  );
}
