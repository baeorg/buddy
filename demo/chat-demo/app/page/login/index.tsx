import { useEffect, useState } from "react";
import { useNavigate } from "react-router";
import { api } from "../../../util";
export function meta() {
  return [
    { title: "Chat - Login" },
    { name: "description", content: "Login Page" },
  ];
}

export function Login() {
  const [accountId, setAccountId] = useState("");
  const [token, setToken] = useState("");
  const [baseUrl, setBaseUrl] = useState("http://localhost:8762");
  const navigate = useNavigate();

  const handleSubmit = () => {
    api.setBaseUrl(baseUrl);
    api.register({ id: Number(accountId), token: token }).then((res) => {
      if (res === "Created") {
        localStorage.setItem("Buddy_AccountId", accountId);
        localStorage.setItem("Buddy_Token", token);
        navigate("/chat");
      } else {
        alert("Failed to register");
      }
    });
  };

  useEffect(() => {
    const accountId = localStorage.getItem("Buddy_AccountId");
    const token = localStorage.getItem("Buddy_Token");
    setAccountId(accountId || "");
    setToken(token || "");
  }, []);

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="max-w-md w-full p-6 bg-white rounded-lg shadow-lg">
        <div className="text-center mb-8">
          <img
            src="/avatar.png"
            alt="Welcome"
            className="mx-auto w-16 h-16 rounded-full"
          />
          <h2 className="mt-4 text-2xl font-bold">Login Chat</h2>
        </div>

        <div className="space-y-4">
          <div>
            <input
              type="text"
              name="api"
              placeholder="API Address"
              className="w-full px-4 py-2 rounded-lg border outline-none focus:ring-black/50 overflow-hidden"
              value={baseUrl}
              onChange={(e) => setBaseUrl(e.target.value)}
            />
          </div>
          <div>
            <input
              type="number"
              name="account_id"
              placeholder="account id"
              className="w-full px-4 py-2 rounded-lg border outline-none focus:ring-black/50 overflow-hidden"
              value={accountId}
              onChange={(e) => setAccountId(e.target.value)}
            />
          </div>
          <div>
            <input
              type="password"
              name="token"
              placeholder="token"
              className="w-full px-4 py-2 rounded-lg border outline-none focus:ring-black/50"
              value={token}
              onChange={(e) => setToken(e.target.value)}
            />
          </div>
          <button
            type="button"
            className="w-full py-3 px-4 bg-black text-white rounded-lg hover:bg-black/80"
            onClick={handleSubmit}
          >
            LOGIN
          </button>
        </div>
      </div>
    </div>
  );
}
