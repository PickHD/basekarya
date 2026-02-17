export interface NotificationPayload {
  id: number;
  user_id: number;
  type:
    | "APPROVED"
    | "REJECTED"
    | "LEAVE_APPROVAL_REQ"
    | "REIMBURSE_APPROVAL_REQ";
  title: string;
  message: string;
  related_id: number;
  is_read: boolean;
  created_at: string;
}

export interface WebSocketMessage {
  type: string;
  payload: any;
}

export interface NotificationResponse {
  data: NotificationPayload[];
}
