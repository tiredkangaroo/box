import type { Mail, Mailbox } from "$lib/types";

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

  createMailbox: (body: {
    server_hostport: string;
    username: string;
    password: string;
    primary_inbox: string;
  }) =>
    request<Mailbox>("/api/mailboxes", {
      method: "POST",
      body: JSON.stringify(body),
    }),

  deleteMailbox: (id: number) =>
    request<{ message: string }>(`/api/mailboxes/${id}`, {
      method: "DELETE",
    }),

  syncMailbox: (id: number) =>
    request<{ message: string }>(`/api/mailboxes/${id}/sync`, {
      method: "POST",
    }),

  listMail: (mailboxId: number) =>
    request<Mail[]>(`/api/mailboxes/${mailboxId}/messages`),

  getMail: (mailboxId: number, mailId: number | string) =>
    request<Mail>(`/api/mailboxes/${mailboxId}/messages/${mailId}`),
};
