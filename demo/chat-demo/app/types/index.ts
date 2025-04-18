export interface Conversation {
  id: number;
  title: string;
  user_ids: number[];
}
export type ChatType = "single" | "group";
export type ConnectionStatus = "connected" | "connecting" | "disconnected";
