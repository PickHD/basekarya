export interface Meta {
  // common
  limit: number;

  // offset specific
  page?: number;
  total_page?: number;
  total_data?: number;

  // cursor specific
  has_next?: boolean;
  next_cursor?: string;
}

export interface ApiResponse<T> {
  message: string;
  data: T;
  meta?: Meta;
  error?: any;
}
