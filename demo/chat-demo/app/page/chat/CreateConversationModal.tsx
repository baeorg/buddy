import { useState } from "react";

type ChatType = "single" | "group";

export function CreateConversationModal({
  isOpen,
  onClose,
  onCreate,
}: {
  isOpen: boolean;
  onClose: () => void;
  onCreate: (type: ChatType) => void;
}) {
  const [selectedType, setSelectedType] = useState<ChatType>("single");

  return isOpen ? (
    <div
      className="fixed inset-0 flex items-center justify-center z-50"
      style={{
        backgroundColor: "rgba(0, 0, 0, 0.7)",
      }}
    >
      <div className="bg-white rounded-lg p-6 w-80">
        <h3 className="text-lg font-semibold mb-4">创建新对话</h3>

        <div className="space-y-3 mb-6">
          <label className="flex items-center gap-2 cursor-pointer">
            <input
              type="radio"
              name="chatType"
              checked={selectedType === "single"}
              onChange={() => setSelectedType("single")}
              className="form-radio text-blue-500"
            />
            <span>单聊</span>
          </label>

          <label className="flex items-center gap-2 cursor-pointer">
            <input
              type="radio"
              name="chatType"
              checked={selectedType === "group"}
              onChange={() => setSelectedType("group")}
              className="form-radio text-blue-500"
            />
            <span>群聊</span>
          </label>
        </div>

        <div className="flex justify-end gap-2">
          <button
            onClick={onClose}
            className="px-4 py-2 text-gray-600 hover:text-gray-800"
          >
            取消
          </button>
          <button
            onClick={() => {
              onCreate(selectedType);
              onClose();
            }}
            className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
          >
            创建
          </button>
        </div>
      </div>
    </div>
  ) : null;
}
