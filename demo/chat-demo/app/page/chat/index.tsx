import { ConversationList } from "./ConversationList";
import { ChatWindow } from "./ChatWindow";
import BuddyWorkerInstance from "util/im-worker/BuddyWorkerInstance";
import { useNavigate, useSearchParams } from "react-router";
import { useEffect, useState, useCallback } from "react";
import { useMessages } from "~/hook/useMessages";
import { IM_CONSTANT } from "util/constant";

interface Conversation {
  id: number;
  title: string;
  user_ids: number[];
}

export function Chat() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();

  const [status, setStatus] = useState<
    "connecting" | "connected" | "disconnected"
  >("disconnected");

  const [accountId, setAccountId] = useState<number | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [username, setUsername] = useState<string | null>(null);

  const [currentConversation, setCurrentConversation] =
    useState<Conversation | null>(null);
  const [conversations, setConversations] = useState<Conversation[]>([]);

  const { messageStore, sendMessage, getConversationMessages } = useMessages(
    accountId ? Number(accountId) : null,
    token || null
  );

  const onReconnect = useCallback(() => {
    if (accountId && token) {
      BuddyWorkerInstance.connect({
        accountId: Number(accountId),
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
    const accountId = Number(searchParams.get("id"));
    const token = searchParams.get("token");
    const username = searchParams.get("username");

    if (!accountId || !token) {
      navigate("/login");
      return;
    }
    setAccountId(accountId);
    setToken(token);
    setUsername(username);

    BuddyWorkerInstance.connect({
      accountId: Number(accountId),
      token: token,
      url: "ws://localhost:8760/ws",
    });

    const msgListener = (data: any) => {
      if (data.to === accountId) {
        if (data.type === "disconnected") setStatus("disconnected");
        if (data.type === "connected") setStatus("connected");
        if (data.type === "connecting") setStatus("connecting");
        if (data.type === "message") {
          if (
            data.kind === IM_CONSTANT.ConvsCreate &&
            data.mesg === "success"
          ) {
            try {
              updateConversation(data.rsp.id);
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
  }, [searchParams, navigate]);

  useEffect(() => {
    console.log("messageStore", messageStore);
  }, [messageStore]);

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
          accountId={Number(accountId)}
          token={token || ""}
        />
      </div>

      <div className="flex-1">
        {currentConversation && accountId && (
          <ChatWindow
            currentConversation={currentConversation}
            messages={getConversationMessages(
              currentConversation.id.toString()
            )}
            onSendMessage={(content) => {
              sendMessage(currentConversation.id.toString(), content);
            }}
            accountId={Number(accountId)}
          />
        )}
      </div>
    </div>
  );
}
