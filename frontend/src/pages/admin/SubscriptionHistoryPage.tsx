import { useAllRequests, useReviewRequest } from "@/features/subscription/hooks/useSubscription";
import type { SubscriptionRequestItem } from "@/features/subscription/types";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Loader2, Check, X } from "lucide-react";
import { cn } from "@/lib/utils";
import { useState } from "react";

function StatusBadge({ status }: { status: string }) {
  return (
    <Badge
      variant="outline"
      className={cn(
        "capitalize text-xs",
        status === "PENDING"
          ? "text-amber-600 border-amber-200 bg-amber-50"
          : status === "APPROVED"
          ? "text-emerald-600 border-emerald-200 bg-emerald-50"
          : "text-red-600 border-red-200 bg-red-50"
      )}
    >
      {status === "PENDING"
        ? "Menunggu"
        : status === "APPROVED"
        ? "Disetujui"
        : "Ditolak"}
    </Badge>
  );
}

export default function SubscriptionHistoryPage() {
  const { data, isLoading } = useAllRequests();
  const reviewRequest = useReviewRequest();

  const [selectedReq, setSelectedReq] = useState<SubscriptionRequestItem | null>(null);
  const [reviewAction, setReviewAction] = useState<"APPROVED" | "REJECTED">("APPROVED");
  const [notes, setNotes] = useState("");
  const [isDialogOpen, setIsDialogOpen] = useState(false);

  const requests = data?.data || [];

  function handleReview(req: SubscriptionRequestItem, action: "APPROVED" | "REJECTED") {
    setSelectedReq(req);
    setReviewAction(action);
    setNotes("");
    setIsDialogOpen(true);
  }

  function confirmReview() {
    if (!selectedReq) return;
    reviewRequest.mutate(
      { id: selectedReq.id, payload: { status: reviewAction, notes } },
      {
        onSuccess: () => {
          setIsDialogOpen(false);
          setSelectedReq(null);
        },
      }
    );
  }

  if (isLoading) {
    return (
      <div className="flex justify-center py-12">
        <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-3xl font-bold tracking-tight">Riwayat Permintaan</h2>
        <p className="text-slate-500">
          Semua permintaan upgrade subscription
        </p>
      </div>

      <Card>
        <CardHeader className="pb-3">
          <CardTitle className="text-lg font-semibold">Request History</CardTitle>
        </CardHeader>
        <CardContent>
          {requests.length === 0 ? (
            <div className="text-center py-12 text-muted-foreground">
              Belum ada permintaan
            </div>
          ) : (
            <>
              <div className="hidden md:block">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Perusahaan</TableHead>
                      <TableHead>Dari</TableHead>
                      <TableHead>Ke</TableHead>
                      <TableHead>Selisih</TableHead>
                      <TableHead>Status</TableHead>
                      <TableHead>Tanggal</TableHead>
                      <TableHead className="text-right">Aksi</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {requests.map((req: SubscriptionRequestItem) => (
                      <TableRow key={req.id}>
                        <TableCell>
                          <div className="font-medium">{req.company_name}</div>
                          <div className="text-xs text-muted-foreground">
                            {req.requested_by_name}
                          </div>
                        </TableCell>
                        <TableCell>{req.current_plan_name}</TableCell>
                        <TableCell className="font-medium">{req.requested_plan_name}</TableCell>
                        <TableCell className="text-xs">
                          +Rp{req.price_difference?.toLocaleString("id-ID")}
                        </TableCell>
                        <TableCell>
                          <StatusBadge status={req.status} />
                        </TableCell>
                        <TableCell className="text-xs text-muted-foreground">
                          {new Date(req.created_at).toLocaleDateString("id-ID", {
                            day: "numeric",
                            month: "short",
                            year: "numeric",
                          })}
                        </TableCell>
                        <TableCell className="text-right">
                          {req.status === "PENDING" && (
                            <div className="flex justify-end gap-1">
                              <Button
                                variant="ghost"
                                size="icon"
                                className="h-8 w-8"
                                onClick={() => handleReview(req, "APPROVED")}
                              >
                                <Check className="h-4 w-4 text-emerald-600" />
                              </Button>
                              <Button
                                variant="ghost"
                                size="icon"
                                className="h-8 w-8"
                                onClick={() => handleReview(req, "REJECTED")}
                              >
                                <X className="h-4 w-4 text-red-600" />
                              </Button>
                            </div>
                          )}
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </div>

              <div className="md:hidden space-y-3">
                {requests.map((req: SubscriptionRequestItem) => (
                  <Card key={req.id} className="p-4">
                    <div className="flex items-start justify-between">
                      <div>
                        <p className="font-medium">{req.company_name}</p>
                        <p className="text-sm text-muted-foreground">
                          {req.requested_by_name}
                        </p>
                      </div>
                      <StatusBadge status={req.status} />
                    </div>
                    <div className="mt-3 space-y-1 text-sm">
                      <p className="text-muted-foreground">
                        Dari: {req.current_plan_name} → Ke: {req.requested_plan_name}
                      </p>
                      <p>+Rp{req.price_difference?.toLocaleString("id-ID")}</p>
                      <p className="text-xs text-muted-foreground">
                        {new Date(req.created_at).toLocaleDateString("id-ID", {
                          day: "numeric",
                          month: "short",
                          year: "numeric",
                        })}
                      </p>
                    </div>
                    {req.status === "PENDING" && (
                      <div className="flex gap-1 mt-3">
                        <Button
                          variant="ghost"
                          size="icon"
                          className="h-8 w-8"
                          onClick={() => handleReview(req, "APPROVED")}
                        >
                          <Check className="h-4 w-4 text-emerald-600" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon"
                          className="h-8 w-8"
                          onClick={() => handleReview(req, "REJECTED")}
                        >
                          <X className="h-4 w-4 text-red-600" />
                        </Button>
                      </div>
                    )}
                  </Card>
                ))}
              </div>
            </>
          )}
        </CardContent>
      </Card>

      <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>
              {reviewAction === "APPROVED"
                ? "Setujui Permintaan Upgrade"
                : "Tolak Permintaan Upgrade"}
            </DialogTitle>
          </DialogHeader>

          {selectedReq && (
            <div className="space-y-3 text-sm">
              <div className="rounded-lg bg-muted p-3 space-y-1">
                <p>
                  <span className="text-muted-foreground">Perusahaan:</span>{" "}
                  {selectedReq.company_name}
                </p>
                <p>
                  <span className="text-muted-foreground">Paket:</span>{" "}
                  {selectedReq.current_plan_name} → {selectedReq.requested_plan_name}
                </p>
              </div>
              <Input
                placeholder="Catatan (opsional)"
                value={notes}
                onChange={(e) => setNotes(e.target.value)}
              />
            </div>
          )}

          <DialogFooter>
            <Button variant="outline" onClick={() => setIsDialogOpen(false)}>
              Batal
            </Button>
            <Button
              variant={reviewAction === "APPROVED" ? "default" : "destructive"}
              onClick={confirmReview}
              disabled={reviewRequest.isPending}
            >
              {reviewRequest.isPending && (
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              )}
              {reviewAction === "APPROVED" ? "Setujui" : "Tolak"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
