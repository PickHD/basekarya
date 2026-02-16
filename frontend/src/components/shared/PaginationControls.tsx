import { Button } from "@/components/ui/button";
import type { Meta } from "@/types/api";
import { ArrowDownCircle, ChevronLeft, ChevronRight, Loader2 } from "lucide-react";

interface PaginationControlsProps {
  meta?: Meta;
  onPageChange?: (page: number) => void;
  onLoadMore?: () => void;
  isLoading?: boolean;
}

export function PaginationControls({
  meta,
  onPageChange,
  onLoadMore,
  isLoading = false,
}: PaginationControlsProps) {
  if (!meta) return null;

  const isCursorMode =
    meta.total_page === undefined && meta.has_next !== undefined;

  if (isCursorMode) {
    if (!meta.has_next) return null;

    return (
      <div className="mt-6 flex justify-center">
        <Button
          variant="secondary"
          className="w-full md:w-auto min-w-[200px]"
          onClick={onLoadMore}
          disabled={isLoading}
        >
          {isLoading ? (
            <>
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              Loading more...
            </>
          ) : (
            <>
              Load More
              <ArrowDownCircle className="ml-2 h-4 w-4" />
            </>
          )}
        </Button>
      </div>
    );
  }

  const currentPage = meta.page || 1;
  const totalPages = meta.total_page || 1;
  const totalData = meta.total_data;

  if (totalPages <= 1 && !totalData) return null;

  return (
    <div className="flex flex-col sm:flex-row items-center justify-between mt-4 px-2 gap-4 sm:gap-0">
      <div className="text-sm text-slate-500 text-center sm:text-left">
        {totalData !== undefined && (
          <span>
            Total <strong>{totalData}</strong> records.{" "}
          </span>
        )}
        Page <span className="font-medium text-slate-900">{currentPage}</span>{" "}
        of <span className="font-medium text-slate-900">{totalPages}</span>
      </div>

      <div className="flex gap-2">
        <Button
          variant="outline"
          size="sm"
          onClick={() => onPageChange?.(currentPage - 1)}
          disabled={currentPage <= 1 || isLoading}
        >
          <ChevronLeft className="h-4 w-4 mr-1" />
          Prev
        </Button>

        <Button
          variant="outline"
          size="sm"
          onClick={() => onPageChange?.(currentPage + 1)}
          disabled={currentPage >= totalPages || isLoading}
        >
          {isLoading ? (
            <Loader2 className="h-4 w-4 animate-spin" />
          ) : (
            <>
              Next
              <ChevronRight className="h-4 w-4 ml-1" />
            </>
          )}
        </Button>
      </div>
    </div>
  );
}
