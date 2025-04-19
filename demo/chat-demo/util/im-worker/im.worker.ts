import { IM_CONSTANT } from "../constant";

interface WSConfig {
  url: string;
  accountId: number;
  token: string;
}

const decodeBase64Response = (base64: string) => {
  if (!base64) {
    console.error("Empty base64 string");
    return null;
  }

  try {
    const decoded = atob(base64);

    try {
      return JSON.parse(decoded);
    } catch (jsonError) {
      return decoded;
    }
  } catch (error) {
    console.error("Failed to decode base64:", error);
    return null;
  }
};

class WebSocketManager {
  private static instances: Map<number, WebSocketClient> = new Map();

  // 通过 accountId 获取实例
  static getInstance(config: WSConfig): WebSocketClient {
    if (!this.instances.has(config.accountId)) {
      this.instances.set(config.accountId, new WebSocketClient(config));
    }
    return this.instances.get(config.accountId)!;
  }

  static hasInstance(accountId: number): boolean {
    return this.instances.has(accountId);
  }

  static removeInstance(accountId: number) {
    const instance = this.instances.get(accountId);
    if (instance) {
      instance.disconnect();
      this.instances.delete(accountId);
    }
  }

  static getActiveConnections(): number[] {
    return Array.from(this.instances.keys());
  }
}

class WebSocketClient {
  private ws: WebSocket | null = null;
  private config: WSConfig;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 3;
  private reconnectTimeout = 3000;
  private isConnecting = false;

  private heartbeatTimer: ReturnType<typeof setInterval> | null = null;
  private missedHeartbeats = 0;
  private readonly maxMissedHeartbeats = 3;
  private readonly heartbeatIntervalTime = 30000;
  private readonly heartbeatTimeout = 3000;

  constructor(config: WSConfig) {
    this.config = config;
  }

  private startHeartbeat() {
    this.stopHeartbeat();

    this.heartbeatTimer = setInterval(() => {
      if (this.ws?.readyState === WebSocket.OPEN) {
        this.missedHeartbeats++;

        // 发送心跳消息
        this.send({ type: "ping", timestamp: Date.now() });

        if (this.missedHeartbeats >= this.maxMissedHeartbeats) {
          self.postMessage({
            type: "heartbeat_failed",
            to: this.config.accountId,
          });
          this.reconnect();
          return;
        }

        setTimeout(() => {
          if (this.missedHeartbeats >= this.maxMissedHeartbeats) {
            self.postMessage({
              type: "heartbeat_timeout",
              to: this.config.accountId,
            });
            this.reconnect();
          }
        }, this.heartbeatTimeout);
      }
    }, this.heartbeatIntervalTime);
  }

  private stopHeartbeat() {
    if (this.heartbeatTimer) {
      clearInterval(this.heartbeatTimer);
      this.heartbeatTimer = null;
    }
    this.missedHeartbeats = 0;
  }

  private handlePong() {
    this.missedHeartbeats = 0;
    self.postMessage({ type: "heartbeat_success" });
  }

  connect() {
    if (this.isConnecting || this.ws?.readyState === WebSocket.OPEN) {
      self.postMessage({ type: "connected", to: this.config.accountId });
      return;
    }

    this.isConnecting = true;
    this.reconnectAttempts = 0;

    this._connect();
  }

  createConversation(payload: any) {
    const msgStr = JSON.stringify(payload);
    const encoder = new TextEncoder();
    const msgBytes = encoder.encode(msgStr);
    this.send({
      kind: 8000,
      reqs: Array.from(msgBytes),
    });
  }

  private reconnect() {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      this.isConnecting = false;
      self.postMessage({ type: "disconnected", to: this.config.accountId });
      return;
    }

    this.reconnectAttempts++;
    self.postMessage({
      type: "connecting",
      payload: { attempt: this.reconnectAttempts },
      to: this.config.accountId,
    });

    setTimeout(() => {
      this._connect();
    }, this.reconnectTimeout * this.reconnectAttempts);
  }

  private _connect() {
    this.ws = new WebSocket(
      this.config.url +
        "?id=" +
        this.config.accountId +
        "&token=" +
        this.config.token
    );

    this.ws.onopen = () => {
      this.isConnecting = false;
      this.reconnectAttempts = 0;
      self.postMessage({ type: "connected", to: this.config.accountId });
      // this.startHeartbeat();
    };

    this.ws.onclose = () => {
      this.reconnect();
    };

    this.ws.onerror = (error) => {
      console.error("WebSocket error", error);
      self.postMessage({
        type: "error",
        payload: error.toString(),
        to: this.config.accountId,
      });
    };

    this.ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        if (data?.kind === IM_CONSTANT.Heartbeat) {
          this.handlePong();
          return;
        }
        self.postMessage({
          type: "message",
          to: this.config.accountId,
          ...data,
          rsp: data.rsp ? decodeBase64Response(data.rsp) : null,
        });
      } catch (error) {
        self.postMessage({
          type: "error",
          payload: "Invalid message format",
          to: this.config.accountId,
        });
      }
    };
  }

  send(message: any) {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(message));
    } else {
      self.postMessage({
        type: "error",
        payload: "WebSocket is not connected",
      });
    }
  }

  disconnect() {
    this.stopHeartbeat();
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }
}

self.onmessage = (
  event: MessageEvent<{
    type: string;
    payload: WSConfig | any;
    accountId?: number;
    token?: string;
    url?: string;
  }>
) => {
  const { type, payload, accountId, token, url } = event.data;
  console.log("receive order", event.data);

  if (!accountId || !type) {
    self.postMessage({
      type: "error",
      payload: "Missing required accountId or type",
    });
    return;
  }

  switch (type) {
    case "connect":
      if (!url || !token) {
        self.postMessage({
          type: "error",
          payload: "Missing required connect parameters",
        });
        return;
      }
      WebSocketManager.getInstance({
        url,
        accountId,
        token,
      }).connect();
      break;

    case "createConversation":
      console.log("createConversation", payload);
      WebSocketManager.getInstance({
        accountId,
      } as WSConfig).createConversation(payload);
      break;

    case "send":
      const content = payload.reqs.payload;
      const encoder = new TextEncoder();
      const contentBytes = encoder.encode(content);
      const sendMsg = {
        ...payload.reqs,
        from_id: String(payload.reqs.from_id),
        convs_id: Number(payload.reqs.convs_id),
        payload: Array.from(contentBytes),
      };
      console.log("sendMsg", sendMsg);
      const reqsStr = JSON.stringify(sendMsg);

      WebSocketManager.getInstance({ accountId } as WSConfig).send({
        ...payload,
        reqs: Array.from(encoder.encode(reqsStr)),
      });
      break;

    case "disconnect":
      WebSocketManager.removeInstance(accountId);
      break;

    case "disconnect_all":
      WebSocketManager.getActiveConnections().forEach((id) => {
        WebSocketManager.removeInstance(id);
      });
      break;
  }
};
