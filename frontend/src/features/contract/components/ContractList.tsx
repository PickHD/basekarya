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

  return (
    <>
      <div className="grid grid-cols-1 gap-4 md:hidden">
        {isLoading ? (
          <div className="flex justify-center p-8">
            <Loader2 className="h-8 w-8 animate-spin text-slate-400" />
          </div>
        ) : data.length === 0 ? (
          <div className="text-center p-8 text-slate-500 border rounded-lg bg-slate-50">
            No contracts found.
          </div>
        ) : (
          data.map((item) => (
            <div
              key={item.id}
              className="flex flex-col rounded-lg border bg-card p-4 shadow-sm space-y-3"
            >
              <div className="flex justify-between items-start">
                <div>
                  <h4 className="font-bold text-slate-800">
                    {item.employee_name || "-"}
                  </h4>
                  <p className="text-xs text-muted-foreground">{item.employee_nik || "-"}</p>
                </div>
                <ContractTypeBadge type={item.contract_type} />
              </div>

              <div className="flex items-center gap-2 text-sm text-slate-700">
                <FileText className="h-4 w-4 text-slate-400" />
                <span className="font-medium">
                  {item.contract_number || "-"}
                </span>
              </div>

              <div className="bg-slate-50 p-3 rounded text-sm grid grid-cols-2 gap-2 text-center">
                <div>
                  <div className="text-xs text-slate-500">Start</div>
                  <div className="font-medium">
                    {format(new Date(item.start_date), "dd MMM yyyy")}
                  </div>
                </div>
                <div>
                  <div className="text-xs text-slate-500">End</div>
                  <div className="font-medium">
                    {item.end_date ? (
                      <span className={new Date(item.end_date) < new Date() ? "text-red-500 font-medium" : ""}>
                        {format(new Date(item.end_date), "dd MMM yyyy")}
                      </span>
                    ) : (
                      <span className="text-muted-foreground">-</span>
                    )}
                  </div>
                </div>
              </div>

              <div className="flex gap-2 w-full pt-2">
                <Button
                  variant="outline"
                  size="sm"
                  className="flex-1"
                  onClick={() => onView(item)}
                >
                  <Eye className="mr-2 h-4 w-4" /> View
                </Button>
                {hasPermission(PERMISSIONS.CREATE_CONTRACT) && (
                  <Button
                    variant="outline"
                    size="sm"
                    className="flex-1"
                    onClick={() => onEdit(item)}
                  >
                    <FileText className="mr-2 h-4 w-4" /> Edit
                  </Button>
                )}
                {hasPermission(PERMISSIONS.UPDATE_CONTRACT) && (
                  <Button
                    variant="outline"
                    size="icon"
                    className="flex-none text-red-600 hover:text-red-700 hover:bg-red-50"
                    onClick={() => setDeleteId(item.id)}
                  >
                    <Trash className="h-4 w-4" />
                  </Button>
                )}
              </div>
            </div>
          ))
        )}
      </div>

      <div className="hidden md:block rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Employee</TableHead>
              <TableHead>Type</TableHead>
              <TableHead>Contract Number</TableHead>
              <TableHead>Start Date</TableHead>
              <TableHead>End Date</TableHead>
              <TableHead className="text-right">Action</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoading ? (
              <TableRow>
                <TableCell colSpan={6} className="h-24 text-center">
                  <div className="flex justify-center items-center gap-2 text-slate-500">
                    <Loader2 className="h-5 w-5 animate-spin" /> Loading data...
                  </div>
                </TableCell>
              </TableRow>
            ) : data.length === 0 ? (
              <TableRow>
                <TableCell colSpan={6} className="text-center py-8 text-slate-500">
                  No contracts found.
                </TableCell>
              </TableRow>
            ) : (
              data.map((item) => (
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
              ))
            )}
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
