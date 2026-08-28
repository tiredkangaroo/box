export interface Mailbox {
  id: number;
  server_hostport: string;
  username: string;
  primary_inbox: string;
}

export interface EmailSummary {
  id: number;
  subject: string;
  from_address: string;
  received_at: string;
}

export interface Email extends EmailSummary {
  body: string;
}

export interface MessageWithMailbox extends EmailSummary {
  mailbox: Mailbox;
}

export interface InboxLoad {
  mailboxes: Mailbox[];
  messages: MessageWithMailbox[];
  failed: number;
}

async function request<T>(url: string, init?: RequestInit): Promise<T> {
  const res = await fetch(url, {
    headers: { "Content-Type": "application/json" },
    ...init,
  });

  if (!res.ok) {
    let message = `Request failed (${res.status})`;
    try {
      const data = await res.json();
      if (data?.error) message = data.error;
    } catch {
      /* body not json */
    }
    throw new Error(message);
  }

  return (await res.json()) as T;
}

export const api = {
  listMailboxes: () => request<Mailbox[]>("/api/mailboxes"),

  createMailbox: (body: { server_hostport: string; username: string; password: string; primary_inbox: string }) =>
    request<Mailbox>("/api/mailboxes", {
      method: "POST",
      body: JSON.stringify(body),
    }),

  syncMailbox: (id: number) =>
    request<{ message: string }>(`/api/mailboxes/${id}/sync`, {
      method: "POST",
    }),

  listMessages: (id: number) => request<EmailSummary[]>(`/api/mailboxes/${id}/messages`),

  getMessage: (mailboxId: number, emailId: number) => request<Email>(`/api/mailboxes/${mailboxId}/messages/${emailId}`),

  listAllMessages: async (): Promise<InboxLoad> => {
    const mailboxes = await api.listMailboxes();
    const settled = await Promise.allSettled(mailboxes.map((m) => api.listMessages(m.id)));
    const messages: MessageWithMailbox[] = [];
    mailboxes.forEach((box, i) => {
      const res = settled[i];
      if (res.status === "fulfilled") {
        for (const msg of res.value) messages.push({ ...msg, mailbox: box });
      }
    });
    messages.sort((a, b) => new Date(b.received_at).getTime() - new Date(a.received_at).getTime());
    return { mailboxes, messages, failed: settled.length - messagesFailed(settled) };
  },
};

function messagesFailed(settled: PromiseSettledResult<EmailSummary[]>[]): number {
  return settled.filter((r) => r.status === "rejected").length;
}
