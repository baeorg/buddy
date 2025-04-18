import axios from "../../axios";

interface MessageParams {
  content: string;
  chatId: string;
}

export default {
  sendMessage: (data: MessageParams) => {
    return axios.post("/messages/send", data);
  },

  getRecentChats: () => {
    return axios.get("/chats/recent");
  },

  getMessages: (chatId: string, params: { page: number; limit: number }) => {
    return axios.post(`/messages/history/${chatId}`, params);
  },
};
