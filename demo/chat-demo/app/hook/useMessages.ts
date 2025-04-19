import { useState, useEffect, useCallback } from "react";
import { IM_CONSTANT } from "util/constant";
import { v4 as uuidv4 } from "uuid";
import BuddyWorkerInstance from "util/im-worker/BuddyWorkerInstance";

export interface Message {
  clientId: string;
  id?: string;
  conversationId: string;
  content: string;
  senderId: number;
  timestamp: number;
}

export interface MessageStore {
  [conversationId: string]: Message[];
}

export function useMessages(accountId: number | null, token: string | null) {
  const [messageStore, setMessageStore] = useState<MessageStore>({});

  const sendMessage = useCallback(
    (conversationId: string, content: string) => {
      if (!accountId || !token) return;

      BuddyWorkerInstance.sendMessage({
        payload: {
          kind: IM_CONSTANT.MesgCreate,
          reqs: {
            from_id: accountId,
            convs_id: conversationId,
            payload: content,
          },
        },
        accountId: accountId,
        token: token,
      });

      setMessageStore((prev) => ({
        ...prev,
        [conversationId]: [
          ...(prev[conversationId] || []),
          {
            clientId: uuidv4(),
            content,
            senderId: accountId,
            timestamp: Date.now(),
            conversationId,
          },
        ],
      }));
    },
    [accountId, token]
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
        switch (data.kind) {
          case IM_CONSTANT.MesgCreate:
            console.log("MesgCreate", data.rsp);
            // addMessage(data.rsp);
            break;
          case IM_CONSTANT.MesgGet:
            console.log("MesgGet", data.rsp);
            // addMessage(data.rsp);
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
