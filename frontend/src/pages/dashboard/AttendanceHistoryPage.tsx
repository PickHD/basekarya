import { useState } from "react";
import { format, isValid } from "date-fns";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Loader2, Calendar, MapPin, Clock, FileText } from "lucide-react";
import { useAttendanceHistory } from "@/features/attendance/hooks/useAttendance";
import type { AttendanceLog } from "@/features/attendance/types";
import { PaginationControls } from "@/components/shared/PaginationControls";

export default function AttendanceHistoryPage() {
  const now = new Date();
  const [month, setMonth] = useState<string>(String(now.getMonth() + 1));
  const [year, setYear] = useState<string>(String(now.getFullYear()));

  const {
    data: response,
    isLoading,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
  } = useAttendanceHistory(parseInt(month), parseInt(year));

  const allLogs = response?.pages.flatMap((page) => page.data) || [];

  const handleFilterChange = (type: "month" | "year", val: string) => {
    if (type === "month") setMonth(val);
    else setYear(val);
  };

  const formatTime = (timeStr?: string) => {
    if (!timeStr) return "-";
    const date = new Date(timeStr);
    if (isNaN(date.getTime())) return "-";
    return date.toLocaleTimeString("id-ID", {
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  const formatDateSafe = (dateStr: string, pattern: string) => {
    if (!dateStr) return "-";
    const date = new Date(dateStr);
    if (!isValid(date)) return "-";
    return format(date, pattern);
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case "PRESENT":
        return "bg-green-100 text-green-700 hover:bg-green-200 border-green-200";
      case "LATE":
        return "bg-red-100 text-red-700 hover:bg-red-200 border-red-200";
      case "EXCUSED":
        return "bg-blue-100 text-blue-700 hover:bg-blue-200 border-blue-200";
      case "SICK":
        return "bg-amber-50 text-amber-600 hover:bg-amber-200 border-white-200";
      default:
        return "bg-slate-100 text-slate-700 border-slate-200";
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h2 className="text-3xl font-bold tracking-tight text-slate-900">
            My Attendance
          </h2>
          <p className="text-slate-500">
            View your monthly attendance records.
          </p>
        </div>

        <div className="flex gap-2">
          <Select
            value={month}
            onValueChange={(val) => handleFilterChange("month", val)}
          >
            <SelectTrigger className="w-[140px]">
              <SelectValue placeholder="Month" />
            </SelectTrigger>
            <SelectContent>
              {Array.from({ length: 12 }, (_, i) => (
                <SelectItem key={i + 1} value={String(i + 1)}>
                  {format(new Date(2024, i, 1), "MMMM")}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>

          <Select
            value={year}
            onValueChange={(val) => handleFilterChange("year", val)}
          >
            <SelectTrigger className="w-[100px]">
              <SelectValue placeholder="Year" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="2024">2024</SelectItem>
              <SelectItem value="2025">2025</SelectItem>
              <SelectItem value="2026">2026</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Calendar className="h-5 w-5 text-slate-500" />
            Records for{" "}
            {format(new Date(parseInt(year), parseInt(month) - 1), "MMMM yyyy")}
          </CardTitle>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="flex justify-center py-10">
              <Loader2 className="h-8 w-8 animate-spin text-slate-400" />
            </div>
          ) : allLogs.length > 0 ? (
            <>
              <div className="rounded-md border">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Date</TableHead>
                      <TableHead>Status</TableHead>
                      <TableHead>Check In</TableHead>
                      <TableHead>Check Out</TableHead>
                      <TableHead>Location</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {allLogs.map((log: AttendanceLog) => (
                      <TableRow key={log.id}>
                        {/* DATE */}
                        <TableCell className="font-medium">
                          {formatDateSafe(log.date, "dd MMM yyyy")}
                          <div className="text-xs text-slate-400 font-normal">
                            {formatDateSafe(log.date, "EEEE")}
                          </div>
                        </TableCell>

                        {/* STATUS */}
                        <TableCell>
                          <Badge
                            variant="outline"
                            className={getStatusColor(log.status)}
                          >
                            {log.status}
                          </Badge>
                          {log.is_suspicious && (
                            <Badge
                              variant="destructive"
                              className="ml-2 text-[10px] px-1 h-5"
                            >
                              FLAGGED
                            </Badge>
                          )}
                        </TableCell>

                        {/* CHECK IN TIME */}
                        <TableCell>
                          {log.check_in_time ? (
                            <div className="flex items-center gap-2">
                              <Clock className="w-4 h-4 text-slate-400" />
                              {formatTime(log.check_in_time)}
                            </div>
                          ) : (
                            <span className="text-slate-400">-</span>
                          )}
                        </TableCell>

                        {/* CHECK OUT TIME */}
                        <TableCell>
                          {log.check_out_time ? (
                            <div className="flex items-center gap-2">
                              <Clock className="w-4 h-4 text-slate-400" />
                              {formatTime(log.check_out_time)}
                            </div>
                          ) : (
                            <Badge
                              variant="outline"
                              className="text-slate-400 border-dashed font-normal"
                            >
                              On Duty
                            </Badge>
                          )}
                        </TableCell>

                        {/* LOCATION */}
                        <TableCell className="max-w-[250px]">
                          <div className="flex flex-col gap-1">
                            {/* Check In Location */}
                            <div className="flex items-start gap-1.5 text-xs text-slate-600">
                              <MapPin className="w-3.5 h-3.5 mt-0.5 text-green-600 shrink-0" />
                              <span
                                className="truncate"
                                title={
                                  log.check_in_address || "No location data"
                                }
                              >
                                {log.check_in_address || "No data"}
                              </span>
                            </div>

                            {/* Check Out Location */}
                            {log.check_out_address && (
                              <div className="flex items-start gap-1.5 text-xs text-slate-400">
                                <MapPin className="w-3.5 h-3.5 mt-0.5 text-orange-400 shrink-0" />
                                <span
                                  className="truncate"
                                  title={log.check_out_address}
                                >
                                  {log.check_out_address}
                                </span>
                              </div>
                            )}

                            {/* Notes Indicator */}
                            {log.notes && (
                              <div className="flex items-center gap-1 text-[10px] text-blue-600 mt-1">
                                <FileText className="w-3 h-3" />
                                <span
                                  className="truncate max-w-[200px]"
                                  title={log.notes}
                                >
                                  Note: {log.notes}
                                </span>
                              </div>
                            )}
                          </div>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </div>

              <PaginationControls
                meta={{
                  limit: 10,
                  has_next: hasNextPage,
                  next_cursor: "managed-by-query",
                }}
                onLoadMore={() => fetchNextPage()}
                isLoading={isFetchingNextPage}
              />
            </>
          ) : (
            <div className="text-center py-10 text-slate-500">
              No attendance records found for this period.
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
