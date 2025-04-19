class BuddyWorkerInstance {
  private static instance: Worker | null = null;

  private static listeners = new Set<(event: MessageEvent) => void>();

  static getInstance(): Worker {
    if (!this.instance) {
      this.instance = new Worker(new URL("./im.worker.ts", import.meta.url), {
        type: "module",
      });

      this.instance.onmessage = (event) => {
        const data = event.data;
        // console.log("receive message from worker", data);
        this.listeners.forEach((listener) => listener(data));
      };
    }
    return this.instance;
  }

  static connect({
    accountId,
    token,
    url,
  }: {
    accountId: number;
    token: string;
    url: string;
  }) {
    if (!accountId || !token || !url) {
      throw new Error("Missing required parameters");
    }
    this.getInstance().postMessage({
      type: "connect",
      accountId,
      token,
      url,
    });
  }

  static addListener(listener: (event: MessageEvent) => void) {
    this.listeners.add(listener);
  }

  static removeListener(listener: (event: MessageEvent) => void) {
    this.listeners.delete(listener);
  }

  static removeAllListeners() {
    this.listeners.clear();
  }

  static cleanup() {
    if (this.instance) {
      this.instance.terminate();
      this.instance = null;
    }
    this.listeners.clear();
  }

  static sendMessage({
    payload,
    accountId,
    token,
  }: {
    payload: any;
    accountId: number;
    token: string;
  }) {
    if (this.instance) {
      this.instance.postMessage({
        type: "send",
        payload,
        accountId: accountId,
        token: token,
      });
    }
  }

  static createConversation({
    accountId,
    token,
    payload,
  }: {
    accountId: number;
    token: string;
    payload: {
      title: string;
      user_ids: number[];
    };
  }) {
    if (this.instance) {
      this.instance.postMessage({
        type: "createConversation",
        accountId,
        token,
        payload,
      });
    }
  }
}

export default BuddyWorkerInstance;
