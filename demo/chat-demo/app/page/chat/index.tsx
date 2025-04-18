import { ConversationList } from "./ConversationList";
import { ChatWindow } from "./ChatWindow";
import BuddyWorkerInstance from "util/im-worker/BuddyWorkerInstance";
import { useNavigate } from "react-router";
import { useEffect, useState, useCallback } from "react";

export function Chat() {
  const navigate = useNavigate();
  const [status, setStatus] = useState<
    "connecting" | "connected" | "disconnected"
  >("disconnected");

  const [accountId, setAccountId] = useState<number | null>(null);
  const [token, setToken] = useState<string | null>(null);

  const onReconnect = useCallback(() => {
    if (accountId && token) {
      BuddyWorkerInstance.connect({
        accountId: accountId,
        token: token,
        url: "ws://localhost:8760",
      });
    }
  }, [accountId, token]);

  useEffect(() => {
    const accountId = Number(localStorage.getItem("Buddy_AccountId"));
    const token = localStorage.getItem("Buddy_Token");
    if (!accountId || !token) {
      navigate("/login");
      return;
    }
    setAccountId(accountId);
    setToken(token);

    BuddyWorkerInstance.connect({
      accountId: Number(accountId),
      token,
      url: "ws://localhost:8760/ws",
    });
    setStatus("connecting");

    const msgListener = (data: any) => {
      console.log(`message to account(${accountId}):`, data);
      if (data.to === accountId) {
        if (data.type === "disconnected") setStatus("disconnected");
        if (data.type === "connected") setStatus("connected");
        if (data.type === "connecting") setStatus("connecting");
      }
    };

    BuddyWorkerInstance.addListener(msgListener);

    return () => {
      BuddyWorkerInstance.removeListener(msgListener);
    };
  }, []);

  useEffect(() => {
    console.log("status", status);
  }, [status]);

  return (
    <div className="h-screen flex">
      <div className="w-80 border-r border-gray-200">
        <ConversationList status={status} onReconnect={onReconnect} />
      </div>

      <div className="flex-1">
        <ChatWindow />
      </div>
    </div>
  );
}
