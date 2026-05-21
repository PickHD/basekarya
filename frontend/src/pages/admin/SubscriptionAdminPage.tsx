import { useState } from "react";
import { Crown, Loader2, CheckCircle2, XCircle, Search } from "lucide-react";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Textarea } from "@/components/ui/textarea";
import {
  usePendingRequests,
  useReviewRequest,
} from "@/features/subscription/hooks/useSubscription";
import type { SubscriptionRequestItem } from "@/features/subscription/types";
import { useDebounce } from "@/hooks/useDebounce";

export default function SubscriptionAdminPage() {
  const { data, isLoading } = usePendingRequests();
  const { mutate: reviewRequest, isPending: isReviewing } = useReviewRequest();
  const [search, setSearch] = useState("");
  const debouncedSearch = useDebounce(search, 500);

  const [selectedRequest, setSelectedRequest] =
    useState<SubscriptionRequestItem | null>(null);
  const [reviewAction, setReviewAction] = useState<"APPROVED" | "REJECTED">(
    "APPROVED"
  );
  const [notes, setNotes] = useState("");
  const [isDialogOpen, setIsDialogOpen] = useState(false);

  const allRequests: SubscriptionRequestItem[] = data?.data || [];

  const filtered = debouncedSearch
    ? allRequests.filter(
        (r) =>
          r.company_name
            .toLowerCase()
            .includes(debouncedSearch.toLowerCase()) ||
          r.requested_plan_name
            .toLowerCase()
            .includes(debouncedSearch.toLowerCase())
      )
    : allRequests;

  function handleReview(
    request: SubscriptionRequestItem,
    action: "APPROVED" | "REJECTED"
  ) {
    setSelectedRequest(request);
    setReviewAction(action);
    setNotes("");
    setIsDialogOpen(true);
  }

  function confirmReview() {
    if (!selectedRequest) return;
    reviewRequest(
      { id: selectedRequest.id, payload: { status: reviewAction, notes } },
      {
        onSettled: () => {
          setIsDialogOpen(false);
          setSelectedRequest(null);
        },
      }
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">
            Manajemen Langganan
          </h2>
          <p className="text-slate-500 mt-1">
            Kelola permintaan upgrade dan pembayaran langganan perusahaan.
          </p>
        </div>
        <Badge variant="outline" className="text-sm px-3 py-1">
          {allRequests.length} permintaan menunggu
        </Badge>
      </div>

      <Card>
        <CardHeader>
          <div className="flex flex-col sm:flex-row justify-between items-center gap-4">
            <CardTitle className="flex items-center gap-2">
              <Crown className="h-5 w-5" />
              Permintaan Upgrade
            </CardTitle>
            <div className="relative w-full sm:w-64">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-slate-400" />
              <Input
                placeholder="Cari perusahaan..."
                className="pl-9"
                value={search}
                onChange={(e) => setSearch(e.target.value)}
              />
            </div>
          </div>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="flex justify-center py-10">
              <Loader2 className="animate-spin h-8 w-8 text-blue-600" />
            </div>
          ) : filtered.length === 0 ? (
            <div className="text-center py-10 text-slate-500">
              {allRequests.length === 0
                ? "Tidak ada permintaan upgrade yang menunggu."
                : "Tidak ada hasil yang cocok."}
            </div>
          ) : (
            <>
              <div className="hidden md:block">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Perusahaan</TableHead>
                      <TableHead>Paket Saat Ini</TableHead>
                      <TableHead>Diminta</TableHead>
                      <TableHead>Selisih Harga</TableHead>
                      <TableHead>Tanggal</TableHead>
                      <TableHead className="text-right">Aksi</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {filtered.map((req) => (
                      <TableRow key={req.id}>
                        <TableCell className="font-medium">
                          {req.company_name}
                        </TableCell>
                        <TableCell>
                          <Badge variant="outline">{req.current_plan_name}</Badge>
                        </TableCell>
                        <TableCell>
                          <Badge className="bg-blue-100 text-blue-800 hover:bg-blue-100">
                            {req.requested_plan_name}
                          </Badge>
                        </TableCell>
                        <TableCell className="font-medium text-blue-700">
                          +
                          {req.price_difference.toLocaleString("id-ID", {
                            style: "currency",
                            currency: "IDR",
                            maximumFractionDigits: 0,
                          })}
                          /bln
                        </TableCell>
                        <TableCell className="text-slate-500">
                          {new Date(req.created_at).toLocaleDateString("id-ID", {
                            day: "numeric",
                            month: "short",
                            year: "numeric",
                          })}
                        </TableCell>
                        <TableCell className="text-right">
                          <div className="flex justify-end gap-2">
                            <Button
                              size="sm"
                              variant="ghost"
                              className="h-8 px-2 text-emerald-600 hover:text-emerald-700 hover:bg-emerald-50"
                              onClick={() => handleReview(req, "APPROVED")}
                              disabled={isReviewing}
                            >
                              <CheckCircle2 className="h-4 w-4 mr-1" />
                              Setujui
                            </Button>
                            <Button
                              size="sm"
                              variant="ghost"
                              className="h-8 px-2 text-red-600 hover:text-red-700 hover:bg-red-50"
                              onClick={() => handleReview(req, "REJECTED")}
                              disabled={isReviewing}
                            >
                              <XCircle className="h-4 w-4 mr-1" />
                              Tolak
                            </Button>
                          </div>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </div>

              <div className="md:hidden space-y-4">
                {filtered.map((req) => (
                  <div
                    key={req.id}
                    className="rounded-lg border p-4 space-y-3"
                  >
                    <div className="flex items-center justify-between">
                      <span className="font-medium">{req.company_name}</span>
                      <Badge
                        variant="outline"
                        className="text-amber-600 capitalize"
                      >
                        Menunggu
                      </Badge>
                    </div>
                    <div className="grid grid-cols-2 gap-2 text-sm">
                      <div>
                        <p className="text-slate-500">Saat Ini</p>
                        <Badge variant="outline">{req.current_plan_name}</Badge>
                      </div>
                      <div>
                        <p className="text-slate-500">Diminta</p>
                        <Badge className="bg-blue-100 text-blue-800 hover:bg-blue-100">
                          {req.requested_plan_name}
                        </Badge>
                      </div>
                    </div>
                    <div className="flex items-center justify-between text-sm">
                      <span className="text-slate-500">Selisih:</span>
                      <span className="font-medium text-blue-700">
                        +
                        {req.price_difference.toLocaleString("id-ID", {
                          style: "currency",
                          currency: "IDR",
                          maximumFractionDigits: 0,
                        })}
                        /bln
                      </span>
                    </div>
                    <div className="flex gap-2 pt-2 border-t">
                      <Button
                        size="sm"
                        className="flex-1"
                        onClick={() => handleReview(req, "APPROVED")}
                        disabled={isReviewing}
                      >
                        <CheckCircle2 className="w-4 h-4 mr-1" />
                        Setujui
                      </Button>
                      <Button
                        size="sm"
                        variant="destructive"
                        className="flex-1"
                        onClick={() => handleReview(req, "REJECTED")}
                        disabled={isReviewing}
                      >
                        <XCircle className="w-4 h-4 mr-1" />
                        Tolak
                      </Button>
                    </div>
                  </div>
                ))}
              </div>
            </>
          )}
        </CardContent>
      </Card>

      <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>
              {reviewAction === "APPROVED"
                ? "Setujui Permintaan Upgrade"
                : "Tolak Permintaan Upgrade"}
            </DialogTitle>
            <DialogDescription>
              {selectedRequest && (
                <>
                  <span className="font-medium">
                    {selectedRequest.company_name}
                  </span>{" "}
                  — {selectedRequest.current_plan_name} →{" "}
                  {selectedRequest.requested_plan_name}
                </>
              )}
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-4 py-2">
            {reviewAction === "APPROVED" && (
              <div className="rounded-lg border border-emerald-200 bg-emerald-50 p-3 text-sm text-emerald-800">
                Perusahaan akan langsung di-upgrade ke paket{" "}
                <strong>{selectedRequest?.requested_plan_name}</strong>.
              </div>
            )}
            <Textarea
              placeholder="Catatan (opsional)..."
              value={notes}
              onChange={(e) => setNotes(e.target.value)}
              rows={3}
            />
          </div>

          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => setIsDialogOpen(false)}
              disabled={isReviewing}
            >
              Batal
            </Button>
            <Button
              variant={
                reviewAction === "APPROVED" ? "default" : "destructive"
              }
              onClick={confirmReview}
              disabled={isReviewing}
            >
              {isReviewing ? (
                <Loader2 className="h-4 w-4 animate-spin mr-1" />
              ) : reviewAction === "APPROVED" ? (
                <CheckCircle2 className="h-4 w-4 mr-1" />
              ) : (
                <XCircle className="h-4 w-4 mr-1" />
              )}
              {reviewAction === "APPROVED" ? "Setujui" : "Tolak"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
