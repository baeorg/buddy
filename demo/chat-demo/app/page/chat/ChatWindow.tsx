export function ChatWindow() {
  return (
    <div className="h-full flex flex-col">
      {/* ============== HEADER ============== */}
      <div className="p-4 border-b border-gray-200">
        <div className="flex items-center gap-3">
          <img
            src="/avatar1.png"
            alt="当前聊天"
            className="w-10 h-10 rounded-full"
          />
          <div>
            <h2 className="font-medium">BAE-DevTeam</h2>
            {/* <p className="text-sm text-gray-500">在线</p> */}
          </div>
        </div>
      </div>

      {/* ============== MESSAGE LIST ============== */}
      <div className="flex-1 overflow-y-auto p-4 space-y-4">
        {MOCK_MESSAGES.map((message) => (
          <div
            key={message.id}
            className={`flex ${
              message.isSelf ? "justify-end" : "justify-start"
            }`}
          >
            <div
              className={`max-w-[70%] rounded-lg p-3 ${
                message.isSelf
                  ? "bg-blue-500 text-white"
                  : "bg-gray-100 text-gray-800"
              }`}
            >
              {message.content}
            </div>
          </div>
        ))}
      </div>

      {/* ============== MESSAGE INPUT ============== */}
      <div className="p-4 border-t border-gray-200">
        <div className="flex items-center gap-2">
          <input
            type="text"
            placeholder="发送消息..."
            className="flex-1 px-4 py-2 rounded-full border border-gray-200 focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <button className="px-4 py-2 bg-blue-500 text-white rounded-full hover:bg-blue-600">
            发送
          </button>
        </div>
      </div>
    </div>
  );
}

const MOCK_MESSAGES = [
  {
    id: 1,
    content: "👋 你好！",
    isSelf: false,
  },
  {
    id: 2,
    content: "你好！有什么我可以帮你的吗？",
    isSelf: true,
  },
];
