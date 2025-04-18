import type { Route } from "./+types/chat";
import { Chat } from "../page/chat";
export function meta({}: Route.MetaArgs) {
  return [{ title: "Chat" }, { name: "description", content: "Chat" }];
}

export default function ChatRoute() {
  return <Chat />;
}
