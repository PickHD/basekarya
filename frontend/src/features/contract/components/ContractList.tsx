import { useState } from "react";
import { format } from "date-fns";
import { MoreHorizontal, FileText, Trash, Eye, Loader2 } from "lucide-react";
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
import { ContractTypeBadge } from "./ContractTypeBadge";
import type { Contract } from "../types";
import { usePermissions } from "@/hooks/usePermissions";
import { PERMISSIONS } from "@/config/permissions";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { useDeleteContract } from "../hooks/useContract";

interface Props {
  data: Contract[];
  isLoading: boolean;
  onView: (contract: Contract) => void;
  onEdit: (contract: Contract) => void;
}

export function ContractList({ data, isLoading, onView, onEdit }: Props) {
  const { hasPermission } = usePermissions();
  const { mutate: deleteContract, isPending: isDeleting } = useDeleteContract();
  const [deleteId, setDeleteId] = useState<number | null>(null);

  const handleDelete = () => {
    if (deleteId) {
      deleteContract(deleteId, {
        onSuccess: () => setDeleteId(null),
      });
    }
  };

  if (isLoading) {
    return (
      <div className="flex justify-center p-8">
        <Loader2 className="h-8 w-8 animate-spin text-slate-400" />
      </div>
    );
  }

  if (data.length === 0) {
    return (
      <div className="text-center p-8 text-slate-500 border rounded-lg bg-slate-50">
        No contracts found.
      </div>
    );
  }

  return (
    <>
      <div className="rounded-md border bg-white">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Employee</TableHead>
              <TableHead>Type</TableHead>
              <TableHead>Contract Number</TableHead>
              <TableHead>Start Date</TableHead>
              <TableHead>End Date</TableHead>
              <TableHead className="text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {data.map((item) => (
              <TableRow key={item.id}>
                <TableCell>
                  <p className="font-medium">{item.employee_name || "-"}</p>
                  <p className="text-xs text-muted-foreground">{item.employee_nik || "-"}</p>
                </TableCell>
                <TableCell>
                  <ContractTypeBadge type={item.contract_type} />
                </TableCell>
                <TableCell>{item.contract_number || "-"}</TableCell>
                <TableCell>{format(new Date(item.start_date), "dd MMM yyyy")}</TableCell>
                <TableCell>
                  {item.end_date ? (
                    <span className={new Date(item.end_date) < new Date() ? "text-red-500 font-medium" : ""}>
                      {format(new Date(item.end_date), "dd MMM yyyy")}
                    </span>
                  ) : (
                    <span className="text-muted-foreground">-</span>
                  )}
                </TableCell>
                <TableCell className="text-right">
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button variant="ghost" className="h-8 w-8 p-0">
                        <MoreHorizontal className="h-4 w-4" />
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                      <DropdownMenuItem onClick={() => onView(item)}>
                        <Eye className="mr-2 h-4 w-4" /> Detail
                      </DropdownMenuItem>
                      {hasPermission(PERMISSIONS.CREATE_CONTRACT) && (
                         <DropdownMenuItem onClick={() => onEdit(item)}>
                           <FileText className="mr-2 h-4 w-4" /> Edit
                         </DropdownMenuItem>
                      )}
                      {hasPermission(PERMISSIONS.UPDATE_CONTRACT) && (
                        <DropdownMenuItem
                          className="text-red-600"
                          onClick={() => setDeleteId(item.id)}
                        >
                          <Trash className="mr-2 h-4 w-4" /> Delete
                        </DropdownMenuItem>
                      )}
                    </DropdownMenuContent>
                  </DropdownMenu>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>

      <AlertDialog open={!!deleteId} onOpenChange={() => setDeleteId(null)}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete Contract?</AlertDialogTitle>
            <AlertDialogDescription>
              This action cannot be undone. The contract data will be deleted from the system.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel disabled={isDeleting}>Cancel</AlertDialogCancel>
            <AlertDialogAction
              className="bg-red-600 hover:bg-red-700"
              onClick={(e) => {
                e.preventDefault();
                handleDelete();
              }}
              disabled={isDeleting}
            >
              {isDeleting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}
