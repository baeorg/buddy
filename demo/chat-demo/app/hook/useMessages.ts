import { useState, useEffect, useCallback } from "react";
import { IM_CONSTANT } from "util/constant";
import BuddyWorkerInstance from "util/im-worker/BuddyWorkerInstance";

export interface Message {
  id: string;
  conversationId: string;
  content: string;
  senderId: number;
  timestamp: number;
  // ... 其他消息属性
}

export interface MessageStore {
  [conversationId: string]: Message[];
}

export function useMessages(accountId: number | null) {
  const [messageStore, setMessageStore] = useState<MessageStore>({});

  const sendMessage = useCallback(
    (conversationId: string, content: string) => {
      if (!accountId) return;

      BuddyWorkerInstance.sendMessage({
        payload: {
          kind: IM_CONSTANT.MesgCreate,
          resq: {
            from_id: accountId,
            convs_id: conversationId,
            payload: content,
          },
        },
        accountId: accountId,
        token: localStorage.getItem("Buddy_Token") || "",
      });
    },
    [accountId]
  );

  const getConversationMessages = useCallback(
    (conversationId: string) => {
      return messageStore[conversationId] || [];
    },
    [messageStore]
  );

  const addMessage = useCallback((message: Message) => {
    setMessageStore((prev) => ({
      ...prev,
      [message.conversationId]: [
        ...(prev[message.conversationId] || []),
        message,
      ],
    }));
  }, []);

  useEffect(() => {
    if (!accountId) return;

    const handleMessage = (data: any) => {
      if (data.type === "message" && data.to === accountId) {
        const payload = data.payload;
        switch (payload.kind) {
          case IM_CONSTANT.MesgCreate:
            console.log("MesgCreate", payload.rsp);
            addMessage(payload.rsp);
            break;
          case IM_CONSTANT.MesgGet:
            const messages = payload.rsp;
            console.log("MesgGet", messages);
            setMessageStore((prev) => ({
              ...prev,
              [payload.convs_id]: messages,
            }));
            break;
        }
      }
    };

    BuddyWorkerInstance.addListener(handleMessage);
    return () => BuddyWorkerInstance.removeListener(handleMessage);
  }, [accountId, addMessage]);

  return {
    messageStore,
    sendMessage,
    getConversationMessages,
  };
}
