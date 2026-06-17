import { useState } from "react";
import { useCompanies, useUpdateCompanyStatus } from "@/features/subscription/hooks/useSubscription";
import type { CompanyListItem } from "@/features/subscription/types";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Search,
  Building2,
  Mail,
  Phone,
  Calendar,
  Users,
  Crown,
} from "lucide-react";
import { cn } from "@/lib/utils";
import { Loader2 } from "lucide-react";

function StatusBadge({ status }: { status: string }) {
  return (
    <Badge
      variant="outline"
      className={cn(
        "capitalize text-xs",
        status === "ACTIVE"
          ? "text-emerald-600 border-emerald-200 bg-emerald-50"
          : status === "PENDING_PAYMENT"
          ? "text-amber-600 border-amber-200 bg-amber-50"
          : "text-red-600 border-red-200 bg-red-50"
      )}
    >
      {status === "ACTIVE"
        ? "Aktif"
        : status === "PENDING_PAYMENT"
        ? "Menunggu Bayar"
        : status === "EXPIRED"
        ? "Kadaluarsa"
        : status}
    </Badge>
  );
}

export default function CompaniesPage() {
  const [search, setSearch] = useState("");
  const [selectedCompany, setSelectedCompany] = useState<CompanyListItem | null>(null);
  const [newStatus, setNewStatus] = useState("");
  const [isDialogOpen, setIsDialogOpen] = useState(false);

  const { data, isLoading } = useCompanies(search);
  const updateStatus = useUpdateCompanyStatus();

  const companies = data?.data || [];

  function handleStatusChange(company: CompanyListItem) {
    setSelectedCompany(company);
    setNewStatus(company.subscription_status);
    setIsDialogOpen(true);
  }

  function confirmStatusChange() {
    if (!selectedCompany || !newStatus) return;
    updateStatus.mutate(
      { id: selectedCompany.id, status: newStatus },
      {
        onSuccess: () => {
          setIsDialogOpen(false);
          setSelectedCompany(null);
        },
      }
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">Perusahaan</h2>
          <p className="text-slate-500">
            Kelola semua perusahaan terdaftar
          </p>
        </div>
      </div>

      <div className="relative max-w-sm">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
        <Input
          placeholder="Cari nama atau email..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="pl-9"
        />
      </div>

      {isLoading ? (
        <div className="flex justify-center py-12">
          <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
        </div>
      ) : companies.length === 0 ? (
        <div className="text-center py-12 text-muted-foreground">
          Tidak ada perusahaan ditemukan
        </div>
      ) : (
        <Card>
          <CardHeader className="pb-3">
            <CardTitle className="text-lg font-semibold">All Companies</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
              {companies.map((company: CompanyListItem) => (
                <div
                  key={company.id}
                  className="rounded-lg border bg-card p-4 space-y-3 hover:shadow-md transition-shadow cursor-pointer"
                  onClick={() => handleStatusChange(company)}
                >
                  <div className="flex items-start justify-between">
                    <div className="flex items-center gap-2">
                      <div className="p-2 rounded-lg bg-primary/10">
                        <Building2 className="h-4 w-4 text-primary" />
                      </div>
                      <div>
                        <p className="font-semibold text-sm">{company.name}</p>
                        <p className="text-xs text-muted-foreground">{company.plan_name}</p>
                      </div>
                    </div>
                    <StatusBadge status={company.subscription_status} />
                  </div>

                  <div className="space-y-1.5 text-xs text-muted-foreground">
                    <div className="flex items-center gap-2">
                      <Mail className="h-3 w-3" />
                      {company.email}
                    </div>
                    <div className="flex items-center gap-2">
                      <Phone className="h-3 w-3" />
                      {company.phone_number}
                    </div>
                    <div className="flex items-center gap-2">
                      <Users className="h-3 w-3" />
                      {company.employee_count} karyawan
                    </div>
                    <div className="flex items-center gap-2">
                      <Calendar className="h-3 w-3" />
                      Terdaftar: {new Date(company.created_at).toLocaleDateString("id-ID")}
                    </div>
                  </div>

                  {company.subscription_expires_at && (
                    <div className="flex items-center gap-2 text-xs text-muted-foreground">
                      <Crown className="h-3 w-3" />
                      Exp: {new Date(company.subscription_expires_at).toLocaleDateString("id-ID")}
                    </div>
                  )}
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      )}

      <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Ubah Status - {selectedCompany?.name}</DialogTitle>
          </DialogHeader>

          {selectedCompany && (
            <div className="space-y-4">
              <div className="rounded-lg bg-muted p-3 text-sm space-y-1">
                <p>
                  <span className="text-muted-foreground">Paket:</span>{" "}
                  {selectedCompany.plan_name}
                </p>
                <p>
                  <span className="text-muted-foreground">Email:</span>{" "}
                  {selectedCompany.email}
                </p>
                <p>
                  <span className="text-muted-foreground">Karyawan:</span>{" "}
                  {selectedCompany.employee_count}
                </p>
              </div>

              <Select value={newStatus} onValueChange={setNewStatus}>
                <SelectTrigger>
                  <SelectValue placeholder="Pilih status" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="ACTIVE">Aktif</SelectItem>
                  <SelectItem value="PENDING_PAYMENT">Menunggu Pembayaran</SelectItem>
                  <SelectItem value="EXPIRED">Kadaluarsa</SelectItem>
                </SelectContent>
              </Select>
            </div>
          )}

          <DialogFooter>
            <Button variant="outline" onClick={() => setIsDialogOpen(false)}>
              Batal
            </Button>
            <Button
              onClick={confirmStatusChange}
              disabled={updateStatus.isPending || newStatus === selectedCompany?.subscription_status}
            >
              {updateStatus.isPending && (
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              )}
              Simpan
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
