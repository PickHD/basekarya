import { useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Loader2, TrendingUp, TrendingDown, Wallet, Activity } from "lucide-react";
import { useFinanceDashboard } from "@/features/finance/hooks/useFinanceDashboard";
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  ArcElement,
  Title,
  Tooltip,
  Legend,
} from "chart.js";
import { Bar, Doughnut } from "react-chartjs-2";

ChartJS.register(CategoryScale, LinearScale, BarElement, ArcElement, Title, Tooltip, Legend);

export const FinanceDashboardView = () => {
  const [startDate, setStartDate] = useState("");
  const [endDate, setEndDate] = useState("");

  const { data, isLoading } = useFinanceDashboard(
    startDate || undefined,
    endDate || undefined
  );

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat("id-ID", {
      style: "currency",
      currency: "IDR",
      minimumFractionDigits: 0,
    }).format(amount);
  };

  const handleFilter = () => {
    setStartDate(startDate);
    setEndDate(endDate);
  };

  const handleReset = () => {
    setStartDate("");
    setEndDate("");
  };

  const monthlyChartData = {
    labels: (data?.monthly_summary || []).map((item) => item.month),
    datasets: [
      {
        label: "Pemasukan",
        data: (data?.monthly_summary || []).map((item) => item.income),
        backgroundColor: "rgba(34, 197, 94, 0.7)",
        borderColor: "rgba(34, 197, 94, 1)",
        borderWidth: 1,
        borderRadius: 4,
      },
      {
        label: "Pengeluaran",
        data: (data?.monthly_summary || []).map((item) => item.expense),
        backgroundColor: "rgba(239, 68, 68, 0.7)",
        borderColor: "rgba(239, 68, 68, 1)",
        borderWidth: 1,
        borderRadius: 4,
      },
    ],
  };

  const incomeBreakdown = (data?.category_breakdown || []).filter(
    (item) => item.type === "INCOME"
  );
  const expenseBreakdown = (data?.category_breakdown || []).filter(
    (item) => item.type === "EXPENSE"
  );

  const incomeChartData = {
    labels: incomeBreakdown.map((item) => item.category_name),
    datasets: [
      {
        data: incomeBreakdown.map((item) => item.total),
        backgroundColor: [
          "rgba(34, 197, 94, 0.8)",
          "rgba(59, 130, 246, 0.8)",
          "rgba(168, 85, 247, 0.8)",
          "rgba(251, 191, 36, 0.8)",
          "rgba(20, 184, 166, 0.8)",
        ],
        borderWidth: 2,
        borderColor: "#fff",
      },
    ],
  };

  const expenseChartData = {
    labels: expenseBreakdown.map((item) => item.category_name),
    datasets: [
      {
        data: expenseBreakdown.map((item) => item.total),
        backgroundColor: [
          "rgba(239, 68, 68, 0.8)",
          "rgba(249, 115, 22, 0.8)",
          "rgba(234, 179, 8, 0.8)",
          "rgba(236, 72, 153, 0.8)",
          "rgba(107, 114, 128, 0.8)",
          "rgba(99, 102, 241, 0.8)",
          "rgba(139, 92, 246, 0.8)",
        ],
        borderWidth: 2,
        borderColor: "#fff",
      },
    ],
  };

  if (isLoading) {
    return (
      <div className="flex justify-center items-center py-20">
        <Loader2 className="animate-spin h-8 w-8 text-blue-600" />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h2 className="text-2xl sm:text-3xl font-bold tracking-tight">
            Finance Dashboard
          </h2>
          <p className="text-sm sm:text-base text-slate-500">
            Ringkasan dan analisis keuangan perusahaan.
          </p>
        </div>
        <div className="flex gap-2 items-center flex-wrap">
          <Input
            type="date"
            value={startDate}
            onChange={(e: any) => setStartDate(e.target.value)}
            className="w-auto"
            placeholder="Dari tanggal"
          />
          <span className="text-slate-400">-</span>
          <Input
            type="date"
            value={endDate}
            onChange={(e: any) => setEndDate(e.target.value)}
            className="w-auto"
            placeholder="Sampai tanggal"
          />
          <Button onClick={handleFilter} size="sm">
            Filter
          </Button>
          <Button onClick={handleReset} variant="outline" size="sm">
            Reset
          </Button>
        </div>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Pemasukan</CardTitle>
            <TrendingUp className="h-4 w-4 text-green-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">
              {formatCurrency(data?.total_income || 0)}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Pengeluaran</CardTitle>
            <TrendingDown className="h-4 w-4 text-red-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-red-600">
              {formatCurrency(data?.total_expense || 0)}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Saldo Bersih</CardTitle>
            <Wallet className="h-4 w-4 text-blue-600" />
          </CardHeader>
          <CardContent>
            <div className={`text-2xl font-bold ${(data?.net_balance || 0) >= 0 ? "text-blue-600" : "text-red-600"}`}>
              {formatCurrency(data?.net_balance || 0)}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Transaksi</CardTitle>
            <Activity className="h-4 w-4 text-slate-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {data?.transaction_count || 0}
            </div>
          </CardContent>
        </Card>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Pemasukan vs Pengeluaran per Bulan</CardTitle>
          </CardHeader>
          <CardContent>
            {(data?.monthly_summary || []).length > 0 ? (
              <div className="h-[300px]">
                <Bar
                  data={monthlyChartData}
                  options={{
                    responsive: true,
                    maintainAspectRatio: false,
                    plugins: {
                      legend: { position: "top" },
                      tooltip: {
                        callbacks: {
                          label: (context) =>
                            `${context.dataset.label}: ${formatCurrency(context.parsed.y)}`,
                        },
                      },
                    },
                    scales: {
                      y: {
                        beginAtZero: true,
                        ticks: {
                          callback: (value) =>
                            new Intl.NumberFormat("id-ID", {
                              notation: "compact",
                              maximumFractionDigits: 1,
                            }).format(value as number),
                        },
                      },
                    },
                  }}
                />
              </div>
            ) : (
              <div className="flex items-center justify-center h-[300px] text-slate-400">
                Belum ada data bulanan
              </div>
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Breakdown Pengeluaran per Kategori</CardTitle>
          </CardHeader>
          <CardContent>
            {expenseBreakdown.length > 0 ? (
              <div className="h-[300px] flex items-center justify-center">
                <Doughnut
                  data={expenseChartData}
                  options={{
                    responsive: true,
                    maintainAspectRatio: false,
                    plugins: {
                      legend: { position: "right" },
                      tooltip: {
                        callbacks: {
                          label: (context) =>
                            `${context.label}: ${formatCurrency(context.parsed)}`,
                        },
                      },
                    },
                  }}
                />
              </div>
            ) : (
              <div className="flex items-center justify-center h-[300px] text-slate-400">
                Belum ada data pengeluaran
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Breakdown Pemasukan per Kategori</CardTitle>
          </CardHeader>
          <CardContent>
            {incomeBreakdown.length > 0 ? (
              <div className="h-[300px] flex items-center justify-center">
                <Doughnut
                  data={incomeChartData}
                  options={{
                    responsive: true,
                    maintainAspectRatio: false,
                    plugins: {
                      legend: { position: "right" },
                      tooltip: {
                        callbacks: {
                          label: (context) =>
                            `${context.label}: ${formatCurrency(context.parsed)}`,
                        },
                      },
                    },
                  }}
                />
              </div>
            ) : (
              <div className="flex items-center justify-center h-[300px] text-slate-400">
                Belum ada data pemasukan
              </div>
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Transaksi Terbaru</CardTitle>
          </CardHeader>
          <CardContent>
            {(data?.recent_transactions || []).length > 0 ? (
              <div className="space-y-3">
                {(data?.recent_transactions || []).map((tx) => (
                  <div
                    key={tx.id}
                    className="flex items-center justify-between py-2 border-b last:border-0"
                  >
                    <div>
                      <p className="font-medium text-sm">{tx.category_name}</p>
                      <p className="text-xs text-slate-400">
                        {new Date(tx.transaction_date).toLocaleDateString("id-ID", {
                          day: "numeric",
                          month: "short",
                          year: "numeric",
                        })}
                        {" - "}
                        {tx.creator_name}
                      </p>
                    </div>
                    <span
                      className={`font-bold text-sm ${
                        tx.type === "INCOME" ? "text-green-600" : "text-red-600"
                      }`}
                    >
                      {tx.type === "INCOME" ? "+" : "-"}
                      {formatCurrency(tx.amount)}
                    </span>
                  </div>
                ))}
              </div>
            ) : (
              <div className="flex items-center justify-center h-[200px] text-slate-400">
                Belum ada transaksi
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
};
