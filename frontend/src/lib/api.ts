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
};
