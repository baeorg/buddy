import { useState } from "react";
import type { Message } from "~/hook/useMessages";

type Conversation = {
  title: string;
  user_ids: number[];
};

export function ChatWindow({
  currentConversation,
  messages,
  onSendMessage,
  accountId,
}: {
  currentConversation: Conversation;
  messages: Message[];
  onSendMessage: (message: string) => void;
  accountId: number;
}) {
  const [inputValue, setInputValue] = useState("");

  const handleSend = () => {
    if (!inputValue.trim()) return;
    onSendMessage(inputValue);
    setInputValue("");
  };

  return (
    <div className="h-full flex flex-col">
      {/* ============== HEADER ============== */}
      <div className="p-4 border-b border-gray-200">
        <div className="flex items-center gap-3">
          <img
            src={"/avatar1.png"}
            alt="当前聊天"
            className="w-10 h-10 rounded-full"
          />
          <div>
            <h2 className="font-medium">{currentConversation.title}</h2>
            {/* <p className="text-sm text-gray-500">在线</p> */}
          </div>
        </div>
      </div>

      {/* ============== MESSAGE LIST ============== */}
      <div className="flex-1 overflow-y-auto p-4 space-y-4">
        {messages.map((message) => (
          <MessageItem
            key={message.id}
            message={message}
            isMine={message.senderId === accountId}
          />
        ))}
      </div>

      {/* ============== MESSAGE INPUT ============== */}
      <div className="p-4 border-t border-gray-200">
        <div className="flex items-center gap-2">
          <input
            type="text"
            value={inputValue}
            onChange={(e) => setInputValue(e.target.value)}
            onKeyPress={(e) => e.key === "Enter" && handleSend()}
            className="flex-1 px-4 py-2 border rounded"
          />
          <button
            onClick={handleSend}
            className="px-4 py-2 bg-blue-500 text-white rounded"
          >
            发送
          </button>
        </div>
      </div>
    </div>
  );
}

function MessageItem({
  message,
  isMine,
}: {
  message: Message;
  isMine: boolean;
}) {
  return (
    <div className={`flex ${isMine ? "justify-end" : "justify-start"}`}>
      <div
        className={`max-w-[70%] rounded-lg p-3 ${
          isMine ? "bg-blue-500 text-white" : "bg-gray-100 text-gray-800"
        }`}
      >
        {message.content}
      </div>
    </div>
  );
}
