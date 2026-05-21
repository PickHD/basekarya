import { useAllRequests, useReviewRequest } from "@/features/subscription/hooks/useSubscription";
import type { SubscriptionRequestItem } from "@/features/subscription/types";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Loader2 } from "lucide-react";
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
        <h1 className="text-2xl font-bold text-foreground">
          Riwayat Permintaan
        </h1>
        <p className="text-muted-foreground mt-1">
          Semua permintaan upgrade subscription
        </p>
      </div>

      {requests.length === 0 ? (
        <div className="text-center py-12 text-muted-foreground">
          Belum ada permintaan
        </div>
      ) : (
        <div className="rounded-lg border">
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b bg-muted/50">
                  <th className="text-left p-3 font-medium">Perusahaan</th>
                  <th className="text-left p-3 font-medium">Dari</th>
                  <th className="text-left p-3 font-medium">Ke</th>
                  <th className="text-left p-3 font-medium">Selisih</th>
                  <th className="text-left p-3 font-medium">Status</th>
                  <th className="text-left p-3 font-medium">Tanggal</th>
                  <th className="text-left p-3 font-medium">Aksi</th>
                </tr>
              </thead>
              <tbody>
                {requests.map((req) => (
                  <tr key={req.id} className="border-b last:border-0 hover:bg-muted/30">
                    <td className="p-3">
                      <div className="font-medium">{req.company_name}</div>
                      <div className="text-xs text-muted-foreground">
                        {req.requested_by_name}
                      </div>
                    </td>
                    <td className="p-3">{req.current_plan_name}</td>
                    <td className="p-3 font-medium">{req.requested_plan_name}</td>
                    <td className="p-3 text-xs">
                      +Rp{req.price_difference?.toLocaleString("id-ID")}
                    </td>
                    <td className="p-3">
                      <StatusBadge status={req.status} />
                    </td>
                    <td className="p-3 text-xs text-muted-foreground">
                      {new Date(req.created_at).toLocaleDateString("id-ID", {
                        day: "numeric",
                        month: "short",
                        year: "numeric",
                      })}
                    </td>
                    <td className="p-3">
                      {req.status === "PENDING" && (
                        <div className="flex gap-1">
                          <Button
                            size="sm"
                            variant="outline"
                            className="text-xs h-7 text-emerald-600 hover:text-emerald-700"
                            onClick={() => handleReview(req, "APPROVED")}
                          >
                            Setujui
                          </Button>
                          <Button
                            size="sm"
                            variant="outline"
                            className="text-xs h-7 text-red-600 hover:text-red-700"
                            onClick={() => handleReview(req, "REJECTED")}
                          >
                            Tolak
                          </Button>
                        </div>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}

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
