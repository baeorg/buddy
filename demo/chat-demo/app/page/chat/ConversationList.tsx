import { useState } from "react";
import { CreateConversationModal } from "./CreateConversationModal";
import { useNavigate } from "react-router";
import BuddyWorkerInstance from "util/im-worker/BuddyWorkerInstance";

type ChatType = "single" | "group";
type ConnectionStatus = "connected" | "connecting" | "disconnected";

interface ConnectionStatusProps {
  status: ConnectionStatus;
  onReconnect?: () => void;
}
export function ConversationList({
  status,
  onReconnect,
}: {
  status: ConnectionStatus;
  onReconnect?: () => void;
}) {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const navigate = useNavigate();

  const handleCreateConversation = (type: ChatType) => {
    console.log("create conversation", type);
    const accountId = Number(localStorage.getItem("Buddy_AccountId"));
    const token = localStorage.getItem("Buddy_Token");
    if (!accountId || !token) {
      navigate("/login");
      return;
    }

    if (type === "single") {
      BuddyWorkerInstance.createConversation({
        accountId,
        token,
        type: "single",
      });
    }
  };

  return (
    <div className="h-full flex flex-col">
      {/* ============== HEADER ============== */}
      <div className="p-4 border-b border-gray-200 relative">
        <input
          type="text"
          placeholder="搜索"
          className="w-full px-3 py-2 rounded-lg border border-gray-200 focus:outline-none focus:ring-2 focus:ring-blue-500"
        />

        {/* ============== CONNECTION STATUS ============== */}
        <div className="flex justify-end items-center absolute right-8 top-0 bottom-0 ">
          <ConnectionStatus status={status} onReconnect={onReconnect} />
        </div>
      </div>

      {/* ============== CONVERSATION LIST ============== */}
      <div className="flex-1 overflow-y-auto">
        {MOCK_USERS.map((user) => (
          <div
            key={user.id}
            className="flex items-center gap-3 p-4 hover:bg-gray-100 cursor-pointer"
          >
            <img
              src={user.avatar}
              alt={user.name}
              className="w-12 h-12 rounded-full"
            />
            <div className="flex-1 min-w-0">
              <div className="flex items-center justify-between">
                <h3 className="font-medium truncate">{user.name}</h3>
                <span className="text-sm text-gray-500">{user.lastTime}</span>
              </div>
              <p className="text-sm text-gray-500 truncate">
                {user.lastMessage}
              </p>
            </div>
          </div>
        ))}
      </div>

      {/* ============== ADD CONVERSATION BUTTON ============== */}
      <div className="p-2 border-t border-gray-200">
        <button
          onClick={() => setIsModalOpen(true)}
          className="w-12 h-12 rounded-full mx-auto flex items-center justify-center gap-2 py-2 bg-blue-500 text-white hover:bg-blue-600 transition-colors"
        >
          <svg
            className="w-5 h-5"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M12 4v16m8-8H4"
            />
          </svg>
        </button>
      </div>
      <CreateConversationModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onCreate={handleCreateConversation}
      />
    </div>
  );
}

export function ConnectionStatus({
  status,
  onReconnect,
}: ConnectionStatusProps) {
  if (status === "connected") {
    return null;
  }

  return (
    <div className="flex items-center gap-2 text-sm">
      {status === "connecting" ? (
        <>
          <svg
            className="animate-spin h-4 w-4 text-gray-500"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
          >
            <circle
              className="opacity-25"
              cx="12"
              cy="12"
              r="10"
              stroke="currentColor"
              strokeWidth="4"
            />
            <path
              className="opacity-75"
              fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
            />
          </svg>
          <span className="text-gray-500">Connecting...</span>
        </>
      ) : (
        <>
          <button
            onClick={onReconnect}
            className="inline-flex items-center gap-1 text-gray-500 hover:text-gray-700"
          >
            <svg
              className="h-4 w-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
              />
            </svg>
            <span>Reconnect</span>
          </button>
        </>
      )}
    </div>
  );
}

const MOCK_USERS = [
  {
    id: 1,
    name: "BAE-DevTeam",
    avatar: "/avatar1.png",
    lastMessage: "👋 yst ppp",
    lastTime: "11:16",
  },
  {
    id: 2,
    name: "Reelz",
    avatar: "/avatar2.png",
    lastMessage: "🎮 Toosii",
    lastTime: "04:00",
  },
];
