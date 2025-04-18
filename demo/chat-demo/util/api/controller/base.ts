import axios from "../../axios";

interface LoginParams {
  username: string;
  password: string;
}

interface RegisterParams {
  id: number;
  token: string;
}

export default {
  login: (data: LoginParams, options?: any) => {
    return axios.post("/auth/login", data, options);
  },
  register: (data: RegisterParams, options?: any): Promise<string> => {
    return axios.post("/user/im", data, options);
  },
};
