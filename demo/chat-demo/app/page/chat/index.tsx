import { ConversationList } from "./ConversationList";
import { ChatWindow } from "./ChatWindow";
import BuddyWorkerInstance from "util/im-worker/BuddyWorkerInstance";
import { useNavigate } from "react-router";
import { useEffect, useState, useCallback } from "react";
import { useMessages } from "~/hook/useMessages";
import { IM_CONSTANT } from "util/constant";
import { decodeBase64Response } from "util/index";

interface Conversation {
  id: number;
  title: string;
  user_ids: number[];
}

export function Chat() {
  const navigate = useNavigate();
  const [status, setStatus] = useState<
    "connecting" | "connected" | "disconnected"
  >("disconnected");

  const [accountId, setAccountId] = useState<number | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [currentConversation, setCurrentConversation] =
    useState<Conversation | null>(null);

  const [conversations, setConversations] = useState<Conversation[]>([]);

  const { messageStore, sendMessage, getConversationMessages } =
    useMessages(accountId);

  const onReconnect = useCallback(() => {
    if (accountId && token) {
      BuddyWorkerInstance.connect({
        accountId: accountId,
        token: token,
        url: "ws://localhost:8760",
      });
    }
  }, [accountId, token]);

  const updateConversation = (id: number) => {
    setConversations((prev) => {
      const newConversations = [...prev];
      if (newConversations.length > 0) {
        newConversations[newConversations.length - 1] = {
          ...newConversations[newConversations.length - 1],
          id,
        };
      }
      return newConversations;
    });
  };

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
        if (data.type === "message") {
          const payload = data.payload;
          if (
            payload.kind === IM_CONSTANT.ConvsCreate &&
            payload.mesg === "success"
          ) {
            try {
              const rsp = decodeBase64Response(payload.rsp);
              console.log("rsp", rsp);
              updateConversation(rsp.id);
            } catch (error) {
              console.error("Failed to parse response:", error);
            }
          }
        }
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
        <ConversationList
          status={status}
          onReconnect={onReconnect}
          onConversationSelect={setCurrentConversation}
          messageStore={messageStore}
          conversations={conversations}
          setConversations={setConversations}
        />
      </div>

      <div className="flex-1">
        {currentConversation && accountId && (
          <ChatWindow
            currentConversation={currentConversation}
            messages={getConversationMessages(currentConversation.title)}
            onSendMessage={(content) =>
              sendMessage(currentConversation.id.toString(), content)
            }
            accountId={accountId}
          />
        )}
      </div>
    </div>
  );
}
