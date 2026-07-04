import { useState, useRef } from "react";
import { Crown, Loader2, CheckCircle2, XCircle, Search, Eye } from "lucide-react";
import { toast } from "sonner";
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
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  useAllRequests,
  useReviewRequest,
} from "@/features/subscription/hooks/useSubscription";
import type { SubscriptionRequestItem } from "@/features/subscription/types";
import { useDebounce } from "@/hooks/useDebounce";

export default function SubscriptionAdminPage() {
  const { data, isLoading } = useAllRequests();
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
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);
  const [statusFilter, setStatusFilter] = useState<string>("ALL");
  const [detailRequest, setDetailRequest] = useState<SubscriptionRequestItem | null>(null);
  const [isDetailOpen, setIsDetailOpen] = useState(false);

  function cleanupPreview() {
    if (previewUrl) URL.revokeObjectURL(previewUrl);
    setPreviewUrl(null);
    setSelectedFile(null);
  }

  function handleFileSelect(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0];
    if (!file) return;

    if (file.size > 5 * 1024 * 1024) {
      toast.error("File terlalu besar", {
        description: "Ukuran maksimal 5MB.",
      });
      return;
    }

    if (previewUrl) URL.revokeObjectURL(previewUrl);
    setSelectedFile(file);
    setPreviewUrl(URL.createObjectURL(file));
  }

  const allRequests: SubscriptionRequestItem[] = data?.data || [];

  const filtered = allRequests.filter((r) => {
    const matchesSearch =
      !debouncedSearch ||
      r.company_name.toLowerCase().includes(debouncedSearch.toLowerCase()) ||
      r.requested_plan_name.toLowerCase().includes(debouncedSearch.toLowerCase());
    const matchesStatus = statusFilter === "ALL" || r.status === statusFilter;
    return matchesSearch && matchesStatus;
  });

  function openDetail(request: SubscriptionRequestItem) {
    setDetailRequest(request);
    setIsDetailOpen(true);
  }

  function handleReview(
    request: SubscriptionRequestItem,
    action: "APPROVED" | "REJECTED"
  ) {
    setSelectedRequest(request);
    setReviewAction(action);
    setNotes("");
    cleanupPreview();
    setIsDialogOpen(true);
  }

  function confirmReview() {
    if (!selectedRequest) return;

    if (reviewAction === "APPROVED" && (!selectedFile || selectedFile.size > 5 * 1024 * 1024)) {
      toast.error("Bukti pembayaran tidak valid", {
        description: !selectedFile
          ? "Silakan unggah bukti pembayaran sebelum menyetujui."
          : "Ukuran file maksimal 5MB.",
      });
      return;
    }

    reviewRequest(
      {
        id: selectedRequest.id,
        payload: { status: reviewAction, notes },
        file: selectedFile || undefined,
      },
      {
        onSettled: () => {
          cleanupPreview();
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
          {filtered.length} permintaan
        </Badge>
      </div>

      <Card>
        <CardHeader>
          <div className="flex flex-col sm:flex-row justify-between items-center gap-4">
            <CardTitle className="flex items-center gap-2">
              <Crown className="h-5 w-5" />
              Permintaan Upgrade
            </CardTitle>
            <Select value={statusFilter} onValueChange={setStatusFilter}>
              <SelectTrigger className="w-40">
                <SelectValue placeholder="Status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="ALL">Semua</SelectItem>
                <SelectItem value="PENDING">Menunggu</SelectItem>
                <SelectItem value="APPROVED">Disetujui</SelectItem>
                <SelectItem value="REJECTED">Ditolak</SelectItem>
              </SelectContent>
            </Select>
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
                      <TableHead>Status</TableHead>
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
                        <TableCell>
                          <Badge
                            className={
                              req.status === "APPROVED"
                                ? "bg-emerald-100 text-emerald-800 hover:bg-emerald-100"
                                : req.status === "REJECTED"
                                ? "bg-red-100 text-red-800 hover:bg-red-100"
                                : "text-amber-600"
                            }
                            variant={req.status === "PENDING" ? "outline" : undefined}
                          >
                            {req.status === "APPROVED"
                              ? "Disetujui"
                              : req.status === "REJECTED"
                              ? "Ditolak"
                              : "Menunggu"}
                          </Badge>
                        </TableCell>
                        <TableCell className="text-right">
                          <div className="flex justify-end gap-2">
                            {req.status === "PENDING" && (
                              <>
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
                              </>
                            )}
                            <Button
                              size="sm"
                              variant="ghost"
                              className="h-8 px-2 text-slate-500 hover:text-slate-700 hover:bg-slate-100"
                              onClick={() => openDetail(req)}
                            >
                              <Eye className="h-4 w-4 mr-1" />
                              Detail
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
                        className={
                          req.status === "APPROVED"
                            ? "bg-emerald-100 text-emerald-800 hover:bg-emerald-100"
                            : req.status === "REJECTED"
                            ? "bg-red-100 text-red-800 hover:bg-red-100"
                            : "text-amber-600"
                        }
                        variant={req.status === "PENDING" ? "outline" : undefined}
                      >
                        {req.status === "APPROVED"
                          ? "Disetujui"
                          : req.status === "REJECTED"
                          ? "Ditolak"
                          : "Menunggu"}
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
                      {req.status === "PENDING" ? (
                        <>
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
                        </>
                      ) : (
                        <div className="flex-1" />
                      )}
                      <Button
                        size="sm"
                        variant="outline"
                        className="flex-shrink-0"
                        onClick={() => openDetail(req)}
                      >
                        <Eye className="w-4 h-4 mr-1" />
                        Detail
                      </Button>
                    </div>
                  </div>
                ))}
              </div>
            </>
          )}
        </CardContent>
      </Card>

      <Dialog open={isDialogOpen} onOpenChange={(open) => { if (!open) cleanupPreview(); setIsDialogOpen(open); }}>
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
            {reviewAction === "APPROVED" && (
              <div className="space-y-2">
                <label className="text-sm font-medium text-slate-700">
                  Bukti Pembayaran <span className="text-red-500">*</span>
                </label>
                <input
                  ref={fileInputRef}
                  type="file"
                  accept="image/*"
                  onChange={handleFileSelect}
                  className="hidden"
                />
                {!selectedFile ? (
                  <button
                    type="button"
                    onClick={() => fileInputRef.current?.click()}
                    className="w-full rounded-lg border-2 border-dashed border-slate-300 p-6 text-center text-sm text-slate-500 hover:border-blue-400 hover:text-blue-600 transition-colors"
                  >
                    Klik untuk unggah bukti pembayaran
                    <br />
                    <span className="text-xs text-slate-400">Maks. 5MB (JPG, PNG)</span>
                  </button>
                ) : (
                  <div className="relative rounded-lg border border-slate-200 overflow-hidden">
                    <img
                      src={previewUrl!}
                      alt="Preview bukti pembayaran"
                      className="w-full h-40 object-cover"
                    />
                    <button
                      type="button"
                      onClick={() => {
                        cleanupPreview();
                        if (fileInputRef.current) {
                          fileInputRef.current.value = "";
                        }
                      }}
                      className="absolute top-2 right-2 rounded-full bg-white/80 p-1 text-slate-600 hover:bg-white hover:text-red-600 shadow-sm"
                    >
                      <XCircle className="h-4 w-4" />
                    </button>
                  </div>
                )}
              </div>
            )}
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

      <Dialog open={isDetailOpen} onOpenChange={setIsDetailOpen}>
        <DialogContent className="max-w-lg">
          <DialogHeader>
            <DialogTitle>Detail Permintaan</DialogTitle>
          </DialogHeader>

          {detailRequest && (
            <div className="space-y-4 py-2">
              <div className="grid grid-cols-2 gap-3 text-sm">
                <div>
                  <p className="text-slate-500">Perusahaan</p>
                  <p className="font-medium">{detailRequest.company_name}</p>
                </div>
                <div>
                  <p className="text-slate-500">Status</p>
                  <Badge
                    className={
                      detailRequest.status === "APPROVED"
                        ? "bg-emerald-100 text-emerald-800 hover:bg-emerald-100"
                        : detailRequest.status === "REJECTED"
                        ? "bg-red-100 text-red-800 hover:bg-red-100"
                        : "text-amber-600"
                    }
                    variant={detailRequest.status === "PENDING" ? "outline" : undefined}
                  >
                    {detailRequest.status === "APPROVED"
                      ? "Disetujui"
                      : detailRequest.status === "REJECTED"
                      ? "Ditolak"
                      : "Menunggu"}
                  </Badge>
                </div>
                <div>
                  <p className="text-slate-500">Paket Saat Ini</p>
                  <p className="font-medium">{detailRequest.current_plan_name}</p>
                </div>
                <div>
                  <p className="text-slate-500">Paket Diminta</p>
                  <p className="font-medium text-blue-700">{detailRequest.requested_plan_name}</p>
                </div>
                <div>
                  <p className="text-slate-500">Harga Saat Ini</p>
                  <p className="font-medium">
                    {detailRequest.current_plan_price.toLocaleString("id-ID", {
                      style: "currency",
                      currency: "IDR",
                      maximumFractionDigits: 0,
                    })}
                  </p>
                </div>
                <div>
                  <p className="text-slate-500">Harga Diminta</p>
                  <p className="font-medium text-blue-700">
                    {detailRequest.requested_plan_price.toLocaleString("id-ID", {
                      style: "currency",
                      currency: "IDR",
                      maximumFractionDigits: 0,
                    })}
                  </p>
                </div>
                <div>
                  <p className="text-slate-500">Selisih</p>
                  <p className="font-medium text-blue-700">
                    +
                    {detailRequest.price_difference.toLocaleString("id-ID", {
                      style: "currency",
                      currency: "IDR",
                      maximumFractionDigits: 0,
                    })}
                    /bln
                  </p>
                </div>
                <div>
                  <p className="text-slate-500">Tanggal</p>
                  <p className="font-medium">
                    {new Date(detailRequest.created_at).toLocaleDateString("id-ID", {
                      day: "numeric",
                      month: "short",
                      year: "numeric",
                    })}
                  </p>
                </div>
                <div>
                  <p className="text-slate-500">Diminta Oleh</p>
                  <p className="font-medium">{detailRequest.requested_by_name}</p>
                </div>
              </div>

              <div>
                <p className="text-slate-500 text-sm mb-1">Catatan</p>
                <p className="text-sm bg-slate-50 rounded-lg p-3">
                  {detailRequest.notes || "-"}
                </p>
              </div>

              <div>
                <p className="text-slate-500 text-sm mb-2">Bukti Pembayaran</p>
                {detailRequest.proof_payment_url ? (
                  <div className="space-y-2">
                    <a
                      href={detailRequest.proof_payment_url}
                      target="_blank"
                      rel="noopener noreferrer"
                    >
                      <img
                        src={detailRequest.proof_payment_url}
                        alt="Bukti Pembayaran"
                        className="w-full max-h-[300px] object-contain rounded-lg border border-slate-200 hover:opacity-90 transition-opacity cursor-pointer"
                      />
                    </a>
                    <a
                      href={detailRequest.proof_payment_url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="text-xs text-blue-600 hover:text-blue-800 underline inline-block"
                    >
                      Lihat Gambar Penuh
                    </a>
                  </div>
                ) : (
                  <p className="text-sm text-slate-400 italic">
                    Tidak ada bukti pembayaran.
                  </p>
                )}
              </div>
            </div>
          )}

          <DialogFooter>
            <Button variant="outline" onClick={() => setIsDetailOpen(false)}>
              Tutup
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
